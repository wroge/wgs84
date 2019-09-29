package wgs84

func Datum() GeodeticDatum {
	return NewDatum(6378137, 298.257223563)
}

func LonLat() CoordinateReferenceSystem {
	return CoordinateReferenceSystem{
		GeodeticDatum:    Datum(),
		CoordinateSystem: lonLat(),
	}
}

func WebMercator() CoordinateReferenceSystem {
	return CoordinateReferenceSystem{
		GeodeticDatum:    Datum(),
		CoordinateSystem: webMercator(),
	}
}
