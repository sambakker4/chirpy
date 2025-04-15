package auth

import (
	"testing"
	"time"
	"net/http"

	"github.com/google/uuid"
)

func TestJwt(t *testing.T) {
	cases := []struct {
		TokenSecret string
		ID          uuid.UUID
	}{
		{
			TokenSecret: "aasdfasdf312413241!#$$%!%$",
			ID:          uuid.MustParse("55555555-5555-5555-5555-555555555555"),
		},
		{
			TokenSecret: "",
			ID:          uuid.MustParse("44444444-4444-4444-4444-444444444444"),
		},
		{
			TokenSecret: "!@%$#!%@!#$#@!DSFAS#@!$!@#$!#@$!@a",
			ID:          uuid.MustParse("33333333-3333-3333-3333-333333333333"),
		},
	}

	for _, testcase := range cases {
		token, err := MakeJWT(
			testcase.ID, testcase.TokenSecret, time.Second*15,
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
	token, err := MakeJWT(uuid.New(), "asdfadsf", time.Second*-15)
	if err != nil {
		t.Errorf("Errof creating token: %v", err)
	}

	_, err = ValidateJWT(token, "asdfadsf")
	if err == nil {
		t.Errorf("Token not timing out")
	}
}

func TestGetBearerToken(t *testing.T) {
	cases := []struct {
		headerName string
		authInfo   string
		result     string
	}{
		{
			headerName: "Authorization",
			authInfo:   "Bearer adsfasdf",
			result:     "adsfasdf",
		},
		{
			headerName: "NotRight",
			authInfo:   "Bearer asdf",
			result:     "error",
		},
		{
			headerName: "Authorization",
			authInfo:   "Beareradfa",
			result:     "error",
		},
	}

	for _, testcase := range cases {
		header := http.Header{}
		header.Set(testcase.headerName, testcase.authInfo)

		if testcase.result == "error" {
			_, err := GetBearerToken(header)	
			if err == nil {
				t.Errorf("Expect error but got nil")
			}
		} else {
			token, err := GetBearerToken(header)
			if err != nil {
				t.Errorf("Error: %v", err)
			}
			
			if token != testcase.result {
				t.Errorf("Expected: %v, and got: %v", testcase.result, token)
			}
		}
	}
}
