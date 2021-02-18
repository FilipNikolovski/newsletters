package storage

import "github.com/mailbadger/app/entities"

// CreateUnsubscribeEvent creates a record for unsubscribed subscriber in the database.
func (db *store) CreateUnsubscribeEvent(us *entities.UnsubscribeEvents) error {
	return db.Create(us).Error
}