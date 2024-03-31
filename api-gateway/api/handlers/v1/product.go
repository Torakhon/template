package v1

import (
	"api-gateway/api/handlers/models"
	pb "api-gateway/genproto/product"
	l "api-gateway/pkg/logger"
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
	"net/http"
	"time"
)

// CreateProduct ...
// @Security ApiKeyAuth
// @Summary Product
// @Tags Product
// @Accept json
// @Produce json
// @Param CreateProduct body models.Product true "Create Product"
// @Success 200 {object} models.Product
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/user/create_product [comment]
func (h *HandlerV1) CreateProduct(c *gin.Context) {
	var (
		body        models.Product
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
	product, err := h.serviceManager.ProductService().CreateProduct(ctx, &pb.CreateProductReq{
		Id:     body.Id,
		Name:   body.Name,
		Prays:  body.Prays,
		Amount: body.Amount,
	})
	if err != nil {
		h.log.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, &models.Product{
		Id:     product.Id,
		Name:   product.Name,
		Prays:  product.Prays,
		Amount: product.Amount,
	})
}

// GetProduct ...
// @Security ApiKeyAuth
// @Summary GetProduct
// @Tags Product
// @Accept json
// @Produce json
// @Param name path string true "name"
// @Success 200 {object} models.Product
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/user/get_product [get]
func (h *HandlerV1) GetProduct(c *gin.Context) {
	name := c.Param("name")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()
	product, err := h.serviceManager.ProductService().GetProduct(ctx, &pb.GetProductReq{
		Name: name,
	})
	if err != nil {
		h.log.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, &pb.Product{
		Id:     product.Id,
		Name:   product.Name,
		Prays:  product.Prays,
		Amount: product.Amount,
	})
}
