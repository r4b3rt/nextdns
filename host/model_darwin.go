package host

import (
	"bytes"
	"os/exec"
)

func Model() string {
	if b, err := exec.Command("sysctl", "-n", "hw.model").Output(); err == nil && len(b) > 0 {
		return "Apple " + string(bytes.TrimSpace(b))
	}
	return ""
}
