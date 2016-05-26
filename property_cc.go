package main

import (
	"errors"
	"fmt"
	//"github.com/hyperledger/fabric/core/chaincode/shim"
	"/tmp/cache/go1.6.2/go/src/github.com/hyperledger/fabric/core/chaincode/shim"
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

type BeneficialOwner struct {
	Name    string `json:"name"`
	Percent string `json:"percent"`
}

//==============================================================================================================================
//	Property - Defines the details for a Property object. JSON on right tells it what JSON fields to map to
//			  that element when reading a JSON object into the struct e.g. JSON make -> Struct Make.
//==============================================================================================================================
type Property struct {
	Folio_ID         string `json:"folio_id"`
	LegalOwner       string `json:"legalOwner"`
	BeneficialOwners []BeneficialOwner `json:"beneficialOwners"`
	Address          string `json:"address"`
	Status           int    `json:"status"`
}

type AllProperties struct {
	Properties []Property `json:"properties"`
}

type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func (c *PropertyChainCode) responseObject(function string, respMessage string, repStatus string) ([]byte, error) {
	var res Response

	/*message := "\"Message\":\"" + respMessage + "\", "
	status := "\"Status\":" + string(repStatus)
	jsonResp := "{" + message + status + "}"*/

	res.Status = repStatus
	res.Message = respMessage

	/*err := json.Unmarshal([]byte(jsonResp), &res)
	if err != nil {
		fmt.Println("UnmarshalError Invalid JSON object on responseObject().!! ", err)
	}*/
	respAsBytes, _ := json.Marshal(res)
	if repStatus != "0" {
		// Status != 0 is error case
		fmt.Printf(function + " function ERROR-ResponseJSON:\n", string(respAsBytes))
		return nil, errors.New(string(respAsBytes))
	} else {
		// Status == 0 is non error case
		fmt.Printf(function + " function ResponseJSON:\n", string(respAsBytes))
		return respAsBytes, nil
	}
}

func (c *PropertyChainCode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	return c.responseObject("Init", "Successfully Initialized", "0")
}

//==============================================================================================================================
// saveProperty - Writes to the ledger the Vehicle struct passed in a JSON format. Uses the shim file's
//				  method 'PutState'.
//==============================================================================================================================
func (c *PropertyChainCode) saveProperty(stub *shim.ChaincodeStub, p Property) ([]byte, error) {
	bytes, err := json.Marshal(p)
	if err != nil {
		fmt.Println("MarshalError while registering property", err)
		c.responseObject("saveProperty", "MarshalError while registering property", "99")
	}

	err = stub.PutState(p.Folio_ID, bytes)

	if err != nil {
		fmt.Println("Error while registering property", err)
		c.responseObject("saveProperty", "Error while registering property", "99")
	}

	allPropAsBytes, err := stub.GetState("allProps")
	var props AllProperties
	json.Unmarshal(allPropAsBytes, &props)

	//fmt.Printf("Query Response (GET ALL PROPS SIZE-BEFORE):\n", len(props.Properties))
	props.Properties = append(props.Properties, p)
	//fmt.Printf("Query Response (GET ALL PROPS SIZE-AFTER):\n", len(props.Properties))

	jsonAsBytes, _ := json.Marshal(props)
	err = stub.PutState("allProps", jsonAsBytes)
	if err != nil {
		fmt.Println("Error while PutState for allProps ", err)
		c.responseObject("saveProperty", "Error while PutState for allProps", "99")
	}

	return c.responseObject("saveProperty", "Successfully property registered", "0")
}

func (c *PropertyChainCode) Register(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	// Variables to define the JSON
	var p Property

	address := "\"Address\":\"" + args[0] + "\", "
	folio_ID := "\"Folio_ID\":\"" + args[1] + "\", "
	legalOwner := "\"LegalOwner\":\"" + args[2] + "\", "
	//beneficialOwners := "\"BeneficialOwners\":\"" + bo + "\", "
	status := "\"Status\":0"

	//fmt.Println("*** Calling Register()- Property args[3]:%s\n", args[3]);

	// Concatenates the variables to create the total JSON object
	//property_json := "{" + folio_ID + legalOwner + beneficialOwners + address + status + "}"
	property_json := "{" + folio_ID + legalOwner + address + status + "}"

	fmt.Println("*** Calling Register()- Property JSON:%s\n", property_json);


	// matched = true if the folio_ID passed fits format of "1/12345"
	matched, err := regexp.Match("^[0-9]{1}[/.][0-9]{5}$", []byte(args[1]))
	//fmt.Println("*** Calling Register()- Property Is it matched:%s\n", matched);

	if err != nil {
		//return nil, errors.New("Invalid Folio Identifier")
		return c.responseObject("Register", "Invalid Folio Identifier", "99")
	}

	if folio_ID == "" ||
	matched == false {
		//return nil, errors.New("Not-valid Folio Identifier value")
		return c.responseObject("Register", "Not-valid Folio Identifier value", "99")
	}

	// Convert the JSON defined above into a vehicle object for go
	err = json.Unmarshal([]byte(property_json), &p)

	// BeneficialOwners formation as input request
	for i := 3; i < len(args); {
		/*
				fmt.Println("################ index:", i)
				fmt.Println("################ NAME:", args[i])
				fmt.Println("################ Percent:", args[i + 1])*/
		if (args[i] != "") &&  (args[i + 1] != "") {
			var bo BeneficialOwner
			bo.Name = args[i]
			bo.Percent = args[i + 1]
			p.BeneficialOwners = append(p.BeneficialOwners, bo)
			i = i + 2
		}
	}
	//fmt.Println("*** Calling Register()- Property *****:%s\n", p);

	if err != nil {
		fmt.Println("UnmarshalError Invalid JSON object..!! ", err)
		//return nil, errors.New("Invalid JSON object")
		return c.responseObject("Register", "Invalid JSON object..!!", "99")
	}

	// If not an error then a record exists so cant create a new Property with this folio_ID as it must be unique
	record, err := stub.GetState(p.Folio_ID)

	if record != nil {
		//return nil, errors.New("Property already exists")
		return c.responseObject("Register", "Property already exists..!!", "99")
	}

	return c.saveProperty(stub, p)
}

func (c *PropertyChainCode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	// Handle different functions
	if function == "init" {
		return c.Init(stub, "init", args)
	} else if function == "register" {
		return c.Register(stub, args)
	}
	return c.responseObject("Invoke", "Unrecognized function:" + function, "99")
}

// ============================================================================================================================
// Get Transactions for a specific Financial Institution (Inbound and Outbound)
// ============================================================================================================================
func (c *PropertyChainCode) getProperties(stub *shim.ChaincodeStub, searchType string, searchValue string) ([]byte, error) {

	var res AllProperties

	fmt.Println("Looking for Property searchType:" + searchType + " searchValue:" + searchValue);

	//get the AllProperties index
	allPropAsBytes, err := stub.GetState("allProps")
	if err != nil {
		return c.responseObject("Query", "Failed to get all Properties", "99")
	}

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
		case "BeneficialOwnerName":
			for j := range props.Properties[i].BeneficialOwners {
				if props.Properties[i].BeneficialOwners[j].Name == searchValue {
					res.Properties = append(res.Properties, props.Properties[i])
				}
			}
		default:
			fmt.Printf("unrecognized property searchType..!!")
			return c.responseObject("Query", "Unrecognized property searchType..!!", "99")
		}
	}
	resAsBytes, _ := json.Marshal(res)
	fmt.Printf("Search ResponseJSON:\n", string(resAsBytes))
	return resAsBytes, nil

}

