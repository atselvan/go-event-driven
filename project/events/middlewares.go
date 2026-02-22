package events

import (
	"log/slog"
	"time"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/lithammer/shortuuid/v3"
)

func useMiddleware(router *message.Router, logger watermill.LoggerAdapter) {
	// middleware to recover from a panic
	router.AddMiddleware(middleware.Recoverer)

	// middleware for retry
	router.AddMiddleware(middleware.Retry{
		MaxRetries:      10,
		InitialInterval: time.Microsecond * 100,
		MaxInterval:     time.Second,
		Multiplier:      2,
		Logger:          logger,
	}.Middleware)

	// middleware to add a generated correlation ID.
	router.AddMiddleware(func(next message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			correlationID := msg.Metadata.Get("correlation_id")
			if correlationID == "" {
				correlationID = "gen_" + shortuuid.New()
			}

			ctx := log.ContextWithCorrelationID(msg.Context(), correlationID)
			msg.SetContext(ctx)

			return next(msg)
		}
	})

	// middleware to get correlation ID from the current context.
	router.AddMiddleware(func(next message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			correlationID := log.CorrelationIDFromContext(msg.Context())

			ctx := log.ToContext(msg.Context(), slog.With("correlation_id", correlationID))
			msg.SetContext(ctx)

			return next(msg)
		}
	})

	// middleware for logging
	router.AddMiddleware(func(next message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			logger := log.FromContext(msg.Context())

			logger = logger.With(
				"message_id", msg.UUID,
				"payload", string(msg.Payload),
				"metadata", msg.Metadata,
				"handler", message.HandlerNameFromCtx(msg.Context()),
			)
			logger.Info("Handling a message")

			msgs, err := next(msg)
			if err != nil {
				logger = logger.With(
					"message_id", msg.UUID,
					"error", err,
				)
				logger.Info("Error while handling a message")
			}

			return msgs, err
		}
	})
}
