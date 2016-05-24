package main

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"regexp"
	"encoding/json"
)


//==============================================================================================================================
//	 Structure Definitions
//==============================================================================================================================
//	PropertyChainCode - A blank struct for use with Shim
//      (A HyperLedger included go file used for get/put state and other HyperLedger functions)
//==============================================================================================================================
type  PropertyChainCode struct {

}


//==============================================================================================================================
//	Property - Defines the details for a Property object. JSON on right tells it what JSON fields to map to
//			  that element when reading a JSON object into the struct e.g. JSON make -> Struct Make.
//==============================================================================================================================
type Property struct {
	Folio_ID        string `json:"folio_id"`
	LegalOwner      string `json:"legalOwner"`
	BeneficialOwner string `json:"beneficialOwner"`
	Address         string `json:"address"`
	Status          int    `json:"status"`
	/*Suburb    string `json:"suburb"`
	State     string `json:"state"`
	Postcode  string `json:"postcode"`*/
	//Image    string `json:"image"`
}

type AllProperties struct {
	Properties []Property `json:"properties"`
}

func (t*PropertyChainCode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	return nil, nil
}


//==============================================================================================================================
// saveProperty - Writes to the ledger the Vehicle struct passed in a JSON format. Uses the shim file's
//				  method 'PutState'.
//==============================================================================================================================
func (t *PropertyChainCode) saveProperty(stub *shim.ChaincodeStub, p Property) (bool, error) {
	fmt.Println("*** Calling saveProperty() ****");

	bytes, err := json.Marshal(p)
	if err != nil {
		return false, errors.New("Error creating Property record")
	}

	err = stub.PutState(p.Folio_ID, bytes)

	if err != nil {
		fmt.Println("Error while registering property ")
		return false, err
	}

	allPropAsBytes, err := stub.GetState("allProps")
	var props AllProperties
	json.Unmarshal(allPropAsBytes, &props)

	fmt.Printf("Query Response (GET ALL PROPS SIZE-BEFORE):\n", len(props.Properties))
	props.Properties = append(props.Properties, p)

	fmt.Printf("Query Response (GET ALL PROPS SIZE-AFTER):\n", len(props.Properties))

	jsonAsBytes, _ := json.Marshal(props)
	err = stub.PutState("allProps", jsonAsBytes)
	if err != nil {
		fmt.Println("Error while PutState for allProps ")
		return false, err
	}

	newAllPropAsBytes, err := stub.GetState("allProps")
	if err != nil {
		return false, errors.New("Failed to get all Properties")
	}
	fmt.Printf("Query Response (GET ALL PROPS ***):\n", string(newAllPropAsBytes))

	return true, nil
}

func (t *PropertyChainCode) Register(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	fmt.Println("*** Calling Register() ****");

	// Variables to define the JSON
	var p Property

	address := "\"Address\":\"" + args[0] + "\", "
	folio_ID := "\"Folio_ID\":\"" + args[1] + "\", "
	legalOwner := "\"LegalOwner\":\"" + args[2] + "\", "
	beneficialOwner := "\"BeneficialOwner\":\"" + args[3] + "\", "
	/*suburb := "\"Suburb\":\"" + args[3] + "\", "
	state := "\"State\":\"" + args[4] + "\", "
	postCode := "\"Postcode\":\"" + args[5] + "\", "*/
	status := "\"Status\":0"

	// Concatenates the variables to create the total JSON object
	//property_json := "{" + folio_ID + legalOwner + address + suburb + state + postCode + status + "}"
	property_json := "{" + folio_ID + legalOwner + beneficialOwner + address + status + "}"

	//fmt.Println("*** Calling Register()- Property JSON:%s\n", property_json);

	// matched = true if the v5cID passed fits format of two letters followed by seven digits
	// regexp.Match("^[A-z][A-z][0-9]{7}", []byte(folio_ID))
	matched, err := regexp.Match("[0-9]{1}[/.][0-9]{5}", []byte(folio_ID))
	//fmt.Println("*** Calling Register()- Property Is it matched:%s\n", matched);

	if err != nil {
		return nil, errors.New("Invalid Folio Identifier")
	}

	if folio_ID == "" ||
	matched == false {
		return nil, errors.New("Not-valid Folio Identifier value")
	}

	// Convert the JSON defined above into a vehicle object for go
	err = json.Unmarshal([]byte(property_json), &p)

	if err != nil {
		return nil, errors.New("Invalid JSON object")
	}

	// If not an error then a record exists so cant create a new Property with this folio_ID as it must be unique
	record, err := stub.GetState(p.Folio_ID)

	if record != nil {
		return nil, errors.New("Property already exists")
	}

	_, err = t.saveProperty(stub, p)

	fmt.Println("*** Calling Register()- Property saved");
	if err != nil {
		return nil, err
	}

	resp_mssg := "\"mssg\": Property Registered, "
	resp_status := "\"Status\":0"

	jsonResp := "{" + resp_mssg + resp_status + "}"

	regResAsBytes, _ := json.Marshal(jsonResp)
	fmt.Printf("Register ResponseJSON:\n", string(regResAsBytes))
	return regResAsBytes, nil
}

