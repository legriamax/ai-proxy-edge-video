package models

// EdgeConnectCredentials returned edge connection credentials from Chrysalis Cloud
type EdgeConnectCredentials struct {
	ID            string `json:"keyId,omitempty"` // edge key id
	PrivateKeyPem []byte `jso