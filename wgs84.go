//nolint:varnamelen,nonamedreturns,ireturn,gomnd,exhaustivestruct,exhaustruct,cyclop,errname,lll,funlen
package wgs84

import (
	"embed"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"
	"strings"
)

type CRS interface {
	Base() CRS
	Spheroid() Spheroid
	ToBase(float64, float64, float64) (float64, float64, float64)
	FromBase(float64, float64, float64) (float64, float64, float64)
}

func Transform(from, to CRS) Func {
	comparable := reflect.ValueOf(from).Comparable() && reflect.ValueOf(to).Comparable()

	if comparable && from == to {
		return identityFunc
	}

	ancestors := make(map[CRS]int)

	if comparable {
		depth := 0
		for c := from; c != nil; c = c.Base() {
			ancestors[c] = depth
			depth++
		}
	}

	var (
		toBase   []Func
		fromBase []Func
		lca      CRS
	)

	for c := to; c != nil; c = c.Base() {
		if comparable {
			if _, ok := ancestors[c]; ok {
				lca = c
				break
			}
		}

		fromBase = append(fromBase, c.FromBase)
	}

	for c := from; c != lca; c = c.Base() {
		toBase = append(toBase, c.ToBase)
	}

	steps := make([]Func, 0, len(toBase)+len(fromBase))
	for _, f := range toBase {
		if f != nil {
			steps = append(steps, f)
		}
	}

	for i := len(fromBase) - 1; i >= 0; i-- {
		if fromBase[i] != nil {
			steps = append(steps, fromBase[i])
		}
	}

	switch len(steps) {
	case 0:
		return identityFunc
	case 1:
		return steps[0]
	default:
		return chainFunc(steps...)
	}
}

var defaultSpheroid = NewSpheroid(6378137, 298.257223563)

func identityFunc(a, b, c float64) (float64, float64, float64) {
	return a, b, c
}

func chainFunc(steps ...Func) Func {
	switch len(steps) {
	case 0:
		return identityFunc
	case 1:
		return steps[0]
	}

	return func(a, b, c float64) (float64, float64, float64) {
		for _, step := range steps {
			if step != nil {
				a, b, c = step(a, b, c)
			}
		}

		return a, b, c
	}
}

type Func func(float64, float64, float64) (float64, float64, float64)

func (f Func) Round(dec int) Func {
	return func(a, b, c float64) (float64, float64, float64) {
		a, b, c = f(a, b, c)

		return round(a, dec), round(b, dec), round(c, dec)
	}
}

func round(val float64, dec int) float64 {
	factor := math.Pow(10, float64(dec))

	val = math.Round(val*factor) / factor

	if val == -0 {
		return 0
	}

	return val
}

type Error struct {
	Cause error
}

func (e Error) Unwrap() error {
	return e.Cause
}

func (e Error) Error() string {
	return e.Cause.Error()
}

func (Error) Base() CRS {
	return nil
}

func (Error) Spheroid() Spheroid {
	return Spheroid{}
}

func (Error) ToBase(_, _, _ float64) (float64, float64, float64) {
	return math.NaN(), math.NaN(), math.NaN()
}

func (Error) FromBase(_, _, _ float64) (float64, float64, float64) {
	return math.NaN(), math.NaN(), math.NaN()
}

type Spheroid struct {
	A, Fi                                          float64
	A2, F, F2, B, E2, E, E4, E6, Ei, Ei2, Ei3, Ei4 float64
}

func NewSpheroid(a, fi float64) Spheroid {
	s := Spheroid{
		A:  a,
		Fi: fi,
	}

	s.A2 = s.A * s.A
	s.F = 1 / s.Fi
	s.F2 = s.F * s.F
	s.B = s.A * (1 - s.F)
	s.E2 = 2/s.Fi - s.F2
	s.E = math.Sqrt(s.E2)
	s.E4 = s.E2 * s.E2
	s.E6 = s.E4 * s.E2
	s.Ei = (1 - math.Sqrt(1-s.E2)) / (1 + math.Sqrt(1-s.E2))
	s.Ei2 = s.Ei * s.Ei
	s.Ei3 = s.Ei2 * s.Ei
	s.Ei4 = s.Ei3 * s.Ei

	return s
}

func (s Spheroid) ToXYZ(lon, lat, h float64) (x, y, z float64) {
	n := s.A / math.Sqrt(1-s.E2*intPow(math.Sin(radian(lat)), 2))

	x = (n + h) * math.Cos(radian(lon)) * math.Cos(radian(lat))
	y = (n + h) * math.Cos(radian(lat)) * math.Sin(radian(lon))
	z = (n*intPow(s.A*(1-s.F), 2)/(s.A2) + h) * math.Sin(radian(lat))

	return x, y, z
}

func (s Spheroid) FromXYZ(x, y, z float64) (lon, lat, h float64) {
	sd := math.Sqrt(x*x + y*y)
	T := math.Atan(z * s.A / (sd * s.B))
	B := math.Atan((z + s.E2*(s.A2)/s.B*
		intPow(math.Sin(T), 3)) / (sd - s.E2*s.A*intPow(math.Cos(T), 3)))
	n := s.A / math.Sqrt(1-s.E2*intPow(math.Sin(B), 2))
	h = sd/math.Cos(B) - n
	lon = degree(math.Atan2(y, x))
	lat = degree(B)

	return lon, lat, h
}

type base struct{}

func (base) Base() CRS {
	return nil
}

func (base) Spheroid() Spheroid {
	return Spheroid{}
}

func (base) ToBase(x0, y0, z0 float64) (float64, float64, float64) {
	return x0, y0, z0
}

