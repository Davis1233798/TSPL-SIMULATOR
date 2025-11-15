package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"tspl-simulator/models"
)

const DPI = 203 // 標準熱感打印機 DPI

// ParseTSPL 解析 TSPL 指令
func ParseTSPL(tsplCode string) (*models.RenderData, error) {
	renderData := &models.RenderData{
		Elements:  []models.Element{},
		DPI:       DPI,
		Direction: 0,
		Reference: models.Reference{X: 0, Y: 0},
	}

	lines := strings.Split(tsplCode, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, ";") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		command := strings.ToUpper(parts[0])

		switch command {
		case "SIZE":
			if err := parseSize(parts[1:], renderData); err != nil {
				return nil, err
			}
		case "GAP":
			if err := parseGap(parts[1:], renderData); err != nil {
				return nil, err
			}
		case "DIRECTION":
			if err := parseDirection(parts[1:], renderData); err != nil {
				return nil, err
			}
		case "REFERENCE":
			if err := parseReference(parts[1:], renderData); err != nil {
				return nil, err
			}
		case "CLS":
			// 清除緩衝區,不需要特別處理
		case "TEXT":
			if err := parseText(line, renderData); err != nil {
				return nil, err
			}
		case "BARCODE":
			if err := parseBarcode(line, renderData); err != nil {
				return nil, err
			}
		case "QRCODE":
			if err := parseQRCode(line, renderData); err != nil {
				return nil, err
			}
		case "BOX":
			if err := parseBox(parts[1:], renderData); err != nil {
				return nil, err
			}
		case "BAR":
			if err := parseBar(parts[1:], renderData); err != nil {
				return nil, err
			}
		case "PRINT":
			// 列印指令,不需要特別處理
		}
	}

	// 計算畫布尺寸
	renderData.Width = mmToPixels(renderData.LabelSize.Width)
	renderData.Height = mmToPixels(renderData.LabelSize.Height)

	return renderData, nil
}

// parseSize 解析 SIZE 指令
func parseSize(parts []string, renderData *models.RenderData) error {
	if len(parts) < 2 {
		return fmt.Errorf("SIZE 指令參數不足")
	}

	width, unit1, err := parseValueWithUnit(parts[0])
	if err != nil {
		return fmt.Errorf("SIZE 寬度格式錯誤: %v", err)
	}

	height, unit2, err := parseValueWithUnit(parts[1])
	if err != nil {
		return fmt.Errorf("SIZE 高度格式錯誤: %v", err)
	}

	if unit1 != unit2 {
		return fmt.Errorf("SIZE 單位不一致")
	}

	renderData.LabelSize = models.LabelSize{
		Width:  width,
		Height: height,
		Unit:   unit1,
	}

	return nil
}

// parseGap 解析 GAP 指令
func parseGap(parts []string, renderData *models.RenderData) error {
	if len(parts) < 2 {
		return fmt.Errorf("GAP 指令參數不足")
	}

	distance, unit1, err := parseValueWithUnit(parts[0])
	if err != nil {
		return fmt.Errorf("GAP 距離格式錯誤: %v", err)
	}

	offset, unit2, err := parseValueWithUnit(parts[1])
	if err != nil {
		return fmt.Errorf("GAP 偏移格式錯誤: %v", err)
	}

	if unit1 != unit2 {
		return fmt.Errorf("GAP 單位不一致")
	}

	renderData.Gap = models.Gap{
		Distance: distance,
		Offset:   offset,
		Unit:     unit1,
	}

	return nil
}

// parseDirection 解析 DIRECTION 指令
func parseDirection(parts []string, renderData *models.RenderData) error {
	if len(parts) < 1 {
		return fmt.Errorf("DIRECTION 指令參數不足")
	}

	direction, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("DIRECTION 參數格式錯誤: %v", err)
	}

	renderData.Direction = direction
	return nil
}

// parseReference 解析 REFERENCE 指令
func parseReference(parts []string, renderData *models.RenderData) error {
	if len(parts) < 2 {
		return fmt.Errorf("REFERENCE 指令參數不足")
	}

	x, err := strconv.Atoi(strings.TrimSuffix(parts[0], ","))
	if err != nil {
		return fmt.Errorf("REFERENCE X 參數格式錯誤: %v", err)
	}

	y, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("REFERENCE Y 參數格式錯誤: %v", err)
	}

	renderData.Reference = models.Reference{X: x, Y: y}
	return nil
}

// parseText 解析 TEXT 指令
func parseText(line string, renderData *models.RenderData) error {
	re := regexp.MustCompile(`TEXT\s+(\d+),(\d+),"([^"]+)",(\d+),(\d+),(\d+)(?:,(\d+),(\d+))?,(?:"([^"]*)")`)
	matches := re.FindStringSubmatch(line)

	if len(matches) < 10 {
		return fmt.Errorf("TEXT 指令格式錯誤")
	}

	x, _ := strconv.Atoi(matches[1])
	y, _ := strconv.Atoi(matches[2])
	font := matches[3]
	rotation, _ := strconv.Atoi(matches[4])
	xScale, _ := strconv.Atoi(matches[5])
	yScale, _ := strconv.Atoi(matches[6])
	text := matches[9]

	element := models.Element{
		Type: "text",
		X:    x,
		Y:    y,
		Properties: map[string]interface{}{
			"text":     text,
			"font":     font,
			"rotation": rotation,
			"xScale":   xScale,
			"yScale":   yScale,
		},
	}

	renderData.Elements = append(renderData.Elements, element)
	return nil
}

