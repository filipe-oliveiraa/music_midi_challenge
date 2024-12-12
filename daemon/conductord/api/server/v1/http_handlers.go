package v1

import (
	"net/http"

	"crossjoin.com/gorxestra/daemon/conductord/api"
	"crossjoin.com/gorxestra/daemon/conductord/api/server/v1/openapi/generated/model"
	"crossjoin.com/gorxestra/data"
	"crossjoin.com/gorxestra/logging"
	"github.com/labstack/echo/v4"
)

// Handlers is an implementation to the V1 route handler interface
type Handlers struct {
	Node api.NodeInterface
	Log  logging.Logger
}

// AddMusician implements server.ServerInterface.
func (h *Handlers) RegisterMusician(ctx echo.Context) error {
	var musicianDto model.Musician
	err := ctx.Bind(&musicianDto)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	id, err := data.IdFromHex(musicianDto.Id)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	err = h.Node.RegisterMusician(data.Musician{
		Id:      id,
		Address: musicianDto.Address,
	})
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	return ctx.JSON(http.StatusOK, nil)
}

// DeleteMusician implements server.ServerInterface.
func (h *Handlers) UnregisterMusician(ctx echo.Context, idRaw string) error {
	id, err := data.IdFromHex(idRaw)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	err = h.Node.UnregisterMusician(id)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, nil)
}

// PlayMusic implements server.ServerInterface.
func (h *Handlers) PlayMusic(ctx echo.Context, name string) error {
	err := h.Node.PlayMusic(name)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, nil)
}
