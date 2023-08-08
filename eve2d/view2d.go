// Copyright (c) 2019, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eve2d

import (
	"fmt"
	"image"

	"github.com/emer/eve/eve"
	"github.com/goki/gi/svg"
	"github.com/goki/ki/ki"
	"github.com/goki/ki/kit"
)

// View connects a Virtual World with a 2D SVG Scene to visualize the world
type View struct {

	// the root Group node of the virtual world
	World *eve.Group `desc:"the root Group node of the virtual world"`

	// the SVG rendering canvas for visualizing in 2D
	Scene *svg.SVG `desc:"the SVG rendering canvas for visualizing in 2D"`

	// the root Group node in the Scene under which the world is rendered
	Root *svg.Group `desc:"the root Group node in the Scene under which the world is rendered"`

	// library of shapes for bodies -- name matches Body.Vis
	Library map[string]*svg.Group `desc:"library of shapes for bodies -- name matches Body.Vis"`
}

var KiT_View = kit.Types.AddType(&View{}, nil)

// NewView returns a new View that links given world with given scene and root group
func NewView(world *eve.Group, sc *svg.SVG, root *svg.Group) *View {
	vw := &View{World: world, Scene: sc, Root: root}
	vw.Library = make(map[string]*svg.Group)
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

// Image returns the current rendered image
func (vw *View) Image() (*image.RGBA, error) {
	img := vw.Scene.Pixels
	if img == nil {
		return nil, fmt.Errorf("eve2d.View Image: is nil")
	}
	return img, nil
}

///////////////////////////////////////////////////////////////
// Sync, Config

// InitLibrary initializes Scene library with basic Solid shapes
// based on bodies in the virtual world.  More complex visualizations
// can be configured after this.
func (vw *View) InitLibrary(wn eve.Node) {
	bod := wn.AsBody()
	if bod != nil {
		vw.InitLibSolid(bod)
	}
	for idx := range *wn.Children() {
		wk := wn.Child(idx).(eve.Node)
		vw.InitLibrary(wk)
	}
}

// NewInLibrary adds a new item of given name in library
func (vw *View) NewInLibrary(nm string) {
	if vw.Library == nil {
		vw.Library = make(map[string]*svg.Group)
	}
	// vw.Library[nm] =
}

// InitLibSolid initializes Scene library with Solid for given body
func (vw *View) InitLibSolid(bod eve.Body) {
	nm := bod.Name()
	lgp, has := vw.Library[nm]
	if !has {
		lgp = vw.NewInLibrary(nm)
	}
	bod.AsBodyBase().Vis = nm
	var sld *svg.Solid
	// if lgp.HasChildren() {
	// 	sld, has = lgp.Child(0).(*svg.Solid)
	// 	if !has {
	// 		return // some other kind of thing already configured
	// 	}
	// } else {
	sld = svg.AddNewSolid(sc, lgp, nm, "")
	// }
	wt := kit.ShortTypeName(ki.Type(bod.This()))
	switch wt {
	case "eve.Box":
		mnm := "eveBox"
		bx := bod.(*eve.Box)
		bm := sc.MeshByName(mnm)
		if bm == nil {
			bm = svg.AddNewBox(sc, mnm, 1, 1, 1)
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
			cm = svg.AddNewCylinder(sc, mnm, 1, 1, 32, 1, true, true)
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
			cm = svg.AddNewCapsule(sc, mnm, 1, .2, 32, 1)
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
			sm = svg.AddNewSphere(sc, mnm, 1, 32)
		}
		sld.SetMeshName(sc, mnm)
		sld.Pose.Scale.SetScalar(sp.Radius)
		if sp.Color != "" {
			sld.Mat.Color.SetName(sp.Color)
		}
	}
}

// ConfigBodySolid configures a solid for a body with current values
func ConfigBodySolid(bod eve.Body, sld *svg.Solid) {
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
func ConfigView(wn eve.Node, vn svg.Node3D, sc *svg.SVG) {
	wb := wn.AsNodeBase()
	vb := vn.(*svg.Group)
	vb.Pose.Pos = wb.Rel.Pos
	vb.Pose.Quat = wb.Rel.Quat
	bod := wn.AsBody()
	if bod != nil {
		if !vb.HasChildren() {
			sc.AddFmLibrary(bod.AsBodyBase().Vis, vb)
		} else {
			bgp := vb.Child(0)
			if bgp.HasChildren() {
				sld, has := bgp.Child(0).(*svg.Solid)
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
func SyncNode(wn eve.Node, vn svg.Node3D, sc *svg.SVG) bool {
	nm := wn.Name()
	vn.SetName(nm) // guaranteed to be unique
	skids := *wn.Children()
	tnl := make(kit.TypeAndNameList, 0, len(skids))
	for _, skid := range skids {
		tnl.Add(svg.KiT_Group, skid.Name())
	}
	mod, updt := vn.ConfigChildren(tnl)
	modall := mod
	for idx := range skids {
		wk := wn.Child(idx).(eve.Node)
		vk := vn.Child(idx).(svg.Node3D)
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
func UpdatePose(wn eve.Node, vn svg.Node3D) {
	skids := *wn.Children()
	for idx := range skids {
		wk := wn.Child(idx).(eve.Node)
		vk := vn.Child(idx).(svg.Node3D)
		wb := wn.AsNodeBase()
		vb := vn.AsNode3D()
		vb.Pose.Pos = wb.Rel.Pos
		vb.Pose.Quat = wb.Rel.Quat
		UpdatePose(wk, vk)
	}
}
