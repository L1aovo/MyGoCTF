package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type User struct {
	gorm.Model
	ID       uint64 `gorm:"primaryKey;autoIncrement"`
	Username string `gorm:"not null" form:"username" json:"username" binding:"required"`
	Password string `gorm:"not null" form:"password" json:"password" binding:"required"`
	Email    string `gorm:"not null;default:''" form:"email" json:"email" binding:"required"`
	IsAdmin  bool   `gorm:"not null;default:false"`
}

type UserLogin struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func LoadUser(r *gin.Engine) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	err = db.AutoMigrate(&User{})
	if err != nil {
		panic("failed to migrate database")
	}
	apiGroup := r.Group("/api/user")
	{
		// 个人信息
		apiGroup.GET("/profile", JWTAuthMiddleware(), func(c *gin.Context) {
			// 获取参数
			name := c.MustGet("username").(string)
			var user User
			// 查询用户是否存在
			result := db.First(&user, "username = ?", name)

			if result.Error != nil { // 查询错误
				c.JSON(http.StatusOK, gin.H{
					"msg": "failed to get user",
				})
				return
			} else { // 查询成功
				c.JSON(http.StatusOK, gin.H{
					"msg":  "success",
					"user": user,
				})
				return
			}
		})
		// 注册
		apiGroup.POST("/register", func(c *gin.Context) {
			var user User
			// 绑定参数
			if err := c.ShouldBind(&user); err != nil {
				log.Println(err)
				c.JSON(http.StatusOK, gin.H{
					"msg": "invalid request",
				})
				return
			}

			// 查询用户是否已经存在
			result := db.First(&user, "username = ?", user.Username)
			log.Println(result.Error)
			if result.Error != nil { //查询错误 --> 不存在？
				if result.RowsAffected == 0 { // 结果为0行 则创建新用户
					err = db.Create(&user).Error
					if err != nil { // 创建失败
						c.JSON(http.StatusInternalServerError, gin.H{
							"msg": "failed to create user",
						})
						return
					} else { // 创建成功
						c.JSON(http.StatusOK, gin.H{
							"msg": "register success",
						})
						return
					}

				} else { // 结果不为0行 则查询失败
					c.JSON(http.StatusOK, gin.H{
						"msg": "user already exists",
					})
					return
				}
			}
			// 查询成功 --> 用户已经存在
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "user already exists",
			})
		})
		// 登陆
		apiGroup.POST("/login", func(c *gin.Context) {
			var user User
			var userLogin UserLogin
			if err := c.ShouldBind(&userLogin); err != nil {
				log.Println(err)
				c.JSON(http.StatusOK, gin.H{
					"msg": "invalid request",
				})
				return
			}
			// 验证账号密码
			result := db.First(&user, "username = ? AND password = ?", userLogin.Username, userLogin.Password)
			if result.Error != nil {
				if result.RowsAffected == 0 {
					c.JSON(http.StatusOK, gin.H{
						"msg": "username or password error",
					})
					return
				}
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "failed to query user",
				})
				return
			} else {
				fmt.Println(user)
				value, _ := GenToken(user)
				cookie := &http.Cookie{
					Name:     "jwt",
					Value:    value,
					HttpOnly: true,
					Path:     "/",
				}
				http.SetCookie(c.Writer, cookie)
				c.JSON(http.StatusOK, gin.H{
					"msg": "login success",
				})
				return
			}
		})
		// 登出
		apiGroup.POST("/logout", JWTAuthMiddleware(), func(c *gin.Context) {
			cookie := &http.Cookie{
				Name:     "jwt",
				Value:    "",
				HttpOnly: true,
				MaxAge:   -1,
				Path:     "/",
			}
			http.SetCookie(c.Writer, cookie)
			c.JSON(http.StatusOK, gin.H{
				"msg": "logout success",
			})
			return
		})
		// 更新用户信息
		apiGroup.POST("/update", JWTAuthMiddleware(), func(c *gin.Context) {
			// 获取参数
			var user User
			name := c.MustGet("username").(string)
			// 查询用户是否存在
			result := db.First(&user, "username = ?", name)
			// 绑定更新参数
			if err := c.ShouldBind(&user); err != nil {
				c.JSON(http.StatusOK, gin.H{
					"msg": "invalid request",
				})
				return
			}
			if result.Error != nil { // 查询错误
				c.JSON(http.StatusOK, gin.H{
					"msg": "failed to get user",
				})
				return
			} else { // 查询成功
				// 更新用户信息
				err = db.Debug().Model(&user).Where("username = ?", name).Updates(&user).Error
				if err != nil {
					log.Println(err)
					c.JSON(http.StatusOK, gin.H{
						"msg": "failed to update user",
					})
					return
				} else {
					c.JSON(http.StatusOK, gin.H{
						"msg": "update success",
					})
					return
				}
			}
		})
	}
}
