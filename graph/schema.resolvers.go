package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"database/sql"
	"log"
	"whosdriving-be/data_interface"
	"whosdriving-be/graph/generated"
	"whosdriving-be/graph/model"
)

// FindOrCreateUser is the resolver for the findOrCreateUser field.
func (r *mutationResolver) FindOrCreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	tx, err := r.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	lCtx := data_interface.LuwContext{Conn: r.DB, Tx: tx}
	user, err := data_interface.FindUser(ctx, &lCtx, &input.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			user, err = data_interface.CreateUser(ctx, &lCtx, &input)
		}
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return user, nil
}

// ChangeUserRole is the resolver for the changeUserRole field.
func (r *mutationResolver) ChangeUserRole(ctx context.Context, input model.NewRole) (*model.User, error) {
	tx, err := r.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	lCtx := data_interface.LuwContext{Conn: r.DB, Tx: tx}
	user, err := data_interface.FindUser(ctx, &lCtx, &input.Email)
	if err != nil {
		return nil, err
	}

	user.Role = input.Role
	usr, err := data_interface.UpdateUser(ctx, &lCtx, user)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return usr, err
}

// AddRotation is the resolver for the addRotation field.
func (r *mutationResolver) AddRotation(ctx context.Context, input model.NewRotation) (*model.Rotation, error) {
	tx, err := r.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	lCtx := data_interface.LuwContext{Conn: r.DB, Tx: tx}
	rotation, err := data_interface.CreateRotation(ctx, &lCtx, &input)
	if err != nil {
		return nil, err
	}

	log.Printf("Create successfuly new rotation %T", rotation)

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return rotation, nil
}

// AddRide is the resolver for the addRide field.
func (r *mutationResolver) AddRide(ctx context.Context, input model.NewRide) (*model.Ride, error) {
	tx, err := r.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	lCtx := data_interface.LuwContext{Conn: r.DB, Tx: tx}
	ride, err := data_interface.AddRide(ctx, &lCtx, &input)
	if err != nil {
		return nil, err
	}

	log.Printf("Create successfuly new ride %T", ride)

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return ride, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, email string) (*model.User, error) {
	tx, err := r.DB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  true,
	})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	lCtx := data_interface.LuwContext{Conn: r.DB, Tx: tx}
	user, err := data_interface.FindUser(ctx, &lCtx, &email)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Rotations is the resolver for the rotations field.
func (r *queryResolver) Rotations(ctx context.Context, email *string) ([]*model.Rotation, error) {
	tx, err := r.DB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  true,
	})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	lCtx := data_interface.LuwContext{Conn: r.DB, Tx: tx}
	rotations, err := data_interface.FindRotations(ctx, &lCtx, email)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return rotations, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
