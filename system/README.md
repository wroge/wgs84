# (Coordinate)System

[![GoDoc](http://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/github.com/wroge/wgs84/system)
[![GoWalker](https://img.shields.io/badge/Go_Walker-Doc-blue.svg?style=flat-square)](https://gowalker.org/github.com/wroge/wgs84/system)

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
ToXYZ(a, b, c float64, sph Spheroid) (x, y, z float64)
FromXYZ(x, y, z float64, sph Spheroid) (a, b, c float64)
```

Back to [WGS84](https://github.com/wroge/wgs84).