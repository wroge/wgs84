package wgs84

func RGF93() rgf93 {
	return rgf93{}
}

type rgf93 struct{}

func (rgf93) Datum() GeodeticDatum {
	return ETRS89().Datum()
}

func (r rgf93) FranceLambert() CoordinateReferenceSystem {
	return CoordinateReferenceSystem{
		GeodeticDatum:    r.Datum(),
		CoordinateSystem: lambertConformalConic2SP(3, 46.5, 49, 44, 700000, 6600000),
	}
}

func (r rgf93) CC(lat float64) CoordinateReferenceSystem {
	return CoordinateReferenceSystem{
		GeodeticDatum:    r.Datum(),
		CoordinateSystem: lambertConformalConic2SP(3, lat, lat-0.75, lat+0.75, 1700000, 2200000+(lat-43)*1000000),
	}
}
