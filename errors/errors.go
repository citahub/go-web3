/*
Copyright 2016-2017 Cryptape Technologies LLC.

This program is free software: you can redistribute it
and/or modify it under the terms of the GNU General Public
License as published by the Free Software Foundation,
either version 3 of the License, or (at your option) any
later version.

This program is distributed in the hope that it will be
useful, but WITHOUT ANY WARRANTY; without even the implied
warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR
PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
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
