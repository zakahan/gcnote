// -------------------------------------------------
// Package middleware
// Author: hanzhi
// Date: 2024/12/10
// -------------------------------------------------

package middleware

import (
	"fmt"
	"testing"
)

func TestJwt(t *testing.T) {
	tokenString := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzM5MDQ5NTUsIm5hbWUiOiJoYW56aGkxIiwic3ViIjoiN2NhZjk1NjctOGQ2ZC00NGJlLTliZDMtOWZiNTlkY2EzMTE5In0.D6iYzVtanih898sj6zHBTVDU4m93UVtaNHRlPI_0jhwLwaexhW-wAZrWYfUOt2AKTpgQdaYKhVcqzXFXhunWnicEwAJOUotV8Kg4XNGKgPDJqA0_r02FyPxz32YP602xM4BbScYsRZxw4mkUoTQVAmREYTXv9yvafUuKl74mPrrQux72kLWipe49Y3XB8T6sJ-5ddK5wYQjnQ6iq8dTLAwTT1qcBvQPGtiYoaqJ7JU7bbZ912u4cINtZ8ZhblSpFfkFmw18isTeHaW3tyUGVBCS96SQkJjyjHGQOz_NL5nwPdJw7qmg-HqkR2hNfyqlTRGqA6b3JjOGr8Y6JOu5ahA"

	j := NewJWT()
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		fmt.Printf("jwt Parse err:%v \n", err)
		fmt.Println("登陆信息错误")

		return
	}
	fmt.Println("claims设置成功")
	fmt.Println(claims)
	//zap.S().Debugf(claims)
	//ctx.Set("claims", claims)
	//ctx.Next()
}
