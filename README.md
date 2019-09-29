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
transform := wgs84.LonLat().To(wgs84.ETRS89{}.UTM(32))
east, north, h := transform(9, 52, 0)

transform = wgs84.EPSG(25832).To(wgs84.EPSG(4326))
lon, lat, h = transform(500000, 5761038, 0)

austria := wgs84.NewDatum(6377397.155, 299.1528128).Helmert(577.326, 90.129, 463.919, 5.137, 1.474, 5.297, 2.4232)
austria_lambert := austria.LambertConformalConic2SP(13.33333333333333, 74.5,49,46,400000,400000)
```

### Features

- Helmert Transformation
- Web Mercator
- Lambert Conformal Conic
- Transverse Mercator (UTM)
- ...
- Easily expandable through simple interfaces
