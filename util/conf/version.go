package conf

import (
	"fmt"
	"strconv"
)

type Version struct {
	// Major version number
	Major int

	// Minor version number
	Minor int

	// Build Number
	BuildNumber int

	// Hash of commit the build is based on
	CommitHash string

	// Branch the build is based on
	Branch string

	// Branch-derived release channel the build is based on
	Channel string
}

var currentVersion = Version{
	Major:       convertToInt(VersionMajor),
	Minor:       convertToInt(VersionMinor),
	BuildNumber: convertToInt(BuildNumber),
	CommitHash:  CommitHash,
	Branch:      Branch,
	Channel:     Channel,
}

// GetCurrentVersion retrieves a copy of the current global Version structure (for the application)
func GetCurrentVersion() Version {
	return currentVersion
}

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d commitHash:%s branch:%s channel:%s\n",
		v.Major,
		v.Minor,
		v.BuildNumber,
		v.CommitHash,
		v.Branch,
		v.Channel,
	)
}

func convertToInt(val string) int {
	if val == "" {
		return 0
	}
	value, _ := strconv.ParseInt(val, 10, 0)
	return int(value)
}
