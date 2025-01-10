// -------------------------------------------------
// Package middleware
// Author: hanzhi
// Date: 2024/12/9
// -------------------------------------------------

package middleware

import (
	"fmt"
	"gcnote/server/config"
	"gcnote/server/dto"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
)

const Secret = "AllYourBase"

type JWTInfo struct {
	privateKey []byte
	publicKey  []byte
}

func Jwt() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString := context.GetHeader("token")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(Secret), nil
		})
		if err != nil {
			zap.S().Errorf("jwt Parse err:%v", err)
			context.JSON(http.StatusOK, dto.Fail(dto.UserTokenErrCode))
			context.Abort()
			return
		}
		switch {
		case token.Valid:
			context.Next()
		default:
			zap.S().Errorf("jwt Parse err:%+v", err)
			context.JSON(http.StatusOK, dto.Fail(dto.UserTokenErrCode))
			context.Abort()
			return
		}
	}
}

func NewJWT() *JWTInfo {
	pathCfg := config.PathCfg
	privateKeyPath := pathCfg.JwtPrivateKeyPath
	publicKeyPath := pathCfg.JwtPublicKeyPath
	var err error
	var privateKey []byte
	var publicKey []byte

	// 读取私钥，注意路径问题
	privateKey, err = os.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatalf("failed to load private key: %v", err)
	}
	// 读取公钥
	publicKey, err = os.ReadFile(publicKeyPath)
	if err != nil {
		log.Fatalf("failed to load public key: %v", err)
	}
	return &JWTInfo{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

// GenerateJWT 签发 JWT
func (j *JWTInfo) GenerateJWT(claims jwt.MapClaims) (string, error) {
	// 解析 RSA 私钥
	rsaPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM(j.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %v", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// 使用私钥签名
	tokenString, err := token.SignedString(rsaPrivateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return tokenString, nil
}

// ParseToken 解析 JWT
func (j *JWTInfo) ParseToken(tokenString string) (jwt.MapClaims, error) {
	// 解析 RSA 公钥
	rsaPublicKey, err := jwt.ParseRSAPublicKeyFromPEM(j.publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}

	// 验证 token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 确认使用的是 RS256 签名方法
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return rsaPublicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %v", err)
	}

	// token 有效性校验
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}

// VerifyJWT 验证 JWT
func VerifyJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for name, headers := range ctx.Request.Header {
			for _, h := range headers {
				if name == "Token" {
					zap.S().Debugf("Header %v: %v\n", name, h)
				}
			}
		}
		tokenString := ctx.GetHeader("token")
		//zap.S().Debugf("token： %v", tokenString)
		if tokenString == "" {
			zap.S().Debugf("没有tokenString")
			ctx.JSON(http.StatusUnauthorized, map[string]string{
				"msg": "请登录",
			})
			ctx.Abort()
			return
		}

		j := NewJWT()
		claims, err := j.ParseToken(tokenString)

		zap.S().Infof("claims setting：%v", claims)
		if err != nil {
			zap.S().Errorf("jwt Parse err:%v", err)
			ctx.JSON(http.StatusUnauthorized, dto.Fail(dto.UserTokenErrCode))
			ctx.Abort()
			return
		}
		//zap.S().Debugf(claims)
		ctx.Set("claims", claims)
		ctx.Next()
	}
}
