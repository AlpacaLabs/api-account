package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	accountV1 "github.com/AlpacaLabs/protorepo-account-go/alpacalabs/account/v1"

	"github.com/jackc/pgx/v4"

	"github.com/AlpacaLabs/api-account/internal/db/entities"
	paginationV1 "github.com/AlpacaLabs/protorepo-pagination-go/alpacalabs/pagination/v1"
)

var (
	ErrNotFound = status.Error(codes.NotFound, "entity not found")
)

type Transaction interface {
	GetAccountByID(ctx context.Context, accountID string) (*entities.Account, error)

	GetEmailAddressByEmailAddress(ctx context.Context, emailAddress string) (*accountV1.EmailAddress, error)

	GetEmailAddressByID(ctx context.Context, id string) (*accountV1.EmailAddress, error)
	GetPhoneNumberByID(ctx context.Context, id string) (*accountV1.PhoneNumber, error)

	GetDeletedEmailAddressByID(ctx context.Context) (*accountV1.EmailAddress, error)
	UpdateEmailAddress(ctx context.Context) error
	DeleteEmailAddress(ctx context.Context, id string) error
	CreateEmailAddress(ctx context.Context, e entities.EmailAddress) error
	GetEmailAddresses(ctx context.Context, request paginationV1.CursorRequest) ([]*accountV1.EmailAddress, error)

	GetEmailAddressesForAccount(ctx context.Context, accountID string, cursorRequest paginationV1.CursorRequest) ([]*accountV1.EmailAddress, error)
	GetPhoneNumbersForAccount(ctx context.Context, accountID string, cursorRequest paginationV1.CursorRequest) ([]*accountV1.PhoneNumber, error)

	EmailIsConfirmed(ctx context.Context, emailAddress string) (bool, error)
	EmailExists(ctx context.Context, emailAddress string) (bool, error)
	CountEmail(ctx context.Context, emailAddress string) (int, error)
	GetConfirmedEmailAddress(ctx context.Context) (*accountV1.EmailAddress, error)
	GetPhoneNumberByPhoneNumber(ctx context.Context, phoneNumber string) (*accountV1.PhoneNumber, error)
}

type txImpl struct {
	tx pgx.Tx
}

