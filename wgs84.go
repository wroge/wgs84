package wgs84

import (
	"errors"
	"fmt"
	"math"
)

var (
	ErrInvalidDatum = errors.New("invalid datum")
	ErrOutOfBounds  = errors.New("coordinate outside transformation area of use")
	ErrGridNotFound = errors.New("grid not found")
	ErrCodeNotFound = errors.New("epsg code not found")

	nullGrid = "null"

	World = BoundingBox{-180, -90, 180, 90}

	Airy1830          = Spheroid{6377563.396, 299.3249646}
	Bessel1841        = Spheroid{6377397.155, 299.1528128}
	GRS80             = Spheroid{6378137, 298.257222101}
	International1924 = Spheroid{6378388, 297}
	AiryModified1849  = Spheroid{6377340.189, 299.3249646}
	CGCS2000          = Spheroid{6378137, 298.257222101}
	Clarke1880        = Spheroid{6378249.2, 293.466021293627}
	Clarke1866        = Spheroid{6378206.4, 294.978698213898}

	WGS84 = Datum{
		Spheroid: Spheroid{6378137, 298.257223563},
		Transformations: []Transformation{
			{
				BoundingBox: World,
			},
		},
	}

	IRENET95 = Datum{
		Spheroid: GRS80,
		Transformations: []Transformation{
			{
				BoundingBox: BoundingBox{-10.56, 51.39, -5.34, 55.43},
			},
		},
	}

	RGF93 = Datum{
		Spheroid: GRS80,
		Transformations: []Transformation{
			{
				BoundingBox: BoundingBox{-9.86, 41.15, 10.38, 51.56},
			},
		},
	}

	ETRS89 = Datum{
		Spheroid: GRS80,
		Transformations: []Transformation{
			{
				BoundingBox: BoundingBox{-16.1, 33.26, 38.01, 84.73},
			},
		},
	}

	TM65 = Datum{
		Spheroid: AiryModified1849,
		Transformations: []Transformation{
			{
				Accuracy:    1,
				Helmert:     Helmert{482.5, -130.6, 564.6, -1.042, -0.214, -0.631, 8.15},
				BoundingBox: BoundingBox{-10.56, 51.39, -5.34, 55.43},
			},
		},
	}

	TM75 = Datum{
		Spheroid: AiryModified1849,
		Transformations: []Transformation{
			{
				Accuracy:    1,
				Helmert:     Helmert{482.5, -130.6, 564.6, -1.042, -0.214, -0.631, 8.15},
				BoundingBox: BoundingBox{-10.56, 51.39, -5.34, 55.43},
			},
		},
	}

	ChinaGeodeticCoordinateSystem2000 = Datum{
		Spheroid: CGCS2000,
		Transformations: []Transformation{
			{
				BoundingBox: BoundingBox{73.62, 16.7, 134.77, 53.56},
			},
		},
	}

	MGI = Datum{
		Spheroid: Bessel1841,
		Transformations: []Transformation{
			{
				Helmert:     Helmert{577.326, 90.129, 463.919, 5.137, 1.474, 5.297, 2.4232},
				BoundingBox: BoundingBox{9.53, 46.4, 17.17, 49.02},
			},
		},
	}

	BD72 = Datum{
		Spheroid: International1924,
		Transformations: []Transformation{
			{
				Accuracy:    1,
				Helmert:     Helmert{-106.8686, 52.2978, -103.7239, 0.3366, -0.457, 1.8422, -1.2747},
				BoundingBox: BoundingBox{2.5, 49.5, 6.4, 51.51},
			},
			{
				Accuracy:    1,
				Helmert:     Helmert{-99.059, 53.322, -112.486, -0.419, 0.83, -1.885, -1},
				BoundingBox: BoundingBox{2.5, 49.5, 6.4, 51.51},
			},
			{
				Accuracy:    5,
				Helmert:     Helmert{-125.8, 79.9, -100.5, 0, 0, 0, 0},
				BoundingBox: BoundingBox{2.5, 49.5, 6.4, 51.51},
			},
		},
	}

	NTF = Datum{
		Spheroid: Clarke1880,
		Transformations: []Transformation{
			{
				BoundingBox: BoundingBox{-4.87, 41.31, 9.63, 51.14},
				Helmert:     Helmert{-168, -60, 320, 0, 0, 0, 0},
			},
		},
	}

	CH1903 = Datum{
		Spheroid: Bessel1841,
		Transformations: []Transformation{
			{
				BoundingBox: BoundingBox{5.95, 45.81, 10.5, 47.81},
				Helmert:     Helmert{674.374, 15.056, 405.346, 0, 0, 0, 0},
			},
		},
	}
)