func (base) FromBase(x0, y0, z0 float64) (float64, float64, float64) {
	return x0, y0, z0
}

func Geographic(geocentric CRS, spheroid Spheroid) CRS {
	if geocentric == nil {
		geocentric = base{}
	}

	return geographic{
		b: geocentric,
		s: spheroid,
	}
}

type geographic struct {
	b CRS
	s Spheroid
}

func (b geographic) Base() CRS {
	return b.b
}

func (b geographic) Spheroid() Spheroid {
	return b.s
}

func (b geographic) ToBase(lon, lat, h float64) (x, y, z float64) {
	return b.s.ToXYZ(lon, lat, h)
}

func (b geographic) FromBase(x, y, z float64) (lon, lat, h float64) {
	return b.s.FromXYZ(x, y, z)
}

func Helmert(tx, ty, tz, rx, ry, rz, ds float64) CRS {
	return helmert{
		tx: tx,
		ty: ty,
		tz: tz,
		rx: rx,
		ry: ry,
		rz: rz,
		ds: ds,
	}
}

type helmert struct {
	tx, ty, tz, rx, ry, rz, ds float64
}

func (t helmert) Base() CRS {
	return base{}
}

func (t helmert) Spheroid() Spheroid {
	return Spheroid{}
}

func (t helmert) ToBase(x, y, z float64) (x0, y0, z0 float64) {
	return calcHelmert(x, y, z, t.tx, t.ty, t.tz, t.rx, t.ry, t.rz, t.ds)
}

func (t helmert) FromBase(x0, y0, z0 float64) (x, y, z float64) {
	return calcHelmert(x0, y0, z0, -t.tx, -t.ty, -t.tz, -t.rx, -t.ry, -t.rz, -t.ds)
}

const (
	asec = math.Pi / 648000
	ppm  = 0.000001
)

func calcHelmert(x, y, z, tx, ty, tz, rx, ry, rz, ds float64) (x0, y0, z0 float64) {
	x0 = (1+ds*ppm)*(x+z*ry*asec-y*rz*asec) + tx
	y0 = (1+ds*ppm)*(y+x*rz*asec-z*rx*asec) + ty
	z0 = (1+ds*ppm)*(z+y*rx*asec-x*ry*asec) + tz

	return
}

//go:embed ntv2
var res embed.FS

func loadNTv2(name string, spheroid Spheroid, base CRS) CRS {
	file, err := res.Open("ntv2/" + name)
	if err != nil {
		return Error{Cause: err}
	}

	crs := loadReaderNTv2(file, spheroid, base)

	return crs
}

func loadReaderNTv2(reader io.Reader, spheroid Spheroid, base CRS) CRS {
	if base == nil {
		base = Geographic(nil, defaultSpheroid)
	}

	data := &ntv2{
		base:     base,
		spheroid: spheroid,
	}

	for i := 1; ; i++ {
		set := make([]byte, 16)

		_, err := io.ReadFull(reader, set)
		if err != nil {
			return Error{Cause: err}
		}

		switch toString(set[:8]) {
		case "NUM_OREC":
			data.numOrec = int32(binary.LittleEndian.Uint32(set[8:]))
		case "NUM_SREC":
			data.numSrec = int32(binary.LittleEndian.Uint32(set[8:]))
		case "NUM_FILE":
			data.numFile = int32(binary.LittleEndian.Uint32(set[8:]))
		case "S_LAT":
			data.sLat = toFloat(set[8:])
		case "N_LAT":
			data.nLat = toFloat(set[8:])
		case "E_LONG":
			data.eLong = toFloat(set[8:])
		case "W_LONG":
			data.wLong = toFloat(set[8:])
		case "LAT_INC":
			data.latInc = toFloat(set[8:])
		case "LONG_INC":
			data.longInc = toFloat(set[8:])
		case "GS_COUNT":
			data.gsCount = int32(binary.LittleEndian.Uint32(set[8:]))
		default:
			offset := int(data.numOrec) + int(data.numSrec) + int(data.numFile)
			if i >= int(data.gsCount)+offset {
				return data
			}

			if i >= offset && i < int(data.gsCount)+offset {
				data.values = append(data.values, [4]float32{
					toFloat32(set[0:4]), toFloat32(set[4:8]), toFloat32(set[8:12]), toFloat32(set[12:16]),
				})

				continue
			}
		}
	}
}

func toFloat32(b []byte) float32 {
	i := binary.LittleEndian.Uint32(b)

	return math.Float32frombits(i)
}

func toFloat(b []byte) float64 {
	i := binary.LittleEndian.Uint64(b)

	return math.Float64frombits(i)
}

func toString(b []byte) string {
	return strings.TrimSpace(string(b))
}

type ntv2 struct {
	spheroid Spheroid
	base     CRS
	numOrec  int32
	numSrec  int32
	numFile  int32
	gsCount  int32
	sLat     float64
	nLat     float64
	eLong    float64
	wLong    float64
	latInc   float64
	longInc  float64
	values   [][4]float32
}

func (n ntv2) String() string {
	return fmt.Sprintf("BASE: %v, NUM_OREC: %d, NUM_SREC: %d, NUM_FILE: %d, S_LAT: %f, N_LAT: %f, E_LONG: %f, W_LONG: %f, LAT_INC: %f, LONG_INC: %f, GS_COUNT: %d", n.base, n.numOrec, n.numSrec, n.numFile, n.sLat, n.nLat, n.eLong, n.wLong, n.latInc, n.longInc, n.gsCount)
}

func (n ntv2) Base() CRS {
	return n.base
}

