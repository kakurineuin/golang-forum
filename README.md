使用 Golang、Echo、GROM、React 建立的論壇網站。

Golang
https://golang.org

Echo
https://echo.labstack.com

GROM
http://gorm.io

React
https://reactjs.org

目錄與檔案說明
===
- config：專案配置檔的目錄。
- database：初始化資料存取物件(data access object，DAO)的程式碼目錄。
- error：自訂錯誤的程式碼目錄。
- frontend：前端程式碼的目錄，前端使用 React 開發。
- handler：Echo 的 handler 目錄。
- logger：初始化日誌物件的程式碼目錄。
- middleware：Echo 的 middleware 目錄。
- model：GORM 的 model 目錄，每個 model 對應一張資料表。用來裝請求參數的 struct 也歸類在此目錄。
- service：實作論壇各功能的 service 的目錄，例如有關文章的新增、修改、刪除、查詢等功能會寫在 topic.go。在 Echo 的 handler 中會呼叫對應的 service 去處理請求。寫成 service 是為了將論壇的業務邏輯(文章 CRUD、使用者管理)和 Web 的瑣碎處理(檢核 request、回傳 response)分開，以提供更好的可維護性。
- sql：sql 語句樣板的放置目錄。遇到比較複雜的查詢必須寫成 sql 時，就寫在此目錄下的 template.xml，然後在程式碼中讀取該段 sql 去執行。
- test：專案程式碼的測試目錄，在專案目錄下執行 ginkgo -r 即可執行測試。
- validate：初始化 govalidator 的自訂驗證器的程式碼目錄。govalidator 是檢核 Golang struct 的函式庫，此專案使用此函式庫來檢核請求參數是否正確。
- vendor：依賴函式庫的放置目錄。
- glide.yaml：依賴函式庫的配置檔。
- main.go：專案程式碼的進入點。
