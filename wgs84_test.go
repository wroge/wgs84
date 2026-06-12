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
