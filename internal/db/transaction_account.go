package db

import (
	"context"

	"github.com/AlpacaLabs/api-account/internal/db/entities"
	"github.com/jackc/pgx/v4"
)

type AccountTransaction interface {
	GetAccountByID(ctx context.Context, accountID string) (*entities.Account, error)
}

type accountTxImpl struct {
	tx pgx.Tx
}

func (tx *accountTxImpl) GetAccountByID(ctx context.Context, accountID string) (*entities.Account, error) {
	var e entities.Account

	query := `
SELECT id, created_at, deleted_at, last_modified_at, username, current_password_id, primary_email_address_id 
 FROM account
 WHERE id=$1 
 AND deleted_at IS NULL
`

	row := tx.tx.QueryRow(ctx, query, accountID)
	err := row.Scan(&e.ID, &e.CreatedAt, &e.DeletedAt, &e.LastModifiedAt, &e.Username, &e.CurrentPasswordID, &e.PrimaryEmailAddressID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &e, nil
}
