package api

import (
	_ "api-gateway/api/docs"
	casbinC "api-gateway/api/handlers/middleware/casbin"
	"api-gateway/api/handlers/models"
	v1 "api-gateway/api/handlers/v1"
	"api-gateway/config"
	"api-gateway/pkg/logger"
	"api-gateway/queue/kafka/producer"
	"api-gateway/services"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	_ "github.com/swaggo/files"
	swaggerFiles "github.com/swaggo/files"
	_ "github.com/swaggo/gin-swagger"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"net/http"
)

// Option ...
type Option struct {
	Conf           config.Config
	Logger         logger.Logger
	Enforcer       *casbin.Enforcer
	ServiceManager services.IServiceManager
	Writer         producer.KafkaProducer
}

// New ...
// @title welcome to
// @version 1.0
// @host localhost:9091
func New(option Option) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	//js, _ := json.Marshal("###########################################################################################")
	//
	//err := option.Writer.ProduceMessages("test-topic", js)

	handlerV1 := v1.New(&v1.HandlerV1Config{
		Logger:         option.Logger,
		ServiceManager: option.ServiceManager,
		Enforcer:       option.Enforcer,
		Cfg:            option.Conf,
		Writer:         option.Writer,
	})

	policies := [][]string{}
	for _, policy := range policies {
		_, err := option.Enforcer.AddPolicy(policy)
		if err != nil {
			option.Logger.Error("error during investor enforcer add policies", zap.Error(err))
		}
	}
	_, err := option.Enforcer.AddGroupingPolicy("suAdmin", "user")

	err = option.Enforcer.SavePolicy()
	if err != nil {
		logger.Error(err)
	}

	// @securityDefinitions.apikey ApiKeyAuth
	// @in header
	// @name Authorization

	api := router.Group("/v1")

	//Authorization
	api.Use(casbinC.CheckCasbinPermission(option.Enforcer, option.Conf))
	auth := api.Group("/auth")
	auth.POST("/register", handlerV1.Register)
	auth.POST("/authorization", handlerV1.Authorization)
	auth.POST("/Login", handlerV1.Login)

	//suAdmin
	suAdmin := api.Group("suAdmin")
	suAdmin.GET("/roles", option.ListRoles())
	suAdmin.DELETE("/:role", option.DeleteRole())
	suAdmin.POST("/add-user-role", option.AddPolicy())
	suAdmin.POST("/up_role", handlerV1.UpdateRole)

	//user
	user := api.Group("/user")
	user.GET("/info", handlerV1.GetUser)
	user.POST("up_user", handlerV1.UpUser)
	user.POST("/{page}/{limit}", handlerV1.GetAllUsers)
	user.GET("/get_product", handlerV1.GetProduct)
	user.POST("/create_product", handlerV1.CreateProduct)

	url := ginSwagger.URL("swagger/doc.json")
	api.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return router
}

// ListRoles ...
// @Security ApiKeyAuth
// @Summary Get list of roles
// @Description Get list of roles
// @Tags Super Admin
// @Accept json
// @Produce json
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/suAdmin/roles [GET]
func (h *Option) ListRoles() gin.HandlerFunc {
	return func(c *gin.Context) {
		obj := h.Enforcer.GetAllObjects() // Potential source of nil pointer dereference
		role := h.Enforcer.GetAllSubjects()
		act := h.Enforcer.GetAllActions()
		respModel := &models.ListRolesResponse{
			Roles: role,
			Obj:   obj,
			Act:   act,
		}
		c.JSON(http.StatusOK, respModel)
	}
}

// DeleteRole ...
// @Summary Delete user-role by id
// @Security ApiKeyAuth
// @Description Delete user-role by id
// @Tags Super Admin
// @Accept json
// @Produce json
// @Param role path string true "role"
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/suAdmin/{role} [DELETE]
func (h *Option) DeleteRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.Param("role")

		var rolePolicies [][]string

		for _, p := range h.Enforcer.GetFilteredPolicy(0, role) {

			rolePolicies = append(rolePolicies, p)
		}

		resp, err := h.Enforcer.DeleteRole(role)
		if err != nil {
			h.Logger.Error("rbacHandler/DeleteRole", zap.Error(err))
			c.JSON(http.StatusInternalServerError, models.Error{Message: err.Error()})
			return
		}
		err = h.Enforcer.SavePolicy()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.Error{Err: err})
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}

// AddPolicy ...
// @Summary Create new user-role
// @Description Create new user-role
// @Tags Super Admin
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param create body models.AddRole true "create"
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /v1/suAdmin/add-user-role [POST]
func (h *Option) AddPolicy() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody models.AddRole
		if err := c.ShouldBindJSON(&reqBody); err != nil {
			h.Logger.Error("rbacHandler/CreateUserRole", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
		fmt.Println(reqBody)

		if _, err := h.Enforcer.AddPolicy(reqBody.Role, reqBody.Url, reqBody.Method); err != nil {
			h.Logger.Error("error on grantAccess", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err := h.Enforcer.SavePolicy()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.Error{Err: err})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"role":   reqBody.Role,
			"url":    reqBody.Url,
			"method": reqBody.Method,
		})
	}
}
