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
		return LambertConformalConic2SP{
			LatitudeOfOrigin:  46.5,
			CentralMeridian:   3,
			StandardParallel1: 49,
			StandardParallel2: 44,
			FalseEasting:      700000,
			FalseNorthing:     6600000,
			Geographic:        RGF93,
		}
	case 3035:
		return LambertAzimuthalEqualArea{
			LatitudeOfCenter:  52,
			LongitudeOfCenter: 10,
			FalseEasting:      4321000,
			FalseNorthing:     3210000,
			Geographic:        ETRS89,
		}
	case 3161:
		return LambertConformalConic2SP{
			LatitudeOfOrigin:  0,
			CentralMeridian:   -85,
			StandardParallel1: 44.5,
			StandardParallel2: 53.5,
			FalseEasting:      930000,
			FalseNorthing:     6430000,
			Geographic:        NAD83,
		}
	case 3416:
		return LambertConformalConic2SP{
			CentralMeridian:   13.33333333333333,
			LatitudeOfOrigin:  47.5,
			StandardParallel1: 49,
			StandardParallel2: 46,
			FalseEasting:      400000,
			FalseNorthing:     400000,
			Geographic:        ETRS89,
		}
	case 3875:
		return WebMercator{
			Geographic: WGS84,
		}
	case 4171:
		return RGF93
	case 4258:
		return ETRS89
	case 4269:
		return NAD83
	case 4277:
		return OSGB36
	case 4314:
		return DHDN2001
	case 4326:
		return WGS84
	case 4978:
		return root{}
	case 6355:
		return TransverseMercator{
			LatitudeOfOrigin: 30.5,
			CentralMeridian:  -85.8333333333333,
			ScaleFactor:      0.99996,
			FalseEasting:     200000,
			FalseNorthing:    0,
			Geographic:       NAD83,
		}
	case 6356:
		return TransverseMercator{
			LatitudeOfOrigin: 30,
			CentralMeridian:  -87.5,
			ScaleFactor:      0.999933333,
			FalseEasting:     600000,
			FalseNorthing:    0,
			Geographic:       NAD83,
		}
	case 6414:
		return AlbersConicEqualArea{
			LatitudeOfCenter:  0,
			LongitudeOfCenter: -120,
			StandardParallel1: 34,
			StandardParallel2: 40.5,
			FalseEasting:      0,
			FalseNorthing:     -4000000,
			Geographic:        NAD83,
		}
	case 27700:
		return TransverseMercator{
			CentralMeridian:  -2,
			LatitudeOfOrigin: 49,
			ScaleFactor:      0.9996012717,
			FalseEasting:     400000,
			FalseNorthing:    -100000,
			Geographic:       OSGB36,
		}
	case 31257:
		return TransverseMercator{
			CentralMeridian:  10.33333333333333,
			LatitudeOfOrigin: 0,
			ScaleFactor:      1,
			FalseEasting:     150000,
			FalseNorthing:    -5000000,
			Geographic:       MGI,
		}
	case 31258:
		return TransverseMercator{
			CentralMeridian:  13.33333333333333,
			LatitudeOfOrigin: 0,
			ScaleFactor:      1,
			FalseEasting:     450000,
			FalseNorthing:    -5000000,
			Geographic:       MGI,
		}
	case 31259:
		return TransverseMercator{
			CentralMeridian:  16.33333333333333,
			LatitudeOfOrigin: 0,
			ScaleFactor:      1,
			FalseEasting:     750000,
			FalseNorthing:    -5000000,
			Geographic:       MGI,
		}
	case 31284:
		return TransverseMercator{
			CentralMeridian:  10.33333333333333,
			LatitudeOfOrigin: 0,
			ScaleFactor:      1,
			FalseEasting:     150000,
			FalseNorthing:    0,
			Geographic:       MGI,
		}
	case 31285:
		return TransverseMercator{
			CentralMeridian:  13.33333333333333,
			LatitudeOfOrigin: 0,
			ScaleFactor:      1,
			FalseEasting:     450000,
			FalseNorthing:    0,
			Geographic:       MGI,
		}
	case 31286:
		return TransverseMercator{
			CentralMeridian:  16.33333333333333,
			LatitudeOfOrigin: 0,
			ScaleFactor:      1,
			FalseEasting:     750000,
			FalseNorthing:    0,
			Geographic:       MGI,
		}
	case 31287:
		return LambertConformalConic2SP{
			CentralMeridian:   13.33333333333333,
			LatitudeOfOrigin:  47.5,
			StandardParallel1: 49,
			StandardParallel2: 46,
			FalseEasting:      400000,
			FalseNorthing:     400000,
			Geographic:        MGI,
		}
	case 900913:
		return EPSG(3857)
	}

	if code > 3941 && code < 3951 {
		lat := float64(code - 3900)

		return LambertConformalConic2SP{
			LatitudeOfOrigin:  lat,
			CentralMeridian:   3,
			StandardParallel1: lat - 0.75,
			StandardParallel2: lat + 0.75,
			FalseEasting:      1700000,
			FalseNorthing:     2200000 + (lat-43)*1000000,
			Geographic:        RGF93,
		}
	}

	if code > 25827 && code < 25839 {
		zone := float64(code - 25800)

		return TransverseMercator{
			CentralMeridian:  zone*6 - 183,
			LatitudeOfOrigin: 0,
			ScaleFactor:      0.9996,
			FalseEasting:     500000,
			FalseNorthing:    0,
			Geographic:       ETRS89,
		}
	}

	if code > 31465 && code < 31470 {
		zone := float64(code - 31464)

		return TransverseMercator{
			LatitudeOfOrigin: 0,
			CentralMeridian:  zone * 3,
			ScaleFactor:      1,
			FalseEasting:     zone*1000000 + 500000,
			FalseNorthing:    0,
			Geographic:       DHDN2001,
		}
	}

	if code > 32600 && code < 32661 {
		zone := float64(code - 32600)

		return TransverseMercator{
			LatitudeOfOrigin: 0,
			CentralMeridian:  zone*6 - 183,
			ScaleFactor:      0.9996,
			FalseEasting:     500000,
			FalseNorthing:    0,
		}
	}

	if code > 32700 && code < 32761 {
		zone := code - 32700
		return TransverseMercator{
			CentralMeridian:  float64(zone)*6 - 183,
			LatitudeOfOrigin: 0,
			ScaleFactor:      0.9996,
			FalseEasting:     500000,
			FalseNorthing:    10000000,
		}
	}

	return ErrCodeNotFound{
		Code: code,
	}
}
