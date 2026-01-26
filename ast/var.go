package ast

import (
	"go/ast"
	"go/token"
)

// GetVarDeclsFromFunc 从函数体中提取所有 *ast.ValueSpec
// 提取函数体内所有ValueSpec（包括转换短声明）
func GetVarDeclsFromFunc(body *ast.BlockStmt) map[string]*ast.ValueSpec {
	results := make(map[string]*ast.ValueSpec)

	for _, stmt := range body.List {
		switch node := stmt.(type) {
		case *ast.DeclStmt:
			// 处理var块声明（包含多个ValueSpec）
			if genDecl, ok := node.Decl.(*ast.GenDecl); ok && genDecl.Tok == token.VAR {
				for _, spec := range genDecl.Specs {
					if valueSpec, ok := spec.(*ast.ValueSpec); ok {
						for _, name := range valueSpec.Names {
							results[name.Name] = valueSpec
						}
					}
				}
			}

		case *ast.AssignStmt:
			// 将短声明（:=）转换为ValueSpec
			if node.Tok == token.DEFINE {
				for i, lhs := range node.Lhs {
					if i >= len(node.Rhs) {
						break
					}
					if ident, ok := lhs.(*ast.Ident); ok {
						results[ident.Name] = &ast.ValueSpec{
							Names:  []*ast.Ident{ident},
							Values: []ast.Expr{node.Rhs[i]},
						}
					}
				}
			}
		}
	}
	return results
}

// GetVarDeclsFromFile 查找变量声明
func GetVarDeclsFromFile(file *ast.File) map[string]*ast.ValueSpec {
	results := make(map[string]*ast.ValueSpec)
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.VAR {
			continue
		}
		for _, spec := range genDecl.Specs {
			if valueSpec, ok := spec.(*ast.ValueSpec); ok {
				for _, name := range valueSpec.Names {
					results[name.Name] = valueSpec
				}
			}
		}
	}
	return results
}

// GetVarTypes 收集所有变量声明中的类型信息
func GetVarTypes(node ast.Node) map[string]string {
	results := make(map[string]string)
	ast.Inspect(node, func(n ast.Node) bool {
		switch nn := n.(type) {
		case *ast.AssignStmt:
			// 处理短声明（如 v := &v2_captcha{}）
			if nn.Tok == token.DEFINE {
				for i, lhs := range nn.Lhs {
					if i >= len(nn.Rhs) {
						break
					}
					varName := lhs.(*ast.Ident).Name
					typeName := ResolveVarType(nn.Rhs[i])
					if typeName != "" {
						results[varName] = typeName
					}
				}
			}
		case *ast.ValueSpec:
			// 处理普通声明（如 var v = &v2_captcha{}）
			for i, name := range nn.Names {
				if i >= len(nn.Values) {
					break
				}
				varName := name.Name
				typeName := ResolveVarType(nn.Values[i])
				if typeName != "" {
					results[varName] = typeName
				}
			}
		}
		return true
	})
	return results
}

// ResolveVarType 解析表达式的实际类型名（如 v2_captcha）
func ResolveVarType(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.UnaryExpr:
		// 处理指针类型（如 &v2_captcha{}）
		if t.Op == token.AND {
			return ResolveVarType(t.X)
		}
	case *ast.CompositeLit:
		// 处理结构体初始化（如 v2_captcha{}）
		if ident, ok := t.Type.(*ast.Ident); ok {
			return ident.Name
		}
	case *ast.CallExpr:
		// 处理构造函数（如 new(v2_captcha)）
		if fun, ok := t.Fun.(*ast.Ident); ok && fun.Name == "new" {
			if arg, ok := t.Args[0].(*ast.Ident); ok {
				return arg.Name
			}
		}
	}
	return ""
}
