package api

import (
	"encoding/json"
	"net/http"

	"github.com/amaury95/GetGround-Party/models"
	"github.com/gin-gonic/gin"
)

/*
	Create Guest
*/

type CreateGuestRequest struct {
	AccompanyingGuests int `json:"accompanying_guests"`
}

type CreateGuestResponse struct {
	Name string `json:"name"`
}

func (h *Handler) CreateGuest(g *gin.Context) {
	var body CreateGuestRequest

	// decode body from request
	if err := json.NewDecoder(g.Request.Body).Decode(&body); err != nil {
		g.String(http.StatusInternalServerError, "error decoding body: %v", err)
		return
	}

	// decode name from params
	name := g.Param("name")

	// get guest reservation
	var reservation models.Reservation
	if err := h.db.First(&reservation, "name = ?", name).Error; err != nil {
		g.String(http.StatusInternalServerError, "error retrieving guest reservation: %v", err)
		return
	}

	record := models.Guest{
		Name:               name,
		AccompanyingGuests: body.AccompanyingGuests,
		TableID:            reservation.TableID,
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

	g.JSON(http.StatusCreated, CreateGuestResponse{Name: record.Name})
}

/*
	Get Guests
*/

type GetGuestsResponse struct {
	Guests []models.Guest `json:"guests"`
}

func (h *Handler) GetGuests(g *gin.Context) {
	var elements []models.Guest

	if err := h.db.Find(&elements).Error; err != nil {
		g.String(http.StatusInternalServerError, "error retrieving the guests: %v", err)
		return
	}

	g.JSON(http.StatusOK, GetGuestsResponse{Guests: elements})
}

/*
	Delete Guest
*/

func (h *Handler) DeleteGuest(g *gin.Context) {
	// decode name from params
	name := g.Param("name")

	if err := h.db.Delete(new(models.Guest), "name = ?", name).Error; err != nil {
		g.String(http.StatusInternalServerError, "error deleting guest: %v", err)
		return
	}

	g.Status(http.StatusAccepted)
}
