[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/wroge/wgs84/v2@v2.0.0-alpha.15)

## WGS84 - Coordinate Transformations

```
go get github.com/wroge/wgs84/v2@v2.0.0-alpha.15
```  

### Example

```go
package main

import (
	"fmt"
	"os"

	"github.com/wroge/wgs84/v2"
	// Alternativ to wgs84.RegisterGridFS
	// _ "github.com/wroge/wgs84/v2/grid/osgb36"
)

func main() {
	conv := wgs84.EPSG[4326].TransformTo(wgs84.EPSG[4277])

	east, north, h, err := conv(-2, 50.7, 0)
	if err != nil {
		panic(err)
	}

	fmt.Println(east, north, h)
	// 	-1.9986362310906312 50.69942427880695 -47.549131196923554

    // For mac: brew install proj
	wgs84.RegisterGridFS("", os.DirFS("/opt/homebrew/opt/proj/share/proj"))

	conv = wgs84.EPSG[4326].TransformTo(wgs84.EPSG[4277])

	east, north, h, err = conv(-2, 50.7, 0)
	if err != nil {
		panic(err)
	}

	fmt.Println(east, north, h)
	// -1.998642581025955 50.699434040486324 -9.313225746154785e-10

	conv = wgs84.EPSG[4326].TransformTo(wgs84.EPSG[4277].Filter(func(t wgs84.Transformation) bool { return t.Grid == "" }))

	east, north, h, err = conv(-2, 50.7, 0)
	if err != nil {
		panic(err)
	}

	fmt.Println(east, north, h)
	// -1.9986362310906312 50.69942427880695 -47.549131196923554

	conv = wgs84.EPSG[4326].TransformTo(wgs84.EPSG[4277].Filter(func(t wgs84.Transformation) bool { return t.Accuracy > 2 }))

	east, north, h, err = conv(-2, 50.7, 0)
	if err != nil {
		panic(err)
	}

	fmt.Println(east, north, h)
	// -1.9986420967349676 50.69943498890854 -47.51741572842002
}

// echo -2 50.7 0 | /opt/homebrew/bin/cct -d 6 +proj=pipeline +step +proj=unitconvert +xy_in=deg +xy_out=rad +step +inv +proj=hgridshift +grids=uk_os_OSTN15_Grid_OSGBtoETRS.tif +step +proj=unitconvert +xy_in=rad +xy_out=deg
//      -1.998643       50.699434      0.000000           inf
```

