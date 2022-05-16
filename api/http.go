package api

import (
	"fmt"
	"net/http"

	"github.com/AntonyIS/go-loco/app"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type LocomotiveHandler interface {
	GetLoco(*gin.Context)
	CreateLoco(*gin.Context)
	GetAllLoco(*gin.Context)
	UpdateLoco(*gin.Context)
	DeleteLoco(*gin.Context)
}
type handler struct {
	locomotiveService app.LocomotiveService
}

func NewHandler(locomotiveService app.LocomotiveService) LocomotiveHandler {
	return &handler{locomotiveService: locomotiveService}
}

func (lh *handler) CreateLoco(c *gin.Context) {
	loco := app.Locomotive{}

	if err := c.BindJSON(&loco); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": http.StatusInternalServerError,
			"message":     fmt.Sprintf("error: %s", err),
		})
	}

	if _, err := lh.locomotiveService.CreateLoco(&loco); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": http.StatusInternalServerError,
			"message":     fmt.Sprintf("error: %s", err),
		})
	}
	c.JSON(http.StatusOK, loco)

}
func (lh *handler) GetLoco(c *gin.Context) {
	loco_id := c.Param("loco_id")
	loco, err := lh.locomotiveService.GetLoco(loco_id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status_code": http.StatusNotFound,
			"message":     fmt.Sprintf("error: %s", err),
		})
	}
	c.JSON(http.StatusOK, loco)

}
func (lh *handler) GetAllLoco(c *gin.Context) {
	locos, err := lh.locomotiveService.GetAllLoco()

	if err != nil {
		if errors.Cause(err) == app.ErrorLocomotiveNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status_code": 404,
				"message":     fmt.Sprintf("error: %s", err),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": 500,
			"message":     fmt.Sprintf("error: %s", err),
		})
		return
	}
	c.JSON(http.StatusOK, locos)

}

func (lh *handler) UpdateLoco(c *gin.Context) {
	loco := &app.Locomotive{}
	loco_id := c.Param("loco_id")

	if err := c.ShouldBindJSON(&loco); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status_code": http.StatusBadRequest,
			"message":     fmt.Sprintf("error: %s", err),
		})
		return
	}

	loco.LocoID = loco_id
	loco, err := lh.locomotiveService.UpdateLoco(loco)
	if err != nil {
		if errors.Cause(err) == app.ErrorLocomotiveNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status_code": http.StatusNotFound,
				"message":     fmt.Sprintf("error: %s", err),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": http.StatusInternalServerError,
			"message":     fmt.Sprintf("error: %s", err),
		})
		return
	}
	c.JSON(http.StatusOK, loco)
}

func (lh *handler) DeleteLoco(c *gin.Context) {
	loco_id := c.Param("loco_id")
	err := lh.locomotiveService.DeleteLoco(loco_id)
	if err != nil {
		if errors.Cause(err) == app.ErrorLocomotiveNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status_code": http.StatusNotFound,
				"message":     fmt.Sprintf("error: %s", err),
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"status_code": http.StatusBadRequest,
			"message":     fmt.Sprintf("error: %s", err),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status_code": http.StatusOK,
		"message":     "item deleted successfully",
	})

}
