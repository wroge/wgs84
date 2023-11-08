//nolint:varnamelen,nonamedreturns,gomnd,gochecknoglobals,exhaustivestruct,exhaustruct,lll
package wgs84

import (
	"math"
)

var (
	GRS80             = NewSpheroid(6378137, 298.257222101)
	Airy1830          = NewSpheroid(6377563.396, 299.3249646)
	AiryModified1849  = NewSpheroid(6377340.189, 299.3249646)
	Bessel1841        = NewSpheroid(6377397.155, 299.1528128)
	Clarke1866        = NewSpheroid(6378206.4, 294.9786982139006)
	International1924 = NewSpheroid(6378388, 297)
	CGCS2000          = NewSpheroid(6378137, 298.257222101)
	WGS72             = NewSpheroid(6378135, 298.26)

	Datum = Geographic{
		Spheroid: NewSpheroid(6378137, 298.257223563),
	}
	SJTSK = Geographic{
		Spheroid: Bessel1841,
		Transformation: Helmert{
			Tx: 589,
			Ty: 76,
			Tz: 480,
		},
	}
	RGF93 = Geographic{
		Spheroid: GRS80,
	}
	ETRS89 = Geographic{
		Spheroid:       GRS80,
		Transformation: nil,
	}
	NAD83 = Geographic{
		Spheroid: GRS80,
	}
	OSGB36 = Geographic{
		Spheroid: Airy1830,
		Transformation: Helmert{
			Tx: 446.448,
			Ty: -125.157,
			Tz: 542.06,
			Rx: 0.15,
			Ry: 0.247,
			Rz: 0.842,
			Ds: -20.489,
		},
	}
	MGI = Geographic{
		Spheroid: Bessel1841,
		Transformation: Helmert{
			Tx: 577.326,
			Ty: 90.129,
			Tz: 463.919,
			Rx: 5.137,
			Ry: 1.474,
			Rz: 5.297,
			Ds: 2.4232,
		},
	}
	DHDN2001 = Geographic{
		Spheroid: Bessel1841,
		Transformation: Helmert{
			Tx: 598.1,
			Ty: 73.7,
			Tz: 418.2,
			Rx: 0.202,
			Ry: 0.045,
			Rz: -2.455,
			Ds: 6.7,
		},
	}
	OSNI1952 = Geographic{
		Spheroid: Airy1830,
		Transformation: Helmert{
			Tx: 482.5,
			Ty: -130.6,
			Tz: 564.6,
			Rx: -1.042,
			Ry: -0.214,
			Rz: -0.631,
			Ds: 8.15,
		},
	}
	IRENET95 = Geographic{
		Spheroid: GRS80,
	}
	TM65 = Geographic{
		Spheroid: AiryModified1849,
		Transformation: Helmert{
			Tx: 482.5,
			Ty: -130.6,
			Tz: 564.6,
			Rx: -1.042,
			Ry: -0.214,
			Rz: -0.631,
			Ds: 8.15,
		},
	}
	TM75 = Geographic{
		Spheroid: AiryModified1849,
		Transformation: Helmert{
			Tx: 482.5,
			Ty: -130.6,
			Tz: 564.6,
			Rx: -1.042,
			Ry: -0.214,
			Rz: -0.631,
			Ds: 8.15,
		},
	}
	ED50 = Geographic{
		Spheroid: International1924,
		Transformation: Helmert{
			Tx: -87,
			Ty: -98,
			Tz: -121,
		},
	}
)

func Transform(from, to CoordinateReferenceSystem) Func {
	return func(a, b, c float64) (a2 float64, b2 float64, c2 float64) {
		if from != nil {
			a, b, c = from.ToWGS84(a, b, c)
		}

		if to == nil {
			return a, b, c
		}

		return to.FromWGS84(a, b, c)
	}
}

type Func func(a float64, b float64, c float64) (a2 float64, b2 float64, c2 float64)

type CoordinateReferenceSystem interface {
	FromWGS84(x0, y0, z0 float64) (x, y, z float64)
	ToWGS84(x, y, z float64) (x0, y0, z0 float64)
}

type root struct{}

func (t root) ToWGS84(x, y, z float64) (x0, y0, z0 float64) {
	return x, y, z
}

