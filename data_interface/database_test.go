package data_interface

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"whosdriving-be/graph/model"

	_ "github.com/mattn/go-sqlite3"

	"github.com/stretchr/testify/assert"
)

// Helper function to create easly new connections
func createNewDb(t *testing.T, dbPath string, ddlPath string) (db *sql.DB) {
	os.Remove(dbPath)

	db, err := NewConnection(dbPath)
	if err != nil {
		t.Fatalf("Could't create connection %s - %s", dbPath, err)
	}

	errMigration := Migrate(ddlPath, db)
	if errMigration != nil {
		db.Close()
		t.Fatalf("Migration error %s - %s", ddlPath, errMigration)
	}

	return db
}

func toNewUser(user *model.User) *model.NewUser {
	return &model.NewUser{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Profile:   user.Profile,
	}
}

func TestUser(t *testing.T) {
	firstName, lastName, profile := "test", "domain", "noProfile"
	expectedUser := model.User{
		Email:     "test@domain.com",
		FirstName: &firstName,
		LastName:  &lastName,
		Profile:   &profile,
		Role:      "STANDARD",
	}

	newUser := toNewUser(&expectedUser)

	ctx := context.Background()
	db := createNewDb(t, "../test_user.sqlite3", "../assets/ddl.whosdriving-core")
	if db == nil {
		t.Fatal("Could't create database connexion")
	}
	defer db.Close()

	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		t.Fatalf("Could't create transaction - %s", err)
	}
	defer tx.Rollback()

	// Create database context
	lCtx := LuwContext{Conn: db, Tx: tx}

	// create & find
	user, err := CreateUser(ctx, &lCtx, newUser)
	assert.Nil(t, err, "")
	assert.Equal(t, &expectedUser, user)

	// update & find
	*expectedUser.Profile = "beautifullProfile"
	user, err = UpdateUser(ctx, &lCtx, &expectedUser)
	assert.Nil(t, err, "")
	assert.Equal(t, &expectedUser, user)

	// delete
	_, err = DeleteUser(ctx, &lCtx, &expectedUser)
	assert.Nil(t, err, "")
	user, err = FindUser(ctx, &lCtx, &newUser.Email)
	assert.Equal(t, err, sql.ErrNoRows)
	assert.Nil(t, user, "User successfuly not found")

	err = tx.Commit()
	if err != nil {
		t.Fatalf("Error on commit - %s", err)
	}
}

func TestRotation(t *testing.T) {
	firstName, lastName, profile := "test", "domain", "noProfile"
	expectedCreator := model.User{
		Email:     "test@domain.com",
		FirstName: &firstName,
		LastName:  &lastName,
		Profile:   &profile,
		Role:      "STANDARD",
	}

	firstNameJohn, lastNameJohn, profileJohn := "John", "Smith", ""
	expectedParticipant1 := model.User{
		Email:     "john@domain.com",
		FirstName: &firstNameJohn,
		LastName:  &lastNameJohn,
		Profile:   &profileJohn,
		Role:      "STANDARD",
	}

	expectedRotation := model.Rotation{
		ID:           1,
		Name:         "TestRotation",
		Creator:      &expectedCreator,
		Participants: []*model.User{&expectedParticipant1, &expectedCreator},
		Rides:        nil,
	}

	newRotation := model.NewRotation{
		Name:              "TestRotation",
		EmailCreator:      "test@domain.com",
		EmailParticipants: []string{expectedCreator.Email, expectedParticipant1.Email},
	}

	ctx := context.Background()
	db := createNewDb(t, "../test_rotation.sqlite3", "../assets/ddl.whosdriving-core")
	if db == nil {
		t.Fatal("Couldn't create database connexion")
	}
	defer db.Close()

	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		t.Fatalf("Could't create transaction - %s", err)
	}
	defer tx.Rollback()

	// Create database context
	lCtx := LuwContext{Conn: db, Tx: tx}

	// create & find
	userTest, err := CreateUser(ctx, &lCtx, toNewUser(&expectedCreator))
	assert.Nil(t, err, "")
	assert.Equal(t, &expectedCreator, userTest)

	userJohn, err := CreateUser(ctx, &lCtx, toNewUser(&expectedParticipant1))
	assert.Nil(t, err, "")
	assert.Equal(t, &expectedParticipant1, userJohn)

	rotation, err := CreateRotation(ctx, &lCtx, &newRotation)
	assert.Nil(t, err, "")
	assert.Equal(t, &expectedRotation, rotation)

	updtExpectedRota := expectedRotation
	updtExpectedRota.Name = "Fancy Rotation Name"
	updtExpectedRota.Creator = &expectedParticipant1
	updtRotation, err := UpdateRotation(ctx, &lCtx, &updtExpectedRota)
	assert.Nil(t, err, "")
	assert.Equal(t, &updtExpectedRota, updtRotation)

	_, err = DeleteRotation(ctx, &lCtx, &updtExpectedRota)
	assert.Nil(t, err, "")
	rotation, err = FindRotation(ctx, &lCtx, int64(expectedRotation.ID))
	assert.Equal(t, err, sql.ErrNoRows)
	assert.Nil(t, rotation, "Rotation successfuly not found")

	phoenixRotation, err := CreateRotation(ctx, &lCtx, &newRotation)
	assert.Nil(t, err, "")
	expectedRotation.ID = 2
	assert.Equal(t, &expectedRotation, phoenixRotation)

	err = RemoveRotationParticipants(ctx, &lCtx, int64(expectedRotation.ID), &[]string{expectedRotation.Creator.Email})
	assert.Nil(t, err, "")

	expectedRotation2 := expectedRotation
	expectedRotation2.Participants = []*model.User{&expectedParticipant1}

	rotation, err = FindRotation(ctx, &lCtx, int64(expectedRotation2.ID))
	assert.Nil(t, err, "")
	assert.Equal(t, &expectedRotation2, rotation)

	err = tx.Commit()
	if err != nil {
		t.Fatalf("Error on commit - %s", err)
	}
}

