[![GoDoc](http://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/github.com/wroge/wgs84)
[![Go Report Card](https://goreportcard.com/badge/github.com/wroge/wgs84?style=flat-square)](https://goreportcard.com/report/github.com/wroge/wgs84)
[![GolangCI](https://golangci.com/badges/github.com/wroge/wgs84.svg)](https://golangci.com/r/github.com/wroge/wgs84)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/wroge/wgs84.svg?style=social)](https://github.com/wroge/wgs84/tags)
![gopherbadger-tag-do-not-edit]()

# WGS84

A pure Go package for coordinate transformations.

```go
go get github.com/wroge/wgs84
```

### Usage
```go
east, north, h := wgs84.LonLat{}.To(wgs84.ETRS89{}.UTM(32)).Round(2)(9, 52, 0)
// 500000 5.76103821e+06 0

lon, lat, h := wgs84.ETRS89{}.UTM(32).To(wgs84.EPSG().Code(4326)).Round(3)(500150, 5761200, 0)
// 9.002 52.001 0

// EPSG-Codes covering the coordinate {longitude: 9, latitude: 52}:
codes := wgs84.EPSG().CodesCover(9, 52)
// [4326 25832 900913 4258 4314 4978 32632 3857 31467]
```

[...Calculate EPSG-Code from Unknown Coordinates](https://gist.github.com/wroge/e2160c1483a083997accf49009e7b08a)   
[...Calculate WebMercator Tile from WGS84 Longitude Latitude](https://gist.github.com/wroge/979869ff59046c4d841248c101472783)   
[...Transformation between OSGB36 NationalGrid and WGS84 Geographic Coordinates](https://gist.github.com/wroge/b7cd3c9dda9973b7085a10b09360ea00)   
[...Adding a CoordinateReferenceSystem (MGI AustriaLambert) to the EPSG-Repository](https://gist.github.com/wroge/844743b2756dcb47077eacbf2f129b92)   

### Features

- Helmert Transformation
- Web Mercator
- Lambert Conformal Conic
- Transverse Mercator (UTM)
- EPSG-Code Coverage
- ...
- Easily expandable through simple interfaces
