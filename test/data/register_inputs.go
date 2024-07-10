package data_test

import (
	"net/http"

	user "github.com/DarrelA/starter-go-postgresql/internal/domain/entity"
)

type TestRegisterInput struct {
	TestName           string
	ExpectedStatusCode int
	user.RegisterInput
}

// @NOTES: FirstName and LastName does not accept space
var RegisterInputs = []TestRegisterInput{
	{TestName: "valid account 1", ExpectedStatusCode: http.StatusOK, RegisterInput: user.RegisterInput{
		FirstName: "kangYao", LastName: "tan", Email: "tky@e.com", Password: "I_<3_c_#"}},
	{TestName: "valid account 2", ExpectedStatusCode: http.StatusOK, RegisterInput: user.RegisterInput{
		FirstName: "jieWei", LastName: "low", Email: "ljw@e.com", Password: "i_<3_Java"}},
	{TestName: "valid account 3", ExpectedStatusCode: http.StatusOK, RegisterInput: user.RegisterInput{
		FirstName: "bingHong", LastName: "tan", Email: "tbh@e.com", Password: "1_heArt_VB:)"}},
	{TestName: "valid account 4", ExpectedStatusCode: http.StatusBadRequest, RegisterInput: user.RegisterInput{
		FirstName: "jason", LastName: "the consultant", Email: "j@e.com", Password: "I&asked(a)question^at!the*town-hall:_Why@is9the6air~conditioning%so+cold"}},

	{TestName: "valid account 5 (leading and trailing whitespace)", ExpectedStatusCode: http.StatusOK,
		RegisterInput: user.RegisterInput{FirstName: "  Dessislava  ", LastName: "  Kenyatta  ",
			Email: "  Dessislava.Kenyatta@outlook.com  ", Password: "     Password1!   "}},

	{TestName: "all empty fields", ExpectedStatusCode: http.StatusBadRequest, RegisterInput: user.RegisterInput{
		FirstName: "", LastName: "", Email: "", Password: ""}},
	{TestName: "empty password only", ExpectedStatusCode: http.StatusBadRequest, RegisterInput: user.RegisterInput{
		FirstName: "John", LastName: "Doe", Email: "John_Doe@gmail.com", Password: ""}},
	{TestName: "empty email only", ExpectedStatusCode: http.StatusBadRequest, RegisterInput: user.RegisterInput{
		FirstName: "Jane", LastName: "Smith", Email: "", Password: "Password1!"}},
	{TestName: "empty last name only", ExpectedStatusCode: http.StatusBadRequest, RegisterInput: user.RegisterInput{
		FirstName: "Alice", LastName: "", Email: "Alice@yahoo.com", Password: "Password1!"}},
	{TestName: "empty first name only", ExpectedStatusCode: http.StatusBadRequest, RegisterInput: user.RegisterInput{
		FirstName: "", LastName: "Brown", Email: "Brown@outlook.com", Password: "Password1!"}},
	{TestName: "email is already taken 1", ExpectedStatusCode: http.StatusBadRequest, RegisterInput: user.RegisterInput{
		FirstName: "Emily", LastName: "Clark", Email: "Emily_Clark@gmail.com", Password: "Password1!"}},
	{TestName: "email is already taken 2", ExpectedStatusCode: http.StatusBadRequest, RegisterInput: user.RegisterInput{
		FirstName: "emily", LastName: "clark", Email: "emily_clark@gmail.com", Password: "Password1!"}},
	{TestName: "invalid email (less than 5 characters)", ExpectedStatusCode: http.StatusBadRequest,
		RegisterInput: user.RegisterInput{FirstName: "Jamie", LastName: "Tuna", Email: "@.me", Password: "Password1!"}},

	{TestName: "invalid email (more than 64 characters)", ExpectedStatusCode: http.StatusBadRequest,
		RegisterInput: user.RegisterInput{FirstName: "Jasmine", LastName: "Worth", Email: "Taumatawhakatangihangakoauauotamateaturipukakapikimaungahoronukupokaiwhenu@gmail.com", Password: "Password1!"}},

	{TestName: "invalid email format", ExpectedStatusCode: http.StatusBadRequest, RegisterInput: user.RegisterInput{
		FirstName: "Oliver", LastName: "Jones", Email: "Oliver_Jones", Password: "Password1!"}},

	{TestName: "invalid password (too short)", ExpectedStatusCode: http.StatusBadRequest, RegisterInput: user.RegisterInput{FirstName: "Michael", LastName: "Taylor", Email: "Michael_Taylor@yahoo.com", Password: "Pass1!"}},

	{TestName: "invalid password (no special character)", ExpectedStatusCode: http.StatusBadRequest,
		RegisterInput: user.RegisterInput{FirstName: "Emma", LastName: "Davis",
			Email: "Emma_Davis@outlook.com", Password: "Password1"}},

	{TestName: "invalid password (no number)", ExpectedStatusCode: http.StatusBadRequest,
		RegisterInput: user.RegisterInput{FirstName: "William", LastName: "Martinez", Email: "William_Martinez@gmail.com", Password: "Password!"}},

	{TestName: "invalid first name (less than 2 characters)", ExpectedStatusCode: http.StatusBadRequest,
		RegisterInput: user.RegisterInput{FirstName: "A", LastName: "Anderson", Email: "A_Anderson@yahoo.com", Password: "Password1!"}},

	{TestName: "invalid last name (less than 2 characters)", ExpectedStatusCode: http.StatusBadRequest,
		RegisterInput: user.RegisterInput{FirstName: "Sophia", LastName: "B", Email: "Sophia_B@outlook.com", Password: "Password1!"}},

	{TestName: "invalid first name (more than 50 characters)", ExpectedStatusCode: http.StatusBadRequest,
		RegisterInput: user.RegisterInput{
			FirstName: "Taumatawhakatangihangakoauauotamateaturipukakapikimaungahoronukupokaiwhenu", LastName: "Harris", Email: "LongFirstName@gmail.com", Password: "Password1!"}},

	{TestName: "invalid last name (more than 50 characters)", ExpectedStatusCode: http.StatusBadRequest,
		RegisterInput: user.RegisterInput{
			FirstName: "James", LastName: "Taumatawhakatangihangakoauauotamateaturipukakapikimaungahoronukupokaiwhenu", Email: "LongLastName@gmail.com", Password: "Password1!"}},

	{TestName: "invalid first name (contains non-alphabetic characters)", ExpectedStatusCode: http.StatusBadRequest,
		RegisterInput: user.RegisterInput{FirstName: "John123", LastName: "Walker",
			Email: "John123_Walker@gmail.com", Password: "Password1!"}},

	{TestName: "invalid last name (contains non-alphabetic characters)", ExpectedStatusCode: http.StatusBadRequest,
		RegisterInput: user.RegisterInput{FirstName: "Grace", LastName: "Miller456",
			Email: "Grace_Miller456@yahoo.com", Password: "Password1!"}},

	{TestName: "whitespace only fields", ExpectedStatusCode: http.StatusBadRequest,
		RegisterInput: user.RegisterInput{FirstName: "   ", LastName: "   ", Email: "   ", Password: "   "}},
}
