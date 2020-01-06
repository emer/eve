// Copyright (c) 2019, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eve

import (
	"github.com/goki/gi/mat32"
	"github.com/goki/ki/ki"
	"github.com/goki/ki/kit"
)

// Capsule is a generalized cylinder body shape, with hemispheres at each end,
// with separate radii for top and bottom.
type Capsule struct {
	BodyBase
	Height float32 `desc:"height of the cylinder portion of the capsule"`
	TopRad float32 `desc:"radius of the top hemisphere"`
	BotRad float32 `desc:"radius of the bottom hemisphere"`
}

var KiT_Capsule = kit.Types.AddType(&Capsule{}, CapsuleProps)

var CapsuleProps = ki.Props{
	"EnumType:Flag": ki.KiT_Flags,
}

// AddNewCapsule adds a new capsule of given name, initial position
// and height, radius to given parent.
func AddNewCapsule(parent ki.Ki, name string, pos mat32.Vec3, height, radius float32) *Capsule {
	cp := parent.AddNewChild(KiT_Capsule, name).(*Capsule)
	cp.Initial.Pos = pos
	cp.Height = height
	cp.TopRad = radius
	cp.BotRad = radius
	return cp
}

func (cp *Capsule) SetBBox() {
	th := cp.Height + cp.TopRad + cp.BotRad
	h2 := th / 2
	cp.BBox.SetBounds(mat32.Vec3{-cp.BotRad, -h2, -cp.BotRad}, mat32.Vec3{cp.TopRad, h2, cp.TopRad})
	cp.BBox.XForm(cp.Abs.Quat, cp.Abs.Pos)
}

func (cp *Capsule) InitPhys(par *NodeBase) {
	cp.InitBase(par)
	cp.SetBBox()
}

func (cp *Capsule) UpdatePhys(par *NodeBase) {
	cp.UpdateBase(par)
	cp.SetBBox()
}
