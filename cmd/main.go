package main

import (
	"fmt"

	"github.com/wroge/wgs84"
)

func main() {
	fmt.Println(wgs84.EPSG().Transform(4326, 3035).Round(2)(5, 50, 0))

	fmt.Println(wgs84.ETRS89LambertAzimuthalEqualArea().To(wgs84.LonLat()).Round(2)(3962799.45, 2999718.85, 0))
}
