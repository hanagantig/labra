package entity

import "time"

type Markers []Marker

// Marker holds metadata about a marker type.
type Marker struct {
	ID             int
	Name           string
	ReferenceRange Range
	PrimaryColor   int
}

type Range struct {
	From string
	To   string
}

type BioMarker struct {
	ID   int
	Name string
}

func (b Markers) IDs() []int {
	ids := make([]int, 0, len(b))
	for _, b := range b {
		ids = append(ids, b.ID)
	}
	return ids
}

type MarkerResults []MarkerResult

// MarkerResult holds a single result for a marker, referencing the marker metadata.
type MarkerResult struct {
	ID               int
	CheckupID        int
	Marker           Marker // Embed the marker metadata
	UnrecognizedName string
	Value            float64
	Unit             Unit
	CreatedAt        time.Time
}

type MarkerFilter struct {
	From      time.Time
	To        time.Time
	Names     []string
	CheckupID int
}

// GroupMarkerResults groups results by marker ID, returning a map of marker ID to results.
func (m MarkerResults) Group() map[MarkerResult]MarkerResults {
	grouped := make(map[MarkerResult]MarkerResults)
	for _, res := range m {
		grouped[res] = append(grouped[res], res)
	}
	return grouped
}

// MaxValue returns the maximum value among the results.
func (m MarkerResults) MaxValue() float64 {
	if len(m) == 0 {
		return 0
	}
	maxVal := m[0].Value
	for _, res := range m {
		if res.Value > maxVal {
			maxVal = res.Value
		}
	}
	return maxVal
}

// MinValue returns the minimum value among the results.
func (m MarkerResults) MinValue() float64 {
	minVal := float64(0)
	for _, res := range m {
		if res.Value < minVal {
			minVal = res.Value
		}
	}

	return minVal
}
