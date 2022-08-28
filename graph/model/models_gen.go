// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type NewRide struct {
	IDRotation        string   `json:"idRotation"`
	EmailConductor    string   `json:"emailConductor"`
	EmailParticipants []string `json:"emailParticipants"`
}

type NewRole struct {
	Email string `json:"email"`
	Role  Role   `json:"role"`
}

type NewRotation struct {
	Name              string   `json:"name"`
	EmailCreator      string   `json:"emailCreator"`
	EmailParticipants []string `json:"emailParticipants"`
}

type NewUser struct {
	Email   string  `json:"email"`
	Name    string  `json:"name"`
	Profile *string `json:"profile"`
}

type Ride struct {
	ID           string  `json:"id"`
	Conductor    *User   `json:"conductor"`
	Participants []*User `json:"participants"`
}

type Rotation struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Creator      *User   `json:"creator"`
	Participants []*User `json:"participants"`
	Rides        []*Ride `json:"rides"`
}

type User struct {
	Email   string  `json:"email"`
	Name    string  `json:"name"`
	Profile *string `json:"profile"`
	Role    Role    `json:"role"`
}

type Role string

const (
	RoleAdmin       Role = "ADMIN"
	RoleStandard    Role = "STANDARD"
	RoleUnregistred Role = "UNREGISTRED"
)

var AllRole = []Role{
	RoleAdmin,
	RoleStandard,
	RoleUnregistred,
}

func (e Role) IsValid() bool {
	switch e {
	case RoleAdmin, RoleStandard, RoleUnregistred:
		return true
	}
	return false
}

func (e Role) String() string {
	return string(e)
}

func (e *Role) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Role(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Role", str)
	}
	return nil
}

func (e Role) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
