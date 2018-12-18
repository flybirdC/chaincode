package contractdb

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"log"
	"strconv"
	"strings"
	"time"
)

type BusinessContract interface {

	//get bussiness byte result
	GetByteResult(stub shim.ChaincodeStubInterface, args []string) (string,string,[]byte,error)
}




//format contract
func PutStateContract(stub shim.ChaincodeStubInterface,  contractID string,objectType string,objectkey string,contractJsonBytes interface{}) error {

	//verity only ID
	id,_ := stub.GetState(contractID)
	if id == nil {
		stub.PutState(contractID,[]byte(contractID))
	} else {
		shim.Error("contractID has already setup!")
	}

	var contract Contract

	contract.ContractID = contractID
	contract.ObjectType = objectType
	contract.TimeStampID = time.Now().Unix()
	contract.ContractJson = contractJsonBytes

	contract.ObjectKey = objectkey

	putJsonBytes, err := json.Marshal(&contract)
	if err != nil {
		return fmt.Errorf("json seriallize fail while marshal!")
	}

	err = stub.PutState(contract.ContractID,putJsonBytes)
	if err != nil {
		return fmt.Errorf("fail to putstate while regist contractID")
	}

	err = stub.PutState(contract.ObjectType,putJsonBytes)
	if err != nil {
		return fmt.Errorf("fail to putstate while regist contract type")
	}

	err = stub.PutState(strconv.FormatInt(contract.TimeStampID,10),putJsonBytes)
	if err != nil {
		return fmt.Errorf("fail to putstate ")
	}

	return nil

}

//format modify contract
func ModifyStateContract(stub shim.ChaincodeStubInterface, contractID string,objectType string,objectkey string,contractJsonBytes interface{}) error{
	var contract Contract

	id ,err := stub.GetState(contractID)
	if err != nil {
		return fmt.Errorf("get contractID error!")
	}

	contract.ContractID = string(id)
	contract.ObjectType = objectType
	contract.TimeStampID = time.Now().Unix()
	contract.ContractJson = contractJsonBytes

	contract.ObjectKey = objectkey
	putJsonBytes, err := json.Marshal(&contract)
	if err != nil {
		return fmt.Errorf("json seriallize fail while marshal!")
	}

	err = stub.PutState(contract.ContractID,putJsonBytes)
	if err != nil {
		return fmt.Errorf("fail to putstate while regist contractID")
	}

	err = stub.PutState(contract.ObjectType,putJsonBytes)
	if err != nil {
		return fmt.Errorf("fail to putstate while regist contract type")
	}

	err = stub.PutState(strconv.FormatInt(contract.TimeStampID,10),putJsonBytes)
	if err != nil {
		return fmt.Errorf("fail to putstate ")
	}

	return nil

}


//format revoke contract
func RevokeStateContract(stub shim.ChaincodeStubInterface,contractID string)  error {

	bytes, err := stub.GetState(contractID)
	if err != nil {
		return fmt.Errorf("contractid not found!")
	}
	var contract Contract

	json.Unmarshal(bytes,contract)

	err = stub.DelState(contract.ObjectType)
	if err != nil {
		return fmt.Errorf("delete object state failed")
	}

	err = stub.DelState(strconv.FormatInt(contract.TimeStampID,10))
	if err != nil {
		return fmt.Errorf("delete timestamp state failed")
	}

	err = stub.DelState(string(contractID))
	if err != nil {
		return fmt.Errorf("delete contractID state failed!")
	}


	return nil
}








//------------------------------------------------------------------------------------------
//get contract member name
func GetCreatorName(stub shim.ChaincodeStubInterface) (string,error)  {

	name, err := GetCreator(stub)
	if err != nil {
		return "", nil
	}

	//format name
	memberName := name[(strings.Index(name,"@")+1):strings.LastIndex(name,".example.com")]
	return memberName, nil
}


//get creator
func GetCreator(stub shim.ChaincodeStubInterface) (string, error)  {

	creatorByte, _ := stub.GetCreator()
	certStart := bytes.IndexAny(creatorByte,"-----BEGIN")
	if certStart == -1 {
		fmt.Errorf("No cert found")
	}
	certText := creatorByte[certStart:]
	bl, _ := pem.Decode(certText)
	if bl == nil {
		fmt.Errorf("Could not decode the PEM structure!")
	}

	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		fmt.Errorf("parserCert failed!")
	}
	uname := cert.Subject.CommonName
	return uname, nil
}

//------------------------------------------------------------------------------------------------

