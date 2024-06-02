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

type Point struct {
    X float64
    Y float64
}

// Convert latitude and longitude to Cartesian coordinates
func LatLongToCartesian(lat, long float64) Point {
    // Convert latitude and longitude to Cartesian coordinates using Haversine formula
    // Assume the starting point as (0, 0) in Cartesian coordinates
    lat1, lon1 := 0.0, 0.0
    lat2, lon2 := lat, long

    // Calculate distance between two points using Haversine formula
    distance := Haversine(lat1, lon1, lat2, lon2)

    // Convert distance to Cartesian coordinates
    // For simplicity, assuming Earth's circumference at the equator forms a circle
    // The circumference of the Earth is approximately 40075 km
    // Using a scaling factor to map the distance to Cartesian coordinates
    scaleFactor := 40075.0 / (2 * math.Pi) // Circumference / (2 * Pi)
    x := distance * scaleFactor
    y := 0.0 // Assuming latitude variation only affects x-coordinate

    return Point{X: x, Y: y}
}

func CalculateRectangleArea(p1, p2 Point) float64 {
    // Calculate the side lengths of the rectangle
    width := math.Abs(p1.X - p2.X)
    height := math.Abs(p1.Y - p2.Y)

    // Calculate the area of the rectangle
    area := width * height

    return area
}