package service

import (
	"context"

	"github.com/itering/subscan/model"
)

// GetTransferList gives us the list of transfers
func (s *Service) GetTransferList(page, row int, optionalParams map[string]string) ([]*model.TransferJson, int) {
	c := context.TODO()
	transfers, count := s.dao.GetTransferExtrinsicsList(c, page, row, optionalParams)
	return transfers, count
}
