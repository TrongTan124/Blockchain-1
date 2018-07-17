package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// ProfileChaincode --
type ProfileChaincode struct {
}

// Profile --
type Profile struct {
	UserID  string   `json:"user_id"`
	Class10 Class    `json:"class_10"`
	Class11 Class    `json:"class_11"`
	Class12 Class    `json:"class_12"`
	BC      []string `json:"bc"`
}

// Class --
type Class struct {
	ClassName  string    `json:"class_name"`
	NameSchool string    `json:"name_school"`
	NameHT     string    `json:"name_HT"`
	NameGVCN   string    `json:"name_GVCN"`
	Subjects   []Subject `ison:"subjects"`
	HK         string    `json:"hk"`
	DH         []string  `json:"dh"`
}

// Subject --
type Subject struct {
	NameSubject  string `json:"name_subject"`
	ScoreSubject string `json:"score_subject"`
}

// Init ProfileChaincode
func (t *ProfileChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("user profile Init")
	return shim.Success(nil)
}

func (t *ProfileChaincode) initProfile(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	//   0          1      2
	// "userID", "class", "bc",
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	userID := args[0]
	classNew := args[1]
	bc := args[2]

	// ==== Check if user already exists ====
	userAsBytes, err := stub.GetState(userID)
	if err != nil {
		return shim.Error("Failed to get user: " + err.Error())
	} else if userAsBytes != nil {
		fmt.Println("This profile already exists: " + userID)
		return shim.Error("This profile already exists: " + userID)
	}

	var className string
	var nameSchool string
	var nameHT string
	var nameGVCN string
	var subjects string
	var hk string
	var dh string

	value := strings.Split(classNew, ",")

	className = value[0]
	nameSchool = value[1]
	nameHT = value[2]
	nameGVCN = value[3]
	subjects = value[4]
	hk = value[5]
	dh = value[6]

	var listSubjectNew []Subject

	listSubject := strings.Split(subjects, "&")

	for _, value := range listSubject {
		valueNew := strings.Split(value, "#")
		listSubjectNew = append(listSubjectNew, Subject{valueNew[0], valueNew[1]})
	}

	var dhNew []string

	for _, value := range strings.Split(dh, "#") {
		dhNew = append(dhNew, value)
	}

	class := Class{className, nameSchool, nameHT, nameGVCN, listSubjectNew, hk, dhNew}

	var classA Class
	var classB Class

	var bcNew []string

	for _, value := range strings.Split(bc, "#") {
		bcNew = append(bcNew, value)
	}

	// ==== Create profile object and marshal to JSON ====

	profile := &Profile{userID, class, classA, classB, bcNew}

	profileJSONasBytes, err := json.Marshal(profile)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Save user to state ===
	err = stub.PutState(userID, profileJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *ProfileChaincode) updateProfile(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//   0         1	   2
	// "userID", "class", "bc"
	if len(args) < 4 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	userID := args[0]
	classNew := args[1]
	bc := args[2]
	level := args[3]

	fmt.Println("- start updateProfile ", userID)

	profileAsBytes, err := stub.GetState(userID)
	if err != nil {
		return shim.Error("Failed to get user:" + err.Error())
	} else if profileAsBytes == nil {
		return shim.Error("User does not exist")
	}

	profileOld := &Profile{}

	err = json.Unmarshal(profileAsBytes, &profileOld) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}

	var className string
	var nameSchool string
	var nameHT string
	var nameGVCN string
	var subjects string
	var hk string
	var dh string

	value := strings.Split(classNew, ",")

	className = value[0]
	nameSchool = value[1]
	nameHT = value[2]
	nameGVCN = value[3]
	subjects = value[4]
	hk = value[5]
	dh = value[6]

	var listSubjectNew []Subject

	listSubject := strings.Split(subjects, "&")

	for _, value := range listSubject {
		valueNew := strings.Split(value, "$")
		listSubjectNew = append(listSubjectNew, Subject{valueNew[0], valueNew[1]})
	}

	var dhNew []string

	for _, value := range strings.Split(dh, "#") {
		dhNew = append(dhNew, value)
	}

	class := Class{className, nameSchool, nameHT, nameGVCN, listSubjectNew, hk, dhNew}
	var bcNew []string

	for _, value := range strings.Split(bc, "#") {
		bcNew = append(bcNew, value)
	}

	var profileNew *Profile

	if level == "10" {
		profileNew = &Profile{userID, class, profileOld.Class11, profileOld.Class12, bcNew}
	} else if level == "11" {
		profileNew = &Profile{userID, profileOld.Class10, class, profileOld.Class12, bcNew}
	} else {
		profileNew = &Profile{userID, profileOld.Class10, profileOld.Class11, class, bcNew}
	}

	userJSONasBytes, _ := json.Marshal(profileNew)
	err = stub.PutState(userID, userJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end success")
	return shim.Success(nil)
}

func (t *ProfileChaincode) deleteProfile(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var jsonResp string
	var profileJSON Profile
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	userID := args[0]

	// to maintain the color~name index, we need to read the marble first and get its color
	valAsbytes, err := stub.GetState(userID) //get the marble from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + userID + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Profile does not exist: " + userID + "\"}"
		return shim.Error(jsonResp)
	}

	err = json.Unmarshal([]byte(valAsbytes), &profileJSON)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to decode JSON of: " + userID + "\"}"
		return shim.Error(jsonResp)
	}

	err = stub.DelState(userID) //remove the marble from chaincode state
	if err != nil {
		return shim.Error("Failed to delete state:" + err.Error())
	}
	return shim.Success(nil)
}

func (t *ProfileChaincode) getProfileByID(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	userID := args[0]

	queryString := fmt.Sprintf("{\"selector\":{\"user_id\":\"%s\"}}", userID)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

func (t *ProfileChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("user profile Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "getProfileByID" {
		// get user by id
		return t.getProfileByID(stub, args)
	} else if function == "deleteProfile" {
		// Delete user
		return t.deleteProfile(stub, args)
	} else if function == "updateProfile" {
		// update user
		return t.updateProfile(stub, args)
	} else if function == "initProfile" {
		// create new profile
		return t.initProfile(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\"")
}

func main() {
	err := shim.Start(new(ProfileChaincode))
	if err != nil {
		fmt.Printf("Error starting profile chaincode: %s", err)
	}
}
