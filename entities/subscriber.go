package entities

import (
	"errors"
	"time"

	valid "github.com/asaskevich/govalidator"
)

var ErrSubscriberNameEmpty = errors.New("The subscriber name cannot be empty.")
var ErrEmailInvalid = errors.New("The specified email is not valid.")

//Subscriber represents the subscriber entity
type Subscriber struct {
	Id           int64                `json:"id" gorm:"column:id; primary_key:yes"`
	UserId       int64                `json:"-" gorm:"column:user_id; index"`
	Name         string               `json:"name" gorm:"not null"`
	Email        string               `json:"email" gorm:"not null"`
	Lists        []List               `json:"-" gorm:"many2many:subscribers_lists;"`
	Metadata     []SubscriberMetadata `json:"metadata" gorm:"ForeignKey:SubscriberId"`
	Blacklisted  bool                 `json:"blacklisted"`
	Active       bool                 `json:"active"`
	CreatedAt    time.Time            `json:"created_at"`
	UpdatedAt    time.Time            `json:"updated_at"`
	Errors       map[string]string    `json:"-" sql:"-"`
	TemplateData map[string]string    `json:"-" sql:"-"`
}

//SubscriberMetadata represents the subscriber metadata in a form of a key and value
type SubscriberMetadata struct {
	Id           int64  `gorm:"column:id; primary_key:yes"`
	SubscriberId int64  `gorm:"column:subscriber_id; index"`
	Key          string `gorm:"not null" valid:"alphanum,required"`
	Value        string `gorm:"not null" valid:"alphanum,required"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (s *Subscriber) Normalize() {
	s.TemplateData = make(map[string]string)

	s.TemplateData["name"] = s.Name

	for _, m := range s.Metadata {
		s.TemplateData[m.Key] = m.Value
	}
}

// Validate subscriber properties,
func (s *Subscriber) Validate() bool {
	s.Errors = make(map[string]string)

	if valid.Trim(s.Name, "") == "" {
		s.Errors["name"] = ErrSubscriberNameEmpty.Error()
	}

	if !valid.IsEmail(s.Email) {
		s.Errors["email"] = ErrEmailInvalid.Error()
	}

	return len(s.Errors) == 0
}
