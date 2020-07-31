package main

import (
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// json 변환해서 저장
func putJson(stub shim.ChaincodeStubInterface, k string, v interface{}) error {
	// json 변환
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	err = stub.PutState(k, b)
	if err != nil {
		return err
	}

	return nil
}

func putDocument(stub shim.ChaincodeStubInterface, d DIDDocument) error {
	err := putJson(stub, d.Id, d)
	if err != nil {
		return newBaseError("DIDDocument 저장에 실패했습니다.")
	}
	return stub.SetEvent("PutDocument", []byte(d.Id))
}

func getDocument(stub shim.ChaincodeStubInterface, did string) (*DIDDocument, error) {
	b, err := stub.GetState(did)
	if err != nil {
		return nil, err
	}

	if b == nil {
		return nil, newBaseError("Document 가 없습니다 (did:" + did + ")")
	}

	var doc DIDDocument
	if err = json.Unmarshal(b, &doc); err != nil {
		return nil, err
	}
	return &doc, nil
}