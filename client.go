package main

import (
	"crypto/md5"
	"fmt"
	"net/rpc"
	"sync"
	"time"
)

const Version string = "1.12.2"

type DataPack struct {
	Id    int   // 操作序号
	Opt   uint8 //
	X, Y  int
	Mymap [8][8]int
}

// 加上锁的 地图数据
var (
	mu   sync.Mutex
	data DataPack
)

func KeepInTouch(client *rpc.Client) {
	for {
		mu.Lock()
		dataclone := data
		mu.Unlock()

		TouchServer(client, dataclone, 'T')              // Touch
		time.Sleep(time.Duration(50) * time.Millisecond) // 每隔 0.05s 通讯一次
	}
}

func TouchServer(client *rpc.Client, cur DataPack, option uint8) { // 与服务器同步数据
	cur.Opt = option
	var reply DataPack

	err := client.Call("McServer.FetchClient", cur, &reply)
	if err != nil {
		fmt.Println(nil, err)
	}

	if option == 'Q' { // 结束游戏不用检查数据包，以 client 端的为准
		return
	} else if reply.Opt == 'C' { // 符合
		return
	} else { // 不符合就强制拉回
		mu.Lock()
		data = reply
		mu.Unlock()
	}
}

func main() {
	fmt.Println("Connecting to server...")

	client, err := rpc.Dial("tcp", ":25565")
	if err != nil {
		fmt.Println("Can't connect server: ", err)
		time.Sleep(time.Duration(5) * time.Second) //给人看的时间
		return
	}
	Login(client)

	go KeepInTouch(client) // 每 0.05秒核对一次
	for {
		cls() // 清屏
		fmt.Println("Yon can press WASD to move, or Q to Quit. Please press enter to confirm your choice.")
		fmt.Println("0 for air, 1 for obstacles and 2 for your character.\n")
		PrintMap()

		var ch uint8 = Getchar()
		// fmt.Println("input char:", ch)
		if ch == 'q' || ch == 'Q' { // 结束游戏
			mu.Lock()                      // Lock 掉, 这样其他进程就阻塞了
			TouchServer(client, data, 'Q') // 退出前最后一次请求数据
			client.Close()                 //客户端关闭连接
			fmt.Println("Quit!")
			break
		} else {
			mu.Lock()
			Move(ch) // 正常移动
			dataclone := data
			mu.Unlock()

			go TouchServer(client, dataclone, data.Opt) // 同步数据
		}
	}
	client.Close()
}

func Getchar() uint8 { // 获取输入字符，自动过滤除 WASDQ 以外的字符
	var ch uint8
	for {
		fmt.Scanf("%c", &ch)
		if ch != 'W' && ch != 'w' && ch != 'a' && ch != 'A' && ch != 's' && ch != 'S' && ch != 'D' && ch != 'd' && ch != 'q' && ch != 'Q' {
			continue
		} else {
			break
		}
		// fmt.Println("input char:", ch)
	}
	return ch
}

func Move(ch uint8) { //根据 ch 移动 data
	data.Id++ // 操作序号
	var dx, dy int

	if ch == 'W' || ch == 'w' {
		data.Opt = 'W'
		dx, dy = -1, 0
	} else if ch == 'A' || ch == 'a' {
		data.Opt = 'A'
		dx, dy = 0, -1
	} else if ch == 'S' || ch == 's' {
		data.Opt = 'S'
		dx, dy = 1, 0
	} else if ch == 'D' || ch == 'd' {
		data.Opt = 'D'
		dx, dy = 0, 1
	} else {
		// 前面的已经屏蔽了
		fmt.Println("Are you kidding me?")
	}

	var nx int = data.X + dx
	var ny int = data.Y + dy
	// fmt.Println("(nx, ny): ", nx, ny)

	if nx < 0 || nx >= 8 || ny < 0 || ny >= 8 {
		return // 出界
	} else if data.Mymap[nx][ny] == 1 {
		return // 有障碍物
	} else {
		// move
		data.X, data.Y = nx, ny
	}
}

func Login(client *rpc.Client) { // 登录模块
	var username, passwd string
	for {
		fmt.Print("username: ")
		fmt.Scan(&username)
		// fmt.Println(username)
		fmt.Print("password: ")
		fmt.Scan(&passwd)
		// fmt.Println(passwd)
		usermd5 := fmt.Sprintf("%x", md5.Sum([]byte(passwd+"|"+username)))
		usermd5 = usermd5 + "|" + username + "|" + Version
		fmt.Println(usermd5)

		var pkg DataPack
		err := client.Call("McServer.Login", usermd5, &pkg)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		fmt.Println(pkg)

		if pkg.Opt == 'A' {
			data = pkg
			break
		} else {
			fmt.Println("The user does not exist or your password is wrong!")
		}
	}
	fmt.Println("Login successfully!")
}

// 还没找到很好的清屏函数
func cls() { // 清屏
	// c := exec.Command("clear")
	// c.Stdout = os.Stdout
	// c.Run()
	// c = exec.Command("pause")
	// c.Run()
	// tm.clear()
	fmt.Print("\033[H\033[2J")
}

func PrintMap() { // 打印地图
	mu.Lock()
	defer mu.Unlock()

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if i == data.X && j == data.Y {
				fmt.Print("2 ")
			} else {
				fmt.Print(data.Mymap[i][j], " ")
			}
		}
		fmt.Print("\n")
	}
}
