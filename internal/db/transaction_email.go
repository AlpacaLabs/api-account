package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/AlpacaLabs/api-account/internal/db/entities"
	accountV1 "github.com/AlpacaLabs/protorepo-account-go/alpacalabs/account/v1"
	paginationV1 "github.com/AlpacaLabs/protorepo-pagination-go/alpacalabs/pagination/v1"
	"github.com/jackc/pgx/v4"
)

type EmailTransaction interface {
	CreateEmailAddress(ctx context.Context, e entities.EmailAddress) error
	DeleteEmailAddress(ctx context.Context, id string) error

	GetEmailAddressByEmailAddress(ctx context.Context, emailAddress string) (*accountV1.EmailAddress, error)
	GetEmailAddressByID(ctx context.Context, id string) (*accountV1.EmailAddress, error)

	GetEmailAddresses(ctx context.Context, request paginationV1.CursorRequest) ([]*accountV1.EmailAddress, error)

	GetEmailAddressesForAccount(ctx context.Context, accountID string, cursorRequest paginationV1.CursorRequest) ([]*accountV1.EmailAddress, error)

	EmailIsConfirmed(ctx context.Context, emailAddress string) (bool, error)
	EmailExists(ctx context.Context, emailAddress string) (bool, error)
	CountEmail(ctx context.Context, emailAddress string) (int, error)
	GetConfirmedEmailAddress(ctx context.Context) (*accountV1.EmailAddress, error)
}

type emailTxImpl struct {
	tx pgx.Tx
}

func (tx *emailTxImpl) CreateEmailAddress(ctx context.Context, e entities.EmailAddress) error {
	query := `
INSERT INTO email_address
 (id, account_id, email_address, confirmed, is_primary)
 VALUES($1, $2, $3, $4, $5)
`
	_, err := tx.tx.Exec(ctx, query, e.ID, e.AccountID, e.EmailAddress, e.Confirmed, e.Primary)

	return err
}

func (tx *emailTxImpl) DeleteEmailAddress(ctx context.Context, id string) error {
	_, err := tx.tx.Exec(ctx, "DELETE FROM email_address WHERE id=$1", id)

	return err
}

