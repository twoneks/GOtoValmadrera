package windStation

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Get and parse the webpage to return a detection
func Get() string {
	response, err := http.Get("http://www.meteolecco.it/attuale.php")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Get the response body as a string
	dataInBytes, err := ioutil.ReadAll(response.Body)
	pageContent := string(dataInBytes)

	// Find a substr
	speedIndex := strings.Index(pageContent, "minuto: media ") + 14
	if speedIndex == -1 {
		fmt.Println("No title element found")
	}
	kmhSpeed, _ := strconv.ParseFloat(pageContent[speedIndex:speedIndex+3], 64)
	return fmt.Sprintf("%.0f", kmhSpeed*0.539957)
}
