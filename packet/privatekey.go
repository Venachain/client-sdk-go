package packet

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"

	"git-c.i.wxblockchain.com/PlatONE/src/node/client-sdk-go/platone/keystore"
)

var ErrNullPassphrase = errors.New("passphrase is null")

type Keyfile struct {
	Address    string `json:"address"`
	Json       []byte
	Passphrase string
	privateKey *ecdsa.PrivateKey
}

func NewKeyfile(keyfilePath string) (*Keyfile, error) {
	var keyfile = new(Keyfile)

	keyJson, err := ParseFileToBytes(keyfilePath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(keyJson, keyfile)
	if err != nil {
		return nil, fmt.Errorf(ErrUnmarshalBytesFormat, keyJson, err.Error())
	}

	keyfile.Json = keyJson

	return keyfile, nil
}

func (k *Keyfile) GetPrivateKey() *ecdsa.PrivateKey {
	return k.privateKey
}

// GetPrivateKey gets the private key by decrypting the keystore file
func (k *Keyfile) ParsePrivateKey() error {

	// Decrypt key with passphrase.
	/*
		if k.Passphrase == "" {
			return nil, ErrNullPassphrase
		}*/

	key, err := keystore.DecryptKey(k.Json, k.Passphrase)
	if err != nil {
		return fmt.Errorf("Error decrypting key: %v", err)
	}

	k.privateKey = key.PrivateKey
	return nil
}
