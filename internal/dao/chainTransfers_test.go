package dao

import (
	"context"
	"testing"

	"github.com/itering/subscan/model"
	"github.com/itering/subscan/util"
	"github.com/stretchr/testify/assert"
)

func TestDao_GetTransferList(t *testing.T) {
	ctx := context.TODO()
	optionalParams := map[string]string{}
	transfers, _ := testDao.GetTransferExtrinsicsList(ctx, 0, 100, optionalParams)
	t.Logf("transfers:%+v", transfers)
	expectedTransfers := []*model.TransferJson{
		{
			From:           "DPmDGJjQfnHgfeVrR29wrS5fttauHnrQALimWZFszQZnJTm",
			To:             "EXPn2gLnQ7h5uhbdoYWBRRcdeSzzJnzH7XEq3YaCdtNTfvo",
			Module:         "balances",
			Amount:         util.DecimalFromInterface("1000000000000000000"),
			Hash:           "0x368f61800f8645f67d59baf0602b236ff47952097dcaef3aa026b50ddc8dcea0",
			BlockTimestamp: 1594791900,
			BlockNum:       947689,
			ExtrinsicIndex: "947689-1",
			Success:        true,
			Nounce:         0,
			Fee:            util.DecimalFromInterface("0"),
		},
	}
	assert.Equal(t, expectedTransfers, transfers)
}

func TestDao_TransferExtrinsicAsTransferJson(t *testing.T) {
	ctx := context.TODO()
	transferExtrinsic := &testSignedExtrinsic
	expectedTransferJSON := &model.TransferJson{
		From:           "DPmDGJjQfnHgfeVrR29wrS5fttauHnrQALimWZFszQZnJTm",
		To:             "EXPn2gLnQ7h5uhbdoYWBRRcdeSzzJnzH7XEq3YaCdtNTfvo",
		Module:         "balances",
		Amount:         util.DecimalFromInterface("1000000000000000000"),
		Hash:           "0x368f61800f8645f67d59baf0602b236ff47952097dcaef3aa026b50ddc8dcea0",
		BlockTimestamp: 1594791900,
		BlockNum:       947689,
		ExtrinsicIndex: "947689-1",
		Success:        true,
		Nounce:         0,
	}
	assert.Equal(t, expectedTransferJSON, testDao.TransferExtrinsicAsTransferJSON(ctx, transferExtrinsic))
}
