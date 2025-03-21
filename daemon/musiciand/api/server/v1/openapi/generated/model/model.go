// Package model provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package model

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

// InfoResponse defines model for InfoResponse.
type InfoResponse struct {
	Body Info `json:"body"`
}

// MusicNote defines model for MusicNote.
type MusicNote struct {
	// Note base64 encoded note
	Note string `json:"note"`
}

// PlayJSONRequestBody defines body for Play for application/json ContentType.
type PlayJSONRequestBody = MusicNote
