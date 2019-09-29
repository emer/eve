// Copyright (c) 2019, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package epe

import (
	"github.com/goki/gi/mat32"
	"github.com/goki/ki/ki"
	"github.com/goki/ki/kit"
)

// Phys contains the full specification of a given object's physical properties
// including position, orientation, velocity.
type Phys struct {
	Pos    mat32.Vec3 `desc:"position of center of mass of object"`
	Quat   mat32.Quat `desc:"rotation specified as a Quat"`
	LinVel mat32.Vec3 `desc:"linear velocity"`
	AngVel mat32.Vec3 `desc:"angular velocity"`
	//	RotInertia mat32.Mat3   `desc:"Last calculated rotational inertia matrix in local coords"`
}

var KiT_Phys = kit.Types.AddType(&Phys{}, PhysProps)

// Defaults sets defaults only if current values are nil
func (ps *Phys) Defaults() {
	if ps.Quat.IsNil() {
		ps.Quat.SetIdentity()
	}
}

///////////////////////////////////////////////////////
// 		Moving

// Note: you can just directly add to .Pos too..

// MoveOnAxis moves (translates) the specified distance on the specified local axis,
// relative to the current rotation orientation.
func (ps *Phys) MoveOnAxis(x, y, z, dist float32) {
	ps.Pos.SetAdd(mat32.NewVec3(x, y, z).Normal().MulQuat(ps.Quat).MulScalar(dist))
}

// MoveOnAxisAbs moves (translates) the specified distance on the specified local axis,
// in absolute X,Y,Z coordinates.
func (ps *Phys) MoveOnAxisAbs(x, y, z, dist float32) {
	ps.Pos.SetAdd(mat32.NewVec3(x, y, z).Normal().MulScalar(dist))
}

///////////////////////////////////////////////////////
// 		Rotating

// SetEulerRotation sets the rotation in Euler angles (degrees).
func (ps *Phys) SetEulerRotation(x, y, z float32) {
	ps.Quat.SetFromEuler(mat32.NewVec3(x, y, z).MulScalar(mat32.DegToRadFactor))
}

// SetEulerRotationRad sets the rotation in Euler angles (radians).
func (ps *Phys) SetEulerRotationRad(x, y, z float32) {
	ps.Quat.SetFromEuler(mat32.NewVec3(x, y, z))
}

// EulerRotation returns the current rotation in Euler angles (degrees).
func (ps *Phys) EulerRotation() mat32.Vec3 {
	return mat32.NewEulerAnglesFromQuat(ps.Quat).MulScalar(mat32.RadToDegFactor)
}

// EulerRotationRad returns the current rotation in Euler angles (radians).
func (ps *Phys) EulerRotationRad() mat32.Vec3 {
	return mat32.NewEulerAnglesFromQuat(ps.Quat)
}

// SetAxisRotation sets rotation from local axis and angle in degrees.
func (ps *Phys) SetAxisRotation(x, y, z, angle float32) {
	ps.Quat.SetFromAxisAngle(mat32.NewVec3(x, y, z), mat32.DegToRad(angle))
}

// SetAxisRotationRad sets rotation from local axis and angle in radians.
func (ps *Phys) SetAxisRotationRad(x, y, z, angle float32) {
	ps.Quat.SetFromAxisAngle(mat32.NewVec3(x, y, z), angle)
}

// RotateOnAxis rotates around the specified local axis the specified angle in degrees.
func (ps *Phys) RotateOnAxis(x, y, z, angle float32) {
	ps.Quat.SetMul(mat32.NewQuatAxisAngle(mat32.NewVec3(x, y, z), mat32.DegToRad(angle)))
}

// RotateOnAxisRad rotates around the specified local axis the specified angle in radians.
func (ps *Phys) RotateOnAxisRad(x, y, z, angle float32) {
	ps.Quat.SetMul(mat32.NewQuatAxisAngle(mat32.NewVec3(x, y, z), angle))
}

