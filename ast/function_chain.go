package ast

import (
	"go/ast"
	"strings"
)

// GetFunctionChain 获取完整的函数调用链
func GetFunctionChain(n *ast.CallExpr) string {
	functions := ParseFunctionCallChain(n)
	return strings.Join(reverseSlice(functions), ".")
}

// ParseFunctionCallChain 递归解析链式函数调用，最近的Ident.Name作为包名，最先调用的函数在slice的最前面
func ParseFunctionCallChain(n *ast.CallExpr) []string {
	var methods []string

	// 递归提取方法名
	for {
		// 检查 Fun 是否是 SelectorExpr
		selectorExpr, ok := n.Fun.(*ast.SelectorExpr)
		if !ok {
			break
		}

		// 添加方法名
		methods = append(methods, selectorExpr.Sel.Name)

		// 检查 X 是否是另一个 CallExpr
		nextCallExpr, ok := selectorExpr.X.(*ast.CallExpr)
		if !ok {
			break
		}

		// 继续递归
		n = nextCallExpr
	}
	return methods
}

func reverseSlice[T any](ss []T) []T {
	// Avoid the allocation. If there is one element or less it is already
	// reversed.
	if len(ss) < 2 {
		return ss
	}

	sorted := make([]T, len(ss))
	for i := 0; i < len(ss); i++ {
		sorted[i] = ss[len(ss)-i-1]
	}

	return sorted
}
