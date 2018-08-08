package detector

import (
	"database/sql"
	"encoding/json"
	"log"
	"sync"

	"github.com/twoneks/gotovalma/database"
)

type WindAverage struct {
	Interval int `json:"interval"`
	Average  int `json:"average"`
}

// Detect the average wind speed and store it on db
func Detect(db *sql.DB, today string) []WindAverage {
	var waitGroup sync.WaitGroup
	intervals := [3]int{120, 60, 30}
	averagesChan := make(chan WindAverage, len(intervals))

	for _, interval := range intervals {
		waitGroup.Add(1)
		go func(interval int, from string, averagesChan chan<- WindAverage, db *sql.DB, waitGroup *sync.WaitGroup) {
			defer waitGroup.Done()
			averagesChan <- calculateAverage(interval, from, db)
		}(interval, today, averagesChan, db, &waitGroup)
	}

	waitGroup.Wait()
	close(averagesChan)
	var averages []WindAverage
	for average := range averagesChan {
		averages = append(averages, average)
	}

	jsonAverages, _ := json.Marshal(averages)
	db.Exec(database.InsertStatRecord(string(jsonAverages)))
	if true {
		// Perform here the call
	}

	return averages
}

// UpdateDailyStat setting whether was windy or not
func UpdateDailyStat(db *sql.DB, today string) {
	// Calculate average wind for a session of 120 seconds back from today
	windy := calculateAverage(120, today, db)
	db.Exec(database.UpdateWindyStats(windy.Average > 14))
}

func calculateAverage(interval int, from string, db *sql.DB) WindAverage {
	detections, err := db.Query(database.SelectDailyDetection(interval, from))
	if err != nil {
		panic(err)
	}

	var iterationDetection database.Detection
	var windSpeeds []int
	detectionsCount := 0

	for detections.Next() {
		err := detections.Scan(&iterationDetection.ID, &iterationDetection.Knots, &iterationDetection.Direction, &iterationDetection.Time)
		if err != nil {
			log.Fatal(err)
		}
		windSpeeds = append(windSpeeds, iterationDetection.Knots)
		detectionsCount++
	}

	sum := 0
	for val := range windSpeeds {
		sum += val
	}

	if detectionsCount == 0 {
		return WindAverage{Interval: interval, Average: sum / 1}
	}
	return WindAverage{Interval: interval, Average: sum / detectionsCount}
}
