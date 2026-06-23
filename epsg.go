package wgs84

var EPSG = map[int]CoordinateReferenceSystem{
	4978: {
		CoordinateSystem: Geocentric{},
		Datum:            WGS84,
	},
	4326: {
		CoordinateSystem: Geographic{},
		Datum:            WGS84,
	},
	3857: {
		CoordinateSystem: WebMercator{},
		Datum: Datum{
			Spheroid: WGS84.Spheroid,
			Transformations: []Transformation{
				{
					BoundingBox: BoundingBox{-180.0, -85.06, 180.0, 85.06},
				},
			},
		},
	},
	900913: {
		CoordinateSystem: WebMercator{},
		Datum: Datum{
			Spheroid: WGS84.Spheroid,
			Transformations: []Transformation{
				{
					BoundingBox: BoundingBox{-180.0, -85.06, 180.0, 85.06},
				},
			},
		},
	},
	4267: {
		CoordinateSystem: Geographic{},
		Datum:            NAD27,
	},
	32024: {
		CoordinateSystem: LambertConformalConic2SP{
			Lonf:   -98,
			Latf:   35,
			Sp1:    35.5666666666667,
			Sp2:    36.7666666666667,
			Eastf:  2000000,
			Northf: 0,
		},
		Datum: NAD27,
	},
	4314: {
		CoordinateSystem: Geographic{},
		Datum:            DHDN90,
	},
	4313: {
		CoordinateSystem: Geographic{},
		Datum:            BD72,
	},
	31370: {
		CoordinateSystem: LambertConformalConic2SP{
			Lonf:   4.36748666666667,
			Latf:   90,
			Sp1:    51.1666672333333,
			Sp2:    49.8333339,
			Eastf:  150000.013,
			Northf: 5400088.438,
		},
		Datum: BD72,
	},
	4312: {
		CoordinateSystem: Geographic{},
		Datum:            MGI,
	},
	31255: {
		CoordinateSystem: TransverseMercator{
			Latf:   0,
			Lonf:   13.33333333333333,
			Scale:  1,
			Eastf:  0,
			Northf: -5000000,
		},
		Datum: MGI,
	},
	31257: {
		CoordinateSystem: TransverseMercator{
			Latf:   0,
			Lonf:   10.33333333333333,
			Scale:  1,
			Eastf:  150000,
			Northf: -5000000,
		},
		Datum: MGI,
	},
	31258: {
		CoordinateSystem: TransverseMercator{
			Latf:   0,
			Lonf:   13.33333333333333,
			Scale:  1,
			Eastf:  450000,
			Northf: -5000000,
		},
		Datum: MGI,
	},
	31259: {
		CoordinateSystem: TransverseMercator{
			Latf:   0,
			Lonf:   16.33333333333333,
			Scale:  1,
			Eastf:  750000,
			Northf: -5000000,
		},
		Datum: MGI,
	},
	31284: {
		CoordinateSystem: TransverseMercator{
			Latf:   0,
			Lonf:   10.33333333333333,
			Scale:  1,
			Eastf:  150000,
			Northf: 0,
		},
		Datum: MGI,
	},
	31285: {
		CoordinateSystem: TransverseMercator{
			Latf:   0,
			Lonf:   13.33333333333333,
			Scale:  1,
			Eastf:  450000,
			Northf: 0,
		},
		Datum: MGI,
	},
	31286: {
		CoordinateSystem: TransverseMercator{
			Latf:   0,
			Lonf:   16.33333333333333,
			Scale:  1,
			Eastf:  750000,
			Northf: 0,
		},
		Datum: MGI,
	},
	4277: {
		CoordinateSystem: Geographic{},
		Datum:            OSGB36,
	},
	27700: {
		CoordinateSystem: TransverseMercator{
			Latf:   49,
			Lonf:   -2,
			Scale:  0.9996012717,
			Eastf:  400000,
			Northf: -100000,
		},
		Datum: OSGB36,
	},
	4156: {
		CoordinateSystem: Geographic{},
		Datum:            SJTSK,
	},
	5514: {
		CoordinateSystem: Krovak{
			Lonf:    24.8333333333333,
			Latf:    49.5,
			Azimuth: 30.2881397527778,
			Sp:      78.5,
			Scale:   0.9999,
			Eastf:   0,
			Northf:  0,
		},
		Datum: SJTSK,
	},
	4173: {
		CoordinateSystem: Geographic{},
		Datum:            IRENET95,
	},
	2157: {
		CoordinateSystem: TransverseMercator{
			Latf:   53.5,
			Lonf:   -8,
			Scale:  0.99982,
			Eastf:  600000,
			Northf: 750000,
		},
		Datum: IRENET95,
	},
	2158: {
		CoordinateSystem: TransverseMercator{
			Latf:   0,
			Lonf:   -9,
			Scale:  0.9996,
			Eastf:  500000,
			Northf: 0,
		},
		Datum: IRENET95,
	},
	4171: {
		CoordinateSystem: Geographic{},
		Datum:            RGF93,
	},
	2154: {
		CoordinateSystem: LambertConformalConic2SP{
			Lonf:   3,
			Latf:   46.5,
			Sp1:    49,
			Sp2:    44,
			Eastf:  700000,
			Northf: 6600000,
		},
		Datum: RGF93,
	},
	4258: {
		CoordinateSystem: Geographic{},
		Datum:            ETRS89,
	},
	3035: {
		CoordinateSystem: LambertAzimuthalEqualArea{
			Lonf:   10,
			Latf:   52,
			Eastf:  4321000,
			Northf: 3210000,
		},
		Datum: ETRS89,
	},
	3416: {
		CoordinateSystem: LambertConformalConic2SP{
			Lonf:   13.33333333333333,
			Latf:   47.5,
			Sp1:    49,
			Sp2:    46,
			Eastf:  400000,
			Northf: 400000,
		},
		Datum: ETRS89,
	},
	102109: {
		CoordinateSystem: TransverseMercator{
			Latf:   0,
			Lonf:   15,
			Scale:  0.9999,
			Eastf:  500000,
			Northf: -5000000,
		},
		Datum: ETRS89,
	},
	102157: {
		CoordinateSystem: TransverseMercator{
			Latf:   0,
			Lonf:   21,
			Scale:  0.9999,
			Eastf:  7500000,
			Northf: 0,
		},
		Datum: ETRS89,
	},
	102173: {
		CoordinateSystem: TransverseMercator{
			Latf:   0,
			Lonf:   19,
			Scale:  0.9993,
			Eastf:  500000,
			Northf: -5300000,
		},
		Datum: ETRS89,
	},
	3126: {
		CoordinateSystem: TransverseMercator{
			Latf:   0,
			Lonf:   19,
			Scale:  1,
			Eastf:  500000,
			Northf: 0,
		},
		Datum: ETRS89,
	},
	3127: {
		CoordinateSystem: TransverseMercator{
			Latf:   0,
			Lonf:   20,
			Scale:  1,
			Eastf:  500000,
			Northf: 0,
		},
		Datum: ETRS89,
	},
	3128: {
		CoordinateSystem: TransverseMercator{
			Latf:   0,
			Lonf:   21,
			Scale:  1,
			Eastf:  500000,
			Northf: 0,
		},
		Datum: ETRS89,
	},
	3129: {
		CoordinateSystem: TransverseMercator{
			Latf:   0,
			Lonf:   22,
			Scale:  1,
			Eastf:  500000,
			Northf: 0,
		},
		Datum: ETRS89,
	},
	3130: {
		CoordinateSystem: TransverseMercator{
			Latf:   0,
			Lonf:   23,
			Scale:  1,
			Eastf:  500000,
			Northf: 0,
		},
		Datum: ETRS89,
	},
	3131: {
		CoordinateSystem: TransverseMercator{
			Latf:   0,
			Lonf:   24,
			Scale:  1,
			Eastf:  500000,
			Northf: 0,
		},
		Datum: ETRS89,
	},
	3132: {
		CoordinateSystem: TransverseMercator{
			Latf:   0,
			Lonf:   25,
			Scale:  1,
			Eastf:  500000,
			Northf: 0,
		},
		Datum: ETRS89,
	},
	3133: {
		CoordinateSystem: TransverseMercator{
			Latf:   0,
			Lonf:   26,
			Scale:  1,
			Eastf:  500000,
			Northf: 0,
		},
		Datum: ETRS89,
	},
	3134: {
		CoordinateSystem: TransverseMercator{
			Latf:   0,
			Lonf:   27,
			Scale:  1,
			Eastf:  500000,
			Northf: 0,
		},
		Datum: ETRS89,
	},
	3135: {
		CoordinateSystem: TransverseMercator{
			Latf:   0,
			Lonf:   28,
			Scale:  1,
			Eastf:  500000,
			Northf: 0,
		},
		Datum: ETRS89,
	},
	3136: {
		CoordinateSystem: TransverseMercator{
			Latf:   0,
			Lonf:   29,
			Scale:  1,
			Eastf:  500000,
			Northf: 0,
		},
		Datum: ETRS89,
	},
	3137: {
		CoordinateSystem: TransverseMercator{
			Latf:   0,
			Lonf:   30,
			Scale:  1,
			Eastf:  500000,
			Northf: 0,
		},
		Datum: ETRS89,
	},
	3138: {
		CoordinateSystem: TransverseMercator{
			Latf:   0,
			Lonf:   31,
			Scale:  1,
			Eastf:  500000,
			Northf: 0,
		},
		Datum: ETRS89,
	},
	4230: {
		CoordinateSystem: Geographic{},
		Datum:            ED50,
	},
	23090: {
		CoordinateSystem: TransverseMercator{
			Lonf:   0,
			Latf:   0,
			Scale:  0.9996,
			Eastf:  500000,
			Northf: 0,
		},
		Datum: ED50,
	},
	4269: {
		CoordinateSystem: Geographic{},
		Datum:            NAD83,
	},
	3161: {
		CoordinateSystem: LambertConformalConic2SP{
			Lonf:   -85,
			Latf:   0,
			Sp1:    44.5,
			Sp2:    53.5,
			Eastf:  930000,
			Northf: 6430000,
		},
		Datum: NAD83,
	},
	26917: {
		CoordinateSystem: TransverseMercator{
			Latf:   0,
			Lonf:   -81,
			Scale:  0.9996,
			Eastf:  500000,
			Northf: 0,
		},
		Datum: NAD83,
	},
	4188: {
		CoordinateSystem: Geographic{},
		Datum:            OSNI1952,
	},
	29901: {
		CoordinateSystem: TransverseMercator{
			Latf:   53.5,
			Lonf:   -8,
			Scale:  1,
			Eastf:  200000,
			Northf: 250000,
		},
		Datum: OSNI1952,
	},
	4299: {
		CoordinateSystem: Geographic{},
		Datum:            TM65,
	},
	29902: {
		CoordinateSystem: TransverseMercator{
			Latf:   53.5,
			Lonf:   -8,
			Scale:  1.000035,
			Eastf:  200000,
			Northf: 250000,
		},
		Datum: TM65,
	},
	4300: {
		CoordinateSystem: Geographic{},
		Datum:            TM75,
	},
	29903: {
		CoordinateSystem: TransverseMercator{
			Latf:   53.5,
			Lonf:   -8,
			Scale:  1.000035,
			Eastf:  200000,
			Northf: 250000,
		},
		Datum: TM75,
	},
	4490: {
		CoordinateSystem: Geographic{},
		Datum:            ChinaGeodeticCoordinateSystem2000,
	},
	4549: {
		CoordinateSystem: TransverseMercator{
			Lonf:   120,
			Latf:   0,
			Scale:  1,
			Eastf:  500000,
			Northf: 0,
		},
		Datum: ChinaGeodeticCoordinateSystem2000,
	},
	4807: {
		CoordinateSystem: Geographic{},
		Datum:            NTF,
	},
	27572: {
		CoordinateSystem: LambertConformalConic1SP{
			Lonf:   2.33722917,
			Latf:   46.8,
			Scale:  0.99987742,
			Eastf:  600000,
			Northf: 2200000,
		},
	},
	6414: {
		CoordinateSystem: AlbersConicEqualArea{
			Lonf:   -120,
			Latf:   0,
			Sp1:    34,
			Sp2:    40.5,
			Eastf:  0,
			Northf: -4000000,
		},
		Datum: Datum{
			Spheroid: GRS80,
			Transformations: []Transformation{
				{
					Accuracy:    2,
					BoundingBox: BoundingBox{167.65, 14.92, -63.88, 74.71},
				},
			},
		},
	},
}

