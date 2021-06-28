package api

import (
	"net/http"

	"github.com/getsentry/raven-go"
)

func Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		f := func() {
			if rvr := recover(); rvr != nil && rvr != http.ErrAbortHandler {
				raven.ClearContext()
				raven.SetHttpContext(raven.NewHttp(r))

				HandleAPIError(w, rvr)
			}
		}
		defer f()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
