/*
Copyright 2016-2017 Cryptape Technologies LLC.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package errors

const (
	invalidRequestCode = -32600
	methodNotFoundCode = -32601
	invalidParamsCode  = -32602
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func New(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func IsInvalidRequest(err error) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	return e.Code == invalidRequestCode
}

func IsMethodNotFound(err error) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	return e.Code == methodNotFoundCode
}

func IsInvalidParams(err error) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	return e.Code == invalidParamsCode
}

func IsNull(err error) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	return e.Code == 0
}

func (e *Error) Error() string {
	return e.Message
}
