package helper

import (
	"fmt"
	"reflect"

	"github.com/merico-dev/lake/plugins/core"
	"gorm.io/gorm"
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

	// cache for batch insertion
	cache := make(map[reflect.Type]*subCache)

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
			resultType := reflect.TypeOf(result)
			// get the cache for the specific type
			resultCache := cache[resultType]
			// create one if not exists
			if resultCache == nil {
				resultCache, err = newSubCache(extractor.args.Ctx.GetDb(), resultType, extractor.args.BatchSize)
				if err != nil {
					return err
				}
				cache[resultType] = resultCache
			}
			// records get saved into db when slots were max outed
			err = resultCache.Add(result)
			if err != nil {
				return err
			}
		}
	}

	// close all caches so all rest records get saved into db as well
	for _, resultCache := range cache {
		resultCache.Close()
	}

	return nil
}

var _ core.SubTask = (*ApiExtractor)(nil)

type subCache struct {
	slotType reflect.Type
	// slots can not be []interface{}, because gorm wouldn't take it
	// I'm guessing the reason is the type information lost when converted to interface{}
	slots   reflect.Value
	db      *gorm.DB
	current int
	size    int
}

func newSubCache(db *gorm.DB, slotType reflect.Type, size int) (*subCache, error) {
	return &subCache{
		slotType: slotType,
		slots:    reflect.MakeSlice(reflect.SliceOf(slotType), size, size),
		db:       db,
		size:     size,
	}, nil
}

func (c *subCache) Add(slot interface{}) error {
	// type checking
	if reflect.TypeOf(slot) != c.slotType {
		return fmt.Errorf("sub cache type mismatched")
	}
	if reflect.ValueOf(slot).Kind() != reflect.Ptr {
		return fmt.Errorf("slot is not a pointer")
	}
	//return c.db.Create(slot).Error

	// push into slot
	c.slots.Index(c.current).Set(reflect.ValueOf(slot))
	c.current++
	// flush out into database if max outed
	if c.current == c.size {
		return c.Flush()
	}
	return nil
}

func (c *subCache) Flush() error {
	err := c.db.CreateInBatches(c.slots.Slice(0, c.current).Interface(), c.current).Error
	c.current = 0
	if err != nil {
		return err
	}
	return nil
}

func (c *subCache) Close() error {
	if c.current > 0 {
		return c.Flush()
	}
	return nil
}
