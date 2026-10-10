// Copyright 2023-2026 Ant Investor Ltd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package business

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"connectrpc.com/connect"
	"github.com/pitabwire/frame/v2/data"
)

// PropertyJurisdiction is the profile property naming the legal jurisdiction
// a person is subject to, as an ISO 3166-1 alpha-2 code (e.g. "KE"). It is
// stored upper-case and must name a country the service knows.
const PropertyJurisdiction = "jurisdiction"

var iso3166Alpha2 = regexp.MustCompile(`^[A-Z]{2}$`)

// validateProperties checks and normalises the typed profile properties in
// place. Properties without a declared type pass through unchanged.
func (pb *profileBusiness) validateProperties(ctx context.Context, props data.JSONMap) error {
	raw, ok := props[PropertyJurisdiction]
	if !ok {
		return nil
	}
	code, isString := raw.(string)
	if !isString {
		return connect.NewError(connect.CodeInvalidArgument,
			errors.New("jurisdiction must be an ISO 3166-1 alpha-2 country code"))
	}
	code = strings.ToUpper(strings.TrimSpace(code))
	if !iso3166Alpha2.MatchString(code) {
		return connect.NewError(connect.CodeInvalidArgument,
			fmt.Errorf("jurisdiction %q is not an ISO 3166-1 alpha-2 country code", raw))
	}
	known, err := pb.addressBusiness.IsCountryISO2(ctx, code)
	if err != nil {
		return err
	}
	if !known {
		return connect.NewError(connect.CodeInvalidArgument,
			fmt.Errorf("jurisdiction %q is not a known ISO 3166-1 country", code))
	}
	props[PropertyJurisdiction] = code
	return nil
}
