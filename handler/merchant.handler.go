package handlers

import (
	"database/sql"
	"fmt"
	"math"
	"strconv"
	"strings"
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Name format should be between 2-30 characters",
		})
	}

	// Check merchant category format
	if registerMerchantResult.MerchantCategory == "" || !helpers.ValidateMerchantCategory(registerMerchantResult.MerchantCategory) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Merchant category is not valid",
		})
	}

	// Check image url format
	if registerMerchantResult.ImageUrl == "" || !helpers.ValidateURL(registerMerchantResult.ImageUrl) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Image url is not valid",
		})
	}

	// Check location format
	if !helpers.ValidateLocation(registerMerchantResult.Lat, registerMerchantResult.Lon) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Location is not valid",
		})
	}

	// Insert data
	_, err := conn.Exec("INSERT INTO \"merchant\" (name, \"merchantCategory\",\"imageUrl\",lat, lon) VALUES ($1, $2, $3, $4, $5)", registerMerchantResult.Name, registerMerchantResult.MerchantCategory, registerMerchantResult.ImageUrl, registerMerchantResult.Lat, registerMerchantResult.Lon)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Get merchant id
	err = conn.QueryRow("SELECT id FROM \"merchant\" WHERE name = $1 AND \"merchantCategory\" = $2 AND \"imageUrl\" = $3 AND lat = $4 AND lon = $5", registerMerchantResult.Name, registerMerchantResult.MerchantCategory, registerMerchantResult.ImageUrl, registerMerchantResult.Lat, registerMerchantResult.Lon).Scan(&registerMerchantResult.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
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
	query := `SELECT * FROM merchant WHERE 1 = 1`
	// Add the WHERE clauses for the optional parameters
	if merchantId != "" {
		query += ` AND "id" = '` + merchantId + `'`
	}
	if name != "" {
		query += ` AND LOWER("name") LIKE LOWER('%` + name + `%')`
	}
	if merchantCategory != "" && helpers.ValidateMerchantCategory(merchantCategory) {
		query += ` AND "merchantCategory" = '` + merchantCategory + `'`
	}
	// Add the ORDER BY and LIMIT clauses
	if sortOrder == "asc" || sortOrder == "desc" {
		query += ` ORDER BY "createdAt" ` + sortOrder
	}
	query += ` LIMIT ` + limit + ` OFFSET ` + offset
	rows, err := conn.Query(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	defer rows.Close()

	// Prepare the data
	data := make([]map[string]interface{}, 0)
	for rows.Next() {
		var id, name, merchantCategory, imageUrl string
		var lat, long float64
		var createdAt, updatedAt time.Time
		err = rows.Scan(&id, &name, &merchantCategory, &imageUrl, &lat, &long, &createdAt, &updatedAt)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
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
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
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
		return c.Status(fiber.StatusNotFound).SendString("Merchant not found")
	}

	err := c.BodyParser(&item)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error parsing body",
		})
	}

	// Validate the request body
	if item.Name == "" || len(item.Name) < 2 || len(item.Name) > 30 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid name",
		})
	}
	if item.ProductCategory == "" || !helpers.ValidateMerchantItem(item.ProductCategory) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid product category",
		})
	}
	if item.Price == 0 || item.Price < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid price",
		})
	}
	if item.ImageUrl == "" || !helpers.ValidateURL(item.ImageUrl) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid image URL",
		})
	}

	// Check if merchant exists
	var count int
	err = conn.QueryRow("SELECT COUNT(*) FROM \"merchant\" WHERE id = $1", merchantId).Scan(&count)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if count == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Merchant not found",
		})
	}

	// Insert data
	_, err = conn.Exec("INSERT INTO \"item\" (name, \"productCategory\",\"imageUrl\",price, \"merchantId\") VALUES ($1, $2, $3, $4, $5)", item.Name, item.ProductCategory, item.ImageUrl, item.Price, merchantId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Get item id
	err = conn.QueryRow("SELECT id FROM \"item\" WHERE name = $1 AND \"productCategory\" = $2 AND \"imageUrl\" = $3 AND price = $4 AND \"merchantId\" = $5", item.Name, item.ProductCategory, item.ImageUrl, item.Price, merchantId).Scan(&item.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"itemId": item.ID,
	})
}

