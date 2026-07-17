package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const requestIdentityKey = "requestIdentity"

// RequestIdentity is the authenticated identity verified for the current request.
// It deliberately contains no persistence or profile data.
type RequestIdentity struct {
	UserUUID uuid.UUID
}

// RequestIdentityFrom returns the identity established by Authenticate.
func RequestIdentityFrom(c *fiber.Ctx) (RequestIdentity, bool) {
	identity, ok := c.Locals(requestIdentityKey).(RequestIdentity)
	return identity, ok
}

// SetRequestIdentity stores a verified identity for downstream handlers.
func SetRequestIdentity(c *fiber.Ctx, identity RequestIdentity) {
	c.Locals(requestIdentityKey, identity)
}
