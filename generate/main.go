package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Input struct {
	From int
	To   int
	Dec  int
	MinX float64
	MaxX float64
	MinY float64
	MaxY float64
}

type Output struct {
	Name  string  `json:"name"`
	From  int     `json:"from"`
	InA   float64 `json:"in_a"`
	InB   float64 `json:"in_b"`
	To    int     `json:"to"`
	WantA float64 `json:"want_a"`
	WantB float64 `json:"want_b"`
	Dec   int     `json:"dec"`
}

func main() {
	inputs := []Input{
		{
			From: 4326,
			To:   3857,
			Dec:  5,
			MinX: -180,
			MinY: -85.06,
			MaxX: 180,
			MaxY: 85.06,
		},
		{
			From: 4326,
			To:   4277,
			Dec:  3,
			MinX: -9.01,
			MinY: 49.75,
			MaxX: 2.01,
			MaxY: 61.01,
		},
		{
			From: 4326,
			To:   27700,
			Dec:  3,
			MinX: -9.01,
			MinY: 49.75,
			MaxX: 2.01,
			MaxY: 61.01,
		},
	}

	output := []Output{}

	for _, in := range inputs {
		points := generateGrid(in.MinX, in.MaxX, in.MinY, in.MaxY, in.Dec)

		for index, p := range points {
			cmd := exec.Command("cs2cs", fmt.Sprintf("+init=epsg:%d", in.From), "+to", fmt.Sprintf("+init=epsg:%d", in.To), "-d", strconv.Itoa(in.Dec))

			coords := fmt.Sprintf("%f %f 0 0", p[0], p[1])

			cmd.Stdin = strings.NewReader(coords + "\n")

			var out bytes.Buffer

			cmd.Stdout = &out
			cmd.Stderr = &out

			if err := cmd.Run(); err != nil {
				panic(out.String())
			}

			fields := strings.Fields(out.String())

			wantA, err := strconv.ParseFloat(fields[0], 64)
			if err != nil {
				panic(out.String())
			}

			wantB, err := strconv.ParseFloat(fields[1], 64)
			if err != nil {
				panic(out.String())
			}

			output = append(output, Output{
				Name:  fmt.Sprintf("%d:%d (%d)", in.From, in.To, index+1),
				From:  in.From,
				InA:   p[0],
				InB:   p[1],
				To:    in.To,
				WantA: wantA,
				WantB: wantB,
				Dec:   in.Dec,
			})
		}
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")

	if err := enc.Encode(output); err != nil {
		panic(err)
	}
}

func generateGrid(minX, maxX, minY, maxY float64, dec int) [][2]float64 {
	points := make([][2]float64, 0, 16)

	n := 4
	step := 1.0 / float64(n)
	factor := math.Pow(10, float64(dec))

	round := func(v float64) float64 {
		return math.Round(v*factor) / factor
	}

	for j := 0; j < n; j++ {
		for i := 0; i < n; i++ {
			tx := (float64(i) + 0.5) * step
			ty := (float64(j) + 0.5) * step

			x := minX + tx*(maxX-minX)
			y := minY + ty*(maxY-minY)

			points = append(points, [2]float64{
				round(x),
				round(y),
			})
		}
	}

	return points
}
