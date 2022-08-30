//nolint:nonamedreturns
package wgs84

import "math"

const (
	asec = math.Pi / 648000
	ppm  = 0.000001
)

type helmert struct {
	tx, ty, tz, rx, ry, rz, ds float64
}

func (t helmert) Forward(x, y, z float64) (x0, y0, z0 float64) {
	return calcHelmert(x, y, z, t.tx, t.ty, t.tz, t.rx, t.ry, t.rz, t.ds)
}

func (t helmert) Inverse(x0, y0, z0 float64) (x, y, z float64) {
	return calcHelmert(x0, y0, z0, -t.tx, -t.ty, -t.tz, -t.rx, -t.ry, -t.rz, -t.ds)
}

func calcHelmert(x, y, z, tx, ty, tz, rx, ry, rz, ds float64) (x0, y0, z0 float64) {
	x0 = (1+ds*ppm)*(x+z*ry*asec-y*rz*asec) + tx
	y0 = (1+ds*ppm)*(y+x*rz*asec-z*rx*asec) + ty
	z0 = (1+ds*ppm)*(z+y*rx*asec-x*ry*asec) + tz

	return
}
