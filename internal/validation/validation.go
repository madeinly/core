package validation

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/madeinly/core/models"
)

type Bag struct {
	Errors []models.Error `json:"errors,omitempty"`
}

func New() *Bag { return &Bag{} }

func (b *Bag) Add(field, code, message string) {
	b.Errors = append(b.Errors, models.Error{
		Field:   field,
		Code:    code,
		Message: message,
	})
}

func (b *Bag) HasErrors() bool { return len(b.Errors) > 0 }

func (b *Bag) WriteJSON(w io.Writer) error {
	if !b.HasErrors() {
		return nil
	}
	return json.NewEncoder(w).Encode(b)
}

func (b *Bag) WriteHTTP(w http.ResponseWriter) error {
	if !b.HasErrors() {
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	return b.WriteJSON(w)
}

// Rule validates a single field and returns every error it finds.
type rule func(string) []*models.Error

// Validate runs one rule for one value and collects its errors.
func (b *Bag) Validate(value string, r rule) {
	for _, e := range r(value) {
		if e != nil {
			b.Add(e.Field, e.Code, e.Message)
		}
	}
}
