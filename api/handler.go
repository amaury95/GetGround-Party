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

type RouterConfig struct {
	ShowLogs    bool
	ReleaseMode bool
}

// Router returns the router engine for the api handler
func (h *Handler) Router(config *RouterConfig) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	// setup configurations
	if config.ShowLogs {
		r.Use(gin.Logger())
	}

	if config.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	// tables
	r.GET(`/tables`, h.GetTables)
	r.POST(`/tables`, h.CreateTable)
	r.GET(`/seats_empty`, h.GetSeatsEmpty)

	// reservations
	r.POST(`/guest_list/:name`, h.CreateReservation)
	r.GET(`/guest_list`, h.GetReservations)

	// guests
	r.PUT(`/guests/:name`, h.CreateGuest)
	r.GET(`/guests`, h.GetGuests)
	r.DELETE(`/guests/:name`, h.DeleteGuest)

	return r
}
