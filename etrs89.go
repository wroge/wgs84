package wgs84

func ETRS89() etrs89 {
	return etrs89{}
}

type etrs89 struct{}

func (etrs89) Datum() GeodeticDatum {
	return NewDatum(6378137, 298.257222101)
}

func (e etrs89) LonLat() CoordinateReferenceSystem {
	return e.Datum().LonLat()
}

func (e etrs89) UTM(zone float64) CoordinateReferenceSystem {
	return e.Datum().UTM(zone, true)
}
