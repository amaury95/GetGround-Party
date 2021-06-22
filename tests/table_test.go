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
			Logger: logger.Default.LogMode(logger.Error),
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

	It("succeed creating a valid table", func() {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `tables` (`capacity`) VALUES (?)")).WithArgs(4).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		client.POST(`/tables`).WithJSON(api.CreateTableRequest{Capacity: 4}).
			Expect().Status(http.StatusCreated).
			JSON().Equal(models.Table{ID: 1, Capacity: 4})
	})

	It("fails creating a table with 0 capacity", func() {
		client.POST(`/tables`).WithJSON(api.CreateTableRequest{Capacity: 0}).
			Expect().Status(http.StatusBadRequest)
	})

	It("fails creating a table with negative capacity", func() {
		client.POST(`/tables`).WithJSON(api.CreateTableRequest{Capacity: -1}).
			Expect().Status(http.StatusBadRequest)
	})

	It("retieves populated tables list", func() {
		rows := sqlmock.NewRows([]string{"id", "capacity"}).
			AddRow(1, 5).
			AddRow(2, 4)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `tables`")).
			WillReturnRows(rows)

		client.GET(`/tables`).
			Expect().Status(http.StatusOK).
			JSON().Array().
			Elements(
				models.Table{ID: 1, Capacity: 5},
				models.Table{ID: 2, Capacity: 4},
			)
	})

	It("retrieves the empty tables list", func() {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `tables`")).
			WillReturnRows(sqlmock.NewRows(nil))

		client.GET(`/tables`).
			Expect().Status(http.StatusOK).
			JSON().Array().Empty()
	})

	It("retrieves the empty seats", func() {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT SUM(capacity) FROM `tables`")).
			WillReturnRows(sqlmock.NewRows([]string{"row"}).AddRow(6))

		mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) + SUM(accompanying_guests) FROM `guests`")).
			WillReturnRows(sqlmock.NewRows([]string{"row"}).AddRow(2))

		client.GET(`/seats_empty`).
			Expect().Status(http.StatusOK).
			JSON().Object().Equal(api.GetSeatsEmptyRespose{SeatsEmpty: 4})
	})
})
