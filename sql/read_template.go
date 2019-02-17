package sql

import (
	"fmt"
	"github.com/beevik/etree"
	"os"
	"path/filepath"
)

// SqlTemplate 讀取 template.sql 內容的物件。
var SqlTemplate = make(map[string]string)

func init() {
	pwd, _ := os.Getwd()
	directory := filepath.Base(pwd)
	sqlTemplatePath := ""

	switch directory {
	case "golang-forum", "app":
		sqlTemplatePath = "sql/template.xml"
	default:
		// 執行測試時的路徑。
		sqlTemplatePath = "../../sql/template.xml"
		fmt.Println("============== directory", directory)
	}

	doc := etree.NewDocument()

	if err := doc.ReadFromFile(sqlTemplatePath); err != nil {
		panic(err)
	}

	sqls := doc.SelectElement("Sqls")
	for _, sql := range sqls.SelectElements("Sql") {
		name := sql.SelectAttrValue("name", "")
		SqlTemplate[name] = sql.Text()
	}
}
