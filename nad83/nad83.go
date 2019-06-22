// Package nad83 implements wgs84.CoordinateReferenceSystem's
// in the NAD83 geodetic datum.
package nad83

import (
	"github.com/wroge/wgs84"
	"github.com/wroge/wgs84/spheroid"
	"github.com/wroge/wgs84/system"
	"github.com/wroge/wgs84/transformation"
)

// Austin is equal to EPSG-Code 100002
func Austin() wgs84.CoordinateReferenceSystem {
	return WithSystem(system.LambertConformalConic2SP(-100.333333333333, 29.6666666666667, 31.8833333333333, 30.1166666666667, 2296583.333333, 9842500.0000000))
}

// Vermont is equal to EPSG-Code 32145
func Vermont() wgs84.CoordinateReferenceSystem {
	return WithSystem(system.TransverseMercator(-72.5, 42.5, 0.999964286, 500000, 0))
}

// Tennessee is equal to EPSG-Code 32136
func Tennessee() wgs84.CoordinateReferenceSystem {
	return WithSystem(system.LambertConformalConic2SP(-86, 34.33333333333334, 36.41666666666666, 35.25, 600000, 0))
}

// Nebraska is equal to EPSG-Code 32104
func Nebraska() wgs84.CoordinateReferenceSystem {
	return WithSystem(system.LambertConformalConic2SP(-100, 39.83333333333334, 43, 40, 500000, 0))
}

// Montana is equal to EPSG-Code 32100
func Montana() wgs84.CoordinateReferenceSystem {
	return WithSystem(system.LambertConformalConic2SP(-109.5, 44.25, 49, 45, 600000, 0))
}

// Maryland is equal to EPSG-Code 26985
func Maryland() wgs84.CoordinateReferenceSystem {
	return WithSystem(system.LambertConformalConic2SP(-77, 37.66666666666666, 39.45, 38.3, 400000, 0))
}

// Delaware is equal to EPSG-Code 26957
func Delaware() wgs84.CoordinateReferenceSystem {
	return WithSystem(system.TransverseMercator(-75.41666666666667, 38, 0.999995, 200000, 0))
}

// Connecticut is equal to EPSG-Code 26956
func Connecticut() wgs84.CoordinateReferenceSystem {
	return WithSystem(system.LambertConformalConic2SP(-72.75, 40.83333333333334, 41.86666666666667, 41.2, 304800.6096, 152400.3048))
}

// WithSystem provides any wgs84.System in the NAD83
// geodetic datum.
func WithSystem(sys wgs84.System) wgs84.CoordinateReferenceSystem {
	return wgs84.CoordinateReferenceSystem{
		Spheroid:       spheroid.GRS80(),
		Transformation: transformation.WGS84(),
		System:         sys,
	}
}
