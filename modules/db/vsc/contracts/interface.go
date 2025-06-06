package contracts

import a "vsc-node/modules/aggregate"

type Contracts interface {
	a.Plugin
	RegisterContract(contractId string, args Contract)
	ContractById(contractId string) (Contract, error)
	FindContracts(contractId *string, code *string, offset int, limit int) ([]Contract, error)
}

type ContractState interface {
	a.Plugin
	IngestOutput(inputArgs IngestOutputArgs)
	GetLastOutput(contractId string, height uint64) *ContractOutput
	FindOutputs(id *string, input *string, contract *string, offset int, limit int) ([]ContractOutput, error)
}

type IngestOutputArgs struct {
	Id          string
	ContractId  string
	StateMerkle string

	Inputs  []string
	Results []ContractOutputResult `bson:"results"`

	Metadata map[string]interface{} `bson:"metadata"`

	AnchoredBlock  string
	AnchoredHeight int64
	AnchoredId     string
	AnchoredIndex  int64
}

type ContractOutputResult struct {
	Ret string `json:"ret" bson:"ret"`
	Ok  bool   `json:"ok" bson:"ok"`
}

type ContractOutput struct {
	Id          string                 `json:"id"`
	BlockHeight int64                  `json:"block_height" bson:"block_height"`
	Timestamp   *string                `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
	ContractId  string                 `json:"contract_id" bson:"contract_id"`
	Inputs      []string               `json:"inputs"`
	Metadata    map[string]interface{} `json:"metadata"`
	//This might not be used

	Results     []ContractOutputResult `json:"results" bson:"results"`
	StateMerkle string                 `json:"state_merkle" bson:"state_merkle"`
}
