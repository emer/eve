// Copyright (c) 2019, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package evev

import (
	"fmt"
	"image"

	"github.com/emer/eve/eve"
	"github.com/goki/gi/gi3d"
	"github.com/goki/ki/ki"
	"github.com/goki/ki/kit"
)

// View connects a Virtual World with a Gi3D Scene to visualize the world,
// including ability to render offscreen
type View struct {
	World *eve.Group  `desc:"the root Group node of the virtual world"`
	Scene *gi3d.Scene `desc:"the scene object for visualizing"`
	Root  *gi3d.Group `desc:"the root Group node in the Scene where the world is rendered under"`
}

var KiT_View = kit.Types.AddType(&View{}, nil)

// NewView returns a new View that links given world with given scene and root group
func NewView(world *eve.Group, sc *gi3d.Scene, root *gi3d.Group) *View {
	vw := &View{World: world, Scene: sc, Root: root}
	return vw
}

// InitLibrary initializes Scene library with basic Solid shapes
// based on bodies in the virtual world.  More complex visualizations
// can be configured after this.
func (vw *View) InitLibrary() {
	InitLibrary(vw.World, vw.Scene)
}

// Sync synchronizes the view to the world
func (vw *View) Sync() bool {
	rval := SyncNode(vw.World, vw.Root, vw.Scene)
	return rval
}

// UpdatePose updates the view pose values only from world tree.
// Essential that both trees are already synchronized.
func (vw *View) UpdatePose() {
	UpdatePose(vw.World, vw.Root)
}

// RenderOffNode does an offscreen render using given node
// for the camera position and orientation.
// Current scene camera is saved and restored
func (vw *View) RenderOffNode(node eve.Node, cam *Camera) error {
	sc := vw.Scene
	camnm := "eve-view-renderoff-save"
	sc.SaveCamera(camnm)
	defer sc.SetCamera(camnm)
	sc.Camera.FOV = cam.FOV
	sc.Camera.Near = cam.Near
	sc.Camera.Far = cam.Far
	nb := node.AsNodeBase()
	sc.Camera.Pose.Pos = nb.Abs.Pos
	sc.Camera.Pose.Quat = nb.Abs.Quat
	sc.Camera.Pose.Scale.Set(1, 1, 1)
	sz := sc.Geom.Size
	sc.Geom.Size = cam.Size
	sc.Frame.SetSize(sc.Geom.Size) // nop if same
	ok := sc.RenderOffscreen()
	sc.Geom.Size = sz
	if !ok {
		return fmt.Errorf("could not render scene")
	}
	return nil
}

// Image returns the current rendered image
func (vw *View) Image() (*image.RGBA, error) {
	fr := vw.Scene.Frame
	if fr == nil {
		return nil, fmt.Errorf("eve.View Image: Scene does not have a Frame")
	}
	sy := &vw.Scene.Phong.Sys
	tcmd := sy.MemCmdStart()
	fr.GrabImage(tcmd, 0)
	sy.MemCmdEndSubmitWaitFree()
	gimg, err := fr.Render.Grab.DevGoImage()
	if err == nil {
		return gimg, err
	}
	return nil, err
}

// DepthImage returns the current rendered depth image
func (vw *View) DepthImage() ([]float32, error) {
	fr := vw.Scene.Frame
	if fr == nil {
		return nil, fmt.Errorf("eve.View Image: Scene does not have a Frame")
	}
	sy := &vw.Scene.Phong.Sys
	tcmd := sy.MemCmdStart()
	fr.GrabDepthImage(tcmd)
	sy.MemCmdEndSubmitWaitFree()

	depth, err := fr.Render.DepthImageArray()
	if err == nil {
		return depth, err
	}
	return nil, err
}

///////////////////////////////////////////////////////////////
// Sync, Config

// InitLibrary initializes Scene library with basic Solid shapes
// based on bodies in the virtual world.  More complex visualizations
// can be configured after this.
func InitLibrary(wn eve.Node, sc *gi3d.Scene) {
	bod := wn.AsBody()
	if bod != nil {
		InitLibSolid(bod, sc)
	}
	for idx := range *wn.Children() {
		wk := wn.Child(idx).(eve.Node)
		InitLibrary(wk, sc)
	}
}

