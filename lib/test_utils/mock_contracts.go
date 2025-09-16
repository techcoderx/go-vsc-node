package test_utils

import (
	"fmt"
	"vsc-node/modules/aggregate"
	"vsc-node/modules/db/vsc/contracts"
)

type MockContractDb struct {
	aggregate.Plugin
	Contracts map[string]contracts.Contract
}

func (m *MockContractDb) RegisterContract(contractId string, args contracts.Contract) {
	m.Contracts[contractId] = args
}

func (m *MockContractDb) ContractById(contractId string) (contracts.Contract, error) {
	info, exists := m.Contracts[contractId]
	if !exists {
		return contracts.Contract{}, fmt.Errorf("contract %s does not exist", contractId)
	}
	return info, nil
}

// GraphQL use only, not implemented in mocks
func (m *MockContractDb) FindContracts(contractId *string, code *string, offset int, limit int) ([]contracts.Contract, error) {
	return []contracts.Contract{}, nil
}

type MockContractStateDb struct {
	aggregate.Plugin
	Outputs map[string]contracts.ContractOutput
}

func (m *MockContractStateDb) IngestOutput(inputArgs contracts.IngestOutputArgs) {
	// Convert IngestOutputArgs to ContractOutput
	output := contracts.ContractOutput{
		Id:          inputArgs.Id,
		ContractId:  inputArgs.ContractId,
		StateMerkle: inputArgs.StateMerkle,
		BlockHeight: inputArgs.AnchoredHeight,
		Metadata:    inputArgs.Metadata,
		Inputs:      inputArgs.Inputs,
		Results:     inputArgs.Results,
	}

	m.Outputs[output.Id] = output
}

func (m *MockContractStateDb) GetLastOutput(contractId string, height uint64) (contracts.ContractOutput, error) {
	var lastOutput contracts.ContractOutput
	found := false

	for _, output := range m.Outputs {
		if output.ContractId == contractId && uint64(output.BlockHeight) <= height {
			if !found || output.BlockHeight > lastOutput.BlockHeight {
				lastOutput = output
				found = true
			}
		}
	}

	if !found {
		return contracts.ContractOutput{}, nil // Return empty output and nil error to match original implementation
	}

	return lastOutput, nil
}

func (m *MockContractStateDb) GetOutput(outputId string) *contracts.ContractOutput {
	result := m.Outputs[outputId]
	return &result
}

// GraphQL use only, not implemented in mocks
func (m *MockContractStateDb) FindOutputs(id *string, input *string, contract *string, offset int, limit int) ([]contracts.ContractOutput, error) {
	return []contracts.ContractOutput{}, nil
}
