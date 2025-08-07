//go:build e2e

package e2e_test

import (
	"fmt"
	"strings"
)

func makeGoMimicEndpoint(port uint16, path ...string) string {
	base := fmt.Sprintf("http://0.0.0.0:%d", port)
	return base + "/" + strings.Join(path, "/")
}
