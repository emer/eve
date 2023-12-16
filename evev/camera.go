// Copyright (c) 2019, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package evev

import (
	"image"

	"goki.dev/mat32/v2"
)

// Camera defines the properties of a camera needed for offscreen rendering
type Camera struct {

	// size of image to record
	Size image.Point

	// field of view in degrees
	FOV float32

	// near plane z coordinate
	Near float32 `def:"0.01"`

	// far plane z coordinate
	Far float32 `def:"1000"`

	// maximum distance for depth maps -- anything above is 1 -- this is independent of Near / Far rendering (though must be < Far) and is for normalized depth maps
	MaxD float32 `def:"20"`

	// use the natural log of 1 + depth for normalized depth values in display etc
	LogD bool `def:"true"`

	// number of multi-samples to use for antialising -- 4 is best and default
	MSample int `def:"4"`

	// up direction for camera -- which way is up -- defaults to positive Y axis, and is reset by call to LookAt method
	UpDir mat32.Vec3
}

func (cm *Camera) Defaults() {
	cm.Size = image.Point{320, 180}
	cm.FOV = 30
	cm.Near = .01
	cm.Far = 1000
	cm.MaxD = 20
	cm.LogD = true
	cm.MSample = 4
	cm.UpDir = mat32.Vec3Y
}
