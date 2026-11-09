package main

import (
	"errors"
	"net/http"
	"testing"
)

// A proven caller (holds :id's key) publishing as their own actor is allowed.
func TestAuthorizeOutboxAllowsProvenMatchingActor(t *testing.T) {
	prove := func(body []byte, sig, id string) error { return nil }
	act := Activity{Actor: "https://base/person/alice"}
	status, err := authorizeOutbox(prove, []byte("{}"), "SIG", "alice", act, "https://base/person/alice")
	if err != nil || status != http.StatusOK {
		t.Fatalf("proven + matching actor must pass: status=%d err=%v", status, err)
	}
}

// A request whose signature does not verify against :id's key is rejected 401,
// before any mutation.
func TestAuthorizeOutboxRejectsBadSignature(t *testing.T) {
	prove := func(body []byte, sig, id string) error { return errors.New("not authentic") }
	act := Activity{Actor: "https://base/person/alice"}
	status, err := authorizeOutbox(prove, []byte("{}"), "SIG", "alice", act, "https://base/person/alice")
	if err == nil || status != http.StatusUnauthorized {
		t.Fatalf("bad signature must be 401: status=%d err=%v", status, err)
	}
}

// Even a proven holder of alice's key may not publish an activity whose actor
// is someone else (bob): the published actor is bound to the outbox owner.
func TestAuthorizeOutboxRejectsForgedActor(t *testing.T) {
	prove := func(body []byte, sig, id string) error { return nil }
	act := Activity{Actor: "https://base/person/bob"}
	status, err := authorizeOutbox(prove, []byte("{}"), "SIG", "alice", act, "https://base/person/alice")
	if err == nil || status != http.StatusUnauthorized {
		t.Fatalf("a proven actor cannot publish under another actor: status=%d err=%v", status, err)
	}
}
