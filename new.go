package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type album struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var albums = []album{
	{Id: "1", Username: "jayanth", Password: "jay"},
	{Id: "2", Username: "srini", Password: "sr"},
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}
func postAlbums(c *gin.Context) {
	var newAlbum album
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}
func getAlbumId(c *gin.Context) {
	id := c.Param("id")
	for _, a := range albums {
		if a.Id == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "username not found"})
}