func (tx *emailTxImpl) GetEmailAddressByEmailAddress(ctx context.Context, emailAddress string) (*accountV1.EmailAddress, error) {
	var e entities.EmailAddress

	query := `
SELECT id, created_at, last_modified_at, deleted_at, confirmed, is_primary, email_address, account_id 
 FROM email_address
 WHERE email_address=$1 
 AND deleted_at IS NULL
`

	row := tx.tx.QueryRow(ctx, query, emailAddress)
	err := row.Scan(&e.ID, &e.CreatedAt, &e.LastModifiedAt, &e.DeletedAt, &e.Confirmed, &e.Primary, &e.EmailAddress, &e.AccountID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return e.ToProtobuf(), nil
}

func (tx *emailTxImpl) GetEmailAddressByID(ctx context.Context, id string) (*accountV1.EmailAddress, error) {
	var e entities.EmailAddress

	row := tx.tx.QueryRow(
		ctx,
		"SELECT id, created_at, last_modified_at, deleted_at, confirmed, is_primary, email_address, account_id "+
			"FROM phone_number WHERE id=$1 "+
			"AND deleted_at IS NULL", id)
	err := row.Scan(&e.ID, &e.CreatedAt, &e.LastModifiedAt, &e.DeletedAt, &e.Confirmed, &e.Primary, &e.EmailAddress, &e.AccountID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return e.ToProtobuf(), nil
}

func (tx *emailTxImpl) GetEmailAddresses(ctx context.Context, request paginationV1.CursorRequest) ([]*accountV1.EmailAddress, error) {
	var sortString string
	if len(request.SortClauses) == 0 {
		sortString = "id ASC"
	} else {
		sortString = buildOrderByClause(request)
	}

	queryTemplate := `
SELECT id, created_at, last_modified_at, deleted_at, confirmed, email_address, account_id 
 FROM email_address
 WHERE id > $1
 AND deleted_at IS NULL
 ORDER BY %s
 FETCH FIRST %d ROWS ONLY
`

	query := fmt.Sprintf(queryTemplate, sortString, request.Count)

	rows, err := tx.tx.Query(ctx, query, request.Cursor)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	emailAddresses := []*accountV1.EmailAddress{}

	for rows.Next() {
		var e entities.EmailAddress
		if err := rows.Scan(&e.ID, &e.CreatedAt, &e.LastModifiedAt, &e.DeletedAt, &e.Confirmed, &e.EmailAddress, &e.AccountID); err != nil {
			return nil, err
		}
		emailAddresses = append(emailAddresses, e.ToProtobuf())
	}

	return emailAddresses, nil
}

func (tx *emailTxImpl) GetEmailAddressesForAccount(ctx context.Context, accountID string, cursorRequest paginationV1.CursorRequest) ([]*accountV1.EmailAddress, error) {
	var sortString string
	if len(cursorRequest.SortClauses) == 0 {
		sortString = "id ASC"
	} else {
		sortString = buildOrderByClause(cursorRequest)
	}

	queryTemplate := `
SELECT id, email_address, account_id 
 FROM email_address 
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

	emailAddresses := []*accountV1.EmailAddress{}

	for rows.Next() {
		var e accountV1.EmailAddress
		if err := rows.Scan(&e.Id, &e.EmailAddress, &e.AccountId); err != nil {
			return nil, err
		}
		// TODO do masking in service layer, not db layer
		e.EmailAddress = maskEmail(e.EmailAddress)
		emailAddresses = append(emailAddresses, &e)
	}

	return emailAddresses, nil
}

func (tx *emailTxImpl) EmailIsConfirmed(ctx context.Context, emailAddress string) (bool, error) {
	var count int

	query := `
SELECT COUNT(*) AS count 
 FROM email_address 
 WHERE email_address = $1
 AND confirmed = $2
 AND deleted_at IS NULL
`

	row := tx.tx.QueryRow(ctx, query, emailAddress, true)
	err := row.Scan(&count)

	if err != nil {
		// TODO check for NotFound?
		return false, err
	}
	return count == 1, nil
}

func (tx *emailTxImpl) EmailExists(ctx context.Context, emailAddress string) (bool, error) {
	count, err := tx.CountEmail(ctx, emailAddress)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return count == 1, nil
}

func (tx *emailTxImpl) CountEmail(ctx context.Context, emailAddress string) (int, error) {
	var count int

	query := `
SELECT COUNT(*) AS count 
 FROM email_address 
 WHERE email_address=$1 
 AND deleted_at IS NULL
`

	row := tx.tx.QueryRow(ctx, query, emailAddress)
	err := row.Scan(&count)
	return count, err
}

func (tx *emailTxImpl) GetConfirmedEmailAddress(ctx context.Context) (*accountV1.EmailAddress, error) {
	var e entities.EmailAddress

	query := `
SELECT id, email_address, account_id 
 FROM email_address WHERE email_address=$1 
 AND confirmed=$2 
 AND deleted_at IS NULL
`

	row := tx.tx.QueryRow(ctx, query, e.EmailAddress, true)

	err := row.Scan(&e.ID, &e.EmailAddress, &e.AccountID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return e.ToProtobuf(), nil
}

func maskEmail(emailAddress string) string {
	return getMaskedEmailUser(emailAddress) + "@" + getMaskedEmailHost(emailAddress)
}

func getMaskedEmailUser(emailAddress string) string {
	splits := strings.Split(emailAddress, "@")
	user := splits[0]
	if len(user) == 1 {
		return user[0:1] + strings.Repeat("*", len(user)-1)
	}
	return user[0:2] + strings.Repeat("*", len(user)-2)
}

func getMaskedEmailHost(emailAddress string) string {
	emailSplits := strings.Split(emailAddress, "@")
	host := emailSplits[1]
	splits := strings.Split(host, ".")
	splits[0] = splits[0][0:1] + strings.Repeat("*", len(splits[0])-1)
	return strings.Join(splits, ".")
}
