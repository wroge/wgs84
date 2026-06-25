package wgs84

var ED50 = Datum{
	Spheroid: International1924,
	Transformations: []Transformation{
		{
			Accuracy:    1,
			Helmert:     Helmert{-116.641, -56.931, -110.559, 0.893, 0.921, -0.917, -3.52},
			BoundingBox: BoundingBox{-3.35, 62.0, 38.01, 84.73},
		},
		{
			Accuracy:    1,
			Helmert:     Helmert{-89.5, -93.8, -123.1, 0, 0, -0.156, 1.2},
			BoundingBox: BoundingBox{-16.1, 47.42, 10.86, 63.89},
		},
		{
			Accuracy:    1,
			Grid:        "es_ign_SPED2ETV2.tif",
			GridOptional: true,
			BoundingBox: BoundingBox{-9.37, 35.26, 4.39, 43.82},
		},
		{
			Accuracy:    1,
			Helmert:     Helmert{-74.292, -135.889, -104.967, 0.524, 0.136, -0.61, -3.761},
			BoundingBox: BoundingBox{-9.56, 36.95, -6.19, 42.16},
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-3.35, 65.0, 38.01, 84.73},
		},
		{
			Accuracy:    1,
			Helmert:     Helmert{-90.365, -101.13, -123.384, 0.333, 0.077, 0.894, 1.994},
			BoundingBox: BoundingBox{1.37, 56.08, 10.81, 62.01},
		},
		{
			Accuracy:    1,
			Helmert:     Helmert{-81.1, -89.4, -115.8, 0.485, 0.024, 0.413, -0.54},
			BoundingBox: BoundingBox{7.98, 54.5, 15.28, 57.81},
		},
		{
			Accuracy:    1,
			Helmert:     Helmert{-157.89, -17.16, -78.41, 2.118, 2.697, -1.434, -5.38},
			BoundingBox: BoundingBox{3.34, 53.58, 8.88, 55.92},
		},
		{
			Accuracy:    1,
			Helmert:     Helmert{-116.8, -106.4, -154.4, 0, 0, 0, 0},
			BoundingBox: BoundingBox{-5.42, 36.0, -4.89, 36.16},
		},
		{
			Accuracy:    1.5,
			Helmert:     Helmert{-131, -100.3, -163.4, -1.244, -0.02, -1.144, 9.39},
			BoundingBox: BoundingBox{-7.54, 35.26, 3.39, 43.56},
		},
		{
			Accuracy:    1.5,
			Helmert:     Helmert{-178.4, -83.2, -221.3, 0.54, -0.532, -0.126, 21.2},
			BoundingBox: BoundingBox{-9.37, 41.5, -4.5, 43.82},
		},
		{
			Accuracy:    1.5,
			Helmert:     Helmert{-181.5, -90.3, -187.2, 0.144, 0.492, -0.394, 17.57},
			BoundingBox: BoundingBox{1.12, 38.59, 4.39, 40.15},
		},
		{
			Accuracy:    1.5,
			Helmert:     Helmert{82.981, -99.719, -110.709, -0.104700015651026, 0.0310016003789386, 0.0804020214751182, -0.3143},
			BoundingBox: BoundingBox{-3.35, 65.0, 38.01, 84.73},
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-5.05, 51.03, 10.86, 62.0},
		},
		{
			Accuracy:    2,
			Helmert:     Helmert{-84, -97, -117, 0, 0, 0, 0},
			BoundingBox: BoundingBox{-9.86, 41.15, 10.38, 51.56},
		},
		{
			Accuracy:    2,
			Helmert:     Helmert{-84.1, -101.8, -129.7, 0, 0, 0.468, 1.05},
			BoundingBox: BoundingBox{25.62, 34.42, 44.83, 43.45},
		},
		{
			Accuracy:    2.5,
			Helmert:     Helmert{-112, -110.3, -140.2, 0, 0, 0, 0},
			BoundingBox: BoundingBox{34.88, 29.18, 39.31, 33.38},
		},
		{
			Accuracy:    5,
			Helmert:     Helmert{-86.277, -108.879, -120.181, 0, 0, 0, 0},
			BoundingBox: BoundingBox{-13.87, 34.91, -7.24, 41.88},
		},
		{
			Accuracy:    5,
			Helmert:     Helmert{-87.987, -108.639, -121.593, 0, 0, 0, 0},
			BoundingBox: BoundingBox{-9.56, 36.95, -6.19, 42.16},
		},
		{
			Accuracy:    5,
			Helmert:     Helmert{-83.11, -97.38, -117.22, 0.0276, -0.2167, 0.2147, 0.1218},
			BoundingBox: BoundingBox{2.53, 51.44, 6.37, 55.77},
		},
		{
			Accuracy:    5,
			Helmert:     Helmert{-82.31, -95.23, -114.96, 0, 0, 0, 0},
			BoundingBox: BoundingBox{-12.5, 53.75, -9.49, 55.76},
		},
		{
			Accuracy:    6,
			Helmert:     Helmert{-87, -96, -120, 0, 0, 0, 0},
			BoundingBox: BoundingBox{-4.87, 42.33, 17.17, 57.8},
		},
		{
			Accuracy:    6,
			Helmert:     Helmert{-86, -96, -120, 0, 0, 0, 0},
			BoundingBox: BoundingBox{-10.56, 49.11, 1.84, 60.9},
		},
		{
			Accuracy:    7,
			Helmert:     Helmert{-87, -95, -120, 0, 0, 0, 0},
			BoundingBox: BoundingBox{4.39, 57.9, 31.59, 71.24},
		},
		{
			Accuracy:    9,
			Helmert:     Helmert{-84, -107, -120, 0, 0, 0, 0},
			BoundingBox: BoundingBox{-9.56, 35.26, 3.39, 43.82},
		},
		{
			Accuracy:    10,
			Helmert:     Helmert{-87, -98, -121, 0, 0, 0, 0},
			BoundingBox: BoundingBox{-9.56, 34.88, 31.59, 71.24},
		},
		{
			Accuracy:    10,
			Helmert:     Helmert{-89.05, -87.03, -124.56, 0, 0, 0, 0},
			BoundingBox: BoundingBox{28.03, 31.35, 41.47, 43.45},
		},
		{
			Accuracy:    13,
			Helmert:     Helmert{-130, -117, -151, 0, 0, 0, 0},
			BoundingBox: BoundingBox{24.7, 25.71, 30.0, 31.68},
		},
		{
			Accuracy:    15,
			Helmert:     Helmert{-84, -103, -122.5, 0, 0, 0.554, 0.2263},
			BoundingBox: BoundingBox{24.7, 25.71, 30.0, 31.68},
		},
		{
			Accuracy:    26,
			Helmert:     Helmert{-104, -101, -140, 0, 0, 0, 0},
			BoundingBox: BoundingBox{29.95, 32.88, 35.2, 36.21},
		},
		{
			Accuracy:    35,
			Helmert:     Helmert{-97, -88, -135, 0, 0, 0, 0},
			BoundingBox: BoundingBox{12.36, 36.59, 15.71, 38.35},
		},
		{
			Accuracy:    44,
			Helmert:     Helmert{-84, -95, -130, 0, 0, 0, 0},
			BoundingBox: BoundingBox{19.57, 34.88, 28.3, 41.75},
		},
		{
			Accuracy:    44,
			Helmert:     Helmert{-112, -77, -145, 0, 0, 0, 0},
			BoundingBox: BoundingBox{7.49, 30.23, 13.67, 38.41},
		},
		{
			Accuracy:    44,
			Helmert:     Helmert{-97, -103, -120, 0, 0, 0, 0},
			BoundingBox: BoundingBox{8.08, 38.82, 9.89, 41.31},
		},
		{
			Accuracy:    44,
			Helmert:     Helmert{-107, -88, -149, 0, 0, 0, 0},
			BoundingBox: BoundingBox{14.27, 35.74, 14.63, 36.05},
		},
	},
}
