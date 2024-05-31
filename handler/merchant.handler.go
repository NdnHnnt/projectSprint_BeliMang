package handlers

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
	"database/sql"

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
	if !helpers.ValidateLocation(registerMerchantResult.Lat, registerMerchantResult.Lon) {
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
	// fmt.Println(query)
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
		var createdAt, updatedAt time.Time
		err = rows.Scan(&id, &name, &merchantCategory, &imageUrl, &lat, &long, &createdAt, &updatedAt)
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
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Merchant ID is required",
		})
	}

	// Check if merchant exists
	var count int
	err := conn.QueryRow("SELECT COUNT(*) FROM \"merchant\" WHERE id = $1", merchantId).Scan(&count)
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
	// fmt.Println(query)
	rows, err := conn.Query(query)
	if err != nil {
		log.Println("Failed to execute the query:", err)
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
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
			log.Println("Failed to scan row:", err)
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
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
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"data": data,
		"meta": fiber.Map{
			"limit":  limit,
			"offset": offset,
			"total":  len(data),
		},
	})
}

func MerchantGetNearby(c *fiber.Ctx) error {

	// fmt.Println("masuk")
	locParam := c.Params("loc")
	coords := strings.Split(locParam, ",")
	if len(coords) != 2 {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid location")
	}
	latParam, longParam := coords[0], coords[1]
	// fmt.Println("latParam:", latParam)
	lat, err := strconv.ParseFloat(latParam, 64)
	if err != nil {
		// fmt.Println("Failed to parse lat:", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid latitude")
	}
	// fmt.Println("longParam:", longParam)
	long, err := strconv.ParseFloat(longParam, 64)
	if err != nil {
		// fmt.Println("Failed to parse long:", err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid longitude")
	}

	conn := db.CreateConn()
	// Get the query parameters
	merchantId := c.Query("merchantId", "")
	name := c.Query("name", "")
	merchantCategory := c.Query("merchantCategory", "")
	limit := c.Query("limit", "5")
	offset := c.Query("offset", "0")
	sortOrder := c.Query("createdAt", "desc")

	// Build the base query
	// fmt.Println("query")
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
	fmt.Println(query)
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
		var lat, long, distance float64
		var createdAt, updatedAt time.Time
		err = rows.Scan(&id, &name, &merchantCategory, &imageUrl, &lat, &long, &createdAt, &updatedAt, &distance)
		if err != nil {
			log.Println("Failed to scan row:", err)
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}

		// Query the items for this merchant
		itemRows, err := conn.Query(`SELECT id, name, "productCategory", "imageUrl", price, "createdAt" FROM item WHERE "merchantId" = $1`, id)
		if err != nil {
			log.Println("Failed to execute the item query:", err)
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}

		defer itemRows.Close()

		items := make([]map[string]interface{}, 0)
		for itemRows.Next() {
			var itemId, itemName, itemCategory, itemImage string
			var itemPrice int64
			var itemCreatedAt time.Time
			err = itemRows.Scan(&itemId, &itemName, &itemCategory, &itemImage, &itemPrice, &itemCreatedAt)
			if err != nil {
				log.Println("Failed to scan item row:", err)
				return c.Status(http.StatusInternalServerError).SendString(err.Error())
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
        return c.Status(400).JSON(fiber.Map{
            "message": err.Error()})
    }
    if len(MerchantEstimatePrice.Orders) <= 0 {
        return c.Status(400).JSON(fiber.Map{"message": "order cannot be empty"})
    }
    var startingMerchant models.MerchantModel
    listMerchants := make([]models.MerchantModel, 0)
    totalPrice := 0
    conn := db.CreateConn()
    count := 0
    for _, order := range MerchantEstimatePrice.Orders {

        if len(order.Items) == 0 {
            return c.Status(400).JSON(fiber.Map{"message": "item cannot be empty"})
        }
        // check merchant
        // var merchant models.MerchantModel
				var merchants []models.MerchantModel
				err := conn.Select(&merchants, "SELECT * FROM merchant WHERE id = $1", order.MerchantId)
				if err != nil {
						if err ==sql.ErrNoRows{
							c.JSON(404, gin.H{"message":"merchant not found"})
					return
				}
				fmt.Println(err.Error())
				c.JSON(500,"server error")
				return
			}
				merchant := merchants[0]
        if order.IsStartingPoint {
            if count != 0 {
                return c.Status(400).JSON(fiber.Map{"message": "invalid starting point"})
            }
            count = 1
            startingMerchant = merchant
        } else {
            listMerchants = append(listMerchants, merchant)
        }
        // count item price
        for _, item := range order.Items {
            // check item
            var product models.MerchantItems
            err := conn.QueryRowx("SELECT id, price FROM item WHERE id = $1 AND \"merchantId\" = $2", item.ItemId, merchant.ID).StructScan(&product)
            fmt.Println(item.ItemId, merchant.ID)
            if err != nil {
                if err == sql.ErrNoRows {
                    return c.Status(404).JSON(fiber.Map{"message": "item not found"})
                }
                fmt.Println(err.Error())
                return c.Status(500).SendString("server error")
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
                // fmt.Println(listMerchants[i])
                dist := helpers.Haversine(currentMerchant.Lat, currentMerchant.Lon, listMerchants[i].Lat, listMerchants[i].Lon)
                if dist <= minDist {
                    minDist = dist
                    next = i
                    fmt.Println("minDist: ", minDist, i)
                }
            }
        }
        if next == -1 {
            return c.Status(500).JSON(fiber.Map{"message": "server error"})
        }
        // mark as visited and continue from the nearest merchant
        isVisited[next] = true
        tour = append(tour, listMerchants[next])
        totalDistance += minDist
        currentMerchant = listMerchants[next]
        fmt.Println("total :", totalDistance)
    }
    lastMerchant := tour[len(tour)-1]
    totalDistance += helper.CountHaversine(lastMerchant.Lat, lastMerchant.Lon, estimateForm.UserLocation.Lat, estimateForm.UserLocation.Long)
    fmt.Println("total :", totalDistance)
    // count delivery in minutes
    deliveryTime := math.Round(totalDistance / 40 * 60)
    // insert into database
    query := "INSERT INTO estimate (\"userId\",\"userLat\", \"userLon\", \"totalPrice\", \"estimateDeliveryTime\") VALUES ($1,$2,$3,$4,$5) RETURNING id"
    var estimateId string
    err := conn.QueryRow(query, userId, estimateForm.UserLocation.Lat, estimateForm.UserLocation.Long, totalPrice, deliveryTime).Scan(&estimateId)
    if err != nil {
        fmt.Println(err.Error())
        return c.Status(500).JSON(fiber.Map{"message": "server error"})
    }
    for _, order := range estimateForm.Orders {
        query = "INSERT INTO \"estimateOrder\" (\"estimateId\", \"isStarting\", \"merchantId\") VALUES ($1,$2,$3) RETURNING id"
        var estimateOrderId string
        err := conn.QueryRow(query, estimateId, order.IsStartingPoint, order.MerchantId).Scan(&estimateOrderId)
        if err != nil {
            fmt.Println(err.Error())
            return c.Status(500).JSON(fiber.Map{"message": "server error"})
        }
        for _, item := range order.Items {
            query = "INSERT INTO \"estimateOrderItem\" (\"estimateOrderId\", \"itemId\", quantity) VALUES ($1,$2,$3) RETURNING id"
            var estimateOrderItemId string
            err := conn.QueryRow(query, estimateOrderId, item.ItemId, item.Quantity).Scan(&estimateOrderItemId)
            if err != nil {
                fmt.Println(err.Error())
                return c.Status(500).JSON(fiber.Map{"message": "server error"})
            }
        }
    }
    return c.Status(200).JSON(fiber.Map{"estimatedDeliveryTimeInMinutes": deliveryTime, "totalPrice": totalPrice, "calculatedEstimateId": estimateId})
}

// func MerchantEstimate(c *fiber.Ctx) error {
// 	var MerchantEstimatePrice models.MerchantEstimatePrice
// 	err := c.BodyParser(&MerchantEstimatePrice)
// 	if err != nil {
// 		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
// 			"message": "Error parsing body",
// 		})
// 	}

// 	var MerchantEstimateOrder models.MerchantEstimateOrder


// 	var MerchantEstimateOrderItem models.MerchantEstimateOrderItem


	// return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
	// 	"message": "Not Implemented",
	// })
// }

func MerchantPostOrder(c *fiber.Ctx) error {

	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"message": "Not Implemented",
	})
}

func MerchantGetOrder(c *fiber.Ctx) error {
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
	// fmt.Println(query)
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
		var createdAt, updatedAt time.Time
		err = rows.Scan(&id, &name, &merchantCategory, &imageUrl, &lat, &long, &createdAt, &updatedAt)
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
