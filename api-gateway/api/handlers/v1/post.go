package v1

import (
	"api-gateway/api/handlers/models/commentModel"
	"api-gateway/api/handlers/models/postModel"
	token2 "api-gateway/api/tokens"
	cfg "api-gateway/config"
	commentPb "api-gateway/genproto/comment"
	postPb "api-gateway/genproto/post"
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

// CreatePost ...
// @Security ApiKeyAuth
// @Summary CreatePost
// @Description Viewing a single User by id
// @Tags post
// @Accept json
// @Produce json
// @Param CreatePost body postModel.CreateReq true "create post"
// @Success 200 {object} postModel.Post
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/post/create [post]
func (h *HandlerV1) CreatePost(c *gin.Context) {
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
		body        postModel.CreateReq
		jsonMarshal protojson.MarshalOptions
	)
	jsonMarshal.UseProtoNames = true

	err = c.ShouldBindJSON(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("Failed to bind json", l.Error(err))
		return
	}
	body.UserID = cast.ToString(id)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()
	postId := uuid.New().String()
	res, err := h.serviceManager.PostService().Create(ctx, &postPb.CreateReq{
		Id:       postId,
		Title:    body.Title,
		Content:  body.Content,
		UserId:   body.UserID,
		Category: body.Category,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("message", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, &postPb.Post{
		Id:        res.Id,
		Title:     res.Title,
		Content:   res.Content,
		UserId:    res.UserId,
		Category:  res.Category,
		CreatedAt: res.CreatedAt,
	})

}

// GetPost ...
// @Security ApiKeyAuth
// @Summary GetPost
// @Description Viewing a single User by id
// @Tags post
// @Accept json
// @Produce json
// @Param post_id query string true "post_id"
// @Success 200 {object} postModel.Post
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/post/get_post [get]
func (h *HandlerV1) GetPost(c *gin.Context) {
	post_id := c.Query("post_id")

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

	post, err := h.serviceManager.PostService().GetPost(ctx, &postPb.GetReq{
		PostId: post_id,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("message", l.Error(err))
		return
	}

	_, err = h.serviceManager.PostService().Views(ctx, &postPb.ViewReq{
		PostId: post_id,
		UserId: cast.ToString(id),
	})

	c.JSON(http.StatusOK, &postModel.Post{
		Id:        post.Id,
		Title:     post.Title,
		Content:   post.Content,
		UserId:    post.UserId,
		Category:  post.Category,
		Likes:     post.Likes,
		Dislikes:  post.Dislikes,
		Views:     post.Views,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	})
}

// GetPostWithComment ...
// @Security ApiKeyAuth
// @Summary GetPostWithComment
// @Description Viewing a single User by id
// @Tags post
// @Accept json
// @Produce json
// @Param post_id query string true "post_id"
// @Param page_comment query string true "page_comment"
// @Param limit_comment query string true "limit_comment"
// @Success 200 {object} postModel.Post
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/post/get_with_comment [get]
func (h *HandlerV1) GetPostWithComment(c *gin.Context) {
	post_id := c.Query("post_id")
	comment_page := c.DefaultQuery("page_comment", "1")
	comment_limit := c.DefaultQuery("limit_comment", "10")

	page, err := strconv.Atoi(comment_page)
	if err != nil || page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid page number",
		})
		h.log.Error("invalid page number", l.Error(err))
		return
	}

	limit, err := strconv.Atoi(comment_limit)
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid limit",
		})
		h.log.Error("invalid limit", l.Error(err))
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
	}
	id := claims["id"]

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	post, err := h.serviceManager.PostService().GetPost(ctx, &postPb.GetReq{
		PostId: post_id,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("message", l.Error(err))
		return
	}

	comments, err := h.serviceManager.CommentService().GetCommentsByPostId(ctx, &commentPb.GetByPostIdReq{
		PostId: post_id,
		Limit:  int64(limit),
		Page:   int64(page),
	})

	var com postModel.Post
	for _, comment1 := range comments.Comments {
		comm := commentModel.Comment{
			CommentID: comment1.CommentId,
			PostID:    comment1.PostId,
			UserID:    comment1.UserId,
			Content:   comment1.Content,
			Likes:     comment1.Likes,
		}
		com.Comments = append(com.Comments, comm)
	}

	_, err = h.serviceManager.PostService().Views(ctx, &postPb.ViewReq{
		PostId: post_id,
		UserId: cast.ToString(id),
	})

	c.JSON(http.StatusOK, &postModel.Post{
		Id:        post.Id,
		Title:     post.Title,
		Content:   post.Content,
		UserId:    post.UserId,
		Category:  post.Category,
		Likes:     post.Likes,
		Dislikes:  post.Dislikes,
		Views:     post.Views,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
		Comments:  com.Comments,
	})
}

