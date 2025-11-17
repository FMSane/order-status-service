// catalog_admin_controller.go
package controller

import (
	"log"
	"net/http"
	"order-status-service/internal/middleware"
	"order-status-service/internal/service"
	"strings"

	"github.com/gin-gonic/gin"
)

type CatalogAdminController struct {
	Service     *service.CatalogAdminService
	AuthService *service.AuthService
}

func NewCatalogAdminController(router *gin.Engine, svc *service.CatalogAdminService, auth *service.AuthService) {
	ctrl := &CatalogAdminController{Service: svc, AuthService: auth}

	group := router.Group("/admin/status/catalog")
	group.Use(middleware.AuthMiddleware(ctrl.AuthService))
	group.Use(middleware.AdminOnly())
	{
		group.GET("", ctrl.GetAll)
		group.POST("", ctrl.CreateStatus)
		group.GET("/:id", ctrl.GetByID)
	}
}

// POST /admin/status/catalog
func (ctrl *CatalogAdminController) CreateStatus(c *gin.Context) {
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

	token := strings.TrimPrefix(authHeader, "Bearer ")
	token = strings.TrimSpace(token)

	user, err := ctrl.AuthService.ValidateToken(token)

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

func (ctrl *CatalogAdminController) GetByID(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	token = strings.TrimSpace(token)

	log.Println("Token received in order-status-service:", c.GetHeader("Authorization"))

	user, err := ctrl.AuthService.ValidateToken(token)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	if !ctrl.AuthService.IsAdmin(user) {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}

	id := c.Param("id")

	result, err := ctrl.Service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "status not found"})
		return
	}

	c.JSON(http.StatusOK, result)
}
