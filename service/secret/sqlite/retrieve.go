package sqlite

import (
	"context"

	"github.com/dkotik/kidwords/service/secret"
)

func (s *sqRepository) Retrieve(
	ctx context.Context,
	ID string,
) (result secret.Secret, err error) {

	return result, nil
}

func (s *sqRepository) List(
	ctx context.Context,
	userID string,
) (result []secret.Secret, err error) {
	// rows, err := s.stmtRetrieveAll.QueryContext(ctx, userID)
	// if err != nil {
	// 	return nil, err
	// }
	// defer rows.Close()

	// var created int64
	// for rows.Next() {
	// 	key := &PaperKey{}
	// 	if err := rows.Scan(
	// 		&key.ID,
	// 		&key.Owner,
	// 		&key.Name,
	// 		&key.SaltedHash,
	// 		&created,
	// 	); err != nil {
	// 		return nil, err
	// 	}
	// 	key.Created = time.Unix(created, 0)
	// 	result = append(result, key)
	// }
	// if err = rows.Err(); err != nil {
	// 	return nil, err
	// }

	return result, nil
}
