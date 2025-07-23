package cleaner

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func StartFileCleanupTask(FolderToClean []string, maxAge time.Duration) {
	go func() {
		ticker := time.NewTicker(time.Hour) // 每小时检查一次
		defer ticker.Stop()

		for {
			<-ticker.C
			now := time.Now()
			for _, folder := range FolderToClean {
				err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
					if err != nil || info.IsDir() {
						return nil
					}
					if now.Sub(info.ModTime()) > maxAge {
						err := os.Remove(path)
						if err != nil {
							fmt.Println("Error deleting file:", err.Error())
						}
						fmt.Println("Deleted expired file:", path)
					}
					return nil
				})
				if err != nil {
					fmt.Println("Error walking directory:", err.Error())
				}
			}
		}
	}()
}
