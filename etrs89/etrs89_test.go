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
