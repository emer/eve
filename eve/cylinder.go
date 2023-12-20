// Copyright (c) 2019, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eve

import (
	"goki.dev/ki/v2"
	"goki.dev/mat32/v2"
)

// Cylinder is a generalized cylinder body shape, with separate radii for top and bottom.
// A cone has a zero radius at one end.
type Cylinder struct {
	BodyBase

	// height of the cylinder
	Height float32

	// radius of the top -- set to 0 for a cone
	TopRad float32

	// radius of the bottom
	BotRad float32
}

// NewCone adds a new cone of given name, initial position
// and height, radius to given parent.
func NewCone(parent ki.Ki, name string, pos mat32.Vec3, height, radius float32) *Cylinder {
	cy := parent.NewChild(BoxType, name).(*Cylinder)
	cy.Initial.Pos = pos
	cy.Height = height
	cy.TopRad = 0
	cy.BotRad = radius
	return cy
}

func (cy *Cylinder) SetBBox() {
	h2 := cy.Height / 2
	cy.BBox.SetBounds(mat32.V3(-cy.BotRad, -h2, -cy.BotRad), mat32.V3(cy.TopRad, h2, cy.TopRad))
	cy.BBox.XForm(cy.Abs.Quat, cy.Abs.Pos)
}

func (cy *Cylinder) InitAbs(par *NodeBase) {
	cy.InitAbsBase(par)
	cy.SetBBox()
	cy.BBox.VelNilProject()
}

func (cy *Cylinder) RelToAbs(par *NodeBase) {
	cy.RelToAbsBase(par)
	cy.SetBBox()
	cy.BBox.VelProject(cy.Abs.LinVel, 1)
}

func (cy *Cylinder) StepPhys(step float32) {
	cy.StepPhysBase(step)
	cy.SetBBox()
	cy.BBox.VelProject(cy.Abs.LinVel, step)
}
