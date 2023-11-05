package main

import (
	"fmt"

	"github.com/wroge/wgs84/v2"
)

func main() {
	east, north, h := wgs84.Transform(wgs84.EPSG(4326), wgs84.EPSG(25832), 10.5, 51.5, 0)

	fmt.Println(east, north, h)
}
