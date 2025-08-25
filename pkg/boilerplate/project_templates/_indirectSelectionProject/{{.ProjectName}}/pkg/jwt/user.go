package jwt

import (
	"encoding/json"
	"github.com/pkg/errors"
	"os"
)

// CliUsers represents a group of authorized users that authenticate via CLI.
type CliUsers struct {
	Users   []*CliUser          `json:"users"`
	UserMap map[string]*CliUser `json:"user_map"`
}

// CliUser represents a user who authenticates via CLI.
type CliUser struct {
	Name    string   `json:"name"`
	PubKeys []string `json:"public_keys"`
	Role    string   `json:"role"`
}

// LoadCliUsersFromFile loads CLI users from a JSON file.
func LoadCliUsersFromFile(filepath string) (users *CliUsers, err error) {
	userBytes, err := os.ReadFile(filepath)
	if err != nil {
		err = errors.Wrapf(err, "failed loading file %s", filepath)
		return users, err
	}

	users, err = LoadCliUsersFromBytes(userBytes)
	if err != nil {
		err = errors.Wrapf(err, "failed loading Users from data in %s", filepath)
	}

	return users, err
}

// LoadCliUsersFromBytes loads CLI users from JSON bytes.
func LoadCliUsersFromBytes(userBytes []byte) (users *CliUsers, err error) {
	users = &CliUsers{}

	err = json.Unmarshal(userBytes, users)
	if err != nil {
		err = errors.Wrapf(err, "failed unmarshalling CliUsers")
		return users, err
	}

	return users, err
}