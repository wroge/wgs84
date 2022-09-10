//nolint:varnamelen,ireturn,exhaustivestruct,exhaustruct,nonamedreturns
package wgs84

import (
	"sync"
)

// EPSG returns a Repository for dealing with several EPSG-Codes and
// CoordinateReferenceSystems.
func EPSG() *Repository {
	codes := map[int]CoordinateReferenceSystem{
		4326:   LonLat(),
		4978:   XYZ(),
		3857:   WebMercator(),
		900913: WebMercator(),
		4258:   ETRS89().LonLat(),
		3416:   ETRS89AustriaLambert(),
		3035:   ETRS89LambertAzimuthalEqualArea(),
		31287:  MGIAustriaLambert(),
		31284:  MGIAustriaM28(),
		31285:  MGIAustriaM31(),
		31286:  MGIAustriaM34(),
		31257:  MGIAustriaGKM28(),
		31258:  MGIAustriaGKM31(),
		31259:  MGIAustriaGKM34(),
		4314:   DHDN2001().LonLat(),
		27700:  OSGB36NationalGrid(),
		4277:   OSGB36().LonLat(),
		4171:   RGF93().LonLat(),
		2154:   RGF93FranceLambert(),
		4269:   NAD83().LonLat(),
		6355:   NAD83AlabamaEast(),
		6356:   NAD83AlabamaWest(),
		6414:   NAD83CaliforniaAlbers(),
	}

	for i := 1; i < 61; i++ {
		codes[32600+i] = UTM(float64(i), true)
		codes[32700+i] = UTM(float64(i), false)
	}

	for i := 42; i < 51; i++ {
		codes[3900+i] = RGF93CC(float64(i))
	}

	for i := 2; i < 6; i++ {
		codes[31464+i] = DHDN2001GK(float64(i))
	}

	for i := 28; i < 39; i++ {
		codes[25800+i] = ETRS89UTM(float64(i))
	}

	return &Repository{
		codes: codes,
	}
}

// Repository holds the EPSG-Codes and CoordinateReferenceSystems.
type Repository struct {
	codes map[int]CoordinateReferenceSystem
	mutex sync.Mutex
}

// Code returns a CoordinateReferenceSystem of a specific EPSG-Code.
func (r *Repository) Code(c int) CoordinateReferenceSystem {
	if r.codes == nil {
		return nil
	}

	return r.codes[c]
}

// Code returns a CoordinateReferenceSystem of a specific EPSG-Code.
func (r *Repository) SafeCode(c int) (CoordinateReferenceSystem, error) {
	if r.codes == nil {
		return nil, ErrNoCoordinateReferenceSystem
	}

	return r.codes[c], nil
}

// Add an EPSG-Code to the Repository.
func (r *Repository) Add(c int, crs CoordinateReferenceSystem) {
	if crs == nil {
		return
	}

	if r.codes == nil {
		r.codes = map[int]CoordinateReferenceSystem{}
	}

	r.mutex.Lock()
	r.codes[c] = crs
	r.mutex.Unlock()
}

// Codes returns all available codes.
func (r *Repository) Codes() []int {
	r.mutex.Lock()

	cc := make([]int, len(r.codes))
	n := 0

	for c := range r.codes {
		cc[n] = c
		n++
	}

	r.mutex.Unlock()

	return cc
}

// CodesCover returns all Codes covering a specific geographic WGS84 location.
func (r *Repository) CodesCover(lon, lat float64) []int {
	r.mutex.Lock()

	var cc []int

	for c, crs := range r.codes {
		if crs.Contains(lon, lat) {
			cc = append(cc, c)
		}
	}

	r.mutex.Unlock()

	return cc
}

// Transform transforms coordinates from one EPSG-Code to another.
func (r *Repository) Transform(from, to int) Func {
	return Transform(r.Code(from), r.Code(to))
}

// SafeTransform transforms coordinates from one EPSG-Code to another
// with errors.
func (r *Repository) SafeTransform(from, to int) SafeFunc {
	f, err := r.SafeCode(from)
	if err != nil {
		return func(_, _, _ float64) (_, _, _ float64, err error) {
			return 0, 0, 0, err
		}
	}

	t, err := r.SafeCode(to)
	if err != nil {
		return func(_, _, _ float64) (_, _, _ float64, err error) {
			return 0, 0, 0, err
		}
	}

	return SafeTransform(f, t)
}
