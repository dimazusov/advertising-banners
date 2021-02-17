package internalhttp

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/dimazusov/hw-test/advertising-banners/internal/domain"
	"github.com/gin-gonic/gin"
)

func AddBannerToPlaceHandler(c *gin.Context, app Application) {
	var params addBannerToPlaceParams
	err := c.ShouldBindJSON(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong params"})
		return
	}

	err = app.AddBannerToPlace(context.Background(), params.BannerID, params.PlaceID)
	if err != nil {
		if err = app.LogError(err); err != nil {
			log.Println(err)
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "system error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func DeleteBannerFromPlaceHandler(c *gin.Context, app Application) {
	var params deleteBannerToPlaceParams
	err := c.ShouldBindJSON(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong params"})
		return
	}

	err = app.DeleteBannerFromPlace(context.Background(), params.BannerID, params.PlaceID)
	if err != nil {
		if err = app.LogError(err); err != nil {
			log.Println(err)
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "system error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func CreateEventHandler(c *gin.Context, app Application) {
	var params createEventParams
	err := c.ShouldBindJSON(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong params"})
		return
	}

	event := domain.Event{
		Type:       params.Type,
		BannerID:   params.BannerID,
		PlaceID:    params.PlaceID,
		SocGroupID: params.SocGroupID,
		Time:       uint(time.Now().Unix()),
	}
	err = app.AddEvent(context.Background(), event)
	if err != nil {
		if err = app.LogError(err); err != nil {
			log.Println(err)
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "system error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func BannerHandler(c *gin.Context, app Application) {
	var params bannerParams
	err := c.ShouldBindJSON(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong params"})
		return
	}

	bannerID, err := app.GetBannerForShow(context.Background(), params.PlaceID, params.SocGroupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "system error"})
	}

	c.JSON(http.StatusOK, gin.H{
		"bannerId": bannerID,
	})
}
