package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var HostPort = flag.String("hostport", ":8080", "Host/Port to listen on")
var DBPath = flag.String("db", "./bin/sqlite.db", "Database to load on startup")
var DBData = flag.String("dbdata", "sql/data.sql", "Data to load on startup")

func check(e error) {
	if e != nil {
		fmt.Println("Encountered fatal error - quitting")
		log.Fatal(e.Error())
	}
}

func main() {

	// ------------------------------------------------------------
	// cmd line args
	// ------------------------------------------------------------
	flag.Parse()

	fns := flag.Args() // check for extra arguments
	if len(fns) != 0 {
		log.Fatal("The only argument suppored is hostport\n")
	}

	// ------------------------------------------------------------
	// setup database, router, etc.
	// ------------------------------------------------------------

	router := gin.Default()

	// check db file exists; if not, create
	if _, err := os.Stat(*DBPath); err != nil {
		// does not exist

		file, err := os.Create(*DBPath)
		check(err)
		file.Close()
		log.Println("Created sqlite db file at", *DBPath)

		db, err := sql.Open("sqlite3", *DBPath+"?_foreign_keys=on")
		check(err)

		// create tables

		stmt, err := ioutil.ReadFile("sql/tables.sql")
		check(err)
		if _, err := db.Exec(string(stmt)); err != nil {
			log.Fatal(err.Error())
		}

		// insert test data

		stmt, err = ioutil.ReadFile(*DBData)
		check(err)
		if _, err := db.Exec(string(stmt)); err != nil {
			log.Fatal(err.Error())
		}

		db.Close()

	}
	// open db
	db, err := sql.Open("sqlite3", *DBPath+"?_foreign_keys=on")
	check(err)
	defer db.Close()
	log.Println("Opened sqlite db at", *DBPath)

	// ------------------------------------------------------------
	// api endpoints
	// ------------------------------------------------------------

	// status GET
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK /*200*/, gin.H{
			"status": "success",
			"msg":    "Hello there!",
		})
	})

	// ------------------------------------------------------------
	// specifically requested API functionality

	router.GET("/api/status", func(c *gin.Context) {
		c.JSON(http.StatusOK /*200*/, gin.H{
			"status": "success",
			"msg":    "Hello there!",
		})
	})

	// drafts GET returns all drafts
	router.GET("/api/drafts", getAllDraftsHandler(db))

	// drafts POST adds a new draft
	router.POST("/api/drafts", postDraftHandler(db))

	// ------------------------------------------------------------

	// comments POST returns all comments for a draft
	router.POST("/api/comments", getCommentsForDraftHandler(db))

	// createcomment POST create comment on draft
	router.POST("/api/createcomment", postCommentOnDraftHandler(db))
	// commentcomment POST create comment on comment
	router.POST("/api/commentcomment", postCommentOnCommentHandler(db))

	// reaction GET get reactions on comment
	router.POST("/api/commentreaction", getReactionsOnCommentHandler(db))
	// reaction POST react to comment
	router.POST("/api/reaction", postReactionHandler(db))

	// findindrafts POST search for text in all drafts
	router.POST("/api/findindrafts", findTextInDraftsHandler(db))

	// createuser POST adds a new user
	// router.POST("/api/createuser", createUserHandler(db))

	// ------------------------------------------------------------
	router.Run(*HostPort) // default localhost:8080
}
