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
	"github.com/goki/mat32"
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

// InitLibrary initializes Scene library with basic shapes
// based on bodies in the virtual world.  More complex visualizations
// can be configured after this.
func (vw *View) InitLibrary() {
	vw.InitLibraryBody(vw.World)
}

// Sync synchronizes the view to the world
func (vw *View) Sync() bool {
	rval := vw.SyncNode(vw.World, vw.Root)
	return rval
}

// UpdatePose updates the view pose values only from world tree.
// Essential that both trees are already synchronized.
func (vw *View) UpdatePose() {
	vw.UpdatePoseNode(vw.World, vw.Root)
	vw.Scene.UpdateSig()
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

// NewInLibrary adds a new item of given name in library
func (vw *View) NewInLibrary(nm string) *svg.Group {
	if vw.Library == nil {
		vw.Library = make(map[string]*svg.Group)
	}
	gp := &svg.Group{}
	gp.InitName(gp, nm)
	vw.Library[nm] = gp
	return gp
}

// AddFmLibrary adds shape from library to given group
func (vw *View) AddFmLibrary(nm string, gp *svg.Group) {
	lgp, has := vw.Library[nm]
	if !has {
		return
	}
	gp.AddChild(lgp.Clone())
}

// InitLibraryBody initializes Scene library with basic shapes
// based on bodies in the virtual world.  More complex visualizations
// can be configured after this.
func (vw *View) InitLibraryBody(wn eve.Node) {
	bod := wn.AsBody()
	if bod != nil {
		vw.InitLibShape(bod)
	}
	for idx := range *wn.Children() {
		wk := wn.Child(idx).(eve.Node)
		vw.InitLibraryBody(wk)
	}
}

// InitLibShape initializes Scene library with basic shape for given body
func (vw *View) InitLibShape(bod eve.Body) {
	nm := bod.Name()
	bb := bod.AsBodyBase()
	if bb.Vis == "" {
		bb.Vis = nm
	}
	if _, has := vw.Library[nm]; has {
		return
	}
	lgp := vw.NewInLibrary(nm)
	wt := kit.ShortTypeName(ki.Type(bod.This()))
	switch wt {
	case "eve.Box":
		mnm := "eveBox"
		svg.AddNewRect(lgp, mnm, 0, 0, 1, 1)
	case "eve.Cylinder":
		mnm := "eveCylinder"
		svg.AddNewCircle(lgp, mnm, 0, 0, 1)
	case "eve.Capsule":
		mnm := "eveCapsule"
		svg.AddNewCircle(lgp, mnm, 0, 0, 1)
	case "eve.Sphere":
		mnm := "eveSphere"
		svg.AddNewCircle(lgp, mnm, 0, 0, 1)
	}
}

func NewVec2Fm3(v3 mat32.Vec3) mat32.Vec2 {
	return mat32.NewVec2(v3.X, v3.Y)
}

// ConfigBodyShape configures a shape for a body with current values
func (vw *View) ConfigBodyShape(bod eve.Body, shp svg.NodeSVG) {
	wt := kit.ShortTypeName(ki.Type(bod.This()))
	sb := shp.AsSVGNode()
	switch wt {
	case "eve.Box":
		bx := bod.(*eve.Box)
		shp.SetSize(NewVec2Fm3(bx.Size))
		sb.Pnt.XForm = mat32.Translate2D(-bx.Size.X/2, -bx.Size.Y/2)
		shp.SetProp("transform", sb.Pnt.XForm.String())
		shp.SetProp("stroke-width", 0.05)
		shp.SetProp("fill", "none")
		if bx.Color != "" {
			shp.SetProp("stroke", bx.Color)
		}
	case "eve.Cylinder":
		cy := bod.(*eve.Cylinder)
		shp.SetSize(mat32.NewVec2(cy.BotRad*2, cy.BotRad*2))
		sb.Pnt.XForm = mat32.Translate2D(-cy.BotRad, -cy.BotRad)
		shp.SetProp("transform", sb.Pnt.XForm.String())
		shp.SetProp("stroke-width", 0.05)
		shp.SetProp("fill", "none")
		if cy.Color != "" {
			shp.SetProp("stroke", cy.Color)
		}
	case "eve.Capsule":
		cp := bod.(*eve.Capsule)
		shp.SetSize(mat32.NewVec2(cp.BotRad*2, cp.BotRad*2))
		sb.Pnt.XForm = mat32.Translate2D(-cp.BotRad, -cp.BotRad)
		shp.SetProp("transform", sb.Pnt.XForm.String())
		shp.SetProp("stroke-width", 0.05)
		shp.SetProp("fill", "none")
		if cp.Color != "" {
			shp.SetProp("stroke", cp.Color)
		}
	case "eve.Sphere":
		sp := bod.(*eve.Sphere)
		shp.SetSize(mat32.NewVec2(sp.Radius*2, sp.Radius*2))
		sb.Pnt.XForm = mat32.Translate2D(-sp.Radius, -sp.Radius)
		shp.SetProp("transform", sb.Pnt.XForm.String())
		shp.SetProp("stroke-width", 0.05)
		shp.SetProp("fill", "none")
		if sp.Color != "" {
			shp.SetProp("stroke", sp.Color)
		}
	}
}

// ConfigView configures the view node to properly display world node
func (vw *View) ConfigView(wn eve.Node, vn svg.NodeSVG) {
	wb := wn.AsNodeBase()
	vb := vn.(*svg.Group)
	ps := NewVec2Fm3(wb.Rel.Pos)
	vb.Pnt.XForm = mat32.Translate2D(ps.X, ps.Y)
	vb.SetProp("transform", vb.Pnt.XForm.String())
	// fmt.Printf("wb: %s  pos: %v  vb: %s\n", wb.Name(), ps, vb.Name())
	bod := wn.AsBody()
	if bod == nil {
		return
	}
	if !vb.HasChildren() {
		vw.AddFmLibrary(bod.AsBodyBase().Vis, vb)
	}
	bgp := vb.Child(0)
	if bgp.HasChildren() {
		shp, has := bgp.Child(0).(svg.NodeSVG)
		if has {
			vw.ConfigBodyShape(bod, shp)
		}
	}
}

// SyncNode updates the view tree to match the world tree, using
// ConfigChildren to maximally preserve existing tree elements
// returns true if view tree was modified (elements added / removed etc)
func (vw *View) SyncNode(wn eve.Node, vn svg.NodeSVG) bool {
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
		vk := vn.Child(idx).(svg.NodeSVG)
		vw.ConfigView(wk, vk)
		if wk.HasChildren() {
			kmod := vw.SyncNode(wk, vk)
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
func (vw *View) UpdatePoseNode(wn eve.Node, vn svg.NodeSVG) {
	skids := *wn.Children()
	for idx := range skids {
		wk := wn.Child(idx).(eve.Node)
		vk := vn.Child(idx).(svg.NodeSVG).(*svg.Group)
		wb := wk.AsNodeBase()
		ps := NewVec2Fm3(wb.Rel.Pos)
		vk.Pnt.XForm = mat32.Translate2D(ps.X, ps.Y)
		vk.SetProp("transform", vk.Pnt.XForm.String())
		// fmt.Printf("wk: %s  pos: %v  vk: %s\n", wk.Name(), ps, vk.Child(0).Name())
		vw.UpdatePoseNode(wk, vk)
	}
}
