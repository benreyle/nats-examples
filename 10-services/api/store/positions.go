package store

import (
	"fmt"
	"nats-examples/10-services/helpers"

	"github.com/google/uuid"
)

var positions = make([]*helpers.Position, 0)

func Save(position *helpers.Position) {
	positions = append(positions, position)
}

func List() []*helpers.Position {
	return positions
}

func Get(id uuid.UUID) (*helpers.Position, error) {
	for _, pos := range positions {
		if pos.ID == id {
			return pos, nil
		}
	}

	return nil, fmt.Errorf("%s not found", id.String())
}

func Delete(id uuid.UUID) error {
	for idx, pos := range positions {
		if pos.ID == id {
			positions[idx].Deleted = true
			return nil
		}
	}

	return fmt.Errorf("%s not found", id.String())
}
