package main

import (
	"github.com/gin-gonic/gin"
	"goCTF/routers"
)

func main() {
	//myJwt, err := routers.GenToken("admin")
	//if err != nil {
	//	panic("failed to generate token")
	//}
	//fmt.Println(myJwt)
	//myJwtFake := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWlubm5uIiwiaXNzIjoibXktcHJvamVjdCIsImV4cCI6MTY3ODQyNTU1Mn0.j_rHPNnOZfn46BtloQqI1AJ6RJvkerOn54ylhqRZ718"
	//aa, err := routers.ParseToken(myJwtFake)
	//if err != nil {
	//	panic("failed to Parse token")
	//}
	//fmt.Printf("%v", aa.Username)
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	routers.LoadUser(r)
	routers.LoadAdmin(r)
	routers.LoadChallenge(r)
	if err := r.Run(":9090"); err != nil {
		panic("failed to start server")
	}
}
