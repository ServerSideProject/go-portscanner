package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

var sb strings.Builder

func append(data string) {
	f, _ := os.OpenFile("out.txt", os.O_APPEND|os.O_WRONLY, 0644)
	f.WriteString(data)
	f.Close()

}

func SetConsoleTitle(title string) (int, error) {
	handle, err := syscall.LoadLibrary("Kernel32.dll")
	if err != nil {
		return 0, err
	}
	defer syscall.FreeLibrary(handle)
	proc, err := syscall.GetProcAddress(handle, "SetConsoleTitleW")
	if err != nil {
		return 0, err
	}
	r, _, err := syscall.Syscall(proc, 1, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))), 0, 0)
	return int(r), err
}

func IsOpened(host string, port int) bool {

	timeout := 5 * time.Second
	target := fmt.Sprintf("%s:%d", host, port)

	conn, err := net.DialTimeout("tcp", target, timeout)
	if err != nil {
		return false
	}

	if conn != nil {
		conn.Close()
		return true
	}

	return false
}

func check(ip string, port int, wg *sync.WaitGroup, id int) {
	defer wg.Done()
	openport := strconv.Itoa(port)
	if IsOpened(ip, port) {
		sb.WriteString(ip + ":" + openport + "\n")
		append(ip + ":" + openport + "\n")
	}
}

func main() {

	filename := filepath.Base(os.Args[0])
	if len(os.Args) != 4 {
		fmt.Println("ServerSide PortScanner\nUsage: " + filename + " <ip> <start port> <end port>")
		os.Exit(0)

	}
	host := os.Args[1]
	min, _ := strconv.Atoi(os.Args[2])
	max, _ := strconv.Atoi(os.Args[3])
	max = max + 1
	SetConsoleTitle("ServerSide PortScanner | ip: " + os.Args[1] + " | range ports: " + os.Args[2] + "-" + os.Args[3])
	fmt.Println("ServerSide PortScanner\nStarting")
	var wg sync.WaitGroup
	var a int

	for i := min; i < max; i++ {
		a++
		fmt.Println("Scanning... Current target: "+host, i)
		wg.Add(1)

		go check(host, i, &wg, a)
		time.Sleep(1000 * time.Nanosecond)
	}
	wg.Wait()

	for i := 1; i < 1500; i++ {
		fmt.Println("")
	}
	fmt.Print(sb.String())

}
