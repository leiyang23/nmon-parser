package core

import (
	"fmt"
	"github.com/leiyang23/nmon-parser/util"
	"github.com/xuri/excelize/v2"
	"log/slog"
	"runtime/debug"
	"sort"
)

func AddToExcel(excelFile string, categories map[string]Category, data map[string]*CategoryData) (err error) {

	f, err := excelize.OpenFile(excelFile)
	if err != nil {
		f = excelize.NewFile()
	}
	defer util.Close(f)

	// 对 sheet 名称进行排序
	catList := make([]string, 0)
	for category := range data {
		catList = append(catList, category)
	}
	sort.Strings(catList)

	// 事先创建好需要的 sheet，方便下一步使用 写入流 写入大规模数据
	for _, category := range catList {
		_, err = f.NewSheet(category)
	}

	var sw *excelize.StreamWriter
	for _, category := range catList {
		categoryData := data[category]

		sw, err = f.NewStreamWriter(category)
		if err != nil {
			slog.Error(err.Error())
			return
		}

		header := make([]interface{}, 0)
		header = append(header, "time")
		for _, i := range categories[category].Metrics {
			header = append(header, i)
		}
		err = sw.SetRow("A1", header)
		if err != nil {
			return
		}

		for i, row := range categoryData.MetricLinesData {
			rowData := make([]interface{}, 0, len(header))
			rowData = append(rowData, categoryData.TimeSeries[i])
			for _, j := range row {
				rowData = append(rowData, j)
			}
			err = sw.SetRow(fmt.Sprintf("A%d", i+2), rowData)
			if err != nil {
				return
			}
		}

		// 写入磁盘
		err = sw.Flush()
		if err != nil {
			return err
		}
	}

	f.SetActiveSheet(0)
	_ = f.DeleteSheet("Sheet1")
	return f.SaveAs(excelFile)
}

func AddAAAToExcel(excelFile string, data []AAA) (err error) {
	f, err := excelize.OpenFile(excelFile)
	if err != nil {
		f = excelize.NewFile()
	}
	defer util.Close(f)

	_, err = f.NewSheet("AAA")
	if err != nil {
		return
	}

	for i, row := range data {
		err = f.SetSheetRow("AAA", fmt.Sprintf("A%d", i+1), &[]string{row.Item1, row.Item2})
		if err != nil {
			return
		}
	}

	return f.SaveAs(excelFile)
}

func AddBBBPToExcel(excelFile string, data []BBBP) (err error) {
	f, err := excelize.OpenFile(excelFile)
	if err != nil {
		f = excelize.NewFile()
	}
	defer util.Close(f)

	_, err = f.NewSheet("BBBP")
	if err != nil {
		slog.Error(err.Error())
		return
	}

	for i, row := range data {
		err = f.SetSheetRow("BBBP", fmt.Sprintf("A%d", i+1), &[]string{row.Cmd, row.Output})
		if err != nil {
			return
		}
	}

	debug.FreeOSMemory()

	return f.SaveAs(excelFile)
}
