docker exec cli peer chaincode install -p chaincodedev/chaincode/petdid -n petdid -v 1.0
docker exec cli peer chaincode instantiate -n petdid -v 1.0 -c '{"Args":[]}' -C myc