func (n ntv2) Spheroid() Spheroid {
	return n.spheroid
}

func (n ntv2) ToBase(lon, lat, h float64) (lon2, lat2, h2 float64) {
	slon, slat := n.Shift(lon, lat)

	return lon + slon, lat + slat, h
}

func (n ntv2) FromBase(lon, lat, h float64) (lon2, lat2, h2 float64) {
	qlon, qlat := lon, lat

	for i := 0; i < 10; i++ {
		slon, slat := n.Shift(qlon, qlat)
		newLon := lon - slon
		newLat := lat - slat

		if math.Abs(newLon-qlon) < 1e-12 && math.Abs(newLat-qlat) < 1e-12 {
			break
		}

		qlon, qlat = newLon, newLat
	}

	return qlon, qlat, h
}

func (n ntv2) Shift(lon, lat float64) (float64, float64) {
	fcol := (-lon*3600 - n.eLong) / n.longInc
	frow := (lat*3600 - n.sLat) / n.latInc

	col := math.Floor(fcol)
	row := math.Floor(frow)

	ppr := math.Floor((n.wLong-n.eLong)/n.longInc+0.5) + 1
	ppc := math.Floor((n.nLat-n.sLat)/n.latInc+0.5) + 1

	se := row*ppr + col
	sw := se + 1
	ne := se + ppr
	nw := ne + 1

	col = math.Max(0, math.Min(col, ppr-2))
	row = math.Max(0, math.Min(row, ppc-2))

	if col >= ppr-1 {
		sw = se
		nw = ne
	}

	if row >= ppc-1 {
		ne = se
		nw = sw
	}

	if col <= 0 {
		se = sw
		ne = nw
	}

	if row <= 0 {
		se = ne
		sw = nw
	}

	seIndex := min(max(int(se), 0), len(n.values)-1)
	swIndex := min(max(int(sw), 0), len(n.values)-1)
	neIndex := min(max(int(ne), 0), len(n.values)-1)
	nwIndex := min(max(int(nw), 0), len(n.values)-1)

	sse := n.values[seIndex]
	ssw := n.values[swIndex]
	sne := n.values[neIndex]
	snw := n.values[nwIndex]

	dx := fcol - col
	dy := frow - row

	latsv := (1-dx)*(1-dy)*float64(sse[0]) + dx*(1-dy)*float64(ssw[0]) + (1-dx)*dy*float64(sne[0]) + dx*dy*float64(snw[0])
	lonsv := (1-dx)*(1-dy)*float64(sse[1]) + dx*(1-dy)*float64(ssw[1]) + (1-dx)*dy*float64(sne[1]) + dx*dy*float64(snw[1])

	return -lonsv / 3600, latsv / 3600
}

func WebMercator(base CRS) CRS {
	if base == nil {
		base = Geographic(nil, defaultSpheroid)
	}

	return webMercator{
		base: base,
	}
}

type webMercator struct {
	base CRS
}

func (p webMercator) Base() CRS {
	return p.base
}

func (p webMercator) Spheroid() Spheroid {
	return p.base.Spheroid()
}

func (p webMercator) ToBase(east, north, h float64) (lon, lat, h2 float64) {
	s := p.base.Spheroid()

	D := (-north) / s.A
	phi := math.Pi/2 - 2*math.Atan(math.Pow(math.E, D))
	lambda := east / s.A

	return degree(lambda), degree(phi), h
}

func (p webMercator) FromBase(lon, lat, h float64) (east, north, h2 float64) {
	s := p.base.Spheroid()

	lambda := radian(lon)
	phi := radian(lat)

	east = s.A * lambda
	north = s.A * math.Log(math.Tan(math.Pi/4+phi/2))

	return east, north, h
}

func TransverseMercator(base CRS, lonf, latf, scale, eastf, northf float64) CRS {
	if base == nil {
		base = Geographic(nil, defaultSpheroid)
	}

	s := base.Spheroid()

	phi0 := radian(latf)
	lambda0 := radian(lonf)
	n := s.F / (2 - s.F)
	n2 := n * n
	n3 := n * n * n
	n4 := n * n * n * n
	B := (s.A / (1 + n)) * (1 + n2/4 + n4/64)
	h1 := n/2.0 - (2/3.0)*n2 + (5/16.0)*n3 + (41/180.0)*n4
	h2 := (13/48.0)*n2 - (3/5.0)*n3 + (557/1440.0)*n4
	h3 := (61/240.0)*n3 - (103/140.0)*n4
	h4 := (49561 / 161280.0) * n4
	e := math.Sqrt(s.E2)

	var M0 float64

	switch phi0 {
	case 0:
		M0 = 0
	case math.Pi / 2:
		M0 = B * (math.Pi / 2)
	case -math.Pi / 2:
		M0 = B * (-math.Pi / 2)
	default:
		Q0 := math.Asinh(math.Tan(phi0)) - (e * math.Atanh(e*math.Sin(phi0)))
		xi00 := math.Atan(math.Sinh(Q0))
		xi01 := h1 * math.Sin(2*xi00)
		xi02 := h2 * math.Sin(4*xi00)
		xi03 := h3 * math.Sin(6*xi00)
		xi04 := h4 * math.Sin(8*xi00)
		xi0 := xi00 + xi01 + xi02 + xi03 + xi04
		M0 = B * xi0
	}

	h1i := n/2.0 - (2/3.0)*n2 + (37/96.0)*n3 + (1/360.0)*n4
	h2i := (1/48.0)*n2 - (1/15.0)*n3 + (437/1440.0)*n4
	h3i := (17/480.0)*n3 - (37/840.0)*n4
	h4i := (4397 / 161280.0) * n4

	return transverseMercator{
		base:    base,
		lambdaO: lambda0,
		scale:   scale,
		eastf:   eastf,
		northf:  northf,
		b:       B,
		h1:      h1,
		h2:      h2,
		h3:      h3,
		h4:      h4,
		mO:      M0,
		h1i:     h1i,
		h2i:     h2i,
		h3i:     h3i,
		h4i:     h4i,
	}
}

