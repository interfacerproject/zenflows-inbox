package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const GQL_PERSON_PUBKEY string = "query($id: ID!) {personPubkey(id: $id)}"

// Fills ZenroomData with the public key requested from zenflows (by person ID)
func (data *ZenroomData) requestPublicKey(url string, id string) error {
	query, err := json.Marshal(map[string]interface{}{
		"query": GQL_PERSON_PUBKEY,
		"variables": map[string]string{
			"id": id,
		},
	})
	resp, err := http.Post(url, "application/json", bytes.NewReader(query))
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var result map[string]map[string]string
	json.Unmarshal(body, &result)
	if result["data"]["personPubkey"] == "" {
		return errors.New(string(body))
	}
	data.EdDSAPublicKey = result["data"]["personPubkey"]
	return nil
}
