package validator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ValidationError 驗證錯誤
type ValidationError struct {
	Line    int    `json:"line"`
	Command string `json:"command"`
	Message string `json:"message"`
}

// ValidationResult 驗證結果
type ValidationResult struct {
	Valid  bool               `json:"valid"`
	Errors []ValidationError  `json:"errors,omitempty"`
}

// ValidateTSPL 驗證 TSPL2 語法
func ValidateTSPL(tsplCode string) *ValidationResult {
	result := &ValidationResult{
		Valid:  true,
		Errors: []ValidationError{},
	}

	lines := strings.Split(tsplCode, "\n")
	hasSize := false
	hasPrint := false

	for i, line := range lines {
		line = strings.TrimSpace(line)
		lineNum := i + 1

		// 跳過空行和註解
		if line == "" || strings.HasPrefix(line, ";") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		command := strings.ToUpper(parts[0])

		// 檢查命令是否有效
		if !isValidCommand(command) {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Line:    lineNum,
				Command: command,
				Message: fmt.Sprintf("未知的命令: %s", command),
			})
			continue
		}

		// 檢查特定命令的語法
		switch command {
		case "SIZE":
			hasSize = true
			if err := validateSize(parts[1:], lineNum); err != nil {
				result.Valid = false
				result.Errors = append(result.Errors, *err)
			}

		case "GAP":
			if err := validateGap(parts[1:], lineNum); err != nil {
				result.Valid = false
				result.Errors = append(result.Errors, *err)
			}

		case "DIRECTION":
			if err := validateDirection(parts[1:], lineNum); err != nil {
				result.Valid = false
				result.Errors = append(result.Errors, *err)
			}

		case "REFERENCE":
			if err := validateReference(parts[1:], lineNum); err != nil {
				result.Valid = false
				result.Errors = append(result.Errors, *err)
			}

		case "TEXT":
			if err := validateText(line, lineNum); err != nil {
				result.Valid = false
				result.Errors = append(result.Errors, *err)
			}

		case "BARCODE":
			if err := validateBarcode(line, lineNum); err != nil {
				result.Valid = false
				result.Errors = append(result.Errors, *err)
			}

		case "QRCODE":
			if err := validateQRCode(line, lineNum); err != nil {
				result.Valid = false
				result.Errors = append(result.Errors, *err)
			}

		case "BOX":
			if err := validateBox(parts[1:], lineNum); err != nil {
				result.Valid = false
				result.Errors = append(result.Errors, *err)
			}

		case "BAR":
			if err := validateBar(parts[1:], lineNum); err != nil {
				result.Valid = false
				result.Errors = append(result.Errors, *err)
			}

		case "PRINT":
			hasPrint = true
			if err := validatePrint(parts[1:], lineNum); err != nil {
				result.Valid = false
				result.Errors = append(result.Errors, *err)
			}

		case "DENSITY":
			if err := validateDensity(parts[1:], lineNum); err != nil {
				result.Valid = false
				result.Errors = append(result.Errors, *err)
			}

		case "SPEED":
			if err := validateSpeed(parts[1:], lineNum); err != nil {
				result.Valid = false
				result.Errors = append(result.Errors, *err)
			}
		}
	}

	// 檢查必要的命令
	if !hasSize {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Line:    0,
			Command: "SIZE",
			Message: "缺少必要的 SIZE 命令",
		})
	}

	if !hasPrint {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Line:    0,
			Command: "PRINT",
			Message: "缺少必要的 PRINT 命令",
		})
	}

	return result
}

// isValidCommand 檢查命令是否有效
func isValidCommand(command string) bool {
	validCommands := map[string]bool{
		"SIZE": true, "GAP": true, "DIRECTION": true, "REFERENCE": true,
		"CLS": true, "TEXT": true, "BARCODE": true, "QRCODE": true,
		"BOX": true, "BAR": true, "PRINT": true, "DENSITY": true,
		"SPEED": true, "SET": true, "SHIFT": true, "OFFSET": true,
		"BITMAP": true, "REVERSE": true, "FORMFEED": true, "BACKFEED": true,
		"HOME": true, "SOUND": true, "LIMITFEED": true, "SELFTEST": true,
		"EOP": true, "BLOCK": true, "CODEPAGE": true, "COUNTRY": true,
		"PUTBMP": true, "PUTPCX": true, "DOWNLOAD": true, "ERASE": true,
	}
	return validCommands[command]
}

