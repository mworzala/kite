package mojang

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mworzala/kite/internal/pkg/crypto"
)

const (
	sessionServerUrl  = "https://sessionserver.mojang.com"
	hasJoinedEndpoint = "/session/minecraft/hasJoined"
)

func HasJoined(ctx context.Context, username, serverName string, sharedSecret, publicKey []byte) (*GameProfile, error) {
	authHash := crypto.Sha1([]byte(serverName), sharedSecret, publicKey)
	url := fmt.Sprintf("%s%s?username=%s&serverId=%s", sessionServerUrl, hasJoinedEndpoint, username, authHash)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Valid response if the client was not instructed to do auth (during transfer)
	if resp.StatusCode == 204 {
		return nil, nil
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("received unexpected status: %d", resp.StatusCode)
	}

	var profile GameProfile
	if err = json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, err
	}

	return &profile, nil
}