type Func func(float64, float64, float64) (float64, float64, float64, error)

func (f Func) Round(decA, decB, decC int) Func {
	return func(a, b, c float64) (float64, float64, float64, error) {
		a, b, c, err := f(a, b, c)

		return round(a, decA), round(b, decB), round(c, decC), err
	}
}

type BoundingBox struct {
	MinLon, MinLat, MaxLon, MaxLat float64
}

func (b BoundingBox) Contains(lon, lat float64) bool {
	if lat < b.MinLat || lat > b.MaxLat {
		return false
	}
	if b.MinLon <= b.MaxLon {
		return lon >= b.MinLon && lon <= b.MaxLon
	}

	return lon >= b.MinLon || lon <= b.MaxLon
}

type Helmert struct {
	Tx, Ty, Tz, Rx, Ry, Rz, Ds float64
}

func (h Helmert) String() string {
	return fmt.Sprintf("+towgs84=%s,%s,%s,%s,%s,%s,%s",
		projFloat(h.Tx), projFloat(h.Ty), projFloat(h.Tz), projFloat(h.Rx), projFloat(h.Ry), projFloat(h.Rz), projFloat(h.Ds))
}

func (h Helmert) ToWGS84(x, y, z float64) (x0, y0, z0 float64) {
	return calcHelmert(x, y, z, h.Tx, h.Ty, h.Tz, h.Rx, h.Ry, h.Rz, h.Ds)
}

func (h Helmert) FromWGS84(x0, y0, z0 float64) (x, y, z float64) {
	return calcHelmert(x0, y0, z0, -h.Tx, -h.Ty, -h.Tz, -h.Rx, -h.Ry, -h.Rz, -h.Ds)
}

func calcHelmert(x, y, z, tx, ty, tz, rx, ry, rz, ds float64) (x0, y0, z0 float64) {
	const (
		asec = math.Pi / 648000
		ppm  = 0.000001
	)

	x0 = (1+ds*ppm)*(x+z*ry*asec-y*rz*asec) + tx
	y0 = (1+ds*ppm)*(y+x*rz*asec-z*rx*asec) + ty
	z0 = (1+ds*ppm)*(z+y*rx*asec-x*ry*asec) + tz

	return
}

type Spheroid struct {
	A, Fi float64
}

func (s Spheroid) String() string {
	if s.Fi == 0 {
		return fmt.Sprintf("+a=%s", projFloat(s.A))
	}

	return fmt.Sprintf("+a=%s +rf=%s", projFloat(s.A), projFloat(s.Fi))
}

func (s Spheroid) F() float64 {
	return 1 / s.Fi
}

func (s Spheroid) A2() float64 {
	return s.A * s.A
}

func (s Spheroid) F2() float64 {
	f := s.F()

	return f * f
}

func (s Spheroid) B() float64 {
	return s.A * (1 - s.F())
}

func (s Spheroid) E2() float64 {

	return 2/s.Fi - s.F2()
}

func (s Spheroid) E() float64 {

	return math.Sqrt(s.E2())

}

func (s Spheroid) E4() float64 {
	e2 := s.E2()

	return e2 * e2
}

func (s Spheroid) E6() float64 {
	e2 := s.E2()

	return e2 * e2 * e2
}

func (s Spheroid) Ei() float64 {
	e2 := s.E2()
	t := math.Sqrt(1 - e2)

	return (1 - t) / (1 + t)
}

func (s Spheroid) Ei2() float64 {
	ei := s.Ei()

	return ei * ei
}

func (s Spheroid) Ei3() float64 {
	ei := s.Ei()

	return ei * ei * ei
}

func (s Spheroid) Ei4() float64 {
	ei := s.Ei()

	return ei * ei * ei * ei
}

