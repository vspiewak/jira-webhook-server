package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/valyala/fastjson"
	"net/http"
	"os"
)

var eventToEmoji = map[string]string{
	"jira:issue_created": "💡",
	"jira:issue_updated": "🛠️",
	"jira:issue_deleted": "🔥",
}

func verifyHmacSignature(c *fiber.Ctx, secret string) bool {

	headers := c.GetReqHeaders()["X-Hub-Signature"][0]

	hmacHash := hmac.New(sha256.New, []byte(secret))
	hmacHash.Write(c.Body())
	dataHmac := hmacHash.Sum(nil)
	hmacHex := fmt.Sprintf("sha256=%s", hex.EncodeToString(dataHmac))

	return hmacHex == headers

}

func main() {

	// check env vars
	jiraWebhookSecret, ok := os.LookupEnv("JIRA_WEBHOOK_SECRET")
	if !ok {
		log.Fatal("JIRA_WEBHOOK_SECRET not set")
	}

	// setup server
	app := fiber.New()

	// setup logger
	app.Use(logger.New())

	// on issue event
	app.Post("/webhooks/jira/projects/:projectId/issues/:issueId/on-event", func(c *fiber.Ctx) error {

		log.Info("")

		// parse json body
		v, err := fastjson.ParseBytes(c.Body())
		if err != nil {

			log.Info("request body malformed")
			return c.SendStatus(http.StatusBadRequest)

		}

		// validate hmac signature
		if !verifyHmacSignature(c, jiraWebhookSecret) {

			log.Info("hmac signature not matching")
			return c.SendStatus(http.StatusUnauthorized)

		}

		// get url params
		projectId := c.Params("projectId")
		issueId := c.Params("issueId")

		// get issue infos
		issueKey, _ := v.Get("issue", "key").StringBytes()
		projectKey, _ := v.Get("issue", "fields", "project", "key").StringBytes()
		event, _ := v.Get("webhookEvent").StringBytes()

		// print issue infos
		log.Infof("%s issue %s (%s) in project %s (id: %s)", eventToEmoji[string(event)], issueKey, issueId, projectKey, projectId)

		// for each changes
		changes := v.GetArray("changelog", "items")
		for _, change := range changes {

			// print change infos
			fieldId, _ := change.Get("fieldId").StringBytes()
			from, _ := change.Get("fromString").StringBytes()
			to, _ := change.Get("toString").StringBytes()

			log.Infof("   👉 %s : %s -> %s", fieldId, from, to)

		}

		return c.SendStatus(http.StatusOK)

	})

	// start or fail
	err := app.Listen(":3000")
	if err != nil {
		log.Error(err)
	}

}
