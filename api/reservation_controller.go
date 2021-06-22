package api

import (
	"encoding/json"
	"net/http"

	"github.com/amaury95/GetGround-Party/models"
	"github.com/gin-gonic/gin"
)

type CreateReservationRequest struct {
	Table  int `json:"table"`
	Guests int `json:"accompanying_guests"`
}

type CreateReservationResponse struct {
	Name string `json:"name"`
}

func (h *Handler) CreateReservation(g *gin.Context) {
	var body CreateReservationRequest

	// decode body from request
	if err := json.NewDecoder(g.Request.Body).Decode(&body); err != nil {
		g.String(http.StatusInternalServerError, "error decoding body: %v", err)
		return
	}

	// decode name from params
	name := g.Param("name")

	record := models.Reservation{
		Name:               name,
		AccompanyingGuests: body.Guests,
		TableID:            body.Table,
	}

	// validate model
	if err := record.Validate(h.db); err != nil {
		g.String(http.StatusBadRequest, "error validating reservation: %v", err)
		return
	}

	// create model in the database
	if err := h.db.Create(&record).Error; err != nil {
		g.String(http.StatusInternalServerError, "error creating reservation: %v", err)
		return
	}

	g.JSON(http.StatusCreated, CreateReservationResponse{Name: record.Name})
}

type GetReservationsResponse struct {
	Elements []models.Reservation `json:"guests"`
}

func (h *Handler) GetReservations(g *gin.Context) {
	var elements []models.Reservation

	if err := h.db.Find(&elements).Error; err != nil {
		g.String(http.StatusInternalServerError, "error retrieving reservations: %v", err)
		return
	}

	g.JSON(http.StatusOK, GetReservationsResponse{Elements: elements})
}
