package routers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func JWTAuthMiddlewareAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求中获取 JWT token
		tokenString, err := c.Cookie("jwt")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			return
		}
		msg, err := ParseToken(tokenString)
		// 报错或者不是管理员
		if err != nil || !msg.IsAdmin {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Set("username", msg.Username)
		c.Next()
	}
}

func LoadAdmin(r *gin.Engine) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	err = db.AutoMigrate(&Challenge{})
	if err != nil {
		panic("failed to migrate database")
	}
	admin := r.Group("/api/admin")
	admin.Use(JWTAuthMiddlewareAdmin())
	{
		// 添加题目
		admin.POST("/addChallenge", func(c *gin.Context) {
			var challenge Challenge
			// 绑定参数
			if err := c.ShouldBind(&challenge); err != nil {
				log.Println(err)
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "invalid request",
				})
				return
			}
			// 查询题目是否已经存在
			result := db.First(&challenge, "category = ? and title = ?", challenge.Category, challenge.Title)
			log.Println(result.Error)
			if result.Error != nil { //查询错误 --> 不存在？
				if result.RowsAffected == 0 { // 结果为0行 则创建新题目
					err = db.Create(&challenge).Error
					if err != nil { // 创建失败
						c.JSON(http.StatusInternalServerError, gin.H{
							"msg": "failed to create challenge",
						})
						return
					} else { // 创建成功
						c.JSON(http.StatusOK, gin.H{
							"msg": "success to create challenge",
						})
					}
				}

			} else {
				c.JSON(http.StatusOK, gin.H{
					"msg": "challenge already exists",
				})
				return
			}
		})
	}
}
