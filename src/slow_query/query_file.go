package slow_query

import (
	"bufio"
	"fmt"
	log "github.com/wfxiang08/cyutils/utils/rolling_log"
	"os"
	"time"
)

// 将Rds的Slow Query合并成为一个本地的，符合mysqldumpslow脚本的文件
func SaveToDefaultFormat(contents []string) (string, error) {
	t := time.Now().UnixNano()
	tmpFile := fmt.Sprintf("%d.tmp", t)
	f, err := os.Create(tmpFile)
	if err != nil {
		log.ErrorError(err, "err when open query summary to write")
		return "", err
	}
	// defer os.Remove(tmpFile)

	// 合并Query
	writer := bufio.NewWriter(f)
	for _, content := range contents {
		writer.WriteString(content)
		writer.WriteString("\n")
	}
	err = writer.Flush()
	if err != nil {
		log.ErrorError(err, "err when open query summary to write")
		return "", err
	}
	f.Close()
	return tmpFile, nil
}
