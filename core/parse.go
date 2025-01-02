package core

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/leiyang23/nmon-parser/util"
	"io"
	"os"
	"strconv"
	"strings"
)

type AAA struct {
	Item1 string
	Item2 string
}

func GetAAA(nmonFile string) (data []AAA, err error) {
	f, err := os.Open(nmonFile)
	if err != nil {
		return
	}

	defer util.Close(f)

	buf := make([]byte, 1024*1024)
	_, _ = f.Read(buf)

	data = make([]AAA, 0)
	for _, line := range bytes.Split(buf, []byte("\n")) {
		if bytes.HasPrefix(line, []byte("AAA")) {
			fields := bytes.SplitN(line, []byte(","), 3)
			data = append(data, AAA{Item1: string(fields[1]), Item2: string(fields[2])})
		}
	}

	return
}

type BBBP struct {
	Index  string
	Cmd    string
	Output string
}

func GetBBBP(nmonFile string) (data []BBBP, err error) {
	f, err := os.Open(nmonFile)
	if err != nil {
		return
	}
	defer util.Close(f)

	r := bufio.NewReader(f)
	data = make([]BBBP, 0)
	for {
		line, err1 := r.ReadString('\n')
		line = strings.TrimSpace(line)
		if err1 != nil && err1 != io.EOF {
			return
		}
		if err1 == io.EOF {
			break
		}

		if strings.HasPrefix(line, "ZZZZ") {
			break
		}

		if !strings.HasPrefix(line, "BBBP") {
			continue
		}

		fields := strings.SplitN(line, ",", 4)
		item := BBBP{
			Index: fields[1],
			Cmd:   fields[2],
		}
		if len(fields) == 4 {
			item.Output = fields[3]
		}
		data = append(data, item)
	}

	return
}

type Category struct {
	UniqueID string
	Name     string
	Metrics  []string
}

func GetAllCategory(nmonFile string) (categories map[string]Category, err error) {
	f, err := os.Open(nmonFile)
	if err != nil {
		return
	}

	defer util.Close(f)

	buf := make([]byte, 1024*1024*4)
	_, _ = f.Read(buf)

	categories = make(map[string]Category)
	for _, line := range bytes.Split(buf, []byte("\n")) {
		if bytes.HasPrefix(line, []byte("AAA")) {
			continue
		}
		if bytes.HasPrefix(line, []byte("BBBP")) {
			break
		}

		fields := bytes.SplitN(line, []byte(","), 3)
		categories[string(fields[0])] = Category{
			UniqueID: string(fields[0]),
			Name:     string(fields[1]),
			Metrics:  strings.Split(string(fields[2]), ","),
		}
	}

	return
}

// CategoryData 提取出特定类型的数据
type CategoryData struct {
	TimeSeries      []string
	MetricLinesData [][]float64
}

func GetCategoryData(fileName string, category Category) (data CategoryData, err error) {
	timeSeries := make([]string, 0)
	metricLinesData := make([][]float64, 0)

	f, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer util.Close(f)

	r := bufio.NewReader(f)

	var sectionCount int
	var serialNumber, recordTime string
	for {
		line, err1 := r.ReadString('\n')
		if err1 != nil && err1 != io.EOF {
			err = fmt.Errorf("read %s err: %v", fileName, err1)
			return
		}
		if err1 == io.EOF {
			break
		}

		line = strings.TrimSpace(line)

		// 从读取到 ZZZZ 行，开始进入指标数据提取
		if strings.HasPrefix(line, "ZZZZ") {
			sectionCount++
			serialNumber = strings.Split(line, ",")[1]
			recordTime = strings.Split(line, ",")[2]
			continue
		}
		// 在没有读取到 ZZZZ 行前，数据都跳过
		if sectionCount <= 0 {
			continue
		}

		fields := strings.Split(line, ",")
		if fields[1] != serialNumber {
			continue
		}

		if fields[0] == category.UniqueID {
			var v float64
			lineData := make([]float64, 0)
			for _, field := range fields[2:] {
				v, err = strconv.ParseFloat(field, 64)
				if err != nil {
					//fmt.Println(22, field)
					//return
				}
				lineData = append(lineData, v)
			}
			metricLinesData = append(metricLinesData, lineData)
			timeSeries = append(timeSeries, recordTime)
		}
	}

	data.TimeSeries = timeSeries
	data.MetricLinesData = metricLinesData

	return
}

// GetCategoryAllData 一次性将文件中的指标数据全提取出来
func GetCategoryAllData(fileName string) (data map[string]*CategoryData, err error) {
	data = make(map[string]*CategoryData)

	f, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer util.Close(f)

	r := bufio.NewReader(f)

	var sectionCount int
	var serialNumber, recordTime string
	for {
		line, err1 := r.ReadString('\n')
		if err1 != nil && err1 != io.EOF {
			err = fmt.Errorf("read %s err: %v", fileName, err1)
			return
		}
		if err1 == io.EOF {
			break
		}

		line = strings.TrimSpace(line)

		// 从读取到 ZZZZ 行，开始进入指标数据提取
		if strings.HasPrefix(line, "ZZZZ") {
			sectionCount++
			serialNumber = strings.Split(line, ",")[1]
			recordTime = strings.Split(line, ",")[2]
			continue
		}
		// 在没有读取到 ZZZZ 行前，数据都跳过
		if sectionCount <= 0 {
			continue
		}

		fields := strings.Split(line, ",")
		if fields[1] != serialNumber {
			continue
		}

		var v float64
		lineData := make([]float64, 0)
		for _, field := range fields[2:] {
			v, err = strconv.ParseFloat(field, 64)
			if err != nil {
				//fmt.Println(22, field)
				//return
			}
			lineData = append(lineData, v)
		}

		if _, ok := data[fields[0]]; !ok {
			data[fields[0]] = &CategoryData{
				TimeSeries:      make([]string, 0),
				MetricLinesData: make([][]float64, 0),
			}
		}
		data[fields[0]].TimeSeries = append(data[fields[0]].TimeSeries, recordTime)
		data[fields[0]].MetricLinesData = append(data[fields[0]].MetricLinesData, lineData)
	}

	return
}
