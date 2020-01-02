// Copyright (c) 2018, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"image"
	"log"

	"github.com/emer/eve/eve"
	"github.com/emer/eve/evev"
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
	EmerHt   float32          `desc:"height of emer"`
	MoveStep float32          `desc:"how far to move every step"`
	RotStep  float32          `desc:"how far to rotate every step"`
	Width    float32          `desc:"width of room"`
	Depth    float32          `desc:"depth of room"`
	Height   float32          `desc:"height of room"`
	Thick    float32          `desc:"thickness of walls of room"`
	Camera   evev.Camera      `desc:"offscreen render camera settings"`
	DepthMap giv.ColorMapName `desc:"color map to use for rendering depth map"`
	World    *eve.Group       `view:"-" desc:"world"`
	View     *evev.View       `view:"-" desc:"view of world"`
	Emer     *eve.Group       `view:"-" desc:"emer group"`
	EyeR     eve.Body         `view:"-" desc:"Right eye of emer"`
	Win      *gi.Window       `view:"-" desc:"gui window"`
	SnapImg  *gi.Bitmap       `view:"-" desc:"snapshot bitmap view"`
	DepthImg *gi.Bitmap       `view:"-" desc:"depth map bitmap view"`
	Frame    gpu.Framebuffer  `view:"-" desc:"offscreen render buffer"`
}

func (ev *Env) Defaults() {
	ev.Width = 10
	ev.Depth = 15
	ev.Height = 2
	ev.Thick = 0.2
	ev.EmerHt = 1
	ev.MoveStep = ev.EmerHt * .2
	ev.RotStep = 15
	ev.DepthMap = giv.ColorMapName("ColdHot")
	ev.Camera.Defaults()
	ev.Camera.FOV = 90
}

// MakeWorld constructs a new virtual physics world
func (ev *Env) MakeWorld() {
	ev.World = &eve.Group{}
	ev.World.InitName(ev.World, "RoomWorld")

	MakeRoom(ev.World, "room1", ev.Width, ev.Depth, ev.Height, ev.Thick)
	ev.Emer = MakeEmer(ev.World, ev.EmerHt)
	ev.EyeR = ev.Emer.ChildByName("head", 1).ChildByName("eye-r", 2).(eve.Body)

	ev.World.InitWorld()
}

// InitWorld does init on world and re-syncs
func (ev *Env) InitWorld() {
	ev.World.InitWorld()
	ev.View.Sync()
	ev.Snapshot()
}

// ReMakeWorld rebuilds the world and re-syncs with gui
func (ev *Env) ReMakeWorld() {
	ev.MakeWorld()
	ev.View.World = ev.World
	ev.View.Sync()
	ev.Snapshot()
}

// MakeView makes the view
func (ev *Env) MakeView(sc *gi3d.Scene) {
	wgp := gi3d.AddNewGroup(sc, sc, "world")
	ev.View = evev.NewView(ev.World, sc, wgp)
	ev.View.InitLibrary() // this makes a basic library based on body shapes, sizes
	// at this point the library can be updated to configure custom visualizations
	// for any of the named bodies.
	ev.View.Sync()
}

// Snapshot takes a snapshot from the perspective of Emer's right eye
func (ev *Env) Snapshot() {
	err := ev.View.RenderOffNode(&ev.Frame, ev.EyeR, &ev.Camera)
	if err != nil {
		log.Println(err)
		return
	}
	var img image.Image
	var depth []float32
	oswin.TheApp.RunOnMain(func() {
		tex := ev.Frame.Texture()
		tex.SetBotZero(true)
		img = tex.GrabImage()
		depth = ev.Frame.DepthAll()
	})
	ev.SnapImg.SetImage(img, 0, 0)
	ev.ViewDepth(depth)
	ev.View.Scene.Render2D()
	ev.View.Scene.DirectWinUpload()
}

// ViewDepth updates depth bitmap with depth data
func (ev *Env) ViewDepth(depth []float32) {
	cmap := giv.AvailColorMaps[string(ev.DepthMap)]
	ev.DepthImg.Resize(ev.Camera.Size)
	evev.DepthImage(ev.DepthImg.Pixels, depth, cmap, &ev.Camera)
	ev.DepthImg.UpdateSig()
}

// StepForward moves Emer forward in current facing direction one step, and takes Snapshot
func (ev *Env) StepForward() {
	ev.Emer.Rel.MoveOnAxis(0, 0, 1, -ev.MoveStep)
	ev.World.UpdateWorld()
	ev.View.UpdatePose()
	ev.Snapshot()
}

// StepBackward moves Emer backward in current facing direction one step, and takes Snapshot
func (ev *Env) StepBackward() {
	ev.Emer.Rel.MoveOnAxis(0, 0, 1, ev.MoveStep)
	ev.World.UpdateWorld()
	ev.View.UpdatePose()
	ev.Snapshot()
}

