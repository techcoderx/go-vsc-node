// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package gqlgen

import (
	ledgerDb "vsc-node/modules/db/vsc/ledger"
	"vsc-node/modules/db/vsc/transactions"
	"vsc-node/modules/gql/model"
)

type ContractDiff struct {
	Diff                    *string `json:"diff,omitempty"`
	PreviousContractStateID string  `json:"previousContractStateId"`
}

type ContractOutputFilter struct {
	ByID       *string `json:"byId,omitempty"`
	ByInput    *string `json:"byInput,omitempty"`
	ByContract *string `json:"byContract,omitempty"`
	Offset     *int    `json:"offset,omitempty"`
	Limit      *int    `json:"limit,omitempty"`
}

type FindContractFilter struct {
	ByID   *string `json:"byId,omitempty"`
	ByCode *string `json:"byCode,omitempty"`
	Offset *int    `json:"offset,omitempty"`
	Limit  *int    `json:"limit,omitempty"`
}

type Gas struct {
	Io *int `json:"IO,omitempty"`
}

type LedgerAction struct {
	ID     *string       `json:"id,omitempty"`
	Status *string       `json:"status,omitempty"`
	Type   *string       `json:"type,omitempty"`
	Data   model.Map     `json:"data,omitempty"`
	Asset  *string       `json:"asset,omitempty"`
	Amount *model.Uint64 `json:"amount,omitempty"`
	Memo   *string       `json:"memo,omitempty"`
	To     *string       `json:"to,omitempty"`
}

type LedgerActionsFilter struct {
	ByTxID     *string         `json:"byTxId,omitempty"`
	ByActionID *string         `json:"byActionId,omitempty"`
	ByAccount  *string         `json:"byAccount,omitempty"`
	ByTypes    []string        `json:"byTypes,omitempty"`
	ByAsset    *ledgerDb.Asset `json:"byAsset,omitempty"`
	ByStatus   *string         `json:"byStatus,omitempty"`
	FromBlock  *model.Uint64   `json:"fromBlock,omitempty"`
	ToBlock    *model.Uint64   `json:"toBlock,omitempty"`
	Offset     *int            `json:"offset,omitempty"`
	Limit      *int            `json:"limit,omitempty"`
}

type LedgerTxFilter struct {
	ByToFrom  *string         `json:"byToFrom,omitempty"`
	ByTxID    *string         `json:"byTxId,omitempty"`
	ByTypes   []string        `json:"byTypes,omitempty"`
	ByAsset   *ledgerDb.Asset `json:"byAsset,omitempty"`
	FromBlock *model.Uint64   `json:"fromBlock,omitempty"`
	ToBlock   *model.Uint64   `json:"toBlock,omitempty"`
	Offset    *int            `json:"offset,omitempty"`
	Limit     *int            `json:"limit,omitempty"`
}

type LocalNodeInfo struct {
	VersionID          string       `json:"version_id"`
	GitCommit          string       `json:"git_commit"`
	LastProcessedBlock model.Uint64 `json:"last_processed_block"`
	Epoch              model.Uint64 `json:"epoch"`
}

type Query struct {
}

type TransactionFilter struct {
	ByID           *string                         `json:"byId,omitempty"`
	ByIds          []string                        `json:"byIds,omitempty"`
	ByAccount      *string                         `json:"byAccount,omitempty"`
	ByContract     *string                         `json:"byContract,omitempty"`
	ByStatus       *transactions.TransactionStatus `json:"byStatus,omitempty"`
	ByType         []string                        `json:"byType,omitempty"`
	ByLedgerToFrom *string                         `json:"byLedgerToFrom,omitempty"`
	ByLedgerTypes  []string                        `json:"byLedgerTypes,omitempty"`
	Offset         *int                            `json:"offset,omitempty"`
	Limit          *int                            `json:"limit,omitempty"`
}

type TransactionSubmitResult struct {
	ID *string `json:"id,omitempty"`
}
