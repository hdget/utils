package ast

//
//// GetFunctionChain 获取完整的函数调用链
//func GetFunctionChain(n *ast.CallExpr) string {
//	functions := make([]string, 0)
//	astRecursiveParseFunction(n, functions)
//	return strings.Join(pie.Reverse(functions), ".")
//}
//
//// GetEmbedVarAndRelPath 获取嵌入资源的信息，返回变量名，embed路径
//func GetEmbedVarAndRelPath(n *ast.GenDecl) (string, string, bool) {
//	// 如果是 GenDecl 类型，则可能是 import 或者变量声明等
//	if n.Tok == token.VAR {
//		for _, spec := range n.Specs {
//			if valueSpec, ok := spec.(*ast.ValueSpec); ok {
//				if astIsEmbedFSType(valueSpec.Type) {
//					return valueSpec.Names[0].name, astGetEmbedRelPath(n), true
//				}
//			}
//		}
//	}
//	return "", "", false
//}
//
//// 检查类型是否为 embed.FS
//func IsEmbedFSType(expr ast.Expr) bool {
//	if selectorExpr, ok := expr.(*ast.SelectorExpr); ok {
//		if ident, ok := selectorExpr.X.(*ast.Ident); ok && ident.name == "embed" {
//			if selectorExpr.Sel.name == "FS" {
//				return true
//			}
//		}
//	}
//	return false
//}
//
//// 获取 embed 路径
//func GetEmbedRelPath(n *ast.GenDecl) string {
//	// 如果直接定义变量
//	// //go:embed assets/*
//	// var assets embed.FS
//	if n.Doc != nil {
//		for _, comment := range n.Doc.List {
//			if strings.HasPrefix(comment.Text, "//go:embed") {
//				// 提取路径部分
//				return filepath.Dir(strings.TrimSpace(strings.TrimPrefix(comment.Text, "//go:embed")))
//			}
//		}
//	}
//	// 如果定义在var block中
//	// var (
//	//   //go:embed assets/*
//	//   assets embed.FS
//	// )
//	for _, spec := range n.Specs {
//		if valueSpec, ok := spec.(*ast.ValueSpec); ok {
//			if valueSpec.Doc != nil {
//				for _, comment := range valueSpec.Doc.List {
//					if strings.HasPrefix(comment.Text, "//go:embed") {
//						// 提取路径部分
//						return filepath.Dir(strings.TrimSpace(strings.TrimPrefix(comment.Text, "//go:embed")))
//					}
//				}
//			}
//		}
//	}
//	return ""
//}
//
//// Parse 尝试从源代码中查找嵌入路径, 返回嵌入资源的绝对路径和相对路径
//func ParseEmbed(callerFilePath string) (string, string, error) {
//	// 创建一个新的文件集
//	fset := token.NewFileSet()
//
//	// 解析源文件，同时保留注释
//	f, err := parser.ParseFile(fset, callerFilePath, nil, parser.ParseComments)
//	if err != nil {
//		return "", "", err
//	}
//
//	// 遍历AST节点
//	count := 0
//	var foundVar, foundRelPath, embedAbsPath string
//	ast.Inspect(f, func(node ast.Node) bool {
//		switch n := node.(type) {
//		case *ast.GenDecl:
//			if varName, relPath, ok := astGetEmbedVarAndRelPath(n); ok {
//				foundVar = varName
//				foundRelPath = relPath
//				return false
//			}
//		}
//		count += 1
//		return foundVar == ""
//	})
//
//	fmt.Println("xxxxxxxxxxxxxx:", count)
//
//	if foundVar == "" {
//		return "", "", fmt.Errorf("embed.FS variable declare not found, var: %s", foundVar)
//	}
//
//	// 有可能定义了embed.FS,但是没有指定编译指令//go:embed
//	if foundRelPath == "" {
//		return "", "", fmt.Errorf("//go:embed compiler directive not found, var: %s", foundVar)
//	}
//
//	if foundRelPath == "." {
//		return "", "", fmt.Errorf("//go:embed must specify a directory, var: %s", foundVar)
//	}
//
//	embedAbsPath = filepath.Join(filepath.Dir(callerFilePath), foundRelPath)
//	return embedAbsPath, foundRelPath, nil
//}
