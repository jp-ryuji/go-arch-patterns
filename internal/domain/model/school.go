package model

import (
	"time"

	"github.com/jp-ryuji/go-sample/internal/pkg/id"
)

type School struct {
	ID        string
	Code      string
	Name      string
	Type      SchoolType
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Schools []*School

func NewSchool(
	code string,
	name string,
	schoolType SchoolType,
	t time.Time,
) *School {
	return &School{
		ID:        id.New(),
		Code:      code,
		Name:      name,
		Type:      schoolType,
		CreatedAt: t,
		UpdatedAt: t,
	}
}

type SchoolType string

type SchoolTypes []SchoolType

const (
	SchoolTypeUnknown       SchoolType = "unknown"
	SchoolTypePrivateSchool SchoolType = "private_school"
	SchoolTypePublicSchool  SchoolType = "public_school"
)

func NewSchoolType(s string) SchoolType {
	switch s {
	case SchoolTypePrivateSchool.String(),
		SchoolTypePublicSchool.String():
		return SchoolType(s)
	}
	return SchoolTypeUnknown
}

func (m SchoolType) String() string {
	return string(m)
}

func (m SchoolType) Valid() bool {
	return m != SchoolTypeUnknown && m != ""
}

func (ms SchoolTypes) Slice() []string {
	dst := make([]string, 0, len(ms))
	for _, schoolType := range ms {
		dst = append(dst, schoolType.String())
	}
	return dst
}

func (ms SchoolTypes) Valid() bool {
	for _, schoolType := range ms {
		if !schoolType.Valid() {
			return false
		}
	}
	return true
}
