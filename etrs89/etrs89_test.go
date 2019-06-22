package etrs89

import (
	"fmt"
	"testing"

	"github.com/wroge/wgs84"
	"github.com/wroge/wgs84/system"
)

func ExampleWithSystem() {
	EPSG4258 := WithSystem(system.LonLat())
	lon, lat, h := wgs84.To(EPSG4258)(9, 52, 0)
	fmt.Println(lon, lat, h)
}

func BenchmarkUTM(b *testing.B) {
	for i := 0; i < b.N; i++ {
		UTM(32).From(wgs84.LonLat())(9, 52, 0)
	}
}

func TestUTM(t *testing.T) {
	longitude := 10.0
	latitude := 52.0
	east := 568649.70
	north := 5761510.32

	e, n, _ := UTM(32).From(wgs84.LonLat()).Round(2)(longitude, latitude, 0)
	if e != east || n != north {
		t.Error("UTM is not correct.")
	}
}
