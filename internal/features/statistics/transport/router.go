package transport

import "github.com/gin-gonic/gin"

func RegisterPrivateRouter(r *gin.RouterGroup, h *Handler, mid gin.HandlerFunc) {
	g := r.Group("/accounts", mid)
	{
		g.GET("/:accountID/statistics", h.GetByAccount)
	}
	g = r.Group("/statistics", mid)
	{
		g.GET("", h.GetByUser)
	}
}