// RotBodyLeft rotates emer left and takes Snapshot
func (ev *Env) RotBodyLeft() {
	ev.Emer.Rel.RotateOnAxis(0, 1, 0, ev.RotStep)
	ev.World.UpdateWorld()
	ev.View.UpdatePose()
	ev.Snapshot()
}

// RotBodyRight rotates emer right and takes Snapshot
func (ev *Env) RotBodyRight() {
	ev.Emer.Rel.RotateOnAxis(0, 1, 0, -ev.RotStep)
	ev.World.UpdateWorld()
	ev.View.UpdatePose()
	ev.Snapshot()
}

// RotHeadLeft rotates emer left and takes Snapshot
func (ev *Env) RotHeadLeft() {
	hd := ev.Emer.ChildByName("head", 1).(*eve.Group)
	hd.Rel.RotateOnAxis(0, 1, 0, ev.RotStep)
	ev.World.UpdateWorld()
	ev.View.UpdatePose()
	ev.Snapshot()
}

// RotHeadRight rotates emer right and takes Snapshot
func (ev *Env) RotHeadRight() {
	hd := ev.Emer.ChildByName("head", 1).(*eve.Group)
	hd.Rel.RotateOnAxis(0, 1, 0, -ev.RotStep)
	ev.World.UpdateWorld()
	ev.View.UpdatePose()
	ev.Snapshot()
}

// MakeRoom constructs a new room in given parent group with given params
func MakeRoom(par *eve.Group, name string, width, depth, height, thick float32) *eve.Group {
	rm := eve.AddNewGroup(par, name)
	floor := eve.AddNewBox(rm, "floor", mat32.Vec3{0, -thick / 2, 0}, mat32.Vec3{width, thick, depth})
	floor.Color = "grey"
	bwall := eve.AddNewBox(rm, "back-wall", mat32.Vec3{0, height / 2, -depth / 2}, mat32.Vec3{width, height, thick})
	bwall.Color = "blue"
	lwall := eve.AddNewBox(rm, "left-wall", mat32.Vec3{-width / 2, height / 2, 0}, mat32.Vec3{thick, height, depth})
	lwall.Color = "red"
	rwall := eve.AddNewBox(rm, "right-wall", mat32.Vec3{width / 2, height / 2, 0}, mat32.Vec3{thick, height, depth})
	rwall.Color = "green"
	fwall := eve.AddNewBox(rm, "front-wall", mat32.Vec3{0, height / 2, depth / 2}, mat32.Vec3{width, height, thick})
	fwall.Color = "yellow"
	return rm
}

// MakeEmer constructs a new Emer virtual robot of given height (e.g., 1)
func MakeEmer(par *eve.Group, height float32) *eve.Group {
	emr := eve.AddNewGroup(par, "emer")
	width := height * .4
	depth := height * .15
	body := eve.AddNewBox(emr, "body", mat32.Vec3{0, height / 2, 0}, mat32.Vec3{width, height, depth})
	body.Color = "purple"

	headsz := depth * 1.5
	hhsz := .5 * headsz
	hgp := eve.AddNewGroup(emr, "head")
	hgp.Initial.Pos = mat32.Vec3{0, height + hhsz, 0}

	head := eve.AddNewBox(hgp, "head", mat32.Vec3{0, 0, 0}, mat32.Vec3{headsz, headsz, headsz})
	head.Color = "tan"
	eyesz := headsz * .2
	eyel := eve.AddNewBox(hgp, "eye-l", mat32.Vec3{-hhsz * .6, headsz * .1, -(hhsz + eyesz*.3)}, mat32.Vec3{eyesz, eyesz * .5, eyesz * .2})
	eyel.Color = "green"
	eyer := eve.AddNewBox(hgp, "eye-r", mat32.Vec3{hhsz * .6, headsz * .1, -(hhsz + eyesz*.3)}, mat32.Vec3{eyesz, eyesz * .5, eyesz * .2})
	eyer.Color = "green"
	return emr
}

var TheEnv Env

