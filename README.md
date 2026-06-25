[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/wroge/wgs84/v2@v2.0.0-alpha.18)

## WGS84 - Coordinate Transformations

```
go get github.com/wroge/wgs84/v2@v2.0.0-alpha.18
```  

This *zero-dependency* module supports

- NTv2 grids via ```wgs84.RegisterGridFS``` or submodules (dhdn90, ed50, nad83, osgb36),
- Geocentric, Geographic, WebMercator, TransverseMercator, SwissObliqueMercator, LambertConformalConic1SP, LambertConformalConic2SP, LambertAzimuthalEqualArea, Krovak and AlbersConicEqualArea coordinate systems,
- ~ 212 EPSG-Codes,
- and proj4 coordinate reference system definitions.

### Example

```go
package main

import (
	"fmt"

	_ "github.com/wroge/wgs84/grids/osgb36" // import OSTN15 grid
	"github.com/wroge/wgs84/v2"
)

func main() {
	// 1. convert wgs84 coordinates to webmercator.
	// transform, err := wgs84.Transform(wgs84.EPSG[4326], wgs84.EPSG[3857])
	// transform, err := wgs84.Transform(4326, 3857)
	// transform, err := wgs84.Transform(4326, wgs84.WebMercator{})
	// transform, err := wgs84.Transform(4326, "+proj=merc +datum=WGS84")
	transform, err := wgs84.Transform("+proj=longlat +datum=WGS84", "+proj=merc")
	if err != nil {
		panic(err)
	}

	east, north, h, err := transform.Round(3, 3, 0)(10, 50, 0)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.3f %.3f %.0f\n", east, north, h)
	// 1113194.908 6446275.841 0

	// 2. convert wgs84 coordinates to osgb36 british national grid using imported OSTN15 grid.
	// transform, err = wgs84.Transform(4326, 27700)
	// transform, err = wgs84.Transform(4326, "+proj=tmerc +lat_0=49 +lon_0=-2 +k=0.9996012717 +x_0=400000 +y_0=-100000 +a=6377563.396 +rf=299.3249646 +nadgrids=@uk_os_OSTN15_NTv2_OSGBtoETRS.tif")
	transform, err = wgs84.Transform(4326, wgs84.CoordinateReferenceSystem{
		CoordinateSystem: wgs84.TransverseMercator{
			Latf:   49,
			Lonf:   -2,
			Scale:  0.9996012717,
			Eastf:  400000,
			Northf: -100000,
		},
		Datum: wgs84.OSGB36,
	})
	if err != nil {
		panic(err)
	}

	east, north, h, err = transform.Round(3, 3, 0)(-5, 55, 0)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.3f %.3f %.0f\n", east, north, h)
	// 208215.245 571385.903 0
}
```

