package sequences

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type Uploader struct {
	DB *sql.DB
}

func (u Uploader) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// POST request only
	if r.Method != http.MethodPost {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	compressedSequences, header, err := r.FormFile("sequences")
	if err != nil {
		log.Printf("Error getting form file from upload: %q", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// max file of 1MB
	if header.Size > 1*1024*1024 {
		rw.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}

	rawSequences, err := gzip.NewReader(compressedSequences)
	if err != nil {
		log.Printf("Error creating gzip reader for uploaded sequences: %q", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rawSequences.Close()

	scanner := bufio.NewScanner(rawSequences)

	tx, err := u.DB.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		log.Printf("Error creating transaction for storing sequences: %q", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// first, remove existing transactions
	if _, err := tx.ExecContext(ctx, `DELETE FROM sequences`); err != nil {
		fmt.Printf("Error deleting existing sequences: %q", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	thisSequence := Sequence{}
	first := true
	for scanner.Scan() {
		line := scanner.Bytes()

		switch {
		case len(line) == 0:
		case line[0] == ';':
			continue
		case line[0] == '>':
			// this is the first one, skip inserting previous
			if !first {
				// new sequence, insert the current one...
				if _, err := tx.ExecContext(ctx, `
				INSERT INTO sequences (UUID, description, sequence) VALUES ($1, $2, $3)
			`, uuid.NewString(), thisSequence.Description, thisSequence.Sequence); err != nil {
					log.Printf("Error inserting a sequence from an upload: %q", err)
					log.Print(thisSequence)
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}
			} else {
				first = false
			}

			// clear the existing one
			thisSequence.Description = scanner.Text()
			thisSequence.Sequence = make([]byte, 0)
		default:
			// normalize the sequence to uppercase
			line = bytes.ToUpper(line)
			// extend the sequence...
			thisSequence.Sequence = append(thisSequence.Sequence, line...)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Encountered an error while scanning the sequence file: %q", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		fmt.Printf("Error committing the new sequences: %q", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
