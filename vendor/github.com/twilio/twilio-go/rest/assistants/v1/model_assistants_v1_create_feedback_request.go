/*
 * This code was generated by
 * ___ _ _ _ _ _    _ ____    ____ ____ _    ____ ____ _  _ ____ ____ ____ ___ __   __
 *  |  | | | | |    | |  | __ |  | |__| | __ | __ |___ |\ | |___ |__/ |__|  | |  | |__/
 *  |  |_|_| | |___ | |__|    |__| |  | |    |__] |___ | \| |___ |  \ |  |  | |__| |  \
 *
 * Twilio - Assistants
 * This is the public Twilio REST API.
 *
 * NOTE: This class is auto generated by OpenAPI Generator.
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

package openapi

import (
	"encoding/json"

	"github.com/twilio/twilio-go/client"
)

// AssistantsV1CreateFeedbackRequest struct for AssistantsV1CreateFeedbackRequest
type AssistantsV1CreateFeedbackRequest struct {
	// The message ID.
	MessageId string `json:"message_id,omitempty"`
	// The score to be given(0-1).
	Score float32 `json:"score,omitempty"`
	// The Session ID.
	SessionId string `json:"session_id"`
	// The text to be given as feedback.
	Text string `json:"text,omitempty"`
}

func (response *AssistantsV1CreateFeedbackRequest) UnmarshalJSON(bytes []byte) (err error) {
	raw := struct {
		MessageId string      `json:"message_id"`
		Score     interface{} `json:"score"`
		SessionId string      `json:"session_id"`
		Text      string      `json:"text"`
	}{}

	if err = json.Unmarshal(bytes, &raw); err != nil {
		return err
	}

	*response = AssistantsV1CreateFeedbackRequest{
		MessageId: raw.MessageId,
		SessionId: raw.SessionId,
		Text:      raw.Text,
	}

	responseScore, err := client.UnmarshalFloat32(&raw.Score)
	if err != nil {
		return err
	}
	response.Score = *responseScore

	return
}