package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"sync"
)

// zenroomMutex serializes access to zencode-exec — the binary is not thread-safe
var zenroomMutex sync.Mutex

// ZenResult mirrors the zenroom binding result
type ZenResult struct {
	Output string
	Logs   string
}

// zencodeExec runs zencode-exec with proper pipe handling to avoid deadlock
// on large inputs. Mirrors the pattern used in interfacer-dpp.
func zencodeExec(script, conf, keys, data string) (ZenResult, bool) {
	cmd := exec.Command("zencode-exec")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return ZenResult{Logs: fmt.Sprintf("stdin pipe error: %v", err)}, false
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return ZenResult{Logs: fmt.Sprintf("start error: %v", err)}, false
	}

	go func() {
		defer stdin.Close()
		// Format: conf\nscript_b64\nkeys_b64\ndata_b64\nextra_b64\ncontext_b64\n
		io.WriteString(stdin, conf)
		io.WriteString(stdin, "\n")
		io.WriteString(stdin, base64.StdEncoding.EncodeToString([]byte(script)))
		io.WriteString(stdin, "\n")
		io.WriteString(stdin, base64.StdEncoding.EncodeToString([]byte(keys)))
		io.WriteString(stdin, "\n")
		io.WriteString(stdin, base64.StdEncoding.EncodeToString([]byte(data)))
		io.WriteString(stdin, "\n")
		io.WriteString(stdin, "") // extra
		io.WriteString(stdin, "\n")
		io.WriteString(stdin, "") // context
		io.WriteString(stdin, "\n")
	}()

	err = cmd.Wait()
	return ZenResult{
		Output: stdout.String(),
		Logs:   stderr.String(),
	}, err == nil
}

// zencodeExecLocked is the mutex-protected version for concurrent handlers
func zencodeExecLocked(script, conf, keys, data string) (ZenResult, bool) {
	zenroomMutex.Lock()
	defer zenroomMutex.Unlock()
	return zencodeExec(script, conf, keys, data)
}

// callZencode is a convenience wrapper that matches the old zenroom.ZencodeExec API:
//   result, success := callZencode(script, conf, data, keys)
func callZencode(script, conf, data, keys string) (ZenResult, bool) {
	return zencodeExecLocked(script, conf, keys, data)
}

// ZenroomData holds the input/output for verify_graphql.zen
type ZenroomData struct {
	Gql            string `json:"gql"`
	EdDSASignature string `json:"eddsa_signature"`
	EdDSAPublicKey string `json:"eddsa_public_key"`
}

type ZenroomResult struct {
	Output []string `json:"output"`
}

// isAuth verifies the EdDSA signature using verify_graphql.zen contract
func (data *ZenroomData) isAuth() error {
	jsonData, _ := json.Marshal(data)
	result, success := callZencode(VERIFY, "", string(jsonData), "")
	if !success {
		return errors.New(result.Logs)
	}
	var zenroomResult ZenroomResult
	if err := json.Unmarshal([]byte(result.Output), &zenroomResult); err != nil {
		return err
	}
	if len(zenroomResult.Output) == 0 || zenroomResult.Output[0] != "1" {
		return errors.New("signature is not authentic")
	}
	return nil
}
