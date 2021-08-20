package address

import (
	"fmt"
	"unsafe"

	"github.com/meton888/meton/config"
)

type Compound struct {
	Servers   string
	Zookeeper string
}

func (c *Compound) Init(masterNodes []config.Node) {
	srvAddrBytes := make([]byte, 0, 128)
	zkAddrBytes := make([]byte, 0, 128)

	for i, node := range masterNodes {
		srvAddrBytes = append(srvAddrBytes, node.Address.External...)
		zkAddrBytes = append(zkAddrBytes, node.Address.External...)
		zkAddrBytes = append(zkAddrBytes, ":2181"...)
		if i < len(masterNodes)-1 {
			srvAddrBytes = append(srvAddrBytes, ","...)
			zkAddrBytes = append(zkAddrBytes, ","...)
		}
	}

	c.Servers = *(*string)(unsafe.Pointer(&srvAddrBytes))
	c.Zookeeper = fmt.Sprintf("zk://%s", *(*string)(unsafe.Pointer(&zkAddrBytes)))
}

func SSH(user string, host string, port int) string {
	if port == 0 {
		port = 22
	}
	return fmt.Sprintf("ssh://%s@%s:%d", user, host, port)
}
