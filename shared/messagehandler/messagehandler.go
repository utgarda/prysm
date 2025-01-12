package messagehandler

import (
	"context"
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
)

var log = logrus.WithField("prefix", "message-handler")

// SafelyHandleMessage will recover and log any panic that occurs from the
// function argument.
func SafelyHandleMessage(ctx context.Context, fn func(message proto.Message), msg proto.Message) {
	defer func() {
		if r := recover(); r != nil {
			printedMsg := "message contains no data"
			if msg != nil {
				printedMsg = proto.MarshalTextString(msg)
			}
			log.WithFields(logrus.Fields{
				"r":   r,
				"msg": printedMsg,
			}).Error("Panicked when handling p2p message! Recovering...")

			if ctx == nil {
				return
			}
			if span := trace.FromContext(ctx); span != nil {
				span.SetStatus(trace.Status{
					Code:    trace.StatusCodeInternal,
					Message: fmt.Sprintf("Panic: %v", r),
				})
			}
		}
	}()

	// Fingers crossed that it doesn't panic...
	fn(msg)
}
