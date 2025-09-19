package value

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewEmail(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		args    string
		want    *Email
		wantErr string
	}{
		"ok (valid email)": {
			args:    "test@example.com",
			want:    &Email{value: "test@example.com"},
			wantErr: "",
		},
		"ok (valid email with plus)": {
			args:    "test+label@example.com",
			want:    &Email{value: "test+label@example.com"},
			wantErr: "",
		},
		"ng (invalid email without @)": {
			args:    "testexample.com",
			want:    nil,
			wantErr: "invalid email format",
		},
		"ng (invalid email without domain)": {
			args:    "test@",
			want:    nil,
			wantErr: "invalid email format",
		},
		"ng (empty email)": {
			args:    "",
			want:    nil,
			wantErr: "email cannot be empty",
		},
		"ng (email too long)": {
			args:    "aVeryLongEmailThatExceedsTheMaximumAllowedLengthOfFiftyCharacters@example.com",
			want:    nil,
			wantErr: "email exceeds maximum length",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, err := NewEmail(tt.args)

			if tt.wantErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErr)
				require.Nil(t, got)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want.value, got.value)
		})
	}
}

func TestEmailString(t *testing.T) {
	t.Parallel()

	emailStr := "test@example.com"
	email, err := NewEmail(emailStr)
	require.NoError(t, err)

	if email.String() != emailStr {
		t.Errorf("Email.String() = %v, want %v", email.String(), emailStr)
	}
}

func TestEmailEquals(t *testing.T) {
	t.Parallel()

	email1, err := NewEmail("test@example.com")
	require.NoError(t, err)

	email2, err := NewEmail("test@example.com")
	require.NoError(t, err)

	email3, err := NewEmail("other@example.com")
	require.NoError(t, err)

	if !email1.Equals(email2) {
		t.Error("Expected email1 to equal email2")
	}

	if email1.Equals(email3) {
		t.Error("Expected email1 to not equal email3")
	}

	if email1.Equals(nil) {
		t.Error("Expected email1 to not equal nil")
	}
}
