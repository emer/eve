// Copyright (c) 2018, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/emer/epe/epe"
	"github.com/emer/epe/epev"
	"github.com/goki/gi/gi"
	"github.com/goki/gi/gi3d"
	"github.com/goki/gi/gimain"
	"github.com/goki/gi/giv"
	"github.com/goki/gi/mat32"
	"github.com/goki/gi/units"
	"github.com/goki/ki/ki"
)

func main() {
	gimain.Main(func() {
		mainrun()
	})
}

// MakeWorld constructs a new virtual physics world
func MakeWorld() *epe.Group {
	world := &epe.Group{}
	world.InitName(world, "RoomWorld")

	rm1 := epe.AddNewGroup(world, "room1")

	thick := float32(.05)
	width := float32(10)
	depth := float32(10)
	height := float32(10)
	bwall := epe.AddNewBox(rm1, "back-wall", mat32.Vec3{0, height / 2, -depth / 2}, mat32.Vec3{width, height, thick})
	bwall.Mat.Color = "tan"
	lwall := epe.AddNewBox(rm1, "left-wall", mat32.Vec3{-width / 2, height / 2, 0}, mat32.Vec3{thick, height, depth})
	lwall.Mat.Color = "red"
	rwall := epe.AddNewBox(rm1, "right-wall", mat32.Vec3{width / 2, height / 2, 0}, mat32.Vec3{thick, height, depth})
	rwall.Mat.Color = "green"

	world.InitWorld()
	return world
}

func mainrun() {
	width := 1024
	height := 768

	// turn these on to see a traces of various stages of processing..
	// ki.SignalTrace = true
	// gi.WinEventTrace = true
	// gi3d.Update3DTrace = true
	// gi.Update2DTrace = true

	rec := ki.Node{}          // receiver for events
	rec.InitName(&rec, "rec") // this is essential for root objects not owned by other Ki tree nodes

	gi.SetAppName("gi3d")
	gi.SetAppAbout(`This is a demo of the Emergent Physics Engine.  See <a href="https://github.com/emer/epe">epe on GitHub</a>.
<p>The <a href="https://github.com/emer/epe/blob/master/examples/virtroom/README.md">README</a> page for this example app has further info.</p>`)

	win := gi.NewWindow2D("epe-demo", "Emergent Physics Engine", width, height, true) // true = pixel sizes

	vp := win.WinViewport2D()
	updt := vp.UpdateStart()

	mfr := win.SetMainFrame()
	mfr.SetProp("spacing", units.NewEx(1))

	trow := gi.AddNewLayout(mfr, "trow", gi.LayoutHoriz)
	trow.SetStretchMaxWidth()

	title := gi.AddNewLabel(trow, "title", `This is a demonstration of the
<a href="https://github.com/emer/epe">epe</a> <i>3D</i> Framework<br>
See <a href="https://github.com/emer/epe/blob/master/examples/virtroomd/README.md">README</a> for detailed info and things to try.`)
	title.SetProp("white-space", gi.WhiteSpaceNormal) // wrap
	title.SetProp("text-align", gi.AlignCenter)       // note: this also sets horizontal-align, which controls the "box" that the text is rendered in..
	title.SetProp("vertical-align", gi.AlignCenter)
	title.SetProp("font-size", "x-large")
	title.SetProp("line-height", 1.5)
	title.SetStretchMaxWidth()
	title.SetStretchMaxHeight()

	//////////////////////////////////////////
	//    world

	world := MakeWorld()

	//////////////////////////////////////////
	//    Splitter

	split := gi.AddNewSplitView(mfr, "split")
	split.Dim = gi.X

	tvfr := gi.AddNewFrame(split, "tvfr", gi.LayoutHoriz)
	svfr := gi.AddNewFrame(split, "svfr", gi.LayoutHoriz)
	scfr := gi.AddNewFrame(split, "scfr", gi.LayoutHoriz)
	split.SetSplits(.2, .2, .6)

	tv := giv.AddNewTreeView(tvfr, "tv")
	tv.SetRootNode(world)

	sv := giv.AddNewStructView(svfr, "sv")
	sv.SetStretchMaxWidth()
	sv.SetStretchMaxHeight()
	sv.SetStruct(world)

	tv.TreeViewSig.Connect(sv.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		if data == nil {
			return
		}
		// tvr, _ := send.Embed(giv.KiT_TreeView).(*gi.TreeView) // root is sender
		tvn, _ := data.(ki.Ki).Embed(giv.KiT_TreeView).(*giv.TreeView)
		svr, _ := recv.Embed(giv.KiT_StructView).(*giv.StructView)
		if sig == int64(giv.TreeViewSelected) {
			svr.SetStruct(tvn.SrcNode)
		}
	})

	//////////////////////////////////////////
	//    Scene

	sc := gi3d.AddNewScene(scfr, "scene")
	sc.SetStretchMaxWidth()
	sc.SetStretchMaxHeight()

	// first, add lights, set camera
	sc.BgColor.SetUInt8(230, 230, 255, 255) // sky blue-ish
	gi3d.AddNewAmbientLight(sc, "ambient", 0.3, gi3d.DirectSun)

	dir := gi3d.AddNewDirLight(sc, "dir", 1, gi3d.DirectSun)
	dir.Pos.Set(0, 2, 1) // default: 0,1,1 = above and behind us (we are at 0,0,X)

	wgp := gi3d.AddNewGroup(sc, sc, "world")

	wview := epev.NewView(world, wgp)
	wview.Sync(sc)

	// grtx := gi3d.AddNewTextureFile(sc, "ground", "ground.png")
	// wdtx := gi3d.AddNewTextureFile(sc, "wood", "wood.png")

	// floorp := gi3d.AddNewPlane(sc, "floor-plane", 100, 100)
	// floor := gi3d.AddNewObject(sc, sc, "floor", floorp.Name())
	// floor.Pose.Pos.Set(0, -5, 0)
	// // floor.Mat.Color.SetName("tan")
	// // floor.Mat.Emissive.SetName("brown")
	// floor.Mat.Bright = 2 // .5 for wood / brown
	// floor.Mat.SetTexture(sc, grtx)
	// floor.Mat.Tiling.Repeat.Set(40, 40)

	sc.Camera.Pose.Pos = mat32.Vec3{0, 20, 30}
	sc.Camera.LookAt(mat32.Vec3{0, 5, 0}, mat32.Vec3Y) // defaults to looking at origin

	appnm := gi.AppName()
	mmen := win.MainMenu
	mmen.ConfigMenus([]string{appnm, "File", "Edit", "Window"})

	amen := win.MainMenu.ChildByName(appnm, 0).(*gi.Action)
	amen.Menu.AddAppMenu(win)

	win.MainMenuUpdated()
	vp.UpdateEndNoSig(updt)
	win.StartEventLoop()
}
