// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package auth

import (
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/x/mongo/driver"
)

// PLAIN is the mechanism name for PLAIN.
const PLAIN = "PLAIN"

func newPlainAuthenticator(cred *Cred, _ *http.Client) (Authenticator, error) {
	return &PlainAuthenticator{
		Username: cred.Username,
		Password: cred.Password,
	}, nil
}

// PlainAuthenticator uses the PLAIN algorithm over SASL to authenticate a connection.
type PlainAuthenticator struct {
	Username string
	Password string
}

// Auth authenticates the connection.
func (a *PlainAuthenticator) Auth(ctx context.Context, cfg *Config) error {
	return ConductSaslConversation(ctx, cfg, "$external", &plainSaslClient{
		username: a.Username,
		password: a.Password,
	})
}

// Reauth reauthenticates the connection.
func (a *PlainAuthenticator) Reauth(_ context.Context, _ *driver.AuthConfig) error {
	return newAuthError("Plain authentication does not support reauthentication", nil)
}

type plainSaslClient struct {
	username string
	password string
}

var _ SaslClient = (*plainSaslClient)(nil)

func (c *plainSaslClient) Start() (string, []byte, error) {
	b := []byte("\x00" + c.username + "\x00" + c.password)
	return PLAIN, b, nil
}

func (c *plainSaslClient) Next([]byte) ([]byte, error) {
	return nil, newAuthError("unexpected server challenge", nil)
}

func (c *plainSaslClient) Completed() bool {
	return true
}
