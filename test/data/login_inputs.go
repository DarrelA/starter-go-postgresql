package data_test

import (
	"net/http"

	"github.com/DarrelA/starter-go-postgresql/internal/domains/users"
)

type TestLoginInput struct {
	TestName           string
	ExpectedStatusCode int
	users.LoginInput
}

var LoginInputs = []TestLoginInput{
	{TestName: "valid account", ExpectedStatusCode: http.StatusOK,
		LoginInput: users.LoginInput{Email: "Carlyn_Daniel@gmail.com", Password: "Password1!"}},
	{TestName: "incorrect email", ExpectedStatusCode: http.StatusBadRequest,
		LoginInput: users.LoginInput{Email: "Carlyn_Daniël@gmail.com", Password: "Password1!"}},
	{TestName: "valid account", ExpectedStatusCode: http.StatusBadRequest,
		LoginInput: users.LoginInput{Email: "Emily_Clark@gmail.com", Password: "Password1!"}},
	{TestName: "incorrect email", ExpectedStatusCode: http.StatusBadRequest,
		LoginInput: users.LoginInput{Email: "emily_clarky@gmail.com", Password: "Password1!"}},
	{TestName: "incorrect password", ExpectedStatusCode: http.StatusBadRequest,
		LoginInput: users.LoginInput{Email: "Emily_Clark@gmail.com", Password: "password1!2"}},
	{TestName: "empty email only", ExpectedStatusCode: http.StatusBadRequest,
		LoginInput: users.LoginInput{Email: "", Password: "Password1!"}},
	{TestName: "empty password only", ExpectedStatusCode: http.StatusBadRequest,
		LoginInput: users.LoginInput{Email: "Emily_Clark@gmail.com", Password: ""}},
	{TestName: "empty email and password", ExpectedStatusCode: http.StatusBadRequest,
		LoginInput: users.LoginInput{Email: "", Password: ""}},
	{TestName: "invalid email format", ExpectedStatusCode: http.StatusBadRequest,
		LoginInput: users.LoginInput{Email: "Emily_Clark@com", Password: "Password1!"}},
	{TestName: "SQL injection attempt", ExpectedStatusCode: http.StatusBadRequest,
		LoginInput: users.LoginInput{Email: "Emily_Clark@gmail.com", Password: "' OR '1'='1"}},
	{TestName: "email case sensitivity", ExpectedStatusCode: http.StatusBadRequest,
		LoginInput: users.LoginInput{Email: "EMILY_CLARK@gmail.com", Password: "Password1!"}},
	{TestName: "password case sensitivity", ExpectedStatusCode: http.StatusBadRequest,
		LoginInput: users.LoginInput{Email: "Emily_Clark@gmail.com", Password: "password1!"}},
	{TestName: "email with trailing spaces", ExpectedStatusCode: http.StatusBadRequest,
		LoginInput: users.LoginInput{Email: "Emily_Clark@gmail.com    ", Password: "Password1!"}},
	{TestName: "password with leading spaces", ExpectedStatusCode: http.StatusBadRequest,
		LoginInput: users.LoginInput{Email: "   Emily_Clark@gmail.com", Password: " Password1! "}},
	{TestName: "long email input", ExpectedStatusCode: http.StatusBadRequest,
		LoginInput: users.LoginInput{Email: "this.is.a.very.long.email.address.Taumatawhakatangihangakoauauotamateaturipukakapikimaungahoronukupokaiwhenu@gmail.com", Password: "Password1!"}},
	{TestName: "long password input", ExpectedStatusCode: http.StatusBadRequest,
		LoginInput: users.LoginInput{Email: "Emily_Clark@gmail.com", Password: "Taumatawhakatangihangakoauauotamateaturipukakapikimaungahoronukupokaiwhenu!"}},
	{TestName: "invalid email format", ExpectedStatusCode: http.StatusBadRequest,
		LoginInput: users.LoginInput{Email: "Emily_Clark@!gmail.com", Password: "Password1!"}},
}
