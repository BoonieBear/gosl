// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

// BinEntry holds data of an entry to bin
type BinEntry struct {
	Id int       // object Id
	X  []float64 // entry coordinate (read only)
}

// Bin defines one bin in Bins (holds entries for search)
type Bin struct {
	Idx     int         // index of bin
	Entries []*BinEntry // entries
}

// Bins defines bins to hold entries and speed up search
type Bins struct {
	Ndim int       // space dimension
	Xi   []float64 // [ndim] left/lower-most point
	Xf   []float64 // [ndim] right/upper-most point
	L    []float64 // [ndim] whole box lengths
	S    []float64 // size of bins
	N    []int     // [ndim] number of divisions
	All  []*Bin    // [nbins] all bins (there will be an extra bin row along each dimension)
	tmp  []int     // [ndim] temporary (auxiliary) slice
}

// xi   -- [ndim] initial positions
// xf   -- [ndim] final positions
// ndiv -- number of divisions for the maximum length
func (o *Bins) Init(xi, xf []float64, ndiv int) (err error) {

	// check for out-of-range values
	o.Ndim = len(xi)
	o.Xi = xi
	o.Xf = xf
	if len(xi) != len(xf) || len(xi) < 2 || len(xi) > 3 {
		return chk.Err("sizes of xi and l must be the same and equal to either 2 or 3")
	}

	// allocate length and number of division slices
	o.L = make([]float64, o.Ndim)
	o.S = make([]float64, o.Ndim)
	for k := 0; k < o.Ndim; k++ {
		o.L[k] = o.Xf[k] - o.Xi[k]
		o.S[k] = o.L[k] / float64(ndiv)
	}

	// number of divisions
	o.N = make([]int, o.Ndim)
	nbins := 1
	for k := 0; k < o.Ndim; k++ {
		o.N[k] = int(o.L[k]/o.S[k]) + 1
		nbins *= o.N[k]
	}

	// allocate slices
	o.All = make([]*Bin, nbins)
	o.tmp = make([]int, o.Ndim)
	return
}

// Append adds a new entry {x, id} to the bins structure
func (o *Bins) Append(x []float64, id int) (err error) {
	idx := o.CalcIdx(x)
	if idx < 0 {
		return chk.Err("point %v is out of range", x)
	}
	bin := o.FindBinByIndex(idx)
	if bin == nil {
		return chk.Err("bin index %v is out of range", idx)
	}
	xcopy := utl.DblCopy(x)
	entry := BinEntry{id, xcopy}
	bin.Entries = append(bin.Entries, &entry)
	return
}

// Clear clears all bins
func (o *Bins) Clear() {
	o.All = make([]*Bin, 0)
}

// Find returns the stored id of the entry whose coordinates are closest to x
// returns -1 if out of range or not found
func (o Bins) Find(x []float64) int {

	// index and check
	idx := o.CalcIdx(x)
	if idx < 0 {
		return -1 // out-of-range
	}

	// search for the closest point
	bin := o.FindBinByIndex(idx)
	dmin := math.MaxFloat64
	id_closest := -1
	var entry *BinEntry
	for _, entry = range bin.Entries {
		var d float64
		for k := 0; k < o.Ndim; k++ {
			d += math.Pow(x[k]-entry.X[k], 2)
		}
		if d < dmin {
			dmin = d
			id_closest = entry.Id
		}
	}
	return id_closest
}

// FindBinByIndex finds or allocate new bin corresponding to index idx
func (o Bins) FindBinByIndex(idx int) *Bin {

	// check
	if idx < 0 || idx >= len(o.All) {
		return nil
	}

	// allocate new bin if necessary
	if o.All[idx] == nil {
		o.All[idx] = new(Bin)
		o.All[idx].Idx = idx
	}
	return o.All[idx]
}

// CalcIdx calculates the bin index where the point x is
// returns -1 if out-of-range
func (o Bins) CalcIdx(x []float64) int {
	for k := 0; k < o.Ndim; k++ {
		if x[k] < o.Xi[k] || x[k] > o.Xf[k] {
			return -1
		}
		o.tmp[k] = int((x[k] - o.Xi[k]) / o.S[k])
	}
	idx := o.tmp[0] + o.tmp[1]*o.N[0]
	if o.Ndim > 2 {
		idx += o.tmp[2] * o.N[0] * o.N[1]
	}
	return idx
}

