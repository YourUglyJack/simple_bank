package api

import (
	"net/http"
	db "simple_bank/db/sqlc"

	"github.com/gin-gonic/gin"
)


type createAccountReq struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=CNY JPY"`
}

// NOTE：func (c *Context) JSON(code int, obj any)， obj是被返回的对象，会被序列化为JSON
func (s *Server) createAccount(ctx *gin.Context) {
	var req createAccountReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	arg := db.CreateAccountParams{
		Balance: 0,
		Owner: req.Owner,
		Currency: req.Currency,	
	}

	account, err := s.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	
	ctx.JSON(http.StatusOK, account)

}