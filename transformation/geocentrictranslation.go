package transformation

type GeocentricTranslation struct {
	Tx, Ty, Tz float64
}

func (t GeocentricTranslation) ToWGS84(x, y, z float64) (x0, y0, z0 float64) {
	return x + t.Tx, y + t.Ty, z + t.Tz
}

func (t GeocentricTranslation) FromWGS84(x0, y0, z0 float64) (x, y, z float64) {
	return x0 - t.Tx, y0 - t.Ty, z0 - t.Tz
}
