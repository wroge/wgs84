package wgs84

import (
	"fmt"
	"sync"
)

var crsStore sync.Map

func EPSG(code int) CRS {
	if crs, ok := crsStore.Load(code); ok {
		return crs.(CRS)
	}

	var crs CRS

	switch code {
	case 2154:
		crs = LambertConformalConic2SP(EPSG(4171), 3, 46.5, 49, 44, 700000, 6600000)
	case 2157:
		crs = TransverseMercator(EPSG(4173), -8, 53.5, 0.99982, 600000, 750000)
	case 2158:
		crs = TransverseMercator(EPSG(4173), -9, 0, 0.9996, 500000, 0)
	case 3035:
		crs = LambertAzimuthalEqualArea(EPSG(4258), 10, 52, 4321000, 3210000)
	case 3161:
		crs = LambertConformalConic2SP(EPSG(4269), -85, 0, 44.5, 53.5, 930000, 6430000)
	case 3416:
		crs = LambertConformalConic2SP(EPSG(4258), 13.33333333333333, 47.5, 49, 46, 400000, 400000)
	case 3857:
		crs = WebMercator(EPSG(4326))
	case 4156:
		crs = Geographic(Helmert(589, 76, 480, 0, 0, 0, 0), NewSpheroid(6377397.155, 299.1528128))
	case 4171:
		crs = Geographic(EPSG(4978), NewSpheroid(6378137, 298.257222101))
	case 4173:
		crs = Geographic(EPSG(4978), NewSpheroid(6378137, 298.257222101))
	case 4188:
		crs = Geographic(Helmert(482.5, -130.6, 564.6, -1.042, -0.214, -0.631, 8.15), NewSpheroid(6377563.396, 299.3249646))
	case 4230:
		crs = Geographic(Helmert(-87, -98, -121, 0, 0, 0, 0), NewSpheroid(6378388, 297))
	case 4258:
		crs = Geographic(EPSG(4978), NewSpheroid(6378137, 298.257222101))
	case 4269:
		crs = Geographic(EPSG(4978), NewSpheroid(6378137, 298.257222101))
	case 4277:
		crs = loadNTv2("OSTN15_NTv2_OSGBtoETRS.gsb", NewSpheroid(6377563.396, 299.3249646), EPSG(4326))
	case 4299:
		crs = Geographic(Helmert(482.5, -130.6, 564.6, -1.042, -0.214, -0.631, 8.15), NewSpheroid(6377340.189, 299.3249646))
	case 4300:
		crs = EPSG(4299)
	case 4312:
		crs = Geographic(Helmert(577.326, 90.129, 463.919, 5.137, 1.474, 5.297, 2.4232), NewSpheroid(6377397.155, 299.1528128))
	case 4314:
		crs = loadNTv2("BeTA2007.gsb", NewSpheroid(6377397.155, 299.1528128), EPSG(4326))
	case 4326:
		crs = Geographic(EPSG(4978), NewSpheroid(6378137, 298.257223563))
	case 4490:
		crs = Geographic(EPSG(4978), NewSpheroid(6378137, 298.257222101))
	case 4549:
		crs = TransverseMercator(EPSG(4490), 120, 0, 1, 500000, 0)
	case 4978:
		crs = base{}
	case 5514:
		crs = Krovak(EPSG(4156), 24.8333333333333, 49.5, 30.2881397527778, 78.5, 0.9999, 0, 0)
	case 6318:
		crs = Geographic(EPSG(4978), NewSpheroid(6378137, 298.257222101))
	case 6355:
		crs = TransverseMercator(EPSG(6318), -85.8333333333333, 30.5, 0.99996, 200000, 0)
	case 6356:
		crs = TransverseMercator(EPSG(6318), -87.5, 30, 0.999933333, 600000, 0)
	case 6414:
		crs = AlbersConicEqualArea(EPSG(6318), -120, 0, 34, 40.5, 0, -4000000)
	case 23090:
		crs = TransverseMercator(EPSG(4230), 0, 0, 0.9996, 500000, 0)
	case 27700:
		crs = TransverseMercator(EPSG(4277), -2, 49, 0.9996012717, 400000, -100000)
	case 29901:
		crs = TransverseMercator(EPSG(4188), -8, 53.5, 1, 200000, 250000)
	case 29902:
		crs = TransverseMercator(EPSG(4299), -8, 53.5, 1.000035, 200000, 250000)
	case 29903:
		crs = TransverseMercator(EPSG(4300), -8, 53.5, 1.000035, 200000, 250000)
	case 31257:
		crs = TransverseMercator(EPSG(4312), 10.33333333333333, 0, 1, 150000, -5000000)
	case 31258:
		crs = TransverseMercator(EPSG(4312), 13.33333333333333, 0, 1, 450000, -5000000)
	case 31259:
		crs = TransverseMercator(EPSG(4312), 16.33333333333333, 0, 1, 750000, -5000000)
	case 31284:
		crs = TransverseMercator(EPSG(4312), 10.33333333333333, 0, 1, 150000, 0)
	case 31285:
		crs = TransverseMercator(EPSG(4312), 13.33333333333333, 0, 1, 450000, 0)
	case 31286:
		crs = TransverseMercator(EPSG(4312), 16.33333333333333, 0, 1, 750000, 0)
	case 31287:
		crs = LambertConformalConic2SP(EPSG(4312), 13.33333333333333, 47.5, 49, 46, 400000, 400000)
	case 900913:
		crs = EPSG(3857)
	default:
		if code > 3941 && code < 3951 {
			lat := float64(code - 3900)

			crs = LambertConformalConic2SP(EPSG(4171), 3, lat, lat-0.75, lat+0.75, 1700000, 2200000+(lat-43)*1000000)
		} else if code > 25827 && code < 25839 {
			zone := float64(code - 25800)

			crs = TransverseMercator(EPSG(4258), zone*6-183, 0, 0.9996, 500000, 0)
		} else if code > 31465 && code < 31470 {
			zone := float64(code - 31464)

			crs = TransverseMercator(EPSG(4314), zone*3, 0, 1, zone*1000000+500000, 0)
		} else if code > 32600 && code < 32661 {
			zone := float64(code - 32600)

			crs = TransverseMercator(EPSG(4326), zone*6-183, 0, 0.9996, 500000, 0)
		} else if code > 32700 && code < 32761 {
			zone := code - 32700

			crs = TransverseMercator(EPSG(4326), float64(zone)*6-183, 0, 0.9996, 500000, 10000000)
		}
	}

	if crs == nil {
		return errorCRS{err: fmt.Errorf("epsg code '%d' not found", code)}
	}

	crsStore.Store(code, crs)

	return crs
}
