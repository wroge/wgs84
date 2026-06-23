package dhdn90

import (
	"embed"

	"github.com/wroge/wgs84/v2"
)

//go:embed *.tif
var grids embed.FS

func init() {
	wgs84.RegisterGridFS("", grids)
}
