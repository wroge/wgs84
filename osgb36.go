package wgs84

var OSGB36 = Datum{
	Spheroid: Airy1830,
	Transformations: []Transformation{
		{
			Accuracy:    1,
			BoundingBox: BoundingBox{-9.01, 49.75, 2.01, 61.01},
			Grid:        "uk_os_OSTN15_NTv2_OSGBtoETRS.tif",
			GridOptional: true,
		},
		{
			Accuracy:    2,
			BoundingBox: BoundingBox{-8.82, 49.79, 1.92, 60.94},
			Helmert:     Helmert{446.448, -125.157, 542.06, 0.15, 0.247, 0.842, -20.489},
		},
		{
			Accuracy:    3,
			BoundingBox: BoundingBox{-2.2, 50.53, -1.68, 50.8},
			Helmert:     Helmert{370.936, -108.938, 435.682, 0, 0, 0, 0},
		},
	},
}
