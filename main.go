package main

import (
	"database/sql"
	"time"

	"github.com/robfig/cron"
	"github.com/twoneks/gotovalma/config"
	"github.com/twoneks/gotovalma/database"
	"github.com/twoneks/gotovalma/detector"
	"github.com/twoneks/gotovalma/helpers"
	"github.com/twoneks/gotovalma/windStation"
)

type detection struct {
	id        int
	knots     int
	direction sql.NullString
	time      string
}

func main() {
	configuration := config.Configuration{}
	helpers.LoadConfig(&configuration)

	db := database.Connect()
	defer db.Close()

	cronTab := cron.New()

	// Monitor the actual wind condition for number ho hours set in MonitoringInterval
	cronTab.AddFunc("0 30 4 * * *", func() { startMonitioring(&configuration, db) })
	//Detect wind average at 5.30am
	cronTab.AddFunc("0 30 5 * * *", func() { detector.Detect(db, helpers.TodayAt(5, 30), &configuration) })
	//Update daily stat at 9.00am setting whether was windy or not
	cronTab.AddFunc("0 0 9 * * *", func() { detector.UpdateDailyStat(db, helpers.TodayAt(9, 00)) })

	cronTab.Start()

	done := make(chan bool)
	<-done
}

func startMonitioring(configuration *config.Configuration, db *sql.DB) {
	ticker := time.NewTicker(time.Duration(configuration.MonitoringPollingInterval) * time.Minute)
	timeout := make(chan bool, 1)

	go func(timeout chan bool) {
		time.Sleep(time.Duration(configuration.MonitoringInterval) * time.Hour)
		timeout <- true
	}(timeout)

	go func(db *sql.DB, ticker *time.Ticker) {
		for {
			select {
			case <-ticker.C:
				windStation.TakeOver(db)
			case <-timeout:
				ticker.Stop()
			}
		}
	}(db, ticker)
}