type transverseMercator struct {
	base                  CRS
	lambdaO               float64
	b, h1, h2, h3, h4, mO float64
	h1i, h2i, h3i, h4i    float64
	scale                 float64
	eastf                 float64
	northf                float64
}

func (p transverseMercator) Base() CRS {
	return p.base
}

func (p transverseMercator) Spheroid() Spheroid {
	return p.base.Spheroid()
}

func (p transverseMercator) ToBase(east, north, h float64) (lon, lat, h2 float64) {
	s := p.base.Spheroid()

	etai := (east - p.eastf) / (p.b * p.scale)
	xii := ((north - p.northf) + p.scale*p.mO) / (p.b * p.scale)

	xi1i := p.h1i * math.Sin(2*xii) * math.Cosh(2*etai)
	xi2i := p.h2i * math.Sin(4*xii) * math.Cosh(4*etai)
	xi3i := p.h3i * math.Sin(6*xii) * math.Cosh(6*etai)
	xi4i := p.h4i * math.Sin(8*xii) * math.Cosh(8*etai)

	xi0i := xii - (xi1i + xi2i + xi3i + xi4i)

	eta1i := p.h1i * math.Cos(2*xii) * math.Sinh(2*etai)
	eta2i := p.h2i * math.Cos(4*xii) * math.Sinh(4*etai)
	eta3i := p.h3i * math.Cos(6*xii) * math.Sinh(6*etai)
	eta4i := p.h4i * math.Cos(8*xii) * math.Sinh(8*etai)

	eta0i := etai - (eta1i + eta2i + eta3i + eta4i)

	betai := math.Asin(math.Sin(xi0i) / math.Cosh(eta0i))
	Qi := math.Asinh(math.Tan(betai))
	Qii := Qi + (s.E * math.Atanh(s.E*math.Tanh(Qi)))

	for i := 0; i < 15; i++ {
		newQ := Qi + (s.E * math.Atanh(s.E*math.Tanh(Qii)))
		if math.Abs(newQ-Qii) < 1e-14 {
			break
		}

		Qii = newQ
	}

	phi := math.Atan(math.Sinh(Qii))
	lambda := p.lambdaO + math.Asin(math.Tanh(eta0i)/math.Cos(betai))

	return degree(lambda), degree(phi), h
}

func (p transverseMercator) FromBase(lon, lat, h float64) (east, north, h2 float64) {
	s := p.base.Spheroid()

	phi := radian(lat)
	lambda := radian(lon)

	Q := math.Asinh(math.Tan(phi)) - s.E*math.Atanh(s.E*math.Sin(phi))
	beta := math.Atan(math.Sinh(Q))
	eta0 := math.Atanh(math.Cos(beta) * math.Sin(lambda-p.lambdaO))
	xi0 := math.Asin(math.Sin(beta) * math.Cosh(eta0))

	xi1 := p.h1 * math.Sin(2*xi0) * math.Cosh(2*eta0)
	xi2 := p.h2 * math.Sin(4*xi0) * math.Cosh(4*eta0)
	xi3 := p.h3 * math.Sin(6*xi0) * math.Cosh(6*eta0)
	xi4 := p.h4 * math.Sin(8*xi0) * math.Cosh(8*eta0)

	xi := xi0 + xi1 + xi2 + xi3 + xi4

	eta1 := p.h1 * math.Cos(2*xi0) * math.Sinh(2*eta0)
	eta2 := p.h2 * math.Cos(4*xi0) * math.Sinh(4*eta0)
	eta3 := p.h3 * math.Cos(6*xi0) * math.Sinh(6*eta0)
	eta4 := p.h4 * math.Cos(8*xi0) * math.Sinh(8*eta0)

	eta := eta0 + eta1 + eta2 + eta3 + eta4

	east = p.eastf + p.scale*p.b*eta
	north = p.northf + p.scale*(p.b*xi-p.mO)

	return east, north, h
}

func LambertConformalConic2SP(base CRS, lonf, latf, sp1, sp2, eastf, northf float64) CRS {
	if base == nil {
		base = Geographic(nil, defaultSpheroid)
	}

	s := base.Spheroid()

	phif := radian(latf)
	phi1 := radian(sp1)
	phi2 := radian(sp2)
	lambdaf := radian(lonf)

	tf := math.Tan(math.Pi/4-phif/2) / math.Pow((1-s.E*math.Sin(phif))/(1+s.E*math.Sin(phif)), s.E/2)
	t1 := math.Tan(math.Pi/4-phi1/2) / math.Pow((1-s.E*math.Sin(phi1))/(1+s.E*math.Sin(phi1)), s.E/2)
	t2 := math.Tan(math.Pi/4-phi2/2) / math.Pow((1-s.E*math.Sin(phi2))/(1+s.E*math.Sin(phi2)), s.E/2)

	m1 := math.Cos(phi1) / math.Sqrt(1-s.E2*sin2(phi1))
	m2 := math.Cos(phi2) / math.Sqrt(1-s.E2*sin2(phi2))

	var n float64
	if math.Abs(phi1-phi2) < 1e-14 {
		n = math.Sin(phi1)
	} else {
		n = math.Log(m1/m2) / math.Log(t1/t2)
	}

	f := m1 / (n * math.Pow(t1, n))
	rf := s.A * f * math.Pow(tf, n)

	return lambertConformalConic{
		base:    base,
		lambdaf: lambdaf,
		n:       n,
		f:       f,
		rf:      rf,
		eastf:   eastf,
		northf:  northf,
	}
}

