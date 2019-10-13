package wgs84

// Transform is the core utility of this package. It is used by the To
// methods of each CoordinateReferenceSystem in this package to transform
// coordinates from one CoordinateReferenceSystem to another.
func Transform(from, to CoordinateReferenceSystem) func(a, b, c float64) (a2, b2, c2 float64) {
	return func(a, b, c float64) (a2, b2, c2 float64) {
		x, y, z := from.ToXYZ(a, b, c, from)
		x0, y0, z0 := from.ToWGS84(x, y, z)
		x2, y2, z2 := to.FromWGS84(x0, y0, z0)
		return to.FromXYZ(x2, y2, z2, to)
	}
}
