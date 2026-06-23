package request

import (
	"wallet_chain.com/admin/model"
)

type SysOperationRecordSearch struct {
	model.SysOperationRecord
	PageInfo
}
