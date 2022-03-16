package helper

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/merico-dev/lake/plugins/core"
	"gorm.io/datatypes"
)

// Table structure for raw data storage
type RawData struct {
	ID        uint64 `gorm:"primaryKey"`
	Params    string `gorm:"type:varchar(255);index"`
	Data      datatypes.JSON
	Url       string
	Input     datatypes.JSON
	CreatedAt time.Time
}

// embedded fields for tool layer tables
type RawDataOrigin struct {
	// can be used for flushing outdated records from table
	RawDataParams string `gorm:"column:_raw_data_params;type:varchar(255);index"`
	RawDataTable  string `gorm:"column:_raw_data_table"`
	// can be used for debugging
	RawDataId uint64 `gorm:"column:_raw_data_id"`
	// we can store record index into this field, which is helpful for debugging
	RawDataRemark string `gorm:"column:_raw_data_remark"`
}

type RawDataSubTaskArgs struct {
	Ctx    core.SubTaskContext
	Table  string      `comment:"Raw data table name"`
	Params interface{} `comment:"To identify a set of records with same UrlTemplate, i.e. {SourceId, BoardId} for jira entities"`
}

// Common features for raw data sub tasks
type RawDataSubTask struct {
	args   *RawDataSubTaskArgs
	table  string
	params string
}

func newRawDataSubTask(args RawDataSubTaskArgs) (*RawDataSubTask, error) {
	if args.Ctx == nil {
		return nil, fmt.Errorf("Ctx is required for RawDataSubTask")
	}
	if args.Table == "" {
		return nil, fmt.Errorf("Table is required for RawDataSubTask")
	}
	paramsString := ""
	if args.Params == nil {
		args.Ctx.GetLogger().Warn("Missing `Params` for raw data subtask %s", args.Ctx.GetName())
	} else {
		// TODO: maybe sort it to make it consisitence
		paramsBytes, err := json.Marshal(args.Params)
		if err != nil {
			return nil, err
		}
		paramsString = string(paramsBytes)
	}
	return &RawDataSubTask{
		args:   &args,
		table:  fmt.Sprintf("_raw_%s", args.Table),
		params: paramsString,
	}, nil
}
