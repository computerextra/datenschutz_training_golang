package api

import (
	"net/http"

	"github.com/computerextra/datenschutz_training_golang/db"
	"github.com/gin-gonic/gin"
)

// Define access checkers
func BookReadOnly(login db.Login) bool {
	return login.Permission.AccessBook == db.ReadOnly || login.Permission.AccessBook == db.Write
}

func BookWrite(login db.Login) bool {
	return login.Permission.AccessBook == db.Write
}

func LoginReadOnly(login db.Login) bool {
	return login.Permission.AccessAdmin == db.ReadOnly || login.Permission.AccessAdmin == db.Write
}

func LoginWrite(login db.Login) bool {
	return login.Permission.AccessAdmin == db.Write
}

// Define auth middleware
func AuthMiddleware(checker func(db.Login) bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// send auth header
		c.Header("WWW-Authenticate", "Basic realm=\"Restricted\"")

		username, password, ok := c.Request.BasicAuth()
		// fmt.Println(username, password)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Lookup user in your database
		user, found := db.FetchLogin(username)
		if found != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Compare the password
		// err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		err := user.Password != password
		if err {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Check user's permission
		if !checker(user) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			return
		}

		// Set the user in the context
		c.Set("user", user)

		c.Next()
	}
}
