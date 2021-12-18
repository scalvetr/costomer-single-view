package controllers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"math"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"service/config"
	"service/models"
)

func GetAllCustomers(c *fiber.Ctx) error {
	log.Printf("GetAllCustomers\n")
	customerCollection := config.MI.DB.Collection("customers")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var customers []models.Customer

	filter := bson.M{}
	findOptions := options.Find()

	if s := c.Query("s"); s != "" {
		filter = bson.M{
			"$or": []bson.M{
				{
					"movieName": bson.M{
						"$regex": primitive.Regex{
							Pattern: s,
							Options: "i",
						},
					},
				},
				{
					"customer": bson.M{
						"$regex": primitive.Regex{
							Pattern: s,
							Options: "i",
						},
					},
				},
			},
		}
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limitVal, _ := strconv.Atoi(c.Query("limit", "10"))
	var limit int64 = int64(limitVal)

	total, _ := customerCollection.CountDocuments(ctx, filter)

	findOptions.SetSkip((int64(page) - 1) * limit)
	findOptions.SetLimit(limit)

	cursor, err := customerCollection.Find(ctx, filter, findOptions)
	defer cursor.Close(ctx)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Customers Not found",
			"error":   err,
		})
	}

	for cursor.Next(ctx) {
		var customer models.Customer
		cursor.Decode(&customer)
		customers = append(customers, customer)
	}

	last := math.Ceil(float64(total / limit))
	if last < 1 && total > 0 {
		last = 1
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":      customers,
		"total":     total,
		"page":      page,
		"last_page": last,
		"limit":     limit,
	})
}

func GetCustomerDetail(c *fiber.Ctx) error {
	log.Printf("GetCustomerDetail\n")
	return getCustomer(c, true)
}
func GetCustomer(c *fiber.Ctx) error {
	log.Printf("GetCustomer\n")
	return getCustomer(c, false)
}

func getCustomer(c *fiber.Ctx, detail bool) error {
	customerCollection := config.MI.DB.Collection("customers")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var customer models.Customer
	var findResult *mongo.Cursor
	var err error

	customerId := c.Params("id")
	if detail {
		findResult, err = customerCollection.Aggregate(ctx, mongo.Pipeline{
			bson.D{{"$match", bson.D{{"_id", customerId}}}},
			bson.D{{"$lookup", bson.D{
				{"from", "cases"},
				{"localField", "_id"},
				{"foreignField", "customer_id"},
				{"as", "cases"},
			}}},
			bson.D{{"$lookup", bson.D{
				{"from", "accounts"},
				{"localField", "_id"},
				{"foreignField", "customer_id"},
				{"as", "accounts"},
			}}},
		})
	} else {
		findResult, err = customerCollection.Find(ctx, bson.M{"_id": customerId})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success":      false,
			"message":      "Customer Not found",
			"error_type":   "DB_ERROR",
			"error_detail": err,
		})
	}
	if !findResult.Next(ctx) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success":    false,
			"message":    "Customer Not found",
			"error_type": "NOT_FOUND",
		})
	}
	if err := findResult.Decode(&customer); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success":      false,
			"message":      "Customer Not found",
			"error_type":   "UNMARSHALL",
			"error_detail": err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(customer)
}
