package transport

import "github.com/gin-gonic/gin"

func RegisterPublicRouter(r *gin.RouterGroup, h *Handler) {
	g := r.Group("/auth")
	{
		g.POST("/register", h.Register)
		g.POST("/login", h.Login)
	}
}

func RegisterPrivateRouter(r *gin.RouterGroup, h *Handler, mid gin.HandlerFunc) {
	g := r.Group("/users", mid)
	{
		g.GET("/me", h.GetMe)
		g.PATCH("/me", h.UpdateMe)
		g.DELETE("/me", h.DeleteMe)
	}
}
