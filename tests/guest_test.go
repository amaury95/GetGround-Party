package tests_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"regexp"
	"time"

	"github.com/amaury95/GetGround-Party/api"
	"github.com/amaury95/GetGround-Party/models"
	"github.com/gavv/httpexpect"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/DATA-DOG/go-sqlmock"
)

var _ = Describe("Guest controller", func() {
	var (
		mock   sqlmock.Sqlmock
		server *httptest.Server
		client *httpexpect.Expect
	)

	BeforeEach(func() {
		var (
			db  *sql.DB
			err error
		)

		// get database mock
		db, mock, err = sqlmock.New()
		Expect(err).NotTo(HaveOccurred())

		// mock database connection
		gdb, err := gorm.Open(mysql.New(mysql.Config{
			Conn:                      db,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		Expect(err).NotTo(HaveOccurred())

		// create handler with mock connection
		handler := new(api.Handler).WithConnection(gdb)

		// setup test server
		server = httptest.NewServer(handler.Router(&api.RouterConfig{
			ReleaseMode: true,
		}))

		// setup http expect
		client = httpexpect.New(GinkgoT(), server.URL)
	})

	AfterEach(func() {
		// close server
		server.Close()

		// make sure all expectations were met
		err := mock.ExpectationsWereMet()
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("registers a guest in an empty table", func() {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `reservations` WHERE name = ? ORDER BY `reservations`.`name` LIMIT 1")).WithArgs("username").
			WillReturnRows(sqlmock.NewRows([]string{"name", "accompanying_guests", "table_id"}).AddRow("username", 5, 1))

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `tables` WHERE `tables`.`id` = ?")).WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "capacity"}).AddRow(1, 6))

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `guests` WHERE `guests`.`table_id` = ?")).WithArgs(1).
			WillReturnRows(sqlmock.NewRows(nil))

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `guests` (`name`,`accompanying_guests`,`table_id`,`created_at`) VALUES (?,?,?,?)")).
			WithArgs("username", 5, 1, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		client.PUT(`/guests/username`).WithJSON(api.CreateGuestRequest{AccompanyingGuests: 5}).
			Expect().Status(http.StatusCreated).
			JSON().Equal(api.CreateGuestResponse{Name: "username"})
	})

	It("registers a guest in a not empty table", func() {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `reservations` WHERE name = ? ORDER BY `reservations`.`name` LIMIT 1")).WithArgs("username").
			WillReturnRows(sqlmock.NewRows([]string{"name", "accompanying_guests", "table_id"}).AddRow("username", 5, 1))

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `tables` WHERE `tables`.`id` = ?")).WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "capacity"}).AddRow(1, 9))

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `guests` WHERE `guests`.`table_id` = ?")).WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"name", "accompanying_guests", "table_id"}).
				AddRow("lastname", 2, 1))

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `guests` (`name`,`accompanying_guests`,`table_id`,`created_at`) VALUES (?,?,?,?)")).
			WithArgs("username", 5, 1, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		client.PUT(`/guests/username`).WithJSON(api.CreateGuestRequest{AccompanyingGuests: 5}).
			Expect().Status(http.StatusCreated).
			JSON().Equal(api.CreateGuestResponse{Name: "username"})
	})

	It("fails registering a guest with a name shorter than 6 characters", func() {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `reservations` WHERE name = ? ORDER BY `reservations`.`name` LIMIT 1")).WithArgs("user").
			WillReturnRows(sqlmock.NewRows([]string{"name", "accompanying_guests", "table_id"}).AddRow("user", 5, 1))

		client.PUT(`/guests/user`).WithJSON(api.CreateGuestRequest{AccompanyingGuests: 5}).
			Expect().Status(http.StatusBadRequest)
	})

	It("fails registering a guest with negative accompanying", func() {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `reservations` WHERE name = ? ORDER BY `reservations`.`name` LIMIT 1")).WithArgs("username").
			WillReturnRows(sqlmock.NewRows([]string{"name", "accompanying_guests", "table_id"}).AddRow("username", 5, 1))

		client.PUT(`/guests/username`).WithJSON(api.CreateGuestRequest{AccompanyingGuests: -5}).
			Expect().Status(http.StatusBadRequest)
	})

	It("fails registering a guest for an accompanying bigger than total capacity", func() {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `reservations` WHERE name = ? ORDER BY `reservations`.`name` LIMIT 1")).WithArgs("username").
			WillReturnRows(sqlmock.NewRows([]string{"name", "accompanying_guests", "table_id"}).AddRow("username", 5, 1))

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `tables` WHERE `tables`.`id` = ?")).WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "capacity"}).AddRow(1, 5))

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `guests` WHERE `guests`.`table_id` = ?")).WithArgs(1).
			WillReturnRows(sqlmock.NewRows(nil))

		mock.ExpectRollback()

		client.PUT(`/guests/username`).WithJSON(api.CreateGuestRequest{AccompanyingGuests: 5}).
			Expect().Status(http.StatusInternalServerError)
	})

	It("fails registering a guest for an accompanying bigger than available capacity", func() {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `reservations` WHERE name = ? ORDER BY `reservations`.`name` LIMIT 1")).WithArgs("username").
			WillReturnRows(sqlmock.NewRows([]string{"name", "accompanying_guests", "table_id"}).AddRow("username", 5, 1))

		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `tables` WHERE `tables`.`id` = ?")).WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "capacity"}).AddRow(1, 8))

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `guests` WHERE `guests`.`table_id` = ?")).WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"name", "accompanying_guests", "table_id"}).
				AddRow("lastname", 2, 1))

		mock.ExpectRollback()

		client.PUT(`/guests/username`).WithJSON(api.CreateGuestRequest{AccompanyingGuests: 5}).
			Expect().Status(http.StatusInternalServerError)
	})

	It("deletes a guest from the registry", func() {
		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `guests` WHERE name = ?")).WithArgs("username").
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		client.DELETE(`/guests/username`).Expect().Status(http.StatusAccepted)
	})

	It("retrieves the populated guests list", func() {
		date := time.Now()

		rows := sqlmock.
			NewRows([]string{"name", "accompanying_guests", "table_id", "created_at"}).
			AddRow("user01", 3, 1, date).
			AddRow("user02", 4, 2, date)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `guests`")).
			WillReturnRows(rows)

		resp := api.GetGuestsResponse{
			Guests: []models.Guest{
				{Name: "user01", AccompanyingGuests: 3, TableID: 1, CreatedAt: date},
				{Name: "user02", AccompanyingGuests: 4, TableID: 2, CreatedAt: date},
			},
		}

		client.GET(`/guests`).
			Expect().Status(http.StatusOK).
			JSON().Equal(resp)
	})

	It("retrieves an empty guests list", func() {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `guests`")).
			WillReturnRows(sqlmock.NewRows(nil))

		resp := api.GetReservationsResponse{
			Guests: []models.Reservation{},
		}

		client.GET(`/guests`).
			Expect().Status(http.StatusOK).
			JSON().Equal(resp)
	})
})
