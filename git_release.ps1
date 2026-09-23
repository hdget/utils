# 分模块发布脚本:每个模块目录独立递增版本号并打 tag。
#
# 为什么必须「只递增、绝不复用」:Go 模块代理(proxy.golang.org 及一切镜像)
# 对同一个版本号只抓取一次并永久缓存。把已发布的 tag 指向另一个提交,
# 下游拿到的仍是旧内容,且可能与 go.sum 冲突报 checksum mismatch。
#
# 用法:
#   ./git_release.ps1 -Module panic              # 发布子模块 panic,patch +1
#   ./git_release.ps1 -Module . -Bump minor      # 发布根模块,minor +1
#   ./git_release.ps1 -Module json -DryRun       # 只演练,不创建也不推送

param(
    # '.' 代表根模块 github.com/hdget/utils,其余为子模块目录名
    [Parameter(Mandatory = $true)][string]$Module,
    [ValidateSet('patch', 'minor', 'major')][string]$Bump = 'patch',
    [switch]$DryRun
)
$ErrorActionPreference = 'Stop'

# 标签命名规则:子模块带目录前缀(panic/v0.0.2),根模块为裸版本(v0.2.5)。
# 前缀同时保证根模块的版本查询不会误命中子模块标签。
$Prefix = if ($Module -eq '.') { 'v' } else { "$Module/v" }
$Pattern = "^$([regex]::Escape($Prefix))(\d+)\.(\d+)\.(\d+)$"

# 嵌套模块目录清单:这些目录各自有 go.mod,因而被父模块的 zip 排除。
$Nested = @(git ls-files '*go.mod' | ForEach-Object { Split-Path $_ -Parent } | Where-Object { $_ })

# 取出该模块已有的最高版本号;从未发布过则从 0.0.1 起。
$Current = git tag -l "$Prefix*" | ForEach-Object {
    if ($_ -match $Pattern) { [pscustomobject]@{ Major = [int]$Matches[1]; Minor = [int]$Matches[2]; Build = [int]$Matches[3] } }
} | Sort-Object Major, Minor, Build -Descending | Select-Object -First 1

if (-not $Current) {
    $Next = [pscustomobject]@{ Major = 0; Minor = 0; Build = 1 }
} else {
    $Next = switch ($Bump) {
        'major' { @{ Major = $Current.Major + 1; Minor = 0; Build = 0 } }
        'minor' { @{ Major = $Current.Major; Minor = $Current.Minor + 1; Build = 0 } }
        'patch' { @{ Major = $Current.Major; Minor = $Current.Minor; Build = $Current.Build + 1 } }
    }
}
$Tag = '{0}{1}.{2}.{3}' -f $Prefix, $Next.Major, $Next.Minor, $Next.Build
$ModulePath = if ($Module -eq '.') { 'github.com/hdget/utils' } else { "github.com/hdget/utils/$Module" }

Write-Host "模块  : $ModulePath" -ForegroundColor Cyan
Write-Host "当前  : $(if ($Current) { '{0}{1}.{2}.{3}' -f $Prefix, $Current.Major, $Current.Minor, $Current.Build } else { '(该模块尚未发布过)' })" -ForegroundColor Cyan
Write-Host "发布  : $Tag" -ForegroundColor Green

# 守卫 1:拒绝复用已存在的版本号。
if (git rev-parse -q --verify "refs/tags/$Tag") {
    throw "标签 $Tag 已存在。Go 模块版本号不可复用,请改用 -Bump minor 或 -Bump major。"
}

# 守卫 2:该模块涉及的工作区必须干净 —— tag 指向的是已提交内容,
# 脏工作区会发布出与你机器上不一致的包。根模块需排除各子模块目录。
if ($Module -eq '.') {
    $Dirty = git status --porcelain | ForEach-Object {
        $p = $_.Substring(3)
        if ($p -match '^"(.*)"$') { $p = $Matches[1] }
        if ($p -match ' -> ') { $p = ($p -split ' -> ')[-1] }
        if (@($Nested | Where-Object { $p -like "$_/*" }).Count -eq 0) { $p }
    }
} else {
    $Dirty = git status --porcelain -- $Module | ForEach-Object { $_.Substring(3) }
}
if (@($Dirty).Count -gt 0) {
    throw "模块 $Module 有未提交改动,请先提交再发布:`n  $(@($Dirty) -join "`n  ")"
}

# 守卫 3:发布前必须能独立构建,避免把编不过的版本推给下游。
Push-Location $Module
try {
    go build ./... 2>&1 | ForEach-Object { Write-Host "  $_" }
    if ($LASTEXITCODE -ne 0) { throw "$Module 构建失败,已中止发布(不要发布编不过的版本)" }
} finally {
    Pop-Location
}

if ($DryRun) {
    Write-Host "[DryRun] 检查通过,未创建标签 $Tag" -ForegroundColor Yellow
    return
}

git tag -a $Tag -m "release $Tag"
git push origin $Tag
Write-Host "已发布 $Tag" -ForegroundColor Green
Write-Host "下游接入: go get $ModulePath@$Tag" -ForegroundColor Cyan
