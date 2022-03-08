package helper

import (
	"fmt"
	"reflect"

	"github.com/merico-dev/lake/plugins/core"
)

// Accept raw json body and params, return list of entities that need to be stored
type RawDataExtractor func(row *RawData) ([]interface{}, error)

type ApiExtractorArgs struct {
	RawDataSubTaskArgs
	Params    interface{}
	RowData   interface{}
	Extract   RawDataExtractor
	BatchSize int
}

type ApiExtractor struct {
	*RawDataSubTask
	args *ApiExtractorArgs
}

func NewApiExtractor(args ApiExtractorArgs) (*ApiExtractor, error) {
	// process args
	rawDataSubTask, err := newRawDataSubTask(args.RawDataSubTaskArgs)
	if err != nil {
		return nil, err
	}
	if args.BatchSize == 0 {
		args.BatchSize = 500
	}
	return &ApiExtractor{
		RawDataSubTask: rawDataSubTask,
		args:           &args,
	}, nil
}

func (extractor *ApiExtractor) Execute() error {
	// load data from database
	db := extractor.args.Ctx.GetDb()
	cursor, err := db.Table(extractor.table).Order("id ASC").Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	row := &RawData{}

	// batch insertion divider
	RAW_DATA_ORIGIN := "RawDataOrigin"
	divider := NewBatchInsertDivider(db, extractor.args.BatchSize)
	divider.OnNewBatchInsert(func(rowType reflect.Type) error {
		// check if row type has RawDataOrigin
		if rawDataOrigin, ok := rowType.Elem().FieldByName(RAW_DATA_ORIGIN); ok {
			if (rawDataOrigin.Type != reflect.TypeOf(RawDataOrigin{})) {
				return fmt.Errorf("type %s must nested RawDataOrigin struct", rowType.Name())
			}
		} else {
			return fmt.Errorf("type %s must nested RawDataOrigin struct", rowType.Name())
		}
		// delete old data
		return db.Delete(
			reflect.New(rowType).Interface(),
			"_raw_data_table = ? AND _raw_data_params = ?",
			extractor.args.Table, extractor.params,
		).Error
	})

	// prgress
	extractor.args.Ctx.SetProgress(0, -1)
	// iterate all rows
	for cursor.Next() {
		err = db.ScanRows(cursor, row)
		if err != nil {
			return err
		}

		results, err := extractor.args.Extract(row)
		if err != nil {
			return err
		}

		for _, result := range results {
			// get the batch operator for the specific type
			batch, err := divider.ForType(reflect.TypeOf(result))
			if err != nil {
				return err
			}
			// set raw data origin field
			reflect.ValueOf(result).Elem().FieldByName(RAW_DATA_ORIGIN).Set(reflect.ValueOf(RawDataOrigin{
				RawDataTable:  extractor.args.Table,
				RawDataId:     row.ID,
				RawDataParams: row.Params,
			}))
			// records get saved into db when slots were max outed
			err = batch.Add(result)
			if err != nil {
				return err
			}
		}
		extractor.args.Ctx.IncProgress(1)
	}

	// save the last batches
	divider.Close()

	return nil
}

var _ core.SubTask = (*ApiExtractor)(nil)