func (s Spheroid) ToXYZ(lon, lat, h float64) (x, y, z float64) {
	n := s.A / math.Sqrt(1-s.E2()*intPow(math.Sin(radian(lat)), 2))

	x = (n + h) * math.Cos(radian(lon)) * math.Cos(radian(lat))
	y = (n + h) * math.Cos(radian(lat)) * math.Sin(radian(lon))
	z = (n*intPow(s.A*(1-s.F()), 2)/(s.A2()) + h) * math.Sin(radian(lat))

	return x, y, z
}

func (s Spheroid) FromXYZ(x, y, z float64) (lon, lat, h float64) {
	sd := math.Sqrt(x*x + y*y)
	T := math.Atan(z * s.A / (sd * s.B()))
	B := math.Atan((z + s.E2()*(s.A2())/s.B()*
		intPow(math.Sin(T), 3)) / (sd - s.E2()*s.A*intPow(math.Cos(T), 3)))
	n := s.A / math.Sqrt(1-s.E2()*intPow(math.Sin(B), 2))
	h = sd/math.Cos(B) - n
	lon = degree(math.Atan2(y, x))
	lat = degree(B)

	return lon, lat, h
}

type CoordinateSystem interface {
	ToGeographic(a, b, c float64, s Spheroid) (lon, lat, h float64)
	FromGeographic(lon, lat, h float64, s Spheroid) (a, b, c float64)
}

type CoordinateReferenceSystem struct {
	CoordinateSystem CoordinateSystem
	Datum            Datum
}

func (crs CoordinateReferenceSystem) String() string {
	return fmt.Sprintf("%s %s", crs.CoordinateSystem, crs.Datum)
}

func (crs CoordinateReferenceSystem) Load(lon, lat float64) (CoordinateReferenceSystem, error) {
	datum, err := crs.Datum.Load(lon, lat)
	if err != nil {
		return CoordinateReferenceSystem{}, err
	}

	return CoordinateReferenceSystem{
		CoordinateSystem: crs.CoordinateSystem,
		Datum:            datum,
	}, nil
}

type CRS interface {
	string | int | CoordinateReferenceSystem
}

type CS interface {
	Geocentric | Geographic | WebMercator | TransverseMercator | SwissObliqueMercator | LambertConformalConic1SP | LambertConformalConic2SP | LambertAzimuthalEqualArea | Krovak | AlbersConicEqualArea
}

func Transform[F CRS, T CRS | CS](from F, to T) (Func, error) {
	var (
		fromCRS, toCRS CoordinateReferenceSystem
		ok             bool
		err            error
	)

	switch f := any(from).(type) {
	case string:
		fromCRS, err = ParseProj(f)
		if err != nil {
			return nil, err
		}
	case int:
		fromCRS, ok = EPSG[f]
		if !ok {
			return nil, ErrCodeNotFound
		}
	case CoordinateReferenceSystem:
		fromCRS = f
	default:
		return nil, fmt.Errorf("invalid crs")
	}

	switch t := any(to).(type) {
	case string:
		toCRS, err = ParseProj(t)
		if err != nil {
			cs, err2 := parseCoordinateSystem(t)
			if err2 != nil {
				return nil, err
			}

			toCRS = CoordinateReferenceSystem{
				Datum:            fromCRS.Datum,
				CoordinateSystem: cs,
			}
		}
	case int:
		toCRS, ok = EPSG[t]
		if !ok {
			return nil, ErrCodeNotFound
		}
	case CoordinateReferenceSystem:
		toCRS = t
	case CoordinateSystem:
		toCRS = CoordinateReferenceSystem{
			Datum:            fromCRS.Datum,
			CoordinateSystem: t,
		}
	default:
		return nil, fmt.Errorf("invalid crs")
	}

	return fromCRS.TransformTo(toCRS), nil
}

