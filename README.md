# hdutils

Go 通用工具集。

> **本仓库是多模块(multi-module)仓库。** 每个子目录都是一个独立的 Go 模块,拥有独立版本号,
> 需要**单独 require**。这是刻意设计:让下游只为真正用到的包付出依赖代价
> (例如用 `panic` 的人不必继承 `sql`/`neo4j` 的驱动依赖)。

## 关键:子包不会随父模块一起发布

Go 打包父模块时,会**排除任何自带 `go.mod` 的子目录**。因此 `github.com/hdget/utils` 的版本里
**不包含** `panic`、`sql` 等子包。只 require 父模块却去 import 子包,会得到:

```
module github.com/hdget/utils@latest found (v0.2.4), but does not contain package github.com/hdget/utils/panic
```

正确做法是单独引入子模块:

```bash
go get github.com/hdget/utils/panic@panic/v0.0.1
```

同理,针对父模块的 `replace` **不会**传递到子模块(非通配 replace 只匹配精确模块路径),
子模块需要各自的 `replace` 条目。

## 模块与版本

| 模块 | 引入方式 | 最新版本 |
| --- | --- | --- |
| 根模块 `utils` | `go get github.com/hdget/utils` | `v0.2.4` |
| `ast` | `go get github.com/hdget/utils/ast` | `ast/v0.0.1` |
| `cmp` | `go get github.com/hdget/utils/cmp` | `cmp/v0.0.2` |
| `currency` | `go get github.com/hdget/utils/currency` | `currency/v0.0.1` |
| `encoding` | `go get github.com/hdget/utils/encoding` | `encoding/v0.0.1` |
| `hash` | `go get github.com/hdget/utils/hash` | `hash/v0.0.2` |
| `json` | `go get github.com/hdget/utils/json` | `json/v0.0.2` |
| `neo4j` | `go get github.com/hdget/utils/neo4j` | `neo4j/v0.0.1` |
| `paginator` | `go get github.com/hdget/utils/paginator` | `paginator/v0.0.2` |
| `panic` | `go get github.com/hdget/utils/panic` | `panic/v0.0.1` |
| `parallel` | `go get github.com/hdget/utils/parallel` | `parallel/v0.0.1` |
| `reflect` | `go get github.com/hdget/utils/reflect` | `reflect/v0.0.1` |
| `sql` | `go get github.com/hdget/utils/sql` | `sql/v0.0.7` |
| `text` | `go get github.com/hdget/utils/text` | `text/v0.0.4` |
| `time` | `go get github.com/hdget/utils/time` | `time/v0.0.5` |

`logger` 不是独立模块(无 `go.mod`),它随**根模块**发布,直接 `import "github.com/hdget/utils/logger"` 即可。

> ⚠️ `cmp`、`hash`、`json` 的 **`v0.0.1` 永久不可用**:它们的 `go.mod` 漏了对根模块的 `require`,
> 单独构建会报 `cannot find module providing package github.com/hdget/utils`。Go 模块版本一经发布
> 不可修改或撤回内容,所以这三个 `v0.0.1` 标签会一直留在代理缓存里 —— 使用或升级时请固定到 `v0.0.2` 及以上。

## 发布新版本

```powershell
./git_release.ps1 -Module panic              # 发布子模块,patch 位 +1
./git_release.ps1 -Module . -Bump minor      # 发布根模块,minor 位 +1
./git_release.ps1 -Module json -DryRun       # 只演练:检查版本号、脏工作区、能否构建
```

标签命名规则:子模块为 `<目录>/vX.Y.Z`,根模块为裸 `vX.Y.Z`。

脚本会自动取该模块已有最高版本并递增,并在打标签前完成三项检查:
版本号未被占用、该模块工作区已提交、该模块能独立构建。

**版本号只能递增,绝不能复用。** Go 模块代理(`proxy.golang.org` 及一切镜像)对同一版本号只抓取
一次并永久缓存,把已发布的标签指向新提交不会让下游拿到新代码,反而可能与 `go.sum` 冲突报
`checksum mismatch`。这也是旧版 `git_tag.ps1`(删标签重建同名标签)被移除的原因。
