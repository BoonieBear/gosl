// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// plt contains functions for plotting
package plt

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

// default directory and temporary file name for python commands
const TEMPORARY = "/tmp/pltgosl.py"

// buffer holding Python commands
var bufferPy bytes.Buffer

// buffer holding Python extra artists commands
var bufferEa bytes.Buffer

// init resets the buffers, in case the user doesn't do this
func init() {
	Reset()
}

// Reset resets drawing buffer (i.e. Python temporary file data)
func Reset() {
	bufferPy.Reset()
	bufferEa.Reset()
	io.Ff(&bufferPy, pythonHeader)
}

// PyCmds adds Python commands to be called when plotting
func PyCmds(text string) {
	io.Ff(&bufferPy, text)
}

// PyFile loads Python file and copy its contents to temporary buffer
func PyFile(filename string) (err error) {
	b, err := io.ReadFile(filename)
	if err != nil {
		return
	}
	io.Ff(&bufferPy, string(b))
	return
}

// DoubleYscale duplicates y-scale
func DoubleYscale(ylabelOrEmpty string) {
	io.Ff(&bufferPy, "plt.gca().twinx()\n")
	if ylabelOrEmpty != "" {
		io.Ff(&bufferPy, "plt.gca().set_ylabel('%s')\n", ylabelOrEmpty)
	}
}

// SetXlog sets x-scale to be log
func SetXlog() {
	io.Ff(&bufferPy, "plt.gca().set_xscale('log')\n")
}

// SetYlog sets y-scale to be log
func SetYlog() {
	io.Ff(&bufferPy, "plt.gca().set_yscale('log')\n")
}

// SetXnticks sets number of ticks along x
func SetXnticks(num int) {
	if num == 0 {
		io.Ff(&bufferPy, "plt.gca().get_xaxis().set_ticks([])\n")
	} else {
		io.Ff(&bufferPy, "plt.gca().get_xaxis().set_major_locator(tck.MaxNLocator(%d))\n", num)
	}
}

// SetYnticks sets number of ticks along y
func SetYnticks(num int) {
	if num == 0 {
		io.Ff(&bufferPy, "plt.gca().get_yaxis().set_ticks([])\n")
	} else {
		io.Ff(&bufferPy, "plt.gca().get_yaxis().set_major_locator(tck.MaxNLocator(%d))\n", num)
	}
}

// SetTicksX sets ticks along x
func SetTicksX(majorEvery, minorEvery float64, majorFmt string) {
	n := bufferPy.Len()
	io.Ff(&bufferPy, "majorLocator%d = tck.MultipleLocator(%g)\n", n, majorEvery)
	io.Ff(&bufferPy, "minorLocator%d = tck.MultipleLocator(%g)\n", n, minorEvery)
	io.Ff(&bufferPy, "majorFormatter%d = tck.FormatStrFormatter('%s')\n", n, majorFmt)
	io.Ff(&bufferPy, "plt.gca().xaxis.set_major_locator(majorLocator%d)\n", n)
	io.Ff(&bufferPy, "plt.gca().xaxis.set_minor_locator(minorLocator%d)\n", n)
	io.Ff(&bufferPy, "plt.gca().xaxis.set_major_formatter(majorFormatter%d)\n", n)
}

// SetTicksY sets ticks along y
func SetTicksY(majorEvery, minorEvery float64, majorFmt string) {
	n := bufferPy.Len()
	io.Ff(&bufferPy, "majorLocator%d = tck.MultipleLocator(%g)\n", n, majorEvery)
	io.Ff(&bufferPy, "minorLocator%d = tck.MultipleLocator(%g)\n", n, minorEvery)
	io.Ff(&bufferPy, "majorFormatter%d = tck.FormatStrFormatter('%s')\n", n, majorFmt)
	io.Ff(&bufferPy, "plt.gca().yaxis.set_major_locator(majorLocator%d)\n", n)
	io.Ff(&bufferPy, "plt.gca().yaxis.set_minor_locator(minorLocator%d)\n", n)
	io.Ff(&bufferPy, "plt.gca().yaxis.set_major_formatter(majorFormatter%d)\n", n)
}

