//nolint:gomnd,nonamedreturns,errname,funlen,gocyclo,cyclop,ireturn
package wgs84

import (
	"fmt"
)

type ErrCodeNotFound struct {
	Code int
}

func (e ErrCodeNotFound) Error() string {
	return fmt.Sprintf("epsg code '%d' not found", e.Code)
}

func (e ErrCodeNotFound) ToWGS84(x, y, z float64) (x0, y0, z0 float64) {
	return x, y, z
}

func (e ErrCodeNotFound) FromWGS84(x, y, z float64) (x0, y0, z0 float64) {
	return x, y, z
}

func EPSG(code int) CoordinateReferenceSystem {
	switch code {
	case 2154:
		return RGF93.LambertConformalConic2SP(3, 46.5, 49, 44, 700000, 6600000)
	case 2157:
		return IRENET95.TransverseMercator(-8, 53.5, 0.99982, 600000, 750000)
	case 2158:
		return IRENET95.TransverseMercator(-9, 0, 0.9996, 500000, 0)
	case 3035:
		return ETRS89.LambertAzimuthalEqualArea(10, 52, 4321000, 3210000)
	case 3161:
		return NAD83.LambertConformalConic2SP(-85, 0, 44.5, 53.5, 930000, 6430000)
	case 3416:
		return ETRS89.LambertConformalConic2SP(13.33333333333333, 47.5, 49, 46, 400000, 400000)
	case 3875:
		return Datum.WebMercator()
	case 4156:
		return SJTSK
	case 4171:
		return RGF93
	case 4188:
		return OSNI1952
	case 4230:
		return ED50
	case 4258:
		return ETRS89
	case 4269:
		return NAD83
	case 4277:
		return OSGB36
	case 4299:
		return TM65
	case 4300:
		return TM75
	case 4312:
		return MGI
	case 4314:
		return DHDN2001
	case 4326:
		return Datum
	case 4978:
		return root{}
	case 5514:
		return SJTSK.Krovak(24.8333333333333, 49.5, 30.2881397527778, 78.5, 0.9999, 0, 0)
	case 6355:
		return NAD83.TransverseMercator(-85.8333333333333, 30.5, 0.99996, 200000, 0)
	case 6356:
		return NAD83.TransverseMercator(-87.5, 30, 0.999933333, 600000, 0)
	case 6414:
		return NAD83.AlbersConicEqualArea(-120, 0, 34, 40.5, 0, -4000000)
	case 23090:
		return ED50.TransverseMercator(0, 0, 0.9996, 500000, 0)
	case 27700:
		return OSGB36.TransverseMercator(-2, 49, 0.9996012717, 400000, -100000)
	case 29901:
		return OSNI1952.TransverseMercator(-8, 53.5, 1, 200000, 250000)
	case 29902:
		return TM65.TransverseMercator(-8, 53.5, 1.000035, 200000, 250000)
	case 29903:
		return TM75.TransverseMercator(-8, 53.5, 1.000035, 200000, 250000)
	case 31257:
		return MGI.TransverseMercator(10.33333333333333, 0, 1, 150000, -5000000)
	case 31258:
		return MGI.TransverseMercator(13.33333333333333, 0, 1, 450000, -5000000)
	case 31259:
		return MGI.TransverseMercator(16.33333333333333, 0, 1, 750000, -5000000)
	case 31284:
		return MGI.TransverseMercator(10.33333333333333, 0, 1, 150000, 0)
	case 31285:
		return MGI.TransverseMercator(13.33333333333333, 0, 1, 450000, 0)
	case 31286:
		return MGI.TransverseMercator(16.33333333333333, 0, 1, 750000, 0)
	case 31287:
		return MGI.LambertConformalConic2SP(13.33333333333333, 47.5, 49, 46, 400000, 400000)
	case 900913:
		return Datum.WebMercator()
	}

	if code > 3941 && code < 3951 {
		lat := float64(code - 3900)

		return RGF93.LambertConformalConic2SP(3, lat, lat-0.75, lat+0.75, 1700000, 2200000+(lat-43)*1000000)
	}

	if code > 25827 && code < 25839 {
		zone := float64(code - 25800)

		return ETRS89.TransverseMercator(zone*6-183, 0, 0.9996, 500000, 0)
	}

	if code > 31465 && code < 31470 {
		zone := float64(code - 31464)

		return DHDN2001.TransverseMercator(zone*3, 0, 1, zone*1000000+500000, 0)
	}

	if code > 32600 && code < 32661 {
		zone := float64(code - 32600)

		return Datum.TransverseMercator(zone*6-183, 0, 0.9996, 500000, 0)
	}

	if code > 32700 && code < 32761 {
		zone := code - 32700

		return Datum.TransverseMercator(float64(zone)*6-183, 0, 0.9996, 500000, 10000000)
	}

	return ErrCodeNotFound{
		Code: code,
	}
}
