package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr ValidationErrors
	}{
		{
			in: User{
				ID:     "id",
				Name:   "name",
				Age:    11,
				Email:  "mail@mail.ru",
				Role:   "admin",
				Phones: []string{"12345678900", "12345678900"},
			},
			expectedErr: ValidationErrors{
				{
					Field: "ID",
					Err:   ErrStrNotEqualLen,
				},
				{
					Field: "Age",
					Err:   ErrIntMin,
				},
			},
		},
		{
			in: User{
				ID:     "qweqweqweqweqweqwe123123112312312312",
				Name:   "name",
				Age:    88,
				Email:  "mail@mail.ru",
				Role:   "test",
				Phones: []string{},
			},
			expectedErr: ValidationErrors{
				{
					Field: "Age",
					Err:   ErrIntMax,
				},
				{
					Field: "Role",
					Err:   ErrStrNotInArr,
				},
			},
		},
		{
			in: User{
				ID:     "id",
				Name:   "name1",
				Age:    18,
				Email:  "email",
				Role:   "admin",
				Phones: []string{"12345678900", "12345678900"},
			},
			expectedErr: ValidationErrors{
				{
					Field: "ID",
					Err:   ErrStrNotEqualLen,
				},
				{
					Field: "Email",
					Err:   ErrStrRegMatch,
				},
			},
		},
		{
			in: User{
				ID:     "id",
				Name:   "name",
				Age:    18,
				Email:  "mail@mail.ru",
				Role:   "admin",
				Phones: []string{"11", "12345678900"},
			},
			expectedErr: ValidationErrors{
				{
					Field: "ID",
					Err:   ErrStrNotEqualLen,
				},
				{
					Field: "Phones",
					Err:   ErrStrNotEqualLen,
				},
			},
		},
		{
			in: User{
				ID:     "qweqweqweqweqweqwe123123112312312312",
				Name:   "name",
				Age:    18,
				Email:  "mail@mail.ru",
				Role:   "admin",
				Phones: []string{"12345678900"},
			},
			expectedErr: ValidationErrors{},
		},
		{
			in: App{
				Version: "123",
			},
			expectedErr: ValidationErrors{
				{
					Field: "Version",
					Err:   ErrStrNotEqualLen,
				},
			},
		},
		{
			in: Response{
				Code: 404,
				Body: "bodybodybody",
			},
			expectedErr: ValidationErrors{},
		},
		{
			in: Response{
				Code: 502,
				Body: "bodybodybody",
			},
			expectedErr: ValidationErrors{
				{
					Field: "Code",
					Err:   ErrIntNotInArr,
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			validationErr, err := Validate(tt.in)

			require.NoError(t, err)
			require.Equal(t, len(validationErr), len(tt.expectedErr), "different length of errors")

			for i, v := range validationErr {
				require.Equal(t, v.Field, tt.expectedErr[i].Field, "wrong validation error field")
				require.Truef(t, errors.Is(v.Err, tt.expectedErr[i].Err), "actual error %q", err)
			}

			_ = tt
		})
	}
}
