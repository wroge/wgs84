package wgs84

// Transform is the core utility of this package. It is used by the To
// methods of each CoordinateReferenceSystem in this package to transform
// coordinates from one CoordinateReferenceSystem to another.
func Transform(from, to CoordinateReferenceSystem) func(a, b, c float64) (a2, b2, c2 float64) {
	return func(a, b, c float64) (a2, b2, c2 float64) {
		s := spheroid(from)
		if from != nil {
			a, b, c = from.ToXYZ(a, b, c, s)
			a, b, c = from.ToWGS84(a, b, c)
		}
		s = spheroid(to, from)
		if to != nil {
			a, b, c = to.FromWGS84(a, b, c)
			a, b, c = to.FromXYZ(a, b, c, s)
		}
		return a, b, c
	}
}
