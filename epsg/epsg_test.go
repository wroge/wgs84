package epsg

import (
	"fmt"
	"testing"

	"github.com/wroge/wgs84/etrs89"
	"github.com/wroge/wgs84/system"
)

func ExampleRepository_Add() {
	r := DefaultRepository()
	r.Add(3416, etrs89.WithSystem(
		system.LambertConformalConic2SP(
			13.33333333333333, 47.5, 49, 46, 400000, 400000,
		)),
		-16.1, 32.88, 40.18, 84.17,
	)
	f := r.Transform(4326, 3416)
	if f != nil {
		fmt.Println(f.Round(0)(16, 48, 0))
	}
}

func BenchmarkTransform(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Transform(4326, 3857)(9, 52, 0)
	}
}
