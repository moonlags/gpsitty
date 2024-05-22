package server

import (
	"log"
	"net/http"
)

type securedHanlder func(http.ResponseWriter, *http.Request, string)

func (s *Server) SecuredRoute(next securedHanlder) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.Store.Get(r, "gpsitty")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Printf("failed to get session: %v\n", err)
			return
		}

		id, ok := session.Values["id"]
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next(w, r, id.(string))
	})
}
