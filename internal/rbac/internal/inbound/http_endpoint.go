package inbound

import (
	"encoding/json"

	"github.com/shandysiswandi/gostarter/internal/rbac/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/framework"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
)

var errInvalidBody = goerror.NewInvalidFormat("invalid request body")

type httpEndpoint struct {
	telemetry *telemetry.Telemetry

	createRoleUC       domain.CreateRole
	createPermissionUC domain.CreatePermission
}

func (h *httpEndpoint) CreateRole(c framework.Context) (any, error) {
	ctx, span := h.telemetry.Tracer().Start(c.Context(), "rbac.inbound.httpEndpoint.CreateRole")
	defer span.End()

	var req CreateRoleRequest
	if err := json.NewDecoder(c.Body()).Decode(&req); err != nil {
		return nil, errInvalidBody
	}

	resp, err := h.createRoleUC.Call(ctx, domain.CreateRoleInput{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}

	return CreateRoleResponse{ID: resp.ID}, nil
}

func (h *httpEndpoint) CreatePermission(c framework.Context) (any, error) {
	ctx, span := h.telemetry.Tracer().Start(c.Context(), "rbac.inbound.httpEndpoint.CreatePermission")
	defer span.End()

	var req CreatePermissionRequest
	if err := json.NewDecoder(c.Body()).Decode(&req); err != nil {
		return nil, errInvalidBody
	}

	resp, err := h.createPermissionUC.Call(ctx, domain.CreatePermissionInput{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}

	return CreatePermissionResponse{ID: resp.ID}, nil
}
