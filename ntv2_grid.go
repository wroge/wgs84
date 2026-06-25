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
		return nil, fmt.Errorf("%w: %s", ErrGridNotFound, p)
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

func (g *Grid) ToWGS84(lon, lat float64) (float64, float64, error) {
	dlon, dlat, err := g.Shift(lon, lat)
	if err != nil {
		return 0, 0, err
	}

	return lon + dlon, lat + dlat, nil
}

func (g *Grid) FromWGS84(lon, lat float64) (float64, float64, error) {
	qlon, qlat := lon, lat

	for range 10 {
		dlon, dlat, err := g.Shift(qlon, qlat)
		if err != nil {
			return 0, 0, err
		}

		newLon := lon - dlon
		newLat := lat - dlat

		if math.Abs(newLon-qlon) < 1e-12 && math.Abs(newLat-qlat) < 1e-12 {
			break
		}

		qlon, qlat = newLon, newLat
	}

	return qlon, qlat, nil
}

func (g *Grid) Shift(lon, lat float64) (dlon, dlat float64, err error) {
	idx := g.selectSubgrid(lon, lat)
	if idx < 0 {
		return 0, 0, ErrOutOfBounds
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

const relToleranceHGridShift = 1e-5

func (sg SubGrid) gridEpsilon() float64 {
	return (sg.LongInc + sg.LatInc) * relToleranceHGridShift
}

func (sg SubGrid) contains(phi, lam float64) bool {
	eps := sg.gridEpsilon()
	return phi >= sg.SLat-eps && phi <= sg.NLat+eps &&
		lam >= sg.ELong-eps && lam <= sg.WLong+eps
}

func interpolateIndex(f float64, size int) (idx int, frac float64, ok bool) {
	if math.IsNaN(f) {
		return 0, 0, true
	}

	idx = int(math.Round(math.Floor(f)))
	frac = f - float64(idx)

	if idx < 0 {
		if idx == -1 && frac > 1-10*relToleranceHGridShift {
			idx++
			frac = 0
		} else {
			return 0, 0, false
		}
	} else if idx+1 >= size {
		if idx+1 == size && frac < 10*relToleranceHGridShift {
			idx--
			frac = 1
		} else {
			return 0, 0, false
		}
	}

	return idx, frac, true
}

func (sg SubGrid) shift(lon, lat float64) (dlon, dlat float64, err error) {
	if sg.Columns < 2 || sg.Rows < 2 || len(sg.Values) == 0 {
		return 0, 0, ErrOutOfBounds
	}

	phi := lat * 3600
	lam := -lon * 3600
	if !sg.contains(phi, lam) {
		return 0, 0, ErrOutOfBounds
	}

	fcol := (lam - sg.ELong) / sg.LongInc
	frow := (phi - sg.SLat) / sg.LatInc

	ppr := sg.Columns

	col, dx, ok := interpolateIndex(fcol, ppr)
	if !ok {
		return 0, 0, ErrOutOfBounds
	}

	row, dy, ok := interpolateIndex(frow, sg.Rows)
	if !ok {
		return 0, 0, ErrOutOfBounds
	}

	se := row*ppr + col
	sw := se + 1
	ne := se + ppr
	nw := ne + 1

	sse := sg.Values[se]
	ssw := sg.Values[sw]
	sne := sg.Values[ne]
	snw := sg.Values[nw]

	latsv := (1-dx)*(1-dy)*float64(sse[0]) + dx*(1-dy)*float64(ssw[0]) +
		(1-dx)*dy*float64(sne[0]) + dx*dy*float64(snw[0])
	lonsv := (1-dx)*(1-dy)*float64(sse[1]) + dx*(1-dy)*float64(ssw[1]) +
		(1-dx)*dy*float64(sne[1]) + dx*dy*float64(snw[1])

	return -lonsv / 3600, latsv / 3600, nil
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
