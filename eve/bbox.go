// Copyright (c) 2019, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eve

import "cogentcore.org/core/mat32"

// BBox contains bounding box and other gross object properties
type BBox struct {

	// bounding box in world coords (Axis-Aligned Bounding Box = AABB)
	BBox mat32.Box3

	// velocity-projected bounding box in world coords: extend BBox to include future position of moving bodies -- collision must be made on this basis
	VelBBox mat32.Box3

	// bounding sphere in local coords
	BSphere mat32.Sphere

	// area
	Area float32

	// volume
	Volume float32
}

// SetBounds sets BBox from min, max and updates other factors based on that
func (bb *BBox) SetBounds(min, max mat32.Vec3) {
	bb.BBox.Set(&min, &max)
	bb.UpdateFmBBox()
}

// UpdateFmBBox updates other values from BBox
func (bb *BBox) UpdateFmBBox() {
	bb.BSphere.SetFromBox(bb.BBox)
	sz := bb.BBox.Size()
	bb.Area = 2*sz.X + 2*sz.Y + 2*sz.Z
	bb.Volume = sz.X * sz.Y * sz.Z
}

// XForm transforms bounds with given quat and position offset to convert to world coords
func (bb *BBox) XForm(q mat32.Quat, pos mat32.Vec3) {
	bb.BBox = bb.BBox.MulQuat(q).Translate(pos)
	bb.BSphere.Translate(pos)
}

// VelProject computes the velocity-projected bounding box for given velocity and step size
func (bb *BBox) VelProject(vel mat32.Vec3, step float32) {
	eb := bb.BBox.Translate(vel.MulScalar(step))
	bb.VelBBox = bb.BBox
	bb.VelBBox.ExpandByBox(eb)
}

// VelNilProject is for static items -- just copy the BBox
func (bb *BBox) VelNilProject() {
	bb.VelBBox = bb.BBox
}

// IntersectsVelBox returns true if two velocity-projected bounding boxes intersect
func (bb *BBox) IntersectsVelBox(oth *BBox) bool {
	return bb.VelBBox.IntersectsBox(oth.VelBBox)
}
