package v1

import (
	"api-gateway/api/handlers/models"
	token2 "api-gateway/api/tokens"
	cfg "api-gateway/config"
	pb "api-gateway/genproto/users"
	l "api-gateway/pkg/logger"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"google.golang.org/protobuf/encoding/protojson"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// UpdateRole ...
// @Security ApiKeyAuth
// @Summary Update Role user
// @Description Viewing a single User by id
// @Tags suAdmin
// @Accept json
// @Produce json
// @Param UpdateRole body models.UpdateRolReq true "Update role"
// @Success 200 {object} models.UpdateRolRes
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/suAdmin/up_role [put]
func (h *HandlerV1) UpdateRole(c *gin.Context) {
	var (
		body        models.UpdateRolReq
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

	res, err := h.serviceManager.UserService().UpdateRole(ctx, &pb.UpdateRoleReq{
		Id:      cast.ToString(id),
		NewRole: body.Role,
	})
	if err != nil {
		h.log.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, &models.UpdateRolRes{Status: res.Stats})
}

// DeleteUser ...
// @Security ApiKeyAuth
// @Summary DeleteUser
// @Tags suAdmin
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/suAdmin/delete_user [DELETE]
func (h HandlerV1) DeleteUser(c *gin.Context) {
	id := c.Param("id")

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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()
	_, err := h.serviceManager.UserService().DeleteUser(ctx, &pb.DeleteUserReq{
		Field: "id",
		Value: id,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get user", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true})
}

// AdGetAllUsers ...
// @Security ApiKeyAuth
// @Summary AdGetAllUsers
// @Description Viewing a single User by id
// @Tags suAdmin
// @Accept json
// @Produce json
// @Param page query int true "Page number"
// @Param limit query int true "Items per page"
// @Success 200 {object} models.AdminUsersRes
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/suAdmin/get_all_users [get]
func (h *HandlerV1) AdGetAllUsers(c *gin.Context) {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid page number",
		})
		h.log.Error("invalid page number", l.Error(err))
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid limit",
		})
		h.log.Error("invalid limit", l.Error(err))
		return
	}

	var jsonMarshal protojson.MarshalOptions
	jsonMarshal.UseProtoNames = true

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()
	response, err := h.serviceManager.UserService().GetAllUsers(ctx, &pb.GetAllUsersReq{
		Page:  int64(page),
		Limit: int64(limit),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to list users", l.Error(err))
		return
	}
	var res models.AdminUsersRes
	for _, user := range response.Users {
		r := models.AdminUser{
			Id:        user.Id,
			UserName:  user.UserName,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Role:      user.Role,
			Bio:       user.Bio,
			Website:   user.WebSite,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}

		res.Users = append(res.Users, r)
	}
	c.JSON(http.StatusOK, res)
}
