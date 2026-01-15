// Package roll implements the core Ironsworn dice mechanic.
//
// This package contains pure domain logic only.
// It knows nothing about Telegram, Discord, HTTP, JSON, users, or rate limits.
//
// An Ironsworn roll consists of:
//   - one action die (1d6)
//   - two challenge dice (1d10, 1d10)
//   - an optional integer modifier
//
// The outcome is determined by comparing the action score
// (action die + modifier) against each challenge die separately.
//
// Important canonical rules:
//
//   - Each challenge die is compared independently.
//   - action_score > challenge_die  => win
//   - action_score <= challenge_die => loss
//     (equality counts as a loss in Ironsworn)
//
// Outcomes:
//   - 2 wins  => Success
//   - 1 win   => Partial Success
//   - 0 wins  => Failure
//
// Critical outcomes:
//   - Critical Success: 2 wins AND both challenge dice show the same value
//   - Critical Failure: 0 wins AND both challenge dice show the same value
//
// This package is intentionally small, explicit, and heavily documented.
// It is designed to be auditable and educational.
package roll
