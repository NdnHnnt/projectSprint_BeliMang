package helpers

import (
	"math"
)

func Haversine(lat1 float64, lon1 float64, lat2 float64, lon2 float64) float64 {
	Radius := 6371
	lat1 = lat1 * math.Pi / 180
	lon1 = lon1 * math.Pi / 180
	lat2 = lat2 * math.Pi / 180
	lon2 = lon2 * math.Pi / 180
	latDiff := lat2 - lat1
	lonDiff := lon2 - lon1
	a := math.Sin(latDiff/2)*math.Sin(latDiff/2) +
		math.Cos(lat1)*math.Cos(lat2)*math.Sin(lonDiff/2)*math.Sin(lonDiff/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	km := float64(Radius) * c

	return km
}
