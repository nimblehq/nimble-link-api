package ginutils

import "github.com/gin-gonic/gin"

type ApplicationRouter struct {
	Router gin.IRouter
}

func (r *ApplicationRouter) Middlewares(middlewares ...gin.HandlerFunc) *ApplicationRouter {
	return &ApplicationRouter{
		Router: r.Router.Group("/", middlewares...),
	}
}

func (r *ApplicationRouter) Register(method string, path string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return r.Router.Handle(method, path, handlers...)
}