// SetScientificX sets scientific notation for ticks along x-axis
func SetScientificX(minOrder, maxOrder int) {
	n := bufferPy.Len()
	io.Ff(&bufferPy, "fmt%d = plt.ScalarFormatter(useOffset=True)\n", n)
	io.Ff(&bufferPy, "fmt%d.set_powerlimits((%d,%d))\n", n, minOrder, maxOrder)
	io.Ff(&bufferPy, "plt.gca().xaxis.set_major_formatter(fmt%d)\n", n)
}

// SetScientificY sets scientific notation for ticks along y-axis
func SetScientificY(minOrder, maxOrder int) {
	n := bufferPy.Len()
	io.Ff(&bufferPy, "fmt%d = plt.ScalarFormatter(useOffset=True)\n", n)
	io.Ff(&bufferPy, "fmt%d.set_powerlimits((%d,%d))\n", n, minOrder, maxOrder)
	io.Ff(&bufferPy, "plt.gca().yaxis.set_major_formatter(fmt%d)\n", n)
}

// SetTicksNormal sets normal ticks
func SetTicksNormal() {
	io.Ff(&bufferPy, "plt.gca().ticklabel_format(useOffset=False)\n")
}

// ReplaceAxes substitutes axis frame (see Axes in gosl.py)
//   ex: xDel, yDel := 0.04, 0.04
func ReplaceAxes(xi, yi, xf, yf, xDel, yDel float64, xLab, yLab string, argsArrow, argsText *A) {
	io.Ff(&bufferPy, "plt.axis('off')\n")
	Arrow(xi, yi, xf, yi, argsArrow)
	Arrow(xi, yi, xi, yf, argsArrow)
	Text(xf, yi-xDel, xLab, argsText)
	Text(xi-yDel, yf, yLab, argsText)
}

// AxHline adds horizontal line to axis
func AxHline(y float64, args *A) {
	io.Ff(&bufferPy, "plt.axhline(%g", y)
	updateBufferAndClose(&bufferPy, args, false)
}

// AxVline adds vertical line to axis
func AxVline(x float64, args *A) {
	io.Ff(&bufferPy, "plt.axvline(%g", x)
	updateBufferAndClose(&bufferPy, args, false)
}

// HideBorders hides frame borders
func HideBorders(args *A) {
	hide := getHideList(args)
	if hide != "" {
		io.Ff(&bufferPy, "for spine in %s: plt.gca().spines[spine].set_visible(0)\n", hide)
	}
}

// Annotate adds annotation to plot
func Annotate(x, y float64, txt string, args *A) {
	io.Ff(&bufferPy, "plt.annotate(%q, xy=(%g,%g)", txt, x, y)
	updateBufferAndClose(&bufferPy, args, false)
}

// AnnotateXlabels sets text of xlabels
func AnnotateXlabels(x float64, txt string, args *A) {
	fsz := 7.0
	if args != nil {
		if args.Fsz > 0 {
			fsz = args.Fsz
		}
	}
	io.Ff(&bufferPy, "plt.annotate('%s', xy=(%g, -%g-3), xycoords=('data', 'axes points'), va='top', ha='center', size=%g", txt, x, fsz, fsz)
	updateBufferAndClose(&bufferPy, args, false)
}

// SupTitle sets subplot title
func SupTitle(txt string, args *A) {
	n := bufferPy.Len()
	io.Ff(&bufferPy, "st%d = plt.suptitle(%q", n, txt)
	updateBufferAndClose(&bufferPy, args, false)
	io.Ff(&bufferPy, "addToEA(st%d)\n", n)
}

// Title sets title
func Title(txt string, args *A) {
	io.Ff(&bufferPy, "plt.title(%q", txt)
	updateBufferAndClose(&bufferPy, args, false)
}

// Text adds text to plot
func Text(x, y float64, txt string, args *A) {
	io.Ff(&bufferPy, "plt.text(%g,%g,%q", x, y, txt)
	updateBufferAndClose(&bufferPy, args, false)
}

// Cross adds a vertical and horizontal lines @ (x0,y0) to plot (i.e. large cross)
func Cross(x0, y0 float64, args *A) {
	cl, ls, lw, z := "black", "dashed", 1.2, 0
	if args != nil {
		if args.C != "" {
			cl = args.C
		}
		if args.Lw > 0 {
			lw = args.Lw
		}
		if args.Ls != "" {
			ls = args.Ls
		}
		if args.Z > 0 {
			z = args.Z
		}
	}
	io.Ff(&bufferPy, "plt.axvline(%g, color='%s', linestyle='%s', linewidth=%g, zorder=%d)\n", x0, cl, ls, lw, z)
	io.Ff(&bufferPy, "plt.axhline(%g, color='%s', linestyle='%s', linewidth=%g, zorder=%d)\n", y0, cl, ls, lw, z)
}

