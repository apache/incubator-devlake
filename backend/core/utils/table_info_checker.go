/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"go/ast"
	"go/parser"
	"go/token"
	"sort"
	"strings"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
)

type TableInfoChecker struct {
	tables  map[string]struct{}
	prefix  string
	ignores []string
}

func NewTableInfoChecker(prefix string, ignores []string) *TableInfoChecker {
	return &TableInfoChecker{tables: make(map[string]struct{}), prefix: prefix, ignores: ignores}
}

// The FeedIn function iterates through the model definitions,
// identifies all the tables by parsing the source files,
// and then compares them with the output obtained from the GetTablesInfo function.
func (checker *TableInfoChecker) FeedIn(modelsDir string, f func() []dal.Tabler) {
	set := token.NewFileSet()
	packs, err := parser.ParseDir(set, modelsDir, nil, 0)
	if err != nil {
		panic(err)
	}
	funcs := checker.getTableNameFuncs(packs)
	for _, fun := range funcs {
		s := checker.getTableName(fun)
		if strings.HasPrefix(s, checker.prefix) {
			checker.tables[s] = struct{}{}
		}
	}
	for _, tb := range checker.ignores {
		delete(checker.tables, tb)
	}
	for _, tabler := range f() {
		delete(checker.tables, tabler.TableName())
	}
}

func (checker *TableInfoChecker) Verify() errors.Error {
	for _, tb := range checker.ignores {
		delete(checker.tables, tb)
	}
	if len(checker.tables) == 0 {
		return nil
	}
	tableNames := make([]string, 0, len(checker.tables))
	sb := strings.Builder{}
	sb.WriteString("The following tables are not returned by the GetTablesInfo\n")
	for t := range checker.tables {
		tableNames = append(tableNames, t)
	}
	// sort the table names so that the tables of the same plugin are together
	sort.Strings(tableNames)
	sb.WriteString(strings.Join(tableNames, "\n"))
	return errors.Default.New(sb.String())
}

func (checker *TableInfoChecker) getTableNameFuncs(pks map[string]*ast.Package) []*ast.FuncDecl {
	var funcs []*ast.FuncDecl
	for _, pack := range pks {
		for _, f := range pack.Files {
			for _, d := range f.Decls {
				if fn, isFn := d.(*ast.FuncDecl); isFn && fn.Name.String() == "TableName" {
					funcs = append(funcs, fn)
				}
			}
		}
	}
	return funcs
}

func (checker *TableInfoChecker) getTableName(fn *ast.FuncDecl) string {
	for _, stmt := range fn.Body.List {
		if rs, isReturn := stmt.(*ast.ReturnStmt); isReturn {
			if lit, ok := rs.Results[0].(*ast.BasicLit); ok {
				return strings.Trim(lit.Value, `"`)
			}
		}
	}
	return ""
}
