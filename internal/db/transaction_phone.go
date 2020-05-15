package db

import (
	"context"
	"fmt"

	"github.com/AlpacaLabs/api-account/internal/db/entities"
	accountV1 "github.com/AlpacaLabs/protorepo-account-go/alpacalabs/account/v1"
	paginationV1 "github.com/AlpacaLabs/protorepo-pagination-go/alpacalabs/pagination/v1"
	"github.com/jackc/pgx/v4"
)

type PhoneTransaction interface {
	CreatePhoneNumber(ctx context.Context, e entities.PhoneNumber) error
	DeletePhoneNumber(ctx context.Context, id string) (int, error)
	ConfirmPhoneNumber(ctx context.Context, id string) error

	GetPhoneNumberByID(ctx context.Context, id string) (*accountV1.PhoneNumber, error)
	GetPhoneNumbersForAccount(ctx context.Context, accountID string, cursorRequest paginationV1.CursorRequest) ([]*accountV1.PhoneNumber, error)
	GetPhoneNumberByPhoneNumber(ctx context.Context, phoneNumber string) (*accountV1.PhoneNumber, error)
}

type phoneTxImpl struct {
	tx pgx.Tx
}

func (tx *phoneTxImpl) CreatePhoneNumber(ctx context.Context, e entities.PhoneNumber) error {
	query := `
INSERT INTO phone_number
 (id, account_id, phone_number, confirmed)
 VALUES($1, $2, $3, $4)
`
	_, err := tx.tx.Exec(ctx, query, e.ID, e.AccountID, e.PhoneNumber, e.Confirmed)

	return err
}

func (tx *phoneTxImpl) DeletePhoneNumber(ctx context.Context, id string) (int, error) {
	res, err := tx.tx.Exec(ctx, "DELETE FROM phone_number WHERE id=$1", id)
	if err != nil {
		return 0, err
	}

	return int(res.RowsAffected()), nil
}

func (tx *phoneTxImpl) ConfirmPhoneNumber(ctx context.Context, id string) error {
	_, err := tx.tx.Exec(ctx, "UPDATE phone_number SET confirmed = TRUE WHERE id=$1", id)
	return err
}

func (tx *phoneTxImpl) GetPhoneNumberByID(ctx context.Context, id string) (*accountV1.PhoneNumber, error) {
	var p entities.PhoneNumber

	query := `
SELECT id, created_at, last_modified_at, deleted_at, confirmed, phone_number, account_id 
 FROM phone_number WHERE id=$1 
 AND deleted_at IS NULL
`

	row := tx.tx.QueryRow(ctx, query, id)
	err := row.Scan(&p.ID, &p.CreatedAt, &p.LastModifiedAt, &p.DeletedAt, &p.Confirmed, &p.PhoneNumber, &p.AccountID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return p.ToProtobuf(), nil
}

func (tx *phoneTxImpl) GetPhoneNumbersForAccount(ctx context.Context, accountID string, cursorRequest paginationV1.CursorRequest) ([]*accountV1.PhoneNumber, error) {
	var sortString string
	if len(cursorRequest.SortClauses) == 0 {
		sortString = "id ASC"
	} else {
		sortString = buildOrderByClause(cursorRequest)
	}

	queryTemplate := `
SELECT id, phone_number, account_id, created_at, confirmed
 FROM phone_number 
 WHERE confirmed=TRUE 
 AND account_id=$1 
 AND deleted_at IS NULL 
 AND id > $2
 ORDER BY %s 
 FETCH FIRST %d ROWS ONLY
`
	query := fmt.Sprintf(queryTemplate, sortString, cursorRequest.Count)
	rows, err := tx.tx.Query(ctx, query, accountID, cursorRequest.Cursor)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	phoneNumbers := []*accountV1.PhoneNumber{}

	for rows.Next() {
		var p entities.PhoneNumber
		if err := rows.Scan(&p.ID, &p.PhoneNumber, &p.AccountID, &p.CreatedAt, &p.Confirmed); err != nil {
			return nil, err
		}
		p.PhoneNumber = maskPhoneNumber(p.PhoneNumber)
		phoneNumbers = append(phoneNumbers, p.ToProtobuf())
	}

	return phoneNumbers, nil
}

func (tx *phoneTxImpl) GetPhoneNumberByPhoneNumber(ctx context.Context, phoneNumber string) (*accountV1.PhoneNumber, error) {
	var p entities.PhoneNumber

	query := `
SELECT id, phone_number, account_id, created_at, confirmed
 FROM phone_number WHERE phone_number=$1 
 AND deleted_at IS NULL
`

	err := tx.tx.QueryRow(ctx, query, phoneNumber).Scan(&p.ID, &p.PhoneNumber, &p.AccountID, &p.CreatedAt, &p.Confirmed)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return p.ToProtobuf(), nil
}

func maskPhoneNumber(phoneNumber string) string {
	return phoneNumber[len(phoneNumber)-2:]
}