// SplotGap sets gap between subplots
func SplotGap(w, h float64) {
	io.Ff(&bufferPy, "plt.subplots_adjust(wspace=%g, hspace=%g)\n", w, h)
}

// Subplot adds/sets a subplot
func Subplot(i, j, k int) {
	io.Ff(&bufferPy, "plt.subplot(%d,%d,%d)\n", i, j, k)
}

// Subplot adds/sets a subplot with given indices in I
func SubplotI(I []int) {
	if len(I) != 3 {
		return
	}
	io.Ff(&bufferPy, "plt.subplot(%d,%d,%d)\n", I[0], I[1], I[2])
}

// SetHspace sets horizontal space between subplots
func SetHspace(hspace float64) {
	io.Ff(&bufferPy, "plt.subplots_adjust(hspace=%g)\n", hspace)
}

// SetVspace sets vertical space between subplots
func SetVspace(vspace float64) {
	io.Ff(&bufferPy, "plt.subplots_adjust(vspace=%g)\n", vspace)
}

// Equal sets same scale for both axes
func Equal() {
	io.Ff(&bufferPy, "plt.axis('equal')\n")
}

// AxisOff hides axes
func AxisOff() {
	io.Ff(&bufferPy, "plt.axis('off')\n")
}

// SetAxis sets axes limits
func SetAxis(xmin, xmax, ymin, ymax float64) {
	io.Ff(&bufferPy, "plt.axis([%g, %g, %g, %g])\n", xmin, xmax, ymin, ymax)
}

// AxisXmin sets minimum x
func AxisXmin(xmin float64) {
	io.Ff(&bufferPy, "plt.axis([%g, plt.axis()[1], plt.axis()[2], plt.axis()[3]])\n", xmin)
}

// AxisXmax sets maximum x
func AxisXmax(xmax float64) {
	io.Ff(&bufferPy, "plt.axis([plt.axis()[0], %g, plt.axis()[2], plt.axis()[3]])\n", xmax)
}

// AxisYmin sets minimum y
func AxisYmin(ymin float64) {
	io.Ff(&bufferPy, "plt.axis([plt.axis()[0], plt.axis()[1], %g, plt.axis()[3]])\n", ymin)
}

// AxisYmax sets maximum y
func AxisYmax(ymax float64) {
	io.Ff(&bufferPy, "plt.axis([plt.axis()[0], plt.axis()[1], plt.axis()[2], %g])\n", ymax)
}

// AxisXrange sets x-range (i.e. limits)
func AxisXrange(xmin, xmax float64) {
	io.Ff(&bufferPy, "plt.axis([%g, %g, plt.axis()[2], plt.axis()[3]])\n", xmin, xmax)
}

// AxisYrange sets y-range (i.e. limits)
func AxisYrange(ymin, ymax float64) {
	io.Ff(&bufferPy, "plt.axis([plt.axis()[0], plt.axis()[1], %g, %g])\n", ymin, ymax)
}

// AxisRange sets x and y ranges (i.e. limits)
func AxisRange(xmin, xmax, ymin, ymax float64) {
	io.Ff(&bufferPy, "plt.axis([%g, %g, %g, %g])\n", xmin, xmax, ymin, ymax)
}

// AxisRange3d sets x, y, and z ranges (i.e. limits)
func AxisRange3d(xmin, xmax, ymin, ymax, zmin, zmax float64) {
	io.Ff(&bufferPy, "plt.gca().set_xlim3d(%g,%g)\ngca().set_ylim3d(%g,%g)\ngca().set_zlim3d(%g,%g)\n", xmin, xmax, ymin, ymax, zmin, zmax)
}

// AxisLims sets x and y limits
func AxisLims(lims []float64) {
	io.Ff(&bufferPy, "plt.axis([%g, %g, %g, %g])\n", lims[0], lims[1], lims[2], lims[3])
}

// Plot plots x-y series
func Plot(x, y []float64, args *A) (sx, sy string) {
	n := bufferPy.Len()
	sx = io.Sf("x%d", n)
	sy = io.Sf("y%d", n)
	gen2Arrays(&bufferPy, sx, sy, x, y)
	io.Ff(&bufferPy, "plt.plot(%s,%s", sx, sy)
	updateBufferAndClose(&bufferPy, args, false)
	return
}

