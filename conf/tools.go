package conf

import (
	"github.com/google/uuid"
)

func getRequiredString() string {
	u, _ := uuid.NewV7()
	return "$" + u.String()
}
