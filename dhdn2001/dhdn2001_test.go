package dhdn2001

import (
	"testing"

	"github.com/wroge/wgs84"
)

func TestGK(t *testing.T) {
	longitude := 10.0
	latitude := 52.0
	east := 3568749.80
	north := 5763376.67

	e, n, _ := GK(3).From(wgs84.LonLat()).Round(2)(longitude, latitude, 0)
	if e != east || n != north {
		t.Error("GK is not correct.")
	}
}
