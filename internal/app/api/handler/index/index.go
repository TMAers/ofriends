package indexhandler

import (
	"io"
	"net/http"
	"os"

	"ofriends/internal/pkg/respond"
)

type (
	// Handler is index web handler
	Handler struct{}
)

// New return new index web handler
func New() *Handler {
	return &Handler{}
}

// Index render index page
func (web *Handler) Index(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open("web/index.html")
	if err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	if _, err := io.Copy(w, f); err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
}
