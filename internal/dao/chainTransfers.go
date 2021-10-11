package dao

import (
	"context"
	"fmt"

	"github.com/itering/subscan/model"
	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/address"
)

// GetTransferExtrinsicsList provides the list of transfers by looking at extrinsics
func (d *Dao) GetTransferExtrinsicsList(c context.Context, page, row int, optionalParams map[string]string) ([]model.TransferJson, int) {
	var extrinsics []model.ChainExtrinsic
	var count int

	blockNum, _ := d.GetFillBestBlockNum(c)
	for index := blockNum / model.SplitTableBlockNum; index >= 0; index-- {
		var tableData []model.ChainExtrinsic
		var tableCount int
		queryOrigin := d.db.Model(model.ChainExtrinsic{BlockNum: index * model.SplitTableBlockNum}).Where("call_module = ? AND  call_module_function = ?", "balances", "transfer")

		queryWhere := d.TransferParamBasedFilters(c, optionalParams)
		for _, w := range queryWhere {
			queryOrigin = queryOrigin.Where(w)
		}

		queryOrigin.Count(&tableCount)

		if tableCount == 0 {
			continue
		}
		preCount := count
		count += tableCount
		if len(extrinsics) >= row {
			continue
		}
		query := queryOrigin.Order("block_num desc").Offset(page*row - preCount).Limit(row - len(extrinsics)).Scan(&tableData)
		if query == nil || query.Error != nil || query.RecordNotFound() {
			continue
		}
		extrinsics = append(extrinsics, tableData...)

	}
	var transferList []model.TransferJson
	for _, extrinsic := range extrinsics {
		transferList = append(transferList, *d.TransferExtrinsicAsTransferJSON(c, &extrinsic))
	}

	return transferList, count
}

// TransferExtrinsicAsTransferJSON formats the extrinsic info for a transfer to that of a transfer type
func (d *Dao) TransferExtrinsicAsTransferJSON(c context.Context, e *model.ChainExtrinsic) *model.TransferJson {
	tj := &model.TransferJson{
		BlockNum:       e.BlockNum,
		BlockTimestamp: e.BlockTimestamp,
		ExtrinsicIndex: e.ExtrinsicIndex,
		From:           address.SS58Address(e.AccountId),
		Success:        e.Success,
		Fee:            e.Fee,
		Hash:           e.ExtrinsicHash,
		Module:         e.CallModule,
		Nounce:         e.Nonce,
	}
	var params []model.ExtrinsicParam
	util.UnmarshalAny(&params, e.Params)
	for _, param := range params {
		if param.Name == "value" && param.Type == "Compact<Balance>" {
			tj.Amount = util.DecimalFromInterface(param.Value)
		}
		if param.Name == "dest" && param.Type == "Address" {
			tj.To = address.SS58Address(param.Value.(string))
		}
	}
	return tj
}

// TransferParamBasedFilters defines addional query conditions for fetching transfer type extrinsic entries
func (d *Dao) TransferParamBasedFilters(c context.Context, optionalParams map[string]string) []string {

	var query []string
	if transferAddress, ok := optionalParams[util.TransferAddress]; ok {
		query = append(query, fmt.Sprintf("account_id = '%s'", transferAddress))
	}
	if transferFromBlock, ok := optionalParams[util.TransferFromBlock]; ok {
		query = append(query, fmt.Sprintf("block_num >= '%d'", util.StringToInt(transferFromBlock)))
	}
	if transferToBlock, ok := optionalParams[util.TransferToBlock]; ok {
		query = append(query, fmt.Sprintf("block_num < '%d'", util.StringToInt(transferToBlock)))
	}
	return query
}
