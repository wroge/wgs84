package wgs84

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"regexp"
	"strings"
)

// ParseTaggedImageFileFormat reads a PROJ Geodetic TIFF grid (GTG) into the same runtime model as ParseGrid.
func ParseTaggedImageFileFormat(r io.Reader) (*Grid, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	p := &tifParser{data: data}
	if err := p.readHeader(); err != nil {
		return nil, err
	}

	var subgrids []SubGrid
	for ifdOff := p.firstIFD; ifdOff != 0; {
		sg, next, err := p.readIFD(ifdOff)
		if err != nil {
			return nil, err
		}
		if err := sg.validate(); err != nil {
			return nil, err
		}
		subgrids = append(subgrids, sg)
		ifdOff = next
	}

	if len(subgrids) == 0 {
		return nil, fmt.Errorf("tiff: no image directories")
	}

	return &Grid{SubGrids: subgrids}, nil
}

type tifParser struct {
	data     []byte
	order    binary.ByteOrder
	firstIFD uint32
}

func (p *tifParser) readHeader() error {
	if len(p.data) < 8 {
		return fmt.Errorf("tiff: file too short")
	}
	switch string(p.data[:2]) {
	case "II":
		p.order = binary.LittleEndian
	case "MM":
		p.order = binary.BigEndian
	default:
		return fmt.Errorf("tiff: invalid byte order")
	}
	if p.u16(2) != 42 {
		return fmt.Errorf("tiff: invalid magic")
	}
	p.firstIFD = p.u32(4)
	return nil
}

func (p *tifParser) readIFD(off uint32) (SubGrid, uint32, error) {
	if int(off)+2 > len(p.data) {
		return SubGrid{}, 0, fmt.Errorf("tiff: IFD out of range")
	}
	n := int(p.u16(int(off)))
	base := int(off) + 2
	need := base + n*12 + 4
	if need > len(p.data) {
		return SubGrid{}, 0, fmt.Errorf("tiff: IFD truncated")
	}

	tags := map[uint16]ifdEntry{}
	for i := range n {
		eo := base + i*12
		e := ifdEntry{
			tag:   p.u16(eo),
			typ:   p.u16(eo + 2),
			count: p.u32(eo + 4),
			val:   p.u32(eo + 8),
		}
		tags[e.tag] = e
	}
	next := p.u32(base + n*12)

	width := int(p.tagLong(tags, 256))
	height := int(p.tagLong(tags, 257))
	if width <= 0 || height <= 0 {
		return SubGrid{}, next, fmt.Errorf("tiff: invalid dimensions %dx%d", width, height)
	}

	meta := parseGDALMetadata(p.tagString(tags, 42112))
	if typ := meta["TYPE"]; typ != "" && typ != "HORIZONTAL_OFFSET" {
		return SubGrid{}, next, fmt.Errorf("tiff: unsupported TYPE %q", typ)
	}

	samples := int(p.tagLong(tags, 277))
	if samples < 2 {
		return SubGrid{}, next, fmt.Errorf("tiff: need at least 2 samples, got %d", samples)
	}

	bands, err := p.readBands(tags, width, height, samples)
	if err != nil {
		return SubGrid{}, next, err
	}

	latBand := bands[0]
	lonBand := bands[1]

	posLon := meta["positive_value"]
	if posLon == "" {
		posLon = "east"
	}

	// GTG tiepoint + pixel scale (degrees). Pixel (0,0) center at upper-left node.
	tie := p.tagDoubles(tags, 33922, 6)
	scale := p.tagDoubles(tags, 33550, 3)
	if len(tie) < 6 || len(scale) < 2 {
		return SubGrid{}, next, fmt.Errorf("tiff: missing ModelTiepointTag/ModelPixelScaleTag")
	}

	lonUL := tie[3] // degrees east
	latUL := tie[4] // degrees north
	dLon := scale[0]
	dLat := scale[1]
	if dLon == 0 || dLat == 0 {
		return SubGrid{}, next, fmt.Errorf("tiff: zero pixel scale")
	}

	// Outer grid edge aligned with first/last pixel centers (Grid node grid).
	lonWest := lonUL
	lonEast := lonUL + float64(width-1)*dLon
	latNorth := latUL
	latSouth := latUL - float64(height-1)*dLat

	sg := SubGrid{
		Name:    meta["grid_name"],
		Parent:  meta["parent_grid_name"],
		Columns: width,
		Rows:    height,
		SLat:    latSouth * 3600,
		NLat:    latNorth * 3600,
		ELong:   (-lonEast) * 3600, // eastern edge (min lam = -lon×3600)
		WLong:   (-lonWest) * 3600, // western edge (max lam)
		LatInc:  math.Abs(dLat) * 3600,
		LongInc: math.Abs(dLon) * 3600,
		Values:  make([][2]float32, width*height),
	}

	// Map GTG pixels → Grid node order (row 0 = south, column 0 = east).
	for row := range height {
		tiffRow := height - 1 - row
		for col := range width {
			iTiff := tiffRow*width + col
			iNtv2 := row*width + (width - 1 - col)

			latOff := latBand[iTiff]
			lonOff := lonBand[iTiff]

			switch posLon {
			case "east":
				lonOff = -lonOff // match Grid west-positive storage
			case "west":
				// already west-positive
			default:
				return SubGrid{}, next, fmt.Errorf("tiff: unknown positive_value %q", posLon)
			}

			sg.Values[iNtv2] = [2]float32{latOff, lonOff}
		}
	}

	return sg, next, nil
}

