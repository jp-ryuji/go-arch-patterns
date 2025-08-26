package model

import (
	"testing"
	"time"

	"github.com/jp-ryuji/go-sample/internal/domain/model/value"
	"github.com/stretchr/testify/require"
)

// okCaseValues holds the test values for the valid student case
type okCaseValues struct {
	authUID          string
	email            string
	displayName      string
	firstName        string
	lastName         string
	firstNameKana    string
	lastNameKana     string
	profileImagePath string
	time             time.Time
}

func TestNewStudent(t *testing.T) {
	t.Parallel()

	// Define test values for the ok case
	okValues := okCaseValues{
		authUID:          "auth-uid",
		email:            "test@example.com",
		displayName:      "Display Name",
		firstName:        "First",
		lastName:         "Last",
		firstNameKana:    "FirstKana",
		lastNameKana:     "LastKana",
		profileImagePath: "profile.jpg",
		time:             time.Now(),
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
				authUID:          okValues.authUID,
				email:            "invalid-email",
				displayName:      okValues.displayName,
				firstName:        okValues.firstName,
				lastName:         okValues.lastName,
				firstNameKana:    okValues.firstNameKana,
				lastNameKana:     okValues.lastNameKana,
				profileImagePath: okValues.profileImagePath,
				time:             okValues.time,
			},
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			student, err := NewStudent(
				tt.values.authUID,
				tt.values.email,
				tt.values.displayName,
				tt.values.firstName,
				tt.values.lastName,
				tt.values.firstNameKana,
				tt.values.lastNameKana,
				tt.values.profileImagePath,
				tt.values.time,
			)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, student)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, student)
			require.NotEmpty(t, student.ID)
			require.Equal(t, tt.values.authUID, student.AuthUID)

			// Create an email VO to compare with the student's email
			expectedEmail, err := value.NewEmail(tt.values.email)
			require.NoError(t, err)
			require.True(t, student.Email.Equals(expectedEmail))

			require.Equal(t, tt.values.displayName, student.DisplayName)
			require.Equal(t, tt.values.firstName, student.FirstName)
			require.Equal(t, tt.values.lastName, student.LastName)
			require.Equal(t, tt.values.firstNameKana, student.FirstNameKana)
			require.Equal(t, tt.values.lastNameKana, student.LastNameKana)
			require.Equal(t, tt.values.profileImagePath, student.ProfileImagePath)
			require.Equal(t, tt.values.time, student.CreatedAt)
			require.Equal(t, tt.values.time, student.UpdatedAt)
		})
	}
}
