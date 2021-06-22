package models

import (
	"fmt"

	"gorm.io/gorm"
)

/*
Table is the object mapping to the table record into the database

It is composed of the attibutes:
	 - Capacity: capacity of guests

It is related to the following models:
	 - Reservations (many-to-one)
	 - Guests       (many-to-one)
*/
type Table struct {
	ID       int `gorm:"primarykey" json:"id"`
	Capacity int `json:"capacity"`

	Reservations []Reservation `json:"reservations,omitempty"`
	Guests       []Guest       `json:"guests,omitempty"`
}

// Validate table fields.
func (t *Table) Validate(db *gorm.DB) error {
	if t.Capacity <= 0 {
		return fmt.Errorf(`capacity "%d" is not valid`, t.Capacity)
	}

	return nil
}

func (t *Table) BeforeCreate(tx *gorm.DB) error {
	if err := t.Validate(tx); err != nil {
		return fmt.Errorf("error creating the table: %v", err)
	}

	return nil
}
