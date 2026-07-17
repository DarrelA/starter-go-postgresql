// Package http provides the Google OAuth HTTP handlers used by the parent
// transport package. Login starts the authorization flow and Callback validates
// state, exchanges the authorization code, resolves the provider identity to a
// local user, and establishes a cookie-based application session.
package http
