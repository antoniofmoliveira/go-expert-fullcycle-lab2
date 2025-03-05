package dto

import (
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/antoniofmoliveira/go-expert-fullcycle-lab1/src/internal/shared"
)

type Cep struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
}

// NewCep creates a new Cep instance with the provided details and validates it.
// Returns the created Cep instance if valid, or an error if validation fails.
func NewCep(cep, state, city, neighborhood, street string) (*Cep, error) {
	c := &Cep{
		Cep:          cep,
		State:        state,
		City:         city,
		Neighborhood: neighborhood,
		Street:       street,
	}
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return c, nil
}

// ToJson returns the json representation of the Cep, or an error if the cep is invalid.
func (c *Cep) ToJson() (string, error) {
	if err := c.Validate(); err != nil {
		return "", err
	}
	j, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(j), nil
}

// Validate validates the Cep fields and returns an error if any of them are invalid.
// It checks if the cep is valid, uf is a valid short state name,
// city, neighborhood and street are not empty.
func (c *Cep) Validate() error {
	if _, err := shared.ValidateCep(c.Cep); err != nil {
		return err
	}
	if !shared.ValidateStateShort(c.State) {
		return errors.New("state not found")
	}
	if c.City == "" || c.Neighborhood == "" || c.Street == "" {
		return errors.New("city, neighborhood and street must not be empty")
	}
	return nil
}

// LogValue returns a slog.Value representing the Cep instance.
// It includes fields such as cep, street, neighborhood, city, and state
// in a grouped format for logging purposes.
func (c *Cep) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("cep", c.Cep),
		slog.String("street", c.Street),
		slog.String("neighborhood", c.Neighborhood),
		slog.String("city", c.City),
		slog.String("state", c.State),
	)
}
