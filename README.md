[![GoDoc](http://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/github.com/wroge/wgs84)
[![GoWalker](https://img.shields.io/badge/Go_Walker-Doc-blue.svg?style=flat-square)](https://gowalker.org/github.com/wroge/wgs84)
[![Go Report Card](https://goreportcard.com/badge/github.com/wroge/wgs84?style=flat-square)](https://goreportcard.com/report/github.com/wroge/wgs84)
[![GolangCI](https://golangci.com/badges/github.com/wroge/wgs84.svg)](https://golangci.com/r/github.com/wroge/wgs84)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/wroge/wgs84.svg?style=social)](https://github.com/wroge/wgs84/tags)
# WGS84

A pure Go package for coordinate transformations.

```go
go get github.com/wroge/wgs84
```

### Usage
```go
east, north, h := wgs84.LonLat{}.To(wgs84.ETRS89{}.UTM(32))(9, 52, 0)
// 500000 5.76103821304426e+06 6.50063157081604e-05

lon, lat, h := wgs84.ETRS89{}.UTM(32).To(wgs84.EPSG().Code(4326))(500000, 5761038, 0)
// 9 51.99999808768612 -6.501004099845886e-05

codes := wgs84.EPSG().CodesContain(9, 52)
// [31468 31469 4978 900913 25832 32632 3857 31466 4326 31467 32732]
```

### Features

- Helmert Transformation
- Web Mercator
- Lambert Conformal Conic
- Transverse Mercator (UTM)
- EPSG-Code Coverage
- ...
- Easily expandable through simple interfaces
