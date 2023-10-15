package providers

import "github.com/gin-gonic/gin"

type GinWrapper interface {
	ShouldBindJSON(any) error
	JSON(int, any)
	Param(string) string
	Inner() *gin.Context
}

type GinWrapperProvider interface {
	Wrap(ctx *gin.Context) GinWrapper
}

type ginWrapperProviderImpl struct{}

func NewGinWrapperProvider() GinWrapperProvider {
	return &ginWrapperProviderImpl{}
}

func (g *ginWrapperProviderImpl) Wrap(ctx *gin.Context) GinWrapper {
	return &ginWrapperImpl{ctx}
}

type ginWrapperImpl struct {
	ctx *gin.Context
}

func (g *ginWrapperImpl) ShouldBindJSON(obj any) error {
	return g.ctx.ShouldBindJSON(obj)
}

func (g *ginWrapperImpl) JSON(code int, obj any) {
	g.ctx.JSON(code, obj)
}

func (g *ginWrapperImpl) Param(key string) string {
	return g.ctx.Param(key)
}

func (g *ginWrapperImpl) Inner() *gin.Context {
	return g.ctx
}
