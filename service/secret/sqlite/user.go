package sqlite

import (
	"context"
	"errors"
	"strings"
	"time"
	"uuid"

	"github.com/dkotik/kidwords/service/secret"
	lib "modernc.org/sqlite/lib"
	"zombiezen.com/go/sqlite"
)

var (
	ErrDuplicateUser      = errors.New("duplicate user")
	ErrUserNotFound       = errors.New("user not found")
	ErrDuplicateUserName  = errors.New("duplicate user name")
	ErrDuplicateUserEmail = errors.New("duplicate user email")
)

type User struct {
	ID            string
	Name          string
	PasswordHash  string
	Email         string
	EmailVerified bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ActiveAt      time.Time
}

func (u *User) GetID() string {
	return u.ID
}

func (u *User) GetName() string {
	return u.Name
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) GetEmailVerified() bool {
	return u.EmailVerified
}

func (u *User) GetPasswordHash() string {
	return u.PasswordHash
}

func (r *sqRepository) CreateUser(ctx context.Context, u *User) (err error) {
	defer r.BindContext(ctx)()
	if err = r.stmtCreateUser.Reset(); err != nil {
		return err
	}
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}
	if u.UpdatedAt.IsZero() {
		u.UpdatedAt = u.CreatedAt
	}

	r.stmtCreateUser.BindText(1, u.ID)
	r.stmtCreateUser.BindText(2, u.Name)
	r.stmtCreateUser.BindText(3, u.PasswordHash)
	r.stmtCreateUser.BindText(4, u.Email)
	r.stmtCreateUser.BindBool(5, u.EmailVerified)
	r.stmtCreateUser.BindText(6, encodeTime(u.CreatedAt))
	r.stmtCreateUser.BindText(7, encodeTime(u.UpdatedAt))
	r.stmtCreateUser.BindText(8, encodeTime(u.ActiveAt))

	ok := false
	for {
		ok, err = r.stmtCreateUser.Step()
		if err != nil {
			switch code := sqlite.ErrCode(err); code {
			case lib.SQLITE_OK:
				return nil
			case lib.SQLITE_CONSTRAINT_PRIMARYKEY:
				return ErrDuplicateUser
			case lib.SQLITE_CONSTRAINT_UNIQUE:
				// if strings.HasSuffix(err.Error(), "name") {
				// 	return ErrDuplicateUserName
				// }
				if strings.HasSuffix(err.Error(), "email") {
					return ErrDuplicateUserEmail
				}
				return ErrDuplicateUserName
			default:
				return err
			}
		}
		if !ok {
			break
		}
	}
	// if r.conn.LastInsertRowID() == 0 {
	// 	return errors.New("never created")
	// }
	// fmt.Println("===============", r.conn.LastInsertRowID(), "||||")
	return nil
}

func (r *sqRepository) newUserFromRow(u *User, stmt *sqlite.Stmt) (err error) {
	u.ID = stmt.ColumnText(0)
	u.Name = stmt.ColumnText(1)
	u.PasswordHash = stmt.ColumnText(2)
	u.Email = stmt.ColumnText(3)
	u.EmailVerified = stmt.ColumnBool(4)
	u.CreatedAt, err = decodeTime(stmt.ColumnText(5))
	if err != nil {
		return err
	}
	u.UpdatedAt, err = decodeTime(stmt.ColumnText(6))
	if err != nil {
		return err
	}
	u.ActiveAt, err = decodeTime(stmt.ColumnText(7))
	if err != nil {
		return err
	}
	return nil
}

func (r *sqRepository) RetrieveUser(
	ctx context.Context,
	ID string,
) (u *User, err error) {
	u = &User{}
	defer r.BindContext(ctx)()
	if err := r.stmtRetrieveUser.Reset(); err != nil {
		return u, err
	}
	r.stmtRetrieveUser.BindText(1, ID)
	ok := false
	for {
		ok, err = r.stmtRetrieveUser.Step()
		if err != nil {
			return u, err
		}
		if !ok {
			break
		}
		err = r.newUserFromRow(u, r.stmtRetrieveUser)
		if err != nil {
			return u, err
		}
	}
	if u.ID == "" {
		return u, ErrUserNotFound
	}
	return u, nil
}

func (r *sqRepository) RetrieveUserByName(
	ctx context.Context,
	name string,
) (_ secret.User, err error) {
	u := &User{}
	defer r.BindContext(ctx)()
	if err := r.stmtRetrieveUserByName.Reset(); err != nil {
		return u, err
	}
	r.stmtRetrieveUserByName.BindText(1, name)
	ok := false
	for {
		ok, err = r.stmtRetrieveUserByName.Step()
		if err != nil {
			return u, err
		}
		if !ok {
			break
		}
		// fmt.Printf("retrieved user: %+v\n", u)
		err = r.newUserFromRow(u, r.stmtRetrieveUserByName)
		if err != nil {
			return u, err
		}
	}
	if u.ID == "" {
		return u, ErrUserNotFound
	}
	return u, nil
}

func (r *sqRepository) RetrieveUserByEmailAddress(
	ctx context.Context,
	email string,
) (_ secret.User, err error) {
	u := &User{}
	defer r.BindContext(ctx)()
	if err := r.stmtRetrieveUserByEmail.Reset(); err != nil {
		return u, err
	}
	r.stmtRetrieveUserByEmail.BindText(1, email)
	ok := false
	for {
		ok, err = r.stmtRetrieveUserByEmail.Step()
		if err != nil {
			return u, err
		}
		if !ok {
			break
		}
		err = r.newUserFromRow(u, r.stmtRetrieveUserByEmail)
		if err != nil {
			return u, err
		}
	}
	if u.ID == "" {
		return u, ErrUserNotFound
	}
	return u, nil
}

func (r *sqRepository) UpdateUser(ctx context.Context, u *User) (err error) {
	defer r.BindContext(ctx)()
	if err := r.stmtUpdateUser.Reset(); err != nil {
		return err
	}
	u.UpdatedAt = time.Now()
	r.stmtUpdateUser.BindText(1, u.ID)
	r.stmtUpdateUser.BindText(2, u.Name)
	r.stmtUpdateUser.BindText(3, u.PasswordHash)
	r.stmtUpdateUser.BindText(4, u.Email)
	r.stmtUpdateUser.BindBool(5, u.EmailVerified)
	r.stmtUpdateUser.BindText(6, encodeTime(u.CreatedAt))
	r.stmtUpdateUser.BindText(7, encodeTime(u.UpdatedAt))
	r.stmtUpdateUser.BindText(8, encodeTime(u.ActiveAt))

	ok := false
	for {
		ok, err = r.stmtUpdateUser.Step()
		if err != nil {
			return err
		}
		if !ok {
			break
		}
	}
	return nil
}

func (r *sqRepository) MarkUserAsActive(ctx context.Context, ID string) (err error) {
	defer r.BindContext(ctx)()
	user, err := r.RetrieveUser(ctx, ID)
	if err != nil {
		return err
	}
	user.ActiveAt = time.Now()
	if err := r.UpdateUser(ctx, user); err != nil {
		return err
	}
	return nil
}

func (r *sqRepository) DeleteUser(ctx context.Context, ID string) (err error) {
	defer r.BindContext(ctx)()
	if err := r.stmtDeleteUser.Reset(); err != nil {
		return err
	}
	r.stmtDeleteUser.BindText(1, ID)
	ok := false
	for {
		ok, err = r.stmtDeleteUser.Step()
		if err != nil {
			return err
		}
		if !ok {
			break
		}
	}
	return nil
}
