package mytrace

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strings"
)

func init() {
	// 设置日志前缀
	log.SetPrefix("TRACE: ")
	// 设置日志输出位置为标准输出
	log.SetOutput(new(Logger))

	// 打印调用链
	// trace()
}

// Logger 结构体实现了 io.Writer 接口的 Write 方法
type Logger struct{}

func (l *Logger) Write(data []byte) (int, error) {
	fmt.Print(string(data))
	// 在这里你可以将日志信息输出到任何地方，比如文件、终端等
	return len(data), nil
}

var (
	deep     int  = 4
	fullinfo bool = false
)

func SetTraceDeep(d int) {
	deep = d
}

func SetFullInfo(b bool) {
	fullinfo = b
}

// trace 打印调用链
func Trace() {
	// 遍历调用栈，跳过前两个调用者（trace 和 runtime.Callers）
	// 这里假设要打印的是第三层调用
	log.Println("begin trace")
	for i := 1; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		if i >= deep {
			break
		}
		if fullinfo {
			// 获取调用函数名
			funcName := runtime.FuncForPC(pc).Name()

			log.Printf("Function: %s, File: %s, Line: %d\n", funcName, file, line)
		} else {
			// 获取调用函数名
			// funcName := runtime.FuncForPC(pc).Name()
			funcName := filepath.Base(runtime.FuncForPC(pc).Name())
			// 提取函数名部分
			if index := strings.LastIndex(funcName, "."); index != -1 {
				funcName = funcName[index+1:]
			}
			file = filepath.Base(file)
			log.Printf("Function: %s, File: %s, Line: %d\n", funcName, file, line)
		}

	}
	log.Println("end trace")
}
