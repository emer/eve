// Copyright (c) 2019, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eve

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

func (gp *Group) InitAbs(par *NodeBase) {
	gp.InitAbsBase(par)
}

func (gp *Group) RelToAbs(par *NodeBase) {
	gp.RelToAbsBase(par) // yes we can move groups
}

func (gp *Group) StepPhys(step float32) {
	// groups do NOT update physics
}

func (gp *Group) GroupBBox() {
	hasDyn := false
	gp.BBox.BBox.SetEmpty()
	gp.BBox.VelBBox.SetEmpty()
	for _, kid := range gp.Kids {
		nii, ni := KiToNode(kid)
		if nii == nil {
			continue
		}
		gp.BBox.BBox.ExpandByBox(ni.BBox.BBox)
		gp.BBox.VelBBox.ExpandByBox(ni.BBox.VelBBox)
		if nii.IsDynamic() {
			hasDyn = true
		}
	}
	gp.SetFlagState(hasDyn, int(Dynamic))
}

// WorldDynGroupBBox does a GroupBBox on all dynamic nodes
func (gp *Group) WorldDynGroupBBox() {
	gp.FuncDownMeLast(0, gp.This(),
		func(k ki.Ki, level int, d interface{}) bool {
			nii, _ := KiToNode(k)
			if nii == nil {
				return false // going into a different type of thing, bail
			}
			if !nii.IsDynamic() {
				return false
			}
			return true
		},
		func(k ki.Ki, level int, d interface{}) bool {
			nii, _ := KiToNode(k)
			if nii == nil {
				return false // going into a different type of thing, bail
			}
			if !nii.IsDynamic() {
				return false
			}
			nii.GroupBBox()
			return true
		})
}

// WorldInit does the full tree InitAbs and GroupBBox updates
func (gp *Group) WorldInit() {
	gp.FuncDownMeFirst(0, gp.This(), func(k ki.Ki, level int, d interface{}) bool {
		nii, _ := KiToNode(k)
		if nii == nil {
			return false // going into a different type of thing, bail
		}
		_, pi := KiToNode(k.Parent())
		nii.InitAbs(pi)
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

// WorldRelToAbs does a full RelToAbs update for all Dynamic groups, for
// Scripted mode updates with manual updating of Rel values.
func (gp *Group) WorldRelToAbs() {
	gp.FuncDownMeFirst(0, gp.This(), func(k ki.Ki, level int, d interface{}) bool {
		nii, _ := KiToNode(k)
		if nii == nil {
			return false // going into a different type of thing, bail
		}
		if !nii.IsDynamic() {
			return false
		}
		_, pi := KiToNode(k.Parent())
		nii.RelToAbs(pi)
		return true
	})

	gp.WorldDynGroupBBox()
}

// WorldStepPhys does a full StepPhys update for all Dynamic nodes, for
// either physics or scripted mode, based on current velocities.
func (gp *Group) WorldStepPhys(step float32) {
	gp.FuncDownMeFirst(0, gp.This(), func(k ki.Ki, level int, d interface{}) bool {
		nii, _ := KiToNode(k)
		if nii == nil {
			return false // going into a different type of thing, bail
		}
		if !nii.IsDynamic() {
			return false
		}
		nii.StepPhys(step)
		return true
	})

	gp.WorldDynGroupBBox()
}

const (
	// DynsTopGps is passed to WorldCollide when all dynamic objects are in separate top groups
	DynsTopGps = true

	// DynsSubGps is passed to WorldCollide when all dynamic objects are in separate groups under top
	// level (i.e., one level deeper)
	DynsSubGps
)

// WorldCollide does first pass filtering step of collision detection
// based on separate dynamic vs. dynamic and dynamic vs. static groups.
// If dynTop is true, then each Dynamic group is separate at the top level --
// otherwise they are organized at the next group level.
// Contacts are organized by dynamic group, when non-nil, for easier
// processing.
func (gp *Group) WorldCollide(dynTop bool) []Contacts {
	var stats []Node
	var dyns []Node
	for _, kid := range gp.Kids {
		nii, _ := KiToNode(kid)
		if nii == nil {
			continue
		}
		if nii.IsDynamic() {
			dyns = append(dyns, nii)
		} else {
			stats = append(stats, nii)
		}
	}

	var sdyns []Node
	if !dynTop {
		for _, d := range dyns {
			for _, dk := range *d.Children() {
				nii, _ := KiToNode(dk)
				if nii == nil {
					continue
				}
				sdyns = append(sdyns, nii)
			}
		}
		dyns = sdyns
	}

	var cts []Contacts
	for i, d := range dyns {
		var dct Contacts
		for _, s := range stats {
			cc := BodyVelBBoxIntersects(d, s)
			dct = append(dct, cc...)
		}
		for di := 0; di < i; di++ {
			od := dyns[di]
			cc := BodyVelBBoxIntersects(d, od)
			dct = append(dct, cc...)
		}
		if len(dct) > 0 {
			cts = append(cts, dct)
		}
	}
	return cts
}

// GroupProps define the ToolBar and MenuBar for StructView
var GroupProps = ki.Props{
	"EnumType:Flag": KiT_NodeFlags,
	"ToolBar": ki.PropSlice{
		{"WorldInit", ki.Props{
			"desc": "initialize all elements in the world.",
			"icon": "reset",
		}},
	},
}
