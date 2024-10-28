package main

import (
    "log"
    "github.com/gin-gonic/gin"
    "asset-management-api/api/handlers"  
)

func main() {
    log.Println("Starting Asset Management API...")
    
    // Initialize 
    api, err := handlers.SetupAPI()
    if err != nil {
        log.Fatalf("Failed to set up API: %v", err)
    }

    router := gin.Default()
    
    // Add CORS middleware
    router.Use(func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    })

    // Define routes
    router.POST("/assets", api.CreateAsset)
    router.GET("/assets/:msisdn", api.ReadAsset)
    router.PUT("/assets/:msisdn", api.UpdateAsset)
    router.DELETE("/assets/:msisdn", api.DeleteAsset)
    router.GET("/assets", api.GetAllAssets)
    router.GET("/assets/:msisdn/history", api.GetAssetHistory)

    log.Println("Server starting on port 8080...")
    router.Run(":8080")
}