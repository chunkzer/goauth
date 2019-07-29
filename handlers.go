package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/rolandovlz/goauth/models"
	"golang.org/x/crypto/bcrypt"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	user := &models.User{}

	err = json.Unmarshal(body, user)
	if err != nil {
		log.Fatal(err)
	}

	row := db.QueryRow("SELECT password FROM users WHERE username=$1;", user.Username)

	var hashFromDb string

	switch err = row.Scan(&hashFromDb); err {
	case sql.ErrNoRows:
		log.Println("Not found!")
		http.Redirect(w, r, "/login", http.StatusBadRequest)
	case nil:
		if err := bcrypt.CompareHashAndPassword([]byte(hashFromDb), []byte(user.Password)); err != nil {
			http.Redirect(w, r, "/login", http.StatusUnauthorized)
		} else {
			expirationTime := time.Now().Add(5 * time.Minute)
			claims := &models.Claims{
				Username: user.Username,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: expirationTime.Unix(),
				},
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString(jwtKey)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusInternalServerError)
			}
			http.SetCookie(w, &http.Cookie{
				Name:    "token",
				Value:   tokenString,
				Expires: expirationTime,
			})

			http.Redirect(w, r, "/welcome", http.StatusOK)
		}
	default:
		log.Fatal(err)
		http.Redirect(w, r, "/login", http.StatusInternalServerError)
	}
}

func signUpHandler(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	user := &models.User{}

	err = json.Unmarshal(body, user)
	if err != nil {
		log.Fatal(err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	_, err = db.Exec("INSERT INTO users (username, password) VALUES ($1, $2);", user.Username, string(hash))
	if err != nil {
		log.Fatal(err)
		http.Redirect(w, r, "/signup", http.StatusBadRequest)
	}

	log.Println("Added new user to db:", user.Username)
	http.Redirect(w, r, "/login", http.StatusCreated)
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenString := c.Value
	claims := &models.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Write([]byte(fmt.Sprintf("Welcome!")))
}
