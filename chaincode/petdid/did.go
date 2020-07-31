package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

const DefaultContext = "https://w3id.org/future-method/v1"
const DidPrefix = "did:pet:"

type DIDDocument struct {
	Context   string      `json:"@context"`
	Id        string      `json:"id"`
	PublicKey []PublicKey `json:"publicKey"`
	Service   []Service   `json:"service"`
}

type PublicKey struct {
	Id              string `json:"id"`
	Type            string `json:"type"`
	PublicKeyBase64 string `json:"publicKeyBase64"`
	Created         int    `json:"created"`
}

type Service struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	ServiceEndpoint string `json:"serviceEndpoint"`
}

func createDIDCocument(stub shim.ChaincodeStubInterface, alg string, pubKey string) (*DIDDocument, error) {
	id, err := createDID(pubKey)
	fmt.Printf("did:%s\n", id)
	if err != nil {
		return nil, err
	}

	// 중복 체크
	b, err := stub.GetState(id)
	if err != nil {
		return nil, err
	}
	if b != nil {
		return nil, newBaseError("already registerd id!")
	}

	ts, err := stub.GetTxTimestamp()
	if err != nil {
		return nil, err
	}

	publicKeys := make([]PublicKey, 1)
	publicKeys[0] = PublicKey{
		Id: id,
		Type: alg,
		PublicKeyBase64: pubKey,
		Created: int(ts.GetSeconds()),
	}
	services := make([]Service, 0)
	doc := DIDDocument{
		Context: DefaultContext,
		Id: id,
		PublicKey: publicKeys,
		Service: services,
	}
	return &doc, nil
}

func createDID(pubKey string) (string, error) {
	pub, err := UnmarshalPubkey(pubKey)
	if err != nil {
		return "", err
	}
	return DidPrefix+pubToAddress(pub), nil
}

