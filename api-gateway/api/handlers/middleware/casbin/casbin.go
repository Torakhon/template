package casbin

import (
	tokens "api-gateway/api/tokens"
	"api-gateway/config"
	cfg "api-gateway/config"
	"errors"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
	"strings"
)

type Handler struct {
	config   config.Config
	enforcer *casbin.Enforcer
}

func CheckCasbinPermission(casbin *casbin.Enforcer, conf config.Config) gin.HandlerFunc {
	casbHandler := &Handler{
		config:   conf,
		enforcer: casbin,
	}
	return func(ctx *gin.Context) {
		allowed, err := casbHandler.CheckPermission(ctx.Request)
		if err != nil {
			ctx.AbortWithError(http.StatusUnauthorized, err)
		}

		if !allowed {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "permission denied",
			})
		}
	}
}

func (h *Handler) NewAuthorized() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, _ := h.GetRole(ctx.Request)

		sub := role

		obj := ctx.Request.URL.Path
		act := ctx.Request.Method

		e, err := casbin.NewEnforcer("config/auth.conf", "config/auth.csv")
		if err != nil {
			return
		}
		t, err := e.Enforce(sub, obj, act)

		if err != nil {
			return
		}
		if t {
			fmt.Println(t)
			ctx.Next()
			return
		}
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "denied permission",
		})
	}
}

func (h *Handler) GetRole(ctx *http.Request) (string, int) {
	var t string
	token := ctx.Header.Get("authorization")
	if token == "" {
		return "unauthorazed", http.StatusUnauthorized
	} else if strings.Contains(token, "Bearer") {
		t = strings.TrimPrefix(token, "Bearer ")
	} else {
		t = token
	}

	claims, err := tokens.ExtractClaim(t, []byte(cfg.Load().SigningKey))
	if err != nil {
		return "unauthorazed", http.StatusUnauthorized
	}

	return cast.ToString(claims["role"]), 0
}

func (h *Handler) CheckPermission(r *http.Request) (bool, error) {
	role, status := h.GetRole(r)

	if role == "unauthorazed" {
		return true, nil
	}

	if status != 0 {
		return false, errors.New(role)
	}
	method := r.Method
	action := r.URL.Path

	c, err := h.enforcer.Enforce(role, action, method)
	if err != nil {
		return false, err
	}

	return c, nil
}
