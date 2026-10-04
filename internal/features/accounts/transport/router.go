package transport

import "github.com/gin-gonic/gin"

func RegisterPrivateRouter(r *gin.RouterGroup, h *Handler, mid gin.HandlerFunc) {
	g := r.Group("/accounts", mid)
	{
		g.POST("", h.Create)
		g.GET("", h.GetAllOfUser)
		g.GET("/:id", h.GetByID)
		g.PATCH("/:id", h.UpdateByID)
		g.DELETE("/:id", h.DeleteByID)
	}
}
