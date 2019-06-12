# EPSG

Store and search 
[WGS84](https://github.com/wroge/wgs84) CoordinateReferenceSystem`s as
[EPSG](https://en.wikipedia.org/wiki/Spatial_reference_system#Identifiers) Codes.

[![GoDoc](https://godoc.org/github.com/wroge/wgs84/epsg?status.svg)](https://godoc.org/github.com/wroge/wgs84/epsg)

```go
import "github.com/wroge/wgs84/epsg"

coordinate transformation:
t := epsg.Transform(4326, 25832)
if t != nil {
    fmt.Println(t(9, 52, 0))
}

EPSG-Codes covering a specific Lon/Lat WGS84 coordinate:
epsg.Codes(9, 52)
// [4978 3857 4326 31467 25832 900913 32632]

Get a wgs84.CoordinateReferenceSystem:
crs, ok := epsg.Code(4326)

Add an EPSG-Code:
r := epsg.DefaultRepository()
r.Add(...)
crs, ok := r.Code(...)
```

Back to [WGS84](https://github.com/wroge/wgs84).
