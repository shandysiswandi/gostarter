package inbound

import (
	"encoding/json"
	"strconv"

	"github.com/shandysiswandi/gostarter/internal/rbac/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/framework"
	"github.com/shandysiswandi/gostarter/pkg/goerror"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
)

var (
	errInvalidBody       = goerror.NewInvalidFormat("invalid request body")
	errFailedParseToUint = goerror.NewInvalidFormat("failed parse id to uint")
)

type httpEndpoint struct {
	telemetry *telemetry.Telemetry

	createRoleUC       domain.CreateRole
	findRoleUC         domain.FindRole
	createPermissionUC domain.CreatePermission
	findPermissionUC   domain.FindPermission
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

func (h *httpEndpoint) FindRole(c framework.Context) (any, error) {
	ctx, span := h.telemetry.Tracer().Start(c.Context(), "rbac.inbound.httpEndpoint.FindRole")
	defer span.End()

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return nil, errFailedParseToUint
	}

	resp, err := h.findRoleUC.Call(ctx, domain.FindRoleInput{ID: id})
	if err != nil {
		return nil, err
	}

	return Role{
		ID:          id,
		Name:        resp.Name,
		Description: resp.Description,
	}, nil
}

func (h *httpEndpoint) FindPermission(c framework.Context) (any, error) {
	ctx, span := h.telemetry.Tracer().Start(c.Context(), "rbac.inbound.httpEndpoint.FindPermission")
	defer span.End()

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return nil, errFailedParseToUint
	}

	resp, err := h.findPermissionUC.Call(ctx, domain.FindPermissionInput{ID: id})
	if err != nil {
		return nil, err
	}

	return Permission{
		ID:          id,
		Name:        resp.Name,
		Description: resp.Description,
	}, nil
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
