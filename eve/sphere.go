// Copyright (c) 2019, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eve

import (
	"github.com/goki/ki/ki"
	"github.com/goki/ki/kit"
	"github.com/goki/mat32"
)

// Sphere is a spherical body shape.
type Sphere struct {
	BodyBase
	Radius float32 `desc:"radius"`
}

var KiT_Sphere = kit.Types.AddType(&Sphere{}, SphereProps)

var SphereProps = ki.Props{
	"EnumType:Flag": KiT_NodeFlags,
}

// AddNewSphere adds a new sphere of given name, initial position
// and radius.
func AddNewSphere(parent ki.Ki, name string, pos mat32.Vec3, radius float32) *Sphere {
	sp := parent.AddNewChild(KiT_Sphere, name).(*Sphere)
	sp.Initial.Pos = pos
	sp.Radius = radius
	return sp
}

func (sp *Sphere) SetBBox() {
	sp.BBox.SetBounds(mat32.Vec3{-sp.Radius, -sp.Radius, -sp.Radius}, mat32.Vec3{sp.Radius, sp.Radius, sp.Radius})
	sp.BBox.XForm(sp.Abs.Quat, sp.Abs.Pos)
}

func (sp *Sphere) InitAbs(par *NodeBase) {
	sp.InitAbsBase(par)
	sp.SetBBox()
	sp.BBox.VelNilProject()
}

func (sp *Sphere) RelToAbs(par *NodeBase) {
	sp.RelToAbsBase(par)
	sp.SetBBox()
	sp.BBox.VelProject(sp.Abs.LinVel, 1)
}

func (sp *Sphere) StepPhys(step float32) {
	sp.StepPhysBase(step)
	sp.SetBBox()
	sp.BBox.VelProject(sp.Abs.LinVel, step)
}