func MerchantGetItem(c *fiber.Ctx) error {
	conn := db.CreateConn() // Use the global db connection
	// Get the query parameters
	merchantId := c.Params("merchantId")
	itemId := c.Query("itemId", "")
	name := c.Query("name", "")
	productCategory := c.Query("productCategory", "")
	limit := c.Query("limit", "5")
	offset := c.Query("offset", "0")
	sortOrder := c.Query("createdAt", "desc")

	// Check if merchantId is empty
	if merchantId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Merchant ID is required",
		})
	}

	// Check if merchant exists
	var count int
	err := conn.QueryRow("SELECT COUNT(*) FROM \"merchant\" WHERE id = $1", merchantId).Scan(&count)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if count == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Merchant not found",
		})
	}

	// Build the base query
	query := `SELECT * FROM item WHERE "merchantId" = '` + merchantId + `'`
	// Add the WHERE clauses for the optional parameters
	if itemId != "" {
		query += ` AND "id" = '` + itemId + `'`
	}
	if name != "" {
		query += ` AND LOWER("name") LIKE LOWER('%` + name + `%')`
	}
	if productCategory != "" && helpers.ValidateMerchantItem(productCategory) {
		query += ` AND "productCategory" = '` + productCategory + `'`
	}

	// Add the ORDER BY and LIMIT clauses
	if sortOrder == "asc" || sortOrder == "desc" {
		query += ` ORDER BY "createdAt" ` + sortOrder
	}
	query += ` LIMIT ` + limit + ` OFFSET ` + offset
	rows, err := conn.Query(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	defer rows.Close()

	// Prepare the data
	data := make([]map[string]interface{}, 0)
	for rows.Next() {
		var id, merchantId, name, productCategory, imageUrl string
		var price int64
		var createdAt, updatedAt time.Time
		err = rows.Scan(&id, &name, &productCategory, &imageUrl, &price, &merchantId, &createdAt, &updatedAt)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		data = append(data, fiber.Map{
			"itemId":          id,
			"name":            name,
			"productCategory": productCategory,
			"price":           price,
			"imageUrl":        imageUrl,
			"createdAt":       createdAt.Format(time.RFC3339Nano),
		})
	}

	// Return the results as JSON
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": data,
		"meta": fiber.Map{
			"limit":  limit,
			"offset": offset,
			"total":  len(data),
		},
	})
}

