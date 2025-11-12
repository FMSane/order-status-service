package controller

import (
	"net/http"
	"order-status-service/internal/service"

	"github.com/gin-gonic/gin"
)

type CatalogAdminController struct {
	Service     *service.CatalogAdminService
	AuthService *service.AuthService
}

func NewCatalogAdminController(router *gin.Engine, svc *service.CatalogAdminService, auth *service.AuthService) {
	ctrl := &CatalogAdminController{Service: svc, AuthService: auth}

	group := router.Group("/admin/status/catalog")
	{
		group.GET("", ctrl.GetAll)
		group.POST("", ctrl.CreateStatus)
	}
}

// POST /admin/status/catalog
func (ctrl *CatalogAdminController) CreateStatus(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}

	user, err := ctrl.AuthService.ValidateToken(authHeader)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	if !ctrl.AuthService.IsAdmin(user) {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}

	var body struct {
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&body); err != nil || body.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	if err := ctrl.Service.CreateStatus(body.Name); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "status added to catalog"})
}

// GET /admin/status/catalog
func (ctrl *CatalogAdminController) GetAll(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}

	user, err := ctrl.AuthService.ValidateToken(authHeader)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	if !ctrl.AuthService.IsAdmin(user) {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}

	statuses, err := ctrl.Service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, statuses)
}
