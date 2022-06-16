package ipgroup

import (
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/component-base/pkg/core"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"

	"github.com/marmotedu/iam/pkg/log"
)

// Get get an user by the user identifier.
func (i *IPGroupController) Get(c *gin.Context) {
	log.L(c).Info("get ipgroup function called.")

	ipgroup, err := i.srv.IPGroups().Get(c, c.Param("ipgroup_name"), metav1.GetOptions{})
	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, user)
}
