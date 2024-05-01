package main

import (
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/beevik/ntp"
	"github.com/gofiber/fiber/v2"
)

func main() {
	port := os.Getenv("PORT")
	ntpServers := os.Getenv("NTP_SERVERS")
	if ntpServers == "" {
		ntpServers = "0.pool.ntp.org,1.pool.ntp.org,2.pool.ntp.org,3.pool.ntp.org"
	}
	ntpServerList := strings.Split(ntpServers, ",")

	if port == "" {
		port = "3000"
	}

	app := fiber.New()

	app.Get("/timestamp", func(c *fiber.Ctx) error {
		// Fetch the current time from an NTP server
		rand.Shuffle(len(ntpServerList), func(i, j int) {
			ntpServerList[i], ntpServerList[j] = ntpServerList[j], ntpServerList[i]
		})
		ntpServer := ntpServerList[0]
		ntpTime, err := ntp.Time(ntpServer)
		if err != nil {
			log.Printf("Error fetching time from NTP server: %s with error: %v", ntpServer, err)
			// If there's an error fetching time from NTP server, fallback to local time
			ntpTime = time.Now()
		}

		// Convert the time to the desired timezone
		tz := c.Query("tz", "UTC")
		location, err := time.LoadLocation(tz)
		if err != nil {
			log.Printf("Error loading timezone: %v", err)
			// If there's an error loading timezone, fallback to UTC
			location = time.UTC
		}
		ntpTime = ntpTime.In(location)

		// Return the timestamp as a JSON response
		return c.JSON(fiber.Map{
			"timestamp":  ntpTime.Format(time.RFC3339),
			"timezone":   location.String(),
			"ntp_server": ntpServer,
		})
	})

	// Start the server on port 3000
	app.Listen(":" + port)
}
