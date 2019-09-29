package wgs84

type ETRS89 struct{}

func (ETRS89) Datum() GeodeticDatum {
	return NewDatum(6378137, 298.257222101)
}

func (e ETRS89) LonLat() CoordinateReferenceSystem {
	return e.Datum().LonLat()
}

func (e ETRS89) UTM(zone float64) CoordinateReferenceSystem {
	return e.Datum().UTM(zone, true)
}
