package main

import (
  "fmt"
  "encoding/json"
  "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// create property transfer contract

type PropertyTransferSmartContract struct {
	contractapi.Contract //contract api interface
}


//define property object 

type Property struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Area int `json:"area"`
	OwnerName string `json:"ownername"`
	Value int `json:"value"`
}

//define functions

//add property
func(pc *PropertyTransferSmartContract) AddProperty(ctx contractapi.TransactionContextInterface,id string,name string,area int,ownerName string,value int) error {

	propertyJSON,err := ctx.GetStub().GetState(id) //check if property exists in blockchain

	if err!=nil{
		return fmt.Errorf("Failed to read data from world state",err)
	}

	if propertyJSON!=nil {
		return fmt.Errorf("the property %s already exists",id)
	}

	prop:=Property{
		ID:id,
		Name:name,
		Area:area,
		OwnerName:ownerName,
		Value:value,
	}

	propertyBytes,err:=json.Marshal(prop) //convert to a variable and check if successful
	if err!=nil{
		return err
	}

	return ctx.GetStub().PutState(id,propertyBytes) // put data in blockchain
}

// get all properties

// GetAllAssets returns all assets found in world state
func (pc *PropertyTransferSmartContract) QueryAllProperties(ctx contractapi.TransactionContextInterface) ([]*Property, error) {
  // range query with empty string for startKey and endKey does an
  // open-ended query of all assets in the chaincode namespace.
  propertyIterator, err := ctx.GetStub().GetStateByRange("", "")
  if err != nil {
    return nil, err
  }
  defer propertyIterator.Close()

  var properties []*Property 

  for propertyIterator.HasNext() {
    propertyResponse, err := propertyIterator.Next()
    if err != nil {
      return nil, err
    }

    var property *Property
    err = json.Unmarshal(propertyResponse.Value, &property)
    if err != nil {
      return nil, err
    }
    properties = append(properties, property)
  }

  return properties, nil
}

// AssetExists returns true when asset with given ID exists in world state
func (pc *PropertyTransferSmartContract) QueryPropertyById(ctx contractapi.TransactionContextInterface, id string) (*Property, error) {
	propertyJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
	  return nil, fmt.Errorf("failed to read from world state: %v", err)
	}

	if propertyJSON == nil {
		return nil,fmt.Errorf("the property %s does not exist",id)
	}

	var property *Property
	err=json.Unmarshal(propertyJSON,&property)

	if err!=nil {
		return nil,err
	}
  
	return property, nil
  }
  
  // TransferAsset updates the owner field of asset with given id in world state.
func (pc *PropertyTransferSmartContract) transferProperty(ctx contractapi.TransactionContextInterface, id string, newOwner string) error {
	property, err := pc.QueryPropertyById(ctx, id)
	if err != nil {
	  return err
	}
  
	property.OwnerName = newOwner
	propertyJSON, err := json.Marshal(property)
	if err != nil {
	  return err
	}
  
	return ctx.GetStub().PutState(id, propertyJSON)
  }

  func main() {
	propTransferSmartContract := new(PropertyTransferSmartContract)

	cc,err:=contractapi.NewChaincode(propTransferSmartContract)
	if err != nil {
	  panic(err.Error())
	}
  
	if err := cc.Start(); err != nil {
	  panic(err.Error())
	}
  }