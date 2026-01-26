package ast

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

type Package struct {
	Name    string
	Files   map[string]*ast.File
	FileSet *token.FileSet
}

// GetPackageImportPaths 获取包导入路径
func GetPackageImportPaths(f *ast.File) map[string]string {
	importMap := make(map[string]string)
	for _, imp := range f.Imports {
		pkgNames := make([]string, 0)
		if imp.Name != nil {
			pkgNames = []string{imp.Name.Name} // 处理别名导入，如 `import alias "math/rand"`
		} else {
			// 提取完整路径（去掉引号）
			pkgPath := strings.Trim(imp.Path.Value, `"`)

			// 获取包名（路径的最后一部分）
			lastPart := pkgPath[strings.LastIndex(pkgPath, "/")+1:]
			// HOTFIX: 有时候定义包名会只使用横杠后的部分，例如: lib-dapr只会用dapr
			pkgNames = append(pkgNames, lastPart)

			possiblePkgName := lastPart[strings.LastIndex(lastPart, "-")+1:]
			if possiblePkgName != lastPart {
				pkgNames = append(pkgNames, possiblePkgName)
			}
		}

		for _, pkgName := range pkgNames {
			importMap[pkgName] = strings.Trim(imp.Path.Value, `"`)
		}
	}
	return importMap
}

// InspectPackage 从源代码目录中获取包的信息, 包相对路径=>包信息
func InspectPackage(srcDir string, excludeDirs []string) (map[string]*Package, error) {
	// 获取根目录绝对路径
	absRoot, err := filepath.Abs(srcDir)
	if err != nil {
		return nil, err
	}

	// 转换排除目录为绝对路径
	var absExcludeDirs []string
	for _, dir := range excludeDirs {
		absDir := filepath.Join(absRoot, dir)
		absExcludeDirs = append(absExcludeDirs, absDir)
	}

	// 启动工作 goroutine
	var mu sync.Mutex
	var wg sync.WaitGroup

	// 使用 WaitGroup 等待所有解析完成
	fileChan := make(chan string, 100)
	fset := token.NewFileSet()
	pkgRelPath2astPkg := make(map[string]*Package)
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for filePath := range fileChan {
				pkgRelPath, _ := filepath.Rel(absRoot, filePath)
				pkgRelPath = filepath.ToSlash(pkgRelPath)

				// 解析单个文件
				pkgName2astPkg, err := parser.ParseDir(fset, filePath, nil, parser.ParseComments)
				if err != nil {
					continue
				}

				mu.Lock()
				for pkgName, astPkg := range pkgName2astPkg {
					pkgRelPath2astPkg[pkgRelPath] = &Package{
						Name:    pkgName,
						Files:   astPkg.Files,
						FileSet: fset,
					}
					break
				}
				mu.Unlock()
			}
		}()
	}

	// 遍历所有子目录
	err = filepath.Walk(absRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 检查是否在排除目录中
		for _, excludeDir := range absExcludeDirs {
			if strings.HasPrefix(path, excludeDir) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		// 只处理 .go 文件
		if info.IsDir() {
			fileChan <- path
		}

		return nil
	})

	close(fileChan)
	wg.Wait()

	if err != nil {
		return nil, errors.Wrapf(err, "traverse dir, srcDir: %s", srcDir)
	}

	return pkgRelPath2astPkg, nil
}
