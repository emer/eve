// Copyright (c) 2019, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package epe

import (
	"github.com/goki/ki/ki"
	"github.com/goki/ki/kit"
)

// Node is the common interface for all epe nodes
type Node interface {
	ki.Ki

	// NodeType returns the type of node this is (Body, Group, Joint)
	NodeType() NodeTypes

	// AsNodeBase returns a generic NodeBase for our node -- gives generic
	// access to all the base-level data structures without needing interface methods.
	AsNodeBase() *NodeBase

	// AsBody returns a generic Body interface for our node -- nil if not a Body
	AsBody() Body

	// InitPhys sets current physical state parameters from Initial values
	// which are local, relative to parent -- is passed the parent (nil = top).
	// Body nodes should also set their bounding boxes.
	// called in a FuncDownMeFirst traversal.
	InitPhys(par *NodeBase)

	// GroupBBox sets bounding boxes for groups based on groups or bodies.
	// called in a FuncDownMeLast traversal.
	GroupBBox()

	// UpdatePhys updates current world physical state parameters based on propagated
	// updates from higher levels -- is passed the parent (nil = top).
	// Body nodes should also update their bounding boxes.
	// called in a FuncDownMeFirst traversal.
	UpdatePhys(par *NodeBase)
}

// NodeBase is the basic epe node, which has position, rotation, velocity
// and computed bounding boxes, etc.
// There are only three different kinds of Nodes: Group, Body, and Joint
type NodeBase struct {
	ki.Node
	Initial Phys    `view:"inline" desc:"initial position, orientation, velocity in *local* coordinates (relative to parent)"`
	Rel     Phys    `view:"inline" desc:"current relative (local) position, orientation, velocity -- only change these values, as abs values are computed therefrom"`
	Abs     Phys    `inactive:"+" view:"inline" desc:"current absolute (world) position, orientation, velocity"`
	Mass    float32 `desc:"mass of body or aggregate mass of group of bodies (just fyi for groups) -- 0 mass = no dynamics"`
	BBox    BBox    `desc:"bounding box in world coordinates (aggregated for groups)"`
}

func (nb *NodeBase) AsNodeBase() *NodeBase {
	return nb
}

func (nb *NodeBase) AsBody() Body {
	return nil
}

// InitBase is the base-level initialization of basic Phys state from Initial conditions
func (nb *NodeBase) InitBase(par *NodeBase) {
	if nb.Initial.Quat.IsNil() {
		nb.Initial.Quat.SetIdentity()
	}
	nb.Rel = nb.Initial
	if par != nil {
		nb.Abs.Pos = nb.Initial.Pos.MulQuat(nb.Initial.Quat).Add(par.Abs.Pos)
		nb.Abs.LinVel = nb.Initial.LinVel.MulQuat(nb.Initial.Quat).Add(par.Abs.LinVel)
		nb.Abs.AngVel = nb.Initial.AngVel.MulQuat(nb.Initial.Quat).Add(par.Abs.AngVel)
		nb.Abs.Quat = nb.Initial.Quat.Mul(par.Abs.Quat)
	} else {
		nb.Abs.Pos = nb.Initial.Pos
		nb.Abs.LinVel = nb.Initial.LinVel
		nb.Abs.AngVel = nb.Initial.AngVel
		nb.Abs.Quat = nb.Initial.Quat
	}
}

// UpdateBase is the base-level update of Phys state based on current relative values
func (nb *NodeBase) UpdateBase(par *NodeBase) {
	if par != nil {
		nb.Abs.Pos = nb.Rel.Pos.MulQuat(nb.Rel.Quat).Add(par.Abs.Pos)
		nb.Abs.LinVel = nb.Rel.LinVel.MulQuat(nb.Rel.Quat).Add(par.Abs.LinVel)
		nb.Abs.AngVel = nb.Rel.AngVel.MulQuat(nb.Rel.Quat).Add(par.Abs.AngVel)
		nb.Abs.Quat = nb.Rel.Quat.Mul(par.Abs.Quat)
	} else {
		nb.Abs.Pos = nb.Rel.Pos
		nb.Abs.LinVel = nb.Rel.LinVel
		nb.Abs.AngVel = nb.Rel.AngVel
		nb.Abs.Quat = nb.Rel.Quat
	}
}

// KiToNode converts Ki to a Node interface and a Node3DBase obj -- nil if not.
func KiToNode(k ki.Ki) (Node, *NodeBase) {
	if k == nil || k.This() == nil { // this also checks for destroyed
		return nil, nil
	}
	nii, ok := k.(Node)
	if ok {
		return nii, nii.AsNodeBase()
	}
	return nil, nil
}

/////////////////////////////////////////////////////////////////////
// NodeTypes

// NodeTypes is a list of node types
type NodeTypes int

const (
	// note: uppercase required to not conflict with type names
	BODY NodeTypes = iota
	GROUP
	JOINT
	NodeTypesN
)

//go:generate stringer -type=NodeTypes

var KiT_NodeTypes = kit.Enums.AddEnum(NodeTypesN, false, nil)
