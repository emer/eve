// Copyright (c) 2019, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package evev

import (
	"image"
	"log"

	"github.com/goki/gi/gi3d"
	"github.com/goki/vgpu/vgpu"

	vk "github.com/goki/vulkan"
)

// NoGUIScene returns a gi3d Scene initialized and ready to use
// in NoGUI offscreen rendering mode.  Initializes the GPU
// and returns that and the graphics GPU device.
// Must manually call Init3D and Style3D on the Scene prior to
// a RenderOffNode call to grab the image from a specific camera.
func NoGUIScene(nm string) (*gi3d.Scene, *vgpu.GPU, *vgpu.Device, error) {
	if err := vgpu.InitNoDisplay(); err != nil {
		log.Println(err)
		return nil, nil, nil, err
	}
	// vgpu.Debug = true
	gp := vgpu.NewGPU()
	if err := gp.Config(nm, nil); err != nil {
		log.Println(err)
		return nil, nil, nil, err
	}
	dev := &vgpu.Device{}
	if err := dev.Init(gp, vk.QueueGraphicsBit); err != nil { // todo: add wrapper to vgpu
		log.Println(err)
		return nil, nil, nil, err
	}

	sc := &gi3d.Scene{}
	sc.InitName(sc, "scene")
	sc.Defaults()
	sc.MultiSample = 4
	sc.Geom.Size = image.Point{1024, 768}
	sc.ConfigFrameImpl(gp, dev)
	return sc, gp, dev, nil
}
