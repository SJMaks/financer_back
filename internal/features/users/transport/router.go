package transport

import "github.com/gin-gonic/gin"

func RegisterRouter(r *gin.RouterGroup, h *Handler) {
	g := r.Group("/users")
	{
		g.GET("/ping", h.Ping)
	}
}
