package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)


var logger = shim.NewLogger("petdid")


type DIDService interface {
	// DID Document 생성
	Create(stub shim.ChaincodeStubInterface, jwt string, pub string) pb.Response

	// DID Document 조회
	Read(stub shim.ChaincodeStubInterface, did string) pb.Response

	// DID Doucment 에 서비스 정보 등록
	AddService(stub shim.ChaincodeStubInterface, jwt string) pb.Response
}

type PetDIDChaincode struct {

}

func (t *PetDIDChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke - Our entry point for Invocations
// ========================================
func (t *PetDIDChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	return run(stub, t)
}

func run(stub shim.ChaincodeStubInterface, service DIDService) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "version" {
		return shim.Success([]byte("0.9.1.1"))
	} else if function == "create" {
		return service.Create(stub, args[0], args[1])
	} else if function == "addService" {
		return service.AddService(stub, args[0])
	} else if function == "read" {
		return service.Read(stub, args[0])
	} else {
		return shim.Error("method 를 찾을 수 없습니다. (method:" + function + ")")
	}
	return shim.Success(nil)
}

func (t *PetDIDChaincode) Create(stub shim.ChaincodeStubInterface, jwtStr string, pub string) pb.Response {
	token, err := parse(jwtStr, pub)
	if err != nil {
		// verification error, token expired
		return shim.Error("jwt parse error:" + err.Error())
	}

	// MapCalims to CreateJwt
	var claims CreateJwt
	mapClaims := token.Claims.(jwt.MapClaims)
	err = convertToStruct(mapClaims, &claims)
	if err != nil {
		return shim.Error("mapClaims convert error:" + err.Error())
	}

	fmt.Printf("claims:%v\n", claims)
	doc, err := createDIDCocument(stub, token.Header["alg"].(string), claims.Param.PublicKeyBase64)
	if err != nil {
		return shim.Error("create document error:" + err.Error())
	}

	// db 에 저장
	err = putDocument(stub, *doc)
	if err != nil {
		return shim.Error("Document put error:" + err.Error())
	}

	// json 변환
	b, err := json.Marshal(*doc)
	if err != nil {
		return shim.Error("convert json error:" + err.Error())
	}

	return shim.Success(b)
}

func (t *PetDIDChaincode) AddService(stub shim.ChaincodeStubInterface, jwt string) pb.Response {
	return shim.Success(nil)
}

func (t *PetDIDChaincode) Read(stub shim.ChaincodeStubInterface, did string) pb.Response {
	doc, err := getDocument(stub, did)
	if err != nil {
		return shim.Error("getDocument error:" + err.Error())
	}

	// json 변환
	b, err := json.Marshal(*doc)
	if err != nil {
		return shim.Error("convert json error:" + err.Error())
	}

	return shim.Success(b)
}

func main() {
	logger.SetLevel(shim.LogDebug)
	err := shim.Start(new(PetDIDChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
