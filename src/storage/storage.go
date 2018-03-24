package storage

import (
	"os"
	"encoding/csv"
	"router"
)

func New(file string, cache *router.Cache) (*csv.Writer, error) {
	var f *os.File
	var w *csv.Writer
	_, err := os.Stat(file)
	if err == nil {
		f, err := os.OpenFile(file, os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		reader := csv.NewReader(f)
		records, err := reader.ReadAll()
		if err != nil {
			return nil, err
		}
		for _, slice := range records {
			if len(slice) == 2 {
				cache.Data[slice[0]] = slice[1]
			}
		}
		w = csv.NewWriter(f)
	} else {
		f, err := os.Create(file)//创建文件
		if err != nil {
			return nil, err
		}
		f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
		w = csv.NewWriter(f)//创建一个新的写入文件流
	}
	defer f.Close()
	return w, nil
}

func Write(w *csv.Writer, message map[string]string) {
	for k, v := range message {
		w.Write([]string{k, v})
	}
	w.Flush()
}
