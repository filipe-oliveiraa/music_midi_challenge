package openapi

import _ "embed" // for embedding purposes

// SwaggerSpecYAMLEmbed is a string that is pulled from oapi.yaml via go-embed
// for use with the GET /swagger endpoint
//
//go:generate mkdir -p generated
//go:generate mkdir -p generated/model
//go:generate mkdir -p generated/server
//go:generate oapi-codegen -config ./models.cfg.yaml oapi.yaml
//go:generate oapi-codegen --templates ../../../../../../util/templates -config ./server.cfg.yaml oapi.yaml
//go:embed oapi.yaml
var SwaggerSpecYAMLEmbed string