func (c *PropertyChainCode) deleteProperty(stub *shim.ChaincodeStub, folidID string) ([]byte, error) {
	stub.DelState(folidID)

	//get the AllProperties index
	allPropAsBytes, err := stub.GetState("allProps")
	if err != nil {
		return c.responseObject("deleteProperty", "Failed to get all Properties", "99")
	}

	var props, newProps AllProperties
	json.Unmarshal(allPropAsBytes, &props)

	for i := range props.Properties {
		if props.Properties[i].Folio_ID != folidID {
			newProps.Properties = append(newProps.Properties, props.Properties[i])
		}
	}

	jsonAsBytes, _ := json.Marshal(newProps)
	err = stub.PutState("allProps", jsonAsBytes)
	if err != nil {
		fmt.Println("Error while PutState for allProps ", err)
		return c.responseObject("deleteProperty", "Error while PutState for allProps", "99")
	}

	return c.responseObject("deleteProperty", "Successfully deleted the property", "0")

}


//=================================================================================================================================
//	Query - Called on PropertyChainCode query. Takes a function name passed and calls that function. Passes the
//  		initial arguments passed are passed on to the called function.
//=================================================================================================================================
func (c *PropertyChainCode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if function == "search" {
		if len(args) != 2 {
			return c.responseObject("search", "Incorrect number of arguments passed", "99")
		} else {
			return c.getProperties(stub, args[0], args[1])
		}
	} else if function == "delete" {
		if len(args) != 1 {
			return c.responseObject("delete", "Incorrect number of arguments passed", "99")
		} else {
			return c.deleteProperty(stub, args[0])
		}
	}
	return c.responseObject("Query", "Unrecognized function: " + function, "99")
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

