package db

import (
	"context"
	"fmt"
	"time"

	paginationV1 "github.com/AlpacaLabs/protorepo-pagination-go/alpacalabs/pagination/v1"

	"github.com/AlpacaLabs/api-account/internal/db/entities"
	"github.com/jackc/pgx/v4"
)

type AccountTransaction interface {
	GetAccountByID(ctx context.Context, accountID string) (*entities.Account, error)
	GetAccountByUsername(ctx context.Context, username string) (*entities.Account, error)
	GetAccountByEmailAddress(ctx context.Context, emailAddress string) (*entities.Account, error)
	GetAccountByPhoneNumber(ctx context.Context, phoneNumber string) (*entities.Account, error)
	UpdateAccount(ctx context.Context, username, primaryEmailAddressID, accountID string) error
	UpdateCurrentPassword(ctx context.Context, currentPasswordID, accountID string) error
	CreateAccount(ctx context.Context, accountID, username string) error
	GetAccounts(ctx context.Context, cursorRequest paginationV1.CursorRequest) ([]*entities.Account, error)
}

type accountTxImpl struct {
	tx pgx.Tx
}

func (tx *accountTxImpl) GetAccountByID(ctx context.Context, accountID string) (*entities.Account, error) {
	var e entities.Account

	query := `
SELECT 
    id, created_at, last_modified_at, deleted_at, 
    username, current_password_id, primary_email_address_id 
 FROM account
 WHERE id=$1 
 AND deleted_at IS NULL
`

	row := tx.tx.QueryRow(ctx, query, accountID)
	err := row.Scan(&e.ID, &e.CreatedAt, &e.LastModifiedAt, &e.DeletedAt, &e.Username, &e.CurrentPasswordID, &e.PrimaryEmailAddressID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &e, nil
}

func (tx *accountTxImpl) GetAccountByUsername(ctx context.Context, username string) (*entities.Account, error) {
	var e entities.Account

	query := `
SELECT 
    id, created_at, last_modified_at, deleted_at, 
    username, current_password_id, primary_email_address_id
  FROM account 
  WHERE username=$1
  AND deleted_at IS NULL
`
	row := tx.tx.QueryRow(ctx, query, username)
	err := row.Scan(&e.ID, &e.CreatedAt, &e.LastModifiedAt, &e.DeletedAt, &e.Username, &e.CurrentPasswordID, &e.PrimaryEmailAddressID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &e, nil
}

func (tx *accountTxImpl) GetAccountByEmailAddress(ctx context.Context, emailAddress string) (*entities.Account, error) {
	var e entities.Account

	query := `
SELECT
    a.id, a.created_at, a.last_modified_at, a.deleted_at,
    a.username, a.current_password_id, a.primary_email_address_id
  FROM email_address e 
  JOIN account a ON e.account_id = a.id
  WHERE e.email_address=$1 
  AND a.deleted_at IS NULL
`
	row := tx.tx.QueryRow(ctx, query, emailAddress)
	err := row.Scan(&e.ID, &e.CreatedAt, &e.LastModifiedAt, &e.DeletedAt, &e.Username, &e.CurrentPasswordID, &e.PrimaryEmailAddressID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &e, nil
}

func (tx *accountTxImpl) GetAccountByPhoneNumber(ctx context.Context, phoneNumber string) (*entities.Account, error) {
	var e entities.Account

	query := `
SELECT
    a.id, a.created_at, a.last_modified_at, a.deleted_at,
    a.username, a.current_password_id, a.primary_email_address_id
  FROM phone_number p 
  JOIN account a ON p.account_id = a.id
  WHERE p.phone_number=$1 
  AND a.deleted_at IS NULL
`
	row := tx.tx.QueryRow(ctx, query, phoneNumber)
	err := row.Scan(&e.ID, &e.CreatedAt, &e.LastModifiedAt, &e.DeletedAt, &e.Username, &e.CurrentPasswordID, &e.PrimaryEmailAddressID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &e, nil
}

func (tx *accountTxImpl) UpdateAccount(ctx context.Context, username, primaryEmailAddressID, accountID string) error {
	query := `
UPDATE account 
  SET last_modified_at=$1, username=$2, primary_email_address_id=$3
  WHERE id=$4
`

	_, err := tx.tx.Exec(ctx, query,
		time.Now(), username, primaryEmailAddressID, accountID)
	return err
}

func (tx *accountTxImpl) UpdateCurrentPassword(ctx context.Context, currentPasswordID, accountID string) error {
	query := `
UPDATE account 
  SET last_modified_at=$1, current_password_id=$2
  WHERE id=$3
`
	_, err := tx.tx.Exec(ctx, query, time.Now(), currentPasswordID, accountID)
	return err
}

func (tx *accountTxImpl) CreateAccount(ctx context.Context, accountID, username string) error {
	query := `
INSERT INTO account(id, created_at, username)
  VALUES($1, $2, $3)
`
	_, err := tx.tx.Exec(ctx, query,
		accountID, time.Now(), username)
	return err
}

func (tx *accountTxImpl) GetAccounts(ctx context.Context, cursorRequest paginationV1.CursorRequest) ([]*entities.Account, error) {
	var sortString string
	if len(cursorRequest.SortClauses) == 0 {
		sortString = "id ASC"
	} else {
		sortString = buildOrderByClause(cursorRequest)
	}

	queryTemplate := `
SELECT 
    a.id, a.created_at, a.last_modified_at, a.deleted_at, 
    a.username, a.current_password_id, a.primary_email_address_id 
  FROM account a
  WHERE a.id > $1
  ORDER BY a.id %s
  FETCH FIRST %d ROWS ONLY
`

	query := fmt.Sprintf(queryTemplate, sortString, cursorRequest.Count)
	rows, err := tx.tx.Query(ctx, query, cursorRequest.Cursor)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	accounts := []*entities.Account{}

	for rows.Next() {
		var a entities.Account
		if err := rows.Scan(&a.ID, &a.CreatedAt, &a.LastModifiedAt, &a.DeletedAt,
			&a.Username, &a.CurrentPasswordID, &a.PrimaryEmailAddressID); err != nil {
			return nil, err
		}
		accounts = append(accounts, &a)
	}

	return accounts, nil
}
