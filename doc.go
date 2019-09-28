// Copyright (c) 2019, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package vpe is a virtual physics engine written in pure Go, for use in creating
virtual environments for neural network models to grow up in.

Ultimately we hope to figure out how the Bullet simulator works and get that
running here, in a clean and simple implementation.

Incrementally, we will start with a very basic explicitly driven form of
physics that is sufficient to get started, and build from there.

The world is made from Ki-based trees (groups, bodies, joints),
which can be mapped onto corresponding 3D renders using the gi3d
3D rendering framework.  This is facilitated by using the same
gi3d.Pose structure for representing the basic physical position
and orientation of bodies, which in turn uses the full-featured
mat32 math / matrix library.

The basic physics however can be simulated entirely independent of the
graphics rendering.
*/
package vpe
