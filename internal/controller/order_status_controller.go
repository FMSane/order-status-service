package controller

import (
	"net/http"
	"order-status-service/internal/dto"
	"order-status-service/internal/middleware"
	"order-status-service/internal/service"

	"github.com/gin-gonic/gin"
)

type OrderStatusController struct {
	Service     *service.OrderStatusService
	AuthService *service.AuthService
}

func NewOrderStatusController(router *gin.Engine, svc *service.OrderStatusService, authSvc *service.AuthService) {
	ctrl := &OrderStatusController{Service: svc, AuthService: authSvc}

	// Grupo con autenticación
	auth := router.Group("/status")
	auth.Use(middleware.AuthMiddleware(authSvc))

	auth.GET("", ctrl.GetStatusesByUser)
	auth.POST("", ctrl.CreateStatus)
	auth.PUT("/:id", ctrl.UpdateStatus)
	auth.GET("/all", ctrl.GetAllOrderStatuses)
	auth.GET("/filter", ctrl.FilterByStatus)

	// Rutas públicas (sin autenticación)
	// router.GET("/status/all", ctrl.GetAllStatuses)
	router.POST("/status/init", ctrl.InitStatus)
}

func (ctrl *OrderStatusController) GetAllStatuses(c *gin.Context) {
	statuses, err := ctrl.Service.GetAllStatuses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, statuses)
}

func (ctrl *OrderStatusController) GetAllOrderStatuses(c *gin.Context) {
	permissions := c.GetStringSlice("userPermissions")
	isAdmin := false
	for _, p := range permissions {
		if p == "admin" {
			isAdmin = true
			break
		}
	}
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: admin access required"})
		return
	}

	statuses, err := ctrl.Service.GetAllOrderStatuses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, statuses)
}

func (ctrl *OrderStatusController) GetStatusesByUser(c *gin.Context) {
	userID := c.GetString("userID")
	statuses, err := ctrl.Service.GetStatusesByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, statuses)
}

func (ctrl *OrderStatusController) CreateStatus(c *gin.Context) {
	permissions := c.GetStringSlice("userPermissions")
	isAdmin := false
	for _, p := range permissions {
		if p == "admin" {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: admin access required"})
		return
	}

	var req dto.CreateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status, err := ctrl.Service.CreateStatus(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, status)
}

func (ctrl *OrderStatusController) UpdateStatus(c *gin.Context) {
	permissions := c.GetStringSlice("userPermissions")
	isAdmin := false
	for _, p := range permissions {
		if p == "admin" {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: admin access required"})
		return
	}

	id := c.Param("id")
	var req dto.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status, err := ctrl.Service.UpdateStatus(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, status)
}

func (ctrl *OrderStatusController) FilterByStatus(c *gin.Context) {
	status := c.Query("status")
	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing status query param"})
		return
	}

	results, err := ctrl.Service.GetByStatus(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}

func (ctrl *OrderStatusController) InitStatus(c *gin.Context) {
	var req dto.CreateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Status = "Pendiente" // Estado inicial

	status, err := ctrl.Service.CreateStatus(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, status)
}
