package ast

import (
	"fmt"
	"go/ast"
)

// GetExprTypeName 返回带有指针指示的类型名称字符串
func GetExprTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.BasicLit:
		return t.Kind.String() // "INT", "STRING"等
	case *ast.CompositeLit:
		return GetExprTypeName(t.Type)
	case *ast.CallExpr:
		return GetExprTypeName(t.Fun) + "()"
	case *ast.UnaryExpr:
		return fmt.Sprintf("%s%s", t.Op, GetExprTypeName(t.X))
	case *ast.Ident:
		return t.Name // 基础类型
	case *ast.StarExpr:
		return "*" + GetExprTypeName(t.X) // 指针类型加*前缀
	case *ast.SelectorExpr:
		return GetExprTypeName(t.X) + "." + t.Sel.Name // 包.类型
	case *ast.ArrayType:
		return "[]" + GetExprTypeName(t.Elt) // 切片类型
	case *ast.MapType:
		// 特别处理map的value部分，区分指针
		keyType := GetExprTypeName(t.Key)
		valueType := GetExprTypeName(t.Value)
		return fmt.Sprintf("map[%s]%s", keyType, valueType)
	case *ast.InterfaceType:
		return "interface{}" // 接口类型
	default:
		return fmt.Sprintf("%T", expr) // 其他未知类型
	}
}
