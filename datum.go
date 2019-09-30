package wgs84

import "github.com/wroge/wgs84/system"

type GeodeticDatum struct {
	Spheroid       Spheroid
	Transformation Transformation
}

func (d GeodeticDatum) XYZ() CoordinateReferenceSystem {
	return CoordinateReferenceSystem{
		Spheroid:       d.Spheroid,
		Transformation: d.Transformation,
	}
}

func (d GeodeticDatum) LonLat() CoordinateReferenceSystem {
	return CoordinateReferenceSystem{
		Spheroid:       d.Spheroid,
		Transformation: d.Transformation,
		System:         system.LonLat{},
	}
}

func (d GeodeticDatum) WebMercator() CoordinateReferenceSystem {
	return CoordinateReferenceSystem{
		Spheroid:       d.Spheroid,
		Transformation: d.Transformation,
		System:         system.WebMercator{},
	}
}

func (d GeodeticDatum) TransverseMercator(lonf, latf, scale, eastf, northf float64) CoordinateReferenceSystem {
	return CoordinateReferenceSystem{
		Spheroid:       d.Spheroid,
		Transformation: d.Transformation,
		System:         system.TransverseMercator{lonf, latf, scale, eastf, northf},
	}
}

func (d GeodeticDatum) UTM(zone float64, northern bool) CoordinateReferenceSystem {
	return CoordinateReferenceSystem{
		Spheroid:       d.Spheroid,
		Transformation: d.Transformation,
		System:         system.UTM(zone, northern),
	}
}

func (d GeodeticDatum) GK(zone float64) CoordinateReferenceSystem {
	return CoordinateReferenceSystem{
		Spheroid:       d.Spheroid,
		Transformation: d.Transformation,
		System:         system.GK(zone),
	}
}

func (d GeodeticDatum) Mercator(lonf, scale, eastf, northf float64) CoordinateReferenceSystem {
	return CoordinateReferenceSystem{
		Spheroid:       d.Spheroid,
		Transformation: d.Transformation,
		System:         system.Mercator{lonf, scale, eastf, northf},
	}
}

func (d GeodeticDatum) LambertConformalConic1SP(lonf, latf, scale, eastf, northf float64) CoordinateReferenceSystem {
	return CoordinateReferenceSystem{
		Spheroid:       d.Spheroid,
		Transformation: d.Transformation,
		System:         system.LambertConformalConic1SP{lonf, latf, scale, eastf, northf},
	}
}

func (d GeodeticDatum) LambertConformalConic2SP(lonf, latf, lat1, lat2, eastf, northf float64) CoordinateReferenceSystem {
	return CoordinateReferenceSystem{
		Spheroid:       d.Spheroid,
		Transformation: d.Transformation,
		System:         system.LambertConformalConic2SP{lonf, latf, lat1, lat2, eastf, northf},
	}
}

func (d GeodeticDatum) AlbersEqualAreaConic(lonf, latf, lat1, lat2, eastf, northf float64) CoordinateReferenceSystem {
	return CoordinateReferenceSystem{
		Spheroid:       d.Spheroid,
		Transformation: d.Transformation,
		System:         system.AlbersEqualAreaConic{lonf, latf, lat1, lat2, eastf, northf},
	}
}

func (d GeodeticDatum) EquidistantConic(lonf, latf, lat1, lat2, eastf, northf float64) CoordinateReferenceSystem {
	return CoordinateReferenceSystem{
		Spheroid:       d.Spheroid,
		Transformation: d.Transformation,
		System:         system.EquidistantConic{lonf, latf, lat1, lat2, eastf, northf},
	}
}
