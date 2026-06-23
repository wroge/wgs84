package wgs84

var SJTSK = Datum{
	Spheroid: Bessel1841,
	Transformations: []Transformation{
		{
			Accuracy:    1,
			Helmert:     Helmert{572.213, 85.334, 461.94, 4.9732, 1.529, 5.2484, 3.5378},
			BoundingBox: BoundingBox{12.09, 48.58, 18.86, 51.06},
		},
		{
			Accuracy:    1,
			Helmert:     Helmert{570.8, 85.7, 462.8, 4.998, 1.587, 5.261, 3.56},
			BoundingBox: BoundingBox{12.09, 48.58, 18.86, 51.06},
		},
		{
			Accuracy:    1,
			Helmert:     Helmert{485, 169.5, 483.8, 7.786, 4.398, 4.103, 0},
			BoundingBox: BoundingBox{16.84, 47.73, 22.56, 49.61},
		},
		{
			Accuracy:    6,
			Helmert:     Helmert{589.0, 76.0, 480.0, 0, 0, 0, 0},
			BoundingBox: BoundingBox{12.09, 47.73, 22.56, 51.06},
		},
	},
}