// FindAlongSegment gets the ids of entries that lie close to a segment
//  Note: the initial (xi) and final (xf) points on segment defined a bounding box of valid points
func (o Bins) FindAlongSegment(xi, xf []float64, tol float64) []int {

	// auxiliary variables
	var sbins []*Bin // selected bins
	lmax := utl.Max(o.S[0], o.S[1])
	if o.Ndim == 3 {
		lmax = utl.Max(lmax, o.S[2])
	}
	btol := 0.9 * lmax // tolerance for bins
	var p, pi, pf Point
	pi.X = xi[0]
	pf.X = xf[0]
	pi.Y = xi[1]
	pf.Y = xf[1]
	if o.Ndim == 3 {
		pi.Z = xi[2]
		pf.Z = xf[2]
	} else {
		xi = []float64{xi[0], xi[1], -1}
		xf = []float64{xf[0], xf[1], 1}
	}

	// loop along all bins
	var i, j, k int
	var x, y, z float64
	nxy := o.N[0] * o.N[1]
	for idx, bin := range o.All {

		// skip empty bins
		if bin == nil {
			continue
		}

		// coordinates of bin center
		i = idx % o.N[0] // indices representing bin
		j = (idx % nxy) / o.N[0]
		x = o.Xi[0] + float64(i)*o.S[0] // coordinates of bin corner
		y = o.Xi[1] + float64(j)*o.S[1]
		x += o.S[0] / 2.0
		y += o.S[1] / 2.0
		if o.Ndim == 3 {
			k = idx / nxy
			z = o.Xi[2] + float64(k)*o.S[2]
			z += o.S[2] / 2.0
		}

		// check if bin is near line
		p = Point{x, y, z}
		d := DistPointLine(&p, &pi, &pf, tol, false)
		if d <= btol {
			sbins = append(sbins, bin)
		}
	}

	// entries ids
	var ids []int

	// find closest points
	for _, bin := range sbins {
		for _, entry := range bin.Entries {
			x = entry.X[0]
			y = entry.X[1]
			if o.Ndim == 3 {
				z = entry.X[0]
			}
			p := Point{x, y, z}
			d := DistPointLine(&p, &pi, &pf, tol, false)
			if d <= tol {
				if IsPointIn(&p, xi, xf, tol) {
					ids = append(ids, entry.Id)
				}
			}
		}
	}
	return ids
}

func (o Bin) String() string {
	l := io.Sf("{\"idx\":%d, \"entries\":[", o.Idx)
	for i, entry := range o.Entries {
		if i > 0 {
			l += ", "
		}
		l += io.Sf("{\"id\":%d, \"x\":[%g,%g", entry.Id, entry.X[0], entry.X[1])
		if len(entry.X) > 2 {
			l += io.Sf(",%g", entry.X[2])
		}
		l += "]}"
	}
	l += "]}"
	return l
}

func (o Bins) String() string {
	l := "[\n"
	k := 0
	for _, bin := range o.All {
		if bin != nil {
			if k > 0 {
				l += ",\n"
			}
			l += io.Sf("  %v", bin)
			k += 1
		}
	}
	l += "\n]"
	return l
}

// Draw2d draws bins' grid
func (o *Bins) Draw2d(withtxt, withgrid, withentries, setup bool, selBins map[int]bool) {

	if withgrid {
		// horizontal lines
		x := []float64{o.Xi[0], o.Xi[0] + o.L[0] + o.S[0]}
		y := make([]float64, 2)
		for j := 0; j < o.N[1]+1; j++ {
			y[0] = o.Xi[1] + float64(j)*o.S[1]
			y[1] = y[0]
			plt.Plot(x, y, &plt.A{C: "#4f3677"})
		}

		// vertical lines
		y[0] = o.Xi[1]
		y[1] = o.Xi[1] + o.L[1] + o.S[1]
		for i := 0; i < o.N[0]+1; i++ {
			x[0] = o.Xi[0] + float64(i)*o.S[0]
			x[1] = x[0]
			plt.Plot(x, y, &plt.A{C: "#4f3677"})
		}
	}

	// selected bins
	nxy := o.N[0] * o.N[1]
	for idx, _ := range selBins {
		i := idx % o.N[0] // indices representing bin
		j := (idx % nxy) / o.N[0]
		x := o.Xi[0] + float64(i)*o.S[0] // coordinates of bin corner
		y := o.Xi[1] + float64(j)*o.S[1]
		plt.Polyline([][]float64{
			{x, y},
			{x + o.S[0], y},
			{x + o.S[0], y + o.S[1]},
			{x, y + o.S[1]},
		}, &plt.A{Fc: "#fbefdc", Ec: "#8e8371", Lw: 0.5, Closed: true})
	}

	// plot items
	if withentries {
		for _, bin := range o.All {
			if bin == nil {
				continue
			}
			for _, entry := range bin.Entries {
				plt.PlotOne(entry.X[0], entry.X[1], &plt.A{C: "r", M: "."})
			}
		}
	}

	// labels
	if withtxt {
		for j := 0; j < o.N[1]; j++ {
			for i := 0; i < o.N[0]; i++ {
				idx := i + j*o.N[0]
				x := o.Xi[0] + float64(i)*o.S[0] + 0.02*o.S[0]
				y := o.Xi[1] + float64(j)*o.S[1] + 0.02*o.S[1]
				plt.Text(x, y, io.Sf("%d", idx), nil)
			}
		}
	}

	// setup
	if setup {
		plt.Equal()
		plt.AxisRange(o.Xi[0]-0.1, o.Xf[0]+o.S[0]+0.1, o.Xi[1]-0.1, o.Xf[1]+o.S[1]+0.1)
	}
}
