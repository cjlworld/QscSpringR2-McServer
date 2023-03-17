package main

import (
	"bufio"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"time"
)

const Version string = "1.12.2" // 版本号

type McServer struct {
}

type DataPack struct {
	Id    int   // 操作序号
	Opt   uint8 //
	X, Y  int
	Mymap [8][8]int
}

func (this *McServer) Login(md5str string, reply *DataPack) error { // 登录函数
	//fmt.Println(md5str)
	//fmt.Println(usermd5)
	server_md5str := usermd5 + "|" + Version // 加上 Version 判断
	if md5str != server_md5str {
		reply.Opt = 'W'
		fmt.Println("Wrong!")
		return nil
	}
	*reply = data
	reply.Opt = 'A'
	reply.Id = 0
	return nil
}

func Compare(a DataPack, b DataPack) bool { // 比较两个数据包除了 option 的部分，相同返回 true，不同返回 false
	if a.Id != b.Id {
		return false
	} else if a.X != b.X || a.Y != b.Y {
		return false
	} else {
		for i := 0; i < 8; i++ {
			for j := 0; j < 8; j++ {
				if a.Mymap[i][j] != b.Mymap[i][j] {
					return false
				}
			}
		}
	}
	return true
}

func CloseConnection() error { // 延迟关闭服务端
	time.Sleep(time.Duration(2) * time.Second) // 两秒后关
	listener.Close()
	return nil
}

// 不知道服务端要不要加锁
// 听说 RPC 服务有缓冲区来着
func (this *McServer) FetchClient(ClientMap DataPack, reply *DataPack) error {
	reply.Opt = 'C'
	// fmt.Println("Receive a package from client!")
	// 根据 id 判断
	// 若相同，则比对
	// 若是比较早的 id，略过
	// 若是下一个 id，且 命令为 WASDQ，保存
	if data.Opt == 'Q' {
		// 游戏已经结束了
		return nil
	} else if ClientMap.Opt == 'Q' { // 结束游戏，保存，以玩家的数据为准
		data = ClientMap
		SaveUserData()       // 保存数据
		go CloseConnection() // 服务端两秒后关闭连接
		return nil
	} else if ClientMap.Id < data.Id { // 过时的数据包，不比对
		return nil
	} else if ClientMap.Id == data.Id { // 验证的数据包，以服务器的为准
		if Compare(data, ClientMap) { // 数据相同
			reply.Opt = 'C' // Correct!
			return nil
		} else { // 数据不同，以服务端为准
			*reply = data
			reply.Opt = 'F' // Fault
			return nil
		}
	} else { // 客户端的数据 更新 服务端的数据
		// 用服务端的数据验算
		// 事实上，有可能 KeepInTouch 的数据包会先发来，屏蔽掉即可 2023/3/17
		if ClientMap.Opt == 'T' {
			return nil
		}
		Move(ClientMap.Opt)
		if Compare(data, ClientMap) { // 数据相同
			reply.Opt = 'C'
			fmt.Printf("Move to (%d, %d)\n", data.X, data.Y)
			return nil
		} else { // 否则客户端的数据错了
			*reply = data
			reply.Opt = 'F' // Fault
			return nil
		}
	}
	return nil
}

var usermd5 string
var data DataPack // 服务器的地图
var listener net.Listener

func main() {
	LoadUserData() // 从文件读入初始数据

	fmt.Println("Starting server...")

	// 监听TCP连接
	listener, err := net.Listen("tcp", ":25565")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	mcs := new(McServer)
	rpc.Register(mcs)
	rpc.Accept(listener)
	listener.Close()
	time.Sleep(time.Duration(5) * time.Second) // 延时关闭
}

func SaveUserData() { // 保存用户数据
	fmt.Println("Saving user data...")

	file, err := os.OpenFile("userdb.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 777) // 读写混合，没有就新建，清空
	if err != nil {
		fmt.Println("Open File error: ", err)
		return
	}
	defer file.Close() // 别忘了关闭文件

	writer := bufio.NewWriter(file)
	fmt.Fprintln(writer, usermd5)
	fmt.Fprintf(writer, "%d %d\n", data.X, data.Y)
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			fmt.Fprintf(writer, "%d ", data.Mymap[i][j])
		}
		fmt.Fprintf(writer, "\n")
	}
	writer.Flush() // 刷新缓存区
	fmt.Println("Save all!")
}

func LoadUserData() { // 读入用户数据和地图
	fmt.Println("Init user data...")

	file, err := os.Open("userdb.txt")
	if err != nil {
		fmt.Println("Open File error: ", err)
		return
	}
	defer file.Close()

	// 创建一个 bufio.Reader 对象，用于逐行读取文件
	reader := bufio.NewReader(file)

	// Read usermd5
	fmt.Fscanf(reader, "%s\n", &usermd5) // fmt 要加 '\n'
	// usermd5 = usermd5 + "|" + Version, 在判断时加
	fmt.Println(usermd5)

	// Read (x,y)
	fmt.Fscanf(reader, "%d %d\n", &data.X, &data.Y)
	fmt.Println("x: ", data.X)
	fmt.Println("y: ", data.Y)

	// Read map
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			fmt.Fscanf(reader, "%d", &data.Mymap[i][j])
		}
		fmt.Fscanf(reader, "\n")
	}
	fmt.Println(data.Mymap)

	data.Id = 0
	fmt.Println("Successfully init userdata!")
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
		fmt.Printf("%c\n", ch)
		return
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