// PlotOne plots one point @ (x,y)
func PlotOne(x, y float64, args *A) {
	io.Ff(&bufferPy, "plt.plot(%23.15e,%23.15e", x, y)
	updateBufferAndClose(&bufferPy, args, false)
}

// Hist draws histogram
func Hist(x [][]float64, labels []string, args *A) {
	n := bufferPy.Len()
	sx := io.Sf("x%d", n)
	sy := io.Sf("y%d", n)
	genList(&bufferPy, sx, x)
	genStrArray(&bufferPy, sy, labels)
	io.Ff(&bufferPy, "plt.hist(%s,label=%s", sx, sy)
	updateBufferAndClose(&bufferPy, args, true)
}

// ContourF draws filled contour and possibly with a contour of lines (if args.UnoLines=false)
func ContourF(x, y, z [][]float64, args *A) {
	n := bufferPy.Len()
	sx := io.Sf("x%d", n)
	sy := io.Sf("y%d", n)
	sz := io.Sf("z%d", n)
	genMat(&bufferPy, sx, x)
	genMat(&bufferPy, sy, y)
	genMat(&bufferPy, sz, z)
	a, colors, levels := argsContour(args)
	io.Ff(&bufferPy, "c%d = plt.contourf(%s,%s,%s%s%s)\n", n, sx, sy, sz, colors, levels)
	if !a.UnoLines {
		io.Ff(&bufferPy, "cc%d = plt.contour(%s,%s,%s,colors=['k']%s,linewidths=[%g])\n", n, sx, sy, sz, levels, a.Lw)
		if !a.UnoLabels {
			io.Ff(&bufferPy, "plt.clabel(cc%d,inline=%d,fontsize=%g)\n", n, pyBool(!a.UnoInline), a.Fsz)
		}
	}
	if !a.UnoCbar {
		io.Ff(&bufferPy, "cb%d = plt.colorbar(c%d, format='%s')\n", n, n, a.UnumFmt)
		if a.UcbarLbl != "" {
			io.Ff(&bufferPy, "cb%d.ax.set_ylabel('%s')\n", n, a.UcbarLbl)
		}
	}
	if a.UselectC != "" {
		io.Ff(&bufferPy, "ccc%d = plt.contour(%s,%s,%s,colors=['%s'],levels=[%g],linewidths=[%g],linestyles=['-'])\n", n, sx, sy, sz, a.UselectC, a.UselectV, a.UselectLw)
	}
}

// ContourL draws a contour with lines only
func ContourL(x, y, z [][]float64, args *A) {
	n := bufferPy.Len()
	sx := io.Sf("x%d", n)
	sy := io.Sf("y%d", n)
	sz := io.Sf("z%d", n)
	genMat(&bufferPy, sx, x)
	genMat(&bufferPy, sy, y)
	genMat(&bufferPy, sz, z)
	a, colors, levels := argsContour(args)
	io.Ff(&bufferPy, "c%d = plt.contour(%s,%s,%s%s%s)\n", n, sx, sy, sz, colors, levels)
	if !a.UnoLabels {
		io.Ff(&bufferPy, "plt.clabel(c%d,inline=%d,fontsize=%g)\n", n, pyBool(!a.UnoInline), a.Fsz)
	}
	if a.UselectC != "" {
		io.Ff(&bufferPy, "cc%d = plt.contour(%s,%s,%s,colors=['%s'],levels=[%g],linewidths=[%g],linestyles=['-'])\n", n, sx, sy, sz, a.UselectC, a.UselectV, a.UselectLw)
	}
}

// Quiver draws vector field
func Quiver(x, y, gx, gy [][]float64, args *A) {
	n := bufferPy.Len()
	sx := io.Sf("x%d", n)
	sy := io.Sf("y%d", n)
	sgx := io.Sf("gx%d", n)
	sgy := io.Sf("gy%d", n)
	genMat(&bufferPy, sx, x)
	genMat(&bufferPy, sy, y)
	genMat(&bufferPy, sgx, gx)
	genMat(&bufferPy, sgy, gy)
	io.Ff(&bufferPy, "plt.quiver(%s,%s,%s,%s", sx, sy, sgx, sgy)
	updateBufferAndClose(&bufferPy, args, false)
}

