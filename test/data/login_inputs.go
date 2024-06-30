package data_test

import "github.com/DarrelA/starter-go-postgresql/internal/domains/users"

var LoginInputs = []users.LoginInput{

	// valid account
	{Email: "Carlyn_Daniel@gmail.com", Password: "Password1!"},
	// incorrect email
	{Email: "Carlyn_Daniël@gmail.com", Password: "Password1!"},
	// valid account
	{Email: "Emily_Clark@gmail.com", Password: "Password1!"},
	// incorrect email
	{Email: "emily_clarky@gmail.com", Password: "Password1!"},
	// incorrect password
	{Email: "emily_clark@gmail.com", Password: "password1!"},
	// empty email only
	{Email: "", Password: "Password1!"},
	// empty password only
	{Email: "Emily_Clark@gmail.com", Password: ""},
	// empty email and password
	{Email: "", Password: ""},
	// invalid email format
	{Email: "Emily_Clark@com", Password: "Password1!"},
	// SQL injection attempt
	{Email: "Emily_Clark@gmail.com", Password: "' OR '1'='1"},
	// email case sensitivity
	{Email: "EMILY_CLARK@gmail.com", Password: "Password1!"},
	// password case sensitivity
	{Email: "Emily_Clark@gmail.com", Password: "password1!"},
	// email with trailing spaces
	{Email: "Emily_Clark@gmail.com    ", Password: "Password1!"},
	// password with leading/trailing spaces
	{Email: "   Emily_Clark@gmail.com", Password: " Password1! "},
	// long email input
	{Email: "this.is.a.very.long.email.address@example.com", Password: "Password1!"},
	// long password input
	{Email: "Emily_Clark@gmail.com", Password: "Taumatawhakatangihangakoauauotamateaturipukakapikimaungahoronukupokaiwhenu!"},
	// invalid characters in email
	{Email: "Emily_Clark@!gmail.com", Password: "Password1!"},
}
