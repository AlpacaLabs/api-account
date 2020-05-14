package db

import (
	"fmt"
	"strings"

	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	paginationV1 "github.com/AlpacaLabs/protorepo-pagination-go/alpacalabs/pagination/v1"
)

var (
	ErrNotFound = status.Error(codes.NotFound, "entity not found")
)

type Transaction interface {
	AccountTransaction
	EmailTransaction
	PhoneTransaction
}

type txImpl struct {
	accountTxImpl
	emailTxImpl
	phoneTxImpl
}

func newTransaction(tx pgx.Tx) Transaction {
	return &txImpl{
		accountTxImpl: accountTxImpl{
			tx: tx,
		},
		emailTxImpl: emailTxImpl{
			tx: tx,
		},
		phoneTxImpl: phoneTxImpl{
			tx: tx,
		},
	}
}

func buildOrderByClause(request paginationV1.CursorRequest) string {
	var arr []string
	for _, sortClause := range request.SortClauses {
		sortString := sortKeyword(sortClause.Sort)
		arr = append(arr, fmt.Sprintf("%s %s", sortClause.FieldName, sortString))
	}
	return strings.Join(arr, ", ")
}

func sortKeyword(sort paginationV1.Sort) string {
	if sort == paginationV1.Sort_SORT_DESC {
		return "DESC"
	}
	return "ASC"
}
