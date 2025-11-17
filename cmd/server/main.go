// main.go
package main

import (
	"context"
	"log"
	"os"

	"order-status-service/internal/controller"
	"order-status-service/internal/repository"
	"order-status-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found")
	}

	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB")

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("❌ Failed to connect to MongoDB: %v", err)
	}

	db := client.Database(dbName)

	// Inicializamos Gin y servicios base
	router := gin.Default()
	authService := service.NewAuthService()

	// Repositorios
	catalogRepo := repository.NewCatalogRepository(db)
	orderRepo := repository.NewOrderStatusRepository(db)

	// Servicios
	catalogService := service.NewCatalogService(catalogRepo)
	orderStatusService := service.NewOrderStatusService(orderRepo)
	catalogAdminService := service.NewCatalogAdminService(catalogRepo)

	// Precargar estados base (solo si no existen)
	if err := catalogService.SeedDefaultStatuses(); err != nil {
		log.Printf("⚠️ Error al crear estados base: %v", err)
	}

	// Controladores
	controller.NewOrderStatusController(router, orderStatusService, authService)
	controller.NewCatalogAdminController(router, catalogAdminService, authService)

	// Puerto
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("🚀 Order Status Service running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
