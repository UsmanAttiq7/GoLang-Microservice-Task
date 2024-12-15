package model

// Ride represents a ride entity.
type Ride struct {
	ID          int32  // Ride ID
	Source      string // Source location
	Destination string // Destination location
	Distance    int32  // Distance in kilometers
	Cost        int32  // Cost in currency units
}