func (tx *txImpl) GetAccountByID(ctx context.Context, accountID string) (*entities.Account, error) {
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

func (tx *txImpl) GetEmailAddressByEmailAddress(ctx context.Context, emailAddress string) (*accountV1.EmailAddress, error) {
	var e entities.EmailAddress

	query := `
SELECT id, created_at, deleted_at, last_modified_at, confirmed, is_primary, email_address, account_id 
 FROM email_address
 WHERE email_address=$1 
 AND deleted_at IS NULL
`

	row := tx.tx.QueryRow(ctx, query, emailAddress)
	err := row.Scan(&e.ID, &e.CreatedAt, &e.DeletedAt, &e.LastModifiedAt, &e.Confirmed, &e.Primary, &e.EmailAddress, &e.AccountID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return e.ToProtobuf(), nil
}

func (tx *txImpl) GetEmailAddressByID(ctx context.Context, id string) (*accountV1.EmailAddress, error) {
	var e entities.EmailAddress

	row := tx.tx.QueryRow(
		ctx,
		"SELECT id, created_at, deleted_at, last_modified_at, confirmed, is_primary, email_address, account_id "+
			"FROM phone_number WHERE id=$1 "+
			"AND deleted_at IS NULL", id)
	err := row.Scan(&e.ID, &e.CreatedAt, &e.DeletedAt, &e.LastModifiedAt, &e.Confirmed, &e.Primary, &e.EmailAddress, &e.AccountID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return e.ToProtobuf(), nil
}

func (tx *txImpl) GetPhoneNumberByID(ctx context.Context, id string) (*accountV1.PhoneNumber, error) {
	var p entities.PhoneNumber

	row := tx.tx.QueryRow(
		ctx,
		"SELECT id, created_at, deleted_at, last_modified_at, confirmed, phone_number, account_id "+
			"FROM email_address WHERE id=$1 "+
			"AND deleted_at IS NULL", id)
	err := row.Scan(&p.ID, &p.CreatedAt, &p.DeletedAt, &p.LastModifiedAt, &p.Confirmed, &p.PhoneNumber, &p.AccountID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return p.ToProtobuf(), nil
}

func (tx *txImpl) GetDeletedEmailAddressByID(ctx context.Context) (*accountV1.EmailAddress, error) {
	var e entities.EmailAddress

	row := tx.tx.QueryRow(
		ctx,
		"SELECT id, created_at, deleted_at, last_modified_at, confirmed, email_address, account_id "+
			"FROM email_address WHERE id=$1 "+
			"AND deleted_at IS NOT NULL", e.ID)

	err := row.Scan(&e.ID, &e.CreatedAt, &e.DeletedAt, &e.LastModifiedAt, &e.Confirmed, &e.EmailAddress, &e.AccountID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return e.ToProtobuf(), nil
}

// UpdateEmailAddress updates only the confirmation status of an email address.
// TODO rename func
func (tx *txImpl) UpdateEmailAddress(ctx context.Context) error {
	var e entities.EmailAddress

	_, err := tx.tx.Exec(
		ctx,
		"UPDATE email_address SET last_modified_at=$1, confirmed=$2 WHERE id=$3",
		time.Now(), e.Confirmed, e.ID)

	return err
}

func (tx *txImpl) DeleteEmailAddress(ctx context.Context, id string) error {
	_, err := tx.tx.Exec(ctx, "DELETE FROM email_address WHERE id=$1", id)

	return err
}

func (tx *txImpl) CreateEmailAddress(ctx context.Context, e entities.EmailAddress) error {
	_, err := tx.tx.Exec(
		ctx,
		"INSERT INTO email_address(id, account_id, email_address, confirmed, is_primary) VALUES($1, $2, $3, $4, $5)",
		e.ID, e.AccountID, e.EmailAddress, e.Confirmed, e.Primary)

	return err
}

func (tx *txImpl) GetEmailAddresses(ctx context.Context, request paginationV1.CursorRequest) ([]*accountV1.EmailAddress, error) {
	var sortString string
	if len(request.SortClauses) == 0 {
		sortString = "id ASC"
	} else {
		sortString = buildOrderByClause(request)
	}

	rows, err := tx.tx.Query(
		ctx,
		fmt.Sprintf(
			"SELECT id, created_at, deleted_at, last_modified_at, confirmed, email_address, account_id "+
				"FROM email_address "+
				"WHERE id > $1 "+
				"AND deleted_at IS NULL "+
				"ORDER BY %s "+
				"FETCH FIRST %d ROWS ONLY", sortString, request.Count), request.Cursor)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	emailAddresses := []*accountV1.EmailAddress{}

	for rows.Next() {
		var e entities.EmailAddress
		if err := rows.Scan(&e.ID, &e.CreatedAt, &e.DeletedAt, &e.LastModifiedAt, &e.Confirmed, &e.EmailAddress, &e.AccountID); err != nil {
			return nil, err
		}
		emailAddresses = append(emailAddresses, e.ToProtobuf())
	}

	return emailAddresses, nil
}

func (tx *txImpl) GetEmailAddressesForAccount(ctx context.Context, accountID string, cursorRequest paginationV1.CursorRequest) ([]*accountV1.EmailAddress, error) {
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

func (tx *txImpl) GetPhoneNumbersForAccount(ctx context.Context, accountID string, cursorRequest paginationV1.CursorRequest) ([]*accountV1.PhoneNumber, error) {
	var sortString string
	if len(cursorRequest.SortClauses) == 0 {
		sortString = "id ASC"
	} else {
		sortString = buildOrderByClause(cursorRequest)
	}

	queryTemplate := `
SELECT id, phone_number, account_id 
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
		var p accountV1.PhoneNumber
		if err := rows.Scan(&p.Id, &p.PhoneNumber, &p.AccountId); err != nil {
			return nil, err
		}
		p.PhoneNumber = maskPhoneNumber(p.PhoneNumber)
		phoneNumbers = append(phoneNumbers, &p)
	}

	return phoneNumbers, nil
}

func (tx *txImpl) EmailIsConfirmed(ctx context.Context, emailAddress string) (bool, error) {
	var count int
	row := tx.tx.QueryRow(
		ctx,
		"SELECT COUNT(*) AS count "+
			"FROM email_address "+
			"WHERE email_address = $1 "+
			"AND confirmed = $2 "+
			"AND deleted_at IS NULL", emailAddress, true)
	err := row.Scan(&count)

	if err != nil {
		// TODO check for NotFound?
		return false, err
	}
	return count == 1, nil
}

func (tx *txImpl) EmailExists(ctx context.Context, emailAddress string) (bool, error) {
	count, err := tx.CountEmail(ctx, emailAddress)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return count == 1, nil
}

func (tx *txImpl) CountEmail(ctx context.Context, emailAddress string) (int, error) {
	var count int
	row := tx.tx.QueryRow(
		ctx,
		"SELECT COUNT(*) AS count FROM email_address WHERE email_address=$1 AND deleted_at IS NULL", emailAddress)
	err := row.Scan(&count)
	return count, err
}

func (tx *txImpl) GetConfirmedEmailAddress(ctx context.Context) (*accountV1.EmailAddress, error) {
	var e entities.EmailAddress

	row := tx.tx.QueryRow(
		ctx,
		"SELECT id, email_address, account_id "+
			"FROM email_address WHERE email_address=$1 "+
			"AND confirmed=$2 "+
			"AND deleted_at IS NULL", e.EmailAddress, true)

	err := row.Scan(&e.ID, &e.EmailAddress, &e.AccountID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return e.ToProtobuf(), nil
}

func (tx *txImpl) GetPhoneNumberByPhoneNumber(ctx context.Context, phoneNumber string) (*accountV1.PhoneNumber, error) {
	var p accountV1.PhoneNumber

	err := tx.tx.QueryRow(
		ctx,
		"SELECT phone_number, account_id "+
			"FROM phone_number WHERE phone_number=$1 "+
			"AND deleted_at IS NULL", phoneNumber).Scan(&p.PhoneNumber, &p.AccountId)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &p, nil
}

func sortKeyword(sort paginationV1.Sort) string {
	if sort == paginationV1.Sort_SORT_DESC {
		return "DESC"
	}
	return "ASC"
}

func buildOrderByClause(request paginationV1.CursorRequest) string {
	var arr []string
	for _, sortClause := range request.SortClauses {
		sortString := sortKeyword(sortClause.Sort)
		arr = append(arr, fmt.Sprintf("%s %s", sortClause.FieldName, sortString))
	}
	return strings.Join(arr, ", ")
}

func maskPhoneNumber(phoneNumber string) string {
	return phoneNumber[len(phoneNumber)-2:]
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
