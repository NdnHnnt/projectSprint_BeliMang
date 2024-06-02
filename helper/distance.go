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

func LatLongToCartesian(lat, lon float64) (float64, float64, float64) {
    // Convert latitude and longitude from degrees to radians
    lat = lat * math.Pi / 180
    lon = lon * math.Pi / 180

    // Assume the Earth as a sphere with radius 6371 km
    radius := 6371.0

    // Convert latitude and longitude to Cartesian coordinates
    x := radius * math.Cos(lat) * math.Cos(lon)
    y := radius * math.Cos(lat) * math.Sin(lon)
    z := radius * math.Sin(lat)

    return x, y, z
}

func CalculateRectangleArea(x1 float64, y1 float64, x2 float64, y2 float64) bool {
	// Calculate the side lengths of the rectangle
	width := math.Abs(x1 - x2)
	height := math.Abs(y1 - y2)

	// Calculate the area of the rectangle
	area := width * height

	return area <= 3.0
}
