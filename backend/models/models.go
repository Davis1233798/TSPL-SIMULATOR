package models

// RenderRequest 渲染請求
type RenderRequest struct {
	TSPLCode string `json:"tspl_code" binding:"required"`
}

// RenderResponse 渲染回應
type RenderResponse struct {
	Success          bool              `json:"success"`
	Data             *RenderData       `json:"data,omitempty"`
	Error            string            `json:"error,omitempty"`
	ValidationErrors []ValidationError `json:"validation_errors,omitempty"`
}

// ValidationError 驗證錯誤
type ValidationError struct {
	Line    int    `json:"line"`
	Command string `json:"command"`
	Message string `json:"message"`
}

// RenderData 渲染資料
type RenderData struct {
	Width      int            `json:"width"`
	Height     int            `json:"height"`
	Elements   []Element      `json:"elements"`
	LabelSize  LabelSize      `json:"labelSize"`
	Gap        Gap            `json:"gap"`
	Direction  int            `json:"direction"`
	Reference  Reference      `json:"reference"`
	DPI        int            `json:"dpi"`
}

// Element 渲染元素
type Element struct {
	Type       string                 `json:"type"`
	X          int                    `json:"x"`
	Y          int                    `json:"y"`
	Properties map[string]interface{} `json:"properties"`
}

// LabelSize 標籤尺寸
type LabelSize struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Unit   string  `json:"unit"`
}

// Gap 標籤間距
type Gap struct {
	Distance float64 `json:"distance"`
	Offset   float64 `json:"offset"`
	Unit     string  `json:"unit"`
}

// Reference 參考點
type Reference struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// ExampleInfo 範例資訊
type ExampleInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

// ExamplesResponse 範例列表回應
type ExamplesResponse struct {
	Success  bool          `json:"success"`
	Examples []ExampleInfo `json:"examples,omitempty"`
	Error    string        `json:"error,omitempty"`
}

// ExampleDetailResponse 範例詳情回應
type ExampleDetailResponse struct {
	Success bool   `json:"success"`
	Code    string `json:"code,omitempty"`
	Error   string `json:"error,omitempty"`
}

// MQTTMessage MQTT 訊息
type MQTTMessage struct {
	Type      string `json:"type"`
	TSPLCode  string `json:"tspl_code,omitempty"`
	Timestamp int64  `json:"timestamp"`
}

// HealthResponse 健康檢查回應
type HealthResponse struct {
	Status string `json:"status"`
	MQTT   string `json:"mqtt"`
}
