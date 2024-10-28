package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// AssetManagement provides functions for managing assets
type AssetManagement struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up an asset
type Asset struct {
	DEALERID    string  `json:"dealerId"`
	MSISDN      string  `json:"msisdn"`
	MPIN        string  `json:"mpin"`
	BALANCE     float64 `json:"balance"`
	STATUS      string  `json:"status"`
	TRANSAMOUNT float64 `json:"transAmount"`
	TRANSTYPE   string  `json:"transType"`
	REMARKS     string  `json:"remarks"`
}

// InitLedger adds a base set of assets to the ledger
func (s *AssetManagement) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []Asset{
		{DEALERID: "dealer1", MSISDN: "1234567890", MPIN: "1234", BALANCE: 1000.0, STATUS: "active"},
		{DEALERID: "dealer2", MSISDN: "0987654321", MPIN: "4321", BALANCE: 2000.0, STATUS: "active"},
	}

	for _, asset := range assets {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(asset.MSISDN, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// CreateAsset issues a new asset to the world state with given details
func (s *AssetManagement) CreateAsset(ctx contractapi.TransactionContextInterface, dealerId string, msisdn string, mpin string, balance float64) error {
	exists, err := s.AssetExists(ctx, msisdn)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", msisdn)
	}

	asset := Asset{
		DEALERID: dealerId,
		MSISDN:   msisdn,
		MPIN:     mpin,
		BALANCE:  balance,
		STATUS:   "active",
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(msisdn, assetJSON)
}

// ReadAsset returns the asset stored in the world state with given id
func (s *AssetManagement) ReadAsset(ctx contractapi.TransactionContextInterface, msisdn string) (*Asset, error) {
	assetJSON, err := ctx.GetStub().GetState(msisdn)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", msisdn)
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters
func (s *AssetManagement) UpdateAsset(ctx contractapi.TransactionContextInterface, msisdn string, newBalance float64, newStatus string) error {
	asset, err := s.ReadAsset(ctx, msisdn)
	if err != nil {
		return err
	}

	asset.BALANCE = newBalance
	asset.STATUS = newStatus

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(msisdn, assetJSON)
}

// DeleteAsset deletes an given asset from the world state
func (s *AssetManagement) DeleteAsset(ctx contractapi.TransactionContextInterface, msisdn string) error {
	exists, err := s.AssetExists(ctx, msisdn)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", msisdn)
	}

	return ctx.GetStub().DelState(msisdn)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *AssetManagement) AssetExists(ctx contractapi.TransactionContextInterface, msisdn string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(msisdn)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// GetAllAssets returns all assets found in world state
func (s *AssetManagement) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	// range query with empty string for startKey and endKey does an open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

// GetAssetHistory returns the chain of custody for an asset since issuance
func (s *AssetManagement) GetAssetHistory(ctx contractapi.TransactionContextInterface, msisdn string) ([]Asset, error) {
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(msisdn)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var history []Asset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		history = append(history, asset)
	}

	return history, nil
}

func main() {
	assetChaincode, err := contractapi.NewChaincode(&AssetManagement{})
	if err != nil {
		fmt.Printf("Error creating asset-management chaincode: %s", err.Error())
		return
	}

	if err := assetChaincode.Start(); err != nil {
		fmt.Printf("Error starting asset-management chaincode: %s", err.Error())
	}
}