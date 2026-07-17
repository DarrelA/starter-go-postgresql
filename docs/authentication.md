# Authentication

The application combines signed JWTs with Redis-backed session records. JWTs carry identity claims, while Redis makes issued tokens revocable.

## Authentication flow

Three mechanisms protect different stages of the flow:

| Mechanism | Responsibility | Used when |
| --- | --- | --- |
| bcrypt | Hash and verify passwords without storing plaintext credentials | Registration and password login |
| JWT | Identify the user and carry signed access or refresh claims | Protected requests, refresh, and logout |
| CSRF token | Prove that this application intentionally initiated a cookie-authenticated mutation | Refresh and logout |

The end-to-end flow is:

1. Registration hashes the submitted password with bcrypt before PostgreSQL stores the user.
2. Login verifies the password with bcrypt. On success, the server creates access and refresh JWTs, records their token family in Redis, and returns them as `HttpOnly` cookies. It also returns a browser-readable CSRF cookie bound to the refresh-token ID.
3. A protected read such as `GET /me` sends the access JWT. The server verifies its signature and expiry, then confirms the access-token record is still active in Redis. Read-only requests do not require a CSRF token.
4. Refresh and logout send the refresh JWT automatically as a cookie. The client must also copy the CSRF cookie into `X-CSRF-Token`, proving that application code deliberately initiated the request.
5. Refresh atomically consumes the current refresh token and replaces the access, refresh, and CSRF cookies. Logout revokes the token family and clears all three cookies. Reuse of a consumed refresh token revokes the family.

## Password login

Registration hashes passwords with bcrypt before storing users in PostgreSQL. Login retrieves a user by email and compares the submitted password with the stored hash.

Successful login creates two RSA-signed JWTs:

- an access token for protected API requests;
- a refresh token for issuing a new access token.

Each JWT contains the user ID, a unique token ID, a refresh-token family ID, and issued/not-before/expiry timestamps. Redis stores the active family, access-token records, refresh-token state, and family membership with TTLs derived from JWT expiry.

The tokens are delivered only through authentication cookies. Successful login and refresh response bodies contain status metadata and do not expose token values to JavaScript.

The access and refresh RSA key pairs are decoded, parsed, and checked for matching public/private keys during startup. Parsed keys are retained by the token service rather than reparsed for each request.

## Access-token validation

Protected endpoints accept an access token from either:

1. `Authorization: Bearer <token>`; or
2. the `access_token` cookie.

Validation has two required stages:

1. Verify the JWT signature and registered time claims.
2. Confirm that its family and token ID identify an active access-token record in Redis.

Redis is therefore part of session validity, not merely an optional cache. A cryptographically valid JWT is rejected after its Redis record expires or is deleted.

## Refresh and logout

`POST /auth/api/v1/users/refresh` validates the refresh JWT and CSRF token, then atomically marks the current refresh token as rotated and issues replacement access, refresh, and CSRF cookies. Refresh tokens are single-use. Redis retains rotated refresh records as replay tombstones until their original expiry.

If a rotated refresh token is presented again, Redis atomically revokes its entire family and removes the family's active access-token records. Concurrent refresh requests therefore cannot both succeed; detecting the losing request as reuse also revokes the replacement session because the server cannot distinguish benign concurrency from token theft.

`POST /auth/api/v1/users/logout` validates the refresh JWT and CSRF token, revokes the token family, and expires the access, refresh, and CSRF cookies. It does not require an unexpired access token, so a client can log out after its access token expires.

## Cookies

Access and refresh cookies use the configured `Domain` and `Secure` values, require `HttpOnly`, and use `SameSite=Strict` and `Path=/`. The `csrf_token` cookie deliberately remains readable by browser code so clients can copy it into the `X-CSRF-Token` request header.

Cookie-authenticated state changes use a signed double-submit CSRF token. Its HMAC binds a random value to the current refresh-token ID. Refresh and logout require the header and cookie values to match and the signature to validate; successful refresh rotates the CSRF token with the refresh token. CORS permits the custom header only from configured origins.

Production configuration should enable `Secure` and use HTTPS; `HttpOnly` is required for authentication cookies. Cookie settings reduce exposure but do not replace a complete browser-security review.

Request logging records route and timing metadata without authentication request or response bodies. Query strings and referrer values are also omitted so passwords, tokens, OAuth authorization codes, and similar credential material are not copied into application logs.

## Google OAuth

The Google login endpoint generates a random state value, stores it in a short-lived `HttpOnly`, `SameSite=Lax` cookie, and redirects to Google. The callback compares state values, exchanges the authorization code with a timeout, and retrieves a bounded user-info response.

Only a profile with a stable Google subject, an email address, and Google's verified-email flag is accepted. The callback then follows this flow:

1. Look up the provider identity by `google` and Google's stable subject identifier. Email is profile data, not the external identity key.
2. Return its existing local user when the identity is already known.
3. Otherwise, link the identity to the local user with the same verified email, or create a passwordless local user when no such user exists.
4. Create the local user and provider identity in one PostgreSQL transaction. A local user can have at most one identity for a given provider; a mismatched subject is reported as a conflict.
5. Issue the same access, refresh, and CSRF cookies used by password login. The response contains the local user and does not expose Google or application tokens.

OAuth-only users have no password hash and cannot use password login unless a future account-management flow explicitly adds a password.
