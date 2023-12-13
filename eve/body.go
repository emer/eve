// Copyright (c) 2019, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eve

// Body is the common interface for all body types
type Body interface {
	Node

	// AsBodyBase returns the body as a BodyBase
	AsBodyBase() *BodyBase

	// SetDynamic sets the Dynamic flag for this body, indicating that it moves.
	// It is important to collect all dynamic objects into separate top-level group(s)
	// for more efficiently organizing the collision detection process.
	SetDynamic()
}

// BodyBase is the base type for all specific Body types
type BodyBase struct {
	NodeBase

	// rigid body properties, including mass, bounce, friction etc
	Rigid Rigid `desc:"rigid body properties, including mass, bounce, friction etc"`

	// visualization name -- looks up an entry in the scene library that provides the visual representation of this body
	Vis string `desc:"visualization name -- looks up an entry in the scene library that provides the visual representation of this body"`

	// default color of body for basic InitLibrary configuration
	Color string `desc:"default color of body for basic InitLibrary configuration"`
}

func (bb *BodyBase) NodeType() NodeTypes {
	return BODY
}

func (bb *BodyBase) AsBody() Body {
	return bb.This().(Body)
}

func (bb *BodyBase) AsBodyBase() *BodyBase {
	return bb
}

func (bb *BodyBase) GroupBBox() {

}

func (bb *BodyBase) SetDynamic() {
	bb.SetFlag(true, Dynamic)
}
