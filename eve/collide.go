// Copyright (c) 2019, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eve

import (
	"github.com/goki/ki/ki"
	"github.com/goki/mat32"
)

// Contact is one pairwise point of contact between two bodies.
// Contacts are represented in spherical terms relative to the
// spherical BBox of A and B.
type Contact struct {

	// one body
	A Body `desc:"one body"`

	// the other body
	B Body `desc:"the other body"`

	// normal pointing from center of B to center of A
	NormB mat32.Vec3 `desc:"normal pointing from center of B to center of A"`

	// point on spherical shell of B where A is contacting
	PtB mat32.Vec3 `desc:"point on spherical shell of B where A is contacting"`

	// distance from PtB along NormB to contact point on spherical shell of A
	Dist float32 `desc:"distance from PtB along NormB to contact point on spherical shell of A"`
}

// UpdtDist updates the distance information for the contact
func (c *Contact) UpdtDist() {

}

// Contacts is a slice list of contacts
type Contacts []*Contact

// AddNew adds a new contact to the list
func (cs *Contacts) AddNew(a, b Body) *Contact {
	c := &Contact{A: a, B: b}
	*cs = append(*cs, c)
	return c
}

// BodyVelBBoxIntersects returns the list of potential contact nodes between a and b
// (could be the same or different groups) that have intersecting velocity-projected
// bounding boxes.  In general a should be dynamic bodies and b either dynamic or static.
// This is the broad first-pass filtering.
func BodyVelBBoxIntersects(a, b Node) Contacts {
	var cts Contacts
	a.FuncDownMeFirst(0, a.This(), func(k ki.Ki, level int, d interface{}) bool {
		aii, ai := KiToNode(k)
		if aii == nil {
			return false // going into a different type of thing, bail
		}
		if aii.NodeType() != BODY {
			return true
		}
		abod := aii.AsBody() // only consider bodies from a

		b.FuncDownMeFirst(0, b.This(), func(k ki.Ki, level int, d interface{}) bool {
			bii, bi := KiToNode(k)
			if bii == nil {
				return false // going into a different type of thing, bail
			}
			if !ai.BBox.IntersectsVelBox(&bi.BBox) {
				return false // done
			}
			if bii.NodeType() == BODY {
				cts.AddNew(abod, bii.AsBody())
				return false // done
			}
			return true // keep going
		})

		return false
	})
	return cts
}
