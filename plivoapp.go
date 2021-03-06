package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/jithinjk/contactsapp/common"
	"github.com/jithinjk/contactsapp/contacts"
	log "github.com/sirupsen/logrus"
)

// GetHandler handler for GET calls
func GetHandler(c *gin.Context) {
	path1 := c.Param("path1")
	path2 := c.Param("path2")

	if path1 == "all" && path2 == "" {
		contacts.GetAllContacts(c)
	} else if path1 != "" && path2 == "details" {
		contactID := path1
		contacts.GetContact(c, contactID)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No contact found. Incorrect Format."})
		c.Abort()
	}
}

// Migrate migrate schema
func Migrate(db *gorm.DB) {
	db.Debug().AutoMigrate(&contacts.Contact{})
}

func main() {
	// open a db connection
	db := common.Init()
	if db != nil {
		log.Println("DB init error")
	}
	defer db.Close()

	log.Println("Connection Established...")

	db.SingularTable(true)

	//Drops table if already exists
	// db.Debug().DropTableIfExists(&Contact{})

	//Auto create table based on Model
	Migrate(db)

	router := setupRouter()
	router.Run(":8080")
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	v1 := router.Group("/api/v1/")

	v1.Use(contacts.GetRequestID())

	v1.Use(gin.BasicAuth(gin.Accounts{
		"user1": "hello",
		"user2": "world",
		"user3": "gopher",
	}))
	{
		v1.GET("/contacts/:path1", GetHandler)        //      /v1/contacts/all
		v1.GET("/contacts/:path1/:path2", GetHandler) //      /v1/contacts/<id>/details
		v1.GET("/search/name/:name", contacts.GetContactByName) //      /v1/contacts/<id>/details
		v1.GET("/search/email/:email", contacts.GetContactByEmail) //      /v1/contacts/<id>/details
		v1.POST("/create", contacts.CreateContact)
		v1.PUT("/update/:id", contacts.UpdateContact)
		v1.DELETE("/delete/:id", contacts.DeleteContact)
	}

	router.Use(cors.Default())

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"status": http.StatusNotFound, "message": "Page not found"})
	})

	return router
}
