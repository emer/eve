# Emergent Virtual Engine

The *Emergent Virtual Engine* (EVE) is a scenegraph-based physics simulator for creating virtual environments for neural network models to grow up in.

Ultimately we hope to figure out how the Bullet simulator works and get that running here, in a clean and simple implementation.

Incrementally, we will start with a very basic explicitly driven form of physics that is sufficient to get started, and build from there.

The world is made using [GoKi](https://github.com/goki/ki) based trees (groups, bodies, joints).

Rendering can *optionally* be performed using corresponding 3D renders in the `gi3d` 3D rendering framework in the [GoGi](https://github.com/goki/gi) GUI framework, using an `epev.View` object that sync's the two.

We also use the full-featured `gi.mat32` math / matrix library (adapted from the `g3n` 3D game environment package).

# Usual rationalization for reinventing the wheel yet again

* Pure *Go* build environment is fast, delightful, simple, clean, easy-to-read, runs fast, etc.  In contrast, building other systems typically requires something like Qt or other gui dependencies, and we know what that world is like, and why we left it for Go..

* Control, control, control.. we can do what we want in a way that we know will work. 

* Physics is easy.  Seriously, the basic stuff is just a few vectors and matricies and a few simple equations.  Doing the full physics version is considerably more challenging technically but we should be able to figure that out in due course, and really most environments we care about don't really require the full physics anyway.


