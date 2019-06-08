# EPSG

Store and search 
[WGS84](https://github.com/wroge/wgs84) CoordinateReferenceSystem`s as
[EPSG](https://en.wikipedia.org/wiki/Spatial_reference_system#Identifiers) Codes.

[![GoDoc](https://godoc.org/github.com/wroge/wgs84/epsg?status.svg)](https://godoc.org/github.com/wroge/wgs84/epsg)

```go
import "github.com/wroge/wgs84/epsg"

EPSG-Codes covering a specific Lon/Lat WGS84 coordinate:
epsg.Codes(9, 52)
// [3857 31467 900913 32632 4978 4326]

Get a wgs84.CoordinateReferenceSystem:
crs, ok := epsg.Code(4326)

Add an EPSG-Code:
r := epsg.DefaultRepository()
r.Add(...)
crs, ok := r.Code(...)
```

Back to [WGS84](https://github.com/wroge/wgs84).