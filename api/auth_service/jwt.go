package auth_service

import (
	"fmt"
	"io/ioutil"
	"crypto/rsa"
	"net/http"
	"time"
	"strings"
	"github.com/dgrijalva/jwt-go"
)

// The secret key used to sign and verify the JWT
var secretKey []byte

// The public key used to verify the JWT
var publicKey *rsa.PublicKey

// Initialize the secret key and public key
//CALL THIS AT THE INITIALIZATION
func init() {
	// Read the secret key from a file
	key, err := ioutil.ReadFile("private.pem")
	if err != nil {
		panic(err)
	}
	secretKey = key

	// Read the public key from a file
	pem, err := ioutil.ReadFile("public.pem")
	if err != nil {
		panic(err)
	}
	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(pem)
	if err != nil {
		panic(err)
	}
}

// Sign a JWT using the HS256 algorithm
// Generate the claims and pass them over here (claims like username(must) etc)
func signJWT(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// Verify a JWT using the HS256 algorithm
func verifyJWT(tokenString string) (jwt.MapClaims, int) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, 404
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, 200
	}
	return nil,404
}


// A handler function that signs a JWT and sends it to the client
func SignHandler(username string) (string,error){
	// Set the claims for the JWT
	claims := jwt.MapClaims{
		"sub": username,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	// Sign the JWT
	token, err := signJWT(claims)
	if err != nil {
		return "",err
	}

	return token,nil

}


// A handler function that verifies a JWT sent by the client
func VerifyHandler(r *http.Request) (string,int){
	// Get the JWT from the request
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	tokenString := splitToken[1]

	// Verify the JWT
	claims, status := verifyJWT(tokenString)
	if status != 200 {
		return "",status
	}
	
	return claims["sub"].(string),200 

}
