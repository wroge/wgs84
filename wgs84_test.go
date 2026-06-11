package wgs84

import (
	"encoding/json"
	"os"
	"testing"
)

type Input struct {
	Name  string  `json:"name"`
	From  int     `json:"from"`
	InA   float64 `json:"in_a"`
	InB   float64 `json:"in_b"`
	To    int     `json:"to"`
	WantA float64 `json:"want_a"`
	WantB float64 `json:"want_b"`
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

			gotA, gotB, gotC := transform(in.InA, in.InB, 0)

			if gotA != in.WantA {
				t.Errorf("Transform() A = %v, want %v", gotA, in.WantA)
			}
			if gotB != in.WantB {
				t.Errorf("Transform() B = %v, want %v", gotB, in.WantB)
			}

			backtest := Transform(toEPSG, fromEPSG).Round(in.Dec)

			inA, inB, _ := backtest(gotA, gotB, gotC)

			if inA != in.InA {
				t.Errorf("backtest() A = %v, want %v", inA, in.InA)
			}
			if inB != in.InB {
				t.Errorf("backtest() B = %v, want %v", inB, in.InB)
			}
		})
	}
}
