package api

import (
	"database/sql"
	"net/http"

	databases "github.com/auliamarsya/simplebank/databases/sqlc"
	"github.com/gin-gonic/gin"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD IDR EUR"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var request createAccountRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponses(err))
		return
	}

	arg := databases.CreateAccountParams{
		Owner:    request.Owner,
		Balance:  0,
		Currency: request.Currency,
	}

	account, err := server.store.Queries.CreateAccount(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponses(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var request getAccountRequest

	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponses(err))
		return
	}

	account, err := server.store.GetAccount(ctx, request.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponses(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponses(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type listAccountRequest struct {
	Page  int64 `form:"page" binding:"required,min=1"`
	Limit int64 `form:"limit" binding:"required,min=5,max=30"`
}

func (server *Server) listAccount(ctx *gin.Context) {
	var request listAccountRequest

	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponses(err))
		return
	}

	arg := databases.ListAccountsParams{
		Limit:  request.Limit,
		Offset: (request.Page - 1) * request.Limit,
	}

	accounts, err := server.store.ListAccounts(ctx, arg)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponses(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponses(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
