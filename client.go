package main

import (
	"crypto/md5"
	"fmt"
	"net/rpc"
	"os"
	"os/exec"
	"sync"
)

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

func cls() { // 清屏
	command := exec.Command("clear")
	command.Stdout = os.Stdout
	command.Run()
}

func PrintMap() { // 打印地图
	mu.Lock()
	defer mu.Unlock()

	cls()

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

func main() {
	fmt.Println("Connecting to server...")

	client, err := rpc.Dial("tcp", ":25565")
	if err != nil {
		fmt.Println("Error: ", err)
	}
	Login(client)

	for {
		PrintMap()
		var ch uint8
		fmt.Scan(&ch)
		fmt.Println("input char:", ch)

		if ch == 'q' || ch == 'Q' {
			// 结束游戏
			fmt.Println("Quit")
		} else if ch != 'W' && ch != 'w' && ch != 'a' && ch != 'A' && ch != 's' && ch != 'S' && ch != 'D' && ch != 'd' {
			// 输入的不是要操作的数
			continue
		} else {
			Move(ch)
		}
	}
}

func Move(ch uint8) { //根据 ch 移动 data
	mu.Lock()
	defer mu.Unlock()

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
		fmt.Println("Input Error!")
	}

	var nx int = data.X + dx
	var ny int = data.Y + dy
	fmt.Println("(nx, ny): ", nx, ny)

	if nx < 0 || nx >= 8 || ny < 0 || ny >= 8 {
		return // 已经 defer 了，这里应该不用 unlock 了
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
		usermd5 = usermd5 + "|" + username
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
