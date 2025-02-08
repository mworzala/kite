package mojangutil

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/mworzala/kite"
	"github.com/mworzala/kite/pkg/mojang"
	"time"
)

var (
	ErrInvalidNonce = errors.New("nonce did not match expected")
	ErrNoClientAuth = errors.New("client failed to authenticate against session server")
)

// HandleEncryptionResponse performs the expected proxy side steps when after a client has responded to an encryption request.
// 1. Validate returned nonce against the original one sent
// 2. Decrypt shared secret and enable encryption in connection
// 3. Complete server side session server check in after client auth.
//
// If this method does not return an error, the auth flow was successful.
// If ErrInvalidNonce is returned the client did not encrypt the nonce correctly
// If ErrNoClientAuth is returned the client did not do its side of the session server exchange (ie does not have a valid account).
func HandleEncryptionResponse(conn *kite.Conn, keyPair mojang.KeyPair, username string, encryptedVerifyToken, encryptedSharedSecret []byte) (mojang.GameProfile, error) {
	checkNonce, err := keyPair.Decrypt(encryptedVerifyToken)
	if err != nil {
		return mojang.GameProfile{}, fmt.Errorf("failed to decrypt verify token: %w", err)
	} else if !bytes.Equal(conn.GetNonce(), checkNonce) {
		return mojang.GameProfile{}, ErrInvalidNonce
	}

	// Read and write encrypted data
	var sharedSecret []byte
	if sharedSecret, err = keyPair.Decrypt(encryptedSharedSecret); err != nil {
		return mojang.GameProfile{}, fmt.Errorf("failed to decrypt shared secret: %w", err)
	}
	if err = conn.EnableEncryption(sharedSecret); err != nil {
		return mojang.GameProfile{}, err
	}

	// Do serverside auth with session server
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	profile, err := mojang.HasJoined(ctx, username, "", sharedSecret, keyPair.PublicKey())
	if err != nil {
		return mojang.GameProfile{}, fmt.Errorf("failed to complete session server auth: %w", err)
	} else if profile == nil {
		return mojang.GameProfile{}, ErrNoClientAuth
	}

	return *profile, nil
}
