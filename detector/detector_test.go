package detector

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/twoneks/gotovalma/config"
	"github.com/twoneks/gotovalma/helpers"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestCalculateAverage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "knots", "direction", "time"}).
		AddRow(1, 2, "n", time.Now()).
		AddRow(1, 4, "n", time.Now()).
		AddRow(1, 6, "n", time.Now())

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	expected := WindAverage{Interval: 120, Average: 4}
	got := calculateAverage(120, "2018-08-1 10:20", db)
	assert.Equal(t, expected, got, "Return the correct wind speed average.")

	rows = sqlmock.NewRows([]string{"id", "knots", "direction", "time"})
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	expected = WindAverage{Interval: 120, Average: 0}
	got = calculateAverage(120, "2018-08-1 10:20", db)
	assert.Equal(t, expected, got, "Return the correct wind speed average.")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDetect(t *testing.T) {
	os.Setenv("GoEnv", "test")
	configuration := config.Configuration{}
	helpers.LoadConfig(&configuration)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "knots", "direction", "time"}).
		AddRow(1, 2, "n", time.Now())
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"id", "knots", "direction", "time"}).
		AddRow(2, 4, "n", time.Now())
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	got := Detect(db, helpers.TodayAt(5, 30), &configuration)
	expected := make([]WindAverage, 2)
	expected[0] = WindAverage{Interval: 120, Average: 4}
	expected[1] = WindAverage{Interval: 30, Average: 2}

	assert.Contains(t, expected, got[0])
	assert.Contains(t, expected, got[1])

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
