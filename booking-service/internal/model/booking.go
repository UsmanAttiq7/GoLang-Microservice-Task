package model

import "time"

type Booking struct {
	ID        int32
	UserID    int32
	RideID    int32
	Timestamp time.Time
}
