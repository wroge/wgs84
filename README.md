# WGS84

A pure go package for the conversion and transformation of geodesic coordinates.

```go
go get github.com/wroge/wgs84
```

[![GoDoc](http://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/github.com/wroge/wgs84)
[![GoWalker](https://img.shields.io/badge/Go_Walker-Doc-blue.svg?style=flat-square)](https://gowalker.org/github.com/wroge/wgs84)
[![Go Report Card](https://goreportcard.com/badge/github.com/wroge/wgs84?style=flat-square)](https://goreportcard.com/report/github.com/wroge/wgs84)
[![GolangCI](https://golangci.com/badges/github.com/wroge/wgs84.svg)](https://golangci.com/r/github.com/wroge/wgs84)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/wroge/wgs84.svg?style=social)](https://github.com/wroge/wgs84/tags)

Subpackages with predefined coordinate reference systems within a geodetic datum:

- [DHDN2001](https://github.com/wroge/wgs84/tree/master/dhdn2001)
- [ETRS89](https://github.com/wroge/wgs84/tree/master/etrs89)
- [NAD83](https://github.com/wroge/wgs84/tree/master/nad83)
- [OSGB36](https://github.com/wroge/wgs84/tree/master/osgb36)
- [RGF93](https://github.com/wroge/wgs84/tree/master/rgf93)

Package for EPSG-Code support: [EPSG](https://github.com/wroge/wgs84/tree/master/epsg)

## Features

- Helmert Transformation
- Geocentric Translation
- Geocentric/Cartesian (XYZ)
- Geographic (LonLat)
- Web Mercator
- Transverse Mercator (UTM/Gauss-Krueger)
- (Normal) Mercator
- Lambert Conformal Conic 1SP/2SP
- Equidistant Conic
- Albers Equal Area Conic

## Example

```go
package main

import (
	"fmt"

	"github.com/wroge/wgs84"
	"github.com/wroge/wgs84/etrs89"
)

func main() {
	longitude := 9.0
	latitude := 52.0
	height := 0
	
	conversion := wgs84.ToWebMercator().Round(0)
	east, north, height := conversion(longitude, latitude, height)
	// 1001875.000000 6800125.000000 0.000000

	transformation := wgs84.To(etrs89.UTM(32))
	east, north, height = transformation(longitude, latitude, height)
	// 500000.000000 5761038.213044 0.000065
}
```