// validateSize 驗證 SIZE 命令
func validateSize(parts []string, lineNum int) *ValidationError {
	if len(parts) < 2 {
		return &ValidationError{
			Line:    lineNum,
			Command: "SIZE",
			Message: "SIZE 命令需要 2 個參數 (寬度, 高度)",
		}
	}

	// 檢查寬度
	if _, _, err := parseValueWithUnit(parts[0]); err != nil {
		return &ValidationError{
			Line:    lineNum,
			Command: "SIZE",
			Message: fmt.Sprintf("寬度格式錯誤: %v", err),
		}
	}

	// 檢查高度
	if _, _, err := parseValueWithUnit(parts[1]); err != nil {
		return &ValidationError{
			Line:    lineNum,
			Command: "SIZE",
			Message: fmt.Sprintf("高度格式錯誤: %v", err),
		}
	}

	return nil
}

// validateGap 驗證 GAP 命令
func validateGap(parts []string, lineNum int) *ValidationError {
	if len(parts) < 2 {
		return &ValidationError{
			Line:    lineNum,
			Command: "GAP",
			Message: "GAP 命令需要 2 個參數 (間距, 偏移)",
		}
	}

	// 檢查間距
	if _, _, err := parseValueWithUnit(parts[0]); err != nil {
		return &ValidationError{
			Line:    lineNum,
			Command: "GAP",
			Message: fmt.Sprintf("間距格式錯誤: %v", err),
		}
	}

	// 檢查偏移
	if _, _, err := parseValueWithUnit(parts[1]); err != nil {
		return &ValidationError{
			Line:    lineNum,
			Command: "GAP",
			Message: fmt.Sprintf("偏移格式錯誤: %v", err),
		}
	}

	return nil
}

// validateDirection 驗證 DIRECTION 命令
func validateDirection(parts []string, lineNum int) *ValidationError {
	if len(parts) < 1 {
		return &ValidationError{
			Line:    lineNum,
			Command: "DIRECTION",
			Message: "DIRECTION 命令需要 1 個參數 (方向)",
		}
	}

	direction, err := strconv.Atoi(strings.TrimSuffix(parts[0], ","))
	if err != nil {
		return &ValidationError{
			Line:    lineNum,
			Command: "DIRECTION",
			Message: "方向必須是數字",
		}
	}

	if direction < 0 || direction > 3 {
		return &ValidationError{
			Line:    lineNum,
			Command: "DIRECTION",
			Message: "方向必須在 0-3 之間",
		}
	}

	return nil
}

// validateReference 驗證 REFERENCE 命令
func validateReference(parts []string, lineNum int) *ValidationError {
	if len(parts) < 2 {
		return &ValidationError{
			Line:    lineNum,
			Command: "REFERENCE",
			Message: "REFERENCE 命令需要 2 個參數 (x, y)",
		}
	}

	// 檢查 x 座標
	if _, err := strconv.Atoi(strings.TrimSuffix(parts[0], ",")); err != nil {
		return &ValidationError{
			Line:    lineNum,
			Command: "REFERENCE",
			Message: "X 座標必須是數字",
		}
	}

	// 檢查 y 座標
	if _, err := strconv.Atoi(parts[1]); err != nil {
		return &ValidationError{
			Line:    lineNum,
			Command: "REFERENCE",
			Message: "Y 座標必須是數字",
		}
	}

	return nil
}

// validateText 驗證 TEXT 命令
func validateText(line string, lineNum int) *ValidationError {
	re := regexp.MustCompile(`TEXT\s+(\d+),(\d+),"([^"]+)",(\d+),(\d+),(\d+)(?:,(\d+),(\d+))?,(?:"([^"]*)")`)
	if !re.MatchString(line) {
		return &ValidationError{
			Line:    lineNum,
			Command: "TEXT",
			Message: "TEXT 命令格式錯誤。正確格式: TEXT x,y,\"font\",rotation,x-scale,y-scale,\"content\"",
		}
	}
	return nil
}

// validateBarcode 驗證 BARCODE 命令
func validateBarcode(line string, lineNum int) *ValidationError {
	re := regexp.MustCompile(`BARCODE\s+(\d+),(\d+),"([^"]+)",(\d+),(\d+),(\d+),(\d+),(\d+),"([^"]*)"`)
	if !re.MatchString(line) {
		return &ValidationError{
			Line:    lineNum,
			Command: "BARCODE",
			Message: "BARCODE 命令格式錯誤。正確格式: BARCODE x,y,\"type\",height,readable,rotation,narrow,wide,\"code\"",
		}
	}
	return nil
}

