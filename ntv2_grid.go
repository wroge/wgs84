package wgs84

import (
	"fmt"
	"io/fs"
	"math"
	"os"
	"path"
	"sort"
	"strings"
	"sync"
)

type gridFilesystem struct {
	prefix string
	fsys   fs.FS
}

var (
	gridFSMu sync.RWMutex
	gridFS   = []gridFilesystem{} // sorted by prefix length, longest first
)

// RegisterGridFS adds a filesystem that provides grid files. prefix is prepended
// to LoadGrid filenames when opening from fsys (use "" for flat embeds).
func RegisterGridFS(prefix string, fsys fs.FS) {
	if fsys == nil {
		panic("wgs84: RegisterGridFS: fsys is nil")
	}

	if prefix == "/" {
		prefix = ""
	}

	gridFSMu.Lock()
	defer gridFSMu.Unlock()

	gridFS = append(gridFS, gridFilesystem{prefix: prefix, fsys: fsys})
	sortGridFSByPrefixLen()
}

func sortGridFSByPrefixLen() {
	sort.Slice(gridFS, func(i, j int) bool {
		return len(gridFS[i].prefix) > len(gridFS[j].prefix)
	})
}

func gridFilename(p string) string {
	return path.Base(strings.ReplaceAll(p, `\`, `/`))
}

func openGridFile(name string) (fs.File, error) {
	for _, rel := range []string{name, "grid/" + name} {
		if f, err := os.Open(rel); err == nil {
			return f, nil
		}
	}

	gridFSMu.RLock()
	defer gridFSMu.RUnlock()

	for _, reg := range gridFS {
		if f, err := reg.fsys.Open(reg.prefix + name); err == nil {
			return f, nil
		}
	}

	return nil, fmt.Errorf("wgs84: grid %q not found", name)
}

var gridStore sync.Map

func LoadGrid(p string) (*Grid, error) {
	p = gridFilename(p)

	if v, ok := gridStore.Load(p); ok {
		return v.(*Grid), nil
	}

	f, err := openGridFile(p)
	if err != nil {
		return nil, err
	}

	defer f.Close() //nolint:errcheck

	var grid *Grid
	switch strings.ToLower(path.Ext(p)) {
	case ".tif":
		grid, err = ParseTaggedImageFileFormat(f)
	case ".gsb":
		grid, err = ParseGridShiftBinary(f)
	default:
		return nil, fmt.Errorf("wgs84: unknown grid extension %q", path.Ext(p))
	}
	if err != nil {
		return nil, err
	}

	gridStore.Store(p, grid)
	return grid, nil
}

type Grid struct {
	SubGrids []SubGrid
}

type SubGrid struct {
	Name    string
	Parent  string
	Columns int
	Rows    int
	SLat    float64
	NLat    float64
	ELong   float64
	WLong   float64
	LatInc  float64
	LongInc float64
	Values  [][2]float32
}

func (g *Grid) ToWGS84(lon, lat float64) (float64, float64) {
	dlon, dlat := g.Shift(lon, lat)

	return lon + dlon, lat + dlat
}

func (g *Grid) FromWGS84(lon, lat float64) (float64, float64) {
	qlon, qlat := lon, lat

	for range 10 {
		dlon, dlat := g.Shift(qlon, qlat)
		newLon := lon - dlon
		newLat := lat - dlat

		if math.Abs(newLon-qlon) < 1e-12 && math.Abs(newLat-qlat) < 1e-12 {
			break
		}

		qlon, qlat = newLon, newLat
	}

	return qlon, qlat
}

func (g *Grid) Shift(lon, lat float64) (dlon, dlat float64) {
	idx := g.selectSubgrid(lon, lat)
	if idx < 0 {
		return 0, 0
	}

	return g.SubGrids[idx].shift(lon, lat)
}

func (g *Grid) selectSubgrid(lon, lat float64) int {
	lam := -lon * 3600
	phi := lat * 3600

	for i, sg := range g.SubGrids {
		if sg.Parent != "" && sg.Parent != "NONE" {
			continue
		}
		if sg.contains(phi, lam) {
			return g.deepestSubgrid(i, phi, lam)
		}
	}

	return -1
}

func (g *Grid) deepestSubgrid(idx int, phi, lam float64) int {
	best := idx
	name := g.SubGrids[idx].Name

	for i, sg := range g.SubGrids {
		if sg.Parent != name || !sg.contains(phi, lam) {
			continue
		}

		if deeper := g.deepestSubgrid(i, phi, lam); deeper >= 0 {
			best = deeper
		}
	}

	return best
}

func (sg SubGrid) contains(phi, lam float64) bool {
	return phi >= sg.SLat && phi <= sg.NLat && lam >= sg.ELong && lam <= sg.WLong
}

func (sg SubGrid) shift(lon, lat float64) (dlon, dlat float64) {
	if sg.Columns < 2 || sg.Rows < 2 || len(sg.Values) == 0 {
		return 0, 0
	}

	fcol := (-lon*3600 - sg.ELong) / sg.LongInc
	frow := (lat*3600 - sg.SLat) / sg.LatInc

	col := math.Floor(fcol)
	row := math.Floor(frow)

	ppr := float64(sg.Columns)
	ppc := float64(sg.Rows)

	se := row*ppr + col
	sw := se + 1
	ne := se + ppr
	nw := ne + 1

	col = math.Max(0, math.Min(col, ppr-2))
	row = math.Max(0, math.Min(row, ppc-2))

	if col >= ppr-1 {
		sw = se
		nw = ne
	}

	if row >= ppc-1 {
		ne = se
		nw = sw
	}

	if col <= 0 {
		se = sw
		ne = nw
	}

	if row <= 0 {
		se = ne
		sw = nw
	}

	seIndex := clampInt(int(se), 0, len(sg.Values)-1)
	swIndex := clampInt(int(sw), 0, len(sg.Values)-1)
	neIndex := clampInt(int(ne), 0, len(sg.Values)-1)
	nwIndex := clampInt(int(nw), 0, len(sg.Values)-1)

	sse := sg.Values[seIndex]
	ssw := sg.Values[swIndex]
	sne := sg.Values[neIndex]
	snw := sg.Values[nwIndex]

	dx := fcol - col
	dy := frow - row

	latsv := (1-dx)*(1-dy)*float64(sse[0]) + dx*(1-dy)*float64(ssw[0]) +
		(1-dx)*dy*float64(sne[0]) + dx*dy*float64(snw[0])
	lonsv := (1-dx)*(1-dy)*float64(sse[1]) + dx*(1-dy)*float64(ssw[1]) +
		(1-dx)*dy*float64(sne[1]) + dx*dy*float64(snw[1])

	return -lonsv / 3600, latsv / 3600
}

func (sg SubGrid) validate() error {
	want := sg.Columns * sg.Rows
	if want == 0 {
		return fmt.Errorf("subgrid %q: zero grid dimensions", sg.Name)
	}
	if len(sg.Values) != want {
		return fmt.Errorf("subgrid %q: got %d values, want %d", sg.Name, len(sg.Values), want)
	}

	return nil
}
