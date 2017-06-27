package handlers

import (
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r,
		"https://docs.storageos.com/",
		http.StatusMovedPermanently,
	)
}
