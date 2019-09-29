// Copyright (c) 2018, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"image"

	"github.com/emer/epe/epe"
	"github.com/emer/epe/epev"
	"github.com/goki/gi/gi"
	"github.com/goki/gi/gi3d"
	"github.com/goki/gi/gimain"
	"github.com/goki/gi/giv"
	"github.com/goki/gi/mat32"
	"github.com/goki/gi/oswin"
	"github.com/goki/gi/oswin/gpu"
	"github.com/goki/gi/units"
	"github.com/goki/ki/ki"
)

func main() {
	gimain.Main(func() {
		mainrun()
	})
}

// Env encapsulates the virtual environment
type Env struct {
	MoveStep float32     `desc:"how far to move every step"`
	Width    float32     `desc:"width of room"`
	Depth    float32     `desc:"depth of room"`
	Height   float32     `desc:"height of room"`
	Thick    float32     `desc:"thickness of walls of room"`
	EmerHt   float32     `desc:"height of emer"`
	Camera   epev.Camera `desc:"offscreen render camera settings"`
	World    *epe.Group  `view:"-" desc:"world"`
	View     *epev.View  `view:"-" desc:"view of world"`
	Emer     *epe.Group
	EyeR     epe.Body        `view:"Right eye of emer"`
	Win      *gi.Window      `view:"-" desc:"gui window"`
	SnapImg  *gi.Bitmap      `view:"-" desc:"snapshot bitmap view"`
	Frame    gpu.Framebuffer `view:"-" desc:"offscreen render buffer"`
}

func (ev *Env) Defaults() {
	ev.MoveStep = 0.1
	ev.Width = 10
	ev.Depth = 15
	ev.Height = 8
	ev.Thick = 0.2
	ev.EmerHt = 1
	ev.Camera.Defaults()
	ev.Camera.FOV = 90
}

// MakeWorld constructs a new virtual physics world
func (ev *Env) MakeWorld() {
	ev.World = &epe.Group{}
	ev.World.InitName(ev.World, "RoomWorld")

	MakeRoom(ev.World, "room1", ev.Width, ev.Depth, ev.Height, ev.Thick)
	ev.Emer = MakeEmer(ev.World, ev.EmerHt)
	ev.EyeR = ev.Emer.ChildByName("head", 1).ChildByName("eye-r", 2).(epe.Body)

	ev.World.InitWorld()
}

// MakeView makes the view
func (ev *Env) MakeView(sc *gi3d.Scene) {
	wgp := gi3d.AddNewGroup(sc, sc, "world")
	ev.View = epev.NewView(ev.World, sc, wgp)
	ev.View.Sync()
}

// Snapshot takes a snapshot from the perspective of Emer's right eye
func (ev *Env) Snapshot() {
	ev.View.RenderOffNode(&ev.Frame, ev.EyeR, &ev.Camera)
	var img image.Image
	oswin.TheApp.RunOnMain(func() {
		tex := ev.Frame.Texture()
		tex.SetBotZero(true)
		img = tex.GrabImage()
	})
	gi.SaveImage("test.png", img)
	ev.SnapImg.SetImage(img, 0, 0)
}

// StepForward moves Emer forward in current facing direction one step, and takes Snapshot
func (ev *Env) StepForward() {
	ev.Emer.Rel.MoveOnAxis(0, 0, 1, -ev.MoveStep)
	ev.World.UpdateWorld()
	ev.View.Sync() // todo: just pos
	ev.Snapshot()
}

// StepBackward moves Emer backward in current facing direction one step, and takes Snapshot
func (ev *Env) StepBackward() {
	ev.Emer.Rel.MoveOnAxis(0, 0, 1, ev.MoveStep)
	ev.World.UpdateWorld()
	ev.View.Sync()
	ev.Snapshot()
}

// MakeRoom constructs a new room in given parent group with given params
func MakeRoom(par *epe.Group, name string, width, depth, height, thick float32) *epe.Group {
	rm := epe.AddNewGroup(par, name)
	bwall := epe.AddNewBox(rm, "back-wall", mat32.Vec3{0, height / 2, -depth / 2}, mat32.Vec3{width, height, thick})
	bwall.Mat.Color = "blue"
	lwall := epe.AddNewBox(rm, "left-wall", mat32.Vec3{-width / 2, height / 2, 0}, mat32.Vec3{thick, height, depth})
	lwall.Mat.Color = "red"
	rwall := epe.AddNewBox(rm, "right-wall", mat32.Vec3{width / 2, height / 2, 0}, mat32.Vec3{thick, height, depth})
	rwall.Mat.Color = "green"
	fwall := epe.AddNewBox(rm, "front-wall", mat32.Vec3{0, height / 2, depth / 2}, mat32.Vec3{width, height, thick})
	fwall.Mat.Color = "yellow"
	return rm
}

