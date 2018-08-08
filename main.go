package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/robfig/cron"
	"github.com/twoneks/gotovalma/database"
	"github.com/twoneks/gotovalma/detector"
	"github.com/twoneks/gotovalma/windStation"
)

type detection struct {
	id        int
	knots     int
	direction sql.NullString
	time      string
}

func main() {
	fmt.Println("Booting...")
	db := database.Connect()
	defer db.Close()

	// // TODO: Find a way to spot this after x iterations
	ticker := time.NewTicker(3 * 60000 * time.Millisecond)
	//
	go func(db *sql.DB, ticker *time.Ticker) {

		for {
			select {
			case i := <-ticker.C:
				fmt.Print(i)
				knots := windStation.Get()
				_, err := db.Exec(database.WindDetectionInsert(knots))
				if err != nil {
					panic(err)
				}
			case <-time.After(time.Hour * 2):
				ticker.Stop()
			}
		}
	}(db, ticker)

	c := cron.New()
	//Detect wind average at 5.30am
	c.AddFunc("0 30 5 * * *", func() { detector.Detect(db, todayAt(5, 30)) })
	//Update daily stat at 9.00am setting whether was windy or not
	c.AddFunc("0 0 9 * * *", func() { detector.UpdateDailyStat(db, todayAt(9, 00)) })

	c.Start()

	done := make(chan bool)
	<-done
}

func todayAt(hour int, min int) string {
	dateLayout := "2006-01-02 04:05"
	t := time.Now()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, hour, min, 0, t.Location()).Format(dateLayout)
}
