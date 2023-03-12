package routers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"time"
)

type CustomClaims struct {
	// 可根据需要自行添加字段
	Username             string `json:"username"`
	IsAdmin              bool   `json:"is_admin"`
	UserId               uint64 `json:"user_id"`
	jwt.RegisteredClaims        // 内嵌标准的声明
}

const TokenExpireDuration = time.Hour * 24

// CustomSecret 用于加盐的字符串
var CustomSecret = []byte("夏天夏天悄悄过去")

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求中获取 JWT token
		tokenString, err := c.Cookie("jwt")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			return
		}
		msg, err := ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Set("username", msg.Username)
		c.Set("userId", msg.UserId)
		c.Next()
	}
}

// GenRegisteredClaims 使用默认声明创建jwt
// GenToken 生成JWT
func GenToken(username User) (string, error) {
	// 创建一个我们自己的声明
	claims := CustomClaims{
		username.Username, // 自定义字段
		username.IsAdmin,
		username.ID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
			Issuer:    "my-project", // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(CustomSecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*CustomClaims, error) {
	// 解析token
	// 如果是自定义Claim结构体则需要使用 ParseWithClaims 方法
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		// 直接使用标准的Claim则可以直接使用Parse方法
		//token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
		return CustomSecret, nil
	})
	if err != nil {
		return nil, err
	}
	// 对token对象中的Claim进行类型断言
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
