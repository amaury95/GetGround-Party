package tests_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"regexp"

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

var _ = Describe("Table controller", func() {
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

	It("creates a reservation in an empty table", func() {
		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `tables` WHERE `tables`.`id` = ?")).WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "capacity"}).AddRow(1, 6))

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `reservations` WHERE `reservations`.`table_id` = ?")).WithArgs(1).
			WillReturnRows(sqlmock.NewRows(nil))

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `reservations` (`name`,`accompanying_guests`,`table_id`) VALUES (?,?,?)")).
			WithArgs("username", 5, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		client.POST(`/guest_list/username`).WithJSON(api.CreateReservationRequest{Table: 1, Guests: 5}).
			Expect().Status(http.StatusCreated).
			JSON().Equal(api.CreateReservationResponse{Name: "username"})
	})

	It("creates a reservation in a not empty table", func() {
		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `tables` WHERE `tables`.`id` = ?")).WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "capacity"}).AddRow(1, 8))

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `reservations` WHERE `reservations`.`table_id` = ?")).WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"name", "accompanying_guests", "table_id"}).
				AddRow("lastname", 1, 1))

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `reservations` (`name`,`accompanying_guests`,`table_id`) VALUES (?,?,?)")).
			WithArgs("username", 5, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		client.POST(`/guest_list/username`).WithJSON(api.CreateReservationRequest{Table: 1, Guests: 5}).
			Expect().Status(http.StatusCreated).
			JSON().Equal(api.CreateReservationResponse{Name: "username"})
	})

	It("fails creating a reservation with a name shorter than 6 characters", func() {
		client.POST(`/guest_list/user`).WithJSON(api.CreateReservationRequest{Table: 1, Guests: 5}).
			Expect().Status(http.StatusBadRequest)
	})

	It("fails creating a reservation with negative accompanying", func() {
		client.POST(`/guest_list/username`).WithJSON(api.CreateReservationRequest{Table: 1, Guests: -5}).
			Expect().Status(http.StatusBadRequest)
	})

	It("fails creating a reservation for an accompanying bigger than total capacity", func() {
		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `tables` WHERE `tables`.`id` = ?")).WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "capacity"}).AddRow(1, 4))

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `reservations` WHERE `reservations`.`table_id` = ?")).WithArgs(1).
			WillReturnRows(sqlmock.NewRows(nil))

		mock.ExpectRollback()

		client.POST(`/guest_list/username`).WithJSON(api.CreateReservationRequest{Table: 1, Guests: 5}).
			Expect().Status(http.StatusInternalServerError)
	})

	It("fails creating a reservation for an accompanying bigger than available capacity", func() {
		mock.ExpectBegin()

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `tables` WHERE `tables`.`id` = ?")).WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "capacity"}).AddRow(1, 8))

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `reservations` WHERE `reservations`.`table_id` = ?")).WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"name", "accompanying_guests", "table_id"}).
				AddRow("lastname", 2, 1))

		mock.ExpectRollback()

		client.POST(`/guest_list/username`).WithJSON(api.CreateReservationRequest{Table: 1, Guests: 5}).
			Expect().Status(http.StatusInternalServerError)
	})

	It("retrieves the populated reservations list", func() {
		rows := sqlmock.
			NewRows([]string{"name", "accompanying_guests", "table_id"}).
			AddRow("user01", 3, 1).
			AddRow("user02", 4, 2)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `reservations`")).
			WillReturnRows(rows)

		resp := api.GetReservationsResponse{
			Guests: []models.Reservation{
				{Name: "user01", AccompanyingGuests: 3, TableID: 1},
				{Name: "user02", AccompanyingGuests: 4, TableID: 2},
			},
		}

		client.GET(`/guest_list`).
			Expect().Status(http.StatusOK).
			JSON().Equal(resp)
	})

	It("retrieves the empty reservations list", func() {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `reservations`")).
			WillReturnRows(sqlmock.NewRows(nil))

		resp := api.GetReservationsResponse{
			Guests: []models.Reservation{},
		}

		client.GET(`/guest_list`).
			Expect().Status(http.StatusOK).
			JSON().Equal(resp)
	})
})
