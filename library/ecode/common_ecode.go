package ecode

// All common ecode
var (
	OK 						= add(0)  // Success
	InvalidParams			= add(-1) // Unknown Parameters


	RequestErr         = add(-400)
	Unauthorized       = add(-401)
	AccessDenied       = add(-403)
	NothingFound       = add(-404)
	MethodNotAllowed   = add(-405)
	Conflict           = add(-409)
	InvalidCors		   = add(-410)
	Canceled           = add(-498)
	ServerErr          = add(-500)
	ServiceUnavailable = add(-503)
	Deadline           = add(-504)
	LimitExceed        = add(-509)
)
