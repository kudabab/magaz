package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kudabab/market-s/db"
	"github.com/kudabab/market-s/entity"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("my_secret_key")

func HashPassword(password string) (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	fmt.Println(string(hashedPassword))
	return string(hashedPassword), nil
}

func GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &entity.Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	fmt.Println("claims: ", claims)

	/*claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		Username:   username,
	}*/
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func CreateUser(ctx context.Context, user entity.User) (string, error) {

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return "", err
	}

	tokenString, err := GenerateToken(user.Username)
	if err != nil {
		return "", err
	}

	_, err = db.DB.Exec(ctx, "INSERT INTO users (username, password, email, token) VALUES ($1, $2, $3, $4)",
		user.Username, hashedPassword, user.Email, tokenString)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func FindAllUsers(ctx context.Context) ([]entity.User, error) {
	var usersSlice []entity.User
	rows, err := db.DB.Query(ctx, "SELECT id, username, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user entity.User
		if err := rows.Scan(&user.Id, &user.Username, &user.Email); err != nil {
			return nil, err
		}
		usersSlice = append(usersSlice, user)
	}

	return usersSlice, nil
}

func GetUserByUsername(ctx context.Context, username string) (entity.User, error) {
	var user entity.User
	err := db.DB.QueryRow(ctx, "SELECT id, username, email FROM users WHERE username = $1", username).
		Scan(&user.Id, &user.Username, &user.Email)
	if err != nil {
		return user, err
	}
	return user, nil

}
