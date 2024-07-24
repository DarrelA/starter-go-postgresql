package jwt

import (
	"bytes"
	"os"
	"strings"
	"testing"

	uuid "github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestNewTokenService(t *testing.T) {
	for _, test := range tokenTests {
		t.Run(test.name, func(t *testing.T) {
			tokenService := NewTokenService()
			userUUID, err := uuid.NewV7()
			if err != nil {
				t.Fatalf("Failed to create userUUID: %v", err)
			}

			if test.validPrivateKey && test.validPublicKey {
				testToken, err := tokenService.CreateToken(userUUID.String(), testTTL, test.privateKey)
				if err != nil {
					t.Fatalf("Failed to CreateToken: %v", err)
				}

				validatedToken, err := tokenService.ValidateToken(*testToken.Token, test.publicKey)
				if err != nil {
					t.Fatalf("Failed to ValidateToken: %v", err)
				}

				if validatedToken == nil {
					t.Fatalf("Token validation failed")
				}
			}

			var buf bytes.Buffer
			log.Logger = log.Output(&buf)
			zerolog.SetGlobalLevel(zerolog.ErrorLevel)

			if !test.validPrivateKey {
				_, err := tokenService.CreateToken(userUUID.String(), testTTL, test.privateKey)
				if err == nil {
					logOutput := buf.String()
					if !strings.Contains(logOutput, test.expectedErrMsg) {
						t.Errorf("Expected error message '%s' not found in log output: '%s'", test.expectedErrMsg, logOutput)
					}
				}
			}

			if test.validPrivateKey && !test.validPublicKey {
				testToken, err := tokenService.CreateToken(userUUID.String(), testTTL, test.privateKey)
				if err != nil {
					t.Fatalf("Failed to CreateToken: %v", err)
				}

				_, err = tokenService.ValidateToken(*testToken.Token, test.publicKey)
				if err == nil {
					logOutput := buf.String()
					if !strings.Contains(logOutput, test.expectedErrMsg) {
						t.Errorf("Expected error message '%s' not found in log output: '%s'", test.expectedErrMsg, logOutput)
					}
				}
			}

			log.Logger = log.Output(os.Stdout)
		})
	}
}
