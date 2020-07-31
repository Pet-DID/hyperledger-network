package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/sha3"
	"hash"
)

// DID Document 생성 요청 구조체
type CreateJwt struct {
	jwt.StandardClaims
	Param CreateParam `json:"param"`
}

type CreateParam struct {
	Id string `json:"id"`
	Type string `json:"type"`
	PublicKeyBase64 string `json:"publicKeyBase64"`
}

// Service 등록 요청 구조체
type AddServiceJwt struct {
	jwt.StandardClaims
	Param Service `json:"param"`
}

// jwt.MapClaims 를 CreateJwt 등의 struct 로 변환합니다
func convertToStruct(mapClaims jwt.MapClaims, v interface{}) error {
	fmt.Printf("mapClaims:%v\n", mapClaims)
	bb, err := json.Marshal(mapClaims)
	if err != nil {
		return err
	}
	return json.Unmarshal(bb, v)
}

func parse(t string, pub string) (*jwt.Token, error) {
	p, err := UnmarshalPubkey(pub)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return p, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func FromECDSAPub(pub *ecdsa.PublicKey) []byte {
	if pub == nil || pub.X == nil || pub.Y == nil {
		return nil
	}
	return elliptic.Marshal(elliptic.P256(), pub.X, pub.Y)
}

func Base64FromECDSAPub(pub *ecdsa.PublicKey) string {
	b := FromECDSAPub(pub)
	return base64.StdEncoding.EncodeToString(b)
}

// UnmarshalPubkey converts bytes to a secp256r1 public key (P-256)
func UnmarshalPubkey(pub string) (*ecdsa.PublicKey, error) {
	p, err := base64.StdEncoding.DecodeString(pub)
	if err != nil {
		return nil, err
	}

	x, y := elliptic.Unmarshal(elliptic.P256(), p)
	if x == nil {
		return nil, errors.New("invalid public key")
	}
	return &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}, nil
}

// pub -> sha3 -> 뒤에서부터 20 byte
func pubToAddress(pub *ecdsa.PublicKey) string {
	b := FromECDSAPub(pub)
	fmt.Printf("pub.len:%d\n", len(b))
	d := sha3256(b[1:])[12:]
	return hex.EncodeToString(d)
}

type ShaState interface {
	hash.Hash
	Read([]byte) (int, error)
}


func sha3256(data ...[]byte) []byte {
	b := make([]byte, 32)
	d := sha3.New256().(ShaState)
	for _, b := range data {
		d.Write(b)
	}
	d.Read(b)
	return b
}