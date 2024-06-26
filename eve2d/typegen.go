// Code generated by "core generate -add-types"; DO NOT EDIT.

package eve2d

import (
	"cogentcore.org/core/types"
)

var _ = types.AddType(&types.Type{Name: "github.com/emer/eve/v2/eve2d.View", IDName: "view", Doc: "View connects a Virtual World with a 2D SVG Scene to visualize the world", Fields: []types.Field{{Name: "LineWidth", Doc: "width of lines for shape rendering, in normalized units"}, {Name: "Prjn", Doc: "projection matrix for converting 3D to 2D -- resulting X, Y coordinates are used from Vector3"}, {Name: "World", Doc: "the root Group node of the virtual world"}, {Name: "Scene", Doc: "the SVG rendering canvas for visualizing in 2D"}, {Name: "Root", Doc: "the root Group node in the Scene under which the world is rendered"}, {Name: "Library", Doc: "library of shapes for bodies -- name matches Body.Vis"}}})
