package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

//APIDraft is used in several functions that handle drafts
type APIDraft struct {
	Document string `json:"name" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

// getAllDraftsHandler returns the most recent drafts for every document
// router.GET("/api/drafts", getAllDraftsHandler(db))
func getAllDraftsHandler(db *sql.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		// query for drafts
		rows, err := db.Query("SELECT name, content FROM draft WHERE draft_id IN ( SELECT max(draft_id) FROM draft GROUP BY name);")
		if err != nil {
			c.JSON(http.StatusInternalServerError /*500*/, gin.H{"error": err.Error()})
			return
		}

		// format data
		defer rows.Close()
		var drafts []APIDraft
		for rows.Next() { // iterate through rows cursor
			var params APIDraft
			err := rows.Scan(&params.Document, &params.Content)
			if err != nil {
				panic(err) // unhandled null type found
			}
			drafts = append(drafts, params)
		}

		// return data
		c.JSON(http.StatusOK, drafts)
	}
}

// postDraftHandler inserts a new Draft containing "content" with "name"
// router.POST("/api/drafts", postDraftHandler(db))
// accepts json "name" and "content"
func postDraftHandler(db *sql.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		// retrieve parameters
		var params APIDraft
		if err := c.BindJSON(&params); err != nil {
			// c.JSON(http.StatusNotAcceptable /*406*/, gin.H{"error": err.Error()})
			return
		}

		// insert new draft
		sql := `INSERT INTO draft(name, content) VALUES (?, ?)`
		stmt, err := db.Prepare(sql)
		if err != nil {
			c.JSON(http.StatusInternalServerError /*500*/, gin.H{"error": err.Error()})
			return
		}

		_, err = stmt.Exec(params.Document, params.Content)
		if err != nil {
			c.JSON(http.StatusInternalServerError /*500*/, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, params)
	}
}

// APIComment is used by several functions that handle comments
type APIComment struct {
	ID      string         `json:"id" binding:"required"`
	Parent  sql.NullString `json:"parent"` // any type that is nullable must allow nulls
	Content string         `json:"content" binding:"required"`
	User    string         `json:"user" binding:"required"`
}

// postCommentOnDraftHandler inserts new Comment containing "content" on draft with "id"
// router.POST("/api/createcomment", postCommentOnDraftHandler(db))
// accepts json "user" and "id" and "content"
func postCommentOnDraftHandler(db *sql.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		// retrieve parameters
		var params APIComment

		if err := c.BindJSON(&params); err != nil {
			// c.JSON(http.StatusNotAcceptable /*406*/, gin.H{"error": err.Error()})
			return
		}

		// check user auth
		if err := checkUserAuthenticated(c, db, params.User); err != nil {
			return
		}

		// insert new comment on draft

		sql := `INSERT INTO "comment" (draft_id, content, created_by) VALUES (?, ?, ?)`
		stmt, err := db.Prepare(sql)
		if err != nil {
			c.JSON(http.StatusInternalServerError /*500*/, gin.H{"error": err.Error()})
			return
		}

		_, err = stmt.Exec(params.ID, params.Content, params.User)
		if err != nil {
			c.JSON(http.StatusInternalServerError /*500*/, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, params)
	}
}

// postCommentOnCommentHandler inserts new Comment containing "content" on comment with "id"
// router.POST("/api/commentcomment", postCommentOnCommentHandler(db))
// accepts json "user" and "id" and "content"
func postCommentOnCommentHandler(db *sql.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		// retrieve parameters
		var params APIComment

		if err := c.BindJSON(&params); err != nil {
			// c.JSON(http.StatusNotAcceptable /*406*/, gin.H{"error": err.Error()})
			return
		}

		// check user auth
		if err := checkUserAuthenticated(c, db, params.User); err != nil {
			return
		}

		// insert new comment on comment

		sql := `INSERT INTO "comment" (draft_id, parent_id, content, created_by) SELECT draft_id, ?, ?, ? FROM "comment" WHERE comment_id = ?`
		stmt, err := db.Prepare(sql)
		if err != nil {
			c.JSON(http.StatusInternalServerError /*500*/, gin.H{"error": err.Error()})
			return
		}

		_, err = stmt.Exec(params.ID, params.Content, params.User, params.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError /*500*/, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, params)
	}
}

// APIDraftGet is used by getCommentsForDraftHandler to handle json input
type APIDraftGet struct {
	ID string `json:"id" binding:"required"`
}

// getCommentsForDraftHandler returns all comments on a draft with "id"
// router.POST("/api/comments", getCommentsForDraftHandler(db))
// accepts json "id"
func getCommentsForDraftHandler(db *sql.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		// retrieve parameters
		var params APIDraftGet

		if err := c.BindJSON(&params); err != nil {
			// c.JSON(http.StatusNotAcceptable /*406*/, gin.H{"error": err.Error()})
			return
		}

		// query for comments
		rows, err := db.Query(`SELECT comment_id, parent_id, content, created_by FROM "comment" WHERE draft_id = ?`, params.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError /*500*/, gin.H{"error": err.Error()})
			return
		}

		// format data
		defer rows.Close()
		var comments []APIComment
		for rows.Next() { // iterate through rows cursor
			var params APIComment
			err := rows.Scan(&params.ID, &params.Parent, &params.Content, &params.User)
			if err != nil {
				panic(err) // unhandled null type found
			}
			comments = append(comments, params)
		}

		// return data
		c.JSON(http.StatusOK, comments)

	}
}

// APIReaction is used by several functions to handle reactions
type APIReaction struct {
	ID      string `json:"id" binding:"required"`
	Content string `json:"content" binding:"required"`
	User    string `json:"user" binding:"required"`
}

// APIReactionGet is used by getReactionsOnCommentHandler to contain input
type APIReactionGet struct {
	ID string `json:"id" binding:"required"`
}

// getReactionsOnCommentHandler inserts new Reaction of "content" on comment with "id"
// router.POST("/api/commentreactions", getReactionsOnCommentHandler(db))
// accepts json "id" and "content" and "user"
func getReactionsOnCommentHandler(db *sql.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		// retrieve parameters
		var params APIReactionGet

		if err := c.BindJSON(&params); err != nil {
			// c.JSON(http.StatusNotAcceptable /*406*/, gin.H{"error": err.Error()})
			return
		}

		// query for reactions
		sql := `SELECT reaction_id, reaction, created_by FROM reaction WHERE comment_id = ?`
		stmt, err := db.Prepare(sql)
		if err != nil {
			c.JSON(http.StatusInternalServerError /*500*/, gin.H{"error": err.Error()})
			return
		}

		rows, err := stmt.Query(params.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError /*500*/, gin.H{"error": err.Error()})
			return
		}

		// format results
		defer rows.Close()
		var reacts []APIReaction
		for rows.Next() { // iterate through rows cursor
			var params APIReaction
			err := rows.Scan(&params.ID, &params.Content, &params.User)
			if err != nil {
				panic(err) // unhandled null type found
			}
			reacts = append(reacts, params)
		}

		// return data
		c.JSON(http.StatusOK, reacts)

	}
}

// postReactionHandler inserts new Reaction of "content" on comment with "id"
// router.POST("/api/reaction", postReactionHandler(db))
// accepts json "id" and "content" and "user"
func postReactionHandler(db *sql.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		// retrieve parameters
		var params APIReaction

		if err := c.BindJSON(&params); err != nil {
			// c.JSON(http.StatusNotAcceptable /*406*/, gin.H{"error": err.Error()})
			return
		}

		// insert new reaction on comment

		sql := `INSERT INTO "reaction" (comment_id, reaction, created_by) VALUES (?, ?, ?)`
		stmt, err := db.Prepare(sql)
		if err != nil {
			c.JSON(http.StatusInternalServerError /*500*/, gin.H{"error": err.Error()})
			return
		}

		_, err = stmt.Exec(params.ID, params.Content, params.User)
		if err != nil {
			c.JSON(http.StatusInternalServerError /*500*/, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, params)
	}
}

// APIDraftSearch is used in findTextInDraftsHandler as the input data for the request
type APIDraftSearch struct {
	Content string `json:"content" binding:"required"`
}

// findTextInDraftsHandler returns drafts with content matching input
// router.POST("/api/findindrafts", findTextInDraftsHandler(db))
// accepts json "content"
func findTextInDraftsHandler(db *sql.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		// retrieve parameters
		var params APIDraftSearch

		if err := c.BindJSON(&params); err != nil {
			// c.JSON(http.StatusNotAcceptable /*406*/, gin.H{"error": err.Error()})
			return
		}

		search_content := "%" + params.Content + "%"

		// query for drafts with contents like input
		rows, err := db.Query(`SELECT name, content FROM draft WHERE content LIKE ?`, search_content)
		if err != nil {
			c.JSON(http.StatusInternalServerError /*500*/, gin.H{"error": err.Error()})
			return
		}

		// format data
		defer rows.Close()
		var drafts []APIDraft
		for rows.Next() { // iterate through rows cursor
			var params APIDraft
			err := rows.Scan(&params.Document, &params.Content)
			if err != nil {
				panic(err) // unhandled null type found
			}
			drafts = append(drafts, params)
		}

		// return data
		c.JSON(http.StatusOK, drafts)
	}
}
