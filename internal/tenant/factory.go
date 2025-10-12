package tenant

import (
	"time"

	"github.com/jp-ryuji/go-arch-patterns/internal/pkg/id"
)

func NewFixture() *Tenant {
	return New("test-"+id.New(), time.Now())
}

func NewWithCode(code string) *Tenant {
	return New(code, time.Now())
}
