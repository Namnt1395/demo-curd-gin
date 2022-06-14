package router

import (
	"crypto/rsa"
	"demo-curd/config"
	"demo-curd/i18n"
	"demo-curd/util"
	"demo-curd/util/constant"
	"errors"
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/bmatcuk/doublestar/v3"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"reflect"
	"strings"
	"time"
)

const (
	reset            = "\033[0m"
	JWT_IDENTITY_KEY = "id"
	JWT_USER_ID      = "user_id"
	JWT_AUTHORITIES  = "authorities"
)

type Router struct {
	Engine                   *gin.Engine
	AuthMiddleware           *jwt.GinJWTMiddleware
	I18n                     *i18n.I18n
	PrivateKey               *rsa.PrivateKey
	LongRefreshExpTime       time.Duration
	CustomAuthorizedHandlers []CustomAuthorizedHandler
}

type CustomAuthorizedHandler interface {
	Authorize(c *gin.Context, authenticationData interface{}, authorities []interface{}) bool
}

func NewRouterWithoutAuthMw(c config.Config, i18n *i18n.I18n) (*Router, error) {
	return NewRouter(c, i18n, jwt.GinJWTMiddleware{})
}

func NewRouter(c config.Config, i18n *i18n.I18n, jwtMdw jwt.GinJWTMiddleware) (*Router, error) {
	e := gin.New()

	e.RedirectTrailingSlash = true
	e.RedirectFixedPath = true

	//e.Use(gin.Logger())
	e.Use(logger.SetLogger())

	// CORS
	corsMiddleware, err := initCorsMiddleware(c)
	if err != nil {
		return nil, err
	}
	e.Use(corsMiddleware)

	// the jwt middleware
	customAuthorizedHandlers := make([]CustomAuthorizedHandler, 0)
	authMiddleware, err := initJwtMiddleware(c, jwtMdw, customAuthorizedHandlers)
	if err != nil {
		log.Fatal().Err(err).Msg("JWT Error:" + err.Error())
		return nil, err
	}

	// When you use jwt.New(), the function is already automatically called for checking,
	// which means you don't need to call it again.
	errInit := authMiddleware.MiddlewareInit()

	if errInit != nil {
		log.Fatal().Err(errInit).Msg("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
		return nil, errInit
	}

	refreshExpTime, err := time.ParseDuration(c.Jwt.RefreshExpTime)
	if err != nil {
		return nil, err
	}
	return &Router{
		Engine:                   e,
		AuthMiddleware:           authMiddleware,
		I18n:                     i18n,
		LongRefreshExpTime:       refreshExpTime,
		PrivateKey:               privateKey(authMiddleware),
		CustomAuthorizedHandlers: make([]CustomAuthorizedHandler, 0),
	}, nil
}

func initJwtMiddleware(c config.Config, jwtMdw jwt.GinJWTMiddleware, handlers []CustomAuthorizedHandler) (*jwt.GinJWTMiddleware, error) {
	expiredTime, err := time.ParseDuration(c.Jwt.ExpiredTime)
	if err != nil {
		return nil, err
	}
	refreshExpTime, err := time.ParseDuration(c.Jwt.RefreshExpTime)
	if err != nil {
		return nil, err
	}
	defAuthorizedMw := DefAuthorizedMw(c, handlers)

	authenticator := jwtMdw.Authenticator
	payloadFunc := jwtMdw.PayloadFunc
	identityHandler := jwtMdw.IdentityHandler
	if identityHandler == nil {
		identityHandler = defAuthorizedMw.IdentityHandler
	}
	authorizator := jwtMdw.Authorizator
	if authorizator == nil {
		authorizator = defAuthorizedMw.Authorizator
	}
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:            c.Jwt.Realm,
		SigningAlgorithm: c.Jwt.SigningAlg,
		Key:              []byte(c.Jwt.Secret),
		Timeout:          expiredTime,
		MaxRefresh:       refreshExpTime,
		IdentityKey:      JWT_IDENTITY_KEY,
		Authenticator:    authenticator,
		PayloadFunc:      payloadFunc,
		IdentityHandler:  identityHandler,
		Authorizator:     authorizator,
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
	return authMiddleware, err
}

func initCorsMiddleware(c config.Config) (gin.HandlerFunc, error) {
	maxAge, err := time.ParseDuration(c.CORS.MaxAge)
	if err != nil {
		return nil, nil
	}

	return cors.New(cors.Config{
		AllowOrigins:     c.CORS.AllowOrigins,
		AllowMethods:     c.CORS.AllowMethods,
		AllowHeaders:     c.CORS.AllowHeaders,
		ExposeHeaders:    c.CORS.ExposeHeaders,
		AllowCredentials: c.CORS.AllowCredentials,
		MaxAge:           maxAge,
	}), nil
}

func privateKey(privKeyFile *jwt.GinJWTMiddleware) *rsa.PrivateKey {
	value := util.GetUnexportedField(reflect.ValueOf(privKeyFile).Elem().FieldByName("privKey"))
	if value == nil {
		return nil
	}
	return value.(*rsa.PrivateKey)
}

func DefAuthorizedMw(cfg config.Config, handlers []CustomAuthorizedHandler) jwt.GinJWTMiddleware {
	return jwt.GinJWTMiddleware{
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			log.Debug().Msgf("IdentityHandler, userId: %v", claims[JWT_USER_ID])
			return claims[JWT_USER_ID]
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			claims := jwt.ExtractClaims(c)
			authorities := claims[JWT_AUTHORITIES].([]interface{})
			log.Debug().Msgf("Authorizator, identity data: %v", data)
			log.Debug().Msgf("authorities: %v", authorities)
			return HandleAuthorizationWithAuthorities(data, c, cfg, authorities, handlers)
		},
	}
}

