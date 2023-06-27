package charts

import (
	"math/rand"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// generate random data for bar chart
func generateBarItems() []opts.BarData {
	items := make([]opts.BarData, 0)
	for i := 0; i < 7; i++ {
		items = append(items, opts.BarData{Value: rand.Intn(300)})
	}
	return items
}

func TestChart() {
	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "owo",
		Subtitle: "no",
	}))
	bar.SetXAxis([]string{"owo", "uwu", "uwo", "owu"}).
		AddSeries("Category A", generateBarItems()).
		AddSeries("Category B", generateBarItems())

	f, _ := os.Create("bar.html")
	bar.Render(f)
}
func TestELOChart() {

}