// parseBarcode 解析 BARCODE 指令
func parseBarcode(line string, renderData *models.RenderData) error {
	re := regexp.MustCompile(`BARCODE\s+(\d+),(\d+),"([^"]+)",(\d+),(\d+),(\d+),(\d+),(\d+),"([^"]*)"`)
	matches := re.FindStringSubmatch(line)

	if len(matches) < 10 {
		return fmt.Errorf("BARCODE 指令格式錯誤")
	}

	x, _ := strconv.Atoi(matches[1])
	y, _ := strconv.Atoi(matches[2])
	codeType := matches[3]
	height, _ := strconv.Atoi(matches[4])
	readable, _ := strconv.Atoi(matches[5])
	rotation, _ := strconv.Atoi(matches[6])
	narrow, _ := strconv.Atoi(matches[7])
	wide, _ := strconv.Atoi(matches[8])
	code := matches[9]

	element := models.Element{
		Type: "barcode",
		X:    x,
		Y:    y,
		Properties: map[string]interface{}{
			"code":     code,
			"type":     codeType,
			"height":   height,
			"readable": readable,
			"rotation": rotation,
			"narrow":   narrow,
			"wide":     wide,
		},
	}

	renderData.Elements = append(renderData.Elements, element)
	return nil
}

// parseQRCode 解析 QRCODE 指令
func parseQRCode(line string, renderData *models.RenderData) error {
	re := regexp.MustCompile(`QRCODE\s+(\d+),(\d+),([HML]),(\d+),([AM]),(\d+)(?:,(\d+),(\d+),(\d+))?,(?:"([^"]*)")`)
	matches := re.FindStringSubmatch(line)

	if len(matches) < 11 {
		return fmt.Errorf("QRCODE 指令格式錯誤")
	}

	x, _ := strconv.Atoi(matches[1])
	y, _ := strconv.Atoi(matches[2])
	eccLevel := matches[3]
	cellSize, _ := strconv.Atoi(matches[4])
	mode := matches[5]
	rotation, _ := strconv.Atoi(matches[6])
	data := matches[10]

	element := models.Element{
		Type: "qrcode",
		X:    x,
		Y:    y,
		Properties: map[string]interface{}{
			"data":     data,
			"eccLevel": eccLevel,
			"cellSize": cellSize,
			"mode":     mode,
			"rotation": rotation,
		},
	}

	renderData.Elements = append(renderData.Elements, element)
	return nil
}

// parseBox 解析 BOX 指令
func parseBox(parts []string, renderData *models.RenderData) error {
	if len(parts) < 5 {
		return fmt.Errorf("BOX 指令參數不足")
	}

	x, _ := strconv.Atoi(strings.TrimSuffix(parts[0], ","))
	y, _ := strconv.Atoi(strings.TrimSuffix(parts[1], ","))
	endX, _ := strconv.Atoi(strings.TrimSuffix(parts[2], ","))
	endY, _ := strconv.Atoi(strings.TrimSuffix(parts[3], ","))
	thickness, _ := strconv.Atoi(parts[4])

	element := models.Element{
		Type: "box",
		X:    x,
		Y:    y,
		Properties: map[string]interface{}{
			"endX":      endX,
			"endY":      endY,
			"thickness": thickness,
		},
	}

	renderData.Elements = append(renderData.Elements, element)
	return nil
}

// parseBar 解析 BAR 指令
func parseBar(parts []string, renderData *models.RenderData) error {
	if len(parts) < 4 {
		return fmt.Errorf("BAR 指令參數不足")
	}

	x, _ := strconv.Atoi(strings.TrimSuffix(parts[0], ","))
	y, _ := strconv.Atoi(strings.TrimSuffix(parts[1], ","))
	width, _ := strconv.Atoi(strings.TrimSuffix(parts[2], ","))
	height, _ := strconv.Atoi(parts[3])

	element := models.Element{
		Type: "bar",
		X:    x,
		Y:    y,
		Properties: map[string]interface{}{
			"width":  width,
			"height": height,
		},
	}

	renderData.Elements = append(renderData.Elements, element)
	return nil
}

// parseValueWithUnit 解析帶單位的數值
func parseValueWithUnit(s string) (float64, string, error) {
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(s, ",")

	re := regexp.MustCompile(`^([\d.]+)\s*(mm|inch)?$`)
	matches := re.FindStringSubmatch(s)

	if len(matches) < 2 {
		return 0, "", fmt.Errorf("無法解析數值: %s", s)
	}

	value, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0, "", err
	}

	unit := "mm"
	if len(matches) > 2 && matches[2] != "" {
		unit = matches[2]
	}

	return value, unit, nil
}

// mmToPixels 毫米轉像素
func mmToPixels(mm float64) int {
	inches := mm / 25.4
	return int(inches * float64(DPI))
}
