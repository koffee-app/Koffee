package controllers

// Group structure that represents a whole group of routes
type Group struct {
	Prefix string
}

// Route Returns the path of the group + the specified path
func (g *Group) Route(s string) string {
	return g.Prefix + s
}

// New Returns a new Group with the appended string as the prefix
func New(g *Group, s string) *Group {
	return &Group{Prefix: g.Route(s)}
}
