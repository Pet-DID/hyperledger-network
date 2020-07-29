# Hyperledger 네트워크를 이용한 애완동물 DID 서비스 구현
Hyperledger 블록체인 네트워크 상에서 Decentralized Identifiers(DIDs) 표준 Spec 에 따라
애완동물(강아지)의 DID 서비스를 구현함

## 시작하기

### 네트워크
다음 명령을 실행해서 네트워크를 시작할 수 있습니다.
```sh
./start.sh
```
네트워크 구성은 fabric-samples 의 `basic-network` 를 사용하고 있으며, 
Pet-DID 서비스를 이용하기 위한 체인코드 설치까지 진행합니다.

## 체인코드 interface
체인코드에서 invoke/query 할 수 있는 메소드는 다음과 같습니다.
```go
type Chaincode interface {
    // DID Document 생성
    Create(stub shim.ChaincodeStubInterface, jwt string) peer.Response

    // DID Document 조회
    Read(stub shim.ChaincodeStubInterface, did string) peer.Response

    // DID Doucment 에 서비스 정보 등록
    AddService(stub shim.ChaincodeStubInterface, jwt string) peer.Response
}
```




