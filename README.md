# WGS84

A pure Go package for coordinate transformations.

```go
go get github.com/wroge/wgs84
```

### Usage
```go
transform := wgs84.LonLat().To(wgs84.ETRS89().UTM(32))
east, north, h := transform(9, 52, 0)

transform = wgs84.EPSG(25832).To(wgs84.EPSG(4326))
lon, lat, h = transform(500000, 5761038, 0)

austria := wgs84.NewDatum(6377397.155, 299.1528128).Helmert(577.326, 90.129, 463.919, 5.137, 1.474, 5.297, 2.4232)
austria_lambert := austria.LambertConformalConic2SP(13.33333333333333, 74.5,49,46,400000,400000)
```

### Features

- ...
- Helmert Transformation
- Geocentric Translation
- Web Mercator
- Lambert Conformal Conic
- Transverse Mercator (UTM)
- Albers Equal Area Conic
- ...

### Data-Structure

```go
type Spheroid interface {
	A() float64
	Fi() float64
}

type Transformation interface {
	ToWGS84(x, y, z float64) (x0, y0, z0 float64)
	FromWGS84(x0, y0, z0 float64) (x, y, z float64)
}

type GeodeticDatum struct {
	Spheroid       Spheroid
	Transformation Transformation
}

type CoordinateSystem interface {
	ToXYZ(a, b, c float64, s Spheroid) (x, y, z float64)
	FromXYZ(x, y, z float64, s Spheroid) (a, b, c float64)
}

type CoordinateReferenceSystem struct {
	GeodeticDatum    GeodeticDatum
	CoordinateSystem CoordinateSystem
}
```
