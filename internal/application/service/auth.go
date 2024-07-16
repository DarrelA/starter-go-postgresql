package service

import "github.com/gofiber/fiber/v2"

/*
The `AuthService` interface define the contract for authentication-related operations.

Placing the `AuthService` interface in the application layer to keep the domain layer
focused on business logic while the application layer handles orchestration and coordination.

This ensures core business logic remains agnostic to the transport mechanism.
*/
type AuthService interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	RefreshAccessToken(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
}