// InitLibSolid initializes Scene library with Solid for given body
func InitLibSolid(bod eve.Body, sc *gi3d.Scene) {
	nm := bod.Name()
	lgp, has := sc.Library[nm]
	if !has {
		lgp = sc.NewInLibrary(nm)
	}
	bod.AsBodyBase().Vis = nm
	var sld *gi3d.Solid
	if lgp.HasChildren() {
		sld, has = lgp.Child(0).(*gi3d.Solid)
		if !has {
			return // some other kind of thing already configured
		}
	} else {
		sld = gi3d.AddNewSolid(sc, lgp, nm, "")
	}
	wt := kit.ShortTypeName(ki.Type(bod.This()))
	switch wt {
	case "eve.Box":
		mnm := "eveBox"
		bx := bod.(*eve.Box)
		bm := sc.MeshByName(mnm)
		if bm == nil {
			bm = gi3d.AddNewBox(sc, mnm, 1, 1, 1)
		}
		sld.SetMeshName(sc, mnm)
		sld.Pose.Scale = bx.Size
		if bx.Color != "" {
			sld.Mat.Color.SetName(bx.Color)
		}
	case "eve.Cylinder":
		mnm := "eveCylinder"
		cy := bod.(*eve.Cylinder)
		cm := sc.MeshByName(mnm)
		if cm == nil {
			cm = gi3d.AddNewCylinder(sc, mnm, 1, 1, 32, 1, true, true)
		}
		sld.SetMeshName(sc, mnm)
		sld.Pose.Scale.Set(cy.BotRad, cy.Height, cy.BotRad)
		if cy.Color != "" {
			sld.Mat.Color.SetName(cy.Color)
		}
	case "eve.Capsule":
		mnm := "eveCapsule"
		cp := bod.(*eve.Capsule)
		cm := sc.MeshByName(mnm)
		if cm == nil {
			cm = gi3d.AddNewCapsule(sc, mnm, 1, .2, 32, 1)
		}
		sld.SetMeshName(sc, mnm)
		sld.Pose.Scale.Set(cp.BotRad/.2, cp.Height/1.4, cp.BotRad/.2)
		if cp.Color != "" {
			sld.Mat.Color.SetName(cp.Color)
		}
	case "eve.Sphere":
		mnm := "eveSphere"
		sp := bod.(*eve.Sphere)
		sm := sc.MeshByName(mnm)
		if sm == nil {
			sm = gi3d.AddNewSphere(sc, mnm, 1, 32)
		}
		sld.SetMeshName(sc, mnm)
		sld.Pose.Scale.SetScalar(sp.Radius)
		if sp.Color != "" {
			sld.Mat.Color.SetName(sp.Color)
		}
	}
}

// ConfigBodySolid configures a solid for a body with current values
func ConfigBodySolid(bod eve.Body, sld *gi3d.Solid) {
	wt := kit.ShortTypeName(ki.Type(bod.This()))
	switch wt {
	case "eve.Box":
		bx := bod.(*eve.Box)
		sld.Pose.Scale = bx.Size
		if bx.Color != "" {
			sld.Mat.Color.SetName(bx.Color)
		}
	case "eve.Cylinder":
		cy := bod.(*eve.Cylinder)
		sld.Pose.Scale.Set(cy.BotRad, cy.Height, cy.BotRad)
		if cy.Color != "" {
			sld.Mat.Color.SetName(cy.Color)
		}
	case "eve.Capsule":
		cp := bod.(*eve.Capsule)
		sld.Pose.Scale.Set(cp.BotRad/.2, cp.Height/1.4, cp.BotRad/.2)
		if cp.Color != "" {
			sld.Mat.Color.SetName(cp.Color)
		}
	case "eve.Sphere":
		sp := bod.(*eve.Sphere)
		sld.Pose.Scale.SetScalar(sp.Radius)
		if sp.Color != "" {
			sld.Mat.Color.SetName(sp.Color)
		}
	}
}

// ConfigView configures the view node to properly display world node
func ConfigView(wn eve.Node, vn gi3d.Node3D, sc *gi3d.Scene) {
	wb := wn.AsNodeBase()
	vb := vn.(*gi3d.Group)
	vb.Pose.Pos = wb.Rel.Pos
	vb.Pose.Quat = wb.Rel.Quat
	bod := wn.AsBody()
	if bod != nil {
		if !vb.HasChildren() {
			sc.AddFmLibrary(bod.AsBodyBase().Vis, vb)
		} else {
			bgp := vb.Child(0)
			if bgp.HasChildren() {
				sld, has := bgp.Child(0).(*gi3d.Solid)
				if has {
					ConfigBodySolid(bod, sld)
				}
			}
		}
	}
}

// SyncNode updates the view tree to match the world tree, using
// ConfigChildren to maximally preserve existing tree elements
// returns true if view tree was modified (elements added / removed etc)
func SyncNode(wn eve.Node, vn gi3d.Node3D, sc *gi3d.Scene) bool {
	nm := wn.Name()
	vn.SetName(nm) // guaranteed to be unique
	skids := *wn.Children()
	tnl := make(kit.TypeAndNameList, 0, len(skids))
	for _, skid := range skids {
		tnl.Add(gi3d.KiT_Group, skid.Name())
	}
	mod, updt := vn.ConfigChildren(tnl)
	modall := mod
	for idx := range skids {
		wk := wn.Child(idx).(eve.Node)
		vk := vn.Child(idx).(gi3d.Node3D)
		ConfigView(wk, vk, sc)
		if wk.HasChildren() {
			kmod := SyncNode(wk, vk, sc)
			if kmod {
				modall = true
			}
		}
	}
	vn.UpdateEnd(updt)
	return modall
}

///////////////////////////////////////////////////////////////
// UpdatePose

// UpdatePose updates the view pose values only from world tree.
// Essential that both trees are already synchronized.
func UpdatePose(wn eve.Node, vn gi3d.Node3D) {
	skids := *wn.Children()
	for idx := range skids {
		wk := wn.Child(idx).(eve.Node)
		vk := vn.Child(idx).(gi3d.Node3D)
		wb := wn.AsNodeBase()
		vb := vn.AsNode3D()
		vb.Pose.Pos = wb.Rel.Pos
		vb.Pose.Quat = wb.Rel.Quat
		UpdatePose(wk, vk)
	}
}
