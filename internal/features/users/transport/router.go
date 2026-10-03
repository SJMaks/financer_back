package transport

import "github.com/gin-gonic/gin"

func RegisterRouter(r *gin.RouterGroup, h *Handler) {
	g := r.Group("/auth")
	{
		g.POST("/register", h.Register)
		g.POST("/login", h.Login)
	}
}
