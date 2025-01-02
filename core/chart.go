package core

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/render"
	"github.com/go-echarts/go-echarts/v2/types"
)

func LineChart(srcFile string, category Category) (chart render.ChartSnippet, err error) {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros, Width: "1800px", Height: "800px"}),
		charts.WithTitleOpts(
			opts.Title{
				Title:    category.UniqueID,
				Subtitle: "Line chart",
				Bottom:   "0",
			}),
		charts.WithDataZoomOpts(opts.DataZoom{}),
		charts.WithToolboxOpts(
			opts.Toolbox{
				Show:   opts.Bool(true),
				Orient: "horizontal",
				Feature: &opts.ToolBoxFeature{
					SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{Show: opts.Bool(true), Title: "Save"},
				},
			},
		),
	)

	var categoryData CategoryData

	categoryData, err = GetCategoryData(srcFile, category)
	if err != nil {
		return
	}

	// X 轴为时间轴
	line.SetXAxis(categoryData.TimeSeries)

	selectedLegend := make(map[string]bool)
	for i, flag := range category.Metrics {
		lineData := make([]opts.LineData, 0, len(categoryData.MetricLinesData))
		for _, metricLineData := range categoryData.MetricLinesData {
			// fixme: 当出现采集行数据数量小于指标数时，该如何处理？直接丢弃本行吗？还是补0？
			if len(metricLineData) < len(category.Metrics) {
				metricLineData = append(metricLineData, make([]float64, len(category.Metrics)-len(metricLineData))...)
			}

			lineData = append(lineData, opts.LineData{Value: metricLineData[i]})
		}
		line.AddSeries(flag, lineData)
		selectedLegend[flag] = i == 0
	}

	line.SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true)}))
	line.SetGlobalOptions(charts.WithLegendOpts(opts.Legend{Selected: selectedLegend}))

	chart = line.RenderSnippet()
	return
}
