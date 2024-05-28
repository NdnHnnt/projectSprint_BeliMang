package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/NdnHnnt/projectSprint_BeliMang/db"
	helpers "github.com/NdnHnnt/projectSprint_BeliMang/helper"
	models "github.com/NdnHnnt/projectSprint_BeliMang/model"
	"github.com/gofiber/fiber/v2"
)

func MerchantRegister(c *fiber.Ctx) error {
	conn := db.CreateConn()

	var requestBody models.MerchantRequest
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Error parsing body",
		})
	}

	registerMerchantResult := models.MerchantModel{
		Name:             requestBody.Name,
		MerchantCategory: requestBody.MerchantCategory,
		ImageUrl:         requestBody.ImageUrl,
		Lat:              requestBody.Location.Lat,
		Lon:              requestBody.Location.Long,
	}

	// Check name format
	if registerMerchantResult.Name == "" || !helpers.ValidateMerchantName(registerMerchantResult.Name) {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Name format should be between 2-30 characters",
		})
	}

	// Check merchant category format
	if registerMerchantResult.MerchantCategory == "" || !helpers.ValidateMerchantCategory(registerMerchantResult.MerchantCategory) {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Merchant category is not valid",
		})
	}

	// Check image url format
	if registerMerchantResult.ImageUrl == "" || !helpers.ValidateURL(registerMerchantResult.ImageUrl) {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Image url is not valid",
		})
	}

	// Check location format
	if registerMerchantResult.Lat == 0 || registerMerchantResult.Lon == 0 || !helpers.ValidateLocation(registerMerchantResult.Lat, registerMerchantResult.Lon) {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Location is not valid",
		})
	}

	// Insert data
	_, err := conn.Exec("INSERT INTO \"merchant\" (name, \"merchantCategory\",\"imageUrl\",lat, lon) VALUES ($1, $2, $3, $4, $5)", registerMerchantResult.Name, registerMerchantResult.MerchantCategory, registerMerchantResult.ImageUrl, registerMerchantResult.Lat, registerMerchantResult.Lon)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Get merchant id
	err = conn.QueryRow("SELECT id FROM \"merchant\" WHERE name = $1 AND \"merchantCategory\" = $2 AND \"imageUrl\" = $3 AND lat = $4 AND lon = $5", registerMerchantResult.Name, registerMerchantResult.MerchantCategory, registerMerchantResult.ImageUrl, registerMerchantResult.Lat, registerMerchantResult.Lon).Scan(&registerMerchantResult.ID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"merchantId": registerMerchantResult.ID,
	})
}

func MerchantGet(c *fiber.Ctx) error {
	conn := db.CreateConn() // Use the global db connection
	// Get the query parameters
	merchantId := c.Query("merchantId", "")
	name := c.Query("name", "")
	merchantCategory := c.Query("merchantCategory", "")
	limit := c.Query("limit", "5")
	offset := c.Query("offset", "0")
	sortOrder := c.Query("createdAt", "desc")
	// Build the base query
	query := `SELECT * FROM "merchant"`
	// Add the WHERE clauses for the optional parameters
	if merchantId != "" {
		query += ` WHERE "id" = '` + merchantId + `'`
	}
	if name != "" {
		query += ` AND LOWER("name") LIKE LOWER('%` + name + `%')`
	}
	if !helpers.ValidateMerchantCategory(merchantCategory) || merchantCategory != "" {
		query += ` AND "merchantCategory" = '` + merchantCategory + `'`
	}
	// Add the ORDER BY and LIMIT clauses
	if sortOrder == "asc" || sortOrder == "desc" {
		query += ` ORDER BY "createdAt" ` + sortOrder
	}
	query += ` LIMIT ` + limit + ` OFFSET ` + offset
	rows, err := conn.Query(query)
	if err != nil {
		log.Println("Failed to execute the query:", err)
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	defer rows.Close()

	// Prepare the data
	data := make([]map[string]interface{}, 0)
	for rows.Next() {
		var id, name, merchantCategory, imageUrl string
		var lat, long float64
		var createdAt time.Time
		err = rows.Scan(&id, &name, &merchantCategory, &imageUrl, &lat, &long, &createdAt)
		if err != nil {
			log.Println("Failed to scan row:", err)
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}

		data = append(data, fiber.Map{
			"merchantId":       id,
			"name":             name,
			"merchantCategory": merchantCategory,
			"imageUrl":         imageUrl,
			"location": fiber.Map{
				"lat":  lat,
				"long": long,
			},
			"createdAt": createdAt.Format(time.RFC3339Nano),
		})
	}

	// Return the results as JSON
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"data": data,
		"meta": fiber.Map{
			"limit":  limit,
			"offset": offset,
			"total":  len(data),
		},
	})
}

func MerchantRegisterItem(c *fiber.Ctx) error {
	conn := db.CreateConn()
	var item models.MerchantItems

	merchantId := c.Params("merchantId")
	if merchantId == "" {
		return c.Status(http.StatusNotFound).SendString("Merchant not found")
	}

	err := c.BodyParser(&item)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Error parsing body",
		})
	}

	// Validate the request body
	if item.Name == "" || len(item.Name) < 2 || len(item.Name) > 30 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid name",
		})
	}
	if item.ProductCategory == "" || !helpers.ValidateMerchantItem(item.ProductCategory) {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid product category",
		})
	}
	if item.Price == 0 || item.Price < 1 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid price",
		})
	}
	if item.ImageUrl == "" || !helpers.ValidateURL(item.ImageUrl) {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid image URL",
		})
	}

	// Check if merchant exists
	var count int
	err = conn.QueryRow("SELECT COUNT(*) FROM \"merchant\" WHERE id = $1", merchantId).Scan(&count)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if count == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "Merchant not found",
		})
	}

	// Insert data
	_, err = conn.Exec("INSERT INTO \"item\" (name, \"productCategory\",\"imageUrl\",price, \"merchantId\") VALUES ($1, $2, $3, $4, $5)", item.Name, item.ProductCategory, item.ImageUrl, item.Price, merchantId)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Get item id
	err = conn.QueryRow("SELECT id FROM \"item\" WHERE name = $1 AND \"productCategory\" = $2 AND \"imageUrl\" = $3 AND price = $4 AND \"merchantId\" = $5", item.Name, item.ProductCategory, item.ImageUrl, item.Price, merchantId).Scan(&item.ID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"itemId": item.ID,
	})
}

func MerchantGetItem(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"message": "Not Implemented",
	})
}

func MerchantGetNearby(c *fiber.Ctx) error {
	lat := c.Params("lat")
	long := c.Params("long")

	// You can now use lat and long in your function

	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"message": "Not Implemented",
		lat:       lat,
		long:      long,
	})
}
