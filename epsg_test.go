package wgs84_test

import (
	"fmt"
	"testing"

	"github.com/wroge/wgs84"
)

func Test(t *testing.T) {
	epsg := wgs84.EPSG()
	epsg.Add(1111, wgs84.Projection{
		GeodeticDatum: wgs84.Datum{},
	})
	for lon := -185.0; lon < 185; lon += 1.5 {
		for lat := -95.0; lat < 95; lat += 1.5 {
			var from = 4326
			var a, b, c = lon, lat, 0.0
			codes := epsg.CodesCover(lon, lat)
			if codes == nil || len(codes) == 0 {
				continue
			}
			for _, code := range codes {
				a, b, c = wgs84.Transform(epsg.Code(from), epsg.Code(code))(a, b, c)
				from = code
			}
			a, b, c = wgs84.Transform(epsg.Code(from), epsg.Code(4326)).Round(3)(a, b, c)
			if a != lon || b != lat {
				fmt.Println(from, lon, lat, a, b)
				t.Fatal(a, b, c)
			}
		}
	}
}
