//go:build dev
// +build dev

package mock

import (
	"github.com/august-kuhfuss/hotdamn/store"
	"github.com/jaswdr/faker/v2"
)

type mock struct {
	F faker.Faker
}

func NewStore() store.Store {
	return &mock{F: faker.New()}
}
