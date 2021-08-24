package endpoint

import "fmt"

type node struct {
	SSH func(string, string, int) string
}

var Node = &node{
	SSH: func(user string, host string, port int) string {
		if port == 0 {
			port = 22
		}
		return fmt.Sprintf("ssh://%s@%s:%d", user, host, port)
	},
}