func LambertConformalConic1SP(base CRS, lonf, latf, scale, eastf, northf float64) CRS {
	if base == nil {
		base = Geographic(nil, defaultSpheroid)
	}

	s := base.Spheroid()

	phif := radian(latf)
	lambdaf := radian(lonf)

	tf := math.Tan(math.Pi/4-phif/2) / math.Pow((1-s.E*math.Sin(phif))/(1+s.E*math.Sin(phif)), s.E/2)
	m1 := math.Cos(phif) / math.Sqrt(1-s.E2*sin2(phif))

	n := math.Sin(phif)
	f := scale * m1 / (n * math.Pow(tf, n))
	rf := s.A * f * math.Pow(tf, n)

	return lambertConformalConic{
		base:    base,
		lambdaf: lambdaf,
		n:       n,
		f:       f,
		rf:      rf,
		eastf:   eastf,
		northf:  northf,
	}
}

// lambertConformalConic implements the Lambert Conformal Conic projection.
// The forward and inverse formulas are shared by the 1SP and 2SP variants,
// which differ only in how the constants n, f and rf are derived from the
// projection parameters (see the respective constructors).
type lambertConformalConic struct {
	base              CRS
	lambdaf, n, f, rf float64
	eastf             float64
	northf            float64
}

func (p lambertConformalConic) Base() CRS {
	return p.base
}

func (p lambertConformalConic) Spheroid() Spheroid {
	return p.base.Spheroid()
}

func (p lambertConformalConic) ToBase(east, north, h float64) (lon, lat, h2 float64) {
	s := p.base.Spheroid()

	ri := math.Hypot(east-p.eastf, p.rf-(north-p.northf))

	ti := math.Pow(ri/(s.A*p.f), 1/p.n)

	theta := math.Atan2(east-p.eastf, p.rf-(north-p.northf))

	phi := math.Pi/2 - 2*math.Atan(ti)

	for range 6 {
		next := math.Pi/2 - 2*math.Atan(ti*math.Pow((1-s.E*math.Sin(phi))/(1+s.E*math.Sin(phi)), s.E/2))
		if math.Abs(next-phi) < 1e-12 {
			phi = next
			break
		}

		phi = next
	}

	lambda := theta/p.n + p.lambdaf

	return degree(lambda), degree(phi), h
}

func (p lambertConformalConic) FromBase(lon, lat, h float64) (east, north, h2 float64) {
	s := p.base.Spheroid()

	phi := radian(lat)
	lambda := radian(lon)

	t := math.Tan(math.Pi/4-phi/2) / math.Pow((1-s.E*math.Sin(phi))/(1+s.E*math.Sin(phi)), s.E/2)

	r := s.A * p.f * math.Pow(t, p.n)
	theta := p.n * (lambda - p.lambdaf)

	east = p.eastf + r*math.Sin(theta)
	north = p.northf + p.rf - r*math.Cos(theta)

	return east, north, h
}

func AlbersConicEqualArea(base CRS, lonf, latf, sp1, sp2, eastf, northf float64) CRS {
	if base == nil {
		base = Geographic(nil, defaultSpheroid)
	}

	s := base.Spheroid()

	phif := radian(latf)
	phi1 := radian(sp1)
	phi2 := radian(sp2)
	lambdaf := radian(lonf)
	alphaf := (1 - s.E2) * ((math.Sin(phif) / (1 - s.E2*sin2(phif))) - (1/(2*s.E))*math.Log((1-s.E*math.Sin(phif))/(1+s.E*math.Sin(phif))))
	alpha1 := (1 - s.E2) * ((math.Sin(phi1) / (1 - s.E2*sin2(phi1))) - (1/(2*s.E))*math.Log((1-s.E*math.Sin(phi1))/(1+s.E*math.Sin(phi1))))
	alpha2 := (1 - s.E2) * ((math.Sin(phi2) / (1 - s.E2*sin2(phi2))) - (1/(2*s.E))*math.Log((1-s.E*math.Sin(phi2))/(1+s.E*math.Sin(phi2))))

	m1 := math.Cos(phi1) / math.Sqrt(1-s.E2*sin2(phi1))
	m2 := math.Cos(phi2) / math.Sqrt(1-s.E2*sin2(phi2))

	n := (m1*m1 - m2*m2) / (alpha2 - alpha1)
	c := m1*m1 + (n * alpha1)
	rf := (s.A * math.Sqrt(c-n*alphaf)) / n

	return albersConicEqualArea{
		base:    base,
		lambdaf: lambdaf,
		alphaf:  alphaf,
		n:       n,
		c:       c,
		rf:      rf,
		eastf:   eastf,
		northf:  northf,
	}
}

type albersConicEqualArea struct {
	base                      CRS
	lambdaf, alphaf, n, c, rf float64
	eastf                     float64
	northf                    float64
}

func (p albersConicEqualArea) Base() CRS {
	return p.base
}

func (p albersConicEqualArea) Spheroid() Spheroid {
	return p.base.Spheroid()
}

