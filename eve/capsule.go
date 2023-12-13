// Copyright (c) 2019, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eve

import (
	"goki.dev/mat32/v2"
)

// Capsule is a generalized cylinder body shape, with hemispheres at each end,
// with separate radii for top and bottom.
type Capsule struct {
	BodyBase

	// height of the cylinder portion of the capsule
	Height float32 `desc:"height of the cylinder portion of the capsule"`

	// radius of the top hemisphere
	TopRad float32 `desc:"radius of the top hemisphere"`

	// radius of the bottom hemisphere
	BotRad float32 `desc:"radius of the bottom hemisphere"`
}

func (cp *Capsule) SetBBox() {
	th := cp.Height + cp.TopRad + cp.BotRad
	h2 := th / 2
	cp.BBox.SetBounds(mat32.Vec3{-cp.BotRad, -h2, -cp.BotRad}, mat32.Vec3{cp.TopRad, h2, cp.TopRad})
	cp.BBox.XForm(cp.Abs.Quat, cp.Abs.Pos)
}

func (cp *Capsule) InitAbs(par *NodeBase) {
	cp.InitAbsBase(par)
	cp.SetBBox()
	cp.BBox.VelNilProject()
}

func (cp *Capsule) RelToAbs(par *NodeBase) {
	cp.RelToAbsBase(par)
	cp.SetBBox()
	cp.BBox.VelProject(cp.Abs.LinVel, 1)
}

func (cp *Capsule) StepPhys(step float32) {
	cp.StepPhysBase(step)
	cp.SetBBox()
	cp.BBox.VelProject(cp.Abs.LinVel, step)
}