// Grid adds grid to plot
func Grid(args *A) {
	io.Ff(&bufferPy, "plt.grid(")
	updateBufferAndClose(&bufferPy, args, false)
}

// Legend adds legend to plot
func Legend(args *A) {
	loc, ncol, hlen, fsz, frame, out, outX := argsLeg(args)
	n := bufferPy.Len()
	io.Ff(&bufferPy, "h%d, l%d = plt.gca().get_legend_handles_labels()\n", n, n)
	io.Ff(&bufferPy, "if len(h%d) > 0 and len(l%d) > 0:\n", n, n)
	if out == 1 {
		io.Ff(&bufferPy, "    d%d = %s\n", n, outX)
		io.Ff(&bufferPy, "    l%d = plt.legend(bbox_to_anchor=d%d, ncol=%d, handlelength=%g, prop={'size':%g}, loc=3, mode='expand', borderaxespad=0.0, columnspacing=1, handletextpad=0.05)\n", n, n, ncol, hlen, fsz)
		io.Ff(&bufferPy, "    addToEA(l%d)\n", n)
	} else {
		io.Ff(&bufferPy, "    l%d = plt.legend(loc=%s, ncol=%d, handlelength=%g, prop={'size':%g})\n", n, loc, ncol, hlen, fsz)
		io.Ff(&bufferPy, "    addToEA(l%d)\n", n)
	}
	if frame == 0 {
		io.Ff(&bufferPy, "    l%d.get_frame().set_linewidth(0.0)\n", n)
	}
}

// Gll adds grid, labels, and legend to plot
func Gll(xl, yl string, args *A) {
	hide := getHideList(args)
	if hide != "" {
		io.Ff(&bufferPy, "for spine in %s: plt.gca().spines[spine].set_visible(False)\n", hide)
	}
	io.Ff(&bufferPy, "plt.grid(color='grey', zorder=-1000)\n")
	io.Ff(&bufferPy, "plt.xlabel(r'%s')\n", xl)
	io.Ff(&bufferPy, "plt.ylabel(r'%s')\n", yl)
	Legend(args)
}

// Clf clears current figure
func Clf() {
	io.Ff(&bufferPy, "plt.clf()\n")
}

// SetFontSizes sets font sizes
func SetFontSizes(args *A) {
	txt, lbl, leg, xtck, ytck := argsFsz(args)
	io.Ff(&bufferPy, "plt.rcParams.update({\n")
	io.Ff(&bufferPy, "    'font.size'       : %g,\n", txt)
	io.Ff(&bufferPy, "    'axes.labelsize'  : %g,\n", lbl)
	io.Ff(&bufferPy, "    'legend.fontsize' : %g,\n", leg)
	io.Ff(&bufferPy, "    'xtick.labelsize' : %g,\n", xtck)
	io.Ff(&bufferPy, "    'ytick.labelsize' : %g})\n", ytck)
}

// 3D /////////////////////////////////////////////////////////////////////////////////////////////

func get3daxes(doInit bool) (n int) {
	n = bufferPy.Len()
	if doInit {
		io.Ff(&bufferPy, "ax%d = plt.gcf().add_subplot(111, projection='3d')\n", n)
		io.Ff(&bufferPy, "ax%d.set_xlabel('x');ax%d.set_ylabel('y');ax%d.set_zlabel('z')\n", n, n, n)
	} else {
		io.Ff(&bufferPy, "ax%d = plt.gca()\n", n)
	}
	return
}

// Plot3dLine plots 3d line
func Plot3dLine(x, y, z []float64, doInit bool, args *A) {
	n := get3daxes(doInit)
	sx := io.Sf("x%d", n)
	sy := io.Sf("y%d", n)
	sz := io.Sf("z%d", n)
	genArray(&bufferPy, sx, x)
	genArray(&bufferPy, sy, y)
	genArray(&bufferPy, sz, z)
	io.Ff(&bufferPy, "p%d = ax%d.plot(%s,%s,%s", n, n, sx, sy, sz)
	updateBufferAndClose(&bufferPy, args, false)
}

