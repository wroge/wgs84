# (Coordinate)System

[![GoDoc](https://godoc.org/github.com/wroge/wgs84/system?status.svg)](https://godoc.org/github.com/wroge/wgs84/system)

This package supports ...

- Geocentric/Cartesian (XYZ)
- Geographic (LonLat)
- Web Mercator
- Transverse Mercator (UTM/Gauss-Krueger)
- (Normal) Mercator
- Lambert Conformal Conic 1SP/2SP
- Equidistant Conic
- Albers Equal Area Conic

...coordinate systems.

Implementing the following methods ...

```go
ToXYZ(a, b, c float64, sph spheroid.Spheroid) (x, y, z float64)
FromXYZ(x, y, z float64, sph spheroid.Spheroid) (a, b, c float64)
```

Back to [WGS84](https://github.com/wroge/wgs84).