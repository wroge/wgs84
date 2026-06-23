[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/wroge/wgs84/v2@v2.0.0-alpha.16)

## WGS84 - Coordinate Transformations

```
go get github.com/wroge/wgs84/v2@v2.0.0-alpha.16
```  

### Example

```go
package main

import (
	"fmt"
	"os"

	"github.com/wroge/wgs84/v2"
	// alternative to wgs84.RegisterGridFS
	// _ "github.com/wroge/wgs84/grids/osgb36"
)

func main() {
	EPSG4277 := wgs84.EPSG[4277].Load(-2, 50.7)

	fmt.Println(EPSG4277)

	conv := wgs84.EPSG[4326].TransformTo(EPSG4277).Round(9, 9, 1)

	east, north, h, err := conv(-2, 50.7, 0)
	if err != nil {
		panic(err)
	}

	fmt.Println(east, north, h)

	wgs84.RegisterGridFS("", os.DirFS("/opt/homebrew/opt/proj/share/proj"))

	EPSG4277 = wgs84.EPSG[4277].Load(-2, 50.7)

	fmt.Println(EPSG4277)

	conv = wgs84.EPSG[4326].TransformTo(EPSG4277).Round(9, 9, 1)

	east, north, h, err = conv(-2, 50.7, 0)
	if err != nil {
		panic(err)
	}

	fmt.Println(east, north, h)

	EPSG4277 = wgs84.EPSG[4277].Filter(func(t wgs84.Transformation) bool { return t.Grid == "" }).Load(-2, 50.7)

	fmt.Println(EPSG4277)

	conv = wgs84.EPSG[4326].TransformTo(EPSG4277).Round(9, 9, 1)

	east, north, h, err = conv(-2, 50.7, 0)
	if err != nil {
		panic(err)
	}

	fmt.Println(east, north, h)

	EPSG4277 = wgs84.EPSG[4277].Filter(func(t wgs84.Transformation) bool { return t.Accuracy > 2 }).Load(-2, 50.7)

	fmt.Println(EPSG4277)

	conv = wgs84.EPSG[4326].TransformTo(EPSG4277).Round(9, 9, 1)

	east, north, h, err = conv(-2, 50.7, 0)
	if err != nil {
		panic(err)
	}

	fmt.Println(east, north, h)
}

// +proj=longlat +a=6.377563396e+06 +rf=299.3249646 +towgs84=446.448,-125.157,542.06,0.15,0.247,0.842,-20.489
// -1.998636231 50.699424279 -47.5
// +proj=longlat +a=6.377563396e+06 +rf=299.3249646 +nadgrids=uk_os_OSTN15_NTv2_OSGBtoETRS.tif
// -1.998642581 50.69943404 -0
// +proj=longlat +a=6.377563396e+06 +rf=299.3249646 +towgs84=446.448,-125.157,542.06,0.15,0.247,0.842,-20.489
// -1.998636231 50.699424279 -47.5
// +proj=longlat +a=6.377563396e+06 +rf=299.3249646 +towgs84=370.936,-108.938,435.682,0,0,0,0
// -1.998642097 50.699434989 -47.5

// echo "-2 50.7 0" | cs2cs -f "%.9f" \
//   +proj=longlat +datum=WGS84 +to \
//   +proj=longlat +a=6.377563396e+06 +rf=299.3249646 +towgs84=446.448,-125.157,542.06,0.15,0.247,0.842,-20.489
// -1.998636179    50.699424268 0.000000000

// echo "-2 50.7 0" | cs2cs -f "%.9f" \
//   +proj=longlat +datum=WGS84 +to \
//   +proj=longlat +a=6.377563396e+06 +rf=299.3249646 +nadgrids=uk_os_OSTN15_NTv2_OSGBtoETRS.tif
// -1.998642581    50.699434040 0.000000000

// echo "-2 50.7 0" | cs2cs -f "%.9f" \
//   +proj=longlat +datum=WGS84 +to \
//   +proj=longlat +a=6.377563396e+06 +rf=299.3249646 +towgs84=446.448,-125.157,542.06,0.15,0.247,0.842,-20.489
// -1.998636179    50.699424268 0.000000000

// echo "-2 50.7 0" | cs2cs -f "%.9f" \
//   +proj=longlat +datum=WGS84 +to \
//   +proj=longlat +a=6.377563396e+06 +rf=299.3249646 +towgs84=370.936,-108.938,435.682,0,0,0,0
// -1.998642097    50.699434989 0.000000000
```

