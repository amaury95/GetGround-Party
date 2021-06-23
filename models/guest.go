package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

/*
Guest is the object mapping to the guest record into the database

It is composed of the attibutes:
	 - Name: name of the guest
	 - AccompanyingGuests: number of persons that accompany the guest

It is related to the following models:
	 - Table (one-to-many)
*/
type Guest struct {
	Name               string `gorm:"primarykey"`
	AccompanyingGuests int    `json:"accompanying_guests"`

	TableID   int       `json:"table"`
	CreatedAt time.Time `json:"time_arrived"`
}

// TotalGuests amount of accompanying people including the guest
func (g *Guest) TotalGuests() int { return 1 + g.AccompanyingGuests }

// Validate guest reservation fields.
func (g *Guest) Validate(db *gorm.DB) error {
	if len(g.Name) < 6 {
		return fmt.Errorf("name should have at least 6 chatacters length")
	}

	if g.AccompanyingGuests < 0 {
		return fmt.Errorf(`invalid "%d" guests amount`, g.AccompanyingGuests)
	}

	return nil
}

func (g *Guest) BeforeCreate(db *gorm.DB) error {
	if err := g.Validate(db); err != nil {
		return fmt.Errorf("error creating the guest reservation: %v", err)
	}

	// check table exists
	var table Table
	if err := db.Find(&table, g.TableID).Error; err != nil {
		return fmt.Errorf(`error loading table with id "%d": %v`, g.TableID, err)
	}

	// validate table capacity
	var guests []Guest

	// load table reservations
	if err := db.Model(&table).Association("Guests").Find(&guests); err != nil {
		return fmt.Errorf("error loading table guests: %v", err)
	}

	// calculate sum of table reservations
	var sum int
	for _, r := range guests {
		sum += r.TotalGuests()
	}

	if sum+g.TotalGuests() > table.Capacity {
		return fmt.Errorf("table capacity is exceded by: %d", sum+g.TotalGuests()-table.Capacity)
	}

	return nil
}
