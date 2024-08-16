package common

import (
	"fmt"

	"github.com/google/uuid"
)

func GeneratePublicID() string {
	publicId := fmt.Sprintf("cld-%s", uuid.New().String())
	return publicId
}
