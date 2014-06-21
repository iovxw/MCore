/*
__/\\\\\\\\\\\\\____/\\\\\\__________________________________________________________/\\\________/\\\\\\\_______________/\\\____
 _\/\\\/////////\\\_\////\\\__________________________________/\\\__________________/\\\\\______/\\\/////\\\___________/\\\\\____
  _\/\\\_______\/\\\____\/\\\_________________________________\/\\\________________/\\\/\\\_____/\\\____\//\\\________/\\\/\\\____
   _\/\\\\\\\\\\\\\\_____\/\\\_____/\\\____/\\\_____/\\\\\\\\__\/\\\\\\\\_________/\\\/\/\\\____\/\\\_____\/\\\______/\\\/\/\\\____
    _\/\\\/////////\\\____\/\\\____\/\\\___\/\\\___/\\\/////\\\_\/\\\////\\\_____/\\\/__\/\\\____\/\\\_____\/\\\____/\\\/__\/\\\____
     _\/\\\_______\/\\\____\/\\\____\/\\\___\/\\\__/\\\\\\\\\\\__\/\\\\\\\\/____/\\\\\\\\\\\\\\\\_\/\\\_____\/\\\__/\\\\\\\\\\\\\\\\_
      _\/\\\_______\/\\\____\/\\\____\/\\\___\/\\\_\//\\///////___\/\\\///\\\___\///////////\\\//__\//\\\____/\\\__\///////////\\\//__
       _\/\\\\\\\\\\\\\/___/\\\\\\\\\_\//\\\\\\\\\___\//\\\\\\\\\\_\/\\\_\///\\\___________\/\\\_____\///\\\\\\\/_____________\/\\\____
        _\/////////////____\/////////___\/////////_____\//////////__\///____\///____________\///________\///////_______________\///_____
*/
package main

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	var (
		user        string
		name        string
		password    string
		clientToken string
		accessToken string
		id          string
		location    string
		mode        string
		needHelp    bool
	)
	read:
	for i, value := range os.Args {
		//fmt.Println(i,value)
		switch value {
		case "-h":
			needHelp = true
			break read
		case "-u":
			user = os.Args[i+1]
		case "-n":
			name = os.Args[i+1]
		case "-p":
			password = os.Args[i+1]
		case "-c":
			clientToken = os.Args[i+1]
		case "-a":
			accessToken = os.Args[i+1]
		case "-i":
			id = os.Args[i+1]
		case "-l":
			location = os.Args[i+1]
		case "-m":
			mode = os.Args[i+1]
		}
	}
	if needHelp {
		fmt.Println(`  __  __  _____
 |  \/  |/ ____|
 | \  / | |     ___  _ __ ___
 | |\/| | |    / _ \| '__/ _ \
 | |  | | |___| (_) | | |  __/
 |_|  |_|\_____\___/|_|  \___|
====== MCore启动器核心帮助 ======
呃，这里是帮助= =
遇到的问题基本都会在这里列出来

基本参数:
(注意，输入参数时不需要加上<>)

-u <正版登陆名>     当涉及到需要登陆的模式时需要加上
-p <正版密码>       输入登陆名的时候肯定要加上密码
-n <玩家游戏名>     刷新密钥的时候需要
-c <Client Token>   客户端密钥，特定模式需要
-a <Access Toekn>   访问密钥，特定模式需要
-i <ID>             这个就不需要解释了吧
-l <绝对路径>       客户端的绝对路径，包含.minecraft
-m <模式名称>       必须要输入，启用的模式
-h                  查看帮助，就是现在显示的东西

使用例子:

1.获取正版登陆信息:
-mode login -u example@123.com -p password
在程序后面添加以上参数就会返回登陆密钥等信息，当然用户名和密码得替换成玩家的

2.刷新密钥
-mode refresh -n Bluek404 -i 45Jd55w7dw3Gdwd -c wa6dDwdf556df -a Jdw534dwDHHdw2
就会返回刷新的密钥等东西，不过用的人比较少所以这个功能未完善～
就是说暂时不能用

3.获取游戏列表信息
-mode luanch -l /home/bluek404/.minecraft
就是游戏的目录，记住一定要是绝对路径
不懂啥叫绝对路径的请百度
然后记得带上.minecraft文件夹
会返回所有存在的游戏版本&启动所需的Lib文件路径
Lib文件路径什么的直接添加到启动命令就行鸟～
以后还会可以自定义返回信息的格式&顺序的`)
		return
	}
	switch mode {
	case "login":
		if user != "" && password != "" {
			fmt.Println(authenticate(user, password, clientToken))
		} else {
			fmt.Println("缺少参数，无法运行。请使用参数 -h 来查看帮助")
		}
	case "refresh":
		if name != "" && id != "" && clientToken != "" && accessToken != "" {
			fmt.Println(refresh(accessToken, clientToken, id, name))
		} else {
			fmt.Println("缺少参数，无法运行。请使用参数 -h 来查看帮助")
		}
	case "luanch":
		if location != "" {
			for _, value := range luanch(location) {
				fmt.Println(value["name"])
				fmt.Println(value["version"])
				fmt.Println(value["lib"])
			}
		} else {
			fmt.Println("缺少参数，无法运行。请使用参数 -h 来查看帮助")
		}
	default:
		fmt.Println("运行模式错误！请使用参数 -h 来查看帮助")
	}
}

