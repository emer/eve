// Copyright (c) 2019, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eve

import (
	"github.com/goki/ki/ki"
	"github.com/goki/ki/kit"
	"github.com/goki/mat32"
)

// Cylinder is a generalized cylinder body shape, with separate radii for top and bottom.
// A cone has a zero radius at one end.
type Cylinder struct {
	BodyBase
	Height float32 `desc:"height of the cylinder"`
	TopRad float32 `desc:"radius of the top -- set to 0 for a cone"`
	BotRad float32 `desc:"radius of the bottom"`
}

var KiT_Cylinder = kit.Types.AddType(&Cylinder{}, CylinderProps)

var CylinderProps = ki.Props{
	"EnumType:Flag": ki.KiT_Flags,
}

// AddNewCylinder adds a new cylinder of given name, initial position
// and height, radius to given parent.
func AddNewCylinder(parent ki.Ki, name string, pos mat32.Vec3, height, radius float32) *Cylinder {
	cy := parent.AddNewChild(KiT_Cylinder, name).(*Cylinder)
	cy.Initial.Pos = pos
	cy.Height = height
	cy.TopRad = radius
	cy.BotRad = radius
	return cy
}

// AddNewCone adds a new cone of given name, initial position
// and height, radius to given parent.
func AddNewCone(parent ki.Ki, name string, pos mat32.Vec3, height, radius float32) *Cylinder {
	cy := parent.AddNewChild(KiT_Cylinder, name).(*Cylinder)
	cy.Initial.Pos = pos
	cy.Height = height
	cy.TopRad = 0
	cy.BotRad = radius
	return cy
}

func (cy *Cylinder) SetBBox() {
	h2 := cy.Height / 2
	cy.BBox.SetBounds(mat32.Vec3{-cy.BotRad, -h2, -cy.BotRad}, mat32.Vec3{cy.TopRad, h2, cy.TopRad})
	cy.BBox.XForm(cy.Abs.Quat, cy.Abs.Pos)
}

func (cy *Cylinder) InitPhys(par *NodeBase) {
	cy.InitBase(par)
	cy.SetBBox()
}

func (cy *Cylinder) UpdatePhys(par *NodeBase) {
	cy.UpdateBase(par)
	cy.SetBBox()
}
