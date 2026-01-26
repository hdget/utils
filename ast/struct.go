package ast

import (
	"go/ast"
	"go/token"
)

// GetStructInfo 获取结构信息
func GetStructInfo(n *ast.GenDecl) (string, *ast.StructType, bool) {
	// 仅处理类型声明
	if n.Tok == token.TYPE {
		for _, spec := range n.Specs {
			// 如果类型规范是类型别名或类型声明
			if typeSpec, ok := spec.(*ast.TypeSpec); ok {
				// 如果类型规范是结构体类型
				if st, ok := typeSpec.Type.(*ast.StructType); ok {
					return typeSpec.Name.Name, st, true
				}
			}
		}
	}
	return "", nil, false
}