// GetPostByOwnerId ...
// @Security ApiKeyAuth
// @Summary GetPostByOwnerId
// @Description Viewing a single User by id
// @Tags post
// @Accept json
// @Produce json
// @Param page query int true "Page number"
// @Param limit query int true "Items per page"
// @Success 200 {object} postModel.Posts
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/post/get_by_owner_id [get]
func (h HandlerV1) GetPostByOwnerId(c *gin.Context) {
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

	posts, err := h.serviceManager.PostService().SearchPost(ctx, &postPb.SearchReq{
		Field: "user_id",
		Value: cast.ToString(id),
		Page:  int32(page),
		Limit: int32(limit),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("message", l.Error(err))
		return
	}
	var postss postModel.Posts
	for _, post := range posts.Posts {
		_, err = h.serviceManager.PostService().Views(ctx, &postPb.ViewReq{
			PostId: post.Id,
			UserId: cast.ToString(id),
		})
		if err != nil {
			c.JSON(404, gin.H{
				"error": err.Error(),
			})
			h.log.Error("message", l.Error(err))
			return
		}
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
			UpdatedAt: post.UpdatedAt,
		}
		postss.Posts = append(postss.Posts, ps)

	}

	c.JSON(http.StatusOK, posts)
}

// SearchPost ...
// @Security ApiKeyAuth
// @Summary SearchPost
// @Description Viewing a single User by id
// @Tags post
// @Accept json
// @Produce json
// @Param SearchPost body postModel.SearchReq true "search post"
// @Success 200 {object} postModel.SearchReq
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/post/search_post [post]
func (h *HandlerV1) SearchPost(c *gin.Context) {
	var (
		body        postModel.SearchReq
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
	}
	id := claims["id"]

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	posts, err := h.serviceManager.PostService().SearchPost(ctx, &postPb.SearchReq{
		Field: body.Field,
		Value: body.Value,
		Page:  body.Page,
		Limit: body.Limit,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("message", l.Error(err))
		return
	}
	var postss postModel.Posts
	for _, post := range posts.Posts {
		_, err = h.serviceManager.PostService().Views(ctx, &postPb.ViewReq{
			PostId: post.Id,
			UserId: cast.ToString(id),
		})
		if err != nil {
			c.JSON(404, gin.H{
				"error": err.Error(),
			})
			h.log.Error("message", l.Error(err))
			return
		}
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
			UpdatedAt: post.UpdatedAt,
		}
		postss.Posts = append(postss.Posts, ps)

	}

	c.JSON(http.StatusOK, postss.Posts)
}

// UpdatePost ...
// @Security ApiKeyAuth
// @Summary UpdatePost
// @Description Viewing a single User by id
// @Tags post
// @Accept json
// @Produce json
// @Param UpdatePost body postModel.UpdatePostReq true "update post"
// @Success 200 {object} postModel.Post
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/post/update_post [put]
func (h *HandlerV1) UpdatePost(c *gin.Context) {
	var (
		body        postModel.UpdatePostReq
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	res, err := h.serviceManager.PostService().UpdatePost(ctx, &postPb.UpdatePostReq{
		Id:       body.ID,
		Title:    body.Title,
		Content:  body.Content,
		Category: body.Category,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("message", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, &postModel.Post{
		Id:       res.Id,
		Title:    res.Title,
		Content:  res.Content,
		UserId:   res.UserId,
		Category: res.Category,
	})
}

// DeletePost ...
// @Security ApiKeyAuth
// @Summary DeletePost
// @Description Viewing a single User by id
// @Tags post
// @Accept json
// @Produce json
// @Param post_id query string true "post_id"
// @Success 200 {object} postModel.Status
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/post/delete_post [delete]
func (h *HandlerV1) DeletePost(c *gin.Context) {
	post_id := c.Query("post_id")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	res, err := h.serviceManager.PostService().DeletePost(ctx, &postPb.DeletePostReq{
		Id: post_id,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("message", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, &postModel.Status{Status: res.Status})
}

// ClickLike ...
// @Security ApiKeyAuth
// @Summary ClickLike
// @Description Viewing a single User by id
// @Tags post
// @Accept json
// @Produce json
// @Param post_id query string true "post_id"
// @Success 200 {object} postModel.Status
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/post/click_like [post]
func (h *HandlerV1) ClickLike(c *gin.Context) {
	post_id := c.Query("post_id")
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

	res, err := h.serviceManager.PostService().PostClickLike(ctx, &postPb.ClickReq{
		PostId: post_id,
		UserId: cast.ToString(id),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("message", l.Error(err))
		return
	}
	c.JSON(http.StatusOK, &postModel.Status{Status: res.Like})

}

// ClickDisLike ...
// @Security ApiKeyAuth
// @Summary ClickDisLike
// @Description Viewing a single User by id
// @Tags post
// @Accept json
// @Produce json
// @Param post_id query string true "post_id"
// @Success 200 {object} postModel.Status
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/post/click_dislike [post]
func (h *HandlerV1) ClickDisLike(c *gin.Context) {
	post_id := c.Query("post_id")
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

	res, err := h.serviceManager.PostService().PostClickDisLike(ctx, &postPb.ClickReq{
		PostId: post_id,
		UserId: cast.ToString(id),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("message", l.Error(err))
		return
	}
	c.JSON(http.StatusOK, &postModel.Status{Status: res.Like})
}
