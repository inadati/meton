package docker

import (
	"fmt"
	"net/http"

	"github.com/docker/cli/cli/connhelper"
	"github.com/docker/docker/client"
)

func Client(sshAddr string) (*client.Client, error) {
	helper, err := connhelper.GetConnectionHelper(sshAddr)
	if err != nil {
		return nil, fmt.Errorf("GetConnectionHelper fail")
	}

	return client.NewClientWithOpts(
		client.WithAPIVersionNegotiation(),
		client.WithHTTPClient(&http.Client{
			Transport: &http.Transport{
				DialContext: helper.Dialer,
			},
		}),
		client.WithHost(helper.Host),
		client.WithDialContext(helper.Dialer),
	)
}
