package model

type Ride struct {
	Source      string // Source location
	Destination string // Destination location
	Distance    int32  // Distance in kilometers
	Cost        int32  // Cost in currency units
}
