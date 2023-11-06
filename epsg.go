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
		return Projected{
			Geographic: EPSG(4171).(Geographic),
			Projection: LambertConformalConic2SP{
				LatitudeOfOrigin:  46.5,
				CentralMeridian:   3,
				StandardParallel1: 49,
				StandardParallel2: 44,
				FalseEasting:      700000,
				FalseNorthing:     6600000,
			},
		}
	case 3035:
		return Projected{
			Geographic: EPSG(4258).(Geographic),
			Projection: LambertAzimuthalEqualArea{
				LatitudeOfCenter:  52,
				LongitudeOfCenter: 10,
				FalseEasting:      4321000,
				FalseNorthing:     3210000,
			},
		}
	case 3161:
		return Projected{
			Geographic: EPSG(4269).(Geographic),
			Projection: LambertConformalConic2SP{
				LatitudeOfOrigin:  0,
				CentralMeridian:   -85,
				StandardParallel1: 44.5,
				StandardParallel2: 53.5,
				FalseEasting:      930000,
				FalseNorthing:     6430000,
			},
		}
	case 3416:
		return Projected{
			Geographic: EPSG(4258).(Geographic),
			Projection: LambertConformalConic2SP{
				CentralMeridian:   13.33333333333333,
				LatitudeOfOrigin:  47.5,
				StandardParallel1: 49,
				StandardParallel2: 46,
				FalseEasting:      400000,
				FalseNorthing:     400000,
			},
		}
	case 3875:
		return Projected{
			Geographic: EPSG(4326).(Geographic),
			Projection: WebMercator{},
		}
	case 4156:
		return Geographic{
			Spheroid: Bessel1841,
			Geocentric: Helmert{
				Tx: 589,
				Ty: 76,
				Tz: 480,
			},
		}
	case 4171:
		return Geographic{
			Spheroid: GRS80,
		}
	case 4258:
		return Geographic{
			Spheroid: GRS80,
		}
	case 4269:
		return Geographic{
			Spheroid: GRS80,
		}
	case 4277:
		return Geographic{
			Spheroid: Airy1830,
			Geocentric: Helmert{
				Tx: 446.448,
				Ty: -125.157,
				Tz: 542.06,
				Rx: 0.15,
				Ry: 0.247,
				Rz: 0.842,
				Ds: -20.489,
			},
		}
	case 4312:
		return Geographic{
			Spheroid: Bessel1841,
			Geocentric: Helmert{
				Tx: 577.326,
				Ty: 90.129,
				Tz: 463.919,
				Rx: 5.137,
				Ry: 1.474,
				Rz: 5.297,
				Ds: 2.4232,
			},
		}
	case 4314:
		return Geographic{
			Spheroid: Bessel1841,
			Geocentric: Helmert{
				Tx: 598.1,
				Ty: 73.7,
				Tz: 418.2,
				Rx: 0.202,
				Ry: 0.045,
				Rz: -2.455,
				Ds: 6.7,
			},
		}
	case 4326:
		return Geographic{
			Spheroid: WGS84,
		}
	case 4978:
		return root{}
	case 5514:
		return Projected{
			Geographic: EPSG(4156).(Geographic),
			Projection: Krovak{
				LatitudeOfCenter:       49.5,
				LongitudeOfCenter:      24.8333333333333,
				Azimuth:                30.2881397527778,
				PseudoStandardParallel: 78.5,
				ScaleFactor:            0.9999,
			},
		}
	case 6355:
		return Projected{
			Geographic: EPSG(4269).(Geographic),
			Projection: TransverseMercator{
				LatitudeOfOrigin: 30.5,
				CentralMeridian:  -85.8333333333333,
				ScaleFactor:      0.99996,
				FalseEasting:     200000,
				FalseNorthing:    0,
			},
		}
	case 6356:
		return Projected{
			Geographic: EPSG(4269).(Geographic),
			Projection: TransverseMercator{
				LatitudeOfOrigin: 30,
				CentralMeridian:  -87.5,
				ScaleFactor:      0.999933333,
				FalseEasting:     600000,
				FalseNorthing:    0,
			},
		}
	case 6414:
		return Projected{
			Geographic: EPSG(4269).(Geographic),
			Projection: AlbersConicEqualArea{
				LatitudeOfCenter:  0,
				LongitudeOfCenter: -120,
				StandardParallel1: 34,
				StandardParallel2: 40.5,
				FalseEasting:      0,
				FalseNorthing:     -4000000,
			},
		}
	case 27700:
		return Projected{
			Geographic: EPSG(4277).(Geographic),
			Projection: TransverseMercator{
				CentralMeridian:  -2,
				LatitudeOfOrigin: 49,
				ScaleFactor:      0.9996012717,
				FalseEasting:     400000,
				FalseNorthing:    -100000,
			},
		}
	case 31257:
		return Projected{
			Geographic: EPSG(4312).(Geographic),
			Projection: TransverseMercator{
				CentralMeridian:  10.33333333333333,
				LatitudeOfOrigin: 0,
				ScaleFactor:      1,
				FalseEasting:     150000,
				FalseNorthing:    -5000000,
			},
		}
	case 31258:
		return Projected{
			Geographic: EPSG(4312).(Geographic),
			Projection: TransverseMercator{
				CentralMeridian:  13.33333333333333,
				LatitudeOfOrigin: 0,
				ScaleFactor:      1,
				FalseEasting:     450000,
				FalseNorthing:    -5000000,
			},
		}
	case 31259:
		return Projected{
			Geographic: EPSG(4312).(Geographic),
			Projection: TransverseMercator{
				CentralMeridian:  16.33333333333333,
				LatitudeOfOrigin: 0,
				ScaleFactor:      1,
				FalseEasting:     750000,
				FalseNorthing:    -5000000,
			},
		}
	case 31284:
		return Projected{
			Geographic: EPSG(4312).(Geographic),
			Projection: TransverseMercator{
				CentralMeridian:  10.33333333333333,
				LatitudeOfOrigin: 0,
				ScaleFactor:      1,
				FalseEasting:     150000,
				FalseNorthing:    0,
			},
		}
	case 31285:
		return Projected{
			Geographic: EPSG(4312).(Geographic),
			Projection: TransverseMercator{
				CentralMeridian:  13.33333333333333,
				LatitudeOfOrigin: 0,
				ScaleFactor:      1,
				FalseEasting:     450000,
				FalseNorthing:    0,
			},
		}
	case 31286:
		return Projected{
			Geographic: EPSG(4312).(Geographic),
			Projection: TransverseMercator{
				CentralMeridian:  16.33333333333333,
				LatitudeOfOrigin: 0,
				ScaleFactor:      1,
				FalseEasting:     750000,
				FalseNorthing:    0,
			},
		}
	case 31287:
		return Projected{
			Geographic: EPSG(4312).(Geographic),
			Projection: LambertConformalConic2SP{
				CentralMeridian:   13.33333333333333,
				LatitudeOfOrigin:  47.5,
				StandardParallel1: 49,
				StandardParallel2: 46,
				FalseEasting:      400000,
				FalseNorthing:     400000,
			},
		}
	case 900913:
		return EPSG(3857)
	}

	if code > 3941 && code < 3951 {
		lat := float64(code - 3900)

		return Projected{
			Geographic: EPSG(4171).(Geographic),
			Projection: LambertConformalConic2SP{
				LatitudeOfOrigin:  lat,
				CentralMeridian:   3,
				StandardParallel1: lat - 0.75,
				StandardParallel2: lat + 0.75,
				FalseEasting:      1700000,
				FalseNorthing:     2200000 + (lat-43)*1000000,
			},
		}
	}

	if code > 25827 && code < 25839 {
		zone := float64(code - 25800)

		return Projected{
			Geographic: EPSG(4258).(Geographic),
			Projection: TransverseMercator{
				CentralMeridian:  zone*6 - 183,
				LatitudeOfOrigin: 0,
				ScaleFactor:      0.9996,
				FalseEasting:     500000,
				FalseNorthing:    0,
			},
		}
	}

	if code > 31465 && code < 31470 {
		zone := float64(code - 31464)

		return Projected{
			Geographic: EPSG(4314).(Geographic),
			Projection: TransverseMercator{
				LatitudeOfOrigin: 0,
				CentralMeridian:  zone * 3,
				ScaleFactor:      1,
				FalseEasting:     zone*1000000 + 500000,
				FalseNorthing:    0,
			},
		}
	}

	if code > 32600 && code < 32661 {
		zone := float64(code - 32600)

		return Projected{
			Geographic: EPSG(4326).(Geographic),
			Projection: TransverseMercator{
				LatitudeOfOrigin: 0,
				CentralMeridian:  zone*6 - 183,
				ScaleFactor:      0.9996,
				FalseEasting:     500000,
				FalseNorthing:    0,
			},
		}
	}

	if code > 32700 && code < 32761 {
		zone := code - 32700

		return Projected{
			Geographic: EPSG(4326).(Geographic),
			Projection: TransverseMercator{
				CentralMeridian:  float64(zone)*6 - 183,
				LatitudeOfOrigin: 0,
				ScaleFactor:      0.9996,
				FalseEasting:     500000,
				FalseNorthing:    10000000,
			},
		}
	}

	return ErrCodeNotFound{
		Code: code,
	}
}
