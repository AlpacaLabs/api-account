package service

import (
	"context"
)

// UpdateEmailAddress updates the email address's confirmation status.
// This is usually done when a user clicks the confirmation link
// in an email they receive.
func (s Service) UpdateEmailAddress(ctx context.Context) {
	// Check if entity exists for email address
	// If not, return NotFound
	// Update the email's confirmation status
	// Return new entity in response
}
