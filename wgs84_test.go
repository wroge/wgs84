package wgs84

import (
	"encoding/json"
	"os"
	"testing"
)

type Input struct {
	Name  string  `json:"name"`
	From  int     `json:"from"`
	FromA float64 `json:"from_a"`
	FromB float64 `json:"from_b"`
	To    int     `json:"to"`
	ToA   float64 `json:"to_a"`
	ToB   float64 `json:"to_b"`
	Dec   int     `json:"dec"`
}

func TestTransform(t *testing.T) {
	file, err := os.Open("./generate/data.json")
	if err != nil {
		t.Fatal(err)

		return
	}

	//nolint:errcheck
	defer file.Close()

	var inputs []Input

	err = json.NewDecoder(file).Decode(&inputs)
	if err != nil {
		t.Fatal(err)

		return
	}

	for _, in := range inputs {
		t.Run(in.Name, func(t *testing.T) {
			fromEPSG := EPSG(in.From)
			toEPSG := EPSG(in.To)

			transform := Transform(fromEPSG, toEPSG).Round(in.Dec)

			gotA, gotB, gotC := transform(in.FromA, in.FromB, 0)

			if gotA != in.ToA {
				t.Errorf("Transform() A = %v, want %v", gotA, in.ToA)
			}
			if gotB != in.ToB {
				t.Errorf("Transform() B = %v, want %v", gotB, in.ToB)
			}

			backtest := Transform(toEPSG, fromEPSG).Round(in.Dec)

			inA, inB, _ := backtest(gotA, gotB, gotC)

			if inA != in.FromA {
				t.Errorf("backtest() A = %v, want %v", inA, in.FromA)
			}
			if inB != in.FromB {
				t.Errorf("backtest() B = %v, want %v", inB, in.FromB)
			}
		})
	}
}

// TestLambertConformalConic2SP checks the projection against the authoritative
// worked example for EPSG method 9802 (https://epsg.io/9802-method):
// NAD27 / Texas South Central, Clarke 1866, computed in US survey feet.
func TestLambertConformalConic2SP(t *testing.T) {
	// Clarke 1866 semi-major axis 6378206.400 m expressed in US survey feet.
	clarke1866 := NewSpheroid(20925832.16, 294.97870)

	texas := LambertConformalConic2SP(
		Geographic(nil, clarke1866),
		-99,        // longitude of false origin (99°W)
		27+50.0/60, // latitude of false origin (27°50'N)
		28+23.0/60, // first standard parallel (28°23'N)
		30+17.0/60, // second standard parallel (30°17'N)
		2000000,    // false easting (US survey feet)
		0,          // false northing
	)

	t.Run("FromBase forward", func(t *testing.T) {
		east, north, _ := texas.FromBase(-96, 28.5, 0)
		if got := round(east, 2); got != 2963503.91 {
			t.Errorf("east = %v, want 2963503.91", got)
		}
		if got := round(north, 2); got != 254759.80 {
			t.Errorf("north = %v, want 254759.80", got)
		}
	})

	t.Run("ToBase inverse", func(t *testing.T) {
		lon, lat, _ := texas.ToBase(2963503.91, 254759.80, 0)
		if got := round(lon, 6); got != -96 {
			t.Errorf("lon = %v, want -96", got)
		}
		if got := round(lat, 6); got != 28.5 {
			t.Errorf("lat = %v, want 28.5", got)
		}
	})
}
