// Copyright (c) 2019, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eve

//go:generate goki generate -add-types

import (
	"goki.dev/ki/v2"
	"goki.dev/mat32/v2"
)

// Node is the common interface for all eve nodes
type Node interface {
	ki.Ki

	// NodeType returns the type of node this is (Body, Group, Joint)
	NodeType() NodeTypes

	// AsNodeBase returns a generic NodeBase for our node -- gives generic
	// access to all the base-level data structures without needing interface methods.
	AsNodeBase() *NodeBase

	// AsBody returns a generic Body interface for our node -- nil if not a Body
	AsBody() Body

	// IsDynamic returns true if node has Dynamic flag set -- otherwise static
	// Groups that contain dynamic objects set their dynamic flags.
	IsDynamic() bool

	// GroupBBox sets bounding boxes for groups based on groups or bodies.
	// called in a FuncDownMeLast traversal.
	GroupBBox()

	// InitAbs sets current Abs physical state parameters from Initial values
	// which are local, relative to parent -- is passed the parent (nil = top).
	// Body nodes should also set their bounding boxes.
	// Called in a FuncDownMeFirst traversal.
	InitAbs(par *NodeBase)

	// RelToAbs updates current world Abs physical state parameters
	// based on Rel values added to updated Abs values at higher levels.
	// Abs.LinVel is updated from the resulting change from prior position.
	// This is useful for manual updating of relative positions (scripted movement).
	// It is passed the parent (nil = top).
	// Body nodes should also update their bounding boxes.
	// Called in a FuncDownMeFirst traversal.
	RelToAbs(par *NodeBase)

	// StepPhys computes one update of the world Abs physical state parameters,
	// using *current* velocities -- add forces prior to calling.
	// Use this for physics-based state updates.
	// Body nodes should also update their bounding boxes.
	StepPhys(step float32)
}

// NodeBase is the basic eve node, which has position, rotation, velocity
// and computed bounding boxes, etc.
// There are only three different kinds of Nodes: Group, Body, and Joint
type NodeBase struct {
	ki.Node

	// initial position, orientation, velocity in *local* coordinates (relative to parent)
	Initial Phys `view:"inline"`

	// current relative (local) position, orientation, velocity -- only change these values, as abs values are computed therefrom
	Rel Phys `view:"inline"`

	// current absolute (world) position, orientation, velocity
	Abs Phys `inactive:"+" view:"inline"`

	// bounding box in world coordinates (aggregated for groups)
	BBox BBox
}

func (nb *NodeBase) AsNodeBase() *NodeBase {
	return nb
}

func (nb *NodeBase) AsBody() Body {
	return nil
}

func (nb *NodeBase) IsDynamic() bool {
	return nb.Is(Dynamic)
}

// SetInitPos sets the initial position
func (nb *NodeBase) SetInitPos(pos mat32.Vec3) *NodeBase {
	nb.Initial.Pos = pos
	return nb
}

// SetInitQuat sets the initial rotation as a Quaternion
func (nb *NodeBase) SetInitQuat(quat mat32.Quat) *NodeBase {
	nb.Initial.Quat = quat
	return nb
}

// SetInitLinVel sets the initial linear velocity
func (nb *NodeBase) SetInitLinVel(vel mat32.Vec3) *NodeBase {
	nb.Initial.LinVel = vel
	return nb
}

// SetInitAngVel sets the initial angular velocity
func (nb *NodeBase) SetInitAngVel(vel mat32.Vec3) *NodeBase {
	nb.Initial.AngVel = vel
	return nb
}

// InitAbsBase is the base-level version of InitAbs -- most nodes call this.
// InitAbs sets current Abs physical state parameters from Initial values
// which are local, relative to parent -- is passed the parent (nil = top).
// Body nodes should also set their bounding boxes.
// Called in a FuncDownMeFirst traversal.
func (nb *NodeBase) InitAbsBase(par *NodeBase) {
	if nb.Initial.Quat.IsNil() {
		nb.Initial.Quat.SetIdentity()
	}
	nb.Rel = nb.Initial
	if par != nil {
		nb.Abs.FromRel(&nb.Initial, &par.Abs)
	} else {
		nb.Abs = nb.Initial
	}
}

// RelToAbsBase is the base-level version of RelToAbs -- most nodes call this.
// note: Group WorldRelToAbs ensures only called on Dynamic nodes.
// RelToAbs updates current world Abs physical state parameters
// based on Rel values added to updated Abs values at higher levels.
// Abs.LinVel is updated from the resulting change from prior position.
// This is useful for manual updating of relative positions (scripted movement).
// It is passed the parent (nil = top).
// Body nodes should also update their bounding boxes.
// Called in a FuncDownMeFirst traversal.
func (nb *NodeBase) RelToAbsBase(par *NodeBase) {
	ppos := nb.Abs.Pos
	if par != nil {
		nb.Abs.FromRel(&nb.Rel, &par.Abs)
	} else {
		nb.Abs = nb.Rel
	}
	nb.Abs.LinVel = nb.Abs.Pos.Sub(ppos) // needed for VelBBox prjn
}

// StepPhysBase is base-level version of StepPhys -- most nodes call this.
// note: Group WorldRelToAbs ensures only called on Dynamic nodes.
// Computes one update of the world Abs physical state parameters,
// using *current* velocities -- add forces prior to calling.
// Use this for physics-based state updates.
// Body nodes should also update their bounding boxes.
func (nb *NodeBase) StepPhysBase(step float32) {
	nb.Abs.StepByAngVel(step)
	nb.Abs.StepByLinVel(step)
}

// AsNode converts Ki to a Node interface and a Node3DBase obj -- nil if not.
func AsNode(k ki.Ki) (Node, *NodeBase) {
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
type NodeTypes int32 //enums:enum

const (
	// note: uppercase required to not conflict with type names
	BODY NodeTypes = iota
	GROUP
	JOINT
)

//////////////////////////////////////////////////////////////////////
// NodeFlags

// NodeFlags define eve node bitflags -- uses ki Flags field (64 bit capacity)
type NodeFlags ki.Flags //enums:bitflag

const (
	// Dynamic means that this node can move -- if not so marked, it is
	// a Static node.  Any top-level group that is not Dynamic is immediately
	// pruned from further consideration, so top-level groups should be
	// separated into Dynamic and Static nodes at the start.
	Dynamic NodeFlags = NodeFlags(ki.FlagsN) + iota
)
