// Copyright (c) 2019, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package epev

import (
	"fmt"
	"reflect"

	"github.com/emer/epe/epe"
	"github.com/goki/gi/gi3d"
	"github.com/goki/gi/oswin/gpu"
	"github.com/goki/ki/kit"
)

// View connects a Virtual World with a Gi3D Scene to visualize the world,
// including ability to render offscreen
type View struct {
	World *epe.Group  `desc:"the root Group node of the virtual world"`
	Scene *gi3d.Scene `desc:"the scene object for visualizing"`
	Root  *gi3d.Group `desc:"the root Group node in the Scene where the world is rendered under"`
}

var KiT_View = kit.Types.AddType(&View{}, nil)

// NewView returns a new View that links given world with given scene and root group
func NewView(world *epe.Group, sc *gi3d.Scene, root *gi3d.Group) *View {
	vw := &View{World: world, Scene: sc, Root: root}
	return vw
}

// Sync synchronizes the view to the world
func (vw *View) Sync() bool {
	return SyncNode(vw.World, vw.Root, vw.Scene)
}

// RenderOffNode does an offscreen render using given node for the camera position and
// orientation, and given parameters for field-of-view (in degrees, e.g., 30)
// near plane (.01 default) and far plane (1000 default).
// and multisampling number (4 = default for good antialiasing, 0 if not hardware accelerated).
// Current scene camera is saved and restored
func (vw *View) RenderOffNode(frame *gpu.Framebuffer, node epe.Node, cam *Camera) error {
	sc := vw.Scene
	camnm := "epe-view-renderoff-save"
	sc.SaveCamera(camnm)
	defer sc.SetCamera(camnm)
	sc.Camera.FOV = cam.FOV
	sc.Camera.Near = cam.Near
	sc.Camera.Far = cam.Far
	nb := node.AsNodeBase()
	sc.Camera.Pose.Pos = nb.Abs.Pos
	sc.Camera.Pose.Quat = nb.Abs.Quat
	sc.Camera.Pose.Scale.Set(1, 1, 1)
	if !sc.ActivateOffFrame(frame, "epe-view", cam.Size, cam.MSample) {
		return fmt.Errorf("could not activate offscreen framebuffer")
	}
	if !sc.RenderOffFrame() {
		return fmt.Errorf("could not render to offscreen framebuffer")
	}
	return nil
}

///////////////////////////////////////////////////////////////
// Sync, Config, etc

// SyncNode updates the view tree to match the world tree, using
// ConfigChildren to maximally preserve existing tree elements
// returns true if view tree was modified (elements added / removed etc)
func SyncNode(wn epe.Node, vn gi3d.Node3D, sc *gi3d.Scene) bool {
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
		ConfigView(wk, vk, sc)
		kmod := SyncNode(wk, vk, sc)
		if kmod {
			modall = true
		}
	}
	vn.UpdateEnd(updt)
	return modall
}

// ConfigView configures the view node to properly display world node
func ConfigView(wn epe.Node, vn gi3d.Node3D, sc *gi3d.Scene) {
	wb := wn.AsNodeBase()
	vb := vn.AsNode3D()
	vb.Pose.Pos = wb.Rel.Pos
	vb.Pose.Quat = wb.Rel.Quat
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