func HandleAuthorizationWithAuthorities(data interface{}, c *gin.Context, cfg config.Config, authorities []interface{}, handlers []CustomAuthorizedHandler) bool {
	if len(cfg.Security.AuthorizedRequests) > 0 {
		for _, req := range cfg.Security.AuthorizedRequests {
			if len(req.Urls) > 0 {
				for _, url := range req.Urls {
					auth, matchUrlOrMet := authorizePerUrl(data, c, url, req, authorities, handlers)
					if !matchUrlOrMet {
						continue
					}
					return auth
				}
			}
		}
	}
	return false
}

func authorizePerUrl(data interface{}, c *gin.Context, url string, req config.ConfigAuthorizedRequests, authorities []interface{}, handlers []CustomAuthorizedHandler) (bool, bool) {
	arrUrl := strings.Split(url, ":")
	pathMatched, err := doublestar.Match(arrUrl[0], c.FullPath())
	util.Must(err)
	if pathMatched {
		if len(arrUrl) <= 1 {
			return false, false
		}

		methodPatched, err := doublestar.Match(arrUrl[1], c.Request.Method)
		util.Must(err)
		if !methodPatched {
			return false, false
		}
		if req.Access == constant.AccessHasPermission {
			if auth := authorizeHasPermission(req, authorities); auth {
				return true, true
			}
		} else if req.Access == constant.AccessHasRole {
			if auth := authorizeHasRole(req, authorities); auth {
				return true, true
			}
		} else if req.Access == constant.AccessPermitAll {
			return true, true
		} else if req.Access == constant.AccessDenyAll {
			return false, true
		} else if req.Access == constant.AccessCustom {
			if handlers != nil && len(handlers) > 0 {
				for _, h := range handlers {
					if auth := h.Authorize(c, data, authorities); auth {
						return true, true
					}
				}
			}
		} else {
			panic(errors.New("Invalid access type, must be has permission, has role, permit all or deny all"))
		}
	}
	return false, false
}

func authorizeHasPermission(req config.ConfigAuthorizedRequests, authorities []interface{}) bool {
	for _, p := range req.Permissions {
		_, find := util.FindStringInGeneric(authorities, p)
		if find {
			return true
		}
	}
	return false
}

func authorizeHasRole(v config.ConfigAuthorizedRequests, authorities []interface{}) bool {
	for _, p := range v.Roles {
		_, find := util.FindStringInGeneric(authorities, p)
		if find {
			return true
		}
	}
	return false
}

func (r *Router) RegisterCustomAuthorizedHandler(cah CustomAuthorizedHandler) {
	r.CustomAuthorizedHandlers = append(r.CustomAuthorizedHandlers, cah)
}

func (r *Router) InitSwagger(c config.Config) {
	url := ginSwagger.URL(fmt.Sprintf("%v/swagger/doc.json", c.Swagger.Url))
	r.Engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}
