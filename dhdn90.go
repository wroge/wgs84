package wgs84

var DHDN90 = Datum{
	Spheroid: Bessel1841,
	Transformations: []Transformation{
		{
			Accuracy:    1,
			BoundingBox: BoundingBox{5.86, 47.27, 15.04, 55.09},
			Grid:        "de_adv_BETA2007.tif",
			GridOptional: true,
		},
		{
			Accuracy:    2,
			Helmert:     Helmert{612.4, 77, 440.2, -0.054, 0.057, -2.797, 2.55},
			BoundingBox: BoundingBox{9.92, 50.2, 15.04, 54.74},
		},
		{
			Accuracy:    3,
			Helmert:     Helmert{598.1, 73.7, 418.2, 0.202, 0.045, -2.455, 6.7},
			BoundingBox: BoundingBox{5.86, 47.27, 13.84, 55.09},
		},
		{
			Accuracy:    5,
			Helmert:     Helmert{582, 105, 414, -1.04, -0.35, 3.08, 8.3},
			BoundingBox: BoundingBox{5.86, 47.27, 13.84, 55.09},
		},
	},
}
