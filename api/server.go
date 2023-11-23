package api

import (
	databases "github.com/auliamarsya/simplebank/databases/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *databases.Store
	router *gin.Engine
}

func NewServer(store *databases.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponses(err error) gin.H {
	return gin.H{"error": err.Error()}
}
