package data_interface

import (
	"context"
	"database/sql"
	"log"
	"whosdriving-be/graph/model"
)

func FindRotation(ctx context.Context, lCtx *LuwContext, id int64) (*model.Rotation, error) {
	const q string = `select id, name, creatorEmail 
						from Rotations r 
					   	where r.id = ? and r.deleteTmstmp is null`

	rotation := new(model.Rotation)
	var creatorEmail sql.NullString

	if err := lCtx.Tx.QueryRowContext(ctx, q, &id).Scan(&rotation.ID,
		&rotation.Name,
		&creatorEmail); err != nil {
		return nil, err
	}

	if !creatorEmail.Valid {
		return nil, nil
	}

	log.Printf("User lookup creator %s", creatorEmail.String)
	creator, err := FindUser(ctx, lCtx, &creatorEmail.String)
	if err != nil {
		return nil, err
	}
	rotation.Creator = creator

	participants, err := FindRotationParticipants(ctx, lCtx, id)
	if err != nil {
		return nil, err
	}
	rotation.Participants = participants

	rides, err := FindRides(ctx, lCtx, int64(rotation.ID))
	if err != nil {
		return nil, err
	}
	rotation.Rides = rides

	return rotation, nil
}

func FindRotations(ctx context.Context, lCtx *LuwContext, email *string) ([]*model.Rotation, error) {
	const q string = `select id
						from Rotations r 
					   	where r.creatorEmail = ? and r.deleteTmstmp is null`

	rows, err := lCtx.Tx.QueryContext(ctx, q, email)
	switch {
	case err == sql.ErrNoRows:
		rows.Close()
		return nil, nil
	case err != nil:
		return nil, err
	}
	defer rows.Close()

	rotations := make([]*model.Rotation, 0)

	for rows.Next() {
		var id sql.NullInt64
		if rows.Err() != sql.ErrNoRows {
			if err := rows.Scan(&id); err != nil {
				// Check for a scan error.
				return nil, err
			}
			if id.Valid {
				rotation, err := FindRotation(ctx, lCtx, id.Int64)
				if err != nil {
					return nil, err
				}
				rotations = append(rotations, rotation)
			}
		}
	}

	return rotations, nil
}

func CreateRotation(ctx context.Context, lCtx *LuwContext, newRot *model.NewRotation) (*model.Rotation, error) {
	const q string = `INSERT INTO Rotations(name, creatorEmail, createTmstmp, lstUpdTmstmp, deleteTmstmp) 
						VALUES (?, ?, DATETIME('now'), DATETIME('now'), null)`

	stmt, err := lCtx.Tx.PrepareContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, newRot.Name, newRot.EmailCreator)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	log.Printf("Create rotation %s assign id %d", newRot.Name, id)
	err = CreateRotationParticipants(ctx, lCtx, id, &newRot.EmailParticipants)
	if err != nil {
		return nil, err
	}

	return FindRotation(ctx, lCtx, id)
}

func UpdateRotation(ctx context.Context, lCtx *LuwContext, rotation *model.Rotation) (*model.Rotation, error) {
	const q string = `UPDATE Rotations set name=?, creatorEmail=?, lstUpdTmstmp=DATETIME('now') 
				WHERE id=? and deleteTmstmp is null`

	stmt, err := lCtx.Tx.PrepareContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, rotation.Name, rotation.Creator.Email, rotation.ID)
	if err != nil {
		return nil, err
	}

	log.Printf("Update rotation id %d", rotation.ID)
	return FindRotation(ctx, lCtx, int64(rotation.ID))
}

func DeleteRotation(ctx context.Context, lCtx *LuwContext, rotation *model.Rotation) (*model.Rotation, error) {
	const q string = `UPDATE Rotations set lstUpdTmstmp=DATETIME('now'), deleteTmstmp=DATETIME('now')
				WHERE id=? and deleteTmstmp is null`

	stmt, err := lCtx.Tx.PrepareContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, rotation.ID)
	if err != nil {
		return nil, err
	}

	log.Printf("Delete rotation id %d", rotation.ID)
	return rotation, nil
}

func FindRotationParticipants(ctx context.Context, lCtx *LuwContext, id int64) ([]*model.User, error) {
	const q string = `select p.email 
						from Rotations r left join RotationParticipants p on p.rotationId = r.Id 
					   where r.id = ? and r.deleteTmstmp is null
					   order by p.email`

	rows, err := lCtx.Tx.QueryContext(ctx, q, id)
	switch {
	case err == sql.ErrNoRows:
		rows.Close()
		return nil, nil
	case err != nil:
		return nil, err
	}
	defer rows.Close()

	participantEmails := make([]string, 0)

	for rows.Next() {
		var email sql.NullString
		if rows.Err() != sql.ErrNoRows {
			if err := rows.Scan(&email); err != nil {
				// Check for a scan error.
				return nil, err
			}
			if email.Valid {
				participantEmails = append(participantEmails, email.String)
			}
		}
	}

	if len(participantEmails) == 0 {
		return nil, nil
	}

	log.Printf("User lookup partitipants %s", participantEmails)
	participants, err := FindUsers(ctx, lCtx, &participantEmails)
	if err != nil {
		return nil, err
	}
	return participants, nil
}

func CreateRotationParticipants(ctx context.Context, lCtx *LuwContext, rotationId int64, participantsEmails *[]string) error {
	const q string = `INSERT INTO RotationParticipants (rotationId, email) VALUES (?, ?)`

	stmt, err := lCtx.Tx.PrepareContext(ctx, q)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, participantEmail := range *participantsEmails {
		log.Printf("Add participant for rotation %d - %s", rotationId, participantEmail)
		_, err = stmt.ExecContext(ctx, rotationId, participantEmail)
		if err != nil {
			return err
		}
	}

	return nil
}

func RemoveRotationParticipants(ctx context.Context, lCtx *LuwContext, rotationId int64, participantsEmails *[]string) error {
	const q string = `DELETE from RotationParticipants where rotationId=? and email=?`

	stmt, err := lCtx.Tx.PrepareContext(ctx, q)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, participantEmail := range *participantsEmails {
		log.Printf("Remove participant for rotation %d - %s", rotationId, participantEmail)
		_, err = stmt.ExecContext(ctx, rotationId, participantEmail)
		if err != nil {
			return err
		}
	}

	return nil
}
