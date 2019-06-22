package nad83

import (
	"testing"

	"github.com/wroge/wgs84"
)

func TestConnecticut(t *testing.T) {
	longitude := -80.0
	latitude := 50.0
	east := -220357.68
	north := 1196767.67

	e, n, _ := Connecticut().From(wgs84.LonLat()).Round(2)(longitude, latitude, 0)
	if e != east || n != north {
		t.Error("Connecticut is not correct.")
	}
}
