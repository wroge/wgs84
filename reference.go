package wgs84

type CoordinateReferenceSystem struct {
	GeodeticDatum    GeodeticDatum
	CoordinateSystem CoordinateSystem
}

func (crs CoordinateReferenceSystem) To(crs2 CoordinateReferenceSystem) func(a, b, c float64) (a2, b2, c2 float64) {
	if crs.GeodeticDatum.Spheroid == nil {
		crs.GeodeticDatum.Spheroid = Datum().Spheroid
	}
	if crs2.GeodeticDatum.Spheroid == nil {
		crs2.GeodeticDatum.Spheroid = Datum().Spheroid
	}
	return func(a, b, c float64) (a2, b2, c2 float64) {
		if crs.CoordinateSystem != nil {
			a, b, c = crs.CoordinateSystem.ToXYZ(a, b, c, crs.GeodeticDatum.Spheroid)
		}
		if crs.GeodeticDatum.Transformation != nil {
			a, b, c = crs.GeodeticDatum.Transformation.ToWGS84(a, b, c)
		}
		if crs2.GeodeticDatum.Transformation != nil {
			a, b, c = crs2.GeodeticDatum.Transformation.FromWGS84(a, b, c)

		}
		if crs2.CoordinateSystem != nil {
			a, b, c = crs2.CoordinateSystem.FromXYZ(a, b, c, crs2.GeodeticDatum.Spheroid)
		}
		return a, b, c
	}
}