func (p albersConicEqualArea) ToBase(east, north, h float64) (lon, lat, h2 float64) {
	s := p.base.Spheroid()

	ri := math.Sqrt(intPow(east-p.eastf, 2) + intPow(p.rf-(north-p.northf), 2))
	alphai := (p.c - (intPow(ri, 2) * intPow(p.n, 2) / s.A2)) / p.n
	betai := math.Asin(alphai / (1 - ((1 - s.E2) / (2 * s.E) * math.Log((1-s.E)/(1+s.E)))))

	var theta float64

	if p.n > 0 {
		theta = math.Atan2((east - p.eastf), (p.rf - (north - p.northf)))
	} else {
		theta = math.Atan2(-(east - p.eastf), -(p.rf - (north - p.northf)))
	}

	phi := betai + ((s.E2/3 + 31*s.E4/180 + 517*s.E6/5040) * math.Sin(2*betai)) + ((23*s.E4/360 + 251*s.E6/3780) * math.Sin(4*betai)) + ((761 * s.E6 / 45360) * math.Sin(6*betai))
	lambda := p.lambdaf + (theta / p.n)

	return degree(lambda), degree(phi), h
}

func (p albersConicEqualArea) FromBase(lon, lat, h float64) (east, north, h2 float64) {
	s := p.base.Spheroid()

	lambda := radian(lon)
	phi := radian(lat)

	alpha := (1 - s.E2) * ((math.Sin(phi) / (1 - s.E2*sin2(phi))) - (1/(2*s.E))*math.Log((1-s.E*math.Sin(phi))/(1+s.E*math.Sin(phi))))

	theta := p.n * (lambda - p.lambdaf)
	r := (s.A * math.Sqrt(p.c-p.n*alpha)) / p.n
	rf := (s.A * math.Sqrt(p.c-p.n*p.alphaf)) / p.n

	east = p.eastf + r*math.Sin(theta)
	north = p.northf + rf - r*math.Cos(theta)

	return east, north, h
}

func LambertAzimuthalEqualArea(base CRS, lonf, latf, eastf, northf float64) CRS {
	if base == nil {
		base = Geographic(nil, defaultSpheroid)
	}

	s := base.Spheroid()

	phi0 := radian(latf)
	lambda0 := radian(lonf)

	q0 := (1 - s.E2) * ((math.Sin(phi0) / (1 - s.E2*sin2(phi0))) - (1 / (2 * s.E) * math.Log((1-s.E*math.Sin(phi0))/(1+s.E*math.Sin(phi0)))))
	qp := (1 - s.E2) * ((1 / (1 - s.E2)) - ((1 / (2 * s.E)) * math.Log((1-s.E)/(1+s.E))))

	beta0 := math.Asin(q0 / qp)
	rq := s.A * math.Sqrt(qp/2)
	G := s.A * (math.Cos(phi0) / math.Sqrt(1-s.E2*sin2(phi0))) / (rq * math.Cos(beta0))

	return lambertAzimuthalEqualArea{
		base:    base,
		phi0:    phi0,
		lambda0: lambda0,
		q0:      q0,
		qp:      qp,
		beta0:   beta0,
		rq:      rq,
		g:       G,
		eastf:   eastf,
		northf:  northf,
	}
}

type lambertAzimuthalEqualArea struct {
	base                                CRS
	phi0, lambda0, q0, qp, beta0, rq, g float64
	eastf                               float64
	northf                              float64
}

func (p lambertAzimuthalEqualArea) Base() CRS {
	return p.base
}

func (p lambertAzimuthalEqualArea) Spheroid() Spheroid {
	return p.base.Spheroid()
}

func (p lambertAzimuthalEqualArea) ToBase(east, north, h float64) (lon, lat, h2 float64) {
	s := p.base.Spheroid()

	rho := math.Sqrt(intPow((east-p.eastf)/p.g, 2) + intPow(p.g*(north-p.northf), 2))
	c := 2 * math.Asin(rho/(2*p.rq))
	betai := math.Asin((math.Cos(c) * math.Sin(p.beta0)) + ((p.g * (north - p.northf) * math.Sin(c) * math.Cos(p.beta0)) / rho))

	phi := betai + ((s.E2/3 + (31*s.E4)/180 + (517*s.E6)/5040) * math.Sin(2*betai)) +
		((23*s.E4)/360+(251*s.E6)/3780)*math.Sin(4*betai) +
		((761*s.E6)/45360)*math.Sin(6*betai)
	lambda := p.lambda0 + math.Atan2((east-p.eastf)*math.Sin(c), (p.g*rho*math.Cos(p.beta0)*math.Cos(c)-intPow(p.g, 2)*(north-p.northf)*math.Sin(p.beta0)*math.Sin(c)))

	return degree(lambda), degree(phi), h
}

func (p lambertAzimuthalEqualArea) FromBase(lon, lat, h float64) (east, north, h2 float64) {
	s := p.base.Spheroid()

	phi := radian(lat)
	lambda := radian(lon)

	q := (1 - s.E2) * ((math.Sin(phi) / (1 - s.E2*sin2(phi))) - (1 / (2 * s.E) * math.Log((1-s.E*math.Sin(phi))/(1+s.E*math.Sin(phi)))))

	beta0 := math.Asin(p.q0 / p.qp)
	rq := s.A * math.Sqrt(p.qp/2)
	g := s.A * (math.Cos(p.phi0) / math.Sqrt(1-s.E2*sin2(p.phi0))) / (rq * math.Cos(beta0))

	beta := math.Asin(q / p.qp)
	b := rq * math.Sqrt(2/(1+math.Sin(beta0)*math.Sin(beta)+(math.Cos(beta0)*math.Cos(beta)*math.Cos(lambda-p.lambda0))))

	east = p.eastf + ((b * g) * (math.Cos(beta) * math.Sin(lambda-p.lambda0)))
	north = p.northf + (b/g)*((math.Cos(beta0)*math.Sin(beta))-(math.Sin(beta0)*math.Cos(beta)*math.Cos(lambda-p.lambda0)))

	return east, north, h
}

