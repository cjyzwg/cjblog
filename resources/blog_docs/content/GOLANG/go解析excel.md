---
title: go解析excel
date: 2022-05-01 15:27:00
categories:
  - GOLANG
---
#### 时间格式转化
1. import ("github.com/360EntSecGroup-Skylar/excelize/v2")
2.  go.mod require (
   github.com/360EntSecGroup-Skylar/excelize/v2 v2.0.2)
3. 
```go
style, _ = xlsx.NewStyle(`{"number_format":15}`)
xlsx.SetCellStyle("Sheet1", "E7", "E7", style)
e7 := xlsx.GetCellValue("Sheet1", "E7")
println(e7) // 29-Jun-17
col,_:=excelize.ColumnNumberToName(5)//将索引转成字母
style, _ := xlsx.NewStyle(`{"number_format":15}`)
xlsx.SetCellStyle(oldsheet, col+strconv.Itoa(keyrow), col+strconv.Itoa(keyrow), style)
formatday ,_:=xlsx.GetCellValue(oldsheet, col+strconv.Itoa(keyrow))
```

#### 时区问题：
```go
var loc = time.FixedZone("CST", 8*3600)       // 东八
if runtime.GOOS !="windows" {
   loc, _ = time.LoadLocation("Asia/Shanghai") //获取时区
}

log.Println("2",loc,runtime.GOOS)
starttmp, _ := time.ParseInLocation(timeLayout, start, loc)


时间格式化输出的时区为东八区北京时间，无需关系系统所在时区。
在Go语言上，go语言的time.Now()返回的是当地时区时间，直接用：
time.Now().Format("2006-01-02 15:04:05")

输出的是当地时区时间。
go语言并没有全局设置时区这么一个东西，每次输出时间都需要调用一个In()函数改变时区：
var cstSh, _ = time.LoadLocation("Asia/Shanghai") //上海
fmt.Println("SH : ", time.Now().In(cstSh).Format("2006-01-02 15:04:05"))

在windows系统上，没有安装go语言环境的情况下，time.LoadLocation会加载失败。
var cstZone = time.FixedZone("CST", 8*3600)       // 东八
fmt.Println("SH : ", time.Now().In(cstZone).Format("2006-01-02 15:04:05"))

最好的办法是用time.FixedZone
```

#### 获取单元格的背景颜色：
```go
  package main
  import (
   "fmt"
   "github.com/360EntSecGroup-Skylar/excelize"
  )
  func main() {
    xlsx, _ := excelize.OpenFile("Book1.xlsx")
    fmt.Println(getCellBgColor(xlsx, "Sheet1", "C1"))
  }
  func getCellBgColor(xlsx *excelize.File, sheet, axix string) string {
    styleID := xlsx.GetCellStyle(sheet, axix)
    fillID := xlsx.Styles.CellXfs.Xf[styleID].FillID
    fgColor := xlsx.Styles.Fills.Fill[fillID].PatternFill.FgColor
    if fgColor.Theme != nil {
        srgbClr := xlsx.Theme.ThemeElements.ClrScheme.Children[*fgColor.Theme].SrgbClr.Val
        return excelize.ThemeColor(srgbClr, fgColor.Tint)
    }
    return fgColor.RGB
  }
  ```