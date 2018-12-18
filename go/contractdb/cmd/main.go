package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/examples/chaincode/go/contractdb"
)

func main() {
	err := shim.Start(new(contractdb.Contract))
	if err != nil {
		fmt.Errorf("contract chaincode start error = %s",err)
	}
}