func (t *PropertyChainCode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "register" {
		return t.Register(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation")
}

// ============================================================================================================================
// Get Transactions for a specific Financial Institution (Inbound and Outbound)
// ============================================================================================================================
func (t *PropertyChainCode) getProperties(stub *shim.ChaincodeStub, searchType string, searchValue string) ([]byte, error) {

	var res AllProperties

	fmt.Println("Start find getProperties..!")
	fmt.Println("Looking for Property searchType:" + searchType + " searchValue:" + searchValue);

	//get the AllProperties index
	allPropAsBytes, err := stub.GetState("allProps")
	if err != nil {
		return nil, errors.New("Failed to get all Properties")
	}
	//fmt.Printf("Query Response (SAI TEST ***):\n", string(allPropAsBytes))

	var props AllProperties
	json.Unmarshal(allPropAsBytes, &props)

	for i := range props.Properties {

		switch searchType {
		case "ALL":
			res.Properties = append(res.Properties, props.Properties[i])
		case "Folio_ID":
			if props.Properties[i].Folio_ID == searchValue {
				res.Properties = append(res.Properties, props.Properties[i])
			}
		case "Address":
			if props.Properties[i].Address == searchValue {
				res.Properties = append(res.Properties, props.Properties[i])
			}
		case "LegalOwner":
			if props.Properties[i].LegalOwner == searchValue {
				res.Properties = append(res.Properties, props.Properties[i])
			}
		case "BeneficialOwner":
			if props.Properties[i].BeneficialOwner == searchValue {
				res.Properties = append(res.Properties, props.Properties[i])
			}
		default:
			fmt.Printf("unrecognized property searchType..!!")
		}
	}
	resAsBytes, _ := json.Marshal(res)
	fmt.Printf("Search ResponseJSON:\n", string(resAsBytes))
	return resAsBytes, nil

}

//=================================================================================================================================
//	Query - Called on PropertyChainCode query. Takes a function name passed and calls that function. Passes the
//  		initial arguments passed are passed on to the called function.
//=================================================================================================================================
func (t *PropertyChainCode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("*** Calling Query() ****");

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments passed")
	}

	if function == "search" {
		return t.getProperties(stub, args[0], args[1])
	}

	/*if function != "query" {
		return nil, errors.New("Invalid query function name. Expecting \"query\"")
	}

	var Folio_ID string // Entities
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the person to query")
	}

	Folio_ID = args[0]

	// Get the state from the ledger
	valAsbytes, err := stub.GetState(Folio_ID)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + Folio_ID + "\"}"
		return nil, errors.New(jsonResp)
	}

	if valAsbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + Folio_ID + "\"}"
		return nil, errors.New(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + Folio_ID + "\",\"Amount\":\"" + string(valAsbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	//return jsonResp, nil

	keysIter, err := stub.RangeQueryState("", "")
	if err != nil {
		return nil, fmt.Errorf("keys operation failed. Error accessing state: %s", err)
	}
	defer keysIter.Close()

	var keys []string
	for keysIter.HasNext() {
		key, _, err := keysIter.Next()
		if err != nil {
			return nil, fmt.Errorf("keys operation failed. Error accessing state: %s", err)
		}
		keys = append(keys, key)
	}

	jsonKeys, err := json.Marshal(keys)
	if err != nil {
		return nil, fmt.Errorf("keys operation failed. Error marshaling JSON: %s", err)
	}

	fmt.Printf("Query Response-jsonKeys:%s\n", jsonKeys)
	return jsonKeys, nil
	*/

	return nil, nil

}

//=================================================================================================================================
//	 Main - main - Starts up the chaincode
//=================================================================================================================================
func main() {
	err := shim.Start(new(PropertyChainCode))
	if err != nil {
		fmt.Printf("Error starting PropertyChainCode: %s", err)
	}
}

