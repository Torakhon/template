package v1

import (
	"api-gateway/api/handlers/models"
	token2 "api-gateway/api/tokens"
	cfg "api-gateway/config"
	pb "api-gateway/genproto/users"
	"api-gateway/pkg/etc"
	l "api-gateway/pkg/logger"
	"api-gateway/pkg/utils"
	"context"
	_ "fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"google.golang.org/protobuf/encoding/protojson"
	"net/http"
	"strings"
	"time"
)

// GetUser ...
// @Security ApiKeyAuth
// @Summary GetUser
// @Description Viewing a single User by id
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} models.User
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/user/info [get]
func (h *HandlerV1) GetUser(c *gin.Context) {

	token := c.GetHeader("Authorization")

	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid token"})
		return
	}

	claims, err := token2.ExtractClaim(token, []byte(cfg.Load().SigningKey))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "invalid token",
		})
	}
	id := claims["id"]
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()
	response, err := h.serviceManager.UserService().GetUser(
		ctx, &pb.GetUserReq{
			Id: cast.ToString(id),
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get user", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, models.User{
		Id:          response.Id,
		Name:        response.Name,
		LastName:    response.LastName,
		Role:        response.Role,
		Email:       response.Email,
		PhoneNumber: response.PhoneNumber,
	})
}

// UpUser ...
// @Security ApiKeyAuth
// @Summary GetUser
// @Description Viewing a single User by id
// @Tags user
// @Accept json
// @Produce json
// @Param UpUser body models.UpdateUserReq true "Up User"
// @Success 200 {object} models.UpdateUserRes
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/user/up_user [comment]
func (h *HandlerV1) UpUser(c *gin.Context) {
	var (
		body        models.UpdateUserReq
		jsonMarshal protojson.MarshalOptions
	)
	jsonMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}

	token := c.GetHeader("Authorization")
	if token == "" {
		token = c.GetHeader("Login")
	}
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid token"})
		return
	} else if strings.Contains(token, "Bearer") {
		token = strings.TrimPrefix(token, "Bearer ")
	}

	claims, err := token2.ExtractClaim(token, []byte(cfg.Load().SigningKey))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "invalid token",
		})
		return
	}
	if body.Password != "" || body.Password != "string" {
		passwordHSH, err := etc.HashPassword(body.Password)
		if err != nil {
			h.log.Error("error hsh password")
		}
		body.Password = passwordHSH
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()
	response, err := h.serviceManager.UserService().UpdateUser(ctx, &pb.UpdateUserReq{
		Name:        body.Name,
		LastName:    body.LastName,
		Password:    body.Password,
		PhoneNumber: body.PhoneNumber,
		Id:          cast.ToString(claims["id"]),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get user", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, models.UpdateUserRes{
		Name:        response.Name,
		LastName:    response.LastName,
		Role:        response.Role,
		PhoneNumber: response.PhoneNumber,
	})
}

// GetAllUsers ...
// @Security ApiKeyAuth
// @Summary GetAllUsers
// @Description Viewing a single User by id
// @Tags user
// @Accept json
// @Produce json
// @Param GetAllUsers body models.UsersReq true "Get All Users"
// @Success 200 {object} models.UsersRes
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/user/{page}/{limit} [comment]
func (h *HandlerV1) GetAllUsers(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	params, errStr := utils.ParseQueryParams(queryParams)
	if errStr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errStr[0],
		})
		h.log.Error("failed to parse query params json" + errStr[0])
		return
	}

	var jsonMarshal protojson.MarshalOptions
	jsonMarshal.UseProtoNames = true

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.UserService().GetAllUsers(ctx, &pb.GetAllUsersReq{
		Page:  params.Page,
		Limit: params.Limit,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to list users", l.Error(err))
		return
	}
	var res []models.AllUsers
	for _, i := range response.Users {
		r := models.AllUsers{}
		r.Name = i.Name
		r.LastName = i.LastName
		r.Email = i.Email
		r.PhoneNumber = i.PhoneNumber
		res = append(res, r)
	}
	c.JSON(http.StatusOK, &models.UsersRes{Users: res})
}
