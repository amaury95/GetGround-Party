/*
Package api holds the structure for the api handler that is going to be used to run and test the application.
*/
package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler is the structure that holds the connection context for the application
type Handler struct {
	db *gorm.DB
}

// Connection is the connection context getter
func (h *Handler) Connection() *gorm.DB { return h.db }

// WithConnection sets the given connection as context of the handler and return it
func (h *Handler) WithConnection(conn *gorm.DB) *Handler {
	h.db = conn
	return h
}

// Router returns the router engine for the api handler
func (h *Handler) Router() *gin.Engine {
	r := gin.Default()

	// tables
	r.GET(`/tables`, h.GetTables)
	r.POST(`/tables`, h.CreateTable)

	// reservations
	r.POST(`/guest_list/:name`, h.CreateReservation)
	r.GET(`/guest_list`, h.GetReservations)

	// guests
	r.PUT(`/guests/:name`, h.CreateGuest)
	r.GET(`/guests`, h.GetGuests)
	r.DELETE(`/guests/:name`, h.DeleteGuest)

	return r
}
