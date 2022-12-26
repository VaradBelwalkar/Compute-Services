package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

// The secret key used to sign and verify the JWT
var secretKey []byte

// The public key used to verify the JWT
var publicKey *rsa.PublicKey

// Initialize the secret key and public key
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
func signJWT(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// Verify a JWT using the HS256 algorithm
func verifyJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("Invalid JWT")
}

// A handler function that signs a JWT and sends it to the client
func signHandler(w http.ResponseWriter, r *http.Request) {
	// Set the claims for the JWT
	claims := jwt.MapClaims{
		"sub": "user1",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	// Sign the JWT
	token, err := signJWT(claims)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the JWT to the client
	w.Write([]byte(token))
}


// A handler function that verifies a JWT sent by the client
func verifyHandler(w http.ResponseWriter, r *http.Request) {
	// Get the JWT from the request
	tokenString := r.FormValue("token")

	// Verify the JWT
	claims, err := verifyJWT(tokenString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Write a response
	w.Write([]byte("JWT verified"))
}
