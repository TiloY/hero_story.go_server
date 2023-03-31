package log

import "os"

type dailyFileWriter struct {
}

func (w *dailyFileWriter) Write(byteArray []byte) (n int, err error) {
	if nil == byteArray ||
		len(byteArray) <= 0 {
		return 0, nil
	}

	_, _ = os.Stdout.Write(byteArray)

	return len(byteArray), nil
}
