package product_controller

import (
	"go-postgres/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name": "Juang Sabit",
		"bio":  "do your best",
	})
}

func FindProducts(c *gin.Context) {
	var products []models.Product
	results := models.DB.Debug().Find(&products)
	if results.Error != nil {
		c.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "count": len(products), "data": products})
}

type Activity struct {
	gorm.Model
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Username    string `gorm:"type:varchar(255)" json:"username"`
	Information string `gorm:"type:varchar(255)" json:"information"`
	CreatedAt   *time.Time
}

func FindProductById(c *gin.Context) {
	var product models.Product
	var activity []Activity
	// var res interface{}

	id := c.Param("id")

	if err := models.DB.Debug().Unscoped().First(&product, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "error", "message": "Data tidak ditemukan"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
			return
		}
	}

	// models.DB.Debug().Unscoped() for ignore where deleteAt is null
	if err := models.DB.Debug().Unscoped().Select("activities.id, username, information, activities.created_at").Joins("JOIN users on users.id = activities.user_id").Where("table_name = ? AND table_id = ?", "product", id).Find(&activity).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "error", "message": "Data tidak ditemukan"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": product, "activity": activity})
}

func CreateProduct(c *gin.Context) {
	var product models.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	models.DB.Debug().Create(&product)

	str_id, _ := c.Get("user_id")
	user_id := str_id.(uint)

	SaveActivity(user_id, "create product", "product", product.ID)
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Product has been created"})
}

func UpdateProduct(c *gin.Context) {
	var product models.Product
	id := c.Param("id")

	if err := c.ShouldBindJSON(&product); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if models.DB.Debug().Model(&product).Where("id = ?", id).Updates(&product).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Product not found"})
		return
	}

	str_id, _ := c.Get("user_id")
	user_id := str_id.(uint)
	id_uint, _ := strconv.Atoi(id)

	SaveActivity(user_id, "update product", "product", uint(id_uint))
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Data berhasil diperbarui"})
}

func DeleteProduct(c *gin.Context) {

	var product models.Product

	id := c.Param("id")
	if models.DB.Debug().Delete(&product, id).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Data tidak ditemukan"})
		return
	}

	str_id, _ := c.Get("user_id")
	user_id := str_id.(uint)
	id_uint, _ := strconv.Atoi(id)

	SaveActivity(user_id, "deleted product", "product", uint(id_uint))
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Data berhasil dihapus"})
}

func SaveActivity(user_id uint, info string, tableName string, tableID uint) {
	obj := models.Activity{UserID: user_id, Information: info, TableName: tableName, TableID: tableID}
	models.DB.Debug().Create(&obj)
}
