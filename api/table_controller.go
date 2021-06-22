package api

import (
	"encoding/json"
	"net/http"

	"github.com/amaury95/GetGround-Party/models"
	"github.com/gin-gonic/gin"
)

/*
	Create Table
*/

type CreateTableRequest struct {
	Capacity int `json:"capacity"`
}

// CreateTable creates a table with the given capacity
func (h *Handler) CreateTable(g *gin.Context) {
	var body CreateTableRequest

	// decode input from request
	if err := json.NewDecoder(g.Request.Body).Decode(&body); err != nil {
		g.String(http.StatusInternalServerError, "error decoding body: %v", err)
		return
	}

	record := models.Table{
		Capacity: body.Capacity,
	}

	// validate model
	if err := record.Validate(h.db); err != nil {
		g.String(http.StatusBadRequest, "error validating table: %v", err)
		return
	}

	// create model in the database
	if err := h.db.Create(&record).Error; err != nil {
		g.String(http.StatusInternalServerError, "error creating table: %v", err)
		return
	}

	g.JSON(http.StatusCreated, record)
}

/*
	Get Tables
*/

// GetTables returns a list of the existing tables on the database.
// TODO: pagination
func (h *Handler) GetTables(g *gin.Context) {
	var tables []models.Table

	if err := h.db.Find(&tables).Error; err != nil {
		g.String(http.StatusInternalServerError, "error retrieving tables: %v", err)
		return
	}

	g.JSON(http.StatusOK, tables)
}

/*
	Get Seats Empty
*/
type GetSeatsEmptyRespose struct {
	SeatsEmpty int `json:"seats_empty"`
}

// GetSeatsEmpty calculate the total availability of the party
func (h *Handler) GetSeatsEmpty(g *gin.Context) {
	// SET @availability := (SELECT SUM(capacity) FROM tables);
	// SELECT @availability - COUNT(*) - SUM(accompanying_guests) FROM guests;

	var capacity, occupied int

	if err := h.db.Select("SUM(capacity)").Table("tables").Scan(&capacity).Error; err != nil {
		g.String(http.StatusInternalServerError, "error calculating capacity: %v", err)
		return
	}

	if err := h.db.Select("COUNT(*) + SUM(accompanying_guests)").Table("guests").Scan(&occupied).Error; err != nil {
		g.String(http.StatusInternalServerError, "error getting occupancy: %v", err)
		return
	}

	g.JSON(http.StatusOK, GetSeatsEmptyRespose{SeatsEmpty: capacity - occupied})
}
