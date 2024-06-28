package data_test

import "github.com/DarrelA/starter-go-postgresql/internal/domains/users"

var RegisterInputs = []users.RegisterInput{
	{FirstName: "kang yao", LastName: "tan", Email: "tky@e.com", Password: "I_love_c_sharp"},
	{FirstName: "jie wei", LastName: "low", Email: "ljw@e.com", Password: "i_love_java"},
	{FirstName: "bing hong", LastName: "tan", Email: "tbh@e.com", Password: "i_heArt_VB"},
	{FirstName: "jason", LastName: "the consultant", Email: "j@e.com", Password: "I&asked(a)question^at!the*town-hall:_Why@is9the6air~conditioning%so+cold"},

	// all empty fields
	{FirstName: "", LastName: "", Email: "", Password: ""},
	// empty password only
	{FirstName: "John", LastName: "Doe", Email: "John_Doe@gmail.com", Password: ""},
	// empty email only
	{FirstName: "Jane", LastName: "Smith", Email: "", Password: "Password1!"},
	// empty last name only
	{FirstName: "Alice", LastName: "", Email: "Alice@yahoo.com", Password: "Password1!"},
	// empty first name only
	{FirstName: "", LastName: "Brown", Email: "Brown@outlook.com", Password: "Password1!"},
	// valid input
	{FirstName: "Emily", LastName: "Clark", Email: "Emily_Clark@gmail.com", Password: "Password1!"},
	// email is already taken
	{FirstName: "emily", LastName: "clark", Email: "emily_clark@gmail.com", Password: "Password1!"},
	// invalid email format
	{FirstName: "Oliver", LastName: "Jones", Email: "Oliver_Jones", Password: "Password1!"},
	// invalid password (too short)
	{FirstName: "Michael", LastName: "Taylor", Email: "Michael_Taylor@yahoo.com", Password: "Pass1!"},
	// invalid password (no special character)
	{FirstName: "Emma", LastName: "Davis", Email: "Emma_Davis@outlook.com", Password: "Password1"},
	// invalid password (no number)
	{FirstName: "William", LastName: "Martinez", Email: "William_Martinez@gmail.com", Password: "Password!"},
	// invalid first name (less than 2 characters)
	{FirstName: "A", LastName: "Anderson", Email: "A_Anderson@yahoo.com", Password: "Password1!"},
	// invalid last name (less than 2 characters)
	{FirstName: "Sophia", LastName: "B", Email: "Sophia_B@outlook.com", Password: "Password1!"},
	// invalid first name (more than 50 characters)
	{FirstName: "Taumatawhakatangihangakoauauotamateaturipukakapikimaungahoronukupokaiwhenu" + "12345678901234567890", LastName: "Harris", Email: "LongFirstName@gmail.com", Password: "Password1!"},
	// invalid last name (more than 50 characters)
	{FirstName: "James", LastName: "Taumatawhakatangihangakoauauotamateaturipukakapikimaungahoronukupokaiwhenu" + "12345678901234567890", Email: "LongLastName@gmail.com", Password: "Password1!"},
	// invalid first name (contains non-alphabetic characters)
	{FirstName: "John123", LastName: "Walker", Email: "John123_Walker@gmail.com", Password: "Password1!"},
	// invalid last name (contains non-alphabetic characters)
	{FirstName: "Grace", LastName: "Miller456", Email: "Grace_Miller456@yahoo.com", Password: "Password1!"},
}
