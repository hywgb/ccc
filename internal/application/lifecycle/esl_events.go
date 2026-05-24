package lifecycle

import (
	"context"
	"strings"

	"github.com/divord97/ccc/internal/domain/call"
	"github.com/divord97/ccc/internal/infrastructure/esl"
)

// HandleESLEvent reacts to FreeSWITCH events received by the ESL listener and
// drives the call state machine. Without this, calls never transition out of
// `active` once the customer hangs up because nothing else observes the channel.
func (s *Service) HandleESLEvent(ctx context.Context, ev esl.Event) {
	if ev.ChannelUUID == "" {
		return
	}
	c, err := s.callSvc.FindByChannelUUID(ctx, ev.ChannelUUID)
	if err != nil || c == nil {
		return
	}

	switch ev.Name {
	case "CHANNEL_ANSWER":
		if c.AnsweredAt != nil || c.AgentUserID == nil {
			return
		}
		_, _, _ = s.AnswerCall(ctx, c.ID, *c.AgentUserID)

	case "CHANNEL_HANGUP", "CHANNEL_HANGUP_COMPLETE":
		if c.Status == call.CallStatusCompleted {
			return
		}
		_, _ = s.EndCall(ctx, c.ID, mapHangupCause(ev.HangupCause))
	}
}

// mapHangupCause translates FreeSWITCH `Hangup-Cause` strings into the domain
// HangupReason enum. Unknown causes fall back to NORMAL so downstream consumers
// always see a populated reason.
func mapHangupCause(cause string) call.HangupReason {
	switch strings.ToUpper(cause) {
	case "NORMAL_CLEARING":
		return call.HangupNormal
	case "USER_BUSY":
		return call.HangupBusy
	case "NO_ANSWER", "NO_USER_RESPONSE":
		return call.HangupNoAnswer
	case "CALL_REJECTED":
		return call.HangupReject
	case "ORIGINATOR_CANCEL", "RECOVERY_ON_TIMER_EXPIRE":
		return call.HangupAbandon
	case "EXCHANGE_ROUTING_ERROR", "DESTINATION_OUT_OF_ORDER":
		return call.HangupSystemError
	default:
		return call.HangupNormal
	}
}
