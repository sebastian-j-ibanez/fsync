package protocol

import (
	"crypto/rand"

	"github.com/cloudflare/circl/hpke"
)

// Setup two-way hpke encryption for listener
// Returns opener, sealer, and potential error
func (s *SocketHandler) setupServerEncryption() (hpke.Opener, hpke.Sealer, error) {
	// Initialize hpke suite
	kemID := hpke.KEM_P384_HKDF_SHA384
	kdfID := hpke.KDF_HKDF_SHA384
	aeadID := hpke.AEAD_AES256GCM
	suite := hpke.NewSuite(kemID, kdfID, aeadID)

	// Generate key pair
	publicServer, privateServer, err := kemID.Scheme().GenerateKeyPair()
	if err != nil {
		return nil, nil, err
	}

	// Marhsall and send public key
	b, err := publicServer.MarshalBinary()
	if err != nil {
		return nil, nil, err
	}
	err = s.Enc.Encode(b)
	if err != nil {
		return nil, nil, err
	}

	// Receive and unmarshal client's public key
	var publicClientBytes []byte
	err = s.Dec.Decode(&publicClientBytes)
	if err != nil {
		return nil, nil, err
	}
	publicClient, err := kemID.Scheme().UnmarshalBinaryPublicKey(publicClientBytes)
	if err != nil {
		return nil, nil, err
	}

	// Init sender and receiver
	sender, err := suite.NewSender(publicClient, []byte{})
	if err != nil {
		return nil, nil, err
	}
	receiver, err := suite.NewReceiver(privateServer, []byte{})
	if err != nil {
		return nil, nil, err
	}

	// Receive client's encapsulated key
	var clientEnc []byte
	err = s.Dec.Decode(&clientEnc)
	if err != nil {
		return nil, nil, err
	}

	// Setup sealer and send encapsulated key to client
	serverEnc, sealer, err := sender.Setup(rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	err = s.Enc.Encode(serverEnc)
	if err != nil {
		return nil, nil, err
	}

	// Setup opener from client's encapsulated key
	opener, err := receiver.Setup(clientEnc)
	if err != nil {
		return nil, nil, err
	}

	s.Opener = opener
	s.Sealer = sealer

	return opener, sealer, nil
}

// Setup two-way hpke encryption for sender
// Returns opener, sealer, and potential error
func (s *SocketHandler) setupClientEncryption() (hpke.Opener, hpke.Sealer, error) {
	// Initialize hpke suite
	kemID := hpke.KEM_P384_HKDF_SHA384
	kdfID := hpke.KDF_HKDF_SHA384
	aeadID := hpke.AEAD_AES256GCM
	suite := hpke.NewSuite(kemID, kdfID, aeadID)

	// Generate key pair
	publicClient, privateClient, err := kemID.Scheme().GenerateKeyPair()
	if err != nil {
		return nil, nil, err
	}

	// Receive and unmarshal server's public key
	var publicServerBytes []byte
	err = s.Dec.Decode(&publicServerBytes)
	if err != nil {
		return nil, nil, err
	}
	publicServer, err := kemID.Scheme().UnmarshalBinaryPublicKey(publicServerBytes)
	if err != nil {
		return nil, nil, err
	}

	// Send public key to server
	pk, err := publicClient.MarshalBinary()
	if err != nil {
		return nil, nil, err
	}
	err = s.Enc.Encode(pk)
	if err != nil {
		return nil, nil, err
	}

	// Init sender and receiver
	sender, err := suite.NewSender(publicServer, []byte{})
	if err != nil {
		return nil, nil, err
	}
	receiver, err := suite.NewReceiver(privateClient, []byte{})
	if err != nil {
		return nil, nil, err
	}

	// Setup sealer and send encapsulated key to server
	clientEnc, sealer, err := sender.Setup(rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	err = s.Enc.Encode(clientEnc)
	if err != nil {
		return nil, nil, err
	}

	// Receive server's encapsulated key
	var serverEnc []byte
	err = s.Dec.Decode(&serverEnc)
	if err != nil {
		return nil, nil, err
	}

	// Setup opener from server's encapsulated key
	opener, err := receiver.Setup(serverEnc)
	if err != nil {
		return nil, nil, err
	}

	s.Opener = opener
	s.Sealer = sealer

	return opener, sealer, nil
}