// Plot3dPoints plots 3d points
func Plot3dPoints(x, y, z []float64, doInit bool, args *A) {
	n := get3daxes(doInit)
	sx := io.Sf("x%d", n)
	sy := io.Sf("y%d", n)
	sz := io.Sf("z%d", n)
	genArray(&bufferPy, sx, x)
	genArray(&bufferPy, sy, y)
	genArray(&bufferPy, sz, z)
	io.Ff(&bufferPy, "p%d = ax%d.scatter(%s,%s,%s", n, n, sx, sy, sz)
	updateBufferAndClose(&bufferPy, args, false)
}

// Wireframe draws wireframe
func Wireframe(x, y, z [][]float64, doInit bool, args *A) {
	n := get3daxes(doInit)
	sx := io.Sf("x%d", n)
	sy := io.Sf("y%d", n)
	sz := io.Sf("z%d", n)
	genMat(&bufferPy, sx, x)
	genMat(&bufferPy, sy, y)
	genMat(&bufferPy, sz, z)
	io.Ff(&bufferPy, "p%d = ax%d.plot_wireframe(%s,%s,%s", n, n, sx, sy, sz)
	updateBufferAndClose(&bufferPy, args, false)
}

// Surface draws surface
func Surface(x, y, z [][]float64, doInit bool, args *A) {
	n := get3daxes(doInit)
	sx := io.Sf("x%d", n)
	sy := io.Sf("y%d", n)
	sz := io.Sf("z%d", n)
	genMat(&bufferPy, sx, x)
	genMat(&bufferPy, sy, y)
	genMat(&bufferPy, sz, z)
	io.Ff(&bufferPy, "p%d = ax%d.plot_surface(%s,%s,%s", n, n, sx, sy, sz)
	updateBufferAndClose(&bufferPy, args, false)
}

// Camera sets camera in 3d graph
func Camera(elev, azim float64, args *A) {
	io.Ff(&bufferPy, "plt.gca().view_init(elev=%g, azim=%g", elev, azim)
	updateBufferAndClose(&bufferPy, args, false)
}

// AxDist sets distance in 3d graph
func AxDist(dist float64) {
	io.Ff(&bufferPy, "plt.gca().dist = %g\n", dist)
}

// functions to save figure ///////////////////////////////////////////////////////////////////////

// SetForPng prepares plot for saving PNG figure
func SetForPng(prop, widpt float64, dpi int, args *A) {
	txt, lbl, leg, xtck, ytck := argsFsz(args)
	Reset()
	width := widpt / 72.27 // width in inches
	height := width * prop // height in inches
	io.Ff(&bufferPy, "plt.rcdefaults()\n")
	io.Ff(&bufferPy, "plt.rcParams.update({\n")
	io.Ff(&bufferPy, "    'figure.figsize'  : [%d,%d],\n", int(width), int(height))
	io.Ff(&bufferPy, "    'savefig.dpi'     : %d,\n", dpi)
	io.Ff(&bufferPy, "    'font.size'       : %g,\n", txt)
	io.Ff(&bufferPy, "    'axes.labelsize'  : %g,\n", lbl)
	io.Ff(&bufferPy, "    'legend.fontsize' : %g,\n", leg)
	io.Ff(&bufferPy, "    'xtick.labelsize' : %g,\n", xtck)
	io.Ff(&bufferPy, "    'ytick.labelsize' : %g})\n", ytck)
}

// SetForEps prepares plot for saving EPS figure
func SetForEps(prop, widpt float64, args *A) {
	txt, lbl, leg, xtck, ytck := argsFsz(args)
	Reset()
	width := widpt / 72.27 // width in inches
	height := width * prop // height in inches
	io.Ff(&bufferPy, "plt.rcdefaults()\n")
	io.Ff(&bufferPy, "plt.rcParams.update({\n")
	io.Ff(&bufferPy, "    'figure.figsize'     : [%d,%d],\n", int(width), int(height))
	io.Ff(&bufferPy, "    'font.size'          : %g,\n", txt)
	io.Ff(&bufferPy, "    'axes.labelsize'     : %g,\n", lbl)
	io.Ff(&bufferPy, "    'legend.fontsize'    : %g,\n", leg)
	io.Ff(&bufferPy, "    'xtick.labelsize'    : %g,\n", xtck)
	io.Ff(&bufferPy, "    'ytick.labelsize'    : %g,\n", ytck)
	io.Ff(&bufferPy, "    'backend'            : 'ps',\n")
	io.Ff(&bufferPy, "    'text.usetex'        : True,\n")  // very IMPORTANT to avoid Type 3 fonts
	io.Ff(&bufferPy, "    'ps.useafm'          : True,\n")  // very IMPORTANT to avoid Type 3 fonts
	io.Ff(&bufferPy, "    'pdf.use14corefonts' : True})\n") // very IMPORTANT to avoid Type 3 fonts
}

