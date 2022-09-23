package data_interface

import (
	"context"
	"database/sql"
	"log"
	"whosdriving-be/graph/model"
)

func FindUser(ctx context.Context, lCtx *LuwContext, email *string) (*model.User, error) {
	const q string = `select email, firstname, lastname, profile, RefRole.RefName
						from Users left join RefRole on Users.roleCd = RefRole.RefCd 
						where email = ? and deleteTmstmp is null`
	user := new(model.User)
	if err := lCtx.Tx.QueryRowContext(ctx, q, &email).Scan(&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Profile,
		&user.Role); err != nil {
		return nil, err
	}
	return user, nil
}

func FindUsers(ctx context.Context, lCtx *LuwContext, emails *[]string) ([]*model.User, error) {
	users := make([]*model.User, 0, len(*emails))
	for _, email := range *emails {
		user, err := FindUser(ctx, lCtx, &email)
		if err != nil {
			if err != sql.ErrNoRows {
				log.Printf("Error: when start loading user %s - %s", email, err)
				return nil, err
			}
			log.Printf("No user found for %s", email)
			continue
		}
		users = append(users, user)
	}
	return users, nil
}

func CreateUser(ctx context.Context, lCtx *LuwContext, newUser *model.NewUser) (*model.User, error) {
	const q string = `INSERT INTO Users(email, firstname, lastname, profile, roleCd, createTmstmp, lstUpdTmstmp, deleteTmstmp) 
	VALUES (?, ?, ?, ?, (select RefCd from RefRole where RefName='STANDARD'), DATETIME('now'), DATETIME('now'), null)`

	stmt, err := lCtx.Tx.PrepareContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, newUser.Email, newUser.FirstName, newUser.LastName, newUser.Profile)
	if err != nil {
		return nil, err
	}

	return FindUser(ctx, lCtx, &newUser.Email)
}

func UpdateUser(ctx context.Context, lCtx *LuwContext, user *model.User) (*model.User, error) {
	const q string = `UPDATE Users set firstname=?, lastname=?, profile=?, roleCd=(select refCd from RefRole where RefName=?), lstUpdTmstmp=DATETIME('now') 
				WHERE email=? and deleteTmstmp is null`

	stmt, err := lCtx.Tx.PrepareContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, user.FirstName, user.LastName, user.Profile, user.Role, user.Email)
	if err != nil {
		return nil, err
	}

	return FindUser(ctx, lCtx, &user.Email)
}

func DeleteUser(ctx context.Context, lCtx *LuwContext, user *model.User) (*model.User, error) {
	const q string = `UPDATE Users set deleteTmstmp=DATETIME('now'), lstUpdTmstmp=DATETIME('now') WHERE email=? and deleteTmstmp is null`

	stmt, err := lCtx.Tx.PrepareContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, user.Email)
	if err != nil {
		return nil, err
	}

	return user, nil
}
