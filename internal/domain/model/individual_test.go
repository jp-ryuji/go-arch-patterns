package model

import (
	"testing"
	"time"

	"github.com/aarondl/null/v9"
	"github.com/jp-ryuji/go-sample/internal/domain/model/value"
	"github.com/stretchr/testify/require"
)

// okCaseValues holds the test values for the valid individual case
type okCaseValues struct {
	email     string
	tenantID  string
	firstName null.String
	lastName  null.String
	time      time.Time
}

func TestNewIndividual(t *testing.T) {
	t.Parallel()

	// Define test values for the ok case
	okValues := okCaseValues{
		email:     "test@example.com",
		tenantID:  "tenant-123",
		firstName: null.StringFrom("First"),
		lastName:  null.StringFrom("Last"),
		time:      time.Now(),
	}

	tests := map[string]struct {
		values  okCaseValues
		wantErr bool
	}{
		"ok": {
			values:  okValues,
			wantErr: false,
		},
		"ng (invalid email)": {
			values: okCaseValues{
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

			individual, err := NewIndividual(
				tt.values.tenantID,
				tt.values.email,
				tt.values.firstName,
				tt.values.lastName,
				tt.values.time,
			)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, individual)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, individual)
			require.NotEmpty(t, individual.ID)

			// Create an email VO to compare with the individual's email
			expectedEmail, err := value.NewEmail(tt.values.email)
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
