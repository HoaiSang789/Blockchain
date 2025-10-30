package main

import (
    "fmt"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract định nghĩa cấu trúc chaincode
type SmartContract struct {
    contractapi.Contract
}

// InitLedger khởi tạo ledger với dữ liệu mẫu
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
    sampleData := map[string]string{
        "asset1": "value1",
        "asset2": "value2",
        "asset3": "value3",
    }

    for key, value := range sampleData {
        err := ctx.GetStub().PutState(key, []byte(value))
        if err != nil {
            return fmt.Errorf("Lỗi khởi tạo ledger: %s", err.Error())
        }
    }
    return nil
}

// Set thêm một cặp key-value vào ledger
func (s *SmartContract) Set(ctx contractapi.TransactionContextInterface, key string, value string) error {
    if key == "" || value == "" {
        return fmt.Errorf("Key và value không được để trống")
    }
    return ctx.GetStub().PutState(key, []byte(value))
}

// Get lấy giá trị của một key từ ledger
func (s *SmartContract) Get(ctx contractapi.TransactionContextInterface, key string) (string, error) {
    value, err := ctx.GetStub().GetState(key)
    if err != nil {
        return "", fmt.Errorf("Không thể lấy dữ liệu: %s", err.Error())
    }
    if value == nil {
        return "", fmt.Errorf("Key không tồn tại")
    }
    return string(value), nil
}

// Update cập nhật giá trị mới cho một key có sẵn
func (s *SmartContract) Update(ctx contractapi.TransactionContextInterface, key string, newValue string) error {
    exists, err := s.AssetExists(ctx, key)
    if err != nil {
        return err
    }
    if !exists {
        return fmt.Errorf("Không thể cập nhật: Key %s không tồn tại", key)
    }
    return ctx.GetStub().PutState(key, []byte(newValue))
}

// Delete xóa một key khỏi ledger
func (s *SmartContract) Delete(ctx contractapi.TransactionContextInterface, key string) error {
    if key == "" {
        return fmt.Errorf("Key không được để trống")
    }
    exists, err := s.AssetExists(ctx, key)
    if err != nil {
        return err
    }
    if !exists {
        return fmt.Errorf("Không thể xóa: Key %s không tồn tại", key)
    }
    return ctx.GetStub().DelState(key)
}

// AssetExists kiểm tra xem key có tồn tại trong ledger hay không
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, key string) (bool, error) {
    data, err := ctx.GetStub().GetState(key)
    if err != nil {
        return false, fmt.Errorf("Không thể kiểm tra key: %s", err.Error())
    }
    return data != nil, nil
}

// Hàm main để khởi động chaincode
func main() {
    chaincode, err := contractapi.NewChaincode(new(SmartContract))
    if err != nil {
        fmt.Printf("Lỗi tạo chaincode: %s", err.Error())
        return
    }

    if err := chaincode.Start(); err != nil {
        fmt.Printf("Lỗi khởi động chaincode: %s", err.Error())
    }
}
