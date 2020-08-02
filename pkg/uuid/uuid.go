package uuid

import "github.com/gofrs/uuid"

func RandToken() (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	return id.String(), nil
}
