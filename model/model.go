package model

import "time"

type (
	// Kind identify the kind of message data
	Kind int8

	// Message struct
	Message struct {
		Kind Kind   `json:"kind,omitempty"`
		File string `json:"file,omitempty"`
		Name string `json:"name,omitempty"`
	}

	// Postback struct
	Postback struct {
		Name      string    `json:"name,omitempty"`
		FullName  string    `json:"full_name,omitempty"`
		Timestamp time.Time `json:"timestamp,omitempty"`
	}
)

const (
	// Base64Kind kind. Identify the message data is Base64.
	Base64Kind Kind = iota
	// URLKind kind. Identify the message data is URL.
	URLKind
)

// Enabled returns true if the given kind equals this kind.
func (k Kind) Enabled(kind Kind) bool {
	return kind == k
}
