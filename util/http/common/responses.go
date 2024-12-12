package common

// InfoResponse is the response to 'GET /info'
//
// swagger:response InfoResponse
type InfoResponse struct {
	// in: body
	Body Info
}

// GetError allows InfoResponse to satisfy the APIV1Response interface, even
// though it can never return an error and is not versioned
func (r InfoResponse) GetError() error {
	return nil
}

// BuildVersion defines model for BuildVersion.
type BuildVersion struct {
	// Branch Branch the build is based on
	Branch string `json:"branch"`

	// BuildNumber Gorxestra's minor version number
	BuildNumber int `json:"build_number"`

	// Channel Branch the build is based on
	Channel string `json:"channel"`

	// CommitHash Hash of commit the build is based on
	CommitHash string `json:"commit_hash"`

	// Major Gorxestra's major version number
	Major int `json:"major"`

	// Minor Gorxestra's minor version number
	Minor int `json:"minor"`
}

// Error defines model for Error.
type Error struct {
	// Error Error message
	Error string `json:"error"`
}

// Info defines model for Info.
type Info struct {
	Build BuildVersion `json:"build"`

	// Versions returns a list of supported protocol versions ( i.e. v1, v2 etc. )
	Versions []string `json:"versions"`
}
