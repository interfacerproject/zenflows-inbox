package main

import (
	b64 "encoding/base64"
	"fmt"
	"net/http"
)

// actorProver proves that `sig` over `body` was produced by the key registered
// in zenflows for actor `id`. It is injected so authorizeOutbox is unit-testable
// without a live zenflows or a Zenroom process.
type actorProver func(body []byte, sig, id string) error

// authorizeOutbox gates a publish to /person/:id/outbox. It fails closed:
//   - the signature must verify against :id's registered key (else 401), so the
//     caller is proven to hold that actor's key;
//   - the published activity's actor must be that same actor (else 401), so a
//     proven holder cannot publish under someone else's identity.
// It returns the HTTP status to use and a non-nil error when the request is
// rejected; (http.StatusOK, nil) means proceed.
func authorizeOutbox(prove actorProver, body []byte, sig, id string, activity Activity, baseUrl string) (int, error) {
	if err := prove(body, sig, id); err != nil {
		return http.StatusUnauthorized, err
	}
	if activity.Actor != baseUrl {
		return http.StatusUnauthorized,
			fmt.Errorf("activity actor %q does not match outbox owner %q", activity.Actor, baseUrl)
	}
	return http.StatusOK, nil
}

// proveActor is the real actorProver: it verifies the signed body against the
// key zenflows holds for `id`, reusing the same path the message endpoints use.
func (inbox *Inbox) proveActor(body []byte, sig, id string) error {
	zenroomData := ZenroomData{
		Gql:            b64.StdEncoding.EncodeToString(body),
		EdDSASignature: sig,
	}
	if err := zenroomData.requestPublicKey(inbox.zfUrl, id); err != nil {
		return err
	}
	return zenroomData.isAuth()
}
