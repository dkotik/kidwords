package sqlite

import (
	"context"
	"fmt"
	"time"

	"github.com/dkotik/kidwords/service/secret"
)

func (r *sqRepository) Retrieve(
	ctx context.Context,
	ID string,
) (result secret.Secret, err error) {
	defer r.BindContext(ctx)()
	if err = r.stmtRetrieve.Reset(); err != nil {
		return result, err
	}
	result.ID = ID
	r.stmtRetrieve.BindText(1, ID)

	ok := false
	t := time.Time{}
	for {
		ok, err = r.stmtRetrieve.Step()
		if err != nil || !ok {
			break
		}
		result.UserID = r.stmtRetrieve.ColumnText(0)
		result.Name = r.stmtRetrieve.ColumnText(1)
		result.Type = r.stmtRetrieve.ColumnText(2)
		result.SaltedHash = r.stmtRetrieve.ColumnText(3)
		t, err = decodeTime(r.stmtRetrieve.ColumnText(4))
		if err != nil {
			err = fmt.Errorf("unable to decode created_at time: %w", err)
			break
		}
		result.CreatedAt = t
		t, err = decodeTime(r.stmtRetrieve.ColumnText(5))
		if err != nil {
			err = fmt.Errorf("unable to decode updated_at time: %w", err)
			break
		}
		result.UpdatedAt = t
		if !r.stmtRetrieve.ColumnIsNull(6) {
			t, err = decodeTime(r.stmtRetrieve.ColumnText(6))
			if err != nil {
				err = fmt.Errorf("unable to decode last_accepted_at time: %w", err)
				break
			}
			result.LastAcceptedAt = t
		}
	}
	if result.Type == "" && result.CreatedAt.IsZero() {
		return result, secret.ErrNotFound
	}
	// fmt.Printf("%+v\n", result)
	return result, err
}

func (r *sqRepository) List(
	ctx context.Context,
	userID string,
) (secrets []secret.Secret, err error) {
	defer r.BindContext(ctx)()
	if err = r.stmtList.Reset(); err != nil {
		return nil, err
	}
	r.stmtList.BindText(1, userID)

	ok := false
	var t time.Time
	result := secret.Secret{}
	for {
		ok, err = r.stmtList.Step()
		if err != nil || !ok {
			break
		}
		result.UserID = userID
		result.ID = r.stmtList.ColumnText(0)
		result.Name = r.stmtList.ColumnText(1)
		result.Type = r.stmtList.ColumnText(2)
		result.SaltedHash = r.stmtList.ColumnText(3)
		t, err = decodeTime(r.stmtList.ColumnText(4))
		if err != nil {
			err = fmt.Errorf("unable to decode created_at time: %w", err)
			break
		}
		result.CreatedAt = t
		t, err = decodeTime(r.stmtList.ColumnText(5))
		if err != nil {
			err = fmt.Errorf("unable to decode updated_at time: %w", err)
			break
		}
		result.UpdatedAt = t
		if !r.stmtList.ColumnIsNull(6) {
			t, err = decodeTime(r.stmtList.ColumnText(6))
			if err != nil {
				err = fmt.Errorf("unable to decode last_accepted_at time: %w", err)
				break
			}
			result.LastAcceptedAt = t
		}
		secrets = append(secrets, result)
	}
	return secrets, err
}
