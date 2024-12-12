package conf

// Do NOT remove or rename these constants - they are inspected by build tools
// to generate the build tag and update package name
/* Build time variables set through -ldflags */

var (
	// VersionMajor is the Major semantic version number (#.y.z)
	// changed when first public release (0.y.z -> 1.y.z)
	// and when backwards compatibility is broken.
	VersionMajor string

	// VersionMinor is the Minor semantic version number (x.#.z)
	// changed when backwards-compatible features are introduced.
	// Not enforced until after initial public release (x > 0).
	VersionMinor string

	// BuildNumber is the monotonic build number, currently based on the date and hour-of-day.
	// It will be set to a build number by the build tools
	BuildNumber string

	// CommitHash is the git commit id in effect when the build was created.
	// It will be set by the build tools
	CommitHash string

	// Branch is the git branch in effect when the build was created.
	// It will be set by the build tools
	Branch string

	// Channel is the computed release channel based on the Branch in effect when the build was created.
	// It will be set by the build tools
	Channel string
)
