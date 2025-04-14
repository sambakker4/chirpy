package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJwt(t *testing.T) {
	cases := []struct {
		TokenSecret string
		ID          uuid.UUID
	} {
		{
			TokenSecret: "aasdfasdf312413241!#$$%!%$",
			ID: uuid.MustParse("55555555-5555-5555-5555-555555555555"),
		},
		{
			TokenSecret: "",
			ID: uuid.MustParse("44444444-4444-4444-4444-444444444444"),
		},
		{
			TokenSecret: "!@%$#!%@!#$#@!DSFAS#@!$!@#$!#@$!@a",
			ID: uuid.MustParse("33333333-3333-3333-3333-333333333333"),
		},
	}

	for _, testcase := range cases {
		token, err := MakeJWT(
			testcase.ID, testcase.TokenSecret, time.Second * 15,
		)		

		if err != nil {
			t.Errorf("Error creating token: %v", err)
		}

		userID, err := ValidateJWT(token, testcase.TokenSecret)
		if err != nil {
			t.Errorf("Error creating token: %v", err)
		}

		if userID != testcase.ID {
			t.Errorf("Returned validated id: %s, and original id: %s, not equal", 
				userID.String(),
				testcase.ID.String(),
			)
		}
	}

	// test expiring
	token, err := MakeJWT(uuid.New(), "asdfadsf", time.Second * -15)		
	if err != nil {
		t.Errorf("Errof creating token: %v", err)
	}

	_, err = ValidateJWT(token, "asdfadsf")
	if err == nil {
		t.Errorf("Token not timing out")
	}
}
