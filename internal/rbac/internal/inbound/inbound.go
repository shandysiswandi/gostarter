package inbound

import (
	"net/http"

	"github.com/shandysiswandi/gostarter/internal/rbac/internal/domain"
	"github.com/shandysiswandi/gostarter/pkg/framework"
	"github.com/shandysiswandi/gostarter/pkg/telemetry"
)

type Inbound struct {
	Router    *framework.Router
	Telemetry *telemetry.Telemetry
	//
	CreateRole       domain.CreateRole
	FindRole         domain.FindRole
	CreatePermission domain.CreatePermission
	FindPermission   domain.FindPermission
}

func (in Inbound) RegisterRBACServiceServer() {
	he := &httpEndpoint{
		telemetry: in.Telemetry,
		//
		createRoleUC:       in.CreateRole,
		findRoleUC:         in.FindRole,
		createPermissionUC: in.CreatePermission,
		findPermissionUC:   in.FindPermission,
	}

	in.Router.Endpoint(http.MethodPost, "/rbac/roles", he.CreateRole)
	in.Router.Endpoint(http.MethodGet, "/rbac/roles/:id", he.FindRole)
	//
	in.Router.Endpoint(http.MethodPost, "/rbac/permissions", he.CreatePermission)
	in.Router.Endpoint(http.MethodGet, "/rbac/permissions/:id", he.FindPermission)
}