func MerchantGetNearby(c *fiber.Ctx) error {

	locParam := c.Params("loc")
	coords := strings.Split(locParam, ",")
	if len(coords) != 2 {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid location")
	}
	latParam, longParam := coords[0], coords[1]
	lat, err := strconv.ParseFloat(latParam, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid latitude")
	}
	long, err := strconv.ParseFloat(longParam, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid longitude")
	}

	conn := db.CreateConn()
	merchantId := c.Query("merchantId", "")
	name := c.Query("name", "")
	merchantCategory := c.Query("merchantCategory", "")
	limit := c.Query("limit", "5")
	offset := c.Query("offset", "0")
	sortOrder := c.Query("createdAt", "desc")

	query := fmt.Sprintf(`
	SELECT *,
	    (6371 * acos(cos(radians(%f)) * cos(radians(lat)) * cos(radians(lon) - radians(%f)) + sin(radians(%f)) * sin(radians(lat)))) AS distance
	FROM merchant
	WHERE 1 = 1`, lat, long, lat)

	// Add the WHERE clauses for the optional parameters
	if merchantId != "" {
		query += ` AND "id" = '` + merchantId + `'`
	}
	if name != "" {
		query += ` AND LOWER("name") LIKE LOWER('%` + name + `%')`
	}
	if merchantCategory != "" && helpers.ValidateMerchantCategory(merchantCategory) {
		query += ` AND "merchantCategory" = '` + merchantCategory + `'`
	}
	query += ` ORDER by "distance" asc `
	// Add the ORDER BY and LIMIT clauses
	if sortOrder == "asc" || sortOrder == "desc" {
		query += ` ,"createdAt" ` + sortOrder
	}
	query += ` LIMIT ` + limit + ` OFFSET ` + offset
	rows, err := conn.Query(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	defer rows.Close()

	// Prepare the data
	data := make([]map[string]interface{}, 0)
	for rows.Next() {
		var id, name, merchantCategory, imageUrl string
		var lat, long, distance float64
		var createdAt, updatedAt time.Time
		err = rows.Scan(&id, &name, &merchantCategory, &imageUrl, &lat, &long, &createdAt, &updatedAt, &distance)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		// Query the items for this merchant
		itemRows, err := conn.Query(`SELECT id, name, "productCategory", "imageUrl", price, "createdAt" FROM item WHERE "merchantId" = $1`, id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		defer itemRows.Close()

		items := make([]map[string]interface{}, 0)
		for itemRows.Next() {
			var itemId, itemName, itemCategory, itemImage string
			var itemPrice int64
			var itemCreatedAt time.Time
			err = itemRows.Scan(&itemId, &itemName, &itemCategory, &itemImage, &itemPrice, &itemCreatedAt)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
			}

			items = append(items, fiber.Map{
				"itemId":          itemId,
				"name":            itemName,
				"productCategory": itemCategory,
				"price":           itemPrice,
				"imageUrl":        itemImage,
				"createdAt":       itemCreatedAt.Format(time.RFC3339Nano),
			})
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
			"items":     items,
			"createdAt": createdAt.Format(time.RFC3339Nano),
			// "distance": distance,
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"lat":  lat,
		"long": long,
		"data": data,
	})
}

func MerchantEstimate(c *fiber.Ctx) error {
	var MerchantEstimatePrice models.MerchantEstimatePrice
	if err := c.BodyParser(&MerchantEstimatePrice); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error()})
	}
	if len(MerchantEstimatePrice.Orders) <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "order cannot be empty"})
	}
	var startingMerchant models.MerchantModel
	listMerchants := make([]models.MerchantModel, 0)
	totalPrice := 0
	conn := db.CreateConn()
	count := 0
	for _, order := range MerchantEstimatePrice.Orders {

		if len(order.Items) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "item cannot be empty"})
		}

		var merchants []models.MerchantModel
		err := conn.Select(&merchants, "SELECT * FROM merchant WHERE id = $1", order.MerchantId)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Merchant not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		merchant := merchants[0]
		// Increment count if order is a starting point
		if order.IsStartingPoint {
			count++
			startingMerchant = merchant
		} else if (order.IsStartingPoint && count > 1) || count > 1 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Only one starting point is allowed"})
		} else {
			listMerchants = append(listMerchants, merchant)
		}

		// Convert the latitude and longitude of the first point to Cartesian coordinates
		x1, y1, _ := helpers.LatLongToCartesian(MerchantEstimatePrice.UserLocation.Lat, MerchantEstimatePrice.UserLocation.Long)

		// Convert the latitude and longitude of the second point to Cartesian coordinates
		x2, y2, _ := helpers.LatLongToCartesian(merchant.Lat, merchant.Lon)

		// Check if the area of the rectangle formed by the two points is less than or equal to 3 km²
		if !helpers.CalculateRectangleArea(x1, y1, x2, y2) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Merchant area exceeds the limit."})
		}

		// count item price
		for _, item := range order.Items {
			// check item
			var product models.MerchantItems
			err := conn.QueryRowx("SELECT id, price FROM item WHERE id = $1 AND \"merchantId\" = $2", item.ItemId, merchant.ID).StructScan(&product)
			if err != nil {
				if err == sql.ErrNoRows {
					return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Item not found"})
				}
				return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
			}
			totalPrice += item.Quantity * product.Price
		}
	}

	//TSP
	totalDistance := 0.0
	n := len(listMerchants)
	isVisited := make([]bool, n)
	tour := make([]models.MerchantModel, 0)
	currentMerchant := startingMerchant
	// make first tour from isStartingPoint
	tour = append(tour, currentMerchant)
	for len(tour) <= n {
		next := -1
		minDist := math.Inf(1)

		for i := 0; i < n; i++ {
            if !isVisited[i] {
                // Check if the merchant is within the acceptable distance
                merchantLocation := helpers.LatLongToCartesian(listMerchants[i].Lat, listMerchants[i].Lon)
                userLocation := helpers.LatLongToCartesian(MerchantEstimatePrice.UserLocation.Lat, MerchantEstimatePrice.UserLocation.Long)
                area := helpers.CalculateRectangleArea(userLocation, merchantLocation)
                if area > 3 {
                    return c.Status(400).JSON(fiber.Map{"message": "Merchant is too far >3km2"})
                }

                // Calculate the distance for the TSP
                dist := helpers.Haversine(currentMerchant.Lat, currentMerchant.Lon, listMerchants[i].Lat, listMerchants[i].Lon)
                if dist <= minDist {
                    minDist = dist
                    next = i
                }
            }
        }

		// BEFORE CARTESIAN!!!
		// for i := 0; i < n; i++ {
		// 	if !isVisited[i] {
		// 		// fmt.Println(listMerchants[i])
		// 		dist := helpers.Haversine(currentMerchant.Lat, currentMerchant.Lon, listMerchants[i].Lat, listMerchants[i].Lon)
		// 		if dist <= minDist {
		// 			minDist = dist
		// 			next = i
		// 			// fmt.Println("minDist: ", minDist, i)
		// 		}
		// 	}
		// }
		// BEFORE CARTESIAN!!!
		if next == -1 {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": next})
		}
		// mark as visited and continue from the nearest merchant
		isVisited[next] = true
		tour = append(tour, listMerchants[next])
		totalDistance += minDist
		currentMerchant = listMerchants[next]
	}
	lastMerchant := tour[len(tour)-1]
	totalDistance += helpers.Haversine(lastMerchant.Lat, lastMerchant.Lon, MerchantEstimatePrice.UserLocation.Lat, MerchantEstimatePrice.UserLocation.Long)
	fmt.Println("totalDistance: ", totalDistance)
	// count delivery in minutes
	deliveryTime := math.Ceil(totalDistance / 40 * 60)

	// insert into database
	query := "INSERT INTO estimate (\"userId\",\"userLat\", \"userLon\", \"totalPrice\", \"estimateDeliveryTime\") VALUES ($1,$2,$3,$4,$5) RETURNING id"
	var estimateId string
	userId := c.Locals("userId")
	err := conn.QueryRow(query, userId, MerchantEstimatePrice.UserLocation.Lat, MerchantEstimatePrice.UserLocation.Long, totalPrice, deliveryTime).Scan(&estimateId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	for _, order := range MerchantEstimatePrice.Orders {
		query = "INSERT INTO \"estimateOrder\" (\"estimateId\", \"merchantId\") VALUES ($1,$2) RETURNING id"
		var estimateOrderId string
		err := conn.QueryRow(query, estimateId, order.MerchantId).Scan(&estimateOrderId)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		for _, item := range order.Items {
			query = "INSERT INTO \"estimateOrderItem\" (\"estimateOrderId\", \"itemId\", quantity) VALUES ($1,$2,$3) RETURNING id"
			var estimateOrderItemId string
			err := conn.QueryRow(query, estimateOrderId, item.ItemId, item.Quantity).Scan(&estimateOrderItemId)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
			}
		}
	}
	return c.Status(200).JSON(fiber.Map{"estimatedDeliveryTimeInMinutes": deliveryTime, "totalPrice": totalPrice, "calculatedEstimateId": estimateId})
}

