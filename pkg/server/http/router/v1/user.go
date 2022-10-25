package v1

import (
	engine "fin-asg/config/gin"
	"fin-asg/pkg/domain/user"
	"fin-asg/pkg/server/http/middleware"
	"fin-asg/pkg/server/http/router"
	"github.com/gin-gonic/gin"
)

type UserRouterImpl struct {
	ginEngine   engine.HttpServer
	routerGroup *gin.RouterGroup
	userHandler user.UserHandler
	authMiddleware middleware.AuthMiddleware
}

func NewUserRouter(ginEngine engine.HttpServer, userHandler user.UserHandler) router.Router {

	// setiap yang /v1/user
	// harus melakukan pengecheckan auth
	// sehingga kita bisa meletakkan middleware di dalam group kita

	routerGroup := ginEngine.GetGin().Group("mygram/v1/user")
	return &UserRouterImpl{ginEngine: ginEngine, routerGroup: routerGroup, userHandler: userHandler}
}

func (u *UserRouterImpl) get() {
	// all path for get method are here
	u.routerGroup.GET(":id", u.userHandler.GetUserByIdHdl)
}

func (u *UserRouterImpl) post() {
	// all path for post method are here
	u.routerGroup.POST("/create", u.userHandler.InsertUserHdl)
}

func (u *UserRouterImpl) put() {
	u.routerGroup.PUT("", u.authMiddleware.CheckJWTAuth, u.userHandler.UpdateUserHdl)
}

func(u *UserRouterImpl) delete() {
	u.routerGroup.PUT("", u.authMiddleware.CheckJWTAuth, u.userHandler.UpdateUserHdl)
}

func (u *UserRouterImpl) Routers() {
	u.get()
	u.post()
	u.put()
	u.delete()
}
