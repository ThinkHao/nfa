package ipgroup

import (
	srvv1 "github.com/marmotedu/iam/internal/apiserver/service/v1"
	"github.com/marmotedu/iam/internal/apiserver/store"
)

// UserController create a user handler used to handle request for user resource.
type IPGroupController struct {
	srv srvv1.Service
}

// NewUserController creates a user handler.
func NewIPGroupController(store store.Factory) *IPGroupController {
	return &IPGroupController{
		srv: srvv1.NewService(store),
	}
}