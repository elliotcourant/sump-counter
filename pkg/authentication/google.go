package authentication

import (
	"context"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"os"
)

const (
	GoogleApplicationCredentialsEnvironment = "GOOGLE_APPLICATION_CREDENTIALS"
	GoogleSheetsScope                       = "https://www.googleapis.com/auth/spreadsheets"
)

func GetGoogleCredentialsFromEnv() (*google.Credentials, error) {
	path := os.Getenv(GoogleApplicationCredentialsEnvironment)
	if len(path) == 0 {
		return nil, errors.Errorf("%s environment variable is not set", GoogleApplicationCredentialsEnvironment)
	}

	json, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read credentials file: %s", path)
	}

	return GetGoogleCredentialsFromJSON(json)
}

func GetGoogleCredentialsFromJSON(json []byte) (*google.Credentials, error) {
	return google.CredentialsFromJSON(context.Background(), json, GoogleSheetsScope)
}
