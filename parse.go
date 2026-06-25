package wgs84

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseProj(proj string) (CoordinateReferenceSystem, error) {
	for pair := range strings.FieldsSeq(proj) {
		key, val, ok := strings.Cut(pair, "=epsg:")
		if !ok {
			continue
		}

		switch key {
		case "+init":
			i, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return CoordinateReferenceSystem{}, err
			}

			epsg, ok := EPSG[int(i)]
			if !ok {
				return CoordinateReferenceSystem{}, ErrCodeNotFound
			}

			return epsg, nil
		}
	}

	datum, err := parseDatum(proj)
	if err != nil {
		return CoordinateReferenceSystem{}, err
	}

	cs, err := parseCoordinateSystem(proj)
	if err != nil {
		return CoordinateReferenceSystem{}, err
	}

	return CoordinateReferenceSystem{
		Datum:            datum,
		CoordinateSystem: cs,
	}, nil
}

func parseDatum(proj string) (Datum, error) {
	var (
		S   Spheroid
		H   Helmert
		G   string
		err error
	)

	for pair := range strings.FieldsSeq(proj) {
		key, val, ok := strings.Cut(pair, "=")
		if !ok {
			continue
		}

		switch key {
		case "+ellps":
			switch val {
			case "GRS80", "grs80":
				S = GRS80
			case "bessel", "amersfoort":
				S = Bessel1841
			case "intl":
				S = International1924
			case "airy":
				S = Airy1830
			case "clrk66":
				S = Clarke1866
			case "clrk80ign":
				S = Clarke1880
			}
		case "+datum":
			switch val {
			case "WGS84", "wgs84":
				return WGS84, nil
			}
		case "+a":
			S.A, err = strconv.ParseFloat(val, 64)
			if err != nil {
				return Datum{}, err
			}
		case "+rf":
			S.Fi, err = strconv.ParseFloat(val, 64)
			if err != nil {
				return Datum{}, err
			}
		case "+nadgrids":
			G = val
		case "+towgs84":
			parts := strings.Split(val, ",")

			if len(parts) > 7 {
				return Datum{}, fmt.Errorf("invalid number of values in +towgs84: %d", len(parts))
			}

			for i, part := range parts {
				switch i {
				case 0:
					H.Tx, err = strconv.ParseFloat(part, 64)
					if err != nil {
						return Datum{}, err
					}
				case 1:
					H.Ty, err = strconv.ParseFloat(part, 64)
					if err != nil {
						return Datum{}, err
					}
				case 2:
					H.Tz, err = strconv.ParseFloat(part, 64)
					if err != nil {
						return Datum{}, err
					}
				case 3:
					H.Rx, err = strconv.ParseFloat(part, 64)
					if err != nil {
						return Datum{}, err
					}
				case 4:
					H.Ry, err = strconv.ParseFloat(part, 64)
					if err != nil {
						return Datum{}, err
					}
				case 5:
					H.Rz, err = strconv.ParseFloat(part, 64)
					if err != nil {
						return Datum{}, err
					}
				case 6:
					H.Ds, err = strconv.ParseFloat(part, 64)
					if err != nil {
						return Datum{}, err
					}
				}
			}
		}
	}

	if S.A == 0 {
		return Datum{}, ErrInvalidDatum
	}

	if S.Fi == 0 {
		S.Fi = WGS84.Spheroid.Fi
	}

	if G != "" && G != "@null" {
		var transformations []Transformation

		for part := range strings.SplitSeq(G, ",") {
			t, ok := parseGridEntry(part)
			if ok {
				transformations = append(transformations, t)
			}
		}

		if len(transformations) == 0 {
			return Datum{}, ErrInvalidDatum
		}

		return Datum{
			Spheroid:        S,
			Transformations: transformations,
		}, nil
	}

	return Datum{
		Spheroid: S,
		Transformations: []Transformation{
			{
				BoundingBox: World,
				Helmert:     H,
			},
		},
	}, nil
}

func parseGridEntry(s string) (Transformation, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return Transformation{}, false
	}

	t := Transformation{BoundingBox: World}
	if s == nullGrid {
		t.Grid = nullGrid
		return t, true
	}

	if strings.HasPrefix(s, "@") {
		t.GridOptional = true
		s = strings.TrimSpace(strings.TrimPrefix(s, "@"))
		if s == "" {
			return Transformation{}, false
		}
	}

	if s == nullGrid {
		t.Grid = nullGrid
		return t, true
	}

	t.Grid = s
	return t, true
}

