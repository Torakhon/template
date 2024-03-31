package v1

import (
	"api-gateway/api/handlers/models"
	"api-gateway/api/handlers/models/postModel"
	token2 "api-gateway/api/tokens"
	cfg "api-gateway/config"
	postPb "api-gateway/genproto/post"
	pb "api-gateway/genproto/users"
	"api-gateway/pkg/etc"
	l "api-gateway/pkg/logger"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"google.golang.org/protobuf/encoding/protojson"
	"net/http"
	"strconv"
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
// @Router /v1/user/get_user [get]
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
	res, err := h.serviceManager.UserService().GetUser(
		ctx, &pb.GetUserReq{
			Field: "id",
			Value: cast.ToString(id),
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get user", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, models.User{
		UserName:  res.UserName,
		FirstName: res.FirstName,
		LastName:  res.LastName,
		Email:     res.Email,
		Role:      res.Role,
		Bio:       res.Bio,
		Website:   res.WebSite,
	})
}

// GetUserWithPosts ...
// @Security ApiKeyAuth
// @Summary GetUserWithPosts
// @Description Viewing a single User by id
// @Tags user
// @Accept json
// @Produce json
// @Param page query int true "Page number"
// @Param limit query int true "Items per page"
// @Success 200 {object} models.User
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/user/get_user_posts [get]
func (h *HandlerV1) GetUserWithPosts(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

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
		return
	}
	id := claims["id"]

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	res, err := h.serviceManager.UserService().GetUser(
		ctx, &pb.GetUserReq{
			Field: "id",
			Value: cast.ToString(id),
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get user", l.Error(err))
		return
	}
	posts, err := h.serviceManager.PostService().SearchPost(ctx, &postPb.SearchReq{
		Field: "user_id",
		Value: res.Id,
		Page:  int32(page),
		Limit: int32(limit),
	})
	fmt.Println(err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to search posts", l.Error(err))
		return
	}

	var userPosts []postModel.Post

	for _, post := range posts.Posts {
		ps := postModel.Post{
			Id:        post.Id,
			Title:     post.Title,
			Content:   post.Content,
			UserId:    post.UserId,
			Category:  post.Category,
			Likes:     post.Likes,
			Dislikes:  post.Dislikes,
			Views:     post.Views,
			CreatedAt: post.CreatedAt,
		}
		userPosts = append(userPosts, ps)
	}

	c.JSON(http.StatusOK, models.User{
		UserName:  res.UserName,
		FirstName: res.FirstName,
		LastName:  res.LastName,
		Email:     res.Email,
		Role:      res.Role,
		Bio:       res.Bio,
		Website:   res.WebSite,
		Posts:     userPosts,
	})
}

// UpdateUser ...
// @Security ApiKeyAuth
// @Summary UpdateUser
// @Description Viewing a single User by id
// @Tags user
// @Accept json
// @Produce json
// @Param UpUser body models.UpdateUser true "Up User"
// @Success 200 {object} models.UpdateUserRes
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/user/up_user [put]
func (h *HandlerV1) UpdateUser(c *gin.Context) {
	var (
		body        models.UpdateUser
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
	res, err := h.serviceManager.UserService().UpdateUser(ctx, &pb.UpdateUserReq{
		Id:        cast.ToString(claims["id"]),
		UserName:  body.UserName,
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Password:  body.Password,
		Bio:       body.Bio,
		WebSite:   body.Website,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get user", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, models.UpdateUser{
		UserName:  res.UserName,
		FirstName: res.FirstName,
		LastName:  res.LastName,
		Bio:       res.Bio,
		Website:   res.WebSite,
	})
}

// GetAllUsers ...
// @Security ApiKeyAuth
// @Summary GetAllUsers
// @Description Viewing a single User by id
// @Tags user
// @Accept json
// @Produce json
// @Param page query int true "Page number"
// @Param limit query int true "Items per page"
// @Success 200 {object} models.UsersRes
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/user/get_all_users [get]
func (h *HandlerV1) GetAllUsers(c *gin.Context) {
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
	var res models.Users
	for _, user := range response.Users {
		r := models.User{
			UserName:  user.UserName,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Role:      user.Role,
			Bio:       user.Bio,
			Website:   user.WebSite,
		}

		res.Users = append(res.Users, r)
	}
	c.JSON(http.StatusOK, res)
}

func (h HandlerV1) UpdateEmail(c *gin.Context) {

}
