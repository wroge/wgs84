package wgs84

import (
	"math"
	"testing"
)

func TestTransformFromWGS84SelectsTransformationByLocation(t *testing.T) {
	// North Sea: not covered by the first ED50 transformation in the list.
	lon, lat := 1.71, 58.73

	to, err := Transform(4326, 23090)
	if err != nil {
		t.Fatal(err)
	}

	east, north, _, err := to(lon, lat, 0)
	if err != nil {
		t.Fatal(err)
	}

	// PROJ 9.5, EPSG:4326 to EPSG:23090
	if math.Abs(east-599101.467) > 0.01 || math.Abs(north-6511469.108) > 0.01 {
		t.Errorf("got %.3f %.3f, want 599101.467 6511469.108", east, north)
	}

	from, err := Transform(23090, 4326)
	if err != nil {
		t.Fatal(err)
	}

	backLon, backLat, _, err := from(east, north, 0)
	if err != nil {
		t.Fatal(err)
	}

	if math.Abs(backLon-lon) > 1e-6 || math.Abs(backLat-lat) > 1e-6 {
		t.Errorf("round trip gives %.9f %.9f, want %.9f %.9f", backLon, backLat, lon, lat)
	}
}
