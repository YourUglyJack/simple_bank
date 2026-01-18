package api

import (
	db "simple_bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store *db.Store
	router *gin.Engine
}


func NewServer(store *db.Store) *Server {
	server := &Server{
		store: store,
	}

	router := gin.Default()

	accounts := router.Group("/accounts")
	{
		accounts.POST("", server.createAccount)
		accounts.GET("/:id", server.getAccount)
		accounts.GET("", server.listAccount)
	}

	server.router = router
	return server
}

func (server *Server) Start(addr string) error {
	return server.router.Run(addr)
}


func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}