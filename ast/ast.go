package ast

type ASTUtils interface {
	InspectFunction(srcPath string, fnParams, fnResults []string, annotationPrefix string) ([]*Function, error) // 从源代码目录中获取fnParams和fnResults匹配的函数的信息,并解析函数对应的注解
	InspectPackage(rootDir string, excludeDirs []string) (map[string]*Package, error)                           // 从源代码中获取包的相对路径到包信息的map
}