func (crs CoordinateReferenceSystem) TransformTo(to CoordinateReferenceSystem) Func {
	if crs.Datum.Equals(to.Datum) {
		return func(a, b, c float64) (float64, float64, float64, error) {
			lon, lat, h := crs.CoordinateSystem.ToGeographic(a, b, c, crs.Datum.Spheroid)

			a, b, c = to.CoordinateSystem.FromGeographic(lon, lat, h, crs.Datum.Spheroid)

			return a, b, c, nil
		}
	}

	return func(a, b, c float64) (float64, float64, float64, error) {
		lon, lat, h := crs.CoordinateSystem.ToGeographic(a, b, c, crs.Datum.Spheroid)

		a0, b0, c0, geocent, err := crs.Datum.toWGS84(lon, lat, h)
		if err != nil {
			return 0, 0, 0, err
		}

		var lon2, lat2, h2 float64

		if geocent {
			lon2, lat2, h2, err = to.Datum.fromGeocentricWGS84(a0, b0, c0)
		} else {
			lon2, lat2, h2, err = to.Datum.fromGeographicWGS84(a0, b0, c0)
		}

		if err != nil {
			return 0, 0, 0, err
		}

		a, b, c = to.CoordinateSystem.FromGeographic(lon2, lat2, h2, to.Datum.Spheroid)

		return a, b, c, nil
	}
}

type Transformation struct {
	EPSG         int
	Accuracy     float64
	Grid         string
	GridOptional bool
	Helmert      Helmert
	BoundingBox  BoundingBox
}

func (t Transformation) String() string {
	if t.Grid != "" {
		grid := t.Grid
		if t.GridOptional {
			grid = "@" + grid
		}

		return fmt.Sprintf("+nadgrids=%s", grid)
	}

	if (t.Helmert == Helmert{}) {
		return ""
	}

	return t.Helmert.String()
}

func (t Transformation) fromGeocentricWGS84(x0, y0, z0 float64, s Spheroid) (lon, lat, h float64, err error) {
	if t.Grid == nullGrid {
		lon, lat, h = WGS84.Spheroid.FromXYZ(x0, y0, z0)

		return lon, lat, h, nil
	}

	if t.Grid != "" {
		grid, err := LoadGrid(t.Grid)
		if err != nil {
			return 0, 0, 0, err
		}

		lon0, lat0, h0 := WGS84.Spheroid.FromXYZ(x0, y0, z0)

		lon, lat, err = grid.FromWGS84(lon0, lat0)
		if err != nil {
			return 0, 0, 0, err
		}

		return lon, lat, h0, nil
	}

	x, y, z := t.Helmert.FromWGS84(x0, y0, z0)

	lon, lat, h = s.FromXYZ(x, y, z)

	return lon, lat, h, nil
}

func (t Transformation) fromGeographicWGS84(lon0, lat0, h0 float64, s Spheroid) (lon, lat, h float64, err error) {
	if !t.BoundingBox.Contains(lon0, lat0) {
		return 0, 0, 0, ErrOutOfBounds
	}

	if t.Grid == nullGrid {
		return lon0, lat0, h0, nil
	}

	if t.Grid != "" {
		grid, err := LoadGrid(t.Grid)
		if err != nil {
			return 0, 0, 0, err
		}

		lon, lat, err = grid.FromWGS84(lon0, lat0)
		if err != nil {
			return 0, 0, 0, err
		}

		return lon, lat, h0, nil
	}

	x0, y0, z0 := WGS84.Spheroid.ToXYZ(lon0, lat0, h0)

	x, y, z := t.Helmert.FromWGS84(x0, y0, z0)

	lon, lat, h = s.FromXYZ(x, y, z)

	return lon, lat, h, nil
}

func (t Transformation) toWGS84(lon, lat, h float64, s Spheroid) (a0, b0, c0 float64, geocent bool, err error) {
	if t.Grid == nullGrid {
		if !t.BoundingBox.Contains(lon, lat) {
			return 0, 0, 0, false, ErrOutOfBounds
		}

		return lon, lat, h, false, nil
	}

	if t.Grid != "" {
		grid, err := LoadGrid(t.Grid)
		if err != nil {
			return 0, 0, 0, false, err
		}

		lon0, lat0, err := grid.ToWGS84(lon, lat)
		if err != nil {
			return 0, 0, 0, false, err
		}

		if !t.BoundingBox.Contains(lon0, lat0) {
			return 0, 0, 0, false, ErrOutOfBounds
		}

		return lon0, lat0, h, false, nil
	}

	x, y, z := s.ToXYZ(lon, lat, h)

	x0, y0, z0 := t.Helmert.ToWGS84(x, y, z)

	lon0, lat0, _ := WGS84.Spheroid.FromXYZ(x0, y0, z0)

	if !t.BoundingBox.Contains(lon0, lat0) {
		return 0, 0, 0, false, ErrOutOfBounds
	}

	return x0, y0, z0, true, nil
}

