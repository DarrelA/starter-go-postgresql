package service

import "github.com/gofiber/fiber/v2"

/*
This interface defines a contract for the `GetUserRecord` method.

By defining an `UserService` interface  in the application layer,
the implementation details is abstracted away from the transport layer.

This allows for easier testing and flexibility in changing implementations without affecting the transport layer.

The interface itself can remain framework agnostic,
ensuring that the core business logic does not depend on specific web frameworks.
*/
type UserService interface {
	GetUserRecord(c *fiber.Ctx) error
}
