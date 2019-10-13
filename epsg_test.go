package wgs84_test

import (
	"testing"

	"github.com/wroge/wgs84"
)

func Test(t *testing.T) {
	epsg := wgs84.EPSG()
	var from = 4326
	var a, b, c = 9.0, 52.0, 0.0
	for _, code := range epsg.CodesCover(a, b) {
		a, b, c = wgs84.Transform(epsg.Code(from), epsg.Code(code))(a, b, c)
		from = code
	}
	a, b, c = wgs84.Transform(epsg.Code(from), epsg.Code(4326)).Round(0)(a, b, c)
	if a != 9 || b != 52 || c != 0 {
		t.Fatal(a, b, c)
	}
}