func parseCoordinateSystem(proj string) (CoordinateSystem, error) {
	var (
		toCS                             func() CoordinateSystem
		Lonf, Latf, Scale, Eastf, Northf float64
		Sp1, Sp2                         float64
		Alpha                            float64
		Zone                             int64
		PrimeMeridian                    float64
		South                            = false
		err                              error
	)

	for pair := range strings.FieldsSeq(proj) {
		key, val, _ := strings.Cut(pair, "=")

		switch key {
		case "+south":
			South = true
		case "+zone":
			Zone, err = strconv.ParseInt(val, 10, 64)
			if err != nil {
				return nil, err
			}
		case "+pm":
			switch val {
			case "paris":
				PrimeMeridian = 2.33722917
			}
		case "+lat_0":
			Latf, err = strconv.ParseFloat(val, 64)
			if err != nil {
				return nil, err
			}
		case "+lon_0":
			Lonf, err = strconv.ParseFloat(val, 64)
			if err != nil {
				return nil, err
			}
		case "+k_0", "+k":
			Scale, err = strconv.ParseFloat(val, 64)
			if err != nil {
				return nil, err
			}
		case "+x_0":
			Eastf, err = strconv.ParseFloat(val, 64)
			if err != nil {
				return nil, err
			}
		case "+y_0":
			Northf, err = strconv.ParseFloat(val, 64)
			if err != nil {
				return nil, err
			}
		case "+lat_1":
			Sp1, err = strconv.ParseFloat(val, 64)
			if err != nil {
				return nil, err
			}
		case "+lat_2":
			Sp2, err = strconv.ParseFloat(val, 64)
			if err != nil {
				return nil, err
			}
		case "+alpha":
			Alpha, err = strconv.ParseFloat(val, 64)
			if err != nil {
				return nil, err
			}
		case "+proj":
			switch val {
			case "geocent":
				toCS = func() CoordinateSystem {
					return Geocentric{}
				}
			case "longlat":
				toCS = func() CoordinateSystem {
					return Geographic{}
				}
			case "merc":
				toCS = func() CoordinateSystem {
					return WebMercator{}
				}
			case "utm":
				toCS = func() CoordinateSystem {
					if South {
						Northf = 10000000
					}

					return TransverseMercator{
						Lonf:   float64(Zone)*6 - 183,
						Latf:   Latf,
						Scale:  0.9996,
						Eastf:  500000,
						Northf: Northf,
					}
				}
			case "tmerc":
				toCS = func() CoordinateSystem {
					if Scale == 0 {
						Scale = 1
					}

					return TransverseMercator{
						Lonf:   Lonf,
						Latf:   Latf,
						Scale:  Scale,
						Eastf:  Eastf,
						Northf: Northf,
					}
				}
			case "aea":
				toCS = func() CoordinateSystem {
					return AlbersConicEqualArea{
						Lonf:   Lonf,
						Latf:   Latf,
						Sp1:    Sp1,
						Sp2:    Sp2,
						Eastf:  Eastf,
						Northf: Northf,
					}
				}
			case "somerc":
				toCS = func() CoordinateSystem {
					return SwissObliqueMercator{
						Lonf:   Lonf,
						Latf:   Latf,
						Scale:  Scale,
						Eastf:  Eastf,
						Northf: Northf,
					}
				}
			case "krovak":
				toCS = func() CoordinateSystem {
					return Krovak{
						Lonf:   Lonf,
						Latf:   Latf,
						Alpha:  Alpha,
						Scale:  Scale,
						Eastf:  Eastf,
						Northf: Northf,
					}
				}
			case "laea":
				toCS = func() CoordinateSystem {
					return LambertAzimuthalEqualArea{
						Lonf:   Lonf,
						Latf:   Latf,
						Eastf:  Eastf,
						Northf: Northf,
					}
				}
			case "lcc":
				toCS = func() CoordinateSystem {
					if Sp2 != 0 {
						return LambertConformalConic2SP{
							Lonf:   Lonf,
							Latf:   Latf,
							Sp1:    Sp1,
							Sp2:    Sp2,
							Eastf:  Eastf,
							Northf: Northf,
						}
					}

					return LambertConformalConic1SP{
						Lonf:   Lonf,
						Latf:   Latf,
						Scale:  Scale,
						Eastf:  Eastf,
						Northf: Northf,
					}
				}
			default:
				return nil, fmt.Errorf("proj '%s' is not implemented", val)
			}
		}
	}

	if Lonf == 0 && PrimeMeridian != 0 {
		Lonf = PrimeMeridian
	}

	if toCS == nil {
		return nil, fmt.Errorf("proj '%s' is not implemented", proj)
	}

	return toCS(), nil
}