func init() {
	for code := 25828; code < 25839; code++ {
		zone := float64(code - 25800)

		EPSG[code] = CoordinateReferenceSystem{
			CoordinateSystem: TransverseMercator{
				Lonf:   zone*6 - 183,
				Latf:   0,
				Scale:  0.9996,
				Eastf:  500000,
				Northf: 0,
			},
			Datum: ETRS89,
		}
	}

	for code := 31466; code < 31470; code++ {
		zone := float64(code - 31464)

		EPSG[code] = CoordinateReferenceSystem{
			CoordinateSystem: TransverseMercator{
				Lonf:   zone * 3,
				Latf:   0,
				Scale:  1,
				Eastf:  zone*1000000 + 500000,
				Northf: 0,
			},
			Datum: DHDN90,
		}
	}

	for code := 32601; code < 32661; code++ {
		zone := float64(code - 32600)

		EPSG[code] = CoordinateReferenceSystem{
			CoordinateSystem: TransverseMercator{
				Lonf:   zone*6 - 183,
				Latf:   0,
				Scale:  0.9996,
				Eastf:  500000,
				Northf: 0,
			},
			Datum: WGS84,
		}
	}

	for code := 32701; code < 32761; code++ {
		zone := float64(code - 32700)

		EPSG[code] = CoordinateReferenceSystem{
			CoordinateSystem: TransverseMercator{
				Lonf:   zone*6 - 183,
				Latf:   0,
				Scale:  0.9996,
				Eastf:  500000,
				Northf: 10000000,
			},
			Datum: WGS84,
		}
	}

	for code := 3942; code < 3951; code++ {
		lat := float64(code - 3900)

		EPSG[code] = CoordinateReferenceSystem{
			CoordinateSystem: LambertConformalConic2SP{
				Lonf:   3,
				Latf:   lat,
				Sp1:    lat - 0.75,
				Sp2:    lat + 0.75,
				Eastf:  1700000,
				Northf: 2200000 + (lat-43)*1000000,
			},
			Datum: RGF93,
		}
	}
}