// validateQRCode 驗證 QRCODE 命令
func validateQRCode(line string, lineNum int) *ValidationError {
	re := regexp.MustCompile(`QRCODE\s+(\d+),(\d+),([HML]),(\d+),([AM]),(\d+)(?:,(\d+),(\d+),(\d+))?,(?:"([^"]*)")`)
	if !re.MatchString(line) {
		return &ValidationError{
			Line:    lineNum,
			Command: "QRCODE",
			Message: "QRCODE 命令格式錯誤。正確格式: QRCODE x,y,ecc,size,mode,rotation,\"data\"",
		}
	}
	return nil
}

// validateBox 驗證 BOX 命令
func validateBox(parts []string, lineNum int) *ValidationError {
	if len(parts) < 5 {
		return &ValidationError{
			Line:    lineNum,
			Command: "BOX",
			Message: "BOX 命令需要 5 個參數 (x, y, x_end, y_end, thickness)",
		}
	}

	// 驗證所有參數都是數字
	for i, param := range parts[:5] {
		if _, err := strconv.Atoi(strings.TrimSuffix(param, ",")); err != nil {
			return &ValidationError{
				Line:    lineNum,
				Command: "BOX",
				Message: fmt.Sprintf("第 %d 個參數必須是數字", i+1),
			}
		}
	}

	return nil
}

// validateBar 驗證 BAR 命令
func validateBar(parts []string, lineNum int) *ValidationError {
	if len(parts) < 4 {
		return &ValidationError{
			Line:    lineNum,
			Command: "BAR",
			Message: "BAR 命令需要 4 個參數 (x, y, width, height)",
		}
	}

	// 驗證所有參數都是數字
	for i, param := range parts[:4] {
		if _, err := strconv.Atoi(strings.TrimSuffix(param, ",")); err != nil {
			return &ValidationError{
				Line:    lineNum,
				Command: "BAR",
				Message: fmt.Sprintf("第 %d 個參數必須是數字", i+1),
			}
		}
	}

	return nil
}

// validatePrint 驗證 PRINT 命令
func validatePrint(parts []string, lineNum int) *ValidationError {
	if len(parts) > 2 {
		return &ValidationError{
			Line:    lineNum,
			Command: "PRINT",
			Message: "PRINT 命令最多接受 2 個參數 (數量, 副本)",
		}
	}

	// 如果有參數,檢查是否為數字
	for i, param := range parts {
		if _, err := strconv.Atoi(strings.TrimSuffix(param, ",")); err != nil {
			return &ValidationError{
				Line:    lineNum,
				Command: "PRINT",
				Message: fmt.Sprintf("第 %d 個參數必須是數字", i+1),
			}
		}
	}

	return nil
}

// validateDensity 驗證 DENSITY 命令
func validateDensity(parts []string, lineNum int) *ValidationError {
	if len(parts) < 1 {
		return &ValidationError{
			Line:    lineNum,
			Command: "DENSITY",
			Message: "DENSITY 命令需要 1 個參數 (濃度)",
		}
	}

	density, err := strconv.Atoi(strings.TrimSuffix(parts[0], ","))
	if err != nil {
		return &ValidationError{
			Line:    lineNum,
			Command: "DENSITY",
			Message: "濃度必須是數字",
		}
	}

	if density < 0 || density > 15 {
		return &ValidationError{
			Line:    lineNum,
			Command: "DENSITY",
			Message: "濃度必須在 0-15 之間",
		}
	}

	return nil
}

// validateSpeed 驗證 SPEED 命令
func validateSpeed(parts []string, lineNum int) *ValidationError {
	if len(parts) < 1 {
		return &ValidationError{
			Line:    lineNum,
			Command: "SPEED",
			Message: "SPEED 命令需要 1 個參數 (速度)",
		}
	}

	speed, err := strconv.ParseFloat(strings.TrimSuffix(parts[0], ","), 64)
	if err != nil {
		return &ValidationError{
			Line:    lineNum,
			Command: "SPEED",
			Message: "速度必須是數字",
		}
	}

	// 一般速度範圍 1-14 英吋/秒
	if speed < 1 || speed > 14 {
		return &ValidationError{
			Line:    lineNum,
			Command: "SPEED",
			Message: "速度建議在 1-14 之間",
		}
	}

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
