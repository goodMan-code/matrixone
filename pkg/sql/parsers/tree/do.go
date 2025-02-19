// Copyright 2021 Matrix Origin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tree

// Do statement
type Do struct {
	statementImpl
	Exprs []Expr
}

func (node *Do) Format(ctx *FmtCtx) {
	ctx.WriteString("do ")
	for _, e := range node.Exprs {
		e.Format(ctx)
	}
}

func (node *Do) GetStatementType() string { return "Do" }
func (node *Do) GetQueryType() string     { return QueryTypeDCL }
