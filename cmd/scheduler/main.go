package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/mailbadger/app/emails"
	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/mode"
	awssqs "github.com/mailbadger/app/sqs"
	"github.com/mailbadger/app/storage"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	mode.SetModeFromEnv()

	lvl, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		lvl = logrus.InfoLevel
	}

	logrus.SetLevel(lvl)
	if mode.IsProd() {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
	logrus.SetOutput(os.Stdout)

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		logrus.WithError(err).Fatal("AWS configuration error")
	}

	client := sqs.NewFromConfig(cfg)

	queueName := awssqs.CampaignerTopic
	gQInput := &sqs.GetQueueUrlInput{
		QueueName: &queueName,
	}
	// Get URL of queue
	urlResult, err := client.GetQueueUrl(ctx, gQInput)
	if err != nil {
		logrus.WithError(err).Fatal("Got an error getting the queue URL")
	}

	queueURL := urlResult.QueueUrl
	p := awssqs.NewPublisher(
		client,
	)

	driver := os.Getenv("DATABASE_DRIVER")
	conf := storage.MakeConfigFromEnv(driver)
	s := storage.New(driver, conf)

	now := time.Now()
	err = job(ctx, s, p, queueURL, now)
	if err != nil {
		logrus.WithField("time", now).WithError(err).Error("failed to start campaign scheduler job")
	}
	end := time.Since(now)

	logrus.Infof("Scheduler started at %v and took %v to finish", now, end)

}

func job(ctx context.Context, s storage.Storage, p awssqs.Publisher, queueURL *string, time time.Time) error {
	scheduledCampaigns, err := s.GetScheduledCampaigns(time)
	if err != nil {
		return fmt.Errorf("failed to get scheduled campaigns: %w", err)
	}

	for _, cs := range scheduledCampaigns {

		logEntry := logrus.WithFields(logrus.Fields{
			"campaign_id": cs.CampaignID,
			"user_id":     cs.UserID,
		})

		u, err := s.GetUser(cs.UserID)
		if err != nil {
			logEntry.WithError(err).Error("failed to get user.")
			continue
		}
		campaign, err := s.GetCampaign(cs.CampaignID, u.ID)
		if err != nil {
			logEntry.WithError(err).Error("failed to get campaign.")
			continue
		}
		if campaign.Status != entities.StatusScheduled {
			logEntry.WithError(err).Warn("campaign status is not 'scheduled'")
			continue
		}

		template, err := s.GetTemplate(campaign.BaseTemplate.ID, u.ID)
		if err != nil {
			logEntry.WithField("template_id", campaign.BaseTemplate.ID).WithError(err).Error("failed to get template.")
			continue
		}
		templateData, err := cs.GetMetadata()
		if err != nil {
			logEntry.WithError(err).Error("failed to unmarshal default template data.")
			continue
		}
		err = template.ValidateData(templateData)
		if err != nil {
			logEntry.WithError(err).Error("failed to validate template data.")
			continue
		}

		sesKeys, err := s.GetSesKeys(u.ID)
		if err != nil {
			logEntry.WithError(err).Error("failed to get ses keys.")
			continue
		}

		segmentIDs, err := cs.GetSegmentIDs()
		if err != nil {
			logEntry.WithError(err).Error("failed to unmarshal segment ids.")
			continue
		}

		lists, err := s.GetSegmentsByIDs(u.ID, segmentIDs)
		if err != nil || len(lists) == 0 {
			logEntry.WithField("segment_ids", segmentIDs).WithError(err).Error("failed to get segments by ids.")
			continue
		}

		sender, err := emails.NewSesSender(sesKeys.AccessKey, sesKeys.SecretKey, sesKeys.Region)
		if err != nil {
			logEntry.WithError(err).Error("failed to create new ses sender.")
			continue
		}

		_, err = sender.DescribeConfigurationSet(&ses.DescribeConfigurationSetInput{
			ConfigurationSetName: aws.String(emails.ConfigurationSetName),
		})

		params := &entities.CampaignerTopicParams{
			EventID:                cs.ID,
			CampaignID:             cs.CampaignID,
			SegmentIDs:             segmentIDs,
			TemplateData:           templateData,
			Source:                 fmt.Sprintf("%s <%s>", cs.FromName, cs.Source),
			UserID:                 u.ID,
			UserUUID:               u.UUID,
			ConfigurationSetExists: err == nil,
			SesKeys:                *sesKeys,
		}
		paramsByte, err := json.Marshal(params)
		if err != nil {
			logEntry.WithError(err).Error("failed to marshal params for campaigner.")
			continue
		}
		err = p.SendMessage(ctx, queueURL, paramsByte)
		if err != nil {
			logEntry.WithError(err).Error("failed to publish campaign to campaigner.")
			continue
		}
		campaign.Status = entities.StatusSending
		err = s.UpdateCampaign(campaign)
		if err != nil {
			logEntry.WithError(err).Error("failed to update status of campaign.")
			continue
		}
	}

	return nil

}