func (ev *Env) ConfigGui() {
	width := 1024
	height := 768

	gi.SetAppName("virtroom")
	gi.SetAppAbout(`This is a demo of the Emergent Virtual Engine.  See <a href="https://github.com/emer/eve">eve on GitHub</a>.
<p>The <a href="https://github.com/emer/eve/blob/master/examples/virtroom/README.md">README</a> page for this example app has further info.</p>`)

	win := gi.NewMainWindow("virtroom", "Emergent Virtual Engine", width, height)
	ev.Win = win

	vp := win.WinViewport2D()
	updt := vp.UpdateStart()

	mfr := win.SetMainFrame()
	mfr.SetProp("spacing", units.NewEx(1))

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
	split.SetSplits(.1, .2, .2, .5)

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

	scvw := gi3d.AddNewSceneView(scfr, "sceneview")
	scvw.SetStretchMaxWidth()
	scvw.SetStretchMaxHeight()
	scvw.Config()
	sc := scvw.Scene()

	// first, add lights, set camera
	sc.BgColor.SetUInt8(230, 230, 255, 255) // sky blue-ish
	gi3d.AddNewAmbientLight(sc, "ambient", 0.3, gi3d.DirectSun)

	dir := gi3d.AddNewDirLight(sc, "dir", 1, gi3d.DirectSun)
	dir.Pos.Set(0, 2, 1) // default: 0,1,1 = above and behind us (we are at 0,0,X)

	ev.MakeView(sc)

	// grtx := gi3d.AddNewTextureFile(sc, "ground", "ground.png")
	// wdtx := gi3d.AddNewTextureFile(sc, "wood", "wood.png")

	// floorp := gi3d.AddNewPlane(sc, "floor-plane", 100, 100)
	// floor := gi3d.AddNewSolid(sc, sc, "floor", floorp.Name())
	// floor.Pose.Pos.Set(0, -5, 0)
	// // floor.Mat.Color.SetName("tan")
	// // floor.Mat.Emissive.SetName("brown")
	// floor.Mat.Bright = 2 // .5 for wood / brown
	// floor.Mat.SetTexture(sc, grtx)
	// floor.Mat.Tiling.Reveat.Set(40, 40)

	sc.Camera.Pose.Pos = mat32.Vec3{0, 40, 3.5}
	sc.Camera.LookAt(mat32.Vec3{0, 5, 0}, mat32.Vec3Y)
	sc.SaveCamera("3")

	sc.Camera.Pose.Pos = mat32.Vec3{0, 20, 30}
	sc.Camera.LookAt(mat32.Vec3{0, 5, 0}, mat32.Vec3Y)
	sc.SaveCamera("2")

	sc.Camera.Pose.Pos = mat32.Vec3{-.86, .97, 2.7}
	sc.Camera.LookAt(mat32.Vec3{0, .8, 0}, mat32.Vec3Y)
	sc.SaveCamera("1")
	sc.SaveCamera("default")

	//////////////////////////////////////////
	//    Bitmap

	imfr.Lay = gi.LayoutVert
	gi.AddNewLabel(imfr, "lab-img", "Right Eye Image:")
	ev.SnapImg = gi.AddNewBitmap(imfr, "snap-img")
	ev.SnapImg.Resize(ev.Camera.Size)
	ev.SnapImg.LayoutToImgSize()
	ev.SnapImg.SetProp("vertical-align", gi.AlignTop)

	gi.AddNewLabel(imfr, "lab-depth", "Right Eye Depth:")
	ev.DepthImg = gi.AddNewBitmap(imfr, "depth-img")
	ev.DepthImg.Resize(ev.Camera.Size)
	ev.DepthImg.LayoutToImgSize()
	ev.DepthImg.SetProp("vertical-align", gi.AlignTop)

	tbar.AddAction(gi.ActOpts{Label: "Edit Env", Icon: "edit", Tooltip: "Edit the settings for the environment."}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		sv.SetStruct(ev)
	})
	tbar.AddAction(gi.ActOpts{Label: "Init", Icon: "update", Tooltip: "Initialize virtual world -- go back to starting positions etc."}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		ev.InitWorld()
	})
	tbar.AddAction(gi.ActOpts{Label: "Make", Icon: "update", Tooltip: "Re-make virtual world -- do this if you have changed any of the world parameters."}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		ev.ReMakeWorld()
	})
	tbar.AddAction(gi.ActOpts{Label: "Snap", Icon: "file-image", Tooltip: "Take a snapshot from perspective of the right eye of emer virtual robot."}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		ev.Snapshot()
	})
	tbar.AddSeparator("mv-sep")
	tbar.AddAction(gi.ActOpts{Label: "Fwd", Icon: "wedge-up", Tooltip: "Take a step Forward."}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		ev.StepForward()
	})
	tbar.AddAction(gi.ActOpts{Label: "Bkw", Icon: "wedge-down", Tooltip: "Take a step Backward."}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		ev.StepBackward()
	})
	tbar.AddAction(gi.ActOpts{Label: "Body Left", Icon: "wedge-left", Tooltip: "Rotate body left."}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		ev.RotBodyLeft()
	})
	tbar.AddAction(gi.ActOpts{Label: "Body Right", Icon: "wedge-right", Tooltip: "Rotate body right."}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		ev.RotBodyRight()
	})
	tbar.AddAction(gi.ActOpts{Label: "Head Left", Icon: "wedge-left", Tooltip: "Rotate body left."}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		ev.RotHeadLeft()
	})
	tbar.AddAction(gi.ActOpts{Label: "Head Right", Icon: "wedge-right", Tooltip: "Rotate body right."}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		ev.RotHeadRight()
	})
	tbar.AddSeparator("rm-sep")
	tbar.AddAction(gi.ActOpts{Label: "README", Icon: "file-markdown", Tooltip: "Open browser on README."}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		gi.OpenURL("https://github.com/emer/eve/blob/master/examples/virtroom/README.md")
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
