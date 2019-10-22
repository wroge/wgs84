package wgs84

import (
	"math"
)

func Helmert(a, fi, tx, ty, tz, rx, ry, rz, ds float64) Datum {
	return Datum{
		Spheroid: spheroid{a: a, fi: fi},
		Transformation: helmert{
			tx: tx,
			ty: ty,
			tz: tz,
			rx: rx,
			ry: ry,
			rz: rz,
			ds: ds,
		},
	}
}

func WGS84() Datum {
	return Datum{
		Spheroid: spheroid{
			a:  6378137,
			fi: 298.257223563,
		},
		Area: AreaFunc(func(lon, lat float64) bool {
			if math.Abs(lon) > 180 || math.Abs(lat) > 90 {
				return false
			}
			return true
		}),
	}
}

func ETRS89() Datum {
	return Datum{
		Spheroid: spheroid{
			a:  6378137,
			fi: 298.257222101,
		},
		Area: AreaFunc(func(lon, lat float64) bool {
			if lon < -16.1 || lon > 40.18 || lat < 32.88 || lat > 84.17 {
				return false
			}
			return true
		}),
	}
}

func OSGB36() Datum {
	return Datum{
		Spheroid: spheroid{
			a:  6377563.396,
			fi: 299.3249646,
		},
		Transformation: helmert{
			tx: 446.448,
			ty: -125.157,
			tz: 542.06,
			rx: 0.15,
			ry: 0.247,
			rz: 0.842,
			ds: -20.489,
		},
		Area: AreaFunc(func(lon, lat float64) bool {
			if lon < -8.82 || lon > 1.92 || lat < 49.79 || lat > 60.94 {
				return false
			}
			return true
		}),
	}
}

func DHDN2001() Datum {
	return Datum{
		Spheroid: spheroid{
			a:  6377397.155,
			fi: 299.1528128,
		},
		Transformation: helmert{
			tx: 598.1,
			ty: 73.7,
			tz: 418.2,
			rx: 0.202,
			ry: 0.045,
			rz: -2.455,
			ds: 6.7,
		},
		Area: AreaFunc(func(lon, lat float64) bool {
			if lon < 5.87 || lon > 13.84 || lat < 47.27 || lat > 55.09 {
				return false
			}
			return true
		}),
	}
}

func RGF93() Datum {
	return Datum{
		Spheroid: spheroid{
			a:  6378137,
			fi: 298.257222101,
		},
		Area: AreaFunc(func(lon, lat float64) bool {
			if lon < -9.86 || lon > 10.38 || lat < 41.15 || lat > 51.56 {
				return false
			}
			return true
		}),
	}
}

func NAD83() Datum {
	return Datum{
		Spheroid: spheroid{
			a:  6378137,
			fi: 298.257222101,
		},
		Area: AreaFunc(func(lon, lat float64) bool {
			if lon < -172.54 || lon > -47.74 || lat < 23.81 || lat > 86.46 {
				return false
			}
			return true
		}),
	}
}

type Datum struct {
	Spheroid       Spheroid
	Transformation Transformation
	Area           Area
}

func (d Datum) Contains(lon, lat float64) bool {
	if math.Abs(lon) > 180 || math.Abs(lat) > 90 {
		return false
	}
	if d.Area != nil {
		return d.Area.Contains(lon, lat)
	}
	return true
}

func (d Datum) A() float64 {
	if d.Spheroid == nil {
		return 6378137
	}
	return d.Spheroid.A()
}

func (d Datum) Fi() float64 {
	if d.Spheroid == nil {
		return 298.257223563
	}
	return d.Spheroid.Fi()
}

func (d Datum) Forward(x, y, z float64) (x0, y0, z0 float64) {
	if d.Transformation == nil {
		return x, y, z
	}
	return d.Transformation.Forward(x, y, z)
}

func (d Datum) Inverse(x0, y0, z0 float64) (x, y, z float64) {
	if d.Transformation == nil {
		return x0, y0, z0
	}
	return d.Transformation.Inverse(x0, y0, z0)
}

func (d Datum) XYZ() GeocentricReferenceSystem {
	return GeocentricReferenceSystem{
		Datum: d,
	}
}

func (d Datum) LonLat() GeographicReferenceSystem {
	return GeographicReferenceSystem{
		Datum: d,
	}
}

func (d Datum) WebMercator() ProjectedReferenceSystem {
	return ProjectedReferenceSystem{
		Datum:      d,
		Projection: webMercator{},
	}
}

func (d Datum) TransverseMercator(lonf, latf, scale, eastf, northf float64) ProjectedReferenceSystem {
	return ProjectedReferenceSystem{
		Datum: d,
		Projection: transverseMercator{
			lonf:   lonf,
			latf:   latf,
			scale:  scale,
			eastf:  eastf,
			northf: northf,
		},
	}
}

func (d Datum) LambertConformalConic2SP(lonf, latf, lat1, lat2, eastf, northf float64) ProjectedReferenceSystem {
	return ProjectedReferenceSystem{
		Datum: d,
		Projection: lambertConformalConic2SP{
			lonf:   lonf,
			latf:   latf,
			lat1:   lat1,
			lat2:   lat2,
			eastf:  eastf,
			northf: northf,
		},
	}
}

func (d Datum) AlbersEqualAreaConic(lonf, latf, lat1, lat2, eastf, northf float64) ProjectedReferenceSystem {
	return ProjectedReferenceSystem{
		Datum: d,
		Projection: albersEqualAreaConic{
			lonf:   lonf,
			latf:   latf,
			lat1:   lat1,
			lat2:   lat2,
			eastf:  eastf,
			northf: northf,
		},
	}
}
