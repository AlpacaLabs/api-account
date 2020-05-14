package entities

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-sql/sqlexp"
	"github.com/guregu/null"
)

type Account struct {
	ID                    string
	CreatedAt             time.Time
	LastModifiedAt        time.Time
	DeletedAt             null.Time
	Username              null.String
	CurrentPasswordID     null.String
	PrimaryEmailAddressID null.String
}

func (p *Account) GetAccountByUsername(q sqlexp.Querier) error {
	return q.QueryRowContext(
		context.TODO(),
		"SELECT p.id, p.created_timestamp, p.deleted_timestamp, p.last_modified_timestamp, "+
			"p.username, p.current_password_id, p.primary_email_address_id "+
			"FROM Account p "+
			"WHERE p.username=$1 "+
			"AND p.deleted_timestamp IS NULL", p.Username).Scan(&p.ID, &p.CreatedAt, &p.DeletedAt, &p.LastModifiedAt,
		&p.Username, &p.CurrentPasswordID, &p.PrimaryEmailAddressID)
}

func GetAccountByEmailAddress(q sqlexp.Querier, emailAddress string) (*Account, error) {
	p := &Account{}
	return p, q.QueryRowContext(
		context.TODO(),
		"SELECT p.id, p.created_timestamp, p.deleted_timestamp, p.last_modified_timestamp, "+
			"p.username, p.current_password_id, p.primary_email_address_id "+
			"FROM email_address e JOIN Account p ON e.accountId = p.id "+
			"WHERE e.email_address=$1 "+
			"AND p.deleted_timestamp IS NULL", emailAddress).Scan(&p.ID, &p.CreatedAt, &p.DeletedAt, &p.LastModifiedAt,
		&p.Username, &p.CurrentPasswordID, &p.PrimaryEmailAddressID)
}

func GetAccountByPhoneNumber(q sqlexp.Querier, phoneNumber string) (*Account, error) {
	p := &Account{}
	return p, q.QueryRowContext(
		context.TODO(),
		"SELECT p.id, p.created_timestamp, p.deleted_timestamp, p.last_modified_timestamp, "+
			"p.username, p.current_password_id, p.primary_email_address_id "+
			"FROM phone_number pn JOIN Account p ON pn.accountId = p.id "+
			"WHERE pn.phone_number=$1 "+
			"AND p.deleted_timestamp IS NULL", phoneNumber).Scan(&p.ID, &p.CreatedAt, &p.DeletedAt, &p.LastModifiedAt,
		&p.Username, &p.CurrentPasswordID, &p.PrimaryEmailAddressID)
}

func (p *Account) GetDeletedAccount(q sqlexp.Querier) error {
	return q.QueryRowContext(
		context.TODO(),
		"SELECT p.id, p.created_timestamp, p.deleted_timestamp, p.last_modified_timestamp, "+
			"p.username, p.current_password_id, p.primary_email_address_id "+
			"FROM Account p "+
			"WHERE p.id=$1 "+
			"AND p.deleted_timestamp IS NOT NULL", p.ID).Scan(&p.ID, &p.CreatedAt, &p.DeletedAt, &p.LastModifiedAt,
		&p.Username, &p.CurrentPasswordID, &p.PrimaryEmailAddressID)
}

func (p *Account) GetAccount(q sqlexp.Querier) error {
	return q.QueryRowContext(
		context.TODO(),
		"SELECT p.id, p.created_timestamp, p.deleted_timestamp, p.last_modified_timestamp, "+
			"p.username, p.current_password_id, p.primary_email_address_id "+
			"FROM Account p "+
			"WHERE p.id=$1 "+
			"AND p.deleted_timestamp IS NULL", p.ID).Scan(&p.ID, &p.CreatedAt, &p.DeletedAt, &p.LastModifiedAt,
		&p.Username, &p.CurrentPasswordID, &p.PrimaryEmailAddressID)
}

func (p *Account) UpdateAccount(q sqlexp.Querier) error {
	_, err := q.ExecContext(
		context.TODO(),
		"UPDATE Account "+
			"SET last_modified_timestamp=$1, username=$2, primary_email_address_id=$3 "+
			"WHERE id=$4",
		time.Now(), p.Username, p.PrimaryEmailAddressID, p.ID)
	return err
}

func (p *Account) DeleteAccount(q sqlexp.Querier) error {
	_, err := q.ExecContext(
		context.TODO(),
		"UPDATE Account SET deleted_timestamp=$1 WHERE id=$2",
		time.Now(), p.ID)
	return err
}

func (p *Account) CreateAccount(q sqlexp.Querier) error {
	_, err := q.ExecContext(
		context.TODO(),
		"INSERT INTO Account(id, created_timestamp, username) VALUES($1, $2, $3)",
		p.ID, time.Now(), p.Username)

	if err != nil {
		return err
	}

	return nil
}

func GetAccounts(q sqlexp.Querier, cursor string, sort string, count int) ([]Account, error) {
	rows, err := q.QueryContext(
		context.TODO(),
		fmt.Sprintf(
			"SELECT p.id, p.created_timestamp, p.deleted_timestamp, p.last_modified_timestamp, "+
				"p.username, p.current_password_id, p.primary_email_address_id "+
				"FROM Account p "+
				"WHERE p.id > $1 "+
				"ORDER BY p.id %s "+
				"FETCH FIRST %d ROWS ONLY", sort, count), cursor)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	accounts := []Account{}

	for rows.Next() {
		var p Account
		if err := rows.Scan(&p.ID, &p.CreatedAt, &p.DeletedAt, &p.LastModifiedAt,
			&p.Username, &p.CurrentPasswordID, &p.PrimaryEmailAddressID); err != nil {
			return nil, err
		}
		accounts = append(accounts, p)
	}

	return accounts, nil
}

func (p *Account) Exists(q sqlexp.Querier) (bool, error) {
	count, err := p.Count(q)
	if err != nil {
		return false, err
	}
	return count == 1, nil
}

func (p *Account) Count(q sqlexp.Querier) (int, error) {
	var count int
	row := q.QueryRowContext(context.TODO(), "SELECT COUNT(*) AS count FROM Account WHERE id=$1 AND deleted_timestamp IS NULL", p.ID)
	err := row.Scan(&count)
	return count, err
}

func (p *Account) UpdateCurrentPassword(q sqlexp.Querier, currentPasswordID string) error {
	_, err := q.ExecContext(
		context.TODO(),
		"UPDATE Account "+
			"SET last_modified_timestamp=$1, current_password_id=$2 "+
			"WHERE id=$3",
		time.Now(), currentPasswordID, p.ID)
	return err
}