type ifdEntry struct {
	tag, typ uint16
	count    uint32
	val      uint32
}

func (p *tifParser) u16(off int) uint16 {
	return p.order.Uint16(p.data[off:])
}

func (p *tifParser) u32(off int) uint32 {
	return p.order.Uint32(p.data[off:])
}

func (p *tifParser) ifdValueSize(typ uint16) int {
	switch typ {
	case 3, 8: // SHORT, SBYTE
		return 2
	case 4, 9, 11: // LONG, SLONG, FLOAT
		return 4
	default:
		return 0
	}
}

func (p *tifParser) ifdInline(e ifdEntry) bool {
	sz := p.ifdValueSize(e.typ)
	return sz > 0 && int(e.count)*sz <= 4
}

func (p *tifParser) ifdInlineU16(e ifdEntry, i int) uint16 {
	var b [4]byte
	p.order.PutUint32(b[:], e.val)

	return p.order.Uint16(b[i*2 : i*2+2])
}

func (p *tifParser) tagLong(tags map[uint16]ifdEntry, tag uint16) uint32 {
	e, ok := tags[tag]
	if !ok {
		return 0
	}
	switch e.typ {
	case 3: // SHORT
		if e.count == 1 {
			return uint32(uint16(e.val))
		}
		if p.ifdInline(e) {
			return uint32(p.ifdInlineU16(e, 0))
		}
		return uint32(p.u16(int(e.val)))
	case 4: // LONG
		if e.count == 1 {
			return e.val
		}
		if p.ifdInline(e) {
			return e.val
		}
		return p.u32(int(e.val))
	default:
		return 0
	}
}

func (p *tifParser) tagString(tags map[uint16]ifdEntry, tag uint16) string {
	e, ok := tags[tag]
	if !ok || e.typ != 2 {
		return ""
	}
	off := int(e.val)
	end := off + int(e.count)
	if off < 0 || end > len(p.data) {
		return ""
	}
	return strings.TrimRight(string(p.data[off:end]), "\x00")
}

func (p *tifParser) tagDoubles(tags map[uint16]ifdEntry, tag uint16, want int) []float64 {
	e, ok := tags[tag]
	if !ok {
		return nil
	}
	off := int(e.val)

	out := make([]float64, 0, want)
	switch e.typ {
	case 5: // RATIONAL
		for i := 0; i < int(e.count) && i < want; i++ {
			o := off + i*8
			num := float64(p.order.Uint32(p.data[o:]))
			den := float64(p.order.Uint32(p.data[o+4:]))
			if den == 0 {
				return nil
			}
			out = append(out, num/den)
		}
	case 12: // DOUBLE
		for i := 0; i < int(e.count) && i < want; i++ {
			o := off + i*8
			out = append(out, math.Float64frombits(p.order.Uint64(p.data[o:])))
		}
	default:
		return nil
	}
	return out
}

func (p *tifParser) readBands(tags map[uint16]ifdEntry, width, height, samples int) ([][]float32, error) {
	planar := int(p.tagLong(tags, 284)) // 1=contig, 2=separate
	if planar == 0 {
		planar = 1
	}
	comp := int(p.tagLong(tags, 259))
	predictor := int(p.tagLong(tags, 317))
	bps := int(p.tagLong(tags, 258))
	if bps == 0 {
		bps = 32
	}
	sfmt := int(p.tagLong(tags, 339)) // SampleFormat (338 is ExtraSamples)
	if sfmt == 0 && bps == 32 {
		sfmt = 3 // IEEE float
	}
	if bps != 32 || sfmt != 3 {
		return nil, fmt.Errorf("tiff: only Float32 supported (bps=%d fmt=%d)", bps, sfmt)
	}

	n := width * height
	bands := make([][]float32, samples)
	for i := range bands {
		bands[i] = make([]float32, n)
	}

	if p.isTiled(tags) {
		return p.readTiledBands(tags, width, height, samples, comp, predictor, planar, bands)
	}

	return p.readStripBands(tags, width, height, samples, comp, predictor, planar, bands)
}