func Krovak(base CRS, lonf, latf, azimuth, sp, scale, eastf, northf float64) CRS {
	if base == nil {
		base = Geographic(nil, defaultSpheroid)
	}

	s := base.Spheroid()

	phic := radian(latf)
	lambda0 := radian(lonf)
	phip := radian(sp)
	alphac := radian(azimuth)

	A := s.A * math.Sqrt(1-s.E2) / (1 - s.E2*sin2(phic))
	B := math.Sqrt(1 + (s.E2 * intPow(math.Cos(phic), 4) / (1 - s.E2)))
	gamma0 := math.Asin(math.Sin(phic) / B)
	t0 := math.Tan(math.Pi/4+gamma0/2) * math.Pow((1+s.E*math.Sin(phic))/(1-s.E*math.Sin(phic)), s.E*B/2) / math.Pow(math.Tan(math.Pi/4+phic/2), B)
	n := math.Sin(phip)
	r0 := scale * A / math.Tan(phip)

	return krovak{
		base:    base,
		lambda0: lambda0,
		phip:    phip,
		alphac:  alphac,
		b:       B,
		t0:      t0,
		n:       n,
		r0:      r0,
		eastf:   eastf,
		northf:  northf,
	}
}

type krovak struct {
	base                                CRS
	lambda0, phip, alphac, b, t0, n, r0 float64
	eastf                               float64
	northf                              float64
}

func (p krovak) Base() CRS {
	return p.base
}

func (p krovak) Spheroid() Spheroid {
	return p.base.Spheroid()
}

func (p krovak) FromBase(lon, lat, h float64) (east, north, h2 float64) {
	s := p.base.Spheroid()

	phi := radian(lat)
	lambda := radian(lon)

	U := 2 * (math.Atan(p.t0*math.Pow(math.Tan(phi/2+math.Pi/4), p.b)/math.Pow((1+s.E*math.Sin(phi))/(1-s.E*math.Sin(phi)), s.E*p.b/2)) - math.Pi/4)
	V := p.b * (p.lambda0 - lambda)
	T := math.Asin(math.Cos(p.alphac)*math.Sin(U) + math.Sin(p.alphac)*math.Cos(U)*math.Cos(V))
	D := math.Asin(math.Cos(U) * math.Sin(V) / math.Cos(T))
	theta := p.n * D
	r := p.r0 * math.Pow(math.Tan(math.Pi/4+p.phip/2), p.n) / math.Pow(math.Tan(T/2+math.Pi/4), p.n)
	Xp := r * math.Cos(theta)
	Yp := r * math.Sin(theta)

	return -(Yp + p.eastf), -(Xp + p.northf), h
}

func (p krovak) ToBase(east, north, h float64) (lon, lat, h2 float64) {
	s := p.base.Spheroid()

	Xpi := (-north) - p.northf
	Ypi := (-east) - p.eastf
	ri := math.Sqrt(intPow(Xpi, 2) + intPow(Ypi, 2))
	thetai := math.Atan2(Ypi, Xpi)
	di := thetai / math.Sin(p.phip)
	ti := 2 * (math.Atan(math.Pow(p.r0/ri, 1/p.n)*math.Tan(math.Pi/4+p.phip/2)) - math.Pi/4)
	ui := math.Asin(math.Cos(p.alphac)*math.Sin(ti) - math.Sin(p.alphac)*math.Cos(ti)*math.Cos(di))
	vi := math.Asin(math.Cos(ti) * math.Sin(di) / math.Cos(ui))

	phi := ui

	for range 3 {
		phi = 2 * (math.Atan(math.Pow(p.t0, -1/p.b)*math.Pow(math.Tan(ui/2+math.Pi/4), 1/p.b)*math.Pow((1+s.E*math.Sin(phi))/(1-s.E*math.Sin(phi)), s.E/2)) - math.Pi/4)
	}

	lambda := p.lambda0 - vi/p.b

	return degree(lambda), degree(phi), h
}

// -------------------------
// helpers
// -------------------------

func sign(x float64) float64 {
	if x < 0 {
		return -1
	}
	return 1
}

func clamp(x, min, max float64) float64 {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}

func HotineObliqueMercatorAzimuthCenter(
	base CRS,
	lonc, latc, azimuth, gammac, k0, eastf, northf float64,
) CRS {
	if base == nil {
		base = Geographic(nil, defaultSpheroid)
	}

	s := base.Spheroid()
	a := s.A
	e := s.E
	e2 := s.E2

	phiC := radian(latc)
	lambdaC := radian(lonc)
	alphaC := radian(azimuth)
	gammaC := radian(gammac)

	sinφ := math.Sin(phiC)
	cosφ := math.Cos(phiC)

	sLat := sign(sinφ)

	B := math.Sqrt(1 + (e2*math.Pow(cosφ, 4))/(1-e2))

	A := a * k0 * math.Sqrt(1-e2) * B /
		(1 - e2*sinφ*sinφ)

	t0 := math.Tan(math.Pi/4-phiC/2) /
		math.Pow((1-e*sinφ)/(1+e*sinφ), e/2)

	D := (B * math.Sqrt(1-e2)) /
		(cosφ * math.Sqrt(1-e2*sinφ*sinφ))

	D2 := D * D
	if D < 1 {
		D2 = 1
	}

	sqrtTerm := math.Sqrt(D2 - 1)

	F := D + sqrtTerm*sLat
	G := (F - 1/F) / 2

	x := math.Sin(alphaC) / D
	x = clamp(x, -1, 1)

	gamma0 := math.Asin(x)

	lonArg := clamp(G*math.Tan(gamma0), -1, 1)
	lon0 := lambdaC - math.Asin(lonArg)/B

	H := F * math.Pow(t0, B)

	var uc float64
	if math.Abs(math.Abs(alphaC)-math.Pi/2) < 1e-12 {
		uc = A * (lambdaC - lon0)
	} else {
		uc = (A / B) * math.Atan2(sqrtTerm, math.Cos(alphaC)) * sLat
	}

	return hotineObliqueMercatorAzimuthCenter{
		base:   base,
		A:      A,
		B:      B,
		t0:     t0,
		H:      H,
		D:      D,
		D2:     D2,
		F:      F,
		G:      G,
		gamma0: gamma0,
		lon0:   lon0,
		uc:     uc,
		eastf:  eastf,
		northf: northf,
		alphaC: alphaC,
		gammaC: gammaC,
		sLat:   sLat,
	}
}

