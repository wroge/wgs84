package wgs84

import (
	"github.com/wroge/wgs84/spheroid"
	"github.com/wroge/wgs84/system"
)

// By default handled as WGS84
type GeodeticDatum struct {
	Spheroid       Spheroid
	Transformation Transformation
}

func (d GeodeticDatum) XYZ() CoordinateReferenceSystem {
	if d.Spheroid == nil {
		d.Spheroid = spheroid.WGS84()
	}
	return CoordinateReferenceSystem{
		Spheroid:       d.Spheroid,
		Transformation: d.Transformation,
	}
}

func (d GeodeticDatum) LonLat() CoordinateReferenceSystem {
	if d.Spheroid == nil {
		d.Spheroid = spheroid.WGS84()
	}
	return CoordinateReferenceSystem{
		Spheroid:       d.Spheroid,
		Transformation: d.Transformation,
		System:         system.LonLat{},
	}
}

func (d GeodeticDatum) WebMercator() CoordinateReferenceSystem {
	if d.Spheroid == nil {
		d.Spheroid = spheroid.WGS84()
	}
	return CoordinateReferenceSystem{
		Spheroid:       d.Spheroid,
		Transformation: d.Transformation,
		System:         system.WebMercator{},
	}
}

func (d GeodeticDatum) TransverseMercator(lonf, latf, scale, eastf, northf float64) CoordinateReferenceSystem {
	if d.Spheroid == nil {
		d.Spheroid = spheroid.WGS84()
	}
	return CoordinateReferenceSystem{
		Spheroid:       d.Spheroid,
		Transformation: d.Transformation,
		System: system.TransverseMercator{
			Lonf:   lonf,
			Latf:   latf,
			Scale:  scale,
			Eastf:  eastf,
			Northf: northf,
		},
	}
}

func (d GeodeticDatum) UTM(zone float64, northern bool) CoordinateReferenceSystem {
	if d.Spheroid == nil {
		d.Spheroid = spheroid.WGS84()
	}
	return CoordinateReferenceSystem{
		Spheroid:       d.Spheroid,
		Transformation: d.Transformation,
		System:         system.UTM(zone, northern),
	}
}

func (d GeodeticDatum) GK(zone float64) CoordinateReferenceSystem {
	if d.Spheroid == nil {
		d.Spheroid = spheroid.WGS84()
	}
	return CoordinateReferenceSystem{
		Spheroid:       d.Spheroid,
		Transformation: d.Transformation,
		System:         system.GK(zone),
	}
}

func (d GeodeticDatum) Mercator(lonf, scale, eastf, northf float64) CoordinateReferenceSystem {
	if d.Spheroid == nil {
		d.Spheroid = spheroid.WGS84()
	}
	return CoordinateReferenceSystem{
		Spheroid:       d.Spheroid,
		Transformation: d.Transformation,
		System: system.Mercator{
			Lonf:   lonf,
			Scale:  scale,
			Eastf:  eastf,
			Northf: northf,
		},
	}
}

func (d GeodeticDatum) LambertConformalConic1SP(lonf, latf, scale, eastf, northf float64) CoordinateReferenceSystem {
	if d.Spheroid == nil {
		d.Spheroid = spheroid.WGS84()
	}
	return CoordinateReferenceSystem{
		Spheroid:       d.Spheroid,
		Transformation: d.Transformation,
		System: system.LambertConformalConic1SP{
			Lonf:   lonf,
			Latf:   latf,
			Scale:  scale,
			Eastf:  eastf,
			Northf: northf,
		},
	}
}

func (d GeodeticDatum) LambertConformalConic2SP(lonf, latf, lat1, lat2, eastf, northf float64) CoordinateReferenceSystem {
	if d.Spheroid == nil {
		d.Spheroid = spheroid.WGS84()
	}
	return CoordinateReferenceSystem{
		Spheroid:       d.Spheroid,
		Transformation: d.Transformation,
		System: system.LambertConformalConic2SP{
			Lonf:   lonf,
			Latf:   latf,
			Lat1:   lat1,
			Lat2:   lat2,
			Eastf:  eastf,
			Northf: northf,
		},
	}
}

func (d GeodeticDatum) AlbersEqualAreaConic(lonf, latf, lat1, lat2, eastf, northf float64) CoordinateReferenceSystem {
	if d.Spheroid == nil {
		d.Spheroid = spheroid.WGS84()
	}
	return CoordinateReferenceSystem{
		Spheroid:       d.Spheroid,
		Transformation: d.Transformation,
		System: system.AlbersEqualAreaConic{
			Lonf:   lonf,
			Latf:   latf,
			Lat1:   lat1,
			Lat2:   lat2,
			Eastf:  eastf,
			Northf: northf,
		},
	}
}

func (d GeodeticDatum) EquidistantConic(lonf, latf, lat1, lat2, eastf, northf float64) CoordinateReferenceSystem {
	if d.Spheroid == nil {
		d.Spheroid = spheroid.WGS84()
	}
	return CoordinateReferenceSystem{
		Spheroid:       d.Spheroid,
		Transformation: d.Transformation,
		System: system.EquidistantConic{
			Lonf:   lonf,
			Latf:   latf,
			Lat1:   lat1,
			Lat2:   lat2,
			Eastf:  eastf,
			Northf: northf,
		},
	}
}
