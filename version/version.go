package version

import (
	"fmt"
)

const (
	GitVersion = "Unknown"

	AppVersion = "0.0.1"

	AppName = "my-app"

	HumanName = "MY APP"
)

func Version() string {
	return fmt.Sprintf("git=%s , app=%s", GitVersion, AppVersion)
}
