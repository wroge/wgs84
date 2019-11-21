package wgs84_test

import (
	"fmt"
	"testing"

	"github.com/wroge/wgs84"
)

func Test(t *testing.T) {
	epsg := wgs84.EPSG()
	epsg.Add(1110, nil)
	epsg.Add(1111, wgs84.GeocentricReferenceSystem{})
	epsg.Add(1112, wgs84.GeographicReferenceSystem{})
	p := wgs84.ProjectedReferenceSystem{}
	p.Area = wgs84.AreaFunc(nil)
	epsg.Add(1113, p)
	epsg.Add(1114, wgs84.Helmert(6378137, 298.257223563, 0, 0, 0, 0, 0, 0, 0).LonLat())
	for lon := -185.0; lon < 185; lon += 1.5 {
		for lat := -95.0; lat < 95; lat += 1.5 {
			codes := epsg.CodesCover(lon, lat)
			if len(codes) == 0 {
				continue
			}
			var from = 4326
			var a, b, c = lon, lat, 0.0
			for _, code := range codes {
				a, b, c = wgs84.Transform(epsg.Code(from), epsg.Code(code))(a, b, c)
				from = code
			}
			a2, b2, _ := epsg.Transform(from, 4326).Round(2)(a, b, c)
			a, b, _, err := epsg.SafeTransform(from, 4326).Round(2)(a, b, c)
			if a != lon || b != lat || err != nil || a2 != a || b2 != b {
				fmt.Println(from, lon, lat, a, b)
				panic("Failed")
			}
		}
	}
}
