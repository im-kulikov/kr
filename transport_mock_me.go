package kr

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"sync"
	"testing"

	"golang.org/x/crypto/ssh"
)

var testSK *rsa.PrivateKey
var testPK ssh.PublicKey
var testMe *Profile
var testMeMutex sync.Mutex

func TestMe(t *testing.T) (profile Profile, sk *rsa.PrivateKey, pk ssh.PublicKey) {
	testMeMutex.Lock()
	defer testMeMutex.Unlock()
	var err error
	if testMe == nil {
		testSK, err = rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			t.Fatal(err)
		}
		testPK, err = ssh.NewPublicKey(&testSK.PublicKey)
		if err != nil {
			t.Fatal(err)
		}
		testMe = &Profile{
			SSHWirePublicKey: testPK.Marshal(),
			Email:            "kevin@krypt.co",
		}
	}
	return *testMe, testSK, testPK
}

type ResponseTransport struct {
	ImmediatePairTransport
	*testing.T
	sync.Mutex
	responses [][]byte
}

func (t *ResponseTransport) SendMessage(ps *PairingSecret, m []byte) (err error) {
	t.Lock()
	defer t.Unlock()
	var request Request
	json.Unmarshal(m, &request)
	if request.MeRequest != nil {
		testMe, _, _ := TestMe(t.T)
		resp, err := json.Marshal(Response{
			RequestID: request.RequestID,
			MeResponse: &MeResponse{
				Me: testMe,
			},
		})
		if err != nil {
			log.Fatal(err)
		}
		t.responses = append(t.responses, resp)
	}
	if request.SignRequest != nil {
	}
	return
}

func (t *ResponseTransport) Read(ps *PairingSecret) (ciphertexts [][]byte, err error) {
	pairCiphertexts, err := t.ImmediatePairTransport.Read(ps)
	ciphertexts = append(ciphertexts, pairCiphertexts...)
	t.Lock()
	defer t.Unlock()
	for _, responseBytes := range t.responses {
		ctxt, err := ps.EncryptMessage(responseBytes)
		if err != nil {
			log.Fatal(err)
		}
		ciphertexts = append(ciphertexts, ctxt)
	}
	t.responses = [][]byte{}

	return
}
