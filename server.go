package main

import (
	"bufio"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"strings"
)

type McServer struct {
}

type DataPack struct {
	Id    int   // 操作序号
	Opt   uint8 //
	X, Y  int
	Mymap [8][8]int
}

func (this *McServer) Login(md5str string, pkg *DataPack) error {
	//fmt.Println(md5str)
	//fmt.Println(usermd5)
	if md5str != usermd5 {
		pkg.Opt = 'W'
		fmt.Println("Wrong!")
		return nil
	}
	*pkg = data
	pkg.Opt = 'A'
	pkg.Id = 0
	return nil
}

var usermd5 string
var data DataPack

func main() {
	InitUserData() // 从文件读入初始数据

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
}

func InitUserData() {
	fmt.Println("Init user data...")

	file, err := os.Open("userdb.txt")
	if err != nil {
		fmt.Println("Open File error: ", err)
		return
	}
	defer file.Close()

	// 创建一个 bufio.Scanner 对象，用于逐行读取文件
	scanner := bufio.NewScanner(file)

	// Read usermd5
	scanner.Scan()
	usermd5 = scanner.Text()
	fmt.Println(usermd5)

	// Read (x,y)
	scanner.Scan()
	parts := strings.Split(scanner.Text(), " ")
	data.X, _ = strconv.Atoi(parts[0])
	data.Y, _ = strconv.Atoi(parts[1])
	fmt.Println("x: ", data.X)
	fmt.Println("y: ", data.Y)

	// Read map
	var i = 0
	for scanner.Scan() {
		line := scanner.Text()

		var j = 0
		for _, s := range strings.Split(line, " ") {
			//fmt.Println(s)
			if s == "" {
				continue
			}
			num, _ := strconv.Atoi(s)
			data.Mymap[i][j] = num
			j++
		}
		i++
	}
	// fmt.Println(data.Mymap)

	data.Id = 0
	fmt.Println("Successfully init userdata!")
}
