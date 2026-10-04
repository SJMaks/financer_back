package transport

import "github.com/gin-gonic/gin"

func RegisterPrivateRouter(r *gin.RouterGroup, h *Handler, mid gin.HandlerFunc) {
	g := r.Group("/accounts/:accountID/operations", mid)
	{
		g.POST("", h.Create)
		g.GET("", h.GetAllOfAccount)
		g.GET("/:operationID", h.GetByID)
		g.PATCH("/:operationID", h.UpdateByID)
		g.DELETE("/:operationID", h.DeleteByID)
	}
}
