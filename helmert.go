package wgs84

import "math"

// Helmert implements the Transformation interface.
type Helmert struct {
	Tx, Ty, Tz, Rx, Ry, Rz, Ds float64
}

// ToWGS84 transforms geocentric coordinates to WGS84 geocentric
// coordinates.
func (t Helmert) ToWGS84(x, y, z float64) (x0, y0, z0 float64) {
	return helmert(x, y, z, t.Tx, t.Ty, t.Tz, t.Rx, t.Ry, t.Rz, t.Ds)
}

// FromWGS84 transforms WGS84 geocentric coordinates to geocentric
// coordinates.
func (t Helmert) FromWGS84(x0, y0, z0 float64) (x, y, z float64) {
	return helmert(x0, y0, z0, -t.Tx, -t.Ty, -t.Tz, -t.Rx, -t.Ry, -t.Rz, -t.Ds)
}

func helmert(x, y, z, tx, ty, tz, rx, ry, rz, ds float64) (x0, y0, z0 float64) {
	asec := math.Pi / 648000
	ppm := 0.000001
	x0 = (1+ds*ppm)*(x+z*ry*asec-y*rz*asec) + tx
	y0 = (1+ds*ppm)*(y+x*rz*asec-z*rx*asec) + ty
	z0 = (1+ds*ppm)*(z+y*rx*asec-x*ry*asec) + tz
	return
}
