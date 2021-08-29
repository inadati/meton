package endpoint

import "fmt"

type NodeRecipe struct {}

var Node = &NodeRecipe{}

func (r *NodeRecipe) SSH(user string, host string, port int) string {
	if port == 0 {
		port = 22
	}
	return fmt.Sprintf("ssh://%s@%s:%d", user, host, port)
}
