package renter

import (
	"testing"
	"time"

	"github.com/aarondl/null/v9"
	"github.com/stretchr/testify/require"
)

type individualCase struct {
	email     string
	tenantID  string
	firstName null.String
	lastName  null.String
	time      time.Time
}

func TestNewIndividual(t *testing.T) {
	t.Parallel()

	okValues := individualCase{
		email:     "test@example.com",
		tenantID:  "tenant-123",
		firstName: null.StringFrom("First"),
		lastName:  null.StringFrom("Last"),
		time:      time.Now(),
	}

	tests := map[string]struct {
		values  individualCase
		wantErr bool
	}{
		"ok": {
			values:  okValues,
			wantErr: false,
		},
		"ng (invalid email)": {
			values: individualCase{
				tenantID:  okValues.tenantID,
				email:     "invalid-email",
				firstName: okValues.firstName,
				lastName:  okValues.lastName,
				time:      okValues.time,
			},
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			emailVO, err := NewEmail(tt.values.email)
			if err != nil {
				if tt.wantErr {
					return
				}
				t.Fatalf("failed to create email VO: %v", err)
			}

			individual := NewIndividual(
				"renter-123",
				tt.values.tenantID,
				*emailVO,
				tt.values.firstName,
				tt.values.lastName,
				tt.values.time,
			)

			if tt.wantErr {
				t.Fatal("expected error but got none")
			}

			require.NotNil(t, individual)
			require.NotEmpty(t, individual.RenterID)

			expectedEmail, err := NewEmail(tt.values.email)
			require.NoError(t, err)
			require.True(t, individual.Email.Equals(expectedEmail))
			require.Equal(t, tt.values.tenantID, individual.TenantID)
			require.Equal(t, tt.values.firstName, individual.FirstName)
			require.Equal(t, tt.values.lastName, individual.LastName)
			require.Equal(t, tt.values.time, individual.CreatedAt)
			require.Equal(t, tt.values.time, individual.UpdatedAt)
		})
	}
}
