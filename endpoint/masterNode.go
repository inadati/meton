package endpoint

import (
	"fmt"
	"unsafe"

	"github.com/meton888/meton/config"
)

type masterNode struct {
	AddrCollection func(masterNodes []config.Node) (string, string)
}

var MasterNode = &masterNode{
	AddrCollection: func(masterNodes []config.Node) (string, string) {
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
		return *(*string)(unsafe.Pointer(&srvAddrBytes)), fmt.Sprintf("zk://%s", *(*string)(unsafe.Pointer(&zkAddrBytes)))
	},
}
