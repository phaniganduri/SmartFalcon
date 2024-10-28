package handlers

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    
    "github.com/gin-gonic/gin"
)

func (api *AssetManagementAPI) CreateAsset(c *gin.Context) {
    var asset Asset
    if err := c.BindJSON(&asset); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    if asset.MSISDN == "" || asset.DealerID == "" || asset.MPIN == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
        return
    }

    _, err := api.Contract.SubmitTransaction("CreateAsset", 
        asset.DealerID, 
        asset.MSISDN, 
        asset.MPIN, 
        fmt.Sprintf("%.2f", asset.Balance))

    if err != nil {
        log.Printf("Error creating asset: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create asset: %v", err)})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "Asset created successfully",
        "asset": asset,
    })
}

func (api *AssetManagementAPI) ReadAsset(c *gin.Context) {
    msisdn := c.Param("msisdn")
    if msisdn == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "MSISDN is required"})
        return
    }

    result, err := api.Contract.EvaluateTransaction("ReadAsset", msisdn)
    if err != nil {
        log.Printf("Error reading asset: %v", err)
        c.JSON(http.StatusNotFound, gin.H{"error": "Asset not found"})
        return
    }

    var asset Asset
    if err := json.Unmarshal(result, &asset); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse asset data"})
        return
    }

    c.JSON(http.StatusOK, asset)
}

func (api *AssetManagementAPI) UpdateAsset(c *gin.Context) {
    msisdn := c.Param("msisdn")
    if msisdn == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "MSISDN is required"})
        return
    }

    var updateData struct {
        Balance float64 `json:"balance"`
        Status  string  `json:"status"`
    }

    if err := c.BindJSON(&updateData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    _, err := api.Contract.SubmitTransaction("UpdateAsset", 
        msisdn,
        fmt.Sprintf("%.2f", updateData.Balance),
        updateData.Status)

    if err != nil {
        log.Printf("Error updating asset: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update asset: %v", err)})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Asset updated successfully"})
}

func (api *AssetManagementAPI) DeleteAsset(c *gin.Context) {
    msisdn := c.Param("msisdn")
    if msisdn == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "MSISDN is required"})
        return
    }

    _, err := api.Contract.SubmitTransaction("DeleteAsset", msisdn)
    if err != nil {
        log.Printf("Error deleting asset: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to delete asset: %v", err)})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Asset deleted successfully"})
}

func (api *AssetManagementAPI) GetAllAssets(c *gin.Context) {
    result, err := api.Contract.EvaluateTransaction("GetAllAssets")
    if err != nil {
        log.Printf("Error getting all assets: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get assets"})
        return
    }

    var assets []*Asset
    if err := json.Unmarshal(result, &assets); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse assets data"})
        return
    }

    c.JSON(http.StatusOK, assets)
}

func (api *AssetManagementAPI) GetAssetHistory(c *gin.Context) {
    msisdn := c.Param("msisdn")
    if msisdn == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "MSISDN is required"})
        return
    }

    result, err := api.Contract.EvaluateTransaction("GetAssetHistory", msisdn)
    if err != nil {
        log.Printf("Error getting asset history: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get asset history"})
        return
    }

    var history []Asset
    if err := json.Unmarshal(result, &history); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse history data"})
        return
    }

    c.JSON(http.StatusOK, history)
}