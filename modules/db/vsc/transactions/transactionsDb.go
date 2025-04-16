package transactions

import (
	"context"
	"strings"
	"time"
	"vsc-node/modules/db"
	"vsc-node/modules/db/vsc"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type transactions struct {
	*db.Collection
}

func New(d *vsc.VscDb) Transactions {
	return &transactions{db.NewCollection(d.DbInstance, "transaction_pool")}
}

func (e *transactions) Init() error {
	err := e.Collection.Init()
	if err != nil {
		return err
	}

	return nil
}

func (e *transactions) Ingest(offTx IngestTransactionUpdate) error {
	ctx := context.Background()

	queryy := bson.M{
		"id": offTx.Id,
	}

	findResult := e.FindOne(ctx, bson.M{
		"id": offTx.Id,
	})

	opts := options.Update().SetUpsert(true)
	setOp := bson.M{
		"anchr_block":    offTx.AnchoredBlock,
		"anchr_id":       offTx.AnchoredId,
		"anchr_height":   offTx.AnchoredHeight,
		"anchr_index":    offTx.AnchoredIndex,
		"anchr_opidx":    offTx.AnchoredOpIdx,
		"type":           offTx.Type,
		"data":           offTx.Tx,
		"required_auths": offTx.RequiredAuths,
		"nonce":          offTx.Nonce,
		"rc_limit":       offTx.RcLimit,
	}
	if findResult.Err() != nil {
		setOp["first_seen"] = time.Now()
		//Prevents case of reprocessing/reindexing
		if offTx.Status != "" {
			setOp["status"] = offTx.Status
		} else {
			setOp["status"] = "UNCONFIRMED"
		}
	} else {
		//If it already exists do nothing
		if offTx.Status != "" {
			setOp["status"] = offTx.Status
		}
	}
	_, err := e.UpdateOne(ctx, queryy, bson.M{
		"$set": setOp,
	}, opts)

	return err
}

func (e *transactions) SetOutput(sOut SetResultUpdate) {
	query := bson.M{
		"id": sOut.Id,
	}
	ctx := context.Background()

	update := bson.M{}

	if sOut.Output != nil {
		update["output"] = sOut.Output
	}
	if sOut.Ledger != nil {
		update["ledger"] = sOut.Ledger
	}

	e.FindOneAndUpdate(ctx, query, bson.M{
		"$set": update,
	})
}

func (e *transactions) GetTransaction(id string) *TransactionRecord {
	query := bson.M{
		"id": id,
	}
	ctx := context.Background()
	findResult := e.FindOne(ctx, query)

	if findResult.Err() != nil {
		return nil
	}
	record := TransactionRecord{}
	err := findResult.Decode(&record)
	if err != nil {
		return nil
	}
	return &record
}

func (e *transactions) FindTransactions(id *string, account *string, contract *string, status *string, byType *string, ledgerToFrom *string, ledgerTypes []string, offset int, limit int) ([]TransactionRecord, error) {
	filters := bson.D{}
	if id != nil {
		filters = append(filters, bson.E{Key: "id", Value: *id})
	}
	if account != nil {
		filters = append(filters, bson.E{Key: "required_auths", Value: *account})
	}
	if contract != nil {
		filters = append(filters, bson.E{Key: "data.contract_id", Value: *contract})
	}
	if status != nil {
		filters = append(filters, bson.E{Key: "status", Value: strings.ToUpper(*status)})
	}
	if byType != nil {
		filters = append(filters, bson.E{Key: "data.type", Value: *byType})
	}
	if ledgerToFrom != nil {
		filters = append(filters, bson.E{Key: "$or", Value: bson.A{
			bson.D{{Key: "ledger.from", Value: *ledgerToFrom}},
			bson.D{{Key: "ledger.to", Value: *ledgerToFrom}},
		}})
	}
	if len(ledgerTypes) > 0 {
		ledgerTypeFilter := bson.A{}
		for _, t := range ledgerTypes {
			ledgerTypeFilter = append(ledgerTypeFilter, bson.D{{Key: "ledger.type", Value: t}})
		}
		filters = append(filters, bson.E{Key: "$or", Value: ledgerTypeFilter})
	}
	pipe := mongo.Pipeline{
		{{Key: "$match", Value: filters}},
		// Join with hive_blocks
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "hive_blocks"},
			{Key: "localField", Value: "anchr_height"},
			{Key: "foreignField", Value: "block.block_number"},
			{Key: "as", Value: "block_info"},
		}}},
		// Unwind the joined array
		{{Key: "$unwind", Value: "$block_info"}},
		// Add timestamp field
		{{Key: "$addFields", Value: bson.D{
			{Key: "anchr_ts", Value: "$block_info.block.timestamp"},
		}}},
		// Remove temporary field
		{{Key: "$project", Value: bson.D{
			{Key: "block_info", Value: 0},
		}}},
		// Sorting
		{{Key: "$sort", Value: bson.D{{Key: "anchr_height", Value: -1}}}},
		// Pagination
		{{Key: "$skip", Value: offset}},
		{{Key: "$limit", Value: limit}},
	}
	cursor, err := e.Aggregate(context.TODO(), pipe)
	if err != nil {
		return []TransactionRecord{}, err
	}
	defer cursor.Close(context.TODO())
	var results []TransactionRecord
	for cursor.Next(context.TODO()) {
		var elem TransactionRecord
		if err := cursor.Decode(&elem); err != nil {
			return []TransactionRecord{}, err
		}
		results = append(results, elem)
	}
	return results, nil
}

// Searches for unconfirmed VSC transactions with no verification
// Provide height for expiration filtering
func (e *transactions) FindUnconfirmedTransactions(height uint64) ([]TransactionRecord, error) {
	query := bson.M{
		"status": "UNCONFIRMED",
		"type":   "vsc",
		"$or": bson.A{
			bson.M{
				"expire_block": bson.M{
					"$exists": false,
				},
			},
			bson.M{
				"expire_block": bson.M{
					"$gt": height,
				},
			},
			bson.M{
				"expire_block": bson.M{
					"$eq": nil,
				},
			},
		},
	}

	ctx := context.Background()
	findResult, _ := e.Find(ctx, query)

	txList := make([]TransactionRecord, 0)
	for findResult.Next(ctx) {
		tx := TransactionRecord{}
		err := findResult.Decode(&tx)

		if err != nil {
			return nil, err
		}
		txList = append(txList, tx)
	}

	return txList, nil
}

// SetStatus of all IDs and ID + Opidx to a specific status
func (e *transactions) SetStatus(ids []string, status string) {

	for _, id := range ids {

		oneResult := e.FindOne(context.Background(), bson.M{
			"id": id,
		})
		var result TransactionRecord
		err := oneResult.Decode(&result)

		//Transaction not indexed (for some reason!)
		if err != nil {
			continue
		}

		e.UpdateMany(context.Background(), bson.M{
			"anchr_height": result.AnchoredHeight,
			"anchr_opidx":  result.AnchoredOpIdx,
		}, bson.M{
			"$set": bson.M{
				"status": status,
			},
		})
	}
}
