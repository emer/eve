// Code generated by "goki generate -add-types"; DO NOT EDIT.

package eve2d

import (
	"goki.dev/gti"
	"goki.dev/ordmap"
)

var _ = gti.AddType(&gti.Type{
	Name:       "github.com/emer/eve/v2/eve2d.View",
	ShortName:  "eve2d.View",
	IDName:     "view",
	Doc:        "View connects a Virtual World with a 2D SVG Scene to visualize the world",
	Directives: gti.Directives{},
	Fields: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"LineWidth", &gti.Field{Name: "LineWidth", Type: "float32", LocalType: "float32", Doc: "width of lines for shape rendering, in normalized units", Directives: gti.Directives{}, Tag: ""}},
		{"Prjn", &gti.Field{Name: "Prjn", Type: "goki.dev/mat32/v2.Mat4", LocalType: "mat32.Mat4", Doc: "projection matrix for converting 3D to 2D -- resulting X, Y coordinates are used from Vec3", Directives: gti.Directives{}, Tag: ""}},
		{"World", &gti.Field{Name: "World", Type: "*github.com/emer/eve/v2/eve.Group", LocalType: "*eve.Group", Doc: "the root Group node of the virtual world", Directives: gti.Directives{}, Tag: ""}},
		{"Scene", &gti.Field{Name: "Scene", Type: "*goki.dev/svg.SVG", LocalType: "*svg.SVG", Doc: "the SVG rendering canvas for visualizing in 2D", Directives: gti.Directives{}, Tag: ""}},
		{"Root", &gti.Field{Name: "Root", Type: "*goki.dev/svg.Group", LocalType: "*svg.Group", Doc: "the root Group node in the Scene under which the world is rendered", Directives: gti.Directives{}, Tag: ""}},
		{"Library", &gti.Field{Name: "Library", Type: "map[string]*goki.dev/svg.Group", LocalType: "map[string]*svg.Group", Doc: "library of shapes for bodies -- name matches Body.Vis", Directives: gti.Directives{}, Tag: ""}},
	}),
	Embeds:  ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}),
	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
})