func authenticate(name_l, password, clientToken_l string) (name, id, accessToken, clientToken string, badlogin bool) {
	x := Login{Agen{"Minecraft", 1}, name_l, password, clientToken_l}
	y, _ := json.Marshal(x) //格式化json
	//fmt.Println("POST：", string(y))

	url := "https://authserver.mojang.com/authenticate"
	b := strings.NewReader(string(y))
	resp, err := http.Post(url, "application/json", b)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(resp)

	body, _ := ioutil.ReadAll(resp.Body) //读取返回的json
	//fmt.Println(string(body))

	js, err := simplejson.NewJson(body) //解析json
	if err != nil {
		panic(err.Error())
	}

	js_map, _ := js.Map() //把解析的json提取到map里

	if resp.Status != "200 OK" {
		//fmt.Println("验证失败")
		badlogin = true
		//fmt.Println(js_map) //如果验证出错，那么打印错误信息（这里就懒得给用户解释详细错误了）
	} else {
		//fmt.Println("验证成功")
		accessToken = js_map["accessToken"].(string)
		//fmt.Println(accessToken)
		clientToken = js_map["clientToken"].(string)
		//fmt.Println(clientToken)
		//fmt.Println(js_map["selectedProfile"])
		//fmt.Println(js_map["availableProfiles"])
		name = js_map["selectedProfile"].(map[string]interface{})["name"].(string)
		//fmt.Println(name)
		id = js_map["selectedProfile"].(map[string]interface{})["id"].(string)
		//fmt.Println(id)
		//fmt.Println(js_map["availableProfiles"].([]interface{})[0].(map[string]interface{})["id"]) //这里的数组是多账号切换用的，暂时先不开发
		//fmt.Println(js_map["availableProfiles"].([]interface{})[0].(map[string]interface{})["name"])
	}
	return
}

func refresh(accessToken_r, clientToken_r, id_r, name_r string) (accessToken, clientToken, id, name string, badlogin bool) {
	x := Refresh{accessToken_r, clientToken_r, SelectedProfile_type{id_r, name_r}}
	y, _ := json.Marshal(x) //格式化json
	//fmt.Println("POST：", string(y))

	url := "https://authserver.mojang.com/refresh"
	b := strings.NewReader(string(y))
	resp, err := http.Post(url, "application/json", b)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(resp)

	body, _ := ioutil.ReadAll(resp.Body) //读取返回的json
	//fmt.Println(string(body))

	js, err := simplejson.NewJson(body) //解析json
	if err != nil {
		panic(err.Error())
	}

	js_map, _ := js.Map() //把解析的json提取到map里

	if resp.Status != "200 OK" {
		//fmt.Println("验证失败")
		badlogin = true
		fmt.Println(js_map) //如果验证出错，那么打印错误信息（这里就懒得给用户解释详细错误了）
	} else {
		//fmt.Println("验证成功")
		//TODO: 待填坑
	}
	return
}

func luanch(location string) (command []map[string]string) {
	launcher_profiles, _ := ioutil.ReadFile(location + "/launcher_profiles.json")
	js, err := simplejson.NewJson(launcher_profiles) //解析json
	if err != nil {
		panic(err.Error())
	}

	js_map, _ := js.Map()                                                     //把解析的json提取到map里
	i := 0                                                                    //计数
	for version, value := range js_map["profiles"].(map[string]interface{}) { //遍历map
		name := value.(map[string]interface{})["name"].(string)
		file := location + "/versions/" + name + "/" + name + ".json"
		if _, err := os.Stat(file); err == nil { //判断json文件是否存在
			command = append(command, map[string]string{
				"name":    name,
				"version": version,
			})
			bi, _ := ioutil.ReadFile(file) //读取这个版本的json文件
			var json_v Json_version
			json.Unmarshal(bi, &json_v)
			for _, value := range json_v.Libraries { //遍历数组
				name := value["name"].(string)
				x := strings.Index(name, ":") - 1
				t1 := strings.Split(string(name[x:len(name)]), ":")
				if value["natives"] != nil {
					natives := (value["natives"].(map[string]interface{}))["linux"]
					if natives != nil {
						t2 := location + "/libraries/" + strings.Replace(string(name[0:x]), ".", "/", -1) + "/" + t1[1] + "/" + t1[2] + "/" + t1[1] + "-" + t1[2] + "-" + natives.(string) + ".jar;"
						command[i]["lib"] = command[i]["lib"] + t2
					}
				} else {
					t2 := location + "/libraries/" + strings.Replace(string(name[0:x]), ".", "/", -1) + "/" + t1[1] + "/" + t1[2] + "/" + t1[1] + "-" + t1[2] + ".jar;"
					command[i]["lib"] = command[i]["lib"] + t2
				}
			}
			i++
		}
	}
	return
}

type Login struct {
	Agent       Agen   `json:"agent"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	ClientToken string `json:"clientToekn"`
}
type Agen struct {
	Name    string `json:"name"`
	Version int    `json:"version"`
}
type Refresh struct {
	AccessToken     string               `json:"accessToken"`
	ClientToken     string               `json:"clientToken"`
	SelectedProfile SelectedProfile_type `json:"selectedProfile"`
}
type SelectedProfile_type struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Json_version struct {
	ID                     string                   `json:"id"`
	Time                   string                   `json:"time"`
	ReleaseTime            string                   `json:"releaseTime"`
	Type                   string                   `json:"type"`
	MinecraftArguments     string                   `json:"minecraftArguments"`
	Libraries              []map[string]interface{} `json:"libraries"`
	MainClass              string                   `json:"mainClass"`
	MinimumLauncherVersion int                      `json:"minimumLauncherVersion"`
	Assets                 string                   `json:"assets"`
}
