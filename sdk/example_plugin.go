/*
本文件实现了一个最简单的接口插件

编译： go build -buildmode=plugin example_plugin.go
得到so文件放到plugins文件夹
注意，假设有多个插件重写同一接口，采用第一个读入的插件.

由于在IDE环境编码的原因，包名为sdk，在实际应用中，包名应该为main
 */
package sdk

import (
	"encoding/json"
	"fmt"
)

// 定义结构体
type Behaviors struct {
	OnEnabled      string
	OnDisabled     string
	RequestHandler []RequestHandler
}
type RequestHandler struct {
	RequestName  string
	FunctionName string
}
type Request struct {
	Method    string
	OperateID int
	Message   string
}

type Response struct {
	Status  int
	Message string
}
// Behavior函数，用于告诉FrozenGo-Daemon 该插件的行为
func Behavior() []byte{
	b,_  := json.Marshal( Behaviors{
		OnEnabled:"ImEnabled", // 在开启时会被自动运行的函数名称 函数只能是func()
		OnDisabled:"ImDisabled", // 在被禁用时运行的函数名称，类型都只能是func()
		RequestHandler:[]RequestHandler{ // 数组，对于Request的解析
			{
				"SayHello", // 接口名
				"SayHello",
				// 要实现该接口的函数名 当请求Method为RequestName时，
				// 无论FGO有没实现该接口，都会执行此接口并直接返回。
				// 函数类型为func([]byte)[]byte
				// 参数为json格式的Request 返回json格式的Response
				// Request会由FGO daemon 为您传入，json格式理论上100%正确
				// 不正确请找Google公司
			},
		},
	})
	return b
}

// 实现Behavior里声明的方法
func ImEnabled() {
	fmt.Println("Kawaii!")
	fmt.Println("Xueluo is so lovely!")
}

func ImDisabled(){
	fmt.Println("Emmm...")
	fmt.Println("Xueluo luo ~~~~~Good bye.")
}

func SayHello(requestJson []byte) []byte{
	req := Request{}
	err := json.Unmarshal(requestJson,&req)
	if err != nil {
		fmt.Println("Hi ,Google~")
		// Who cares?
		// ...
		// ...
		// do something
	}
	// ...
	// ... Handle request and make a Response
	// ...
	b,_ := json.Marshal(Response{0,"Wow I get a fresh XueLuo!"})
	return b
}