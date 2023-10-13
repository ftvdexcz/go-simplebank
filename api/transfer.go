package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/ftvdexcz/simplebank/db/sqlc"
	"github.com/ftvdexcz/simplebank/token"
	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,min=1"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}


	
	fromAccount, valid := server.validAccount(ctx, req.FromAccountID, req.Currency)
	if !valid{
		return 
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if fromAccount.Owner != authPayload.Username{
		err := errors.New("from account does not belong to to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return 
	}
	

	if _, valid = server.validAccount(ctx, req.ToAccountID, req.Currency); !valid{
		return 
	}

	arg := db.TransferTxParams{
		FromAccountId: req.FromAccountID,
		ToAccountId: req.ToAccountID,
		Amount: req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool){
	account, err := server.store.GetAccount(ctx, accountID)

	if err != nil{
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency{
		err := fmt.Errorf("account [%d] mismatch with currency [%s] vs [%s]", accountID, account.Currency, currency)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	return account, true
}