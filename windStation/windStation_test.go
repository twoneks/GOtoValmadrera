package windStation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func TestWindStationWithInvalidBody(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://www.meteolecco.it/attuale.php",
		httpmock.NewStringResponder(200, `<body></body>`))

	assert.Equal(t, -1, getActualWindSpeed(), "Return -1 if wind speed value not present.")
}

func TestWindStationWithBroken(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://www.meteolecco.it/attuale.php",
		httpmock.NewStringResponder(500, `<tr bgcolor="#E2E2E2">                                      <td colspan="2"><div align="left">Ultimo                                          minuto: media 4.2                                          Km/h - max 6.4 Km/h</div></td>                                    </tr>`))

	assert.Equal(t, -1, getActualWindSpeed(), "Return -1 if wind speed value not present.")
}

func TestWindStationWithValidBody(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://www.meteolecco.it/attuale.php",
		httpmock.NewStringResponder(200, `<tr bgcolor="#E2E2E2">                                      <td colspan="2"><div align="left">Ultimo                                          minuto: media 4.2                                          Km/h - max 6.4 Km/h</div></td>                                    </tr>`))
	assert.Equal(t, 2, getActualWindSpeed(), "Return the correct wind speed.")
}

func TestShouldUpdateStats(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectExec("INSERT INTO detections").WillReturnResult(sqlmock.NewResult(1, 1))

	if err = storeDetection(4, db); err != nil {
		t.Errorf("%s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
