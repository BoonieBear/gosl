// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plt

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/utl"
)

func Test_args01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("args01")

	var a A

	// plot and basic options
	a.C = "red"
	a.M = "o"
	a.Ls = "--"
	a.Lw = 1.2
	a.Ms = -1
	a.L = "gosl"
	a.Me = 2
	a.Z = 123
	a.Mec = "blue"
	a.Mew = 0.3
	a.Void = true
	a.NoClip = true

	// shapes
	a.Ha = "center"
	a.Va = "center"
	a.Fc = "magenta"
	a.Ec = "yellow"

	// text and extra arguments
	a.Fsz = 7

	l := a.String(false)
	chk.String(tst, l, "color='red',marker='o',ls='--',lw=1.2,label='gosl',markevery=2,zorder=123,markeredgecolor='blue',mew=0.3,markerfacecolor='none',clip_on=0,facecolor='magenta',edgecolor='yellow',ha='center',va='center',fontsize=7")
}

func Test_args02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("args02")

	a := &A{
		Colors:   []string{"red", "tan", "lime"},
		Htype:    "bar",
		Hstacked: true,
		Hvoid:    true,
		Hnbins:   10,
		Hnormed:  true,
	}

	l := a.String(true)
	chk.String(tst, l, "color=['red','tan','lime'],histtype='bar',stacked=1,fill=0,bins=10,normed=1")
}

func Test_plot01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("plot01")

	if chk.Verbose {

		x := utl.LinSpace(0.0, 1.0, 11)
		y := make([]float64, len(x))
		for i := 0; i < len(x); i++ {
			y[i] = x[i] * x[i]
		}

		SetForPng(0.75, 600, 100, nil)
		//SetForEps(0.75, 600, nil)
		SetFontSizes(&A{Fsz: 20, FszLbl: 20, FszXtck: 10, FszYtck: 10})
		Plot(x, y, &A{L: "first", C: "r", M: "o", Ls: "-", Lw: 2, NoClip: true})
		Plot(y, x, &A{L: "second", C: "b", M: ".", Ls: ":", Lw: 40})
		Text(0.2, 0.8, "HERE", &A{Fsz: 20, Ha: "center", Va: "center"})
		SetTicksX(0.1, 0.01, "%.3f")
		SetTicksY(0.2, 0.1, "%.2f")
		HideBorders(&A{HideR: true, HideT: true})
		Gll(`$\varepsilon$`, `$\sigma$`, &A{
			LegOut:  true,
			LegNcol: 3,
			FszLeg:  14,
			HideR:   true,
		})

		err := SaveD("/tmp/gosl", "t_plot01.png")
		//err := SaveD("/tmp/gosl", "t_plot01.eps")
		if err != nil {
			tst.Errorf("%v", err)
		}
	}
}

func Test_plot02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("plot02")

	if chk.Verbose {

		Reset()
		ReplaceAxes(0, 0, 1, 1, 0.04, 0.04, "the x", "the y", &A{}, &A{})
		Arrow(0, 0, 1, 1, &A{})
		AxHline(0, &A{C: "red"})
		AxVline(0, &A{C: "blue"})
		Annotate(0, 0, "TEST", &A{C: "green"})
		AnnotateXlabels(0, "HERE", &A{Fsz: 10})
		SupTitle("suptitle goes here", &A{C: "red"})
		Title("title goes here", &A{C: "yellow"})
		Text(0, 0, "TEXT", &A{C: "blue"})
		Cross(0.5, 0.5, nil)
		PlotOne(0, 0, &A{M: "*"})

		err := SaveD("/tmp/gosl", "t_plot02.png")
		if err != nil {
			tst.Errorf("%v", err)
		}
	}
}

