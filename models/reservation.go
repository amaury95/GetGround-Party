package models

import (
	"fmt"

	"gorm.io/gorm"
)

/*
Reservation is the object mapping to the guest list record into the database

It is composed of the attibutes:
	 - Name: name of the guest
	 - AccompanyingGuests: number of persons that accompany the guest

It is related to the following models:
	 - Table (one-to-many)
*/
type Reservation struct {
	Name               string `gorm:"primarykey" json:"name"`
	AccompanyingGuests int    `json:"accompanying_guests"`

	TableID int `json:"table"`
}

// Guests amount of accompanying people including the guest
func (r *Reservation) Guests() int { return 1 + r.AccompanyingGuests }

// Validate guest reservation fields.
func (r *Reservation) Validate(db *gorm.DB) error {
	if len(r.Name) < 6 {
		return fmt.Errorf("name should have at least 6 chatacters length")
	}

	return nil
}

func (r *Reservation) BeforeCreate(db *gorm.DB) error {
	if err := r.Validate(db); err != nil {
		return fmt.Errorf("error creating the guest reservation: %v", err)
	}

	// check table exists
	var table Table
	if err := db.Find(&table, r.TableID).Error; err != nil {
		return fmt.Errorf(`error loading table with id "%d": %v`, r.TableID, err)
	}

	// validate table capacity
	var reservations []Reservation

	// load table reservations
	if err := db.Model(&table).Association("Reservations").Find(&reservations); err != nil {
		return fmt.Errorf("error loading table reservations: %v", err)
	}

	// calculate sum of table reservations
	var sum int
	for _, r := range reservations {
		sum += r.Guests()
	}

	if sum+r.Guests() > table.Capacity {
		return fmt.Errorf("table capacity is exceded by: %d", sum+r.Guests()-table.Capacity)
	}

	return nil
}
