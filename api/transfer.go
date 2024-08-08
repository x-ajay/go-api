package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	db "github.com/x-ajay/go-api/db/sqlc"
)

type trasferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int32  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,oneof=USD EUR JPY"`
}

func (server *Server) createTrasfer(c *gin.Context) {
	var req trasferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !server.isValidAccount(c, req.FromAccountID, db.Currency(req.Currency)) {
		return
	}

	if !server.isValidAccount(c, req.ToAccountID, db.Currency(req.Currency)) {
		return
	}

	// transfer the amount from one account to another
	result, err := server.store.TransferTx(c, db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, result)
}

func (server *Server) isValidAccount(c *gin.Context, accountID int64, currency db.Currency) bool {
	account, err := server.store.GetAccount(c, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if account.Currency != currency {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account currency mismatch"})
		return false
	}

	return true
}
