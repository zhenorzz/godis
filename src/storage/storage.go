package storage

import (
	"os"
	"encoding/csv"
	"router"
)
var file = "test.csv"
func Read(cache *router.Cache) error {
	var f *os.File
	var err error
	_, err = os.Stat(file)
	if err == nil {
		f, err = os.OpenFile(file, os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		reader := csv.NewReader(f)
		records, err := reader.ReadAll()
		if err == nil {
			for _, slice := range records {
				if len(slice) == 2 {
					cache.Data[slice[0]] = slice[1]
				}
			}
		}
	} else {
		f, err = os.Create(file)//创建文件
		if err != nil {
			return err
		}
		f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	}
	defer f.Close()
	return nil
}

func Write(data map[string]string) {
	f, err:= os.OpenFile(file, os.O_APPEND, 0666)
	if err != nil {
		return
	}
	w := csv.NewWriter(f)
	for k, v := range data {
		w.Write([]string{k, v})
	}
	w.Flush()
}

func Update(data []string) {
	f, err:= os.OpenFile(file, os.O_RDONLY, 0666)
	if err != nil {
		return
	}
	r := csv.NewReader(f)
	records, _ := r.ReadAll()
	for k, v := range records {
		if v[0] == data[0] {
			records[k][1] = data[1]
		}
	}
	f.Close()
	f,err = os.Create(file)//创建文件
	if err != nil {
		return
	}
	w := csv.NewWriter(f)
	w.WriteAll(records)
	f.Close()
}

func Delete(data string) {
	f, err:= os.OpenFile(file, os.O_RDONLY, 0666)
	if err != nil {
		return
	}
	r := csv.NewReader(f)
	records, _ := r.ReadAll()
	f.Close()
	f,err = os.Create(file)//创建文件
	if err != nil {
		return
	}
	w := csv.NewWriter(f)
	for _, v := range records {
		if v[0] != data {
			w.Write(v)

		}
	}
	w.Flush()
	f.Close()
}