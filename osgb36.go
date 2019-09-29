package wgs84

func OSGB36() osgb36 {
	return osgb36{}
}

type osgb36 struct{}

func (osgb36) Datum() GeodeticDatum {
	return NewDatum(6377563.396, 299.3249646).Helmert(446.448, -125.157, 542.06, 0.15, 0.247, 0.842, -20.489)
}

func (o osgb36) LonLat() CoordinateReferenceSystem {
	return o.Datum().LonLat()
}

func (o osgb36) NationalGrid() CoordinateReferenceSystem {
	return o.Datum().TransverseMercator(-2, 49, 0.9996012717, 400000, -100000)
}
