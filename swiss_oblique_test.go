package wgs84

import (
	"math"
	"testing"
)

func lv03() SwissObliqueMercator {
	return SwissObliqueMercator{
		Latf:   46.9524055555556,
		Lonf:   7.43958333333333,
		Scale:  1,
		Eastf:  600000,
		Northf: 200000,
	}
}

func TestSwissObliqueMercatorFalseOrigin(t *testing.T) {
	som, sp := lv03(), CH1903.Spheroid

	lon, lat, _ := som.ToGeographic(som.Eastf, som.Northf, 0, sp)
	if math.Abs(lon-som.Lonf) > 1e-9 || math.Abs(lat-som.Latf) > 1e-9 {
		t.Errorf("false origin maps to %.9f %.9f, want %.9f %.9f", lon, lat, som.Lonf, som.Latf)
	}
}

func TestSwissObliqueMercatorRoundTrip(t *testing.T) {
	som, sp := lv03(), CH1903.Spheroid

	for lon := 6.0; lon <= 10.5; lon += 0.25 {
		for lat := 45.8; lat <= 47.9; lat += 0.25 {
			east, north, _ := som.FromGeographic(lon, lat, 0, sp)
			back, backLat, _ := som.ToGeographic(east, north, 0, sp)
			if math.Abs(back-lon) > 1e-9 || math.Abs(backLat-lat) > 1e-9 {
				t.Fatalf("%.4f %.4f round-trips to %.9f %.9f", lon, lat, back, backLat)
			}
		}
	}
}
