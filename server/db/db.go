package db

import (
	"context"
	"database/sql"
	"time"
)

func GetEmailsForVerification(db *sql.DB, fileId int64) ([]string, error) {
	query := `call sp_get_emails_for_verification(?)`

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	stmt, err := db.PrepareContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, fileId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	emails := []string{}

	for rows.Next() {
		var email string

		if err := rows.Scan(&email); err != nil {
			return nil, err
		}

		emails = append(emails, email)
	}

	return emails, nil
}

func GetToVerifyCount(db *sql.DB, fileId int64) (int64, error) {
	query := `
	select count(*)
	from emails
	where (file_id = ?) and (error_msg is null or error_msg != '')`

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	stmt, err := db.PrepareContext(ctx, query)

	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, fileId)

	if err = row.Err(); row.Err() != nil {
		return 0, err
	}

	var count int64 = 0
	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	if err != nil {
		return 0, err
	}

	return count, nil
}

