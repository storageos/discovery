package handlers

import (
	"fmt"
	"net/http"
)

func robotsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "User-agent: *\nDisallow: /")
}
