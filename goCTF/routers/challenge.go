package routers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

// Challenge 结构体Id为唯一标识，自增
type Challenge struct {
	gorm.Model
	ID       uint64 `gorm:"primaryKey;autoIncrement"`
	Category string `gorm:"not null" form:"category" json:"category" binding:"required"`
	Title    string `gorm:"not null" form:"title" json:"title" binding:"required"`
	Content  string `gorm:"not null" form:"content" json:"content" binding:"required"`
	Flag     string `gorm:"not null" form:"flag" json:"flag" binding:"required"`
}

type Solved struct {
	gorm.Model
	UserID      uint64    // 用户ID
	ChallengeID uint64    // 题目ID
	Submitted   bool      // 是否提交
	Passed      bool      // 是否通过
	SubmittedAt time.Time // 提交时间
}

type FlagRequest struct {
	Flag string `json:"flag" binding:"required"`
}

func LoadChallenge(r *gin.Engine) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	err = db.AutoMigrate(&Challenge{})
	err = db.AutoMigrate(&Solved{})
	if err != nil {
		panic("failed to migrate database")
	}
	challenge := r.Group("/api/challenge")
	challenge.Use(JWTAuthMiddleware())
	{
		challenge.GET("/list", func(c *gin.Context) {
			//查询数据库并返回所有题目
			var challenges []Challenge
			result := db.Select("id", "category", "title", "content").Find(&challenges)
			if result.Error != nil {
				c.JSON(200, gin.H{
					"code": 400,
					"msg":  "failed",
					"data": "challenges not found",
				})
				return
			} else {
				c.JSON(200, gin.H{
					"code": 200,
					"msg":  "success",
					"data": challenges,
				})
				return
			}
		})
		challenge.POST("/submit", func(c *gin.Context) {
			//提交flag
			var challenge Challenge
			var flag FlagRequest
			log.Println(flag.Flag)
			//获取flag
			if err := c.ShouldBindJSON(&flag); err != nil {
				log.Println(err)
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "invalid request",
				})
				return
			}
			result := db.Where("flag = ?", flag.Flag).First(&challenge)
			if result.Error != nil {
				c.JSON(200, gin.H{
					"code": 400,
					"msg":  "failed",
					"data": "flag incorrect",
				})
				return
			} else {
				// 记录用户和题目的对应关系
				var solved Solved
				solved.UserID = c.MustGet("userId").(uint64)
				solved.ChallengeID = challenge.ID
				solved.Submitted = true
				solved.Passed = true
				solved.SubmittedAt = time.Now()
				db.Create(&solved)

				c.JSON(200, gin.H{
					"code": 200,
					"msg":  "success",
					"data": "flag correct",
				})
				return
			}
		})
	}
}