func (p *tifParser) isTiled(tags map[uint16]ifdEntry) bool {
	if len(p.tagLongArray(tags, 273)) > 0 {
		return false
	}
	if len(p.tagLongArray(tags, 324)) > 1 {
		return true
	}
	if len(p.tagLongArray(tags, 322)) > 1 {
		return true
	}
	return false
}

func (p *tifParser) tileLayout(tags map[uint16]ifdEntry) (offsets, counts []uint32, tileW, tileH int) {
	// Standard TIFF: 322/323 offsets, 324/325 dimensions.
	offsets = p.tagLongArray(tags, 322)
	counts = p.tagLongArray(tags, 323)
	tileW = int(p.tagLong(tags, 324))
	tileH = int(p.tagLong(tags, 325))

	// PROJ GTG: 322/323 hold tile size, 324/325 hold offset tables.
	if len(offsets) <= 1 && len(p.tagLongArray(tags, 324)) > 1 {
		offsets = p.tagLongArray(tags, 324)
		counts = p.tagLongArray(tags, 325)
		tileW = int(p.tagLong(tags, 322))
		tileH = int(p.tagLong(tags, 323))
	}

	return offsets, counts, tileW, tileH
}

func (p *tifParser) readTiledBands(
	tags map[uint16]ifdEntry,
	width, height, samples, comp, predictor, planar int,
	bands [][]float32,
) ([][]float32, error) {
	offsets, counts, tileW, tileH := p.tileLayout(tags)
	if len(offsets) == 0 || len(offsets) != len(counts) || tileW <= 0 || tileH <= 0 {
		return nil, fmt.Errorf("tiff: invalid tile tables")
	}

	tilesX := (width + tileW - 1) / tileW
	tilesY := (height + tileH - 1) / tileH
	tilesPerSample := tilesX * tilesY

	if planar == 2 {
		if len(offsets)%samples != 0 || len(offsets)/samples != tilesPerSample {
			return nil, fmt.Errorf("tiff: tile count %d, want %d×%d", len(offsets), samples, tilesPerSample)
		}
	} else if len(offsets) != tilesPerSample {
		return nil, fmt.Errorf("tiff: tile count %d, want %d", len(offsets), tilesPerSample)
	}

	readSampleTiles := func(sample int) error {
		for ty := range tilesY {
			for tx := range tilesX {
				tw := min(tileW, width-tx*tileW)
				th := min(tileH, height-ty*tileH)
				ti := ty*tilesX + tx

				idx := ti
				if planar == 2 {
					idx = sample*tilesPerSample + ti
				}

				dec, err := p.readChunk(offsets[idx], counts[idx], comp, tileW*tileH*4, predictor, tileW)
				if err != nil {
					return fmt.Errorf("tile (%d,%d) sample %d: %w", tx, ty, sample, err)
				}

				for row := range th {
					for col := range tw {
						src := (row*tileW + col) * 4
						dst := (ty*tileH+row)*width + (tx*tileW + col)
						bands[sample][dst] = math.Float32frombits(p.order.Uint32(dec[src:]))
					}
				}
			}
		}
		return nil
	}

	if planar == 2 {
		for s := range samples {
			if err := readSampleTiles(s); err != nil {
				return nil, err
			}
		}
		return bands, nil
	}

	if err := readSampleTiles(0); err != nil {
		return nil, err
	}
	for i := range width * height {
		for s := 1; s < samples; s++ {
			// contig tiled: not expected for GTG; leave zero
			_ = bands[s][i]
		}
	}
	return bands, nil
}

func (p *tifParser) readStripBands(
	tags map[uint16]ifdEntry,
	width, height, samples, comp, predictor, planar int,
	bands [][]float32,
) ([][]float32, error) {
	offsets := p.tagLongArray(tags, 273)
	counts := p.tagLongArray(tags, 279)
	if len(offsets) == 0 || len(offsets) != len(counts) {
		return nil, fmt.Errorf("tiff: invalid strip tables")
	}

	n := width * height
	rowsPerStrip := int(p.tagLong(tags, 278))
	if rowsPerStrip == 0 {
		rowsPerStrip = height
	}
	stripsPerSample := len(offsets)
	if planar == 2 {
		if len(offsets)%samples != 0 {
			return nil, fmt.Errorf("tiff: strip count %d not divisible by samples %d", len(offsets), samples)
		}
		stripsPerSample = len(offsets) / samples
	}

	readSample := func(sample int) error {
		stripBase := 0
		if planar == 2 {
			stripBase = sample * stripsPerSample
		}

		row := 0
		for si := range stripsPerSample {
			idx := stripBase + si
			stripRows := rowsPerStrip
			if row+stripRows > height {
				stripRows = height - row
			}
			want := width * stripRows * 4

			dec, err := p.readChunk(offsets[idx], counts[idx], comp, want, predictor, width)
			if err != nil {
				return fmt.Errorf("strip %d sample %d: %w", si, sample, err)
			}

			for r := range stripRows {
				for col := range width {
					src := (r*width + col) * 4
					dst := (row+r)*width + col
					bands[sample][dst] = math.Float32frombits(p.order.Uint32(dec[src:]))
				}
			}
			row += stripRows
		}
		return nil
	}

	if planar == 2 {
		for s := range samples {
			if err := readSample(s); err != nil {
				return nil, fmt.Errorf("sample %d: %w", s, err)
			}
		}
		return bands, nil
	}

	dec, err := p.decompressStrips(offsets, counts, comp, n*samples*4, predictor, width*samples)
	if err != nil {
		return nil, err
	}
	for i := range n {
		for s := range samples {
			off := (i*samples + s) * 4
			bands[s][i] = math.Float32frombits(p.order.Uint32(dec[off:]))
		}
	}
	return bands, nil
}

