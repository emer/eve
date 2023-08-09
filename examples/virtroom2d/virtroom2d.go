// Copyright (c) 2019, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math/rand"

	"github.com/emer/eve/eve"
	"github.com/emer/eve/eve2d"
	"github.com/goki/gi/gi"
	"github.com/goki/gi/gimain"
	"github.com/goki/gi/giv"
	"github.com/goki/gi/svg"
	"github.com/goki/gi/units"
	"github.com/goki/ki/ki"
	"github.com/goki/mat32"
)

func main() {
	gimain.Main(mainrun)
}

// Env encapsulates the virtual environment
type Env struct {

	// height of emer
	EmerHt float32 `desc:"height of emer"`

	// how far to move every step
	MoveStep float32 `desc:"how far to move every step"`

	// how far to rotate every step
	RotStep float32 `desc:"how far to rotate every step"`

	// width of room
	Width float32 `desc:"width of room"`

	// depth of room
	Depth float32 `desc:"depth of room"`

	// height of room
	Height float32 `desc:"height of room"`

	// thickness of walls of room
	Thick float32 `desc:"thickness of walls of room"`

	// [view: -] world
	World *eve.Group `view:"-" desc:"world"`

	// [view: -] view of world
	View *eve2d.View `view:"-" desc:"view of world"`

	// [view: -] emer group
	Emer *eve.Group `view:"-" desc:"emer group"`

	// [view: -] contacts from last step, for body
	Contacts eve.Contacts `view:"-" desc:"contacts from last step, for body"`

	// [view: -] gui window
	Win *gi.Window `view:"-" desc:"gui window"`
}

func (ev *Env) Defaults() {
	ev.Width = 10
	ev.Depth = 15
	ev.Height = 2
	ev.Thick = 0.2
	ev.EmerHt = 1
	ev.MoveStep = ev.EmerHt * .2
	ev.RotStep = 15
}

// MakeWorld constructs a new virtual physics world
func (ev *Env) MakeWorld() {
	ev.World = &eve.Group{}
	ev.World.InitName(ev.World, "RoomWorld")

	MakeRoom(ev.World, "room1", ev.Width, ev.Depth, ev.Thick)
	ev.Emer = MakeEmer(ev.World, ev.EmerHt)

	ev.World.WorldInit()
}

// InitWorld does init on world and re-syncs
func (ev *Env) WorldInit() {
	ev.World.WorldInit()
	if ev.View != nil {
		ev.View.Sync()
	}
}

// ReMakeWorld rebuilds the world and re-syncs with gui
func (ev *Env) ReMakeWorld() {
	ev.MakeWorld()
	ev.View.World = ev.World
	if ev.View != nil {
		ev.View.Sync()
	}
}

// MakeView makes the view
func (ev *Env) MakeView(sc *svg.Editor) {
	wgp := svg.AddNewGroup(sc, "world")
	ev.View = eve2d.NewView(ev.World, &sc.SVG, wgp)
	ev.View.InitLibrary() // this makes a basic library based on body shapes, sizes
	// at this point the library can be updated to configure custom visualizations
	// for any of the named bodies.
	ev.View.Sync()
}

// WorldStep does one step of the world
func (ev *Env) WorldStep() {
	ev.World.WorldRelToAbs()
	cts := ev.World.WorldCollide(eve.DynsTopGps)
	ev.Contacts = nil
	for _, cl := range cts {
		if len(cl) > 1 {
			for _, c := range cl {
				if c.A.Name() == "body" {
					ev.Contacts = cl
				}
				fmt.Printf("A: %v  B: %v\n", c.A.Name(), c.B.Name())
			}
		}
	}
	if len(ev.Contacts) > 1 { // turn around
		fmt.Printf("hit wall: turn around!\n")
		rot := 100.0 + 90.0*rand.Float32()
		ev.Emer.Rel.RotateOnAxis(0, 1, 0, rot)
	}
	ev.View.UpdatePose()
}

// StepForward moves Emer forward in current facing direction one step
func (ev *Env) StepForward() {
	ev.Emer.Rel.MoveOnAxis(0, 1, 0, -ev.MoveStep)
	ev.WorldStep()
}

// StepBackward moves Emer backward in current facing direction one step
func (ev *Env) StepBackward() {
	ev.Emer.Rel.MoveOnAxis(0, 1, 0, ev.MoveStep)
	ev.WorldStep()
}

// RotBodyLeft rotates emer left
func (ev *Env) RotBodyLeft() {
	ev.Emer.Rel.RotateOnAxis(0, 0, 1, ev.RotStep)
	ev.WorldStep()
}

// RotBodyRight rotates emer right
func (ev *Env) RotBodyRight() {
	ev.Emer.Rel.RotateOnAxis(0, 0, 1, -ev.RotStep)
	ev.WorldStep()
}

// RotHeadLeft rotates emer left
func (ev *Env) RotHeadLeft() {
	hd := ev.Emer.ChildByName("head", 1).(*eve.Group)
	hd.Rel.RotateOnAxis(0, 1, 0, ev.RotStep)
	ev.WorldStep()
}

