// Package http provides the Fiber HTTP adapter for the authentication service.
// It owns routing, transport validation, cookies, status-code mapping, and the
// conversion between Fiber requests and application services.
//
// # Routes
//
// NewRouter mounts these routes below /auth:
//
//	GET  /auth/health
//	POST /auth/api/v1/users/register
//	POST /auth/api/v1/users/login
//	GET  /auth/api/v1/users/me
//	POST /auth/api/v1/users/refresh
//	POST /auth/api/v1/users/logout
//	GET  /auth/google_login
//	GET  /auth/google_callback
//
// Register and Login accept JSON payloads defined in the application/dto
// package. The protected user routes accept an access token from either an
// access_token cookie or an Authorization: Bearer header. Refresh and logout
// require the refresh_token and csrf_token cookies plus the csrf_token value in
// an X-CSRF-Token header. Each successful refresh rotates both JWTs and the
// CSRF token. Reusing a rotated refresh token revokes its token family.
// Authentication middleware validates the JWT and Redis session and passes a
// typed user UUID to downstream handlers. Endpoints such as /me load the full
// PostgreSQL record only when their response requires it.
//
// Successful responses contain a status field with the value "success".
// Authentication tokens are delivered through HttpOnly cookies and are not
// included in login or refresh response bodies. The non-HttpOnly csrf_token
// cookie is intended to be copied into X-CSRF-Token by browser clients. Most
// failures contain status "fail" and an error object with message and status
// fields. Expected errors are mapped at this boundary; unexpected internal
// errors are logged and returned with a generic message.
//
// Google Callback resolves the stable provider identity to a local user,
// atomically creates or links that identity when necessary, and issues the
// same cookie-based application session used by password login.
package http
