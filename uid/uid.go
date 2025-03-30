package uid

import (
	"github.com/google/uuid"
	"github.com/inovacc/ksuid"
)

func GenerateUUID() string {
	return uuid.NewString()
}

func GenerateKSUID() string {
	return ksuid.NewString()
}
