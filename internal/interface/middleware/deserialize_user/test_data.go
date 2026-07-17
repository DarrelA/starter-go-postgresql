// coverage:ignore file
// Test file
package middleware

const mockExpiresIn = int64(3600)

var deserializerTests = []struct {
	name           string
	hasError       bool
	expectedErrMsg string
	header         string
	cookieName     string
	cookieValue    string
}{
	{
		name: "Authorization header", hasError: true, expectedErrMsg: "please login again",
		header: "Bearer mockInvalidBearerToken",
	},
	{
		name: "No authorization header and empty cookie value", hasError: true, expectedErrMsg: "please login again",
		header: "", cookieName: "access_token", cookieValue: "",
	},
	{
		name: "No authorization header and no cookie", hasError: true, expectedErrMsg: "please login again",
		header: "", cookieName: "", cookieValue: "",
	},
	{
		name: "Authorization cookie", hasError: false,
		cookieName: "access_token", cookieValue: "mockAccessToken",
	},
}
