package wgs84

type Datum struct {
	GeodeticSpheroid GeodeticSpheroid
	Transformation   Transformation
}

func (d Datum) MajorAxis() float64 {
	return spheroid(d.GeodeticSpheroid).MajorAxis()
}

func (d Datum) InverseFlattening() float64 {
	return spheroid(d.GeodeticSpheroid).InverseFlattening()
}

func (d Datum) ToWGS84(x, y, z float64) (x0, y0, z0 float64) {
	return toWGS84(d.Transformation, x, y, z)
}

func (d Datum) FromWGS84(x0, y0, z0 float64) (x, y, z float64) {
	return fromWGS84(d.Transformation, x0, y0, z0)
}

type ETRS89 struct{}

func (d ETRS89) Spheroid() Spheroid {
	return GRS80()
}

func (d ETRS89) MajorAxis() float64 {
	return d.Spheroid().MajorAxis()
}

func (d ETRS89) InverseFlattening() float64 {
	return d.Spheroid().InverseFlattening()
}

func (d ETRS89) ToWGS84(x, y, z float64) (x0, y0, z0 float64) {
	return x, y, z
}

func (d ETRS89) FromWGS84(x0, y0, z0 float64) (x, y, z float64) {
	return x0, y0, z0
}

func (d ETRS89) LonLat() LonLat {
	return LonLat{
		GeodeticDatum: d,
	}
}

func (d ETRS89) UTM(zone float64) TransverseMercator {
	return TransverseMercator{
		GeodeticDatum: d,
		Lonf:          zone*6 - 183,
		Latf:          0,
		Scale:         0.9996,
		Eastf:         500000,
		Northf:        0,
	}
}

type OSGB36 struct{}

func (d OSGB36) Spheroid() Spheroid {
	return Airy()
}

func (d OSGB36) Helmert() Helmert {
	return Helmert{
		Tx: 446.448,
		Ty: -125.157,
		Tz: 542.06,
		Rx: 0.15,
		Ry: 0.247,
		Rz: 0.842,
		Ds: -20.489,
	}
}

func (d OSGB36) MajorAxis() float64 {
	return d.Spheroid().MajorAxis()
}

func (d OSGB36) InverseFlattening() float64 {
	return d.Spheroid().InverseFlattening()
}

func (d OSGB36) ToWGS84(x, y, z float64) (x0, y0, z0 float64) {
	return d.Helmert().ToWGS84(x, y, z)
}

func (d OSGB36) FromWGS84(x0, y0, z0 float64) (x, y, z float64) {
	return d.Helmert().FromWGS84(x0, y0, z0)
}

func (d OSGB36) LonLat() LonLat {
	return LonLat{
		GeodeticDatum: d,
	}
}

func (d OSGB36) NationalGrid() TransverseMercator {
	return TransverseMercator{
		GeodeticDatum: d,
		Lonf:          -2,
		Latf:          49,
		Scale:         0.9996012717,
		Eastf:         400000,
		Northf:        -100000,
	}
}

type DHDN2001 struct{}

func (d DHDN2001) Spheroid() Spheroid {
	return Bessel()
}

func (d DHDN2001) Helmert() Helmert {
	return Helmert{
		Tx: 598.1,
		Ty: 73.7,
		Tz: 418.2,
		Rx: 0.202,
		Ry: 0.045,
		Rz: -2.455,
		Ds: 6.7,
	}
}

func (d DHDN2001) MajorAxis() float64 {
	return d.Spheroid().MajorAxis()
}

func (d DHDN2001) InverseFlattening() float64 {
	return d.Spheroid().InverseFlattening()
}

func (d DHDN2001) ToWGS84(x, y, z float64) (x0, y0, z0 float64) {
	return d.Helmert().ToWGS84(x, y, z)
}

func (d DHDN2001) FromWGS84(x0, y0, z0 float64) (x, y, z float64) {
	return d.Helmert().FromWGS84(x0, y0, z0)
}

func (d DHDN2001) LonLat() LonLat {
	return LonLat{
		GeodeticDatum: d,
	}
}

func (d DHDN2001) GK(zone float64) TransverseMercator {
	return TransverseMercator{
		GeodeticDatum: d,
		Lonf:          zone * 3,
		Latf:          0,
		Scale:         1,
		Eastf:         zone*1000000 + 500000,
		Northf:        0,
	}
}

type RGF93 struct{}

func (d RGF93) Spheroid() Spheroid {
	return GRS80()
}

func (d RGF93) MajorAxis() float64 {
	return d.Spheroid().MajorAxis()
}

func (d RGF93) InverseFlattening() float64 {
	return d.Spheroid().InverseFlattening()
}

func (d RGF93) ToWGS84(x, y, z float64) (x0, y0, z0 float64) {
	return x, y, z
}

func (d RGF93) FromWGS84(x0, y0, z0 float64) (x, y, z float64) {
	return x0, y0, z0
}

func (d RGF93) LonLat() LonLat {
	return LonLat{
		GeodeticDatum: d,
	}
}

func (d RGF93) CC(lat float64) LambertConformalConic2SP {
	return LambertConformalConic2SP{
		GeodeticDatum: d,
		Lonf:          3,
		Latf:          lat,
		Lat1:          lat - 0.75,
		Lat2:          lat + 0.75,
		Eastf:         1700000,
		Northf:        2200000 + (lat-43)*1000000,
	}
}

func (d RGF93) FranceLambert() LambertConformalConic2SP {
	return LambertConformalConic2SP{
		GeodeticDatum: d,
		Lonf:          3,
		Latf:          46.5,
		Lat1:          49,
		Lat2:          44,
		Eastf:         700000,
		Northf:        6600000,
	}
}