func TestRide(t *testing.T) {
	firstName, lastName, profile := "test", "domain", "noProfile"
	expectedCreator := model.User{
		Email:     "test@domain.com",
		FirstName: &firstName,
		LastName:  &lastName,
		Profile:   &profile,
		Role:      "STANDARD",
	}

	firstNameJohn, lastNameJohn, profileJohn := "John", "Smith", ""
	expectedParticipant1 := model.User{
		Email:     "john@domain.com",
		FirstName: &firstNameJohn,
		LastName:  &lastNameJohn,
		Profile:   &profileJohn,
		Role:      "STANDARD",
	}

	expectedRotation := model.Rotation{
		ID:           1,
		Name:         "TestRotation",
		Creator:      &expectedCreator,
		Participants: []*model.User{&expectedParticipant1, &expectedCreator},
		Rides:        nil,
	}

	newRotation := model.NewRotation{
		Name:              "TestRotation",
		EmailCreator:      "test@domain.com",
		EmailParticipants: []string{expectedCreator.Email, expectedParticipant1.Email},
	}

	expectedRide := model.Ride{
		ID:           1,
		Conductor:    &expectedParticipant1,
		Participants: []*model.User{&expectedParticipant1, &expectedCreator},
	}

	newRide := model.NewRide{
		IDRotation:        1,
		EmailConductor:    "john@domain.com",
		EmailParticipants: []string{expectedCreator.Email, expectedParticipant1.Email},
	}

	ctx := context.Background()
	db := createNewDb(t, "../test_ride.sqlite3", "../assets/ddl.whosdriving-core")
	if db == nil {
		t.Fatal("Couldn't create database connexion")
	}
	defer db.Close()

	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		t.Fatalf("Could't create transaction - %s", err)
	}
	defer tx.Rollback()

	// Create database context
	lCtx := LuwContext{Conn: db, Tx: tx}

	// create & find
	userTest, err := CreateUser(ctx, &lCtx, toNewUser(&expectedCreator))
	assert.Nil(t, err, "")
	assert.Equal(t, &expectedCreator, userTest)

	userJohn, err := CreateUser(ctx, &lCtx, toNewUser(&expectedParticipant1))
	assert.Nil(t, err, "")
	assert.Equal(t, &expectedParticipant1, userJohn)

	rotation, err := CreateRotation(ctx, &lCtx, &newRotation)
	assert.Nil(t, err, "")
	assert.Equal(t, &expectedRotation, rotation)

	ride, err := AddRide(ctx, &lCtx, &newRide)
	assert.Nil(t, err, "")
	assert.Equal(t, &expectedRide, ride)

	updtExpectedRide := expectedRide
	updtExpectedRide.Conductor = &expectedCreator
	updtRide, err := UpdateRide(ctx, &lCtx, &updtExpectedRide)
	assert.Nil(t, err, "")
	assert.Equal(t, &updtExpectedRide, updtRide)

	_, err = DeleteRide(ctx, &lCtx, &updtExpectedRide)
	assert.Nil(t, err, "")
	ride, err = FindRide(ctx, &lCtx, int64(expectedRide.ID))
	assert.Equal(t, err, sql.ErrNoRows)
	assert.Nil(t, ride, "Ride successfuly not found")

	phoenixRide, err := AddRide(ctx, &lCtx, &newRide)
	assert.Nil(t, err, "")
	expectedRide.ID = 2
	assert.Equal(t, &expectedRide, phoenixRide)

	err = RemoveRideParticipants(ctx, &lCtx, int64(expectedRide.ID), &[]string{expectedRide.Conductor.Email})
	assert.Nil(t, err, "")

	expectedRide2 := expectedRide
	expectedRide2.Participants = []*model.User{&expectedCreator}

	ride, err = FindRide(ctx, &lCtx, int64(expectedRide2.ID))
	assert.Nil(t, err, "")
	assert.Equal(t, &expectedRide2, ride)

	err = tx.Commit()
	if err != nil {
		t.Fatalf("Error on commit - %s", err)
	}
}
