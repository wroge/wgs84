//nolint:varnamelen,exhaustivestruct,exhaustruct
package wgs84_test

import (
	"testing"

	"github.com/wroge/wgs84"
)

func Test(t *testing.T) {
	t.Parallel()

	epsg := wgs84.EPSG()
	epsg.Add(1110, nil)

	p := wgs84.ProjectedReferenceSystem{
		Area: wgs84.AreaFunc(nil),
	}

	epsg.Add(1113, p)
	epsg.Add(1114, wgs84.Helmert(wgs84.A, wgs84.Fi, 0, 0, 0, 0, 0, 0, 0).LonLat())

	for lon := -185.0; lon < 185; lon += 1.5 {
		for lat := -95.0; lat < 95; lat += 1.5 {
			codes := epsg.CodesCover(lon, lat)
			if len(codes) == 0 {
				continue
			}

			var (
				from    = 4326
				a, b, c = lon, lat, 0.0
			)

			for _, code := range codes {
				a, b, c = wgs84.Transform(epsg.Code(from), epsg.Code(code))(a, b, c)
				from = code
			}

			a2, b2, _ := epsg.Transform(from, 4326).Round(2)(a, b, c)

			a, b, _, err := epsg.SafeTransform(from, 4326).Round(2)(a, b, c)
			if a != lon || b != lat || err != nil || a2 != a || b2 != b {
				t.Fatal("Failed")
			}
		}
	}
}

func Test27572(t *testing.T) {
	t.Parallel()

	// Reference values from PROJ 9.8.1 (same engine as epsg.io)
	cases := []struct{ lon, lat, wantE, wantN float64 }{
		{2.3522, 48.8566, 601152.30, 2428695.90},  // Paris
		{4.8357, 45.7640, 794390.28, 2087943.39},  // Lyon
		{5.3698, 43.2965, 846499.55, 1815214.13},  // Marseille
	}

	transform := wgs84.EPSG().Transform(4326, 27572)

	for _, tc := range cases {
		e, n, _ := transform.Round(2)(tc.lon, tc.lat, 0)
		if e != tc.wantE || n != tc.wantN {
			t.Fatalf("lon=%v lat=%v: got (%v, %v), want (%v, %v)",
				tc.lon, tc.lat, e, n, tc.wantE, tc.wantN)
		}
	}
}

func Test2(t *testing.T) {
	t.Parallel()

	east, north, h := wgs84.To(wgs84.WebMercator())(9, 52, 0)

	lon, lat, h := wgs84.From(wgs84.WebMercator()).Round(1)(east, north, h)
	if lon != 9 ||
		lat != 52 ||
		h != 0 {
		t.Fatal("Failed (2)")
	}
}
