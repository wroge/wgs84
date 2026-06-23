package wgs84

var OSNI1952 = Datum{
	Spheroid: Airy1830,
	Transformations: []Transformation{
		{
			Accuracy:    1,
			Helmert:     Helmert{482.5, -130.6, 564.6, -1.042, -0.214, -0.631, 8.15},
			BoundingBox: BoundingBox{-8.18, 53.96, -5.34, 55.36},
		},
	},
}
