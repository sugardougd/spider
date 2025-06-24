package commands

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"io"
	"math/rand"
	spider2 "spider"
	"strconv"
	"strings"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func noyaCommandFlags(f *spider2.Flags) {
	f.Int(&spider2.Flag{Short: "l", Long: "low", Help: "the low value", Default: 10})
	f.Int(&spider2.Flag{Short: "t", Long: "top", Help: "the top value", Default: 10})
	f.Int(&spider2.Flag{Short: "c", Long: "count", Help: "the count of problem", Default: 10})
	f.String(&spider2.Flag{Short: "f", Long: "file", Help: "the excel file to save", Default: ""})
}

func noyaCommandRun(c *spider2.Context, operators []mathOperator) error {
	low, err := c.FlagValues.Int(c.Command, "low")
	if err != nil {
		return err
	}
	top, err := c.FlagValues.Int(c.Command, "top")
	if err != nil {
		return err
	}
	count, err := c.FlagValues.Int(c.Command, "count")
	if err != nil {
		return err
	}
	file, err := c.FlagValues.String(c.Command, "file")
	if err != nil {
		return err
	}
	out := make(chan string)
	param := mathParams{low, top}
	go calculate(count, operators, param, out)
	if len(file) > 0 {
		saveAsExcel(file, out, c.Spider)
	} else {
		for o := range out {
			c.Spider.Println(o)
		}
	}
	return nil
}

func NoyaCommand() *spider2.Command {
	noya := &spider2.Command{
		Name:        "noya",
		Description: "noya tools",
		Usage:       "noya...",
		Run: func(c *spider2.Context) error {
			c.Spider.PrintCommandHelp(c.Command)
			return nil
		},
	}
	// add
	add := &spider2.Command{
		Name:        "math-add",
		Description: "displays math-add problem",
		Usage:       "noya math-add",
		Flags:       noyaCommandFlags,
		Run: func(c *spider2.Context) error {
			opers := []mathOperator{mathOperatorFunc(mathAddCalculate)}
			return noyaCommandRun(c, opers)
		},
	}

	// keep-add
	addKeep := &spider2.Command{
		Name:        "math-add-keep",
		Description: "displays math-add-keep problem",
		Usage:       "noya math-add-keep",
		Flags:       noyaCommandFlags,
		Run: func(c *spider2.Context) error {
			opers := []mathOperator{mathOperatorFunc(mathKeepAddCalculate)}
			return noyaCommandRun(c, opers)
		},
	}

	// sub
	sub := &spider2.Command{
		Name:        "math-sub",
		Description: "displays math-sub problem",
		Usage:       "noya math-sub",
		Flags:       noyaCommandFlags,
		Run: func(c *spider2.Context) error {
			opers := []mathOperator{mathOperatorFunc(mathSubCalculate)}
			return noyaCommandRun(c, opers)
		},
	}

	// keep-sub
	subKeep := &spider2.Command{
		Name:        "math-sub-keep",
		Description: "displays math-sub-keep problem",
		Usage:       "noya math-sub-keep",
		Flags:       noyaCommandFlags,
		Run: func(c *spider2.Context) error {
			opers := []mathOperator{mathOperatorFunc(mathKeepSubCalculate)}
			return noyaCommandRun(c, opers)
		},
	}

	// add & sub
	addSub := &spider2.Command{
		Name:        "math-add-sub",
		Description: "displays math-add & math-sub problem",
		Usage:       "noya math-add-sub",
		Flags:       noyaCommandFlags,
		Run: func(c *spider2.Context) error {
			opers := []mathOperator{mathOperatorFunc(mathAddCalculate), mathOperatorFunc(mathSubCalculate)}
			return noyaCommandRun(c, opers)
		},
	}

	// keep-add & keep-sub
	addSubKeep := &spider2.Command{
		Name:        "math-add-sub-keep",
		Description: "displays math-addkeep & math-subkeep problem",
		Usage:       "noya math-add-sub-keep",
		Flags:       noyaCommandFlags,
		Run: func(c *spider2.Context) error {
			opers := []mathOperator{mathOperatorFunc(mathKeepAddCalculate),
				mathOperatorFunc(mathKeepSubCalculate),
				mathOperatorFunc(mathKeepAddSubCalculate)}
			return noyaCommandRun(c, opers)
		},
	}

	noya.AddCommand(add)
	noya.AddCommand(addKeep)
	noya.AddCommand(sub)
	noya.AddCommand(subKeep)
	noya.AddCommand(addSub)
	noya.AddCommand(addSubKeep)
	return noya
}

