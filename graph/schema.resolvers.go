package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"whosdriving-be/graph/generated"
	"whosdriving-be/graph/model"
)

// FindOrCreateUser is the resolver for the findOrCreateUser field.
func (r *mutationResolver) FindOrCreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	panic(fmt.Errorf("not implemented: FindOrCreateUser - findOrCreateUser"))
}

// ChangeUserRole is the resolver for the changeUserRole field.
func (r *mutationResolver) ChangeUserRole(ctx context.Context, input model.NewRole) (*model.User, error) {
	panic(fmt.Errorf("not implemented: ChangeUserRole - changeUserRole"))
}

// AddRotation is the resolver for the addRotation field.
func (r *mutationResolver) AddRotation(ctx context.Context, input model.NewRotation) (*model.Rotation, error) {
	panic(fmt.Errorf("not implemented: AddRotation - addRotation"))
}

// AddRide is the resolver for the addRide field.
func (r *mutationResolver) AddRide(ctx context.Context, input model.NewRide) (*model.Ride, error) {
	panic(fmt.Errorf("not implemented: AddRide - addRide"))
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, email string) (*model.User, error) {
	panic(fmt.Errorf("not implemented: User - user"))
}

// Rotations is the resolver for the rotations field.
func (r *queryResolver) Rotations(ctx context.Context, email *string) ([]*model.Rotation, error) {
	panic(fmt.Errorf("not implemented: Rotations - rotations"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
