package helper

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/merico-dev/lake/plugins/core"
)

// sql query and parameters to select the same batch of data, i.e. all issues that came from
// the same jira board
type BatchSelector struct {
	Query      string
	Parameters []interface{}
}

// Accept row from source cursor, return list of entities that need to be stored
type DataConvertHandler func(row interface{}) ([]interface{}, error)

type DataConverterArgs struct {
	Ctx core.SubTaskContext
	// Domain layer entity Id prefix, i.e. `jira:JiraIssue:1`, `github:GithubIssue`
	InputRowType reflect.Type
	// Cursor to a set of Tool Layer Records
	Input          *sql.Rows
	Convert        DataConvertHandler
	BatchSelectors map[reflect.Type]BatchSelector
	BatchSize      int
}

// DataConverter helps you convert Data from Tool Layer Tables to Domain Layer Tables
// It reads rows from specified Iterator, and feed it into `Convter` handler
// you can return arbitrary domain layer entities from this handler, ApiConverter would
// first delete old data by their RawDataOrigin information, and then perform a
// batch save operation for you.
type DataConverter struct {
	args *DataConverterArgs
}

func NewDataConverter(args DataConverterArgs) (*DataConverter, error) {
	// process args
	if args.BatchSize == 0 {
		args.BatchSize = 500
	}
	return &DataConverter{
		args: &args,
	}, nil
}

func (converter *DataConverter) Execute() error {
	// load data from database
	db := converter.args.Ctx.GetDb()

	inputRow := reflect.New(converter.args.InputRowType).Interface()
	// batch insertion divider
	divider := NewBatchSaveDivider(db, converter.args.BatchSize)
	divider.OnNewBatchSave(func(rowType reflect.Type) error {
		sel, ok := converter.args.BatchSelectors[rowType]
		if !ok {
			return fmt.Errorf("must provide BatchSelector for type %s in order to clean up old data", rowType)
		}
		// delete old data
		return db.Delete(
			reflect.New(rowType).Interface(),
			sel.Query,
			sel.Parameters,
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
