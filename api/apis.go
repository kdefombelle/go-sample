// Package api for the application.
//
// Documentation nursery API.
//
//     Schemes: http
//     BasePath: /
//     Version: 1.0.0
//     Host: localhost:3000
//
//     components:
//       securitySchemes:
//         cookieAuth:
//           type: jwtToken
//           in: cookie
//           name: toekn  # cookie name
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - basic
//
//    SecurityDefinitions:
//    basic:
//      type: basic
//
// swagger:meta
package api

import "github.com/kdefombelle/go-sample/http/rest"

// swagger:operation POST /login Signin
// Login
// ---
// tags:
//  - "authn"
// security: []
// consumes:
//  - "application/json"
// parameters:
// - name: requestBody
//   in: body
//   description: login credentials details
//   required: true
//   schema:
//     $ref: '#/definitions/LoginRequest'
// produces:
//  - "application/json"
// responses:
//   '200':
//     description:
//       Successfully authenticated.
//       The session ID is returned in a cookie named `token`.
//       You need to include this cookie in subsequent requests.
//     schema:
//       "$ref": "#/definitions/LoginResponse"

// LoginRequestWrapper wraps LoginRequest for documenting purpose
// swagger:parameters LoginRequest
type LoginRequestWrapper struct {
	// in:body
	Body rest.LoginRequest
}

// LoginResponseWrapper wraps LoginResponse for documenting purpose
// swagger:response LoginRequest
type LoginResponseWrapper struct {
	// in:body
	Body rest.LoginResponse
}
