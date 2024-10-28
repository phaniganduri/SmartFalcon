package handlers

import (
	"fmt"
	"os"
	"log"
    "runtime"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

type AssetManagementAPI struct {
	Contract *gateway.Contract
}

func populateWallet(wallet *gateway.Wallet) error {
    credPath := filepath.Join(
        "..",
        "config",
    )

    certPath := filepath.Join(credPath, "org1.example.com-cert.pem")
    cert, err := os.ReadFile(filepath.Clean(certPath))
    if err != nil {
        return fmt.Errorf("failed to read cert file: %v", err)
    }

    keyPath := filepath.Join(credPath, "org1.example.com-key.pem")
    key, err := os.ReadFile(filepath.Clean(keyPath))
    if err != nil {
        return fmt.Errorf("failed to read key file: %v", err)
    }

    identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))
    err = wallet.Put("appUser", identity)
    if err != nil {
        return fmt.Errorf("failed to put identity in wallet: %v", err)
    }

    return nil
}

func SetupAPI() (*AssetManagementAPI, error) {
    err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
    if err != nil {
        return nil, fmt.Errorf("error setting DISCOVERY_AS_LOCALHOST environment variable: %v", err)
    }

    wallet, err := gateway.NewFileSystemWallet("wallet")
    if err != nil {
        return nil, fmt.Errorf("failed to create wallet: %v", err)
    }

    if !wallet.Exists("appUser") {
        err = populateWallet(wallet)
        if err != nil {
            return nil, fmt.Errorf("failed to populate wallet contents: %v", err)
        }
    }

    _, filename, _, ok := runtime.Caller(0)
    if !ok {
        return nil, fmt.Errorf("failed to get current file path")
    }
    currentDir := filepath.Dir(filename)

    // Construct the path to the config file
    configPath := filepath.Join(currentDir, "..", "config", "connection-org1.yaml")

    // Convert to absolute path
    absConfigPath, err := filepath.Abs(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to get absolute path: %v", err)
    }

    log.Printf("Looking for connection profile at: %s", absConfigPath)

    // Check if the file exists
    if _, err := os.Stat(absConfigPath); os.IsNotExist(err) {
        return nil, fmt.Errorf("connection profile does not exist at %s", absConfigPath)
    }

    gw, err := gateway.Connect(
        gateway.WithConfig(config.FromFile(filepath.Clean(absConfigPath))),
        gateway.WithIdentity(wallet, "appUser"),
    )
    if err != nil {
        return nil, fmt.Errorf("failed to connect to gateway: %v", err)
    }

    network, err := gw.GetNetwork("mychannel")
    if err != nil {
        return nil, fmt.Errorf("failed to get network: %v", err)
    }

    contract := network.GetContract("basic")

    api := &AssetManagementAPI{
        Contract: contract,
    }

    return api, nil
}