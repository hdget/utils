package ast

import "go/ast"

type CallSignature struct {
	FunctionChain      string         // 必传
	Package            string         // 可选
	ArgCount           int            // 可选, -1不去检查
	ArgIndex2Signature map[int]string // 可选, nil不去检查
}

// GetCaller 获取调用者
func GetCaller(n *ast.CallExpr) (string, bool) {
	// 递归查找调用者
	for {
		// 检查 Fun 是否是 SelectorExpr
		selectorExpr, ok := n.Fun.(*ast.SelectorExpr)
		if !ok {
			break
		}

		// 检查 X 是否是另一个 CallExpr
		nextCallExpr, ok := selectorExpr.X.(*ast.CallExpr)
		if !ok {
			// 如果不是 CallExpr，则可能是调用者（如 sdk）
			if ident, ok := selectorExpr.X.(*ast.Ident); ok {
				return ident.Name, true
			}
			break
		}

		// 继续递归
		n = nextCallExpr
	}

	return "", false
}

func MatchCall(n *ast.CallExpr, signature *CallSignature, imports map[string]string) bool {
	if GetFunctionChain(n) != signature.FunctionChain {
		return false
	}

	if signature.ArgCount > 0 {
		if len(n.Args) != signature.ArgCount {
			return false
		}
	}

	// 如果传入了pkg，则调用类似: dapr.New
	if signature.Package != "" {
		caller, found := GetCaller(n)
		if !found {
			return false
		}
		if imports[caller] != signature.Package {
			return false
		}
	}

	for argIndex, s := range signature.ArgIndex2Signature {
		if argIndex >= len(n.Args) {
			return false
		}

		if GetExprTypeName(n.Args[argIndex]) != s {
			return false
		}
	}

	return true
}
