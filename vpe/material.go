// Copyright (c) 2019, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vpe

// Material defines basic material properties (color, texture)
// that can be used for rendering -- for more complete control you
// will need to update rendering node Mat properties directly
type Material struct {
	Color   string
	Texture string
}

// note: could add a map[string]string for other properties
