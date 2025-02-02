package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-crud/initializers"
	"go-crud/models"
	"net/http"
	"os"
	"time"
)

func RequireAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	fmt.Println("TokenString:", tokenString)
	if err != nil {
		fmt.Println("Error:", err)
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			fmt.Println("Expired token")
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		var user models.User
		initializers.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			fmt.Println("USER NOT FOUND")
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		c.Set("user", user)
		c.Next()
	} else {
		fmt.Println("invalid token")
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	fmt.Println("In middleware")
	c.Next()
}
