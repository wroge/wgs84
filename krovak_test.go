package wgs84

import (
	"math"
	"testing"
)

// TestKrovakOfficialExample validates the Krovak projection against the worked
// numerical example in IOGP Geomatics Guidance Note 7-2 (publication 373-07-02),
// section 3.1.2.1 (coordinate operation method 9819).
//
// The example is datum-independent: it projects geographic coordinates on the
// Bessel 1841 ellipsoid directly, so it exercises only the projection math and
// not any datum transformation. EPSG:5514 is the East-North variant, so the
// published Southing/Westing map to (East, North) = (-Westing, -Southing).
func TestKrovakOfficialExample(t *testing.T) {
	// Bessel 1841 with the inverse-flattening exactly as printed in IOGP 373-07-02
	// (299.15281, not the repo's usual 299.1528128) so this reproduces the
	// document's published result from its own stated inputs.
	bessel := NewSpheroid(6377397.155, 299.15281)
	base := Geographic(nil, bessel)

	crs := Krovak(base, 24.8333333333333, 49.5, 30.2881397527778, 78.5, 0.9999, 0, 0)

	// Example point: 50°12'32.442"N, 16°50'59.179"E (Greenwich).
	lat := 50.0 + 12.0/60 + 32.442/3600
	lon := 16.0 + 50.0/60 + 59.179/3600

	// Published result: Southing X = 1050538.64, Westing Y = 568991.00.
	wantE := -568991.00
	wantN := -1050538.64
	const tol = 0.02 // metres; absorbs rounding of the published in/out values

	gotE, gotN, _ := Transform(base, crs)(lon, lat, 0)

	if math.Abs(gotE-wantE) > tol {
		t.Errorf("East = %.3f, want %.2f (±%g)", gotE, wantE, tol)
	}
	if math.Abs(gotN-wantN) > tol {
		t.Errorf("North = %.3f, want %.2f (±%g)", gotN, wantN, tol)
	}

	// Round-trip back to geographic.
	gotLon, gotLat, _ := Transform(crs, base)(gotE, gotN, 0)
	if math.Abs(gotLon-lon) > 1e-7 || math.Abs(gotLat-lat) > 1e-7 {
		t.Errorf("round-trip = (%.8f, %.8f), want (%.8f, %.8f)", gotLon, gotLat, lon, lat)
	}
}
