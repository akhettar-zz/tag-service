# jwt-go-plugin
This is plugin is used to parse JWT tokens and validate the integrity of the token and extract the accountId if token is of a User type. This plugin is designed specifically to work with WS02 JWT tokens. It does the following:
* Validate the JWT token including signature validation
* Checks the claims (AccountId, Issuer)
* Extract the `AccountId` and set in Request header

[![CircleCI](https://ci.shared.astoapp.co.uk/gh/BetaProjectWave/jwt-go-plugin.svg?style=svg&circle-token=b8a236f1d65982e01ff2048b5c85f2a5e319eebc)](https://ci.shared.astoapp.co.uk/gh/BetaProjectWave/jwt-go-plugin)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/251df2198c6347b8ae0a25bf6bc134fd)](https://www.codacy.com?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=BetaProjectWave/jwt-go-plugin&amp;utm_campaign=Badge_Grade)
[![codecov](https://codecov.io/gh/BetaProjectWave/jwt-go-plugin/branch/master/graph/badge.svg?token=NaLhHNGQ59)](https://codecov.io/gh/BetaProjectWave/jwt-go-plugin)

Here is an example of how can this be invoked. It should be attached to all endpoint you would want to authenticate against the user.

```go
// Registers all the routes
func (handler *TagHandler) CreateRouter() *gin.Engine {

	// Create router
	router := gin.New()
	router.Use(logger.Logger())
	router.Use(gin.Recovery())

	router.GET("/tags", jwt.GinJWTMiddleware().MiddlewareFunc(), handler.GetAllTags)
	router.GET("/tags/:id", jwt.GinJWTMiddleware().MiddlewareFunc(), handler.GetTag)
	router.DELETE("/tags/:id", jwt.GinJWTMiddleware().MiddlewareFunc(), handler.DeleteTag)
	router.GET("/health", handler.Health)
	router.POST("/tags", jwt.GinJWTMiddleware().MiddlewareFunc(), handler.CreateTag)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}

```
