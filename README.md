# 如何写go



## 代码组织

Go程序被组织成包。包是同一目录下的源文件的集合，它们被编译在一起。在一个源文件中定义的函数、类型、变量和常量对同一包内的所有其他源文件都是可见的。

一个资源库包含一个或多个模块。一个模块是相关 Go 包的集合，它们被一起发布。一个 Go 代码库通常只包含一个模块，位于代码库的根部。那里有一个名为go.mod的文件声明了模块路径：模块内所有包的导入路径前缀。该模块包含包及其go.mod文件的目录中的软件包，以及该目录的子目录，直到包含另一个go.mod文件的下一个子目录（如果有的话）。

请注意，你不需要在构建之前将你的代码发布到一个远程仓库。一个模块可以在本地定义而不属于一个仓库。然而，组织你的代码是一个好习惯，就像你有一天会发布它一样。

每个模块的路径不仅是其软件包的导入路径前缀，而且还指明了 go 命令应该从哪里下载它。例如，为了下载golang.org/x/tools模块，go命令会查阅`https://golang.org/x/tools`（[这里有更多描述](https://go.dev/cmd/go/#hdr-Remote_import_paths)）所指示的仓库。

导入路径是一个用于导入软件包的字符串。一个包的导入路径是它的模块路径与模块内的子目录相连。例如，模块github.com/google/go-cmp在cmp/目录下包含一个包。该包的导入路径是github.com/google/go-cmp/cmp。标准库中的包没有模块路径前缀。



## 第一个程序

要编译和运行一个简单的程序，首先要选择一个模块路径（我们将使用example/user/hello），并创建一个go.mod文件来声明它。

```shell
$ mkdir hello # 或者从远程仓库拉下来的目录
$ cd hello
$ go mod init example/user/hello
go: creating new go.mod: module example/user/hello
$ cat go.mod
module example/user/hello

go 1.18
$
```

Go源文件中的第一条语句必须是包名。可执行的命令必须始终使用软件包main。

接下来，在该目录下创建一个名为hello.go的文件，包含以下Go代码。

```go
package main

import "fmt"

func main() {
	fmt.Println("Hello, world.")
}
```

现在你可以用go工具构建和安装该程序。

```shell
$ go install example/user/hello
$
```

这个命令建立了hello命令，产生了一个可执行的二进制文件。然后，它将该二进制文件安装为$HOME/go/bin/hello（或者，在Windows下，%USERPROFILE%/go/bin/hello.exe）。

安装目录是由GOPATH和GOBIN环境变量控制的。如果设置了GOBIN，二进制文件将被安装到该目录。如果设置了GOPATH，二进制文件将被安装到GOPATH列表中第一个目录的bin子目录中。否则，二进制文件将被安装到默认的 GOPATH（$HOME/go 或 %USERPROFILE%\go）的 bin 子目录。

你可以使用go env命令，为未来的go命令可移植地设置环境变量的默认值。

```shell
$ go env -w GOBIN=/somewhere/else/bin
$
```

要取消之前由go env -w设置的变量，请使用go env -u。

```shell
$ go env -u GOBIN
$
```

像go install这样的命令在包含当前工作目录的模块的上下文中应用。如果工作目录不在example/user/hello模块内，go install可能会失败。

为了方便起见，go命令接受相对于工作目录的路径，如果没有给出其他路径，则默认为当前工作目录下的软件包。因此，在我们的工作目录下，以下命令都是等价的。

```shell
$ go install example/user/hello
$ go install .
$ go install
```

接下来，让我们运行该程序，以确保它的工作。为了方便起见，我们将把安装目录添加到我们的PATH中，使运行二进制文件变得容易。

```she
# Windows users should consult https://github.com/golang/go/wiki/SettingGOPATH
# for setting %PATH%.
$ export PATH=$PATH:$(dirname $(go list -f '{{.Target}}' .))
$ hello
Hello, world.
$
```

go命令通过请求相应的HTTPS URL并读取嵌入HTML响应中的元数据来定位包含给定模块路径的资源库（见go help importpath）。许多托管服务已经为包含Go代码的仓库提供了元数据，所以让你的模块供他人使用的最简单方法通常是使其模块路径与仓库的URL一致。

### 从模块导入包

让我们编写一个morestrings包，并从hello程序中使用它。首先，为该包创建一个名为$HOME/hello/morestrings的目录，然后在该目录下创建一个名为reverse.go的文件，内容如下。

```go
// Package morestrings 实现了额外的函数来操作UTF-8
// 编码的字符串，超出了标准 "strings "包所提供的内容。
package morestrings

// ReverseRunes 返回其参数字符串从左到右颠倒的符文。
func ReverseRunes(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

```

因为我们的ReverseRunes函数以大写字母开头，所以它被导出，可以在其他导入我们的morestrings包的包中使用。

让我们用go build来测试一下该软件包的编译情况。

```shell
$ cd $HOME/hello/morestrings
$ go build
$
```

这不会产生一个输出文件。相反，它会把编译好的软件包保存在本地的构建缓存中。

在确认了morestrings包的构建之后，让我们在hello程序中使用它。要做到这一点，请修改你原来的$HOME/hello/hello.go以使用morestrings包。

```go
package main

import (
	"fmt"

	"example/user/hello/morestrings"
)

func main() {
	fmt.Println(morestrings.ReverseRunes("!oG ,olleH"))
}
```

安装Hello程序。

```shell
$ go install example/user/hello
```

运行新版本的程序，你应该看到一个新的、反转的信息。

```shell
$ hello
Hello, Go!
```



### 从远程模块导入包

一个导入路径可以描述如何使用Git或Mercurial这样的修订控制系统来获取软件包的源代码。go工具使用这个属性来自动从远程仓库获取软件包。例如，要在你的程序中使用github.com/google/go-cmp/cmp。

```go
package main

import (
	"fmt"

	"example/user/hello/morestrings"
	"github.com/google/go-cmp/cmp"
)

func main() {
	fmt.Println(morestrings.ReverseRunes("!oG ,olleH"))
	fmt.Println(cmp.Diff("Hello World", "Hello Go"))
}
```

现在你有了对外部模块的依赖，你需要下载该模块并在你的go.mod文件中记录其版本。`go mod tidy`命令为导入的软件包添加缺少的模块需求，并删除不再使用的模块的需求。

```shell
$ go mod tidy
go: finding module for package github.com/google/go-cmp/cmp
go: found github.com/google/go-cmp/cmp in github.com/google/go-cmp v0.5.4
$ go install example/user/hello
$ hello
Hello, Go!
  string(
- 	"Hello World",
+ 	"Hello Go",
  )
$ cat go.mod
module example/user/hello

go 1.18

require github.com/google/go-cmp v0.5.7
$
```

模块的依赖性被自动下载到GOPATH环境变量所指示的目录下的pkg/mod子目录中。一个给定版本的模块的下载内容在所有需要该版本的模块之间共享，所以go命令将这些文件和目录标记为只读。要删除所有下载的模块，你可以在go clean中传递-modcache标志。

```shell
$ go clean -modcache
$
```



## 测试

Go有一个由go test命令和测试包组成的轻量级测试框架。

你通过创建一个名称以_test.go结尾的文件来编写测试，该文件包含名为TestXXX的函数，其签名为func（t *testing.T）。测试框架运行每个这样的函数；如果该函数调用一个失败函数，如t.Error或t.Fail，则认为测试失败。

通过创建包含以下Go代码的$HOME/hello/morestrings/reverse_test.go文件，向morestrings包添加一个测试。

```go
package morestrings

import "testing"

func TestReverseRunes(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"Hello, world", "dlrow ,olleH"},
		{"Hello, 世界", "界世 ,olleH"},
		{"", ""},
	}
	for _, c := range cases {
		got := ReverseRunes(c.in)
		if got != c.want {
			t.Errorf("ReverseRunes(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}

```

然后用go test运行该测试。

```shell
$ cd $HOME/hello/morestrings
$ go test
PASS
ok      example/user/hello/morestrings  0.021s
$
```

运行[`go help test`](https://go.dev/cmd/go/#hdr-Test_packages)，更多细节见[测试包文档](https://go.dev/pkg/testing/)。



## 接下来

订阅 [golang-announce](https://groups.google.com/group/golang-announce) 邮件列表，以便在 Go 的新稳定版本发布时获得通知。

请参阅[Effective Go](https://go.dev/doc/effective_go.html)，了解编写清晰、简洁的Go代码的技巧。

参加 [Go 之旅](https://go.dev/tour/)以了解该语言本身。

访问[文档页面](https://go.dev/doc/#articles)，了解有关 Go 语言及其库和工具的一系列深度文章。