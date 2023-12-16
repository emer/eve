// Code generated by "goki generate -add-types"; DO NOT EDIT.

package evev

import (
	"goki.dev/gti"
	"goki.dev/ordmap"
)

var _ = gti.AddType(&gti.Type{
	Name:       "github.com/emer/eve/v2/evev.Camera",
	ShortName:  "evev.Camera",
	IDName:     "camera",
	Doc:        "Camera defines the properties of a camera needed for offscreen rendering",
	Directives: gti.Directives{},
	Fields: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"Size", &gti.Field{Name: "Size", Type: "image.Point", LocalType: "image.Point", Doc: "size of image to record", Directives: gti.Directives{}, Tag: ""}},
		{"FOV", &gti.Field{Name: "FOV", Type: "float32", LocalType: "float32", Doc: "field of view in degrees", Directives: gti.Directives{}, Tag: ""}},
		{"Near", &gti.Field{Name: "Near", Type: "float32", LocalType: "float32", Doc: "near plane z coordinate", Directives: gti.Directives{}, Tag: "def:\"0.01\""}},
		{"Far", &gti.Field{Name: "Far", Type: "float32", LocalType: "float32", Doc: "far plane z coordinate", Directives: gti.Directives{}, Tag: "def:\"1000\""}},
		{"MaxD", &gti.Field{Name: "MaxD", Type: "float32", LocalType: "float32", Doc: "maximum distance for depth maps -- anything above is 1 -- this is independent of Near / Far rendering (though must be < Far) and is for normalized depth maps", Directives: gti.Directives{}, Tag: "def:\"20\""}},
		{"LogD", &gti.Field{Name: "LogD", Type: "bool", LocalType: "bool", Doc: "use the natural log of 1 + depth for normalized depth values in display etc", Directives: gti.Directives{}, Tag: "def:\"true\""}},
		{"MSample", &gti.Field{Name: "MSample", Type: "int", LocalType: "int", Doc: "number of multi-samples to use for antialising -- 4 is best and default", Directives: gti.Directives{}, Tag: "def:\"4\""}},
		{"UpDir", &gti.Field{Name: "UpDir", Type: "goki.dev/mat32/v2.Vec3", LocalType: "mat32.Vec3", Doc: "up direction for camera -- which way is up -- defaults to positive Y axis, and is reset by call to LookAt method", Directives: gti.Directives{}, Tag: ""}},
	}),
	Embeds:  ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}),
	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
})

var _ = gti.AddType(&gti.Type{
	Name:       "github.com/emer/eve/v2/evev.View",
	ShortName:  "evev.View",
	IDName:     "view",
	Doc:        "View connects a Virtual World with a Xyz Scene to visualize the world,\nincluding ability to render offscreen",
	Directives: gti.Directives{},
	Fields: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"World", &gti.Field{Name: "World", Type: "*github.com/emer/eve/v2/eve.Group", LocalType: "*eve.Group", Doc: "the root Group node of the virtual world", Directives: gti.Directives{}, Tag: ""}},
		{"Scene", &gti.Field{Name: "Scene", Type: "*goki.dev/xyz.Scene", LocalType: "*xyz.Scene", Doc: "the scene object for visualizing", Directives: gti.Directives{}, Tag: ""}},
		{"Root", &gti.Field{Name: "Root", Type: "*goki.dev/xyz.Group", LocalType: "*xyz.Group", Doc: "the root Group node in the Scene under which the world is rendered", Directives: gti.Directives{}, Tag: ""}},
	}),
	Embeds:  ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}),
	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
})
