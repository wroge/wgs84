package wgs84_test

import (
	"math"
	"testing"

	"github.com/wroge/wgs84/v2"
)

type Projection interface {
	ToGeographic(east, north float64) (lon, lat float64)
	FromGeographic(lon, lat float64) (east, north float64)
}

type ProjectionTest struct {
	Name        string
	Projection  Projection
	Lon, Lat    float64
	East, North float64
}

func TestProjections(t *testing.T) {
	table := []ProjectionTest{
		{
			Name: "LambertAzimuthalEqualArea",
			Projection: wgs84.LambertAzimuthalEqualArea{
				LatitudeOfCenter:  52,
				LongitudeOfCenter: 10,
				FalseEasting:      4321000,
				FalseNorthing:     3210000,
				Geographic: wgs84.Geographic{
					Spheroid: wgs84.GRS80,
				},
			},
			Lon:   5,
			Lat:   50,
			East:  3962799.45,
			North: 2999718.85,
		},
		{
			Name: "AlbersEqualAreaConic",
			Projection: wgs84.AlbersConicEqualArea{
				LatitudeOfCenter:  toDeg(45, 34, 8.3172),
				LongitudeOfCenter: -toDeg(84, 27, 21.4380),
				StandardParallel1: toDeg(42, 7, 21.9864),
				StandardParallel2: toDeg(49, 0, 54.648),
				FalseEasting:      1000000,
				FalseNorthing:     1000000,
				Geographic: wgs84.Geographic{
					Spheroid: wgs84.GRS80,
				},
			},
			Lat:   toDeg(42, 45, 0),
			Lon:   -toDeg(78, 45, 0),
			East:  1466493.492,
			North: 702903.006,
		},
		{
			Name: "LambertConformalConic2SP",
			Projection: wgs84.LambertConformalConic2SP{
				LatitudeOfOrigin:  toDeg(27, 50, 0),
				CentralMeridian:   -99,
				StandardParallel1: toDeg(28, 23, 0),
				StandardParallel2: toDeg(30, 17, 0),
				FalseEasting:      2000000.00 / (3937.0 / 1200.0),
				FalseNorthing:     0,
				Geographic: wgs84.Geographic{
					Spheroid: wgs84.Clarke1866,
				},
			},
			Lat:   toDeg(28, 30, 0),
			Lon:   -96,
			East:  2963503.91 / (3937.0 / 1200.0),
			North: 254759.80 / (3937.0 / 1200.0),
		},
		{
			Name: "TransverseMercator",
			Projection: wgs84.TransverseMercator{
				LatitudeOfOrigin: 49,
				CentralMeridian:  -2,
				ScaleFactor:      0.9996012717,
				FalseEasting:     400000.00,
				FalseNorthing:    -100000.00,
				Geographic: wgs84.Geographic{
					Spheroid: wgs84.Airy1830,
				},
			},
			Lat:   50.5,
			Lon:   0.5,
			East:  577274.99,
			North: 69740.50,
		},
		{
			Name: "WebMercator",
			Projection: wgs84.WebMercator{
				Geographic: wgs84.WGS84,
			},
			Lat:   toDeg(24, 22, 54.433),
			Lon:   -toDeg(100, 20, 0),
			East:  -11169055.58,
			North: 2800000.00,
		},
	}

	for _, each := range table {
		east, north := each.Projection.FromGeographic(each.Lon, each.Lat)
		if math.IsNaN(east) || math.IsNaN(north) || math.Abs(east-each.East) > 0.01 || math.Abs(north-each.North) > 0.01 {
			t.Fatal(each.Name+" FromGeographic", east, north, each.East, each.North)
		}

		lon, lat := each.Projection.ToGeographic(each.East, each.North)
		if math.IsNaN(lon) || math.IsNaN(lat) || math.Abs(lon-each.Lon) > 0.0000001 || math.Abs(lat-each.Lat) > 0.0000001 {
			t.Fatal(each.Name+" ToGeographic", lon, lat)
		}
	}
}

func toDeg(deg, min, sec float64) float64 {
	return deg + min/60 + sec/3600
}