// RotateEuler rotates by given Euler angles (in degrees) relative to existing rotation.
func (ps *Phys) RotateEuler(x, y, z float32) {
	ps.Quat.SetMul(mat32.NewQuatEuler(mat32.NewVec3(x, y, z).MulScalar(mat32.DegToRadFactor)))
}

// RotateEulerRad rotates by given Euler angles (in radians) relative to existing rotation.
func (ps *Phys) RotateEulerRad(x, y, z, angle float32) {
	ps.Quat.SetMul(mat32.NewQuatEuler(mat32.NewVec3(x, y, z)))
}

// PhysProps define the ToolBar and MenuBar for StructView
var PhysProps = ki.Props{
	"ToolBar": ki.PropSlice{
		{"SetEulerRotation", ki.Props{
			"desc": "Set the local rotation (relative to parent) using Euler angles, in degrees.",
			"icon": "rotate-3d",
			"Args": ki.PropSlice{
				{"Pitch", ki.Props{
					"desc": "rotation up / down along the X axis (in the Y-Z plane), e.g., the altitude (climbing, descending) for motion along the Z depth axis",
				}},
				{"Yaw", ki.Props{
					"desc": "rotation along the Y axis (in the horizontal X-Z plane), e.g., the bearing or direction for motion along the Z depth axis",
				}},
				{"Roll", ki.Props{
					"desc": "rotation along the Z axis (in the X-Y plane), e.g., the bank angle for motion along the Z depth axis",
				}},
			},
		}},
		{"SetAxisRotation", ki.Props{
			"desc": "Set the local rotation (relative to parent) using Axis about which to rotate, and the angle.",
			"icon": "rotate-3d",
			"Args": ki.PropSlice{
				{"X", ki.BlankProp{}},
				{"Y", ki.BlankProp{}},
				{"Z", ki.BlankProp{}},
				{"Angle", ki.BlankProp{}},
			},
		}},
		{"RotateEuler", ki.Props{
			"desc": "rotate (relative to current rotation) using Euler angles, in degrees.",
			"icon": "rotate-3d",
			"Args": ki.PropSlice{
				{"Pitch", ki.Props{
					"desc": "rotation up / down along the X axis (in the Y-Z plane), e.g., the altitude (climbing, descending) for motion along the Z depth axis",
				}},
				{"Yaw", ki.Props{
					"desc": "rotation along the Y axis (in the horizontal X-Z plane), e.g., the bearing or direction for motion along the Z depth axis",
				}},
				{"Roll", ki.Props{
					"desc": "rotation along the Z axis (in the X-Y plane), e.g., the bank angle for motion along the Z depth axis",
				}},
			},
		}},
		{"RotateOnAxis", ki.Props{
			"desc": "Rotate (relative to current rotation) using Axis about which to rotate, and the angle.",
			"icon": "rotate-3d",
			"Args": ki.PropSlice{
				{"X", ki.BlankProp{}},
				{"Y", ki.BlankProp{}},
				{"Z", ki.BlankProp{}},
				{"Angle", ki.BlankProp{}},
			},
		}},
		{"EulerRotation", ki.Props{
			"desc":        "The local rotation (relative to parent) in Euler angles in degrees (X = Pitch, Y = Yaw, Z = Roll)",
			"icon":        "rotate-3d",
			"show-return": "true",
		}},
		{"sep-rot", ki.BlankProp{}},
		{"MoveOnAxis", ki.Props{
			"desc": "Move given distance on given X,Y,Z axis relative to current rotation orientation.",
			"icon": "pan",
			"Args": ki.PropSlice{
				{"X", ki.BlankProp{}},
				{"Y", ki.BlankProp{}},
				{"Z", ki.BlankProp{}},
				{"Dist", ki.BlankProp{}},
			},
		}},
		{"MoveOnAxisAbs", ki.Props{
			"desc": "Move given distance on given X,Y,Z axis in absolute coords, not relative to current rotation orientation.",
			"icon": "pan",
			"Args": ki.PropSlice{
				{"X", ki.BlankProp{}},
				{"Y", ki.BlankProp{}},
				{"Z", ki.BlankProp{}},
				{"Dist", ki.BlankProp{}},
			},
		}},
	},
}