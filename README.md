# WGS84

A pure Go package for coordinate conversion and transformation..

[![GoDoc](https://godoc.org/github.com/wroge/wgs84?status.svg)](https://godoc.org/github.com/wroge/wgs84)

```go
import "github.com/wroge/wgs84"

coordinate conversion:

f := wgs84.ToUTM(32, true)         // wgs84.LonLat().ToUTM(32, true)
east, north, h := f(8.5, 52.5, 0)
// 466058.576840 5816769.501816 0.000000

f = wgs84.WebMercator().ToXYZ().Round(3)
x, y, z := f(1002000, 6800100, 0)
// 3886514.741000 615641.281000 5002793.702000

coordinate transformation:
	
import "github.com/wroge/wgs84/etrs89"

f := wgs84.From(etrs89.UTM(32))
lon, lat, h := f(501000, 5760000, 0)
// 9.014563 51.990665 -0.000065
```

Example of the realization of a geodetic datum (OSGB36):

```go
func OSGB36(sys system.System) wgs84.CoordinateReferenceSystem {
	return wgs84.CoordinateReferenceSystem{
		Spheroid:       spheroid.Airy(),
		Transformation: transformation.WGS84().Helmert(446.448, -125.157, 542.06, 0.15, 0.247, 0.842, -20.489),
		System:         sys,
	}
}
```

Subpackages with predefined coordinate systems for a specific geodetic datum:

- [DHDN2001](https://github.com/wroge/wgs84/tree/master/dhdn2001)
- [ETRS89](https://github.com/wroge/wgs84/tree/master/etrs89)
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