package epsg

import (
	"sync"

	"github.com/wroge/wgs84"
	"github.com/wroge/wgs84/dhdn2001"
	"github.com/wroge/wgs84/etrs89"
	"github.com/wroge/wgs84/osgb36"
	"github.com/wroge/wgs84/rgf93"
)

type epsg struct {
	crs  wgs84.CoordinateReferenceSystem
	bbox boundingBox
}

type boundingBox struct {
	minlon, minlat, maxlon, maxlat float64
}

func (b boundingBox) Contains(lon, lat float64) bool {
	for lon < -180 {
		lon += 360
	}
	for lon > 180 {
		lon -= 360
	}
	for lat < -90 {
		lat += 180
	}
	for lat > 90 {
		lat -= 180
	}
	if b.minlon > b.maxlon {
		if lon < b.minlon && lon > b.maxlon {
			return false
		}
	} else {
		if lon < b.minlon || lon > b.maxlon {
			return false
		}
	}
	if lat < b.minlat || lat > b.maxlat {
		return false
	}
	return true
}

func (b boundingBox) normalize() boundingBox {
	for b.minlon < -180 {
		b.minlon += 360
	}
	for b.maxlon > 180 {
		b.minlon -= 360
	}
	for b.minlat < -90 {
		b.minlat += 180
	}
	for b.maxlat > 90 {
		b.maxlat -= 180
	}
	return b
}

type Repository struct {
	codes map[int]epsg
	m     sync.Mutex
}

func (r *Repository) Add(c int, crs wgs84.CoordinateReferenceSystem, minlon, minlat, maxlon, maxlat float64) {
	r.m.Lock()
	defer r.m.Unlock()
	b := boundingBox{minlon, minlat, maxlon, maxlat}.normalize()
	if r.codes == nil {
		r.codes = map[int]epsg{}
	}
	r.codes[c] = epsg{crs, b}
}

func Codes(lon, lat float64) []int {
	return DefaultRepository().Codes(lon, lat)
}

func (r *Repository) Codes(lon, lat float64) []int {
	r.m.Lock()
	defer r.m.Unlock()
	cc := []int{}
	for c, e := range r.codes {
		if e.bbox.Contains(lon, lat) {
			cc = append(cc, c)
		}
	}
	return cc
}

func Code(c int) (wgs84.CoordinateReferenceSystem, bool) {
	return DefaultRepository().Code(c)
}

func (r *Repository) Code(c int) (wgs84.CoordinateReferenceSystem, bool) {
	r.m.Lock()
	defer r.m.Unlock()
	e, ok := r.codes[c]
	return e.crs, ok
}

func DefaultRepository() *Repository {
	r := &Repository{}
	r.codes = map[int]epsg{
		4326:   {wgs84.LonLat(), boundingBox{-180, -90, 180, 90}},
		3857:   {wgs84.WebMercator(), boundingBox{-180, -85.06, 180, 85.06}},
		900913: {wgs84.WebMercator(), boundingBox{-180, -85.06, 180, 85.06}},
		4978:   {wgs84.XYZ(), boundingBox{-180, -90, 180, 90}},
		2154:   {rgf93.FranceLambert(), boundingBox{-9.86, 41.15, 10.38, 51.56}},
		27700:  {osgb36.NationalGrid(), boundingBox{-8.82, 49.79, 1.92, 60.94}},
		4277:   {osgb36.LonLat(), boundingBox{-8.82, 49.79, 1.92, 60.94}},
	}
	for i := 1; i < 61; i++ {
		r.Add(32600+i, wgs84.UTM(float64(i), true), float64(i)*6-186, 0, float64(i)*6-180, 84)
		r.Add(32700+i, wgs84.UTM(float64(i), false), float64(i)*6-186, -80, float64(i)*6-180, 0)
	}
	for i := 42; i < 51; i++ {
		r.Add(3900+i, rgf93.CC(float64(i)), -9.86, 41.15, 10.38, 51.56)
	}
	for i := 2; i < 6; i++ {
		r.Add(31464+i, dhdn2001.GK(float64(i)), float64(i)*3-1.5, 47.27, float64(i)*3+1.5, 55.09)
	}
	for i := 28; i < 39; i++ {
		r.Add(25800+i, etrs89.UTM(float64(i)), float64(i)*6-186, 32.88, float64(i)*6-186, 84.17)
	}
	return r
}
