package flash

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var store *session.Store

// Severity represents the severity level of a flash message
type Severity string

// Message represents a flash message
type Message struct {
	Text     string   `json:"text"`
	Severity Severity `json:"severity"`
}

type Callback func(*fiber.Ctx)

const flashKey = "flash_messages"

const (
	DEBUG   Severity = "DEBUG"
	INFO    Severity = "INFO"
	SUCCESS Severity = "SUCCESS"
	WARNING Severity = "WARNING"
	ERROR   Severity = "ERROR"
)

// String implements the Stringer interface for Severity
func (s Severity) String() string {
	return string(s)
}

// Setup initializes the flash message system with a session store
func Setup(s *session.Store) {
	store = s
	store.RegisterType([]Message{})
}

// AddMessage adds a flash message to the session
func AddMessage(c *fiber.Ctx, severity Severity, text string) *fiber.Ctx {
	if store == nil {
		log.Warn("flash message system not initialized")
		return c
	}

	sess, err := store.Get(c)
	if err != nil {
		log.Errorf("Error getting session: %v\n", err)
		return c
	}

	messages := getMessages(sess)
	messages = append(messages, Message{
		Text:     text,
		Severity: severity,
	})

	sess.Set(flashKey, messages)
	err = sess.Save()

	if err != nil {
		log.Errorf("Error saving session: %v\n", err)
	}

	return c
}

// getMessages retrieves flash messages from the session
func getMessages(sess *session.Session) []Message {
	if v := sess.Get(flashKey); v != nil {
		if messages, ok := v.([]Message); ok {
			return messages
		}
	}
	return []Message{}
}

// Debug adds a debug message
func Debug(c *fiber.Ctx, text string) {
	AddMessage(c, DEBUG, text)
}

// Info adds an info message
func Info(c *fiber.Ctx, text string) {
	AddMessage(c, INFO, text)
}

// Success adds a success message
func Success(c *fiber.Ctx, text string) {
	AddMessage(c, SUCCESS, text)
}

// Warning adds a warning message
func Warning(c *fiber.Ctx, text string) {
	AddMessage(c, WARNING, text)
}

// Error adds an error message
func Error(c *fiber.Ctx, text string) {
	AddMessage(c, ERROR, text)
}

// Middleware is a Fiber middleware that handles flash messages
func Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if store == nil {
			return fiber.NewError(fiber.StatusInternalServerError, "flash message system not initialized")
		}

		sess, err := store.Get(c)
		if err != nil {
			return err
		}

		if messages := getMessages(sess); len(messages) > 0 {
			if err := c.Bind(fiber.Map{
				"flashMessages": messages,
			}); err != nil {
				return err
			}

			// Clear the messages after consuming them
			sess.Delete(flashKey)
			if err := sess.Save(); err != nil {
				return err
			}
		}

		return c.Next()
	}
}
