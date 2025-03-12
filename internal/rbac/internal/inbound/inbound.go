package inbound

import (
	"net/http"

	"github.com/shandysiswandi/goreng/telemetry"
	"github.com/shandysiswandi/gostarter/internal/rbac/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/framework"
)

type Inbound struct {
	Router    *framework.Router
	Telemetry *telemetry.Telemetry
	//
	CreateRole domain.CreateRole
	FindRole   domain.FindRole
	FetchRole  domain.FetchRole
	UpdateRole domain.UpdateRole
	//
	CreatePermission domain.CreatePermission
	FindPermission   domain.FindPermission
	FetchPermission  domain.FetchPermission
	UpdatePermission domain.UpdatePermission
}

func (in Inbound) RegisterRBACServiceServer() {
	he := &httpEndpoint{
		telemetry: in.Telemetry,
		//
		createRoleUC: in.CreateRole,
		findRoleUC:   in.FindRole,
		fetchRoleUC:  in.FetchRole,
		updateRoleUC: in.UpdateRole,
		//
		createPermissionUC: in.CreatePermission,
		findPermissionUC:   in.FindPermission,
		fetchPermissionUC:  in.FetchPermission,
		updatePermissionUC: in.UpdatePermission,
	}

	in.Router.Endpoint(http.MethodPost, "/rbac/roles", he.CreateRole)
	in.Router.Endpoint(http.MethodGet, "/rbac/roles", he.FetchRole)
	in.Router.Endpoint(http.MethodGet, "/rbac/roles/:id", he.FindRole)
	in.Router.Endpoint(http.MethodPut, "/rbac/roles/:id", he.UpdateRole)
	//
	in.Router.Endpoint(http.MethodPost, "/rbac/permissions", he.CreatePermission)
	in.Router.Endpoint(http.MethodGet, "/rbac/permissions", he.FetchPermission)
	in.Router.Endpoint(http.MethodGet, "/rbac/permissions/:id", he.FindPermission)
	in.Router.Endpoint(http.MethodPut, "/rbac/permissions/:id", he.UpdatePermission)
}