func Test_plot03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("plot03")

	if chk.Verbose {

		// grid size
		xmin, xmax, N := -math.Pi/2.0+0.1, math.Pi/2.0-0.1, 21

		// mesh grid
		X, Y, F := utl.MeshGrid2dF(xmin, xmax, xmin, xmax, N, N, func(x, y float64) (z float64) {
			m := math.Pow(math.Cos(x), 2.0) + math.Pow(math.Cos(y), 2.0)
			z = -math.Pow(m, 2.0)
			return
		})

		// configuration
		a := &A{
			UnumFmt:  "%.1f",
			Lw:       1.5,
			UcbarLbl: "NICE",
			UselectC: "yellow",
			UselectV: -2.5,
		}

		Reset()
		Equal()
		ContourF(X, Y, F, a)

		err := SaveD("/tmp/gosl", "t_plot03.png")
		if err != nil {
			tst.Errorf("%v", err)
		}
	}
}

func Test_plot04(tst *testing.T) {

	//verbose()
	chk.PrintTitle("plot04")

	if chk.Verbose {

		// grid size
		xmin, xmax, N := -math.Pi/2.0+0.1, math.Pi/2.0-0.1, 21

		// mesh grid
		X, Y, F, U, V := utl.MeshGrid2dFG(xmin, xmax, xmin, xmax, N, N, func(x, y float64) (z, u, v float64) {
			m := math.Pow(math.Cos(x), 2.0) + math.Pow(math.Cos(y), 2.0)
			z = -math.Pow(m, 2.0)
			u = 4.0 * math.Cos(x) * math.Sin(x) * m
			v = 4.0 * math.Cos(y) * math.Sin(y) * m
			return
		})

		// configuration
		a := &A{
			Colors:    []string{"cyan", "blue", "yellow", "green"},
			Ulevels:   []float64{-4, -3, -2, -1, 0},
			UnumFmt:   "%.1f",
			UnoLines:  true,
			UnoLabels: true,
			UnoInline: true,
			UnoCbar:   true,
			Lw:        1.5,
			UselectC:  "white",
			UselectV:  -2.5,
		}

		b := &A{
			UcmapIdx: 4,
			UselectC: "black",
			UselectV: -2.5,
		}

		Reset()
		Equal()
		ContourF(X, Y, F, a)
		ContourL(X, Y, F, b)
		Quiver(X, Y, U, V, nil)

		err := SaveD("/tmp/gosl", "t_plot04.png")
		if err != nil {
			tst.Errorf("%v", err)
		}
	}
}

func Test_plot05(tst *testing.T) {

	//verbose()
	chk.PrintTitle("plot05")

	if chk.Verbose {

		X := [][]float64{
			{1, 1, 1, 2, 2, 2, 2, 2, 3, 3, 4, 5, 6}, // first series
			{-1, -1, 0, 1, 2, 3},                    // second series
			{5, 6, 7, 8},                            // third series
		}

		L := []string{
			"first",
			"second",
			"third",
		}

		a := &A{
			Colors:   []string{"red", "tan", "lime"},
			Ec:       "black",
			Lw:       0.5,
			Htype:    "bar",
			Hstacked: true,
			//Hvoid:    true,
			//Hnbins:   10,
			//Hnormed: true,
		}

		Reset()
		Hist(X, L, a)
		Gll("series", "count", nil)

		err := SaveD("/tmp/gosl", "t_plot05.png")
		if err != nil {
			tst.Errorf("%v", err)
		}
	}
}

func Test_plot06(tst *testing.T) {

	//verbose()
	chk.PrintTitle("plot06")

	if chk.Verbose {

		x := []float64{0, 1, 1, 1}
		y := []float64{0, 0, 1, 1}
		z := []float64{0, 0, 0, 1}

		np := 3

		X, Y, Z := utl.MeshGrid2dF(0, 1, 0, 1, np, np, func(x, y float64) float64 {
			return x + y
		})

		U, V, W := utl.MeshGrid2dF(0, 1, 0, 1, np, np, func(u, v float64) float64 {
			return u*u + v*v
		})

		Reset()
		Plot3dLine(x, y, z, true, nil)
		Plot3dPoints(x, y, z, false, nil)
		Wireframe(X, Y, Z, false, nil)
		Surface(U, V, W, false, nil)
		//Camera(elev, azim float64, args *A)
		//SetFontSize(args *A) {
		//Circle(xc, yc, r float64, args *A)

		err := SaveD("/tmp/gosl", "t_plot06.png")
		if err != nil {
			tst.Errorf("%v", err)
		}
	}
}
