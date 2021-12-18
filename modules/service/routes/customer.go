package routes

import (
	"github.com/gofiber/fiber/v2"
	"service/controllers"
)

func CustomersRoute(route fiber.Router) {
	route.Get("/", controllers.GetAllCustomers)
	route.Get("/:id", controllers.GetCustomer)
	route.Get("/:id/detail", controllers.GetCustomerDetail)
}
