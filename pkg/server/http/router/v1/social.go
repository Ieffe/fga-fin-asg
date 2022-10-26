package v1

import (
	engine "fin-asg/config/gin"
	"fin-asg/pkg/domain/social"
	"fin-asg/pkg/server/http/middleware"
	"fin-asg/pkg/server/http/router"

	"github.com/gin-gonic/gin"
)

type SocialRouterImpl struct {
	ginEngine          engine.HttpServer
	routerGroup        *gin.RouterGroup
	socialHandler social.SocialHandler
	authMiddleware     middleware.AuthMiddleware
}

func (p *SocialRouterImpl) get() {
	p.routerGroup.GET("", p.authMiddleware.CheckJWTAuth, p.socialHandler.GetSocialsHdl)
}

func (p *SocialRouterImpl) post() {
	p.routerGroup.POST("", p.authMiddleware.CheckJWTAuth, p.socialHandler.InsertSocialHdl)
}

func (p *SocialRouterImpl) put() {
	p.routerGroup.PUT("/:socialId", p.authMiddleware.CheckJWTAuth, p.socialHandler.UpdateSocialHdl)
}

func (p *SocialRouterImpl) delete() {
	p.routerGroup.DELETE("/:socialId", p.authMiddleware.CheckJWTAuth, p.socialHandler.DeleteSocialHdl)
}

func (p *SocialRouterImpl) Routers() {
	p.get()
	p.post()
	p.put()
	p.delete()
}

func NewSocialRouter(ginEngine engine.HttpServer, socialHandler social.SocialHandler, authMiddleware middleware.AuthMiddleware) router.Router {
	routerGroup := ginEngine.GetGin().Group("/api/mygram/v1/socials")
	return &SocialRouterImpl{ginEngine: ginEngine, routerGroup: routerGroup, socialHandler: socialHandler, authMiddleware: authMiddleware}
}
