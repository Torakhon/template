package v1

import (
	"api-gateway/api/handlers/models"
	pbu "api-gateway/genproto/users"
	"api-gateway/pkg/etc"
	l "api-gateway/pkg/logger"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

// Register ...
// @Summary Register
// @Description Register - Api for registering users
// @Tags Register
// @Accept json
// @Produce json
// @Param Register body models.RegisterModelReq true "RegisterModelReq"
// @Success 200 {object} models.RegisterModelRes
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/auth/register/ [post]
func (h *HandlerV1) Register(c *gin.Context) {
	var (
		body        models.RegisterModelReq
		jsonMarshal protojson.MarshalOptions
	)
	jsonMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	rdb := redis.NewClient(&redis.Options{
		Addr: "redisdb:6379",
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("Failed to bind json", l.Error(err))
		return
	}
	body.Email = strings.TrimSpace(body.Email)
	body.Email = strings.ToLower(body.Email)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	existsEmail, err := h.serviceManager.UserService().CheckUniques(ctx, &pbu.CheckUniqReq{
		Field: "email",
		Value: body.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("Failed to check email uniques")
		return
	}
	if existsEmail.Code == 0 {
		c.JSON(http.StatusConflict, gin.H{
			"error": "This email already in use, please use another email address",
		})
		h.log.Error("failed to check email uniques", l.Error(err))
		return
	}

	byteDate, err := json.Marshal(&body)
	if err != nil {
		log.Fatalln(err)
	}

	err = rdb.Set(context.Background(), body.Email, existsEmail.Code, 0).Err()
	if err != nil {
		log.Fatalln(err)
	}

	err = rdb.Set(context.Background(), body.Email+cast.ToString(existsEmail.Code), byteDate, 0).Err()
	if err != nil {
		h.log.Error(err.Error())
		return
	}

	code := strconv.Itoa(int(existsEmail.Code))

	auth := smtp.PlainAuth("test email", "torakhonoffical@gmail.com", "nfjdryidumjowyyv", "smtp.gmail.com")
	err = smtp.SendMail("smtp.gmail.com:587", auth, "torakhonoffical@gmail.com", []string{body.Email}, []byte(code))
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, true)
}

// Authorization ...
// @Summary Authorization
// @Description Authorization - Api for registering users
// @Tags Register
// @Accept json
// @Produce json
// @Param Register body models.AuthorizationReq true "RegisterModelReq"
// @Failure 200 {object} models.Authorization
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/auth/authorization/ [post]
func (h *HandlerV1) Authorization(c *gin.Context) {
	var (
		body        models.AuthorizationReq
		jsonMarshal protojson.MarshalOptions
	)

	jsonMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)

	rdb := redis.NewClient(&redis.Options{
		Addr: "redisdb:6379",
	})
	defer func(rdb *redis.Client) {
		err := rdb.Close()
		if err != nil {
			h.log.Error(err.Error())
			return
		}
	}(rdb)

	respCode, err := rdb.Get(context.Background(), body.Email).Result()
	if err != nil {
		h.log.Error(err.Error())
	}

	var code int
	if err := json.Unmarshal([]byte(respCode), &code); err != nil {
		h.log.Error(err.Error())
		return
	}
	code1 := body.Code

	if code != code1 {
		c.JSON(http.StatusBadRequest, false)
		return
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
		defer cancel()

		err = rdb.Del(context.Background(), body.Email).Err()
		if err != nil {
			h.log.Error(err.Error())
			return
		}
		var regis models.RegisterModelReq
		respUser, err := rdb.Get(context.Background(), body.Email+cast.ToString(code)).Result()
		if err != nil {
			h.log.Error(err.Error())
		}
		err = rdb.Del(context.Background(), body.Email+cast.ToString(code)).Err()
		if err != nil {
			h.log.Error(err.Error())
			return
		}
		if err := json.Unmarshal([]byte(respUser), &regis); err != nil {
			h.log.Error(err.Error())
			return
		}

		regis.Id = uuid.New().String()

		access, _, err := h.jwthandler.GenerateAuthJWT(regis.Id, "user")
		hashPassword, err := etc.HashPassword(regis.Password)
		if err != nil {
			h.log.Error(err.Error())
			return
		}
		_, err = h.serviceManager.UserService().CreateUser(ctx, &pbu.CreateUserReq{
			Id:        regis.Id,
			UserName:  regis.UserName,
			FirstName: regis.LastName,
			LastName:  regis.LastName,
			Email:     regis.Email,
			Password:  hashPassword,
			Role:      "user",
			Bio:       regis.Bio,
			WebSite:   regis.Website,
		})

		if err != nil {
			h.log.Error(err.Error())
			return
		}
		c.JSON(http.StatusOK, &models.Authorization{
			Token:  access,
			Status: true,
		})
	}
}

// Login ...
// @Summary Login
// @Description Login - Api for registering users
// @Tags Register
// @Accept json
// @Produce json
// @Param Login body models.LoginReq true "Login Req"
// @Success 200 {object} models.RegisterRes
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/auth/Login/ [post]
func (h *HandlerV1) Login(c *gin.Context) {
	var (
		body        models.LoginReq
		jsonMarshal protojson.MarshalOptions
	)
	jsonMarshal.UseProtoNames = true
	err := c.ShouldBindJSON(&body)

	if err != nil {
		h.log.Error(err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	user, err := h.serviceManager.UserService().Login(ctx, &pbu.LoginReq{
		Email: body.Email,
	})
	psw := etc.CheckPasswordHash(body.Password, user.Password)
	if !psw {
		h.log.Error("password or email error")
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "password or email error",
		})
		return
	}

	access, _, err := h.jwthandler.GenerateAuthJWT(user.Id, user.Role)
	c.JSON(http.StatusOK, models.LoginRes{AccessToken: access})
}