func MerchantPostOrder(c *fiber.Ctx) error {

	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"message": "Not Implemented",
	})
}

func MerchantGetOrders(c *fiber.Ctx) error {
	conn := db.CreateConn()
	merchantId := c.Query("merchantId", "")
	name := c.Query("name", "")
	merchantCategory := c.Query("merchantCategory", "")

	query := fmt.Sprintf(`SELECT id, "createdAt" FROM estimate WHERE "userId" = '%s'`, c.Locals("userId"))
	fmt.Println(query)
	rows, err := conn.Query(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	defer rows.Close()

	// Prepare the data
	data := make([]map[string]interface{}, 0)
	for rows.Next() {
		var orderId string
		var orderCreatedAt time.Time
		err = rows.Scan(&orderId, &orderCreatedAt)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		// Query the merchant for this order
		query := fmt.Sprintf(`SELECT "estimateOrder".id, "estimateOrder"."merchantId", merchant.name, merchant."merchantCategory", merchant."imageUrl", merchant.lat, merchant.lon, merchant."createdAt" 
                          FROM "estimateOrder"
                          JOIN merchant ON "estimateOrder"."merchantId" = merchant.id 
                          WHERE "estimateOrder"."estimateId" = '%s'`, orderId)
		// Add the WHERE clauses for the optional parameters
		if merchantId != "" {
			query += ` AND merchant.id = '` + merchantId + `'`
		}
		if name != "" {
			query += ` AND (LOWER(merchant.name) LIKE LOWER('%` + name + `%') 
                    OR EXISTS (SELECT 1 FROM "estimateOrderItem"
                               JOIN item ON "estimateOrderItem"."itemId" = item.id
                               WHERE "estimateOrderItem"."estimateOrderId" = "estimateOrder".id
                               AND LOWER(item.name) LIKE LOWER('%` + name + `%'))) `
		}
		if merchantCategory != "" && helpers.ValidateMerchantCategory(merchantCategory) {
			query += ` AND merchant."merchantCategory" = '` + merchantCategory + `'`
		}
		fmt.Println(query)
		merchantRows, err := conn.Query(query)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		defer merchantRows.Close()

		merchants := make([]map[string]interface{}, 0)
		for merchantRows.Next() {
			var estimateOrderId, merchantId, merchantName, merchantCategory, merchantImageUrl string
			var merchantLat, merchantLong float64
			var merchantCreatedAt time.Time
			err = merchantRows.Scan(&estimateOrderId, &merchantId, &merchantName, &merchantCategory, &merchantImageUrl, &merchantLat, &merchantLong, &merchantCreatedAt)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
			}

			// Query the items for this merchant
			query := fmt.Sprintf(`SELECT item.id, item.name, item."productCategory", item.price, item."imageUrl", item."createdAt", "estimateOrderItem".quantity 
                              FROM "estimateOrderItem"
                              JOIN item ON "estimateOrderItem"."itemId" = item.id 
                              WHERE "estimateOrderItem"."estimateOrderId" = '%s'`, estimateOrderId)
			// Add the WHERE clauses for the optional parameters
			if name != "" {
				query += ` AND LOWER(item.name) LIKE LOWER('%` + name + `%')`
			}
			fmt.Println(query)
			itemRows, err := conn.Query(query)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
			}

			defer itemRows.Close()

			items := make([]map[string]interface{}, 0)
			for itemRows.Next() {
				var itemId, itemName, itemCategory, itemImageUrl string
				var itemPrice, itemQuantity int64
				var itemCreatedAt time.Time
				err = itemRows.Scan(&itemId, &itemName, &itemCategory, &itemPrice, &itemImageUrl, &itemCreatedAt, &itemQuantity)
				if err != nil {
					return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
				}

				items = append(items, fiber.Map{
					"itemId":          itemId,
					"name":            itemName,
					"productCategory": itemCategory,
					"price":           itemPrice,
					"imageUrl":        itemImageUrl,
					"createdAt":       itemCreatedAt.Format(time.RFC3339Nano),
				})
			}

			merchants = append(merchants, fiber.Map{
				"merchantId":       merchantId,
				"name":             merchantName,
				"merchantCategory": merchantCategory,
				"imageUrl":         merchantImageUrl,
				"location": fiber.Map{
					"lat":  merchantLat,
					"long": merchantLong,
				},
				"items":     items,
				"createdAt": merchantCreatedAt.Format(time.RFC3339Nano),
			})
		}
		data = append(data, fiber.Map{
			"orderId":   orderId,
			"merchants": merchants,
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(data)

	// 	conn := db.CreateConn()
	// 	merchantId := c.Query("merchantId", "")
	// 	name := c.Query("name", "")
	// 	merchantCategory := c.Query("merchantCategory", "")

	// 	query := (`SELECT e.id AS "orderId", e."createdAt" AS "orderCreatedAt",
	//   ARRAY_AGG(
	//     JSON_BUILD_OBJECT(
	//       'merchantId', m.id,
	//       'name', m.name,
	//       'merchantCategory', m."merchantCategory",
	//       'imageUrl', m."imageUrl",
	//       'location', JSON_BUILD_OBJECT('lat', m.lat, 'long', m.lon),
	//       'createdAt', m."createdAt",
	//       'items', i.items
	//     )
	//   ) AS merchants
	// FROM estimate e
	// JOIN "estimateOrder" eo ON e.id = eo."estimateId"
	// JOIN merchant m ON eo."merchantId" = m.id
	// LEFT JOIN LATERAL (
	//   SELECT JSON_AGG(
	//     JSON_BUILD_OBJECT(
	//       'itemId', i.id,
	//       'name', i.name,
	//       'productCategory', i."productCategory",
	//       'price', i.price,
	//       'imageUrl', i."imageUrl",
	//       'createdAt', i."createdAt",
	//       'quantity', eoi.quantity
	//     )
	//   ) AS items
	//   FROM "estimateOrderItem" eoi
	//   JOIN item i ON eoi."itemId" = i.id
	//   WHERE eoi."estimateOrderId" = eo.id`)
	// 	if merchantId != "" {
	// 		query += ` AND m.id = '` + merchantId + `'`
	// 	}
	// 	if name != "" {
	// 		query += ` AND (LOWER(m.name) LIKE LOWER('` + name + `') OR LOWER(i.name) LIKE LOWER('` + name + `'))`
	// 	}
	// 	if merchantCategory != "" {
	// 		query += ` AND m."merchantCategory" = '` + merchantCategory + `'`
	// 	}
	// 	query += `) i ON TRUE WHERE e."userId" = '` + fmt.Sprintf("%v", c.Locals("userId")) + `' GROUP BY e.id ORDER BY e."createdAt" DESC`

	// 	fmt.Println(query)
	// 	rows, err := conn.Query(query)
	// 	if err != nil {
	// 		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	// 	}
	// 	defer rows.Close()

	// 	// Prepare the data
	// 	data := make([]map[string]interface{}, 0)
	// 	for rows.Next() {
	// 		var orderId string
	// 		var orderCreatedAt time.Time
	// 		var merchants string
	// 		err = rows.Scan(&orderId, &orderCreatedAt, &merchants)
	// 		if err != nil {
	// 			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	// 		}

	// 		fmt.Println("Merchants:", merchants)

	// 		// Remove the first `{"` and the last `"}` from the JSON string
	// 		if len(merchants) > 2 {
	// 			merchants = merchants[2 : len(merchants)-2]
	// 		}

	// 		// Replace escaped quotes with actual quotes
	// 		merchants = strings.ReplaceAll(merchants, `\"`, `"`)

	// 		fmt.Println("Cleaned Merchants JSON:", merchants)

	// 		// Parse the JSON data
	// 		var jsonData interface{}
	// 		if err := json.Unmarshal([]byte(merchants), &jsonData); err != nil {
	// 			// fmt.Println("Error:", err)
	// 			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	// 		}

	// 		// Pretty print the JSON
	// 		prettyJSON, err := json.MarshalIndent(jsonData, "", "    ")
	// 		if err != nil {
	// 			// fmt.Println("Error:", err)
	// 			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	// 		}

	// 		fmt.Println(string(prettyJSON))
	// 	}

	// 	return c.Status(fiber.StatusOK).JSON(data)

}
