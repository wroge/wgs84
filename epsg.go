package wgs84

import (
	"sync"

	"github.com/wroge/wgs84/system"
)

func EPSG() *Repository {
	codes := map[int]CoordinateReferenceSystem{}
	codes[4326] = LonLat()
	codes[4978] = CoordinateReferenceSystem{}
	codes[3857] = WebMercator()
	codes[900913] = WebMercator()
	codes[27700] = OSGB36NationalGrid()
	codes[4277] = OSGB36LonLat()
	codes[2154] = RGF93FranceLambert()
	for i := 1; i < 61; i++ {
		codes[32600+i] = CoordinateReferenceSystem{System: system.UTM(float64(i), true)}
		codes[32700+i] = CoordinateReferenceSystem{System: system.UTM(float64(i), false)}
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

type Repository struct {
	codes map[int]CoordinateReferenceSystem
	mutex sync.Mutex
}

func (r *Repository) Code(c int) CoordinateReferenceSystem {
	if r.codes == nil {
		return CoordinateReferenceSystem{}
	}
	return r.codes[c]
}

func (r *Repository) Add(c int, crs CoordinateReferenceSystem) {
	if r.codes == nil {
		r.codes = map[int]CoordinateReferenceSystem{}
	}
	r.mutex.Lock()
	r.codes[c] = crs
	r.mutex.Unlock()
}

func (r *Repository) Codes() []int {
	r.mutex.Lock()
	var cc []int
	for c := range r.codes {
		cc = append(cc, c)
	}
	r.mutex.Unlock()
	return cc
}