// Save saves figure
func Save(fname string) error {
	io.Ff(&bufferPy, "plt.savefig(r'%s', bbox_inches='tight', bbox_extra_artists=EXTRA_ARTISTS)\n", fname)
	return run(fname)
}

// SaveD saves figure after creating a directory
func SaveD(dirout, fname string) (err error) {
	err = os.MkdirAll(dirout, 0777)
	if err != nil {
		return chk.Err("cannot create directory to save figure file:\n%v\n", err)
	}
	fn := filepath.Join(dirout, fname)
	io.Ff(&bufferPy, "plt.savefig(r'%s', bbox_inches='tight', bbox_extra_artists=EXTRA_ARTISTS)\n", fn)
	return run(fn)
}

// Show shows figure
func Show() error {
	io.Ff(&bufferPy, "plt.show()\n")
	return run("")
}

// generate arrays and matrices ///////////////////////////////////////////////////////////////////

// genMat generates matrix
func genMat(buf *bytes.Buffer, name string, a [][]float64) {
	io.Ff(buf, "%s=np.array([", name)
	for i, _ := range a {
		io.Ff(buf, "[")
		for j, _ := range a[i] {
			io.Ff(buf, "%g,", a[i][j])
		}
		io.Ff(buf, "],")
	}
	io.Ff(buf, "],dtype=float)\n")
}

// genList generates list
func genList(buf *bytes.Buffer, name string, a [][]float64) {
	io.Ff(buf, "%s=[", name)
	for i, _ := range a {
		io.Ff(buf, "[")
		for j, _ := range a[i] {
			io.Ff(buf, "%g,", a[i][j])
		}
		io.Ff(buf, "],")
	}
	io.Ff(buf, "]\n")
}

// genArray generates the NumPy text corresponding to an array of float point numbers
func genArray(buf *bytes.Buffer, name string, u []float64) {
	io.Ff(buf, "%s=np.array([", name)
	for i, _ := range u {
		io.Ff(buf, "%g,", u[i])
	}
	io.Ff(buf, "],dtype=float)\n")
}

// gen2Arrays generates the NumPy text corresponding to 2 arrays of float point numbers
func gen2Arrays(buf *bytes.Buffer, nameA, nameB string, a, b []float64) {
	genArray(buf, nameA, a)
	genArray(buf, nameB, b)
}

// genStrArray generates the NumPy text corresponding to an array of strings
func genStrArray(buf *bytes.Buffer, name string, u []string) {
	io.Ff(buf, "%s=[", name)
	for i, _ := range u {
		io.Ff(buf, "%q,", u[i])
	}
	io.Ff(buf, "]\n")
}

// call Python ////////////////////////////////////////////////////////////////////////////////////

// run calls Python to generate plot
func run(fn string) (err error) {

	// write file
	io.WriteFile(TEMPORARY, &bufferEa, &bufferPy)

	// set command
	cmd := exec.Command("python", TEMPORARY)
	var out, serr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &serr

	// call Python
	err = cmd.Run()
	if err != nil {
		return chk.Err("call to Python failed:\n%v\n", serr.String())
	}

	// show filename
	if fn != "" {
		io.Pf("file <%s> written\n", fn)
	}

	// show output
	io.Pf("%s", out.String())
	return
}

const pythonHeader = `### file generated by Gosl #################################################
import numpy as np
import matplotlib.pyplot as plt
import matplotlib.ticker as tck
import matplotlib.patches as pat
import matplotlib.path as pth
import matplotlib.patheffects as pff
import matplotlib.lines as lns
import mpl_toolkits.mplot3d as m3d
EXTRA_ARTISTS = []
def addToEA(obj):
    if obj!=None: EXTRA_ARTISTS.append(obj)
COLORMAPS = [plt.cm.bwr, plt.cm.RdBu, plt.cm.hsv, plt.cm.jet, plt.cm.terrain, plt.cm.pink, plt.cm.Greys]
def getCmap(idx): return COLORMAPS[idx %% len(COLORMAPS)]
`