// saveAsExcel 将结果写到excel
func saveAsExcel(file string, data chan string, w io.Writer) {
	if !strings.HasPrefix(file, "/") {
		file = "/Users/leiqian/Downloads/" + file
	}
	excel := excelize.NewFile()
	defer func() {
		if err := excel.Close(); err != nil {
			fmt.Fprintln(w, err)
		}
	}()
	cols := []string{"A", "B"}
	index, row := 0, 1
	// 设置列宽度
	excel.SetColWidth("Sheet1", cols[0], cols[len(cols)-1], 37)
	// 设置单元格的值
	for d := range data {
		if index == len(cols) {
			index = 0
			// 设置行高度
			excel.SetRowHeight("Sheet1", row, 34)
			row++
		}
		excel.SetCellValue("Sheet1", cols[index]+strconv.Itoa(row), d)
		index++
	}

	// 设置单元格的样式
	style, err := excel.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   false,
			Family: "宋体",
			Size:   24,
			Color:  "#000000",
		},
		Border: []excelize.Border{
			{
				Type:  "left",
				Color: "#000000",
				Style: 3,
			},
			{
				Type:  "right",
				Color: "#000000",
				Style: 3,
			},
			{
				Type:  "top",
				Color: "#000000",
				Style: 3,
			},
			{
				Type:  "bottom",
				Color: "#000000",
				Style: 3,
			},
		},
	})
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	excel.SetCellStyle("Sheet1", cols[0]+"1", cols[1]+strconv.Itoa(row), style)

	// 根据指定路径保存文件
	if err := excel.SaveAs(file); err != nil {
		fmt.Fprintln(w, err)
	} else {
		fmt.Fprintln(w, "Saved to", file)
	}
}

func calculate(c int, operators []mathOperator, param mathParams, out chan string) {
	for i := 0; i < c; i++ {
		out <- operators[r.Intn(len(operators))].calculate(param)
	}
	close(out)
}

type mathParams struct {
	low int
	top int
}

type mathOperator interface {
	calculate(param mathParams) string
}

type mathOperatorFunc func(param mathParams) string

func (oper mathOperatorFunc) calculate(param mathParams) string {
	return oper(param)
}

func mathAddCalculate(param mathParams) string {
	num1 := r.Intn(param.top+1-param.low) + param.low // [low,top]
	num2 := r.Intn(num1-1) + 1                        // [1, num1-1]
	num3 := num1 - num2
	return fmt.Sprintf("%d + %d = ", num2, num3)
}

func mathKeepAddCalculate(param mathParams) string {
	num1 := r.Intn(param.top+1-param.low) + param.low // [low,top]
	num2 := r.Intn(num1-3) + 2                        // [2, num1-2]
	num3 := r.Intn(num1-num2-1) + 1
	num4 := num1 - num2 - num3
	return fmt.Sprintf("%d + %d + %d = ", num2, num3, num4)
}

func mathSubCalculate(param mathParams) string {
	num1 := r.Intn(param.top+1-param.low) + param.low // [low,top]
	num2 := r.Intn(num1-1) + 1                        // [1, num1-1]
	return fmt.Sprintf("%d - %d = ", num1, num2)
}

func mathKeepSubCalculate(param mathParams) string {
	num1 := r.Intn(param.top+1-param.low) + param.low // [low,top]
	num2 := r.Intn(num1-3) + 2                        // [2, num1-2]
	num3 := r.Intn(num1-num2-1) + 1
	return fmt.Sprintf("%d - %d - %d = ", num1, num2, num3)
}

func mathKeepAddSubCalculate(param mathParams) string {
	num1 := r.Intn(param.top+1-param.low) + param.low // [low,top]
	num2 := r.Intn(num1-3) + 2                        // [2, num1-2]
	num3 := r.Intn(num1-num2-1) + 1
	num4 := num1 - num2 - num3
	t := r.Intn(4)
	if t == 0 {
		return fmt.Sprintf("%d + %d + %d = ", num2, num3, num4)
	}
	if t == 1 {
		return fmt.Sprintf("%d + %d - %d = ", num2, num3+num4, num4)
	}
	if t == 2 {
		return fmt.Sprintf("%d - %d + %d = ", num2+num3, num3, num4)
	}
	return fmt.Sprintf("%d - %d - %d = ", num1, num2, num3)
}
