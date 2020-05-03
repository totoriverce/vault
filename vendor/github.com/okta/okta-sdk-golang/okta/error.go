/*
 * Copyright 2018 - Present Okta, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package okta

import (
	"fmt"
	"strings"
)

type Error struct {
	ErrorCode    string                   `json:"errorCode,omitempty"`
	ErrorSummary string                   `json:"errorSummary,omitempty"`
	ErrorLink    string                   `json:"errorLink,omitempty"`
	ErrorId      string                   `json:"errorId,omitempty"`
	ErrorCauses  []map[string]interface{} `json:"errorCauses,omitempty"`
}

func (e *Error) Error() string {
	formattedErr := fmt.Sprintf("The API returned an error: %s", e.ErrorSummary)

	if len(e.ErrorCauses) > 0 {
		causes := []string{}

		for _, cause := range e.ErrorCauses {
			for key, val := range cause {
				causes = append(causes, fmt.Sprintf("%s: %v", key, val))
			}
		}
		formattedErr = fmt.Sprintf("%s. Causes: %s", formattedErr, strings.Join(causes, ", "))
	}

	return formattedErr
}
