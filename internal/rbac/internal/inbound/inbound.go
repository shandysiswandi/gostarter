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
	CreatePermission domain.CreatePermission
}

func (in Inbound) RegisterRBACServiceServer() {
	he := &httpEndpoint{
		telemetry: in.Telemetry,
		//
		createRoleUC:       in.CreateRole,
		createPermissionUC: in.CreatePermission,
	}

	in.Router.Endpoint(http.MethodPost, "/rbac/roles", he.CreateRole)
	in.Router.Endpoint(http.MethodPost, "/rbac/permissions", he.CreatePermission)
}
