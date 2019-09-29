// Copyright (c) 2019, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package epe

import (
	"github.com/goki/ki/ki"
	"github.com/goki/ki/kit"
)

// Group is a container of bodies, joints, or other groups
// it should be used strategically to partition the space
// and its BBox is used to optimize tree-based collision detection.
// Use a group for the top-level World node as well.
type Group struct {
	NodeBase
}

var KiT_Group = kit.Types.AddType(&Group{}, GroupProps)

// AddNewGroup adds a new group of given name to given parent
func AddNewGroup(parent ki.Ki, name string) *Group {
	gp := parent.AddNewChild(KiT_Group, name).(*Group)
	return gp
}

func (gp *Group) NodeType() NodeTypes {
	return GROUP
}

func (gp *Group) InitPhys(par *NodeBase) {
	gp.InitBase(par)
}

func (gp *Group) GroupBBox() {
	gp.BBox.BBox.SetEmpty()
	for _, kid := range gp.Kids {
		nii, ni := KiToNode(kid)
		if nii == nil {
			continue
		}
		gp.BBox.BBox.ExpandByPoint(ni.BBox.BBox.Min)
		gp.BBox.BBox.ExpandByPoint(ni.BBox.BBox.Max)
	}
}

// InitWorld does the full tree InitPhys and GroupBBox updates
func (gp *Group) InitWorld() {
	gp.FuncDownMeFirst(0, gp.This(), func(k ki.Ki, level int, d interface{}) bool {
		nii, _ := KiToNode(k)
		if nii == nil {
			return false // going into a different type of thing, bail
		}
		_, pi := KiToNode(k.Parent())
		nii.InitPhys(pi)
		return true
	})

	gp.FuncDownMeLast(0, gp.This(),
		func(k ki.Ki, level int, d interface{}) bool {
			nii, _ := KiToNode(k)
			if nii == nil {
				return false // going into a different type of thing, bail
			}
			return true
		},
		func(k ki.Ki, level int, d interface{}) bool {
			nii, _ := KiToNode(k)
			if nii == nil {
				return false // going into a different type of thing, bail
			}
			nii.GroupBBox()
			return true
		})

}

// GroupProps define the ToolBar and MenuBar for StructView
var GroupProps = ki.Props{
	"ToolBar": ki.PropSlice{
		{"InitWorld", ki.Props{
			"desc": "initialize all elements in the world.",
			"icon": "reset",
		}},
	},
}
