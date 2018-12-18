package contractdb

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type Demo struct {


	ContractID string `json:"contract_id"`
	UserA string `json:"user_a"`
	UserB string `json:"user_b"`
	FileHash string `json:"file_hash"`
	ObjectType string `json:"object_type"`

}

//relize the interface,contract result
func (demo *Demo)  GetByteResult(stub shim.ChaincodeStubInterface, args []string) (string,string,string,interface{},error) {

	if len(args) != 5 {
		return "","","",nil,fmt.Errorf("paramter must be 5!")
	}

	demo.ContractID = args[0]
	demo.ObjectType = args[1]
	demo.UserA = args[2]
	demo.UserB = args[3]
	demo.FileHash = args[4]

	//set objectkey
	objectkey := "evidence"

	err := SetCompositekey(stub,objectkey,[]string{demo.UserA,demo.UserB})
	if err != nil {
		return "","","",nil,fmt.Errorf(err.Error())
	}

	return demo.ContractID,demo.ObjectType,objectkey,demo,nil

}
