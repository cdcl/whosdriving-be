package data_interface

import (
	"context"
	"database/sql"
	"log"
	"whosdriving-be/graph/model"
)

func FindRide(ctx context.Context, lCtx *LuwContext, id int64) (*model.Ride, error) {
	const q string = `select id, riderEmail 
						from rides r 
						where r.id=? and r.deleteTmstmp is null`

	ride := new(model.Ride)
	var riderEmail sql.NullString

	if err := lCtx.Tx.QueryRowContext(ctx, q, &id).Scan(&ride.ID,
		&riderEmail); err != nil {
		return nil, err
	}

	log.Printf("User lookup rider %s", riderEmail.String)
	rider, err := FindUser(ctx, lCtx, &riderEmail.String)
	if err != nil {
		return nil, err
	}
	ride.Conductor = rider

	participants, err := FindRideParticipants(ctx, lCtx, int64(ride.ID))
	if err != nil {
		return nil, err
	}
	ride.Participants = participants

	return ride, nil
}

func FindRides(ctx context.Context, lCtx *LuwContext, rotationId int64) ([]*model.Ride, error) {
	const q string = `select id 
						from rides r 
						where r.rotationId=? and r.deleteTmstmp is null`
	rows, err := lCtx.Tx.QueryContext(ctx, q, rotationId)
	switch {
	case err == sql.ErrNoRows:
		rows.Close()
		return nil, nil
	case err != nil:
		return nil, err
	}
	defer rows.Close()

	rides := make([]*model.Ride, 0)

	for rows.Next() {
		var rideId sql.NullInt64
		if rows.Err() != sql.ErrNoRows {
			if err := rows.Scan(&rideId); err != nil {
				// Check for a scan error.
				return nil, err
			}
			if rideId.Valid {
				ride, err := FindRide(ctx, lCtx, rideId.Int64)
				if err != nil {
					return nil, err
				}
				rides = append(rides, ride)
			}
		}
	}

	if len(rides) == 0 {
		return nil, nil
	}

	return rides, nil
}

func AddRide(ctx context.Context, lCtx *LuwContext, newRide *model.NewRide) (*model.Ride, error) {
	const q string = `INSERT INTO Rides(rotationId, riderEmail, createTmstmp, lstUpdTmstmp, deleteTmstmp) 
						VALUES (?, ?, DATETIME('now'), DATETIME('now'), null)`

	stmt, err := lCtx.Tx.PrepareContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, newRide.IDRotation, newRide.EmailConductor)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	log.Printf("Create ride %s assign id %d", newRide.EmailConductor, id)
	err = CreateRideParticipants(ctx, lCtx, id, &newRide.EmailParticipants)
	if err != nil {
		return nil, err
	}

	return FindRide(ctx, lCtx, id)
}

func UpdateRide(ctx context.Context, lCtx *LuwContext, ride *model.Ride) (*model.Ride, error) {
	const q string = `UPDATE Rides set riderEmail=?, lstUpdTmstmp=DATETIME('now') 
				WHERE id=? and deleteTmstmp is null`

	stmt, err := lCtx.Tx.PrepareContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, ride.Conductor.Email, ride.ID)
	if err != nil {
		return nil, err
	}

	log.Printf("Update rotation id %d", ride.ID)
	return FindRide(ctx, lCtx, int64(ride.ID))
}

func DeleteRide(ctx context.Context, lCtx *LuwContext, ride *model.Ride) (*model.Ride, error) {
	const q string = `UPDATE Rides set lstUpdTmstmp=DATETIME('now'), deleteTmstmp=DATETIME('now')
				WHERE id=? and deleteTmstmp is null`

	stmt, err := lCtx.Tx.PrepareContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, ride.ID)
	if err != nil {
		return nil, err
	}

	log.Printf("Delete rotation id %d", ride.ID)
	return ride, nil
}

func FindRideParticipants(ctx context.Context, lCtx *LuwContext, id int64) ([]*model.User, error) {
	const q string = `select p.email 
						from Rides r left join RideParticipants p on p.rideId = r.Id 
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

	log.Printf("User lookup ride partitipants %s", participantEmails)
	participants, err := FindUsers(ctx, lCtx, &participantEmails)
	if err != nil {
		return nil, err
	}
	return participants, nil
}

func CreateRideParticipants(ctx context.Context, lCtx *LuwContext, rideId int64, participantsEmails *[]string) error {
	const q string = `INSERT INTO RideParticipants (rideId, email) VALUES (?, ?)`

	stmt, err := lCtx.Tx.PrepareContext(ctx, q)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, participantEmail := range *participantsEmails {
		log.Printf("Add participant for rotation %d - %s", rideId, participantEmail)
		_, err = stmt.ExecContext(ctx, rideId, participantEmail)
		if err != nil {
			return err
		}
	}

	return nil
}

func RemoveRideParticipants(ctx context.Context, lCtx *LuwContext, rideId int64, participantsEmails *[]string) error {
	const q string = `DELETE from RideParticipants where rideId=? and email=?`

	stmt, err := lCtx.Tx.PrepareContext(ctx, q)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, participantEmail := range *participantsEmails {
		log.Printf("Remove participant for rotation %d - %s", rideId, participantEmail)
		_, err = stmt.ExecContext(ctx, rideId, participantEmail)
		if err != nil {
			return err
		}
	}

	return nil
}
