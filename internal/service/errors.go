package service

import "errors"

var (
	ErrEmailAlreadyRegisteredByDifferentAccount = errors.New("that email address is already registered by a different account")
	ErrPhoneAlreadyRegisteredByDifferentAccount = errors.New("that phone number is already registered by a different account")

	ErrUnowned = errors.New("you do not own that resource")

	ErrUnregisterPrimaryEmailAddress = errors.New("cannot unregister primary email address")
	ErrUnregisterUnownedEmailAddress = errors.New("cannot unregister email address you do not own")
	ErrUnregisterUnownedPhoneNumber  = errors.New("cannot unregister phone number you do not own")

	ErrNilCursorRequest = errors.New("client must provide non-nil cursor info for pagination")
)
