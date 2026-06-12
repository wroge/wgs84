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
	FromA float64 `json:"from_a"`
	FromB float64 `json:"from_b"`
	To    int     `json:"to"`
	ToA   float64 `json:"to_a"`
	ToB   float64 `json:"to_b"`
	Dec   int     `json:"dec"`
}

func main() {
	inputs := []Input{
		{
			From: 4326,
			To:   27700,
			Dec:  4,
			MinX: -9.01,
			MinY: 49.75,
			MaxX: 2.01,
			MaxY: 61.01,
		},
		{
			From: 4326,
			To:   2056,
			Dec:  3,
			MinX: 5.95,
			MinY: 45.81,
			MaxX: 10.5,
			MaxY: 47.81,
		},
		{
			From: 4326,
			To:   21781,
			Dec:  3,
			MinX: 5.95,
			MinY: 45.81,
			MaxX: 10.5,
			MaxY: 47.81,
		},
		{
			From: 4326,
			To:   2154,
			Dec:  2,
			MinX: -9.86,
			MinY: 41.15,
			MaxX: 10.38,
			MaxY: 51.56,
		},
		{
			From: 4326,
			To:   2157,
			Dec:  2,
			MinX: -10.56,
			MinY: 51.39,
			MaxX: -5.34,
			MaxY: 55.43,
		},
		{
			From: 4326,
			To:   2158,
			Dec:  1,
			MinX: -10.56,
			MinY: 51.39,
			MaxX: -5.34,
			MaxY: 55.43,
		},
		{
			From: 4326,
			To:   3035,
			Dec:  2,
			MinX: -16.1,
			MinY: 33.26,
			MaxX: 38.01,
			MaxY: 84.73,
		},
		{
			From: 4326,
			To:   3126,
			Dec:  2,
			MinX: 19.08,
			MinY: 58.84,
			MaxX: 31.59,
			MaxY: 70.09,
		},
		{
			From: 4326,
			To:   3127,
			Dec:  2,
			MinX: 19.08,
			MinY: 58.84,
			MaxX: 31.59,
			MaxY: 70.09,
		},
		{
			From: 4326,
			To:   3128,
			Dec:  2,
			MinX: 19.08,
			MinY: 58.84,
			MaxX: 31.59,
			MaxY: 70.09,
		},
		{
			From: 4326,
			To:   3129,
			Dec:  1,
			MinX: 19.08,
			MinY: 58.84,
			MaxX: 31.59,
			MaxY: 70.09,
		},
		{
			From: 4326,
			To:   3130,
			Dec:  2,
			MinX: 19.08,
			MinY: 58.84,
			MaxX: 31.59,
			MaxY: 70.09,
		},
		{
			From: 4326,
			To:   3131,
			Dec:  2,
			MinX: 19.08,
			MinY: 58.84,
			MaxX: 31.59,
			MaxY: 70.09,
		},
		{
			From: 4326,
			To:   3132,
			Dec:  1,
			MinX: 19.08,
			MinY: 58.84,
			MaxX: 31.59,
			MaxY: 70.09,
		},
		{
			From: 4326,
			To:   3133,
			Dec:  2,
			MinX: 19.08,
			MinY: 58.84,
			MaxX: 31.59,
			MaxY: 70.09,
		},
		{
			From: 4326,
			To:   3134,
			Dec:  2,
			MinX: 19.08,
			MinY: 58.84,
			MaxX: 31.59,
			MaxY: 70.09,
		},
		{
			From: 4326,
			To:   3135,
			Dec:  2,
			MinX: 19.08,
			MinY: 58.84,
			MaxX: 31.59,
			MaxY: 70.09,
		},
		{
			From: 4326,
			To:   3136,
			Dec:  2,
			MinX: 19.08,
			MinY: 58.84,
			MaxX: 31.59,
			MaxY: 70.09,
		},
		{
			From: 4326,
			To:   3137,
			Dec:  2,
			MinX: 19.08,
			MinY: 58.84,
			MaxX: 31.59,
			MaxY: 70.09,
		},
		{
			From: 4326,
			To:   3138,
			Dec:  1,
			MinX: 19.08,
			MinY: 58.84,
			MaxX: 31.59,
			MaxY: 70.09,
		},
		// {
		// 	From: 4326,
		// 	To:   3161,
		// 	Dec:  0,
		// 	MinX: -172.54,
		// 	MinY: 23.81,
		// 	MaxX: -47.74,
		// 	MaxY: 86.46,
		// },
		{
			From: 4326,
			To:   3416,
			Dec:  1,
			MinX: -16.1,
			MinY: 33.26,
			MaxX: 38.01,
			MaxY: 84.73,
		},
		{
			From: 4326,
			To:   3857,
			Dec:  7,
			MinX: -180,
			MinY: -85.06,
			MaxX: 180,
			MaxY: 85.06,
		},
		{
			From: 4326,
			To:   4156,
			Dec:  2,
			MinX: 12.09,
			MinY: 47.73,
			MaxX: 22.56,
			MaxY: 51.06,
		},
		{
			From: 4326,
			To:   4171,
			Dec:  7,
			MinX: -9.86,
			MinY: 41.15,
			MaxX: 10.38,
			MaxY: 51.56,
		},
		{
			From: 4326,
			To:   4173,
			Dec:  7,
			MinX: -10.56,
			MinY: 51.39,
			MaxX: -5.34,
			MaxY: 55.43,
		},
		{
			From: 4326,
			To:   4188,
			Dec:  5,
			MinX: -8.18,
			MinY: 53.96,
			MaxX: -5.34,
			MaxY: 55.36,
		},
		{
			From: 4326,
			To:   4230,
			Dec:  2,
			MinX: -9.56,
			MinY: 34.88,
			MaxX: 31.59,
			MaxY: 71.24,
		},
		{
			From: 4326,
			To:   4258,
			Dec:  7,
			MinX: -16.1,
			MinY: 33.26,
			MaxX: 38.0,
			MaxY: 84.73,
		},
		// {
		// 	From: 4326,
		// 	To:   4267,
		// 	Dec:  5,
		// 	MinX: -141.01,
		// 	MinY: 40.0,
		// 	MaxX: -44.0,
		// 	MaxY: 83.17,
		// },
		{
			From: 4326,
			To:   4269,
			Dec:  5,
			MinX: -172.54,
			MinY: 23.81,
			MaxX: -47.74,
			MaxY: 86.46,
		},
		{
			From: 4326,
			To:   4277,
			Dec:  6,
			MinX: -9.01,
			MinY: 49.75,
			MaxX: 2.01,
			MaxY: 61.01,
		},
		{
			From: 4326,
			To:   4299,
			Dec:  6,
			MinX: -10.56,
			MinY: 51.39,
			MaxX: -5.34,
			MaxY: 55.43,
		},
		{
			From: 4326,
			To:   4312,
			Dec:  4,
			MinX: 9.53,
			MinY: 46.4,
			MaxX: 17.17,
			MaxY: 49.02,
		},
		{
			From: 4326,
			To:   4313,
			Dec:  4,
			MinX: 2.5,
			MinY: 49.5,
			MaxX: 6.4,
			MaxY: 51.51,
		},
		{
			From: 4326,
			To:   4314,
			Dec:  7,
			MinX: 5.86,
			MinY: 47.27,
			MaxX: 15.04,
			MaxY: 55.09,
		},
		{
			From: 4326,
			To:   4549,
			Dec:  3,
			MinX: 118.5,
			MinY: 24.43,
			MaxX: 121.5,
			MaxY: 53.33,
		},
		// {
		// 	From: 4326,
		// 	To:   5514,
		// 	Dec:  0,
		// 	MinX: 12.09,
		// 	MinY: 47.73,
		// 	MaxX: 22.56,
		// 	MaxY: 51.06,
		// },
	}

	output := []Output{}

	for _, in := range inputs {
		points := generateGrid(in.MinX, in.MaxX, in.MinY, in.MaxY, in.Dec)

		for index, p := range points {
			cmd := exec.Command("cs2cs", fmt.Sprintf("+init=epsg:%d", in.From), "+to", fmt.Sprintf("+init=epsg:%d", in.To), "-d", strconv.Itoa(max(0, in.Dec)))

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
				FromA: p[0],
				FromB: p[1],
				To:    in.To,
				ToA:   wantA,
				ToB:   wantB,
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

	for j := range n {
		for i := range n {
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
