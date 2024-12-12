package v1

import (
	"encoding/base64"

	"crossjoin.com/gorxestra/daemon/musiciand/api"
	"crossjoin.com/gorxestra/daemon/musiciand/api/server/v1/openapi/generated/model"
	"crossjoin.com/gorxestra/logging"
	"github.com/labstack/echo/v4"
)

// Handlers is an implementation to the V1 route handler interface
type Handlers struct {
	Node api.NodeInterface
	Log  logging.Logger
}

func (h *Handlers) Play(ctx echo.Context) error {
	var a model.MusicNote

	ctx.Bind(&a)
	h.Log.With("note", a.Note).Info("received musical note")
	bs, err := base64.RawStdEncoding.DecodeString(a.Note)
	if err != nil {
		h.Log.With("error", err).Error("decoding note")
		return err
	}

	return h.Node.Play(bs)
}
