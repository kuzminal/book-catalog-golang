package secure

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"strings"
)

func checkToken(r *http.Request) error {
	authHeader := r.Header.Get("Authorization")
	tokenString := strings.ReplaceAll(authHeader, "Bearer ", "")
	log.Println(tokenString)
	//tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.u1riaD1rW97opCoAuRCTy4w58Br-Zk-bh7vLiRIsrpU"

	key, err := jwt.ParseRSAPublicKeyFromPEM([]byte("MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAhAZou/UbaPU5O7uSjvS4CmtN6Dk9bY/MlwwvtQ5IjsgduPiRWz4gQtpp6LiG9yvkGnKoQOXYB63N/7sNoqUeMB/AIICY4blFDX+/mWs4n/uGa3APOIItkqLz4E4Dix4UmPxSjd5qg73GjP4yPTH9VQq5kfzcw3ohHGk9RrpeUEE3wmB93uOunNOSLnDHnY/4Ssy8/uKY6Ua6T3dDWLir7EApyPlhlfHbgrWd6vsMIuDBiUwVYCvqtcBFbD1gSWpk0j84CFhqzryrCzECklNCjQBrCH98gxJBrM4zcKWeB8uVKzRa6Qc5tTLasBE9nwzh4aKfDQaGkkzKXhseWXmkAwIDAQAB"))
	if err != nil {
		return fmt.Errorf("validate: parse key: %w", err)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return key, nil
	})

	_, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		fmt.Println(err)
	}
	return nil
}
