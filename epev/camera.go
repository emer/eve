// Copyright (c) 2019, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package epev

import (
	"image"

	"github.com/goki/gi/mat32"
)

// Camera defines the properties of a camera needed for offscreen rendering
type Camera struct {
	Size    image.Point `desc:"size of image to record"`
	FOV     float32     `desc:"field of view in degrees"`
	Near    float32     `def:"0.01" desc:"near plane z coordinate"`
	Far     float32     `def:"1000" desc:"far plane z coordinate"`
	MSample int         `def:"4" desc:"number of multi-samples to use for antialising -- 4 is best and default"`
	UpDir   mat32.Vec3  `desc:"up direction for camera -- which way is up -- defaults to positive Y axis, and is reset by call to LookAt method"`
}

func (cm *Camera) Defaults() {
	cm.Size = image.Point{320, 180}
	cm.FOV = 30
	cm.Near = .01
	cm.Far = 1000
	cm.MSample = 4
	cm.UpDir = mat32.Vec3Y
}
