package helper

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/apache/incubator-devlake/models/common"

	"github.com/apache/incubator-devlake/plugins/core"
)

// Accept row from source cursor, return list of entities that need to be stored
type DataConvertHandler func(row interface{}) ([]interface{}, error)

type DataConverterArgs struct {
	RawDataSubTaskArgs
	// Domain layer entity Id prefix, i.e. `jira:JiraIssue:1`, `github:GithubIssue`
	InputRowType reflect.Type
	// Cursor to a set of Tool Layer Records
	Input     *sql.Rows
	Convert   DataConvertHandler
	BatchSize int
}

// DataConverter helps you convert Data from Tool Layer Tables to Domain Layer Tables
// It reads rows from specified Iterator, and feed it into `Convter` handler
// you can return arbitrary domain layer entities from this handler, ApiConverter would
// first delete old data by their RawDataOrigin information, and then perform a
// batch save operation for you.
type DataConverter struct {
	*RawDataSubTask
	args *DataConverterArgs
}

func NewDataConverter(args DataConverterArgs) (*DataConverter, error) {
	rawDataSubTask, err := newRawDataSubTask(args.RawDataSubTaskArgs)
	if err != nil {
		return nil, err
	}
	// process args
	if args.BatchSize == 0 {
		args.BatchSize = 500
	}
	return &DataConverter{
		RawDataSubTask: rawDataSubTask,
		args:           &args,
	}, nil
}

func (converter *DataConverter) Execute() error {
	// load data from database
	db := converter.args.Ctx.GetDb()

	// batch insertion divider
	RAW_DATA_ORIGIN := "RawDataOrigin"
	divider := NewBatchSaveDivider(db, converter.args.BatchSize)
	divider.OnNewBatchSave(func(rowType reflect.Type) error {
		// check if row type has RawDataOrigin
		if rawDataOrigin, ok := rowType.Elem().FieldByName(RAW_DATA_ORIGIN); ok {
			if (rawDataOrigin.Type != reflect.TypeOf(common.RawDataOrigin{})) {
				return fmt.Errorf("type %s must nested RawDataOrigin struct", rowType.Name())
			}
		} else {
			return fmt.Errorf("type %s must nested RawDataOrigin struct", rowType.Name())
		}
		// delete old data
		return db.Delete(
			reflect.New(rowType).Interface(),
			"_raw_data_table = ? AND _raw_data_params = ?",
			converter.table, converter.params,
		).Error
	})

	// prgress
	converter.args.Ctx.SetProgress(0, -1)

	cursor := converter.args.Input
	defer cursor.Close()
	ctx := converter.args.Ctx.GetContext()
	// iterate all rows
	for cursor.Next() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		inputRow := reflect.New(converter.args.InputRowType).Interface()
		err := db.ScanRows(cursor, inputRow)
		if err != nil {
			return err
		}

		results, err := converter.args.Convert(inputRow)
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
			reflect.ValueOf(result).Elem().FieldByName(RAW_DATA_ORIGIN).
				Set(reflect.ValueOf(inputRow).Elem().FieldByName(RAW_DATA_ORIGIN))
			// records get saved into db when slots were max outed
			err = batch.Add(result)
			if err != nil {
				return err
			}
		}
		converter.args.Ctx.IncProgress(1)
	}

	// save the last batches
	return divider.Close()
}

var _ core.SubTask = (*DataConverter)(nil)
