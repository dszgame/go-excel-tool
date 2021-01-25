package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

// 字段数据
type stDataInfo struct {
	sName       string // 字段名
	sVaueType   string // 字段值类型 int/string/audo
	nArrayLevel int    // 几维数组 table
	bKey        bool   // 是否为key
	sTips       string // 说明
}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func pause() {
	cmd := exec.Command("cmd", "/c", "pause")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func parseValueType() {

}

/**
 * 导出EXCEL到json或者lua
 */
func excelExport(file string, sheetName string, outjson string, outlua string) (int, error) {
	xlsx, err := excelize.OpenFile(file)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	rows := xlsx.GetRows(sheetName)
	colNum := len(rows[1])

	// 前3行是字段定义
	nameLists := make([]*stDataInfo, colNum)
	keyname := ""
	tipsstr := fmt.Sprintf("---@class %v\n", sheetName)
	for j := 0; j < colNum; j++ {
		data := &stDataInfo{
			sTips: rows[0][j], // 中文注释

			sName:       rows[2][j], // 第三行是字段名
			nArrayLevel: 0,
		}

		typeList := strings.Split(rows[1][j], "+") // 第二行是指类型
		typestr := ""
		for _, stype := range typeList {
			switch stype {
			case "string", "int", "auto":
				data.sVaueType = stype
				if len(typestr) == 0 {
					typestr = fmt.Sprintf("---@field %v %v %v\n", data.sName, stype, data.sTips)
				}
				break
			case "table":
				data.nArrayLevel++
				typestr = fmt.Sprintf("---@field %v table %v\n", data.sName, data.sTips)
				break
			case "key":
				keyname = data.sName
				data.bKey = true
				break
			}
		}
		tipsstr += typestr
		nameLists[j] = data

		fmt.Println(data.sTips, "\t", data.sName, "\t", data.sVaueType)
	}
	tipsstr += fmt.Sprintf("---@return table<string, %v>\n", sheetName)

	// 解析数据
	jsonstr := "["
	luastr := tipsstr + "return {\n"
	if len(keyname) > 0 {
		jsonstr = "{"
	}
	rowNum := 0
	for rowIndex, row := range rows {
		if rowIndex < 3 {
			continue
		}
		rowstrlua := ""
		rowstr := "{"
		keystr := ""
		keynum := 0
		for colIndex, colCell := range row {
			if colIndex < colNum {
				if len(colCell) > 0 {
					data := nameLists[colIndex]
					if data == nil {
						continue
					}
					if data.nArrayLevel == 0 {
						switch data.sVaueType {
						case "string":
							rowstr += fmt.Sprintf(`"%v":"%v",`, data.sName, colCell)
							rowstrlua += fmt.Sprintf(`%v="%v",`, data.sName, colCell)
							if keyname == data.sName {
								keystr = colCell
							}
							break
						case "int", "auto":
							num, err := strconv.Atoi(colCell)
							if err == nil { // 数值转换错误
								rowstr += fmt.Sprintf(`"%v":%v,`, data.sName, num)
								rowstrlua += fmt.Sprintf(`%v=%v,`, data.sName, num)
								if keyname == data.sName {
									keynum = num
								}
							} else {
								rowstr += fmt.Sprintf(`"%v":"%v",`, data.sName, colCell)
								rowstrlua += fmt.Sprintf(`%v="%v",`, data.sName, colCell)
								if keyname == data.sName {
									keystr = colCell
								}
							}
							break
						}
					} else if data.nArrayLevel == 1 {
						str := ""
						ary := strings.Split(colCell, "+")
						for _, v := range ary {
							switch data.sVaueType {
							case "string":
								str += fmt.Sprintf(`"%v",`, v)
								break
							case "int", "auto":
								num, err := strconv.Atoi(v)
								if err == nil { // 数值转换错误
									str += fmt.Sprintf(`%v,`, num)
								} else {
									str += fmt.Sprintf(`"%v",`, v)
								}
								break
							}
						}
						if len(str) > 0 {
							str = str[0 : len(str)-1] // 去掉最后一个,
						}
						rowstr += fmt.Sprintf(`"%v":[%v],`, data.sName, str)
						rowstrlua += fmt.Sprintf(`%v={%v},`, data.sName, str)
					} else if data.nArrayLevel == 2 {
						luastr := ""
						jsonstr := ""
						ary1 := strings.Split(colCell, "|")
						for _, v := range ary1 {
							luastr += "{"
							jsonstr += "["
							str := ""
							ary := strings.Split(v, "+")
							for _, v := range ary {
								switch data.sVaueType {
								case "string":
									str += fmt.Sprintf(`"%v",`, v)
									break
								case "int", "auto":
									num, err := strconv.Atoi(v)
									if err == nil { // 数值转换错误
										str += fmt.Sprintf(`%v,`, num)
									} else {
										str += fmt.Sprintf(`"%v",`, v)
									}
									break
								}
							}
							if len(str) > 0 {
								str = str[0 : len(str)-1] // 去掉最后一个,
							}
							luastr += str + "},"
							jsonstr += str + "],"
						}

						if len(jsonstr) > 0 {
							jsonstr = jsonstr[0 : len(jsonstr)-1] // 去掉最后一个,
						}

						rowstr += fmt.Sprintf(`"%v":[%v],`, data.sName, jsonstr)
						rowstrlua += fmt.Sprintf(`%v={%v},`, data.sName, luastr)
					}
				}
			}
		}
		// 处理key
		jsonkeystr := ""

		if len(keystr) == 0 {
			if keynum > 0 {
				keystr = fmt.Sprintf("\t[%v]=", keynum)
				jsonkeystr = fmt.Sprintf(`"%v":`, keynum)
			}
		} else {
			jsonkeystr = fmt.Sprintf(`"%v":`, keystr)
			keystr = "\t" + fmt.Sprintf(`["%v"]=`, keystr)
		}

		// 过滤无效行数据
		if len(rowstrlua) == 0 || (len(keyname) > 0 && len(keystr) == 0 && keynum == 0) {
			continue
		}

		rowstr = rowstr[0 : len(rowstr)-1] // 去掉最后一个,
		jsonstr += jsonkeystr + rowstr + "},"

		luastr += keystr + "{" + rowstrlua + "},\n"
		rowNum++

	}
	jsonstr = jsonstr[0 : len(jsonstr)-1] // 去掉最后一个,

	if len(keyname) > 0 { // 如果有key
		jsonstr += "}"
	} else {
		jsonstr += "]"
	}
	luastr += "}"

	if outjson != "" {
		f, err1 := os.OpenFile(outjson, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
		if err1 != nil {
			return 0, err1
		}
		f.WriteString(jsonstr)
		f.Close()
	}

	if outlua != "" {
		f, err1 := os.OpenFile(outlua, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
		if err1 != nil {
			return 0, err1
		}
		f.WriteString(luastr)
		f.Close()
	}

	return rowNum, nil
}

func main() {
	fConfg, err := excelize.OpenFile("config.xlsx")
	if err != nil {
		fmt.Println(err)
		pause()
		return
	}

	rows := fConfg.GetRows("Sheet1")
	excel_file_path := rows[1][0]
	json_file_path := rows[1][1]
	lua_file_path := rows[1][2]
	for index, row := range rows {
		if index < 3 || row[0] == "" {
			continue
		}
		jsonfile := json_file_path + row[2] + ".json"
		luafile := lua_file_path + row[2] + ".lua"
		if row[3] != "1" {
			jsonfile = ""
		}
		if row[4] != "1" {
			luafile = ""
		}
		fmt.Println("开始导表：", excel_file_path+row[0], "\t", row[1])
		num, err1 := excelExport(excel_file_path+row[0], row[1], jsonfile, luafile)
		if err1 == nil {
			fmt.Println("导表成功：", row[0], "\t", row[1], "\t行数：", num, "\n\n")
		} else {
			fmt.Println("！！！导表失败：", row[0], "\t", row[1], err1.Error())
		}

	}
	pause()
}
