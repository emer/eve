// Copyright (c) 2019, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eve

// Surface defines the physical surface properties of bodies
type Surface struct {
	Friction float32 `desc:"coulomb friction coefficient (mu). 0 = frictionless, 1e22 = infinity = no slipping"`
	Bounce   float32 `desc:"(0-1) how bouncy is the surface (0 = hard, 1 = maximum bouncyness)"`
}