// MakeEmer constructs a new Emer virtual robot of given height (e.g., 1)
func MakeEmer(par *epe.Group, height float32) *epe.Group {
	emr := epe.AddNewGroup(par, "emer")
	width := height * .3
	depth := height * .15
	body := epe.AddNewBox(emr, "body", mat32.Vec3{0, height / 2, 0}, mat32.Vec3{width, height, depth})
	body.Mat.Color = "purple"

	headsz := depth * 1.5
	hhsz := .5 * headsz
	hgp := epe.AddNewGroup(emr, "head")
	hgp.Initial.Pos = mat32.Vec3{0, height + hhsz, 0}

	head := epe.AddNewBox(hgp, "head", mat32.Vec3{0, 0, 0}, mat32.Vec3{headsz, headsz, headsz})
	head.Mat.Color = "tan"
	eyesz := headsz * .2
	eyel := epe.AddNewBox(hgp, "eye-l", mat32.Vec3{-hhsz * .6, headsz * .1, -(hhsz + eyesz*.3)}, mat32.Vec3{eyesz, eyesz * .5, eyesz * .2})
	eyel.Mat.Color = "green"
	eyer := epe.AddNewBox(hgp, "eye-r", mat32.Vec3{hhsz * .6, headsz * .1, -(hhsz + eyesz*.3)}, mat32.Vec3{eyesz, eyesz * .5, eyesz * .2})
	eyer.Mat.Color = "green"
	return emr
}

var TheEnv Env

func (ev *Env) ConfigGui() {
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
	ev.Win = win

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

	ev.MakeWorld()

	tbar := gi.AddNewToolBar(mfr, "main-tbar")
	tbar.SetStretchMaxWidth()
	tbar.Viewport = vp

	//////////////////////////////////////////
	//    Splitter

	split := gi.AddNewSplitView(mfr, "split")
	split.Dim = gi.X

	tvfr := gi.AddNewFrame(split, "tvfr", gi.LayoutHoriz)
	svfr := gi.AddNewFrame(split, "svfr", gi.LayoutHoriz)
	imfr := gi.AddNewFrame(split, "imfr", gi.LayoutHoriz)
	scfr := gi.AddNewFrame(split, "scfr", gi.LayoutHoriz)
	split.SetSplits(.1, .1, .2, .6)

	tv := giv.AddNewTreeView(tvfr, "tv")
	tv.SetRootNode(ev.World)

	sv := giv.AddNewStructView(svfr, "sv")
	sv.SetStretchMaxWidth()
	sv.SetStretchMaxHeight()
	sv.SetStruct(ev)

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

	ev.MakeView(sc)

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
	sc.SaveCamera("2")

	sc.Camera.Pose.Pos = mat32.Vec3{-.86, .97, 2.7}
	sc.Camera.LookAt(mat32.Vec3{0, .8, 0}, mat32.Vec3Y) // defaults to looking at origin
	sc.SaveCamera("1")
	sc.SaveCamera("default")

	//////////////////////////////////////////
	//    Bitmap

	ev.SnapImg = gi.AddNewBitmap(imfr, "snapimg")
	ev.SnapImg.Resize(ev.Camera.Size)
	ev.SnapImg.LayoutToImgSize()
	ev.SnapImg.SetProp("vertical-align", gi.AlignTop)

	tbar.AddAction(gi.ActOpts{Label: "Snap", Icon: "file-image", Tooltip: "Take a snapshot from perspective of the right eye of emer virtual robot."}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		ev.Snapshot()
		vp.SetNeedsFullRender()
	})
	tbar.AddAction(gi.ActOpts{Label: "Fwd", Icon: "wedge-up", Tooltip: "Take a step Forward."}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		ev.StepForward()
		vp.SetNeedsFullRender()
	})
	tbar.AddAction(gi.ActOpts{Label: "Bkw", Icon: "wedge-down", Tooltip: "Take a step Backward."}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		ev.StepBackward()
		vp.SetNeedsFullRender()
	})

	appnm := gi.AppName()
	mmen := win.MainMenu
	mmen.ConfigMenus([]string{appnm, "File", "Edit", "Window"})

	amen := win.MainMenu.ChildByName(appnm, 0).(*gi.Action)
	amen.Menu.AddAppMenu(win)

	win.MainMenuUpdated()
	vp.UpdateEndNoSig(updt)
}

func mainrun() {
	ev := &TheEnv
	ev.Defaults()
	ev.ConfigGui()
	ev.Win.StartEventLoop()
}