type Datum struct {
	Spheroid        Spheroid
	Transformations []Transformation
}

func (d Datum) String() string {
	if len(d.Transformations) == 0 {
		return "invalid datum"
	}

	return fmt.Sprintf("%s %s", d.Spheroid, d.Transformations[0])
}

func (d Datum) Equals(other Datum) bool {
	if d.Spheroid != other.Spheroid {
		return false
	}

	if len(d.Transformations) != len(other.Transformations) {
		return false
	}

	for i := range d.Transformations {
		if d.Transformations[i] != other.Transformations[i] {
			return false
		}
	}

	return true
}

func (d Datum) Load(lon, lat float64) (Datum, error) {
	for _, t := range d.Transformations {
		if t.Grid != "" && t.Grid != nullGrid {
			grid, err := LoadGrid(t.Grid)
			if err != nil {
				if t.GridOptional {
					continue
				}

				return Datum{}, err
			}

			if grid.selectSubgrid(lon, lat) < 0 {
				continue
			}
		}

		if t.BoundingBox.Contains(lon, lat) {
			t.GridOptional = false

			return Datum{
				Spheroid: d.Spheroid,
				Transformations: []Transformation{
					t,
				},
			}, nil
		}
	}

	return Datum{}, ErrOutOfBounds
}

func (d Datum) fromGeocentricWGS84(x0, y0, z0 float64) (lon, lat, h float64, err error) {
	for _, t := range d.Transformations {
		lon, lat, h, err = t.fromGeocentricWGS84(x0, y0, z0, d.Spheroid)
		if err != nil {
			if errors.Is(err, ErrGridNotFound) && t.GridOptional {
				continue
			}

			if errors.Is(err, ErrGridNotFound) {
				return 0, 0, 0, err
			}

			continue
		}

		return lon, lat, h, nil
	}

	return 0, 0, 0, ErrInvalidDatum
}

func (d Datum) fromGeographicWGS84(lon0, lat0, h0 float64) (lon, lat, h float64, err error) {
	for _, t := range d.Transformations {
		lon, lat, h, err = t.fromGeographicWGS84(lon0, lat0, h0, d.Spheroid)
		if err != nil {
			if errors.Is(err, ErrGridNotFound) && t.GridOptional {
				continue
			}

			if errors.Is(err, ErrGridNotFound) {
				return 0, 0, 0, err
			}

			continue
		}

		return lon, lat, h, nil
	}

	return 0, 0, 0, ErrInvalidDatum
}

func (d Datum) toWGS84(lon, lat, h float64) (a0, b0, c0 float64, geocent bool, err error) {
	for _, t := range d.Transformations {
		a0, b0, c0, geocent, err = t.toWGS84(lon, lat, h, d.Spheroid)
		if err != nil {
			if errors.Is(err, ErrGridNotFound) && t.GridOptional {
				continue
			}

			if errors.Is(err, ErrGridNotFound) {
				return 0, 0, 0, false, err
			}

			continue
		}

		if !geocent {
			return a0, b0, h, false, nil
		}

		return a0, b0, c0, true, nil
	}

	return 0, 0, 0, false, ErrInvalidDatum
}

type Geocentric struct{}

func (Geocentric) String() string {
	return "+proj=geocent"
}

func (Geocentric) FromGeographic(lon, lat, h float64, s Spheroid) (a, b, c float64) {
	return s.ToXYZ(lon, lat, h)
}

func (Geocentric) ToGeographic(a, b, c float64, s Spheroid) (lon, lat, h float64) {
	return s.FromXYZ(a, b, c)
}

type Geographic struct{}

func (Geographic) String() string {
	return "+proj=longlat"
}

func (Geographic) FromGeographic(lon, lat, h float64, s Spheroid) (a, b, c float64) {
	return lon, lat, h
}

func (Geographic) ToGeographic(a, b, c float64, s Spheroid) (lon, lat, h float64) {
	return a, b, c
}
