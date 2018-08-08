package database

import (
	"database/sql"
	"time"
)

// Detection map the wind table
type Detection struct {
	ID        int            `db:"id" json:"id"`
	Knots     int            `db:"knots" json:"knots"`
	Direction sql.NullString `db:"direction" json:"direction"`
	Time      time.Time      `db:"time" json:"time"`
}

type Stat struct {
	ID       int         `db:"id"`
	Averages []Detection `db:"averages"`
	Windy    bool        `db:"windi"`
}
