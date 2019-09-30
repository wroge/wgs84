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
east, north, h := wgs84.LonLat().To(wgs84.ETRS89UTM(32))(9, 52, 0)

epsg := wgs84.EPSG()
lon, lat, h := epsg.Code(25832).To(epsg.Code(4326))(500000, 5761038, 0)
```

### Features

- Helmert Transformation
- Web Mercator
- Lambert Conformal Conic
- Transverse Mercator (UTM)
- ...
- Easily expandable through simple interfaces
