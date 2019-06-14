package status

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
)

// Get - gets all statuses from database
func Get(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		c := r.Header.Get("Authorization")
		if c == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		c = strings.Fields(c)[1]

		tknStr := c
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})

		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		page, byPage, sort_by, sort_order, filter_field, filter_value := r.URL.Query()["page"][0], r.URL.Query()["byPage"][0], r.URL.Query()["sort_by"][0], r.URL.Query()["sort_order"][0],
			r.URL.Query()["filter_field"][0], r.URL.Query()["filter_value"][0]
		query := "SELECT * FROM status"

		if filter_value != "" && filter_field != "" {
			query += " WHERE " + filter_field + " LIKE '%" + filter_value + "%' "
		}
		if sort_by != "" {
			query += " ORDER BY " + sort_by
		}

		if sort_order != "" {
			query += " " + sort_order
		}
		if page != "" && byPage != "" {
			pageInt, _ := strconv.Atoi(page)
			byPageInt, _ := strconv.Atoi(byPage)
			offset := pageInt * byPageInt
			query += " LIMIT " + strconv.Itoa(byPageInt) + " OFFSET " + strconv.Itoa(offset)
		}

		statuses := []Status{}
		err = db.Select(&statuses, query)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(statuses)
	}
}

// GetOneByPoints - gets status by points from database
func GetOneByPoints(db *sqlx.DB, points int64) (*Status, error) {

	status := Status{}
	err := db.Get(&status, "SELECT * FROM status WHERE status.from_points <= ? and status.to_points >= ?", points, points)
	if err != nil {
		return &Status{}, err
	}

	return &status, nil

}
