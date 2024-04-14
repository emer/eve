// Copyright (c) 2019, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eve

import (
	"math"

	"cogentcore.org/core/math32"
)

// Phys contains the basic physical properties including position, orientation, velocity.
// These are only the values that can be either relative or absolute -- other physical
// state values such as Mass should go in Rigid.
type Phys struct {

	// position of center of mass of object
	Pos math32.Vec3

	// rotation specified as a Quat
	Quat math32.Quat

	// linear velocity
	LinVel math32.Vec3

	// angular velocity
	AngVel math32.Vec3
}

// Defaults sets defaults only if current values are nil
func (ps *Phys) Defaults() {
	if ps.Quat.IsNil() {
		ps.Quat.SetIdentity()
	}
}

///////////////////////////////////////////////////////
// 	State updates

// FromRel sets state from relative values compared to a parent state
func (ps *Phys) FromRel(rel, par *Phys) {
	ps.Quat = rel.Quat.Mul(par.Quat)
	ps.Pos = rel.Pos.MulQuat(par.Quat).Add(par.Pos)
	ps.LinVel = rel.LinVel.MulQuat(rel.Quat).Add(par.LinVel)
	ps.AngVel = rel.AngVel.MulQuat(rel.Quat).Add(par.AngVel)
}

// AngMotionMax is maximum angular motion that can be taken per update
const AngMotionMax = math.Pi / 4

// StepByAngVel steps the Quat rotation from angular velocity
func (ps *Phys) StepByAngVel(step float32) {
	ang := math32.Sqrt(ps.AngVel.Dot(ps.AngVel))

	// limit the angular motion
	if ang*step > AngMotionMax {
		ang = AngMotionMax / step
	}
	var axis math32.Vec3
	if ang < 0.001 {
		// use Taylor's expansions of sync function
		axis = ps.AngVel.MulScalar(0.5*step - (step*step*step)*0.020833333333*ang*ang)
	} else {
		// sync(fAngle) = sin(c*fAngle)/t
		axis = ps.AngVel.MulScalar(math32.Sin(0.5*ang*step) / ang)
	}
	var dq math32.Quat
	dq.SetFromAxisAngle(axis, ang*step)
	ps.Quat = dq.Mul(ps.Quat)
	ps.Quat.Normalize()
}

// StepByLinVel steps the Pos from the linear velocity
func (ps *Phys) StepByLinVel(step float32) {
	ps.Pos = ps.Pos.Add(ps.LinVel.MulScalar(step))
}

///////////////////////////////////////////////////////
// 		Moving

// Move moves (translates) Pos by given amount, and sets the LinVel to the given
// delta -- this can be useful for Scripted motion to track movement.
func (ps *Phys) Move(delta math32.Vec3) {
	ps.LinVel = delta
	ps.Pos.SetAdd(delta)
}

// MoveOnAxis moves (translates) the specified distance on the specified local axis,
// relative to the current rotation orientation.
// The axis is normalized prior to aplying the distance factor.
// Sets the LinVel to motion vector.
func (ps *Phys) MoveOnAxis(x, y, z, dist float32) {
	ps.LinVel = math32.V3(x, y, z).Normal().MulQuat(ps.Quat).MulScalar(dist)
	ps.Pos.SetAdd(ps.LinVel)
}

// MoveOnAxisAbs moves (translates) the specified distance on the specified local axis,
// in absolute X,Y,Z coordinates (does not apply the Quat rotation factor.
// The axis is normalized prior to aplying the distance factor.
// Sets the LinVel to motion vector.
func (ps *Phys) MoveOnAxisAbs(x, y, z, dist float32) {
	ps.LinVel = math32.V3(x, y, z).Normal().MulScalar(dist)
	ps.Pos.SetAdd(ps.LinVel)
}

///////////////////////////////////////////////////////
// 		Rotating

// SetEulerRotation sets the rotation in Euler angles (degrees).
func (ps *Phys) SetEulerRotation(x, y, z float32) {
	ps.Quat.SetFromEuler(math32.V3(x, y, z).MulScalar(math32.DegToRadFactor))
}

// SetEulerRotationRad sets the rotation in Euler angles (radians).
func (ps *Phys) SetEulerRotationRad(x, y, z float32) {
	ps.Quat.SetFromEuler(math32.V3(x, y, z))
}

// EulerRotation returns the current rotation in Euler angles (degrees).
func (ps *Phys) EulerRotation() math32.Vec3 {
	return ps.Quat.ToEuler().MulScalar(math32.RadToDegFactor)
}

// EulerRotationRad returns the current rotation in Euler angles (radians).
func (ps *Phys) EulerRotationRad() math32.Vec3 {
	return ps.Quat.ToEuler()
}

// SetAxisRotation sets rotation from local axis and angle in degrees.
func (ps *Phys) SetAxisRotation(x, y, z, angle float32) {
	ps.Quat.SetFromAxisAngle(math32.V3(x, y, z), math32.DegToRad(angle))
}