// RotHeadRight rotates emer right
func (ev *Env) RotHeadRight() {
	hd := ev.Emer.ChildByName("head", 1).(*eve.Group)
	hd.Rel.RotateOnAxis(0, 1, 0, -ev.RotStep)
	ev.WorldStep()
}

// MakeRoom constructs a new room in given parent group with given params
func MakeRoom(par *eve.Group, name string, width, depth, thick float32) *eve.Group {
	rm := eve.AddNewGroup(par, name)
	// floor := eve.AddNewBox(rm, "floor", mat32.Vec3{0, -thick / 2, 0}, mat32.Vec3{width, thick, depth})
	// floor.Color = "grey"
	bwall := eve.AddNewBox(rm, "back-wall", mat32.Vec3{0, -depth / 2, 0}, mat32.Vec3{width, thick, 0})
	bwall.Color = "blue"
	lwall := eve.AddNewBox(rm, "left-wall", mat32.Vec3{-width / 2, 0, 0}, mat32.Vec3{thick, depth, 0})
	lwall.Color = "red"
	rwall := eve.AddNewBox(rm, "right-wall", mat32.Vec3{width / 2, 0, 0}, mat32.Vec3{thick, depth, 0})
	rwall.Color = "green"
	fwall := eve.AddNewBox(rm, "front-wall", mat32.Vec3{0, depth / 2, 0}, mat32.Vec3{width, thick, 0})
	fwall.Color = "brown"
	return rm
}

// MakeEmer constructs a new Emer virtual robot of given height (e.g., 1)
func MakeEmer(par *eve.Group, height float32) *eve.Group {
	emr := eve.AddNewGroup(par, "emer")
	width := height * .4
	depth := height * .15
	body := eve.AddNewBox(emr, "body", mat32.Vec3{0, height / 2, 0}, mat32.Vec3{width, height, depth})
	// body := eve.AddNewCapsule(emr, "body", mat32.Vec3{0, height / 2, 0}, height, width/2)
	// body := eve.AddNewCylinder(emr, "body", mat32.Vec3{0, height / 2, 0}, height, width/2)
	body.Color = "purple"
	body.SetDynamic()

	headsz := depth * 1.5
	hhsz := .5 * headsz
	hgp := eve.AddNewGroup(emr, "head")
	hgp.Initial.Pos = mat32.Vec3{0, height + hhsz, 0}

	head := eve.AddNewBox(hgp, "head", mat32.Vec3{0, 0, 0}, mat32.Vec3{headsz, headsz, headsz})
	head.Color = "tan"
	head.SetDynamic()
	eyesz := headsz * .2
	eyel := eve.AddNewBox(hgp, "eye-l", mat32.Vec3{-hhsz * .6, headsz * .1, -(hhsz + eyesz*.3)}, mat32.Vec3{eyesz, eyesz * .5, eyesz * .2})
	eyel.Color = "green"
	eyel.SetDynamic()
	eyer := eve.AddNewBox(hgp, "eye-r", mat32.Vec3{hhsz * .6, headsz * .1, -(hhsz + eyesz*.3)}, mat32.Vec3{eyesz, eyesz * .5, eyesz * .2})
	eyer.Color = "green"
	eyer.SetDynamic()
	return emr
}

var TheEnv Env

func (ev *Env) ConfigGui() {
	width := 1024
	height := 768

	// vgpu.Debug = true

	gi.SetAppName("virtroom2d")
	gi.SetAppAbout(`This is a demo of the Emergent Virtual Engine.  See <a href="https://github.com/emer/eve">eve on GitHub</a>.
<p>The <a href="https://github.com/emer/eve/blob/master/examples/virtroom/README.md">README</a> page for this example app has further info.</p>`)

	win := gi.NewMainWindow("virtroom2d", "Emergent 2D Virtual Engine", width, height)
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
	split.Dim = mat32.X

	tvfr := gi.AddNewFrame(split, "tvfr", gi.LayoutHoriz)
	svfr := gi.AddNewFrame(split, "svfr", gi.LayoutHoriz)
	scfr := gi.AddNewFrame(split, "scfr", gi.LayoutHoriz)
	split.SetSplits(.1, .3, .6)

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

	scvw := svg.AddNewEditor(scfr, "sceneview")
	scvw.Fill = true
	scvw.SetProp("background-color", "white")
	scvw.SetStretchMaxWidth()
	scvw.SetStretchMaxHeight()
	scvw.InitScale()
	scvw.Trans.Set(600, 540)
	scvw.Scale = 60
	scvw.SetTransform()

	ev.MakeView(scvw)

	tbar.AddAction(gi.ActOpts{Label: "Edit Env", Icon: "edit", Tooltip: "Edit the settings for the environment."}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		sv.SetStruct(ev)
	})
	tbar.AddAction(gi.ActOpts{Label: "Init", Icon: "update", Tooltip: "Initialize virtual world -- go back to starting positions etc."}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		ev.WorldInit()
	})
	tbar.AddAction(gi.ActOpts{Label: "Make", Icon: "update", Tooltip: "Re-make virtual world -- do this if you have changed any of the world parameters."}, win.This(), func(recv, send ki.Ki, sig int64, data interface{}) {
		ev.ReMakeWorld()
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
