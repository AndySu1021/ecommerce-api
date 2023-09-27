package vo

import (
	"database/sql/driver"
	"encoding/json"
)

type Temperature int8

const (
	TemperatureNormal Temperature = iota + 1 // 常溫
	TemperatureCool                          // 冷藏
	TemperatureFreeze                        // 冷凍
)

type SortBy string

const (
	SortByDefault SortBy = "default" // 預設
	SortByPrice   SortBy = "price"   // 價格
	SortByTime    SortBy = "time"    // 創建時間
	SortBySales   SortBy = "sales"   // 銷售量
	SortByStock   SortBy = "stock"   // 有存貨
)

type PictureArray []string

func (f *PictureArray) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	return json.Unmarshal(value.([]byte), f)
}

func (f PictureArray) Value() (driver.Value, error) {
	return json.Marshal(&f)
}

type SupportDeliveryMethod []int8

func (f *SupportDeliveryMethod) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	return json.Unmarshal(value.([]byte), f)
}

func (f SupportDeliveryMethod) Value() (driver.Value, error) {
	return json.Marshal(&f)
}

type Extra struct {
}
