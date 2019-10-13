package wgs84_test

import (
	"fmt"
	"testing"

	"github.com/wroge/wgs84"
)

func Test(t *testing.T) {
	for lon := -179.0; lon < 180; lon += 1.5 {
		for lat := -89.0; lat < 90; lat += 1.5 {
			epsg := wgs84.EPSG()
			var from = 4326
			var a, b, c = lon, lat, 0.0
			for _, code := range epsg.CodesCover(lon, lat) {
				a, b, c = wgs84.Transform(epsg.Code(from), epsg.Code(code))(a, b, c)
				from = code
			}
			a, b, c = wgs84.Transform(epsg.Code(from), epsg.Code(4326)).Round(3)(a, b, c)
			if a != lon || b != lat {
				fmt.Println(lon, lat, a, b)
				t.Fatal(a, b, c)
			}
		}
	}
}
