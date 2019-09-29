// Copyright (c) 2019, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vpe

import (
	"github.com/goki/gi/mat32"
	"github.com/goki/ki/ki"
	"github.com/goki/ki/kit"
)

// Box is a rigid body box shape
type Box struct {
	BodyBase
	Size mat32.Vec3 `desc:"size of box in each dimension"`
}

var KiT_Box = kit.Types.AddType(&Box{}, nil)

// AddNewBox adds a new box of given name, initial position and size to given parent
func AddNewBox(parent ki.Ki, name string, pos, size mat32.Vec3) *Box {
	bx := parent.AddNewChild(KiT_Box, name).(*Box)
	bx.Initial.Pos = pos
	bx.Size = size
	return bx
}

func (bx *Box) InitPhys(par *NodeBase) {
	bx.InitBase(par)
	bx.BBox.SetBounds(bx.Size.MulScalar(-.5), bx.Size.MulScalar(.5))
	bx.BBox.XForm(bx.Cur.Quat, bx.Cur.Pos)
}
