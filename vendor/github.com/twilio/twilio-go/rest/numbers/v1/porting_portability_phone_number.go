/*
 * This code was generated by
 * ___ _ _ _ _ _    _ ____    ____ ____ _    ____ ____ _  _ ____ ____ ____ ___ __   __
 *  |  | | | | |    | |  | __ |  | |__| | __ | __ |___ |\ | |___ |__/ |__|  | |  | |__/
 *  |  |_|_| | |___ | |__|    |__| |  | |    |__] |___ | \| |___ |  \ |  |  | |__| |  \
 *
 * Twilio - Numbers
 * This is the public Twilio REST API.
 *
 * NOTE: This class is auto generated by OpenAPI Generator.
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

package openapi

import (
	"encoding/json"
	"net/url"
	"strings"
)

// Optional parameters for the method 'FetchPortingPortability'
type FetchPortingPortabilityParams struct {
	// Account Sid to which the number will be ported. This can be used to determine if a sub account already has the number in its inventory or a different sub account. If this is not provided, the authenticated account will be assumed to be the target account.
	TargetAccountSid *string `json:"TargetAccountSid,omitempty"`
	// Address Sid of customer to which the number will be ported.
	AddressSid *string `json:"AddressSid,omitempty"`
}

func (params *FetchPortingPortabilityParams) SetTargetAccountSid(TargetAccountSid string) *FetchPortingPortabilityParams {
	params.TargetAccountSid = &TargetAccountSid
	return params
}
func (params *FetchPortingPortabilityParams) SetAddressSid(AddressSid string) *FetchPortingPortabilityParams {
	params.AddressSid = &AddressSid
	return params
}

// Check if a single phone number can be ported to Twilio
func (c *ApiService) FetchPortingPortability(PhoneNumber string, params *FetchPortingPortabilityParams) (*NumbersV1PortingPortability, error) {
	path := "/v1/Porting/Portability/PhoneNumber/{PhoneNumber}"
	path = strings.Replace(path, "{"+"PhoneNumber"+"}", PhoneNumber, -1)

	data := url.Values{}
	headers := map[string]interface{}{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	if params != nil && params.TargetAccountSid != nil {
		data.Set("TargetAccountSid", *params.TargetAccountSid)
	}
	if params != nil && params.AddressSid != nil {
		data.Set("AddressSid", *params.AddressSid)
	}

	resp, err := c.requestHandler.Get(c.baseURL+path, data, headers)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	ps := &NumbersV1PortingPortability{}
	if err := json.NewDecoder(resp.Body).Decode(ps); err != nil {
		return nil, err
	}

	return ps, err
}