func (t root) FromWGS84(x0, y0, z0 float64) (x, y, z float64) {
	return x0, y0, z0
}

type Helmert struct {
	Tx, Ty, Tz, Rx, Ry, Rz, Ds float64
}

func (t Helmert) ToWGS84(x, y, z float64) (x0, y0, z0 float64) {
	return calcHelmert(x, y, z, t.Tx, t.Ty, t.Tz, t.Rx, t.Ry, t.Rz, t.Ds)
}

func (t Helmert) FromWGS84(x0, y0, z0 float64) (x, y, z float64) {
	return calcHelmert(x0, y0, z0, -t.Tx, -t.Ty, -t.Tz, -t.Rx, -t.Ry, -t.Rz, -t.Ds)
}

func (t Helmert) Spheroid(a, fi float64) Geographic {
	return Geographic{
		Transformation: t,
		Spheroid:       NewSpheroid(a, fi),
	}
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

type Geographic struct {
	Transformation CoordinateReferenceSystem
	Spheroid       Spheroid
}

func (g Geographic) ToWGS84(lon, lat, h float64) (x0, y0, z0 float64) {
	n := g.Spheroid.A / math.Sqrt(1-g.Spheroid.E2*math.Pow(math.Sin(radian(lat)), 2))

	x := (n + h) * math.Cos(radian(lon)) * math.Cos(radian(lat))
	y := (n + h) * math.Cos(radian(lat)) * math.Sin(radian(lon))
	z := (n*math.Pow(g.Spheroid.A*(1-g.Spheroid.F), 2)/(g.Spheroid.A2) + h) * math.Sin(radian(lat))

	if g.Transformation == nil {
		return x, y, z
	}

	return g.Transformation.ToWGS84(x, y, z)
}

func (g Geographic) FromWGS84(x0, y0, z0 float64) (lon, lat, h float64) {
	if g.Transformation != nil {
		x0, y0, z0 = g.Transformation.FromWGS84(x0, y0, z0)
	}

	sd := math.Sqrt(x0*x0 + y0*y0)
	T := math.Atan(z0 * g.Spheroid.A / (sd * g.Spheroid.B))
	B := math.Atan((z0 + g.Spheroid.E2*(g.Spheroid.A2)/g.Spheroid.B*
		math.Pow(math.Sin(T), 3)) / (sd - g.Spheroid.E2*g.Spheroid.A*math.Pow(math.Cos(T), 3)))
	n := g.Spheroid.A / math.Sqrt(1-g.Spheroid.E2*math.Pow(math.Sin(B), 2))
	h = sd/math.Cos(B) - n
	lon = degree(math.Atan2(y0, x0))
	lat = degree(B)

	return lon, lat, h
}

func (g Geographic) WebMercator() WebMercator {
	return WebMercator{geogr: g}
}

func (g Geographic) TransverseMercator(lonf, latf, scale, eastf, northf float64) TransverseMercator {
	phi0 := radian(latf)
	lambda0 := radian(lonf)
	n := g.Spheroid.F / (2 - g.Spheroid.F)
	n2 := math.Pow(n, 2)
	n3 := math.Pow(n, 3)
	n4 := math.Pow(n, 4)
	B := (g.Spheroid.A / (1 + n)) * (1 + n2/4 + n4/64)
	h1 := n/2.0 - (2/3.0)*n2 + (5/16.0)*n3 + (41/180.0)*n4
	h2 := (13/48.0)*n2 - (3/5.0)*n3 + (557/1440.0)*n4
	h3 := (61/240.0)*n3 - (103/140.0)*n4
	h4 := (49561 / 161280.0) * n4
	e := math.Sqrt(g.Spheroid.E2)

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

	return TransverseMercator{
		geogr:   g,
		phiO:    phi0,
		lambdaO: lambda0,
		scale:   scale,
		eastf:   eastf,
		northf:  northf,
		n:       n,
		n2:      n2,
		n3:      n3,
		n4:      n4,
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

func (g Geographic) LambertConformalConic2SP(lonf, latf, sp1, sp2, eastf, northf float64) LambertConformalConic2SP {
	phif := radian(latf)
	phi1 := radian(sp1)
	phi2 := radian(sp2)
	lambdaf := radian(lonf)

	tf := math.Tan(math.Pi/4-phif/2) / math.Pow((1-g.Spheroid.E*math.Sin(phif))/(1+g.Spheroid.E*math.Sin(phif)), g.Spheroid.E/2)
	t1 := math.Tan(math.Pi/4-phi1/2) / math.Pow((1-g.Spheroid.E*math.Sin(phi1))/(1+g.Spheroid.E*math.Sin(phi1)), g.Spheroid.E/2)
	t2 := math.Tan(math.Pi/4-phi2/2) / math.Pow((1-g.Spheroid.E*math.Sin(phi2))/(1+g.Spheroid.E*math.Sin(phi2)), g.Spheroid.E/2)

	m1 := math.Cos(phi1) / math.Sqrt(1-g.Spheroid.E2*sin2(phi1))
	m2 := math.Cos(phi2) / math.Sqrt(1-g.Spheroid.E2*sin2(phi2))

	n := (math.Log(m1) - math.Log(m2)) / (math.Log(t1) - math.Log(t2))
	f := m1 / (n * math.Pow(t1, n))
	rf := g.Spheroid.A * f * math.Pow(tf, n)

	return LambertConformalConic2SP{
		geogr:   g,
		phif:    phif,
		phi1:    phi1,
		phi2:    phi2,
		lambdaf: lambdaf,
		tf:      tf,
		t1:      t1,
		t2:      t2,
		m1:      m1,
		m2:      m2,
		n:       n,
		f:       f,
		rf:      rf,
		eastf:   eastf,
		northf:  northf,
	}
}

func (g Geographic) AlbersConicEqualArea(lonf, latf, sp1, sp2, eastf, northf float64) AlbersConicEqualArea {
	phif := radian(latf)
	phi1 := radian(sp1)
	phi2 := radian(sp2)
	lambdaf := radian(lonf)
	alphaf := (1 - g.Spheroid.E2) * ((math.Sin(phif) / (1 - g.Spheroid.E2*sin2(phif))) - (1/(2*g.Spheroid.E))*math.Log((1-g.Spheroid.E*math.Sin(phif))/(1+g.Spheroid.E*math.Sin(phif))))
	alpha1 := (1 - g.Spheroid.E2) * ((math.Sin(phi1) / (1 - g.Spheroid.E2*sin2(phi1))) - (1/(2*g.Spheroid.E))*math.Log((1-g.Spheroid.E*math.Sin(phi1))/(1+g.Spheroid.E*math.Sin(phi1))))
	alpha2 := (1 - g.Spheroid.E2) * ((math.Sin(phi2) / (1 - g.Spheroid.E2*sin2(phi2))) - (1/(2*g.Spheroid.E))*math.Log((1-g.Spheroid.E*math.Sin(phi2))/(1+g.Spheroid.E*math.Sin(phi2))))

	m1 := math.Cos(phi1) / math.Sqrt(1-g.Spheroid.E2*sin2(phi1))
	m2 := math.Cos(phi2) / math.Sqrt(1-g.Spheroid.E2*sin2(phi2))

	n := (math.Pow(m1, 2) - math.Pow(m2, 2)) / (alpha2 - alpha1)
	c := math.Pow(m1, 2) + (n * alpha1)
	rf := (g.Spheroid.A * math.Sqrt(c-n*alphaf)) / n

	return AlbersConicEqualArea{
		geogr:   g,
		phif:    phif,
		phi1:    phi1,
		phi2:    phi2,
		lambdaf: lambdaf,
		alphaf:  alphaf,
		alpha1:  alpha1,
		alpha2:  alpha2,
		m1:      m1,
		m2:      m2,
		n:       n,
		c:       c,
		rf:      rf,
		eastf:   eastf,
		northf:  northf,
	}
}

func (g Geographic) LambertAzimuthalEqualArea(lonf, latf, eastf, northf float64) LambertAzimuthalEqualArea {
	phi0 := radian(latf)
	lambda0 := radian(lonf)

	q0 := (1 - g.Spheroid.E2) * ((math.Sin(phi0) / (1 - g.Spheroid.E2*sin2(phi0))) - (1 / (2 * g.Spheroid.E) * math.Log((1-g.Spheroid.E*math.Sin(phi0))/(1+g.Spheroid.E*math.Sin(phi0)))))
	qp := (1 - g.Spheroid.E2) * ((1 / (1 - g.Spheroid.E2)) - ((1 / (2 * g.Spheroid.E)) * math.Log((1-g.Spheroid.E)/(1+g.Spheroid.E))))

	beta0 := math.Asin(q0 / qp)
	rq := g.Spheroid.A * math.Sqrt(qp/2)
	G := g.Spheroid.A * (math.Cos(phi0) / math.Sqrt(1-g.Spheroid.E2*sin2(phi0))) / (rq * math.Cos(beta0))

	return LambertAzimuthalEqualArea{
		geogr:   g,
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

func (g Geographic) Krovak(lonf, latf, azimuth, sp, scale, eastf, northf float64) Krovak {
	phic := radian(latf)
	lambda0 := radian(lonf)
	phip := radian(sp)
	alphac := radian(azimuth)

	A := g.Spheroid.A * math.Sqrt(1-g.Spheroid.E2) / (1 - g.Spheroid.E2*sin2(phic))
	B := math.Sqrt(1 + (g.Spheroid.E2 * math.Pow(math.Cos(phic), 4) / (1 - g.Spheroid.E2)))
	gamma0 := math.Asin(math.Sin(phic) / B)
	t0 := math.Tan(math.Pi/4+gamma0/2) * math.Pow((1+g.Spheroid.E*math.Sin(phic))/(1-g.Spheroid.E*math.Sin(phic)), g.Spheroid.E*B/2) / math.Pow(math.Tan(math.Pi/4+phic/2), B)
	n := math.Sin(phip)
	r0 := scale * A / math.Tan(phip)

	return Krovak{
		geogr:   g,
		phic:    phic,
		lambda0: lambda0,
		phip:    phip,
		alphac:  alphac,
		a:       A,
		b:       B,
		gamma0:  gamma0,
		t0:      t0,
		n:       n,
		r0:      r0,
		eastf:   eastf,
		northf:  northf,
	}
}

type WebMercator struct {
	geogr Geographic
}

func (p WebMercator) ToWGS84(east, north, h float64) (x0, y0, z0 float64) {
	D := (-north) / p.geogr.Spheroid.A
	phi := math.Pi/2 - 2*math.Atan(math.Pow(math.E, D))
	lambda := east / p.geogr.Spheroid.A

	return p.geogr.ToWGS84(radian(lambda), radian(phi), h)
}

func (p WebMercator) FromWGS84(x0, y0, z0 float64) (east, north, h float64) {
	lon, lat, h := p.geogr.FromWGS84(x0, y0, z0)

	lambda := radian(lon)
	phi := radian(lat)

	east = p.geogr.Spheroid.A * lambda
	north = p.geogr.Spheroid.A * math.Log(math.Tan(math.Pi/4+phi/2))

	return east, north, h
}

type TransverseMercator struct {
	geogr                                Geographic
	phiO                                 float64
	lambdaO                              float64
	n, n2, n3, n4, b, h1, h2, h3, h4, mO float64
	h1i, h2i, h3i, h4i                   float64
	scale                                float64
	eastf                                float64
	northf                               float64
}

func (p TransverseMercator) ToWGS84(east, north, h float64) (x0, y0, z0 float64) {
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
	Qii := Qi + (p.geogr.Spheroid.E * math.Atanh(p.geogr.Spheroid.E*math.Tanh(Qi)))

	for i := 0; i < 15; i++ {
		q := Qi + (p.geogr.Spheroid.E * math.Atanh(p.geogr.Spheroid.E*math.Tanh(Qii)))

		Qii = q
	}

	phi := math.Atan(math.Sinh(Qii))
	lambda := p.lambdaO + math.Asin(math.Tanh(eta0i)/math.Cos(betai))

	return p.geogr.ToWGS84(degree(lambda), degree(phi), h)
}

func (p TransverseMercator) FromWGS84(x0, y0, z0 float64) (east, north, h float64) {
	lon, lat, h := p.geogr.FromWGS84(x0, y0, z0)

	phi := radian(lat)
	lambda := radian(lon)

	Q := math.Asinh(math.Tan(phi)) - p.geogr.Spheroid.E*math.Atanh(p.geogr.Spheroid.E*math.Sin(phi))
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

type LambertConformalConic2SP struct {
	geogr                                                   Geographic
	phif, phi1, phi2, lambdaf, tf, t1, t2, m1, m2, n, f, rf float64
	eastf                                                   float64
	northf                                                  float64
}

func (p LambertConformalConic2SP) ToWGS84(east, north, h float64) (x0, y0, z0 float64) {
	ri := math.Sqrt(math.Pow(east-p.eastf, 2) + math.Pow(p.rf-(north-p.northf), 2))
	if p.n < 0 && ri > 0 {
		ri = -ri
	}

	ti := math.Pow(ri/(p.geogr.Spheroid.A*p.f), 1/p.n)

	var theta float64
	if p.n > 0 {
		theta = math.Atan2((east - p.eastf), (p.rf - (north - p.northf)))
	} else {
		theta = math.Atan2(-(east - p.eastf), -(p.rf - (north - p.northf)))
	}

	phi := math.Pi/2 - 2*math.Atan(ti)

	for i := 0; i < 5; i++ {
		phi = math.Pi/2 - 2*math.Atan(ti*math.Pow((1-p.geogr.Spheroid.E*math.Sin(phi))/(1+p.geogr.Spheroid.E*math.Sin(phi)), p.geogr.Spheroid.E/2))
	}

	lambda := theta/p.n + p.lambdaf

	return p.geogr.ToWGS84(degree(lambda), degree(phi), h)
}

func (p LambertConformalConic2SP) FromWGS84(x0, y0, z0 float64) (east, north, h float64) {
	lon, lat, h := p.geogr.FromWGS84(x0, y0, z0)

	phi := radian(lat)
	lambda := radian(lon)

	t := math.Tan(math.Pi/4-phi/2) / math.Pow((1-p.geogr.Spheroid.E*math.Sin(phi))/(1+p.geogr.Spheroid.E*math.Sin(phi)), p.geogr.Spheroid.E/2)
	tf := math.Tan(math.Pi/4-p.phif/2) / math.Pow((1-p.geogr.Spheroid.E*math.Sin(p.phif))/(1+p.geogr.Spheroid.E*math.Sin(p.phif)), p.geogr.Spheroid.E/2)
	t1 := math.Tan(math.Pi/4-p.phi1/2) / math.Pow((1-p.geogr.Spheroid.E*math.Sin(p.phi1))/(1+p.geogr.Spheroid.E*math.Sin(p.phi1)), p.geogr.Spheroid.E/2)
	t2 := math.Tan(math.Pi/4-p.phi2/2) / math.Pow((1-p.geogr.Spheroid.E*math.Sin(p.phi2))/(1+p.geogr.Spheroid.E*math.Sin(p.phi2)), p.geogr.Spheroid.E/2)

	m1 := math.Cos(p.phi1) / math.Sqrt(1-p.geogr.Spheroid.E2*sin2(p.phi1))
	m2 := math.Cos(p.phi2) / math.Sqrt(1-p.geogr.Spheroid.E2*sin2(p.phi2))

	n := (math.Log(m1) - math.Log(m2)) / (math.Log(t1) - math.Log(t2))
	F := m1 / (n * math.Pow(t1, n))
	r := p.geogr.Spheroid.A * F * math.Pow(t, n)
	rf := p.geogr.Spheroid.A * F * math.Pow(tf, n)
	theta := n * (lambda - p.lambdaf)

	east = p.eastf + r*math.Sin(theta)
	north = p.northf + rf - r*math.Cos(theta)

	return east, north, h
}

type AlbersConicEqualArea struct {
	geogr                                                               Geographic
	lambdaf, phif, phi1, phi2, alphaf, alpha1, alpha2, m1, m2, n, c, rf float64
	eastf                                                               float64
	northf                                                              float64
}

func (p AlbersConicEqualArea) ToWGS84(east, north, h float64) (x0, y0, z0 float64) {
	ri := math.Sqrt(math.Pow(east-p.eastf, 2) + math.Pow(p.rf-(north-p.northf), 2))
	alphai := (p.c - (math.Pow(ri, 2) * math.Pow(p.n, 2) / p.geogr.Spheroid.A2)) / p.n
	betai := math.Asin(alphai / (1 - ((1 - p.geogr.Spheroid.E2) / (2 * p.geogr.Spheroid.E) * math.Log((1-p.geogr.Spheroid.E)/(1+p.geogr.Spheroid.E)))))

	var theta float64

	if p.n > 0 {
		theta = math.Atan2((east - p.eastf), (p.rf - (north - p.northf)))
	} else {
		theta = math.Atan2(-(east - p.eastf), -(p.rf - (north - p.northf)))
	}

	phi := betai + ((p.geogr.Spheroid.E2/3 + 31*p.geogr.Spheroid.E4/180 + 517*p.geogr.Spheroid.E6/5040) * math.Sin(2*betai)) + ((23*p.geogr.Spheroid.E4/360 + 251*p.geogr.Spheroid.E6/3780) * math.Sin(4*betai)) + ((761 * p.geogr.Spheroid.E6 / 45360) * math.Sin(6*betai))
	lambda := p.lambdaf + (theta / p.n)

	return p.geogr.ToWGS84(degree(lambda), degree(phi), h)
}

func (p AlbersConicEqualArea) FromWGS84(x0, y0, z0 float64) (east, north, h float64) {
	lon, lat, h := p.geogr.FromWGS84(x0, y0, z0)

	lambda := radian(lon)
	phi := radian(lat)

	alpha := (1 - p.geogr.Spheroid.E2) * ((math.Sin(phi) / (1 - p.geogr.Spheroid.E2*sin2(phi))) - (1/(2*p.geogr.Spheroid.E))*math.Log((1-p.geogr.Spheroid.E*math.Sin(phi))/(1+p.geogr.Spheroid.E*math.Sin(phi))))

	theta := p.n * (lambda - p.lambdaf)
	r := (p.geogr.Spheroid.A * math.Sqrt(p.c-p.n*alpha)) / p.n
	rf := (p.geogr.Spheroid.A * math.Sqrt(p.c-p.n*p.alphaf)) / p.n

	east = p.eastf + r*math.Sin(theta)
	north = p.northf + rf - r*math.Cos(theta)

	return east, north, h
}

type LambertAzimuthalEqualArea struct {
	geogr                               Geographic
	phi0, lambda0, q0, qp, beta0, rq, g float64
	eastf                               float64
	northf                              float64
}

func (p LambertAzimuthalEqualArea) ToWGS84(east, north, h float64) (x0, y0, z0 float64) {
	rho := math.Sqrt(math.Pow((east-p.eastf)/p.g, 2) + math.Pow(p.g*(north-p.northf), 2))
	c := 2 * math.Asin(rho/(2*p.rq))
	betai := math.Asin((math.Cos(c) * math.Sin(p.beta0)) + ((p.g * (north - p.northf) * math.Sin(c) * math.Cos(p.beta0)) / rho))

	phi := betai + ((p.geogr.Spheroid.E2/3 + (31*p.geogr.Spheroid.E4)/180 + (517*p.geogr.Spheroid.E6)/5040) * math.Sin(2*betai)) +
		((23*p.geogr.Spheroid.E4)/360+(251*p.geogr.Spheroid.E6)/3780)*math.Sin(4*betai) +
		((761*p.geogr.Spheroid.E6)/45360)*math.Sin(6*betai)
	lambda := p.lambda0 + math.Atan2((east-p.eastf)*math.Sin(c), (p.g*rho*math.Cos(p.beta0)*math.Cos(c)-math.Pow(p.g, 2)*(north-p.northf)*math.Sin(p.beta0)*math.Sin(c)))

	return p.geogr.ToWGS84(degree(lambda), degree(phi), h)
}

func (p LambertAzimuthalEqualArea) FromWGS84(x0, y0, z0 float64) (east, north, h float64) {
	lon, lat, h := p.geogr.FromWGS84(x0, y0, z0)

	phi := radian(lat)
	lambda := radian(lon)

	q := (1 - p.geogr.Spheroid.E2) * ((math.Sin(phi) / (1 - p.geogr.Spheroid.E2*sin2(phi))) - (1 / (2 * p.geogr.Spheroid.E) * math.Log((1-p.geogr.Spheroid.E*math.Sin(phi))/(1+p.geogr.Spheroid.E*math.Sin(phi)))))

	beta0 := math.Asin(p.q0 / p.qp)
	rq := p.geogr.Spheroid.A * math.Sqrt(p.qp/2)
	g := p.geogr.Spheroid.A * (math.Cos(p.phi0) / math.Sqrt(1-p.geogr.Spheroid.E2*sin2(p.phi0))) / (rq * math.Cos(beta0))

	beta := math.Asin(q / p.qp)
	b := rq * math.Sqrt(2/(1+math.Sin(beta0)*math.Sin(beta)+(math.Cos(beta0)*math.Cos(beta)*math.Cos(lambda-p.lambda0))))

	east = p.eastf + ((b * g) * (math.Cos(beta) * math.Sin(lambda-p.lambda0)))
	north = p.northf + (b/g)*((math.Cos(beta0)*math.Sin(beta))-(math.Sin(beta0)*math.Cos(beta)*math.Cos(lambda-p.lambda0)))

	return east, north, h
}

type Krovak struct {
	geogr                                                Geographic
	phic, lambda0, phip, alphac, a, b, gamma0, t0, n, r0 float64
	eastf                                                float64
	northf                                               float64
}

func (p Krovak) FromWGS84(x0, y0, z0 float64) (east, north, h float64) {
	lon, lat, h := p.geogr.FromWGS84(x0, y0, z0)

	phi := radian(lat)
	lambda := radian(lon)

	U := 2 * (math.Atan(p.t0*math.Pow(math.Tan(phi/2+math.Pi/4), p.b)/math.Pow((1+p.geogr.Spheroid.E*math.Sin(phi))/(1-p.geogr.Spheroid.E*math.Sin(phi)), p.geogr.Spheroid.E*p.b/2)) - math.Pi/4)
	V := p.b * (p.lambda0 - lambda)
	T := math.Asin(math.Cos(p.alphac)*math.Sin(U) + math.Sin(p.alphac)*math.Cos(U)*math.Cos(V))
	D := math.Asin(math.Cos(U) * math.Sin(V) / math.Cos(T))
	theta := p.n * D
	r := p.r0 * math.Pow(math.Tan(math.Pi/4+p.phip/2), p.n) / math.Pow(math.Tan(T/2+math.Pi/4), p.n)
	Xp := r * math.Cos(theta)
	Yp := r * math.Sin(theta)

	return -(Yp + p.eastf), -(Xp + p.northf), h
}

func (p Krovak) ToWGS84(east, north, h float64) (x0, y0, z0 float64) {
	Xpi := (-north) - p.northf
	Ypi := (-east) - p.eastf
	ri := math.Sqrt(math.Pow(Xpi, 2) + math.Pow(Ypi, 2))
	thetai := math.Atan2(Ypi, Xpi)
	di := thetai / math.Sin(p.phip)
	ti := 2 * (math.Atan(math.Pow(p.r0/ri, 1/p.n)*math.Tan(math.Pi/4+p.phip/2)) - math.Pi/4)
	ui := math.Asin(math.Cos(p.alphac)*math.Sin(ti) - math.Sin(p.alphac)*math.Cos(ti)*math.Cos(di))
	vi := math.Asin(math.Cos(ti) * math.Sin(di) / math.Cos(ui))

	phi := ui

	for i := 0; i < 3; i++ {
		phi = 2 * (math.Atan(math.Pow(p.t0, -1/p.b)*math.Pow(math.Tan(ui/2+math.Pi/4), 1/p.b)*math.Pow((1+p.geogr.Spheroid.E*math.Sin(phi))/(1-p.geogr.Spheroid.E*math.Sin(phi)), p.geogr.Spheroid.E/2)) - math.Pi/4)
	}

	lambda := p.lambda0 - vi/p.b

	return p.geogr.ToWGS84(degree(lambda), degree(phi), h)
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

func sin2(r float64) float64 {
	return math.Pow(math.Sin(r), 2)
}

func degree(r float64) float64 {
	return r * 180 / math.Pi
}

func radian(g float64) float64 {
	return g * math.Pi / 180
}
