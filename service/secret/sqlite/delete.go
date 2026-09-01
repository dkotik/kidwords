package sqlite

import "context"

func (r *sqRepository) Delete(ctx context.Context, ID string) (err error) {
	defer r.BindContext(ctx)()
	if err = r.stmtDelete.Reset(); err != nil {
		return err
	}
	r.stmtDelete.BindText(1, ID)
	ok := false
	for {
		ok, err = r.stmtDelete.Step()
		if err != nil || !ok {
			break
		}
	}
	return err
}
