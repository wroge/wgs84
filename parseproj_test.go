package wgs84

import (
	"fmt"
	"math"
	"testing"
)

func floatEq(a, b float64) bool {
	if math.IsNaN(a) && math.IsNaN(b) {
		return true
	}
	if a == b {
		return true
	}
	diff := math.Abs(a - b)
	if diff <= 1e-9 {
		return true
	}
	scale := math.Max(math.Abs(a), math.Abs(b))
	return scale > 0 && diff/scale <= 1e-9
}

func spheroidEq(a, b Spheroid) bool {
	return floatEq(a.A, b.A) && floatEq(a.Fi, b.Fi)
}

func helmertEq(a, b Helmert) bool {
	return floatEq(a.Tx, b.Tx) &&
		floatEq(a.Ty, b.Ty) &&
		floatEq(a.Tz, b.Tz) &&
		floatEq(a.Rx, b.Rx) &&
		floatEq(a.Ry, b.Ry) &&
		floatEq(a.Rz, b.Rz) &&
		floatEq(a.Ds, b.Ds)
}

func transformationEq(a, b Transformation) bool {
	if a.Grid != b.Grid {
		return false
	}

	return helmertEq(a.Helmert, b.Helmert)
}

func datumRoundtripEq(want, got Datum) bool {
	if !spheroidEq(want.Spheroid, got.Spheroid) {
		return false
	}

	wantT := datumFromString(want)
	gotT := datumFromString(got)

	if len(wantT.Transformations) != len(gotT.Transformations) {
		return false
	}

	for i := range wantT.Transformations {
		if !transformationEq(wantT.Transformations[i], gotT.Transformations[i]) {
			return false
		}
	}

	return true
}

func datumFromString(d Datum) Datum {
	if len(d.Transformations) == 0 {
		return Datum{Spheroid: d.Spheroid}
	}

	t := Transformation{}
	first := d.Transformations[0]

	switch {
	case first.Grid != "":
		t.Grid = first.Grid
	case first.Helmert != (Helmert{}):
		t.Helmert = first.Helmert
	}

	return Datum{
		Spheroid:        d.Spheroid,
		Transformations: []Transformation{t},
	}
}

func coordinateSystemEq(want, got CoordinateSystem) bool {
	switch wantCS := want.(type) {
	case Geocentric:
		_, ok := got.(Geocentric)
		return ok
	case Geographic:
		_, ok := got.(Geographic)
		return ok
	case WebMercator:
		_, ok := got.(WebMercator)
		return ok
	case TransverseMercator:
		gotCS, ok := got.(TransverseMercator)
		if !ok {
			return false
		}
		return floatEq(wantCS.Latf, gotCS.Latf) &&
			floatEq(wantCS.Lonf, gotCS.Lonf) &&
			floatEq(wantCS.Scale, gotCS.Scale) &&
			floatEq(wantCS.Eastf, gotCS.Eastf) &&
			floatEq(wantCS.Northf, gotCS.Northf)
	case SwissObliqueMercator:
		gotCS, ok := got.(SwissObliqueMercator)
		if !ok {
			return false
		}
		return floatEq(wantCS.Latf, gotCS.Latf) &&
			floatEq(wantCS.Lonf, gotCS.Lonf) &&
			floatEq(wantCS.Scale, gotCS.Scale) &&
			floatEq(wantCS.Eastf, gotCS.Eastf) &&
			floatEq(wantCS.Northf, gotCS.Northf)
	case Krovak:
		gotCS, ok := got.(Krovak)
		if !ok {
			return false
		}
		return floatEq(wantCS.Latf, gotCS.Latf) &&
			floatEq(wantCS.Lonf, gotCS.Lonf) &&
			floatEq(wantCS.Alpha, gotCS.Alpha) &&
			floatEq(wantCS.Scale, gotCS.Scale) &&
			floatEq(wantCS.Eastf, gotCS.Eastf) &&
			floatEq(wantCS.Northf, gotCS.Northf)
	case AlbersConicEqualArea:
		gotCS, ok := got.(AlbersConicEqualArea)
		if !ok {
			return false
		}
		return floatEq(wantCS.Latf, gotCS.Latf) &&
			floatEq(wantCS.Lonf, gotCS.Lonf) &&
			floatEq(wantCS.Sp1, gotCS.Sp1) &&
			floatEq(wantCS.Sp2, gotCS.Sp2) &&
			floatEq(wantCS.Eastf, gotCS.Eastf) &&
			floatEq(wantCS.Northf, gotCS.Northf)
	case LambertAzimuthalEqualArea:
		gotCS, ok := got.(LambertAzimuthalEqualArea)
		if !ok {
			return false
		}
		return floatEq(wantCS.Latf, gotCS.Latf) &&
			floatEq(wantCS.Lonf, gotCS.Lonf) &&
			floatEq(wantCS.Eastf, gotCS.Eastf) &&
			floatEq(wantCS.Northf, gotCS.Northf)
	case LambertConformalConic2SP:
		gotCS, ok := got.(LambertConformalConic2SP)
		if !ok {
			return false
		}
		return floatEq(wantCS.Latf, gotCS.Latf) &&
			floatEq(wantCS.Lonf, gotCS.Lonf) &&
			floatEq(wantCS.Sp1, gotCS.Sp1) &&
			floatEq(wantCS.Sp2, gotCS.Sp2) &&
			floatEq(wantCS.Eastf, gotCS.Eastf) &&
			floatEq(wantCS.Northf, gotCS.Northf)
	case LambertConformalConic1SP:
		gotCS, ok := got.(LambertConformalConic1SP)
		if !ok {
			return false
		}
		return floatEq(wantCS.Latf, gotCS.Latf) &&
			floatEq(wantCS.Lonf, gotCS.Lonf) &&
			floatEq(wantCS.Scale, gotCS.Scale) &&
			floatEq(wantCS.Eastf, gotCS.Eastf) &&
			floatEq(wantCS.Northf, gotCS.Northf)
	default:
		return false
	}
}

func TestEPSGStaticParseProjRoundtrip(t *testing.T) {
	for code, want := range EPSG {
		t.Run(fmt.Sprintf("EPSG:%d", code), func(t *testing.T) {
			proj := want.String()

			got, err := ParseProj(proj)
			if err != nil {
				t.Fatalf("ParseProj(%q): %v", proj, err)
			}

			if !coordinateSystemEq(want.CoordinateSystem, got.CoordinateSystem) {
				t.Fatalf("coordinate system mismatch:\nwant %T %#v\ngot  %T %#v\nproj %q",
					want.CoordinateSystem, want.CoordinateSystem,
					got.CoordinateSystem, got.CoordinateSystem,
					proj,
				)
			}

			if !datumRoundtripEq(want.Datum, got.Datum) {
				t.Fatalf("datum mismatch:\nwant %+v\ngot  %+v\nproj %q",
					datumFromString(want.Datum),
					datumFromString(got.Datum),
					proj,
				)
			}
		})
	}
}
