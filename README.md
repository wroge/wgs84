[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/wroge/wgs84/v2@v2.0.0-alpha.17)

## WGS84 - Coordinate Transformations

```
go get github.com/wroge/wgs84/v2@v2.0.0-alpha.17
```  

### Example

```go
package main

import (
	"fmt"

	_ "github.com/wroge/wgs84/grids/dhdn90" // import osgb36 grid
	_ "github.com/wroge/wgs84/grids/osgb36" // import dhdn90 grid
	"github.com/wroge/wgs84/v2"
)

func main() {
	// 1. convert wgs84 longitude/latitude to webmercator
	convert := wgs84.EPSG[4326].ConvertTo(wgs84.WebMercator{}).Round(3, 3, 0)

	east, north, h, err := convert(10, 50, 0)

	fmt.Println(east, north, h, err)
	// 1.113194908e+06 6.446275841e+06 0 <nil>

	// 2. transform wgs84 longitude/latitude to british national grid using imported osgb36
	fmt.Println(wgs84.EPSG[27700])
	// +proj=tmerc +lat_0=49 +lon_0=-2 +k=0.9996012717 +x_0=400000 +y_0=-100000 +a=6377563.396 +rf=299.3249646 +nadgrids=uk_os_OSTN15_NTv2_OSGBtoETRS.tif

	transform := wgs84.EPSG[4326].TransformTo(wgs84.EPSG[27700]).Round(3, 3, 0)

	east, north, h, err = transform(-5, 55, 0)

	fmt.Println(east, north, h, err)
	// 208215.245 571385.903 -0 <nil>

	// 3. parse proj4 format and import dhdn90 grid

	crs, err := wgs84.ParseProj("+proj=tmerc +lat_0=0 +lon_0=9 +k=1 +x_0=3500000 +y_0=0 +ellps=bessel +nadgrids=de_adv_BETA2007.tif +units=m +no_defs +type=crs")

	fmt.Println(crs, err)
	// +proj=tmerc +lat_0=0 +lon_0=9 +k=1 +x_0=3500000 +y_0=0 +a=6377397.155 +rf=299.1528128 +nadgrids=de_adv_BETA2007.tif <nil>

	transform = wgs84.EPSG[4326].TransformTo(crs).Round(3, 3, 0)

	east, north, h, err = transform(10, 50, 0)

	fmt.Println(east, north, h, err)
	// 3.5717699e+06 5.540887024e+06 -0 <nil>

	// Is the same as
	fmt.Println(wgs84.EPSG[31467])
	// +proj=tmerc +lat_0=0 +lon_0=9 +k=1 +x_0=3500000 +y_0=0 +a=6377397.155 +rf=299.1528128 +nadgrids=de_adv_BETA2007.tif

	fmt.Println(wgs84.EPSG[4326].TransformTo(wgs84.EPSG[31467]).Round(3, 3, 0)(10, 50, 0))
	// 3.5717699e+06 5.540887024e+06 -0 <nil>
}
```

