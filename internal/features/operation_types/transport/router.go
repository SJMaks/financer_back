package transport

import "github.com/gin-gonic/gin"

func RegisterPrivateRouter(r *gin.RouterGroup, h *Handler, mid gin.HandlerFunc) {
	g := r.Group("/operation_types", mid)
	{
		g.GET("", h.GetAll)
	}
}
