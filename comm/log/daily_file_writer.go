package log

import (
	"github.com/pkg/errors"
	"io"
	"os"
	"path"
	"sync"
	"time"
)

type dailyFileWriter struct {
	// 日志文件名称
	fileName string
	// 上一次写入日期
	lastYearDay int
	// 输出文件
	outputFile *os.File

	fileSwitchLock *sync.Mutex
}

func (w *dailyFileWriter) Write(byteArray []byte) (n int, err error) {
	if nil == byteArray ||
		len(byteArray) <= 0 {
		return 0, nil
	}

	outputFile, err := w.getOutputFile()

	if err != nil {
		return 0, err
	}

	_, _ = os.Stderr.Write(byteArray) // console write
	_, _ = outputFile.Write(byteArray)

	return len(byteArray), nil
}

// 这里可能出现一个并发问题
// 获取输出文件  每天输出一个新的日志文件
func (w *dailyFileWriter) getOutputFile() (io.Writer, error) {
	yearDay := time.Now().YearDay()

	if w.lastYearDay == yearDay &&
		nil != w.outputFile {
		return w.outputFile, nil
	}

	w.fileSwitchLock.Lock()
	defer w.fileSwitchLock.Unlock()

	if w.lastYearDay == yearDay &&
		nil != w.outputFile {
		return w.outputFile, nil
	}

	w.lastYearDay = yearDay

	// 先建立日志目录
	err := os.MkdirAll(path.Dir(w.fileName), os.ModePerm)

	if err != nil {
		return nil, err
	}
	// 定义日志文件名称 = 日志文件名. 日期后缀
	newDailyFile := w.fileName + "." + time.Now().Format("20060102")

	// 打开文件
	outputFile, err := os.OpenFile(
		newDailyFile,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644, // rw-r--r--
	)

	//如果打开的文件失败 就直接不往下走了
	if err != nil ||
		nil == outputFile {
		return nil, errors.Errorf("打开文件 %s 失败, err = %v", newDailyFile, err)
	}
	// 将原始的文件关闭
	if nil != w.outputFile {
		_ = w.outputFile.Close()
	}
	// 将原来的引用赋值给dailyFileWrite
	w.outputFile = outputFile
	return outputFile, nil
}
