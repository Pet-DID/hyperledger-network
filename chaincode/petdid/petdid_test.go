package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"testing"

)

func checkInit(t *testing.T, stub *shim.MockStub) {
	res := stub.MockInit("1", nil)

	if res.Status != shim.OK {
		fmt.Println("Init failed", res.Message)
		t.FailNow()
	}
}

func checkState(t *testing.T, stub *shim.MockStub, name string, value string) {
	bytes := stub.State[name]
	if bytes == nil {
		fmt.Println("State", name, "failed to get Value")
		t.FailNow()
	}
	if string(bytes) != value {
		fmt.Println("State value ", name, " was not ", value, " as expected")
		t.FailNow()
	}
}

func checkInvoke(t *testing.T, stub *shim.MockStub, funcName string, arguments... []byte) pb.Response {
	args := make([][]byte, 0)
	args = append(args, []byte(funcName))
	for _, arg := range arguments {
		args = append(args, arg)
	}

	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Invoke failed", res.Message)
		t.FailNow()
	} else {
		fmt.Println("Invoke Success message:", string(res.Payload))
	}
	return res
}

func checkQuery(t *testing.T, stub *shim.MockStub, funcName string, arguments... []byte) {
	args := make([][]byte, 0)
	args = append(args, []byte(funcName))
	for _, arg := range arguments {
		args = append(args, arg)
	}

	res := stub.MockInvoke("1", args)
	fmt.Println(res.String())
	if res.Status != shim.OK {
		fmt.Println("Query '"+string(args[0])+"' failed ", res.Message)
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("Query '"+string(args[0])+"'", nil)
		t.FailNow()
		//t.SkipNow()
	}
}

func createPrivateKey() (*ecdsa.PrivateKey, error) {
	// ecdsa 256
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	return privateKey, nil
}

// did create 요청 jwt 생성
func Create(privateKey *ecdsa.PrivateKey) (string, string, error) {

	// publickey to base64 string
	publicKeyECDSA, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

		// Create the Claims
	pub := Base64FromECDSAPub(publicKeyECDSA)
	claims := CreateJwt{
		jwt.StandardClaims{
			Issuer: "test",
		},
		CreateParam{
			Id:              "did:pet:0000",
			Type:            "ES256",
			PublicKeyBase64: pub,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	tokenString, err := token.SignedString(privateKey)
	fmt.Printf("%v %v\n", tokenString, err)

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return &privateKey.PublicKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["param"])
	} else {
		fmt.Println(err)
	}

	return tokenString, pub, err
}

// did create 요청 jwt 생성
func AddService(privateKey *ecdsa.PrivateKey) (string, string, error) {

	// publickey to base64 string
	publicKeyECDSA, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	// Create the Claims
	pub := Base64FromECDSAPub(publicKeyECDSA)
	claims := AddServiceJwt{
		jwt.StandardClaims{
			Issuer: "test",
		},
		Service{
			Id:              "did:pet:issuer",
			Name:            "유기견확인",
			ServiceEndpoint: "https://localhost/alone",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	tokenString, err := token.SignedString(privateKey)
	fmt.Printf("%v %v\n", tokenString, err)

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return &privateKey.PublicKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["param"])
	} else {
		fmt.Println(err)
	}

	return tokenString, pub, err
}

func TestCreate(t *testing.T) {
	c := new(PetDIDChaincode)
	stub := shim.NewMockStub("petdid", c)

	// pk 생성
	pk, err := createPrivateKey()
	if err != nil {
		t.Fail()
	}

	// Setting data to invoke
	funcName := "create"
	jwt, pub, err := Create(pk)
	if err != nil {
		t.Fail()
	}

	// Invoke
	checkInit(t, stub)
	res := checkInvoke(t, stub, funcName, []byte(jwt), []byte(pub))
	did := string(res.Payload)

	// query
	funcName = "read"
	checkQuery(t, stub, funcName, []byte(did))
}

func TestAddService(t *testing.T) {
	c := new(PetDIDChaincode)
	stub := shim.NewMockStub("petdid", c)

	// pk 생성
	pk, err := createPrivateKey()
	if err != nil {
		t.Fail()
	}

	// Setting data to invoke
	funcName := "create"
	jwt, pub, err := Create(pk)
	if err != nil {
		t.Fail()
	}

	// Invoke
	checkInit(t, stub)
	res := checkInvoke(t, stub, funcName, []byte(jwt), []byte(pub))
	did := string(res.Payload)

	// AddService
	funcName = "addService"
	jwt, pub, err = AddService(pk)
	if err != nil {
		t.Fail()
	}
	checkInvoke(t, stub, funcName, []byte(did), []byte(jwt))

	// query
	funcName = "read"
	checkQuery(t, stub, funcName, []byte(did))
}




