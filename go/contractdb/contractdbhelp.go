package contractdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"strings"
)

//parameter:contractID and business args, putstate contract,ID must be new
func (contract *Contract) InvokeContract(stub shim.ChaincodeStubInterface, args []string) peer.Response  {



	//demo
	var demo Demo
	contractID, contractType,objectkey,bytesjson,err := demo.GetByteResult(stub,args)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = PutStateContract(stub,contractID,contractType,objectkey,bytesjson)


	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("invoke contract success!"))
}



//parameter:contractID and business args, modify contract state,the first paramter must ID!
func (contract *Contract) ModifyContract(stub shim.ChaincodeStubInterface,args []string) peer.Response   {


	if len(args) != 1 {
		return shim.Error("paramter must be contract ID")
	}

	//demo
	var demo Demo
	contractID, contractType, objectkey,bytesjson,err := demo.GetByteResult(stub,args)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = ModifyStateContract(stub,contractID,contractType,objectkey,bytesjson)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("modify contract success!"))
}




//paramter:contractID, delete contract state

func (contract *Contract) DeleteContract(stub shim.ChaincodeStubInterface, args []string) peer.Response  {

	if len(args) != 1 {
		return shim.Error("only input contractID")
	}
	id := args[0]

	err := RevokeStateContract(stub,id)

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}


//---------------couchdb-------------------------------------------------------------------------

//set couchdb query key
func (contract *Contract) SetQuerykey(stub shim.ChaincodeStubInterface,args []string) peer.Response  {

	if len(args) != 1 {

		return shim.Error("key must string and be 1")
	}

	queryrichkey := args[0]

	stub.PutState("richkey",[]byte(queryrichkey))

	return shim.Success([]byte("set rich key success!"))
}

//set couchdb query keys
func (Contract *Contract) SetQuerykeys(stub shim.ChaincodeStubInterface,args []string) peer.Response {

	if len(args) <2 {
		return shim.Error("couchdb more keys,set more than one key!")
	}

	//strings.Replace(strings.Trim(fmt.Sprint(array_or_slice), "[]"), " ", ",", -1)
	keystring := strings.Replace(strings.Trim(fmt.Sprint(args),"[]")," ",",",-1)

	err := stub.PutState("richkeys",[]byte(keystring))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("richkeys put state success!"))
}


//couchdb rich query,only a key
func (contract *Contract)queryrichkey(stub shim.ChaincodeStubInterface, args []string) peer.Response  {

	if len(args) != 1 {
		return shim.Error("rich parameter must be 1")
	}

	value:= args[0]
	key,err := stub.GetState("richkey")
	if err != nil || key == nil {
		return shim.Error(err.Error())
	}

	queryStr := fmt.Sprintf("{\"selector\":{\"%s\":\"%s\"}}",key,value)

	resultsIterator ,err := stub.GetQueryResult(queryStr)
	if err != nil {
		return shim.Error(err.Error())
	}

	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResult.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		buffer.WriteString(string(queryResult.Value))
		buffer.WriteString("}")

		bArrayMemberAlreadyWritten = true
	}

	buffer.WriteString("]")

	fmt.Printf("- getContractUser queryResult:\n%s\n",buffer.String())
	return shim.Success(buffer.Bytes())

}

//couchdb rich query, keys
/*
{
   "selector": {
      "_id": {
         "$gt": null
      },
      "image_name": {
         "$gt": null
      },
      "image_user": {
         "$gt": null
      }
   }
}
 */
