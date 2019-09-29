package wgs84

func EPSG(code int) CoordinateReferenceSystem {
	switch code {
	case 4326:
		return LonLat()
	case 4978:
		return Datum().XYZ()
	case 3857, 900913:
		return WebMercator()
	case 27700:
		return OSGB36{}.NationalGrid()
	case 4277:
		return OSGB36{}.LonLat()
	case 2154:
		return RGF93{}.FranceLambert()
	}
	if z := code - 32600; z >= 1 && z <= 60 {
		return Datum().UTM(float64(z), true)
	}
	if z := code - 32700; z >= 1 && z <= 60 {
		return Datum().UTM(float64(z), false)
	}
	if z := code - 31464; z >= 2 && z <= 5 {
		return DHDN2001{}.GK(float64(z))
	}
	if z := code - 25800; z >= 28 && z <= 38 {
		return ETRS89{}.UTM(float64(z))
	}
	if lat := code - 3900; lat >= 42 && lat <= 50 {
		return RGF93{}.CC(float64(lat))
	}
	return CoordinateReferenceSystem{}
}
