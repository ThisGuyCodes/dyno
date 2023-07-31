package sequences

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/thisguycodes/dyno/hamming"
)

type HammingQuery struct {
	DB          *sql.DB
	MaxDistance int
}

func (hq HammingQuery) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// POST or GET request only
	switch r.Method {
	case http.MethodPost:
	case http.MethodGet:
		break
	default:
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	incommingSequence := r.FormValue("sequence")

	// empty sequence, or no sequence provided
	if incommingSequence == "" {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// query all existing sequences
	rows, err := hq.DB.QueryContext(ctx, `
		SELECT UUID, description, sequence FROM sequences
	`)
	if err != nil {
		log.Printf("Error querying for existing sequences: %q", err)
	}
	defer rows.Close()

	matches := make([]Sequence, 0)
	for rows.Next() {
		sequence := Sequence{}
		if err := rows.Scan(&sequence.UUID, &sequence.Description, &sequence.Sequence); err != nil {
			log.Printf("Error scanning existing sequences: %q", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		if hamming.Distance(sequence.Sequence, []byte(incommingSequence)) <= hq.MaxDistance {
			matches = append(matches, sequence)
		}
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error calling .Next() on querying for existing sequences: %q", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(matches) == 0 {
		rw.WriteHeader(http.StatusNotFound)
	} else {
		rw.WriteHeader(http.StatusOK)
	}

	encoder := json.NewEncoder(rw)

	encoder.Encode(matches)
}
