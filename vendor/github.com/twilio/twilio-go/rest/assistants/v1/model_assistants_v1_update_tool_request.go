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

// AssistantsV1UpdateToolRequest struct for AssistantsV1UpdateToolRequest
type AssistantsV1UpdateToolRequest struct {
	// The Assistant ID.
	AssistantId string `json:"assistant_id,omitempty"`
	// The description of the tool.
	Description string `json:"description,omitempty"`
	// True if the tool is enabled.
	Enabled bool `json:"enabled,omitempty"`
	// The metadata related to method, url, input_schema to used with the Tool.
	Meta map[string]interface{} `json:"meta,omitempty"`
	// The name of the tool.
	Name   string                          `json:"name,omitempty"`
	Policy AssistantsV1CreatePolicyRequest `json:"policy,omitempty"`
	// The type of the tool.
	Type string `json:"type,omitempty"`
}