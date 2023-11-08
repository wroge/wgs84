//nolint:depguard,cyclop,lll,varnamelen
package wgs84_test

import (
	"math"
	"testing"

	"github.com/wroge/wgs84/v2"
)

type CodeTest struct {
	Name string
	From Coordinate
	To   Coordinate
}

type Coordinate struct {
	EPSG    int
	A, B, C float64
}

func TestCodes(t *testing.T) {
	t.Parallel()

	table := []CodeTest{
		{
			Name: "EPSG(3035,4258)",
			From: Coordinate{
				EPSG: 3035,
				A:    3962799.45,
				B:    2999718.85,
				C:    0,
			},
			To: Coordinate{
				EPSG: 4258,
				A:    5,
				B:    50,
				C:    0,
			},
		},
	}

	for _, each := range table {
		a, b, c := wgs84.Transform(wgs84.EPSG(each.From.EPSG), wgs84.EPSG(each.To.EPSG))(each.From.A, each.From.B, each.From.C)
		if math.IsNaN(a) || math.IsNaN(b) || math.IsNaN(c) || math.Abs(a-each.To.A) > 0.01 || math.Abs(b-each.To.B) > 0.01 || math.Abs(c-each.To.C) > 0.01 {
			t.Fatalf("%s: %d to %d: RESULT=[%f %f %f] EXPECT=[%f %f %f]", each.Name, each.From.EPSG, each.To.EPSG, a, b, c, each.To.A, each.To.B, each.To.C)
		}

		a, b, c = wgs84.Transform(wgs84.EPSG(each.To.EPSG), wgs84.EPSG(each.From.EPSG))(each.To.A, each.To.B, each.To.C)
		if math.IsNaN(a) || math.IsNaN(b) || math.IsNaN(c) || math.Abs(a-each.From.A) > 0.01 || math.Abs(b-each.From.B) > 0.01 || math.Abs(c-each.From.C) > 0.01 {
			t.Fatalf("%s: %d to %d: RESULT=[%f %f %f] EXPECT=[%f %f %f]", each.Name, each.To.EPSG, each.From.EPSG, a, b, c, each.From.A, each.From.B, each.From.C)
		}
	}
}