//set compositekey
func SetCompositekey(stub shim.ChaincodeStubInterface,objectkey string,keys []string) error {

	key, err := stub.CreateCompositeKey(objectkey,keys)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	value := []byte{0x00}
	err = stub.PutState(key,value)

	if err != nil {
		return fmt.Errorf(err.Error())
	}

	return nil
}


//---------------------------------------query--------------------------------------------------------------------
//query by range k-v,data not include endkey, the key must be with state
func QueryByRange(stub shim.ChaincodeStubInterface,startkey string, endkey string) ([]byte,error)  {

	var bytebuffer bytes.Buffer

	//get the result data between startkey and endkey
	resultsIterator, err := stub.GetStateByRange(startkey,endkey)

	if err != nil {
		return nil,fmt.Errorf(err.Error())
	}

	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMenberAlreadyWritten := false

	//iterator data, and construct josn, return data result
	for resultsIterator.HasNext() {

		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil,fmt.Errorf(err.Error())
		}

		if bArrayMenberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResult.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		buffer.WriteString(string(queryResult.Value))
		buffer.WriteString("}")

		bArrayMenberAlreadyWritten =true
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


	return bytebuffer.Bytes(),nil
}

//query contract history  --> get json
func QueryByHistory(stub shim.ChaincodeStubInterface, key string) ([]byte,error) {

	//return the history key-value
	resultsIterator, err := stub.GetHistoryForKey(key)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResult.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(queryResult.Timestamp.Seconds,int64(queryResult.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		buffer.WriteString(string(queryResult.Value))


		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(queryResult.IsDelete))
		buffer.WriteString("\"}")

		bArrayMemberAlreadyWritten = true

	}

	buffer.WriteString("]")

	fmt.Printf("- getContractsByRange queryresult:\n%s\n",buffer.String())



	return buffer.Bytes(), nil
}


//query composite keys by order
func QuerybyCompositekeys(stub shim.ChaincodeStubInterface, args[] string) ([]byte,error)  {

	if len(args) < 1 {
		return nil,fmt.Errorf("no objectType, no query paramter")
	}
	objectKey := args[0]

	keys := append(args[:0],args[1:]...)


	var buffer bytes.Buffer

	resultIterator, err := stub.GetStateByPartialCompositeKey(objectKey,keys)
	defer resultIterator.Close()
	if err != nil {
		return nil,fmt.Errorf(err.Error())
	}
	buffer.WriteString(fmt.Sprintf("{\"object_key\":\"%s\",\"record\":{",objectKey,))
	i := 1

	bArrayMemberAlreadyWritten := false

	for resultIterator.HasNext() {

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		item, _ := resultIterator.Next()

		_,comkeyparts,err :=stub.SplitCompositeKey(item.Key)
		k := len(comkeyparts)
		if err != nil {
			return nil,fmt.Errorf(err.Error())
		}
		buffer.WriteString(fmt.Sprintf("\"%s%d\":\"",objectKey,i))
		for index, value:= range comkeyparts {
			buffer.WriteString(value)
			if index != k-1 {
				buffer.WriteString("-")
			} else {
				buffer.WriteString("\"")
			}
		}

		i++

		bArrayMemberAlreadyWritten = true



	}

	buffer.WriteString("}}")
	return buffer.Bytes(),nil

}


//--------------------------------------query end------------------------------------------------------------------


//date format

//check date and set query datekey
func CheckAndSetDateKey(stub shim.ChaincodeStubInterface,start string,end string)  (string,string) {

	var startstring, endstring string
	startdate,err := time.Parse("2006-01-01 01:02:03 AM",start)
	enddate,err := time.Parse("2006-01-01 01:02:03 AM",end)
	if err != nil{
		log.Fatalf("date error!")
	}

	if startdate.Unix()< 0 {
		log.Fatalf("date format error!")
	}
	indexS := startdate.Unix()
	indexE := enddate.Unix()

	for  {
		bytes , _ := stub.GetState(strconv.FormatInt(indexS,10))
		if bytes != nil {
			startstring = strconv.FormatInt(indexS,10)
			break
		}

		indexS++
		if indexS == indexE {
			return startstring,endstring
			break
		}
	}


	for  {
		bytes, _ := stub.GetState(strconv.FormatInt(indexE,10))
		if bytes != nil {
			endstring = strconv.FormatInt(indexE,10)
			return startstring,endstring
			break
		}

		indexE--
		if indexE==indexS {
			endstring = startstring
			break
		}

	}


	return startstring,endstring

}
