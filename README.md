# wgs84

[![Build Status](https://img.shields.io/travis/wroge/wgs84/master)](https://travis-ci.org/wroge/wgs84)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/wroge/wgs84)
[![Go Report Card](https://goreportcard.com/badge/github.com/wroge/wgs84)](https://goreportcard.com/report/github.com/wroge/wgs84)
![golangci-lint](https://github.com/wroge/wgs84/workflows/golangci-lint/badge.svg)
[![Codecov](https://img.shields.io/codecov/c/gh/wroge/wgs84)](https://codecov.io/gh/wroge/wgs84)
[![tippin.me](https://badgen.net/badge/%E2%9A%A1%EF%B8%8Ftippin.me/@_wroge/yellow)](https://tippin.me/@_wroge)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/wroge/wgs84.svg?style=social)](https://github.com/wroge/wgs84/tags)

A *zero-dependency* Go package for coordinate transformations.

```go
go get github.com/wroge/wgs84
```

## Usage

```go
east, north, h := wgs84.LonLat().To(wgs84.ETRS89UTM(32)).Round(2)(9, 52, 0)
// 500000 5.76103821e+06 0

east, north, h := wgs84.To(wgs84.WebMercator())(9, 52, 0)
// 1.0018754171394621e+06 6.800125454397305e+06 -9.313225746154785e-10

epsg := wgs84.EPSG()

lon, lat, h := wgs84.ETRS89UTM(32).To(epsg.Code(4326)).Round(3)(500150, 5761200, 0)
// 9.002 52.001 0

// EPSG-Codes covering the coordinate {longitude: 9, latitude: 52}:
codes := epsg.CodesCover(9, 52)
// [25832 4314 32632 4978 4258 31467 4326 3857 900913]
```

- [Calculate EPSG-Code from Unknown Coordinates](https://gist.github.com/wroge/e2160c1483a083997accf49009e7b08a)
- [Calculate WebMercator Tile from WGS84 Longitude Latitude](https://gist.github.com/wroge/979869ff59046c4d841248c101472783)
- [Transformation between OSGB36 NationalGrid and WGS84 Geographic Coordinates](https://gist.github.com/wroge/b7cd3c9dda9973b7085a10b09360ea00)
- [Adding a CoordinateReferenceSystem (MGI AustriaLambert) to the EPSG-Repository](https://gist.github.com/wroge/844743b2756dcb47077eacbf2f129b92)

## Features

- Helmert Transformation
- Web Mercator
- Lambert Conformal Conic
- Transverse Mercator (UTM)
- EPSG-Code Coverage
- ...
- Easily expandable through simple [Interfaces](https://github.com/wroge/wgs84/blob/master/interface.go)
