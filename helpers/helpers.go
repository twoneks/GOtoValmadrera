package helpers

import (
	"fmt"
	"os"
	"time"

	"github.com/tkanos/gonfig"
	"github.com/twoneks/gotovalma/config"
)

func TodayAt(hour int, min int) string {
	dateLayout := "2006-01-02 04:05"
	t := time.Now()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, hour, min, 0, t.Location()).Format(dateLayout)
}

func LoadConfig(configuration *config.Configuration) {
	if os.Getenv("GoEnv") == "" {
		configuration.GoEnv = "development"
	} else {
		configuration.GoEnv = os.Getenv("GoEnv")
	}
	err := gonfig.GetConf(fmt.Sprintf("%s/src/github.com/twoneks/gotovalma/config/%s.json", os.Getenv("GOPATH"), configuration.GoEnv), configuration)
	if err != nil {
		panic(err)
	}
}
