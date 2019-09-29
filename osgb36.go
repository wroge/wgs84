package wgs84

type OSGB36 struct{}

func (OSGB36) Datum() GeodeticDatum {
	return NewDatum(6377563.396, 299.3249646).Helmert(446.448, -125.157, 542.06, 0.15, 0.247, 0.842, -20.489)
}

func (o OSGB36) LonLat() CoordinateReferenceSystem {
	return o.Datum().LonLat()
}

func (o OSGB36) NationalGrid() CoordinateReferenceSystem {
	return o.Datum().TransverseMercator(-2, 49, 0.9996012717, 400000, -100000)
}
