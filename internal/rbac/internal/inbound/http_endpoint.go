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

	createRoleUC domain.CreateRole
	findRoleUC   domain.FindRole
	fetchRoleUC  domain.FetchRole
	updateRoleUC domain.UpdateRole
	//
	createPermissionUC domain.CreatePermission
	findPermissionUC   domain.FindPermission
	fetchPermissionUC  domain.FetchPermission
	updatePermissionUC domain.UpdatePermission
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
		ID:          resp.ID,
		Name:        resp.Name,
		Description: resp.Description,
	}, nil
}

func (h *httpEndpoint) FetchRole(c framework.Context) (any, error) {
	ctx, span := h.telemetry.Tracer().Start(c.Context(), "rbac.inbound.httpEndpoint.FetchRole")
	defer span.End()

	resp, err := h.fetchRoleUC.Call(ctx, domain.FetchRoleInput{
		Cursor: c.Query("cursor"),
		Limit:  c.Query("limit"),
		Name:   c.Query("name"),
	})
	if err != nil {
		return nil, err
	}

	roles := make([]Role, 0)
	for _, role := range resp.Roles {
		roles = append(roles, Role{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
		})
	}

	return FetchRoleResponse{
		Roles: roles,
		Pagination: Pagination{
			NextCursor: resp.NextCursor,
			HasMore:    resp.HasMore,
		},
	}, nil
}

func (h *httpEndpoint) UpdateRole(c framework.Context) (any, error) {
	ctx, span := h.telemetry.Tracer().Start(c.Context(), "rbac.inbound.httpEndpoint.UpdateRole")
	defer span.End()

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return nil, errFailedParseToUint
	}

	var req UpdateRoleRequest
	if err := json.NewDecoder(c.Body()).Decode(&req); err != nil {
		return nil, errInvalidBody
	}

	resp, err := h.updateRoleUC.Call(ctx, domain.UpdateRoleInput{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}

	return Role{
		ID:          resp.ID,
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
		ID:          resp.ID,
		Name:        resp.Name,
		Description: resp.Description,
	}, nil
}

func (h *httpEndpoint) FetchPermission(c framework.Context) (any, error) {
	ctx, span := h.telemetry.Tracer().Start(c.Context(), "rbac.inbound.httpEndpoint.FetchPermission")
	defer span.End()

	resp, err := h.fetchPermissionUC.Call(ctx, domain.FetchPermissionInput{
		Cursor: c.Query("cursor"),
		Limit:  c.Query("limit"),
		Name:   c.Query("name"),
	})
	if err != nil {
		return nil, err
	}

	perms := make([]Permission, 0)
	for _, perm := range resp.Permissions {
		perms = append(perms, Permission{
			ID:          perm.ID,
			Name:        perm.Name,
			Description: perm.Description,
		})
	}

	return FetchPermissionResponse{
		Permissions: perms,
		Pagination: Pagination{
			NextCursor: resp.NextCursor,
			HasMore:    resp.HasMore,
		},
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

func (h *httpEndpoint) UpdatePermission(c framework.Context) (any, error) {
	ctx, span := h.telemetry.Tracer().Start(c.Context(), "rbac.inbound.httpEndpoint.UpdatePermission")
	defer span.End()

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return nil, errFailedParseToUint
	}

	var req UpdateRoleRequest
	if err := json.NewDecoder(c.Body()).Decode(&req); err != nil {
		return nil, errInvalidBody
	}

	resp, err := h.updatePermissionUC.Call(ctx, domain.UpdatePermissionInput{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}

	return Permission{
		ID:          resp.ID,
		Name:        resp.Name,
		Description: resp.Description,
	}, nil
}
