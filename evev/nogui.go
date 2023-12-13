// Copyright (c) 2019, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package evev

import (
	"image"
	"log"

	"github.com/goki/vgpu/vgpu"
	"goki.dev/gi/v2/gi3d"

	vk "github.com/goki/vulkan"
)

// NoDisplayGPU Initializes the Vulkan GPU (vgpu) and returns that
// and the graphics GPU device, with given name, without connecting
// to the display.
func NoDisplayGPU(nm string) (*vgpu.GPU, *vgpu.Device, error) {
	if err := vgpu.InitNoDisplay(); err != nil {
		log.Println(err)
		return nil, nil, err
	}
	// vgpu.Debug = true
	gp := vgpu.NewGPU()
	if err := gp.Config(nm, nil); err != nil {
		log.Println(err)
		return nil, nil, err
	}
	dev := &vgpu.Device{}
	if err := dev.Init(gp, vk.QueueGraphicsBit); err != nil { // todo: add wrapper to vgpu
		log.Println(err)
		return nil, nil, err
	}
	return gp, dev, nil
}

// NoDisplayScene returns a gi3d Scene initialized and ready to use
// in NoGUI offscreen rendering mode, using given GPU and device.
// Must manually call Init3D and Style3D on the Scene prior to
// a RenderOffNode call to grab the image from a specific camera.
func NoDisplayScene(nm string, gp *vgpu.GPU, dev *vgpu.Device) *gi3d.Scene {
	sc := &gi3d.Scene{}
	sc.InitName(sc, "scene")
	sc.Defaults()
	sc.MultiSample = 4
	sc.Geom.Size = image.Point{1024, 768}
	sc.ConfigFrameImpl(gp, dev)
	return sc
}