func (p *tifParser) readChunk(offset, count uint32, comp, want, predictor, rowWidth int) ([]byte, error) {
	off := int(offset)
	n := int(count)
	if off < 0 || off+n > len(p.data) {
		return nil, fmt.Errorf("tiff: chunk out of range")
	}

	dec, err := decompressTIFF(p.data[off:off+n], comp)
	if err != nil {
		return nil, err
	}
	if predictor == 3 {
		rowBytes := rowWidth * 4
		for rowOff := 0; rowOff+rowBytes <= len(dec); rowOff += rowBytes {
			undoFloatPredictorRow(dec[rowOff : rowOff+rowBytes])
		}
	}
	if want >= 0 && len(dec) != want {
		return nil, fmt.Errorf("tiff: chunk size %d, want %d", len(dec), want)
	}
	return dec, nil
}

func (p *tifParser) decompressStrips(offsets, counts []uint32, comp, want, predictor, rowWidth int) ([]byte, error) {
	var out []byte
	for i := range offsets {
		dec, err := p.readChunk(offsets[i], counts[i], comp, -1, predictor, rowWidth)
		if err != nil {
			return nil, fmt.Errorf("strip %d: %w", i, err)
		}
		out = append(out, dec...)
	}
	if want >= 0 && len(out) != want {
		return nil, fmt.Errorf("tiff: strip total size %d, want %d", len(out), want)
	}
	return out, nil
}

// undoFloatPredictorRow reverses libtiff predictor 3 for one scanline/tile row.
func undoFloatPredictorRow(row []byte) {
	if len(row) < 8 {
		return
	}
	const bps = 4
	wc := len(row) / bps

	for i := 1; i < len(row); i++ {
		row[i] = row[i] + row[i-1]
	}

	tmp := make([]byte, len(row))
	copy(tmp, row)

	for i := range wc {
		for b := range bps {
			row[i*bps+b] = tmp[(bps-b-1)*wc+i]
		}
	}
}

func (p *tifParser) tagLongArray(tags map[uint16]ifdEntry, tag uint16) []uint32 {
	e, ok := tags[tag]
	if !ok {
		return nil
	}

	if e.count == 1 {
		switch e.typ {
		case 3:
			return []uint32{uint32(uint16(e.val))}
		case 4:
			return []uint32{e.val}
		}
	}

	if p.ifdInline(e) {
		out := make([]uint32, e.count)
		for i := range e.count {
			switch e.typ {
			case 3:
				out[i] = uint32(p.ifdInlineU16(e, int(i)))
			case 4:
				out[i] = e.val
			default:
				return nil
			}
		}

		return out
	}
	off := int(e.val)
	out := make([]uint32, e.count)
	for i := range e.count {
		switch e.typ {
		case 4:
			out[i] = p.u32(off + int(i)*4)
		case 3:
			out[i] = uint32(p.u16(off + int(i)*2))
		default:
			return nil
		}
	}
	return out
}

func decompressTIFF(raw []byte, compression int) ([]byte, error) {
	switch compression {
	case 1: // none
		return raw, nil
	case 8, 32946: // deflate / adobe deflate
		zr, err := zlib.NewReader(bytes.NewReader(raw))
		if err != nil {
			return nil, fmt.Errorf("tiff: deflate: %w", err)
		}

		defer zr.Close() //nolint:errcheck

		return io.ReadAll(zr)
	default:
		return nil, fmt.Errorf("tiff: unsupported compression %d", compression)
	}
}

var gdalItemRE = regexp.MustCompile(`<Item\s+name="([^"]+)"[^>]*>([^<]*)</Item>`)

func parseGDALMetadata(s string) map[string]string {
	out := map[string]string{}
	for _, m := range gdalItemRE.FindAllStringSubmatch(s, -1) {
		out[m[1]] = strings.TrimSpace(m[2])
	}
	return out
}
