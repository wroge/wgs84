package wgs84

func NewDatum(a, fi float64) GeodeticDatum {
	return GeodeticDatum{
		Spheroid: spheroid{a, fi},
	}
}

type GeodeticDatum struct {
	Spheroid       Spheroid
	Transformation Transformation
}

func (gd GeodeticDatum) GeocentricTranslation(tx, ty, tz float64) GeodeticDatum {
	return GeodeticDatum{
		Spheroid:       gd.Spheroid,
		Transformation: geocentricTranslation(gd.Transformation, tx, ty, tz),
	}
}

func (gd GeodeticDatum) Helmert(tx, ty, tz, rx, ry, rz, ds float64) GeodeticDatum {
	return GeodeticDatum{
		Spheroid:       gd.Spheroid,
		Transformation: helmert(gd.Transformation, tx, ty, tz, rx, ry, rz, ds),
	}
}

func (gd GeodeticDatum) XYZ() CoordinateReferenceSystem {
	return CoordinateReferenceSystem{
		GeodeticDatum:    gd,
		CoordinateSystem: nil,
	}
}

func (gd GeodeticDatum) LonLat() CoordinateReferenceSystem {
	return CoordinateReferenceSystem{
		GeodeticDatum:    gd,
		CoordinateSystem: lonLat(),
	}
}

func (gd GeodeticDatum) WebMercator() CoordinateReferenceSystem {
	return CoordinateReferenceSystem{
		GeodeticDatum:    gd,
		CoordinateSystem: webMercator(),
	}
}

func (gd GeodeticDatum) TransverseMercator(lonf, latf, scale, eastf, northf float64) CoordinateReferenceSystem {
	return CoordinateReferenceSystem{
		GeodeticDatum:    gd,
		CoordinateSystem: transverseMercator(lonf, latf, scale, eastf, northf),
	}
}

func (gd GeodeticDatum) UTM(zone float64, northern bool) CoordinateReferenceSystem {
	return CoordinateReferenceSystem{
		GeodeticDatum:    gd,
		CoordinateSystem: utm(zone, northern),
	}
}

func (gd GeodeticDatum) GK(zone float64) CoordinateReferenceSystem {
	return CoordinateReferenceSystem{
		GeodeticDatum:    gd,
		CoordinateSystem: gk(zone),
	}
}

func (gd GeodeticDatum) Mercator(lonf, scale, eastf, northf float64) CoordinateReferenceSystem {
	return CoordinateReferenceSystem{
		GeodeticDatum:    gd,
		CoordinateSystem: mercator(lonf, scale, eastf, northf),
	}
}

func (gd GeodeticDatum) LambertConformalConic1SP(lonf, latf, scale, eastf, northf float64) CoordinateReferenceSystem {
	return CoordinateReferenceSystem{
		GeodeticDatum:    gd,
		CoordinateSystem: lambertConformalConic1SP(lonf, latf, scale, eastf, northf),
	}
}

func (gd GeodeticDatum) LambertConformalConic2SP(lonf, latf, lat1, lat2, eastf, northf float64) CoordinateReferenceSystem {
	return CoordinateReferenceSystem{
		GeodeticDatum:    gd,
		CoordinateSystem: lambertConformalConic2SP(lonf, latf, lat1, lat2, eastf, northf),
	}
}

func (gd GeodeticDatum) AlbersEqualAreaConic(lonf, latf, lat1, lat2, eastf, northf float64) CoordinateReferenceSystem {
	return CoordinateReferenceSystem{
		GeodeticDatum:    gd,
		CoordinateSystem: albersEqualAreaConic(lonf, latf, lat1, lat2, eastf, northf),
	}
}

func (gd GeodeticDatum) EquidistantConic(lonf, latf, lat1, lat2, eastf, northf float64) CoordinateReferenceSystem {
	return CoordinateReferenceSystem{
		GeodeticDatum:    gd,
		CoordinateSystem: equidistantConic(lonf, latf, lat1, lat2, eastf, northf),
	}
}
