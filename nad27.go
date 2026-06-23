package wgs84

var NAD27 = Datum{
	Spheroid: Clarke1866,
	Transformations: []Transformation{
		{
			Accuracy:    1,
			BoundingBox: BoundingBox{-87.01, 18.83, -73.57, 25.51},
			Helmert:     Helmert{2.478, 149.752, 197.726, -0.526, -0.498, 0.501, 0.685},
		},
		{
			Accuracy:    1,
			BoundingBox: BoundingBox{-83.04, 7.15, -77.19, 9.68},
			Helmert:     Helmert{-32.3841359, 180.4090461, 120.8442577, 2.1545854, 0.1498782, -0.5742915, 8.1049164},
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-79.85, 44.99, -57.1, 62.62},
			Grid:        "ca_nrc_NA27SCRS.tif",
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-110.0, 49.0, -101.34, 60.01},
			Grid:        "ca_nrc_SK27-98.tif",
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-100.0, 25.83, -93.5, 34.58},
			Grid:        "us_noaa_ethpgn.tif",
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-106.66, 28.04, -100.0, 36.5},
			Grid:        "us_noaa_wthpgn.tif",
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-87.63, 24.41, -79.97, 31.01},
			Grid:        "us_noaa_FL.tif",
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-113.0, 41.99, -104.04, 49.01},
			Grid:        "us_noaa_emhpgn.tif",
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-124.79, 41.98, -116.47, 49.05},
			Grid:        "us_noaa_WO.tif",
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-90.42, 41.69, -82.13, 48.32},
			Grid:        "us_noaa_mihpgn.tif",
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-124.45, 36.5, -116.54, 42.01},
			Grid:        "us_noaa_cnhpgn.tif",
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-120.0, 34.99, -114.03, 42.0},
			Grid:        "us_noaa_nvhpgn.tif",
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-97.22, 43.49, -89.49, 49.38},
			Grid:        "us_noaa_mnhpgn.tif",
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-109.06, 31.33, -102.99, 37.0},
			Grid:        "us_noaa_nmhpgn.tif",
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-79.77, 40.47, -71.8, 45.02},
			Grid:        "us_noaa_nyhpgn.tif",
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-121.98, 32.53, -114.12, 36.5},
			Grid:        "us_noaa_cshpgn.tif",
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-95.77, 35.98, -89.1, 40.61},
			Grid:        "us_noaa_mohpgn.tif",
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-103.0, 33.62, -94.42, 37.01},
			Grid:        "us_noaa_okhpgn.tif",
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-92.89, 42.48, -86.25, 47.31},
			Grid:        "us_noaa_WI.tif",
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-109.06, 36.98, -102.04, 41.01},
			Grid:        "us_noaa_cohpgn.tif",
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-117.24, 41.99, -113.0, 49.01},
			Grid:        "us_noaa_wmhpgn.tif",
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-111.06, 40.99, -104.05, 45.01},
			Grid:        "us_noaa_wyhpgn.tif",
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-104.06, 39.99, -95.3, 43.01},
			Grid:        "us_noaa_nbhpgn.tif",
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-114.05, 36.99, -109.04, 42.01},
			Grid:        "us_noaa_uthpgn.tif",
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-83.68, 36.54, -75.31, 39.46},
			Grid:        "us_noaa_vahpgn.tif",
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-91.52, 36.97, -87.02, 42.51},
			Grid:        "us_noaa_ilhpgn.tif",
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-104.07, 42.48, -96.43, 45.95},
			Grid:        "us_noaa_sdhpgn.tif",
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-94.05, 28.85, -88.75, 33.03},
			Grid:        "us_noaa_lahpgn.tif",
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-85.61, 30.36, -80.77, 35.01},
			Grid:        "us_noaa_gahpgn.tif",
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-102.06, 36.99, -94.58, 40.01},
			Grid:        "us_noaa_kshpgn.tif",
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-89.57, 36.49, -81.95, 39.15},
			Grid:        "us_noaa_kyhpgn.tif",
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-104.07, 45.93, -96.55, 49.01},
			Grid:        "us_noaa_ndhpgn.tif",
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-96.65, 40.36, -90.14, 43.51},
			Grid:        "us_noaa_iahpgn.tif",
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-91.65, 30.01, -88.09, 35.01},
			Grid:        "us_noaa_mshpgn.tif",
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-88.48, 30.14, -84.89, 35.02},
			Grid:        "us_noaa_alhpgn.tif",
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-94.62, 33.01, -89.64, 36.5},
			Grid:        "us_noaa_arhpgn.tif",
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-82.65, 37.2, -77.72, 40.64},
			Grid:        "us_noaa_wvhpgn.tif",
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-71.09, 43.04, -66.91, 47.47},
			Grid:        "us_noaa_mehpgn.tif",
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-84.83, 38.4, -80.51, 42.33},
			Grid:        "us_noaa_ohhpgn.tif",
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-73.73, 40.98, -69.86, 45.31},
			Grid:        "us_noaa_nehpgn.tif",
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-90.31, 34.98, -81.65, 36.68},
			Grid:        "us_noaa_TN.tif",
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-88.1, 37.77, -84.78, 41.77},
			Grid:        "us_noaa_inhpgn.tif",
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-79.49, 37.97, -74.97, 39.85},
			Grid:        "us_noaa_MD.tif",
		},
		{
			Accuracy:    1.5,
			BoundingBox: BoundingBox{-75.6, 38.87, -73.88, 41.36},
			Grid:        "us_noaa_njhpgn.tif",
		},
		{
			Accuracy:    2,
			BoundingBox: BoundingBox{-141.01, 40.0, -44.0, 83.17},
			Grid:        "ca_nrc_ntv2_0.tif",
		},
		{
			Accuracy:    2.5,
			BoundingBox: BoundingBox{-120.0, 48.99, -109.98, 60.0},
			Grid:        "ca_nrc_ABCSRSV4.tif",
		},
		{
			Accuracy:    3,
			BoundingBox: BoundingBox{-67.75, 40.0, -43.99, 64.21},
			Helmert:     Helmert{-0.991, 1.9072, 0.5129, 0.0257899075194932, 0.0096500989602704, 0.0116599432323421, 0},
		},
		{
			Accuracy:    3,
			BoundingBox: BoundingBox{-87.01, 18.83, -73.57, 25.51},
			Helmert:     Helmert{-4.2, 135.4, 181.9, 0, 0, 0, 0},
		},
		{
			Accuracy:    5,
			BoundingBox: BoundingBox{-129.17, 23.81, -65.69, 49.38},
			Grid:        "us_noaa_conus.tif",
		},
		{
			Accuracy:    5,
			BoundingBox: BoundingBox{167.65, 47.88, -129.99, 74.71},
			Grid:        "us_noaa_alaska.tif",
		},
		{
			Accuracy:    5,
			BoundingBox: BoundingBox{-87.25, 23.82, -81.17, 30.25},
			Helmert:     Helmert{-3, 154, 177, 0, 0, 0, 0},
		},
		{
			Accuracy:    5,
			BoundingBox: BoundingBox{-95.0, 25.61, -87.25, 30.23},
			Helmert:     Helmert{-7, 151, 175, 0, 0, 0, 0},
		},
		{
			Accuracy:    5,
			BoundingBox: BoundingBox{-94.79, 17.85, -89.75, 20.89},
			Helmert:     Helmert{-2, 124.7, 196, 0, 0, 0, 0},
		},
		{
			Accuracy:    5,
			BoundingBox: BoundingBox{-94.33, 20.87, -88.67, 23.01},
			Helmert:     Helmert{0, 125, 196, 0, 0, 0, 0},
		},
		{
			Accuracy:    5,
			BoundingBox: BoundingBox{-97.22, 25.97, -95.0, 28.97},
			Helmert:     Helmert{-7, 151, 178, 0, 0, 0, 0},
		},
		{
			Accuracy:    5,
			BoundingBox: BoundingBox{-98.1, 21.51, -96.89, 22.75},
			Helmert:     Helmert{-8, 125, 190, 0, 0, 0, 0},
		},
		{
			Accuracy:    7,
			BoundingBox: BoundingBox{-124.79, 25.83, -89.64, 49.05},
			Helmert:     Helmert{-8, 159, 175, 0, 0, 0, 0},
		},
		{
			Accuracy:    8,
			BoundingBox: BoundingBox{-136.46, 49.0, -60.72, 83.17},
			Helmert:     Helmert{4, 159, 188, 0, 0, 0, 0},
		},
		{
			Accuracy:    8,
			BoundingBox: BoundingBox{-97.22, 23.82, -81.17, 30.25},
			Helmert:     Helmert{-7, 158, 172, 0, 0, 0, 0},
		},
		{
			Accuracy:    8,
			BoundingBox: BoundingBox{-79.04, 20.86, -72.68, 27.29},
			Helmert:     Helmert{-4, 154, 178, 0, 0, 0, 0},
		},
		{
			Accuracy:    9,
			BoundingBox: BoundingBox{-79.85, 43.41, -52.54, 62.62},
			Helmert:     Helmert{-22, 160, 190, 0, 0, 0, 0},
		},
		{
			Accuracy:    10,
			BoundingBox: BoundingBox{-124.79, 24.41, -66.91, 49.38},
			Helmert:     Helmert{-8, 160, 176, 0, 0, 0, 0},
		},
		{
			Accuracy:    10,
			BoundingBox: BoundingBox{-92.29, 7.98, -82.53, 18.49},
			Helmert:     Helmert{0, 125, 194, 0, 0, 0, 0},
		},
		{
			Accuracy:    10,
			BoundingBox: BoundingBox{-141.01, 59.99, -123.91, 69.7},
			Helmert:     Helmert{-7, 139, 181, 0, 0, 0, 0},
		},
		{
			Accuracy:    11,
			BoundingBox: BoundingBox{-97.22, 24.41, -66.91, 49.38},
			Helmert:     Helmert{-9, 161, 179, 0, 0, 0, 0},
		},
		{
			Accuracy:    12,
			BoundingBox: BoundingBox{-118.47, 14.51, -86.68, 32.72},
			Helmert:     Helmert{-12, 130, 190, 0, 0, 0, 0},
		},
		{
			Accuracy:    12,
			BoundingBox: BoundingBox{-102.0, 41.67, -74.35, 60.01},
			Helmert:     Helmert{-9, 157, 184, 0, 0, 0, 0},
		},
		{
			Accuracy:    12,
			BoundingBox: BoundingBox{-168.26, 54.34, -129.99, 71.4},
			Helmert:     Helmert{-5, 135, 172, 0, 0, 0, 0},
		},
		{
			Accuracy:    13,
			BoundingBox: BoundingBox{-139.04, 48.25, -109.98, 60.01},
			Helmert:     Helmert{-7, 162, 188, 0, 0, 0, 0},
		},
		{
			Accuracy:    15,
			BoundingBox: BoundingBox{-178.3, 51.54, -164.84, 54.34},
			Helmert:     Helmert{-2, 152, 149, 0, 0, 0, 0},
		},
		{
			Accuracy:    16,
			BoundingBox: BoundingBox{-85.01, 13.0, -59.37, 23.25},
			Helmert:     Helmert{-3, 142, 183, 0, 0, 0, 0},
		},
		{
			Accuracy:    18,
			BoundingBox: BoundingBox{172.42, 51.3, 179.86, 53.07},
			Helmert:     Helmert{2, 204, 105, 0, 0, 0, 0},
		},
		{
			Accuracy:    20,
			BoundingBox: BoundingBox{-141.01, 40.0, -44.0, 83.17},
			Helmert:     Helmert{-10, 158, 187, 0, 0, 0, 0},
		},
		{
			Accuracy:    35,
			BoundingBox: BoundingBox{-80.07, 8.82, -79.46, 9.45},
			Helmert:     Helmert{0, 125, 201, 0, 0, 0, 0},
		},
		{
			Accuracy:    44,
			BoundingBox: BoundingBox{-85.01, 19.77, -74.07, 23.25},
			Helmert:     Helmert{-9, 152, 178, 0, 0, 0, 0},
		},
		{
			Accuracy:    44,
			BoundingBox: BoundingBox{-73.29, 75.86, -60.98, 79.2},
			Helmert:     Helmert{11, 114, 195, 0, 0, 0, 0},
		},
		{
			Accuracy:    44,
			BoundingBox: BoundingBox{-74.6, 23.9, -74.37, 24.19},
			Helmert:     Helmert{1, 140, 165, 0, 0, 0, 0},
		},
	},
}