type hotineObliqueMercatorAzimuthCenter struct {
	base CRS

	A, B  float64
	t0, H float64
	D, D2 float64
	F, G  float64

	gamma0 float64
	lon0   float64
	uc     float64

	eastf, northf float64
	alphaC        float64
	gammaC        float64
	sLat          float64
}

func (h hotineObliqueMercatorAzimuthCenter) Base() CRS {
	return h.base
}

func (h hotineObliqueMercatorAzimuthCenter) Spheroid() Spheroid {
	return h.base.Spheroid()
}

func (h hotineObliqueMercatorAzimuthCenter) FromBase(lon, lat, _ float64) (float64, float64, float64) {

	s := h.base.Spheroid()
	e := s.E

	phi := radian(lat)
	lambda := radian(lon)

	t := math.Tan(math.Pi/4-phi/2) /
		math.Pow((1-e*math.Sin(phi))/(1+e*math.Sin(phi)), e/2)

	Q := h.H / math.Pow(t, h.B)

	S := (Q - 1/Q) / 2
	T := (Q + 1/Q) / 2

	V := math.Sin(h.B * (lambda - h.lon0))
	U := (-V*math.Cos(h.gamma0) + S*math.Sin(h.gamma0)) / T

	v := (h.A / (2 * h.B)) * math.Log((1-U)/(1+U))

	u0 := (h.A / h.B) * math.Atan2(
		S*math.Cos(h.gamma0)+V*math.Sin(h.gamma0),
		math.Cos(h.B*(lambda-h.lon0)),
	)

	u := u0 - math.Abs(h.uc)*h.sLat*sign(lambda-h.lon0)

	E := v*math.Cos(h.gammaC) + u*math.Sin(h.gammaC) + h.eastf
	N := u*math.Cos(h.gammaC) - v*math.Sin(h.gammaC) + h.northf

	return E, N, 0
}

func (h hotineObliqueMercatorAzimuthCenter) ToBase(E, N, _ float64) (float64, float64, float64) {
	s := h.base.Spheroid()
	e := s.E

	// ----------------------------
	// 1. inverse rotation (PROJ-style)
	// ----------------------------
	v := (E-h.eastf)*math.Cos(h.gammaC) -
		(N-h.northf)*math.Sin(h.gammaC)

	u := (N-h.northf)*math.Cos(h.gammaC) +
		(E-h.eastf)*math.Sin(h.gammaC)

	// IMPORTANT: undo uc shift correctly (no double-signing here)
	u += math.Abs(h.uc) * h.sLat

	// ----------------------------
	// 2. conformal part
	// ----------------------------
	Qp := math.Exp(-(h.B * v / h.A))

	Sp := 0.5 * (Qp - 1/Qp)
	Tp := 0.5 * (Qp + 1/Qp)

	Vp := math.Sin(h.B * u / h.A)

	Up := (Vp*math.Cos(h.gamma0) + Sp*math.Sin(h.gamma0)) / Tp

	// ----------------------------
	// 3. latitude
	// ----------------------------
	tp := h.H / math.Pow(math.Sqrt((1+Up)/(1-Up)), 1.0/h.B)

	chi := math.Pi/2 - 2*math.Atan(tp)

	lat := chi +
		math.Sin(2*chi)*(e*e/2+5*math.Pow(e, 4)/24+math.Pow(e, 6)/12+13*math.Pow(e, 8)/360) +
		math.Sin(4*chi)*(7*math.Pow(e, 4)/48+29*math.Pow(e, 6)/240+811*math.Pow(e, 8)/11520) +
		math.Sin(6*chi)*(7*math.Pow(e, 6)/120+81*math.Pow(e, 8)/1120) +
		math.Sin(8*chi)*(4279*math.Pow(e, 8)/161280)

	// ----------------------------
	// 4. longitude (THIS is the correct PROJ form)
	// ----------------------------
	lon := h.lon0 -
		(1.0/h.B)*math.Atan2(
			Sp*math.Cos(h.gamma0)-Vp*math.Sin(h.gamma0),
			math.Cos(h.B*u/h.A),
		)

	return degree(lon), degree(lat), 0
}

func sin2(r float64) float64 {
	return intPow(math.Sin(r), 2)
}

func degree(r float64) float64 {
	return r * 180 / math.Pi
}

func radian(g float64) float64 {
	return g * math.Pi / 180
}

func intPow(val float64, times int) float64 {
	result := 1.0

	for range times {
		result *= val
	}

	return result
}
