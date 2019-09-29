// Copyright (c) 2019, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package epev

import (
	"reflect"

	"github.com/emer/epe/epe"
	"github.com/goki/gi/gi3d"
	"github.com/goki/ki/kit"
)

// View connects a Virtual World with a Gi3D visualization thereof
type View struct {
	World *epe.Group  `desc:"the root Group node of the virtual world"`
	Vis   *gi3d.Group `desc:"the root Group node of the visualization"`
}

var KiT_View = kit.Types.AddType(&View{}, nil)

// NewView returns a new View that links given world with given view group
func NewView(world *epe.Group, viewGp *gi3d.Group) *View {
	vw := &View{World: world, Vis: viewGp}
	return vw
}

// Sync synchronizes the view to the world
func (vw *View) Sync(sc *gi3d.Scene) bool {
	return vw.SyncNode(vw.World, vw.Vis, sc)
}

// SyncNode updates the view tree to match the world tree, using
// ConfigChildren to maximally preserve existing tree elements
// returns true if view tree was modified (elements added / removed etc)
func (vw *View) SyncNode(wn epe.Node, vn gi3d.Node3D, sc *gi3d.Scene) bool {
	nm := wn.UniqueName()
	vn.SetNameRaw(nm) // guaranteed to be unique
	vn.SetUniqueName(nm)
	skids := *wn.Children()
	tnl := make(kit.TypeAndNameList, 0, len(skids))
	for _, skid := range skids {
		wt := kit.ShortTypeName(skid.Type())
		var ntyp reflect.Type
		switch wt {
		case "epe.Group":
			ntyp = gi3d.KiT_Group
		// todo could have switch to ignore joints
		default:
			ntyp = gi3d.KiT_Object
		}
		tnl.Add(ntyp, skid.UniqueName())
	}
	mod, updt := vn.ConfigChildren(tnl, false) // false = don't use unique names
	modall := mod
	for idx, _ := range skids {
		wk := wn.Child(idx).(epe.Node)
		vk := vn.Child(idx).(gi3d.Node3D)
		vw.ConfigView(wk, vk, sc)
		kmod := vw.SyncNode(wk, vk, sc)
		if kmod {
			modall = true
		}
	}
	vn.UpdateEnd(updt)
	return modall
}

// ConfigView configures the view node to properly display world node
func (vw *View) ConfigView(wn epe.Node, vn gi3d.Node3D, sc *gi3d.Scene) {
	wb := wn.AsNodeBase()
	vb := vn.AsNode3D()
	vb.Pose.Pos = wb.Cur.Pos
	vb.Pose.Quat = wb.Cur.Quat
	wt := kit.ShortTypeName(wn.Type())
	switch wt {
	case "epe.Box":
		nm := "epeBox"
		bx := wn.(*epe.Box)
		vo := vn.(*gi3d.Object)
		bm := sc.MeshByName(nm)
		if bm == nil {
			bm = gi3d.AddNewBox(sc, nm, 1, 1, 1)
		}
		vo.Mesh = gi3d.MeshName(nm)
		vo.Pose.Scale = bx.Size
		if bx.Mat.Color != "" {
			vo.Mat.Color.SetName(bx.Mat.Color)
		}
		if bx.Mat.Texture != "" {
			vo.Mat.Texture = gi3d.TexName(bx.Mat.Texture) // note: texture must have already been set
		}
	}
}

// todo: UpdatePos ONLY updates positions
