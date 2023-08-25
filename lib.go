package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// checkUserAuthenticated returns an error if user is not an "admin"
func checkUserAuthenticated(c *gin.Context, db *sql.DB, user_id string) error {

	// get token for user from db
	stmt, err := db.Prepare("SELECT token, role FROM user WHERE user_id = ?")
	if err != nil {
		c.JSON(http.StatusInternalServerError /*500*/, gin.H{"error": err.Error()})
		return fmt.Errorf("Error in user authentication")
	}
	defer stmt.Close()

	var token string
	var role string
	err = stmt.QueryRow(user_id).Scan(&token, &role)
	if err != nil {
		c.JSON(http.StatusInternalServerError /*500*/, gin.H{"error": err.Error()})
		return fmt.Errorf("Error in user authentication")
	}

	// check bearer token in header against token in db
	if c.GetHeader("Authorization") != "Bearer "+token {
		c.JSON(http.StatusUnauthorized /*401*/, gin.H{"error": "401 Unauthorized"})
		return fmt.Errorf("User token does not match")
	}

	// check user has perms
	if role != "admin" {
		c.JSON(http.StatusUnauthorized /*401*/, gin.H{"error": "401 Unauthorized"})
		return fmt.Errorf("User does not have permissions for action")
	}

	return nil
}
