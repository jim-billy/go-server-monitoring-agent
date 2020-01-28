// Copyright 2013 Ardan Studios. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package helper provides helper functions for logging and catching panics
// Increase file descriptors
//		https://www.tecmint.com/increase-set-open-file-limits-in-linux/
//		https://linoxide.com/linux-how-to/03-methods-change-number-open-file-limit-linux/

package util

import (
	"fmt"
	"runtime"
	"time"
	"io"
	"encoding/binary"
	"net"
	"syscall"
	"reflect"
	"os"
	"os/exec"
	"log"
	"strings"
)

// loggingOn is a simple flag to turn logging on or off.
var loggingOn = true

// TurnLoggingOff sets the logging flag to off.
func TurnLoggingOff() {
	loggingOn = false
}

func Write(s net.Conn,data []byte) error {
	buf := make([]byte, 4+len(data))                       // 4 字节头部 + 数据长度
	binary.BigEndian.PutUint32(buf[:4], uint32(len(data))) // 写入头部
	copy(buf[4:], data)                                    // 写入数据
	_, err := s.Write(buf)
	if err != nil {
		return err
	}
	return nil
}


func Read(s net.Conn)([]byte, error) {
	header := make([]byte, 4)
	_, err := io.ReadFull(s, header)
	if err != nil {
		return nil, err
	}
	dataLen := binary.BigEndian.Uint32(header)
	data := make([]byte, dataLen)
	_, err = io.ReadFull(s, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func SetLimit() {
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
	rLimit.Cur = rLimit.Max
	fmt.Println("Current fd limit : ",rLimit.Cur,", Max fd limit : ",rLimit.Max)
// 	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
// 		panic(err)
// 	}
}

func WriteToFile(filename string, data string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    _, err = io.WriteString(file, data)
    if err != nil {
        return err
    }
    return file.Sync()
}

func AppendToFile(filename string, data string) {     
    file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        log.Fatalf("failed opening file: %s", err)
    }
    defer file.Close()
 
    len, err := file.WriteString(" The Go language was conceived in September 2007 by Robert Griesemer, Rob Pike, and Ken Thompson at Google.")
    if err != nil {
        log.Fatalf("failed writing to file: %s", err)
    }
    fmt.Printf("\nLength: %d bytes", len)
    fmt.Printf("\nFile Name: %s", file.Name())
}


func GetValueFromMap(key string, mapObj map[string]interface{}){

    val := reflect.ValueOf(mapObj)
    fmt.Println("VALUE = ", val)
    fmt.Println("KIND = ", val.Kind())

    if val.Kind() == reflect.Map {
        for _, e := range val.MapKeys() {
            v := val.MapIndex(e)
            //fmt.Println("value type = ", v.Interface())
            switch t := v.Interface().(type) {
            case int:
                fmt.Println(e, t)
            case string:
                fmt.Println(e, t)
            case bool:
                fmt.Println(e, t)
            /*case *url.UrlStats:
                fmt.Println("=============== ", e, t)
            */
            default:
                fmt.Println("not found")

            }
        }
    }
}

func CountOpenFiles() int {
  out, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("lsof -p %v", os.Getpid())).Output()
  if err != nil {
   log.Fatal(err)
  }
  lines := strings.Split(string(out), "\n")
  return len(lines) - 1
}

func fileExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}

func NowAsUnixMilli() int64 {
    return time.Now().UnixNano() / 1e6
}

func CheckAndCreateDirectory(dirName string) error {
    err := os.MkdirAll(dirName, 0777)
    if err == nil || os.IsExist(err) {
        return nil
    } else {
    	fmt.Println("Directory creation failed with error: " + err.Error())
        os.Exit(1)
    }
    return nil
}

// CatchPanic is used to catch any Panic and log exceptions to Stdout. It will also write the stack trace.
func CatchPanic(err *error, goRoutine string, functionName string) {
	if r := recover(); r != nil {
		// Capture the stack trace
		buf := make([]byte, 10000)
		runtime.Stack(buf, false)

		logException(goRoutine, functionName, "PANIC Defered [%v] : Stack Trace : %v", r, string(buf))

		if err != nil {
			*err = fmt.Errorf("%v", r)
		}
	}
}

// writeStdout is used to write a system message directly to stdout.
func writeStdout(goRoutine string, functionName string, message string) {
	fmt.Printf("======================== %s : %s : %s\n", goRoutine, functionName, message)
}

// writeStdoutf is used to write a formatted system message directly stdout.
func logException(goRoutine string, functionName string, format string, a ...interface{}) {
	writeStdout(goRoutine, functionName, fmt.Sprintf(format, a...))
}