// SetAxisRotationRad sets rotation from local axis and angle in radians.
func (ps *Phys) SetAxisRotationRad(x, y, z, angle float32) {
	ps.Quat.SetFromAxisAngle(math32.V3(x, y, z), angle)
}

// RotateOnAxis rotates around the specified local axis the specified angle in degrees.
func (ps *Phys) RotateOnAxis(x, y, z, angle float32) {
	ps.Quat.SetMul(math32.NewQuatAxisAngle(math32.V3(x, y, z), math32.DegToRad(angle)))
}

// RotateOnAxisRad rotates around the specified local axis the specified angle in radians.
func (ps *Phys) RotateOnAxisRad(x, y, z, angle float32) {
	ps.Quat.SetMul(math32.NewQuatAxisAngle(math32.V3(x, y, z), angle))
}

// RotateEuler rotates by given Euler angles (in degrees) relative to existing rotation.
func (ps *Phys) RotateEuler(x, y, z float32) {
	ps.Quat.SetMul(math32.NewQuatEuler(math32.V3(x, y, z).MulScalar(math32.DegToRadFactor)))
}

// RotateEulerRad rotates by given Euler angles (in radians) relative to existing rotation.
func (ps *Phys) RotateEulerRad(x, y, z, angle float32) {
	ps.Quat.SetMul(math32.NewQuatEuler(math32.V3(x, y, z)))
}

/*

// PhysProps define the ToolBar and MenuBar for StructView
var PhysProps = tree.Props{
	"ToolBar": tree.PropSlice{
		{"SetEulerRotation", tree.Props{
			"desc": "Set the local rotation (relative to parent) using Euler angles, in degrees.",
			"icon": "rotate-3d",
			"Args": tree.PropSlice{
				{"Pitch", tree.Props{
					"desc": "rotation up / down along the X axis (in the Y-Z plane), e.g., the altitude (climbing, descending) for motion along the Z depth axis",
				}},
				{"Yaw", tree.Props{
					"desc": "rotation along the Y axis (in the horizontal X-Z plane), e.g., the bearing or direction for motion along the Z depth axis",
				}},
				{"Roll", tree.Props{
					"desc": "rotation along the Z axis (in the X-Y plane), e.g., the bank angle for motion along the Z depth axis",
				}},
			},
		}},
		{"SetAxisRotation", tree.Props{
			"desc": "Set the local rotation (relative to parent) using Axis about which to rotate, and the angle.",
			"icon": "rotate-3d",
			"Args": tree.PropSlice{
				{"X", tree.BlankProp{}},
				{"Y", tree.BlankProp{}},
				{"Z", tree.BlankProp{}},
				{"Angle", tree.BlankProp{}},
			},
		}},
		{"RotateEuler", tree.Props{
			"desc": "rotate (relative to current rotation) using Euler angles, in degrees.",
			"icon": "rotate-3d",
			"Args": tree.PropSlice{
				{"Pitch", tree.Props{
					"desc": "rotation up / down along the X axis (in the Y-Z plane), e.g., the altitude (climbing, descending) for motion along the Z depth axis",
				}},
				{"Yaw", tree.Props{
					"desc": "rotation along the Y axis (in the horizontal X-Z plane), e.g., the bearing or direction for motion along the Z depth axis",
				}},
				{"Roll", tree.Props{
					"desc": "rotation along the Z axis (in the X-Y plane), e.g., the bank angle for motion along the Z depth axis",
				}},
			},
		}},
		{"RotateOnAxis", tree.Props{
			"desc": "Rotate (relative to current rotation) using Axis about which to rotate, and the angle.",
			"icon": "rotate-3d",
			"Args": tree.PropSlice{
				{"X", tree.BlankProp{}},
				{"Y", tree.BlankProp{}},
				{"Z", tree.BlankProp{}},
				{"Angle", tree.BlankProp{}},
			},
		}},
		{"EulerRotation", tree.Props{
			"desc":        "The local rotation (relative to parent) in Euler angles in degrees (X = Pitch, Y = Yaw, Z = Roll)",
			"icon":        "rotate-3d",
			"show-return": "true",
		}},
		{"sep-rot", tree.BlankProp{}},
		{"MoveOnAxis", tree.Props{
			"desc": "Move given distance on given X,Y,Z axis relative to current rotation orientation.",
			"icon": "pan",
			"Args": tree.PropSlice{
				{"X", tree.BlankProp{}},
				{"Y", tree.BlankProp{}},
				{"Z", tree.BlankProp{}},
				{"Dist", tree.BlankProp{}},
			},
		}},
		{"MoveOnAxisAbs", tree.Props{
			"desc": "Move given distance on given X,Y,Z axis in absolute coords, not relative to current rotation orientation.",
			"icon": "pan",
			"Args": tree.PropSlice{
				{"X", tree.BlankProp{}},
				{"Y", tree.BlankProp{}},
				{"Z", tree.BlankProp{}},
				{"Dist", tree.BlankProp{}},
			},
		}},
	},
}

*/
