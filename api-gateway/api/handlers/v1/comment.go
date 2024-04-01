package v1

import (
	"api-gateway/api/handlers/models/commentModel"
	token2 "api-gateway/api/tokens"
	cfg "api-gateway/config"
	commentPb "api-gateway/genproto/comment"
	l "api-gateway/pkg/logger"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"google.golang.org/protobuf/encoding/protojson"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// CreateComment ...
// @Security ApiKeyAuth
// @Summary CreateComment
// @Description Viewing a single User by id
// @Tags comment
// @Accept json
// @Produce json
// @Param CreateComment body commentModel.CreateReq true "create comment"
// @Success 200 {object} commentModel.Comment
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/user/comment/create [post]
func (h *HandlerV1) CreateComment(c *gin.Context) {
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
	}
	id := claims["id"]

	var (
		body        commentModel.CreateReq
		jsonMarshal protojson.MarshalOptions
	)
	jsonMarshal.UseProtoNames = true

	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	body.CommentID = uuid.New().String()

	res, err := h.serviceManager.CommentService().Create(ctx, &commentPb.CreateReq{
		PostId:    body.PostID,
		UserId:    cast.ToString(id),
		Content:   body.Content,
		CommentId: body.CommentID,
	})

	if err != nil {
		c.JSON(404, gin.H{
			"error": err.Error(),
		})
		h.log.Error("message", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, &commentModel.Comment{
		CommentID: res.CommentId,
		PostID:    res.PostId,
		UserID:    res.UserId,
		Content:   res.Content,
		Likes:     0,
	})

}

// GetCommentsByPostId ...
// @Security ApiKeyAuth
// @Summary GetCommentsByPostId
// @Description Viewing a single User by id
// @Tags comment
// @Accept json
// @Produce json
// @Param post_id query string true "post_id"
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Success 200 {object} commentModel.Comments
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/user/comment/get_comm_by_post_id [get]
func (h *HandlerV1) GetCommentsByPostId(c *gin.Context) {
	post_id := c.Query("post_id")
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	res, err := h.serviceManager.CommentService().GetCommentsByPostId(ctx, &commentPb.GetByPostIdReq{
		PostId: post_id,
		Limit:  int64(limit),
		Page:   int64(page),
	})

	if err != nil {
		c.JSON(404, gin.H{
			"error": err.Error(),
		})
		h.log.Error("message", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, res.Comments)
}

// GetCommentsByOwnerId ...
// @Security ApiKeyAuth
// @Summary GetCommentsByOwnerId
// @Description Viewing a single User by id
// @Tags comment
// @Accept json
// @Produce json
// @Param post_id query string true "post_id"
// @Param page query string true "page"
// @Param limit query string true "limit"
// @Success 200 {object} commentModel.Comments
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/user/comment/get_comm_by_owner_id [get]
func (h *HandlerV1) GetCommentsByOwnerId(c *gin.Context) {
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
	}
	id := claims["id"]

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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	res, err := h.serviceManager.CommentService().GetCommentsByOwnerId(ctx, &commentPb.GetByOwnerIdReq{
		OwnerId: cast.ToString(id),
		Limit:   int64(limit),
		Page:    int64(page),
	})

	if err != nil {
		c.JSON(404, gin.H{
			"error": err.Error(),
		})
		h.log.Error("message", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, &commentPb.GetByIdCommentsRes{
		Comments: res.Comments,
	})
}

// UpdateComment ...
// @Security ApiKeyAuth
// @Summary UpdateComment
// @Description Viewing a single User by id
// @Tags comment
// @Accept json
// @Produce json
// @Param UpdateComment body commentModel.UpdateComment true "update comment"
// @Success 200 {object} commentModel.Comment
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/user/comment/update [put]
func (h *HandlerV1) UpdateComment(c *gin.Context) {
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
	}
	id := claims["id"]

	var (
		body        commentModel.UpdateComment
		jsonMarshal protojson.MarshalOptions
	)
	jsonMarshal.UseProtoNames = true

	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	res, err := h.serviceManager.CommentService().UpdateComment(ctx, &commentPb.UpdateCommentReq{
		CommentId:  body.CommentID,
		UserId:     cast.ToString(id),
		NewContent: body.NewContent,
	})

	if err != nil {
		c.JSON(404, gin.H{
			"error": err.Error(),
		})
		h.log.Error("message", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, &commentModel.Comment{
		CommentID: res.CommentId,
		PostID:    res.PostId,
		UserID:    res.UserId,
		Content:   res.Content,
		Likes:     res.Likes,
	})

}

// DeleteComment ...
// @Security ApiKeyAuth
// @Summary DeleteComment
// @Description Viewing a single User by id
// @Tags comment
// @Accept json
// @Produce json
// @Param comment_id query string true "comment_id"
// @Success 200 {object} commentModel.Status
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/user/comment/delete [delete]
func (h *HandlerV1) DeleteComment(c *gin.Context) {
	comment_id := c.Query("comment_id")

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
	}
	id := claims["id"]

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	res, err := h.serviceManager.CommentService().DeleteComment(ctx, &commentPb.DeleteCommentReq{
		CommentId: comment_id,
		UserId:    cast.ToString(id),
	})
	if err != nil {
		c.JSON(404, gin.H{
			"error": err.Error(),
		})
		h.log.Error("message", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, &commentModel.Status{Status: res.Status})

}

// CommentClickLike ...
// @Security ApiKeyAuth
// @Summary CommentClickLike
// @Description Viewing a single User by id
// @Tags comment
// @Accept json
// @Produce json
// @Param comment_id query string true "comment_id"
// @Success 200 {object} commentModel.Status
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/user/comment/click_like [post]
func (h *HandlerV1) CommentClickLike(c *gin.Context) {
	comment_id := c.Query("comment_id")

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
	}
	id := claims["id"]

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	res, err := h.serviceManager.CommentService().CommentClickLike(ctx, &commentPb.ClickReq{
		CommentId: comment_id,
		UserId:    cast.ToString(id),
	})

	if err != nil {
		c.JSON(404, gin.H{
			"error": err.Error(),
		})
		h.log.Error("message", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, &commentModel.Status{Status: res.Like})
}
