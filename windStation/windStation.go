package windStation

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/twoneks/gotovalma/database"
)

// TakeOver the wind speed on the meteolecco web page and store the data in the database
func TakeOver(db *sql.DB) {
	knots := getActualWindSpeed()
	storeDetection(knots, db)
}

func getActualWindSpeed() int {
	response, err := http.Get("http://www.meteolecco.it/attuale.php")
	if err != nil || response.StatusCode != 200 {
		return -1
	}
	defer response.Body.Close()

	// Get the response body as a string
	dataInBytes, _ := ioutil.ReadAll(response.Body)
	pageContent := string(dataInBytes)

	// Find a substr
	speedIndex := strings.Index(pageContent, "minuto: media ")
	if speedIndex == -1 {
		return -1
	}
	kmhSpeed, _ := strconv.ParseFloat(pageContent[speedIndex+14:speedIndex+14+3], 64)
	knotsSpeed, _ := strconv.ParseFloat(fmt.Sprintf("%.0f", kmhSpeed*0.539957), 64)
	return int(knotsSpeed)
}

func storeDetection(knots int, db *sql.DB) error {
	_, err := db.Exec(database.WindDetectionInsert(strconv.Itoa(knots)))
	if err != nil {
		return err
	}
	return nil
}