func (contract *Contract)queryrichkeys(stub shim.ChaincodeStubInterface,args []string) peer.Response  {

	if len(args) < 1 {
		return shim.Error("parameter nil, error!")
	}

	key,err := stub.GetState("richkeys")
	if err != nil {
		return shim.Error(err.Error())
	}

	keystring := string(key)
	keys := strings.Split(keystring,",")

	if len(keys) != len(args) {
		return shim.Error("key value not pair, value paramter not complete!")
	}

	var querybuffer bytes.Buffer
	querybuffer.WriteString("{\"selector\":{\"_id\":{\"$gt\":null},")

	bArrayMemberAlreadyWritten := false

	for index, key := range keys {

		if args[index] == "" && index == len(keys)-1 {
			bArrayMemberAlreadyWritten = false
		}


		if bArrayMemberAlreadyWritten == true {
			querybuffer.WriteString(",")
		}

		if args[index] != "" {
			querybuffer.WriteString(fmt.Sprintf("\"%s\":\"%s\"",key,args[index]))
			bArrayMemberAlreadyWritten = true

		} else {
			bArrayMemberAlreadyWritten = false
		}


	}

	querybuffer.WriteString("}}")

	resultsIterator ,err := stub.GetQueryResult(querybuffer.String())
	if err != nil {
		return shim.Error(err.Error())
	}

	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten = false
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResult.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		buffer.WriteString(string(queryResult.Value))
		buffer.WriteString("}")

		bArrayMemberAlreadyWritten = true
	}

	buffer.WriteString("]")

	fmt.Printf("- getContractUser queryResult:\n%s\n",buffer.String())

	return shim.Success(buffer.Bytes())

}



//--------------couchdb end-----------------------------------------------------------------------


//----------------query---------------------------------------------------------------------------


//paramter:query key:contractID, query contract state
func (contract *Contract) QueryContractState(stub shim.ChaincodeStubInterface, args []string) peer.Response  {

	var contractID string
	var err error

	if len(args) != 1 {
		return shim.Error("arg must be username!")
	}

	contractID = args[0]
	content, err:= stub.GetState(contractID)

	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + contractID + "\"}"
		return shim.Error(jsonResp)
	}
	if content == nil {
		jsonResp := "{\"Error\":\"Nil record for " + contractID + "\"}"
		return shim.Error(jsonResp)
	}

	json.Unmarshal(content,contract)


	jsonResp := "{\"ID\":\"" + contractID + "\",\"RecordNowState\":\"" + string(content) +"}"
	fmt.Printf("Query Response:%s\n",jsonResp)

	return shim.Success([]byte(jsonResp))

}


//paramter:objectType , query by objectType, return all contracts
func (contract *Contract) querybyContractType(stub shim.ChaincodeStubInterface, args []string) peer.Response  {

	if len(args) != 1 {
		return shim.Error("Contract type must be one")
	}

	result, err := QueryByHistory(stub,args[0])

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(result)
}

//paramter:args[0]=objectkey,args[1],args[2]...=key, query key by order,first key not be nil
func (contract *Contract)querybyComkey(stub shim.ChaincodeStubInterface, args []string) peer.Response  {

	bytes,err := QuerybyCompositekeys(stub,args)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(bytes)

}


//paramter:start date, end date, query by date,return contracts state between start date and end date
//the date fomat:2018-12-16 01:02:03 AM
// --> 12 hours
//consider double time history state, must return time history contract data:key--> getStateHistory(key)
func (contract *Contract) querybyContractDate(stub shim.ChaincodeStubInterface,args []string) peer.Response  {
	if len(args) != 2 {
		return shim.Error("parameter error! start date and end date")
	}

	startkey,endkey := CheckAndSetDateKey(stub,args[0],args[1])

	resultsIterator, err := stub.GetStateByRange(startkey,endkey)
	if err != nil {
		return shim.Error(err.Error())
	}

	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false

	for resultsIterator.HasNext() {

		queryResult, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResult.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		buffer.WriteString(string(queryResult.Value))
		buffer.WriteString("}")

		bArrayMemberAlreadyWritten =true
	}

	//include endkey
	buffer.WriteString(",")

	buffer.WriteString("{\"Key\":")
	buffer.WriteString("\"")
	buffer.WriteString(endkey)
	buffer.WriteString("\"")

	buffer.WriteString(", \"Record\":")
	endvalue, _ := stub.GetState(endkey)
	buffer.WriteString(string(endvalue))
	buffer.WriteString("}")


	buffer.WriteString("]")


	fmt.Printf("- getContractByRange queryResult:\n%s\n",buffer.String())

	return shim.Success(buffer.Bytes())

}






//---------------query end------------------------------------------------------------------------