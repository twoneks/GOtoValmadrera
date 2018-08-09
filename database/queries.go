package database

import "fmt"

func WindDetectionInsert(detection string) string {
	return fmt.Sprintf("INSERT INTO detections VALUES (DEFAULT, %s, null)", detection)
}

func SelectDailyDetection(interval int, from string) string {
	return fmt.Sprintf("SELECT * FROM detections WHERE time BETWEEN '%s'::timestamp - (interval '%vm') AND '%s'::timestamp", from, interval, from)
}

func InsertStatRecord(stats string) string {
	return fmt.Sprintf("INSERT INTO stats values(DEFAULT, '%s', now()::date);", stats)
}

func UpdateWindyStats(windy bool) string {
	return fmt.Sprintf("UPDATE stats SET windy = %t WHERE day = now()::date;", windy)
}
