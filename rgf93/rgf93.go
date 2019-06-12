package rgf93

import (
	"github.com/wroge/wgs84"
	"github.com/wroge/wgs84/spheroid"
	"github.com/wroge/wgs84/system"
	"github.com/wroge/wgs84/transformation"
)

// FranceLambert is equal to EPSG-Code 2154.
func FranceLambert() wgs84.CoordinateReferenceSystem {
	return WithSystem(system.LambertConformalConic2SP(3, 46.5, 49, 44, 700000, 6600000))
}

// CC(42) is equal to EPSG-Code 3942
// CC(50) is equal to EPSG-Code 3950 ...
func CC(latitude float64) wgs84.CoordinateReferenceSystem {
	return WithSystem(system.LambertConformalConic2SP(3, latitude, latitude-0.75, latitude+0.75, 1700000, 2200000+(latitude-43)*1000000))
}

// WithSystem provides any wgs84.System in the RGF93
// geodetic datum.
func WithSystem(sys wgs84.System) wgs84.CoordinateReferenceSystem {
	return wgs84.CoordinateReferenceSystem{
		Spheroid:       spheroid.GRS80(),
		Transformation: transformation.WGS84(),
		System:         sys,
	}
}
