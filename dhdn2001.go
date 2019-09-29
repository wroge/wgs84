package wgs84

type DHDN2001 struct{}

func (DHDN2001) Datum() GeodeticDatum {
	return NewDatum(6377397.155, 299.1528128).Helmert(598.1, 73.7, 418.2, 0.202, 0.045, -2.455, 6.7)
}

func (d DHDN2001) LonLat() CoordinateReferenceSystem {
	return CoordinateReferenceSystem{
		GeodeticDatum:    d.Datum(),
		CoordinateSystem: lonLat(),
	}
}

func (d DHDN2001) GK(zone float64) CoordinateReferenceSystem {
	return CoordinateReferenceSystem{
		GeodeticDatum:    d.Datum(),
		CoordinateSystem: gk(zone),
	}
}
