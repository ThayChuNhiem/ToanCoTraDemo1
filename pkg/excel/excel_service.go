package excel

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
	"toan-co-tra-backend/internal/domain"

	"github.com/xuri/excelize/v2"
)

var phoneRegex = regexp.MustCompile(`^(03|05|07|08|09)\d{8}$`)

// -----------------------------------------------------------------------------
// 1. DỊCH VỤ EXCEL CHO DANH MỤC LỚP HỌC (CLASSES)
// -----------------------------------------------------------------------------

// ExportClassesToExcel xuất toàn bộ danh sách lớp học ra file Excel (.xlsx)
func ExportClassesToExcel(classes []*domain.Class) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheet := "DanhSachLopHoc"
	f.SetSheetName("Sheet1", sheet)

	// Tạo kiểu dáng tiêu đề chuyên nghiệp
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Color:  "FFFFFF",
			Size:   11,
			Family: "Arial",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"10B981"}, // Ngọc lục bảo Emerald
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "D1D5DB", Style: 1},
			{Type: "top", Color: "D1D5DB", Style: 1},
			{Type: "bottom", Color: "D1D5DB", Style: 1},
			{Type: "right", Color: "D1D5DB", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}

	dataStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size:   10,
			Family: "Arial",
		},
		Alignment: &excelize.Alignment{
			Vertical: "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "E5E7EB", Style: 1},
			{Type: "top", Color: "E5E7EB", Style: 1},
			{Type: "bottom", Color: "E5E7EB", Style: 1},
			{Type: "right", Color: "E5E7EB", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}

	headers := []string{
		"Mã Lớp (ID)",
		"Tên Lớp Học",
		"Khối Lớp (1-9)",
		"Hình Thức (online/offline)",
		"Loại Lớp (basic/advanced/high_quality)",
		"Học Phí",
		"Mũi Nhọn (Có/Không)",
		"Mô Tả Khóa Học",
	}

	// Đặt tiêu đề
	for colIdx, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(colIdx+1, 1)
		_ = f.SetCellValue(sheet, cell, h)
		_ = f.SetCellStyle(sheet, cell, cell, headerStyle)
	}
	_ = f.SetRowHeight(sheet, 1, 28)

	// Ghi dữ liệu từng dòng
	for rowIdx, c := range classes {
		row := rowIdx + 2
		popularText := "Không"
		if c.IsPopular {
			popularText = "Có"
		}

		values := []interface{}{
			c.ID,
			c.Name,
			c.Grade,
			c.Model,
			c.Type,
			c.Price,
			popularText,
			c.Desc,
		}

		for colIdx, val := range values {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, row)
			_ = f.SetCellValue(sheet, cell, val)
			_ = f.SetCellStyle(sheet, cell, cell, dataStyle)
		}
		_ = f.SetRowHeight(sheet, row, 22)
	}

	// Tùy chỉnh độ rộng các cột
	colWidths := map[string]float64{
		"A": 18, // Mã Lớp
		"B": 36, // Tên Lớp
		"C": 14, // Khối Lớp
		"D": 22, // Hình Thức
		"E": 28, // Loại Lớp
		"F": 22, // Học Phí
		"G": 18, // Mũi Nhọn
		"H": 45, // Mô Tả
	}
	for col, width := range colWidths {
		_ = f.SetColWidth(sheet, col, col, width)
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GenerateClassTemplateExcel tạo file Excel mẫu (.xlsx) để quản trị viên nhập mới/chỉnh sửa hàng loạt lớp học
func GenerateClassTemplateExcel() ([]byte, error) {
	sampleClasses := []*domain.Class{
		{
			ID:        "class_5_clc",
			Name:      "Toán Tư Duy Lớp 5 CLC Chuyên Sâu",
			Grade:     5,
			Model:     "offline",
			Type:      "high_quality",
			Price:     "1.500.000đ/tháng",
			IsPopular: true,
			Desc:      "Chương trình đào tạo mũi nhọn ôn thi vào các trường chuyên cấp 2 chất lượng cao",
		},
		{
			ID:        "class_4_adv",
			Name:      "Toán Nâng Cao Lớp 4 Bứt Phá",
			Grade:     4,
			Model:     "online",
			Type:      "advanced",
			Price:     "1.200.000đ/tháng",
			IsPopular: false,
			Desc:      "Bồi dưỡng tư duy hình học và số học nâng cao cho học sinh khá giỏi",
		},
		{
			ID:        "class_3_bas",
			Name:      "Toán Nền Tảng Lớp 3 Vững Vàng",
			Grade:     3,
			Model:     "offline",
			Type:      "basic",
			Price:     "1.000.000đ/tháng",
			IsPopular: false,
			Desc:      "Củng cố kiến thức nền tảng sách giáo khoa, rèn tính cẩn thận và tính nhẩm nhanh",
		},
	}
	return ExportClassesToExcel(sampleClasses)
}

// ParseClassesFromExcel đọc file Excel tải lên, kiểm tra tính hợp lệ và trả về danh sách lớp học cùng danh sách lỗi
func ParseClassesFromExcel(r io.Reader) ([]*domain.Class, []string, error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, nil, fmt.Errorf("không thể đọc tệp tin Excel: %v", err)
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, nil, fmt.Errorf("tệp tin Excel không có trang tính nào")
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, nil, fmt.Errorf("lỗi đọc dữ liệu trang tính: %v", err)
	}

	if len(rows) <= 1 {
		return nil, nil, fmt.Errorf("tệp tin Excel không có dữ liệu để nhập (chỉ có dòng tiêu đề hoặc trống)")
	}

	var classes []*domain.Class
	var validationErrors []string

	for i := 1; i < len(rows); i++ {
		rowNum := i + 1
		row := rows[i]

		// Bỏ qua dòng trống hoàn toàn
		if len(row) == 0 || isEmptyRow(row) {
			continue
		}

		getCol := func(idx int) string {
			if idx < len(row) {
				return strings.TrimSpace(row[idx])
			}
			return ""
		}

		id := getCol(0)
		name := getCol(1)
		gradeStr := getCol(2)
		model := strings.ToLower(getCol(3))
		classType := strings.ToLower(getCol(4))
		price := getCol(5)
		popularStr := strings.ToLower(getCol(6))
		desc := getCol(7)

		// 1. Kiểm tra Tên Lớp
		if name == "" {
			validationErrors = append(validationErrors, fmt.Sprintf("Dòng %d: Tên lớp học không được để trống", rowNum))
			continue
		}

		// 2. Tự động sinh ID nếu để trống
		if id == "" {
			id = fmt.Sprintf("class_%d_%d", time.Now().UnixNano(), rowNum)
		}

		// 3. Kiểm tra Khối Lớp
		grade, err := strconv.Atoi(gradeStr)
		if err != nil || grade < 1 || grade > 9 {
			validationErrors = append(validationErrors, fmt.Sprintf("Dòng %d (Lớp '%s'): Khối lớp phải là số nguyên từ 1 đến 9 (nhận được: '%s')", rowNum, name, gradeStr))
			continue
		}

		// 4. Kiểm tra Hình Thức
		if model != "online" && model != "offline" {
			// Hỗ trợ tự động chuẩn hóa tiếng Việt nếu người dùng gõ
			if strings.Contains(model, "on") || strings.Contains(model, "tuyến") {
				model = "online"
			} else if strings.Contains(model, "off") || strings.Contains(model, "trực tiếp") || strings.Contains(model, "tại lớp") {
				model = "offline"
			} else {
				validationErrors = append(validationErrors, fmt.Sprintf("Dòng %d (Lớp '%s'): Hình thức học phải là 'online' hoặc 'offline' (nhận được: '%s')", rowNum, name, getCol(3)))
				continue
			}
		}

		// 5. Kiểm tra Loại Lớp
		if classType != "basic" && classType != "advanced" && classType != "high_quality" {
			if strings.Contains(classType, "cơ bản") || strings.Contains(classType, "co ban") {
				classType = "basic"
			} else if strings.Contains(classType, "nâng cao") || strings.Contains(classType, "nang cao") {
				classType = "advanced"
			} else if strings.Contains(classType, "clc") || strings.Contains(classType, "chất lượng cao") {
				classType = "high_quality"
			} else {
				validationErrors = append(validationErrors, fmt.Sprintf("Dòng %d (Lớp '%s'): Loại lớp phải là 'basic', 'advanced' hoặc 'high_quality' (nhận được: '%s')", rowNum, name, getCol(4)))
				continue
			}
		}

		// Quy tắc nghiệp vụ: CLC chỉ dành cho Khối 5
		if classType == "high_quality" && grade != 5 {
			validationErrors = append(validationErrors, fmt.Sprintf("Dòng %d (Lớp '%s'): Lớp Chất lượng cao (high_quality) chỉ dành riêng cho Khối Lớp 5", rowNum, name))
			continue
		}

		// 6. Kiểm tra Học Phí
		if price == "" {
			validationErrors = append(validationErrors, fmt.Sprintf("Dòng %d (Lớp '%s'): Học phí không được để trống", rowNum, name))
			continue
		}

		// 7. Nhận diện Mũi Nhọn
		isPopular := false
		if popularStr == "có" || popularStr == "co" || popularStr == "yes" || popularStr == "true" || popularStr == "1" || popularStr == "x" {
			isPopular = true
		}

		classes = append(classes, &domain.Class{
			ID:        id,
			Name:      name,
			Grade:     grade,
			Model:     model,
			Type:      classType,
			Price:     price,
			IsPopular: isPopular,
			Desc:      desc,
		})
	}

	return classes, validationErrors, nil
}

// -----------------------------------------------------------------------------
// 2. DỊCH VỤ EXCEL CHO DANH SÁCH HỌC VIÊN / LEADS (Đúng yêu cầu SĐT & Họ tên)
// -----------------------------------------------------------------------------

// ExportLeadsToExcel xuất toàn bộ danh sách học viên đăng ký ra file Excel (.xlsx)
func ExportLeadsToExcel(leads []*domain.Lead) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheet := "DanhSachHocVien"
	f.SetSheetName("Sheet1", sheet)

	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Color:  "FFFFFF",
			Size:   11,
			Family: "Arial",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"2563EB"}, // Xanh lam hoàng gia
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "D1D5DB", Style: 1},
			{Type: "top", Color: "D1D5DB", Style: 1},
			{Type: "bottom", Color: "D1D5DB", Style: 1},
			{Type: "right", Color: "D1D5DB", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}

	dataStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size:   10,
			Family: "Arial",
		},
		Alignment: &excelize.Alignment{
			Vertical: "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "E5E7EB", Style: 1},
			{Type: "top", Color: "E5E7EB", Style: 1},
			{Type: "bottom", Color: "E5E7EB", Style: 1},
			{Type: "right", Color: "E5E7EB", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}

	headers := []string{
		"Mã ID",
		"Họ Tên Học Sinh",
		"Họ Tên Phụ Huynh",
		"Số Điện Thoại",
		"Khối Lớp (1-9)",
		"Hình Thức (online/offline)",
		"Loại Lớp (basic/advanced/high_quality)",
		"Học Lực (excellent/good/average)",
		"Trạng Thái (pending/contacted)",
		"Người Tư Vấn",
		"Ngày Đăng Ký",
	}

	for colIdx, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(colIdx+1, 1)
		_ = f.SetCellValue(sheet, cell, h)
		_ = f.SetCellStyle(sheet, cell, cell, headerStyle)
	}
	_ = f.SetRowHeight(sheet, 1, 28)

	for rowIdx, l := range leads {
		row := rowIdx + 2
		createdStr := l.CreatedAt.Format("02/01/2006 15:04")

		values := []interface{}{
			l.ID,
			l.StudentName,
			l.ParentName,
			l.PhoneNumber,
			l.Grade,
			l.LearningModel,
			l.ClassType,
			l.AcademicPerformance,
			l.Status,
			l.ConsultedBy,
			createdStr,
		}

		for colIdx, val := range values {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, row)
			_ = f.SetCellValue(sheet, cell, val)
			_ = f.SetCellStyle(sheet, cell, cell, dataStyle)
		}
		_ = f.SetRowHeight(sheet, row, 22)
	}

	colWidths := map[string]float64{
		"A": 10, // ID
		"B": 24, // Họ Tên Học Sinh
		"C": 24, // Họ Tên Phụ Huynh
		"D": 18, // Số Điện Thoại
		"E": 14, // Khối Lớp
		"F": 20, // Hình Thức
		"G": 22, // Loại Lớp
		"H": 20, // Học Lực
		"I": 18, // Trạng Thái
		"J": 20, // Người Tư Vấn
		"K": 20, // Ngày Đăng Ký
	}
	for col, width := range colWidths {
		_ = f.SetColWidth(sheet, col, col, width)
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GenerateLeadTemplateExcel tạo file Excel mẫu (.xlsx) để nhập học sinh / phụ huynh hàng loạt
func GenerateLeadTemplateExcel() ([]byte, error) {
	sampleLeads := []*domain.Lead{
		{
			ID:                  1,
			StudentName:         "Nguyễn Minh Đức",
			ParentName:          "Nguyễn Thị Lan",
			PhoneNumber:         "0983459912",
			Grade:               5,
			LearningModel:       "online",
			ClassType:           "basic",
			AcademicPerformance: "good",
			Status:              "pending",
			CreatedAt:           time.Now(),
		},
		{
			ID:                  2,
			StudentName:         "Trần Khánh An",
			ParentName:          "Trần Văn Hùng",
			PhoneNumber:         "0976543210",
			Grade:               4,
			LearningModel:       "offline",
			ClassType:           "advanced",
			AcademicPerformance: "excellent",
			Status:              "contacted",
			ConsultedBy:         "Cô Minh Châu",
			CreatedAt:           time.Now(),
		},
	}
	return ExportLeadsToExcel(sampleLeads)
}

// ParseLeadsFromExcel đọc dữ liệu học sinh từ Excel, validate chặt chẽ theo luật nghiệp vụ domain.Lead
func ParseLeadsFromExcel(r io.Reader) ([]*domain.Lead, []string, error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, nil, fmt.Errorf("không thể đọc tệp tin Excel: %v", err)
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, nil, fmt.Errorf("tệp tin Excel không có trang tính nào")
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, nil, fmt.Errorf("lỗi đọc dữ liệu trang tính: %v", err)
	}

	if len(rows) <= 1 {
		return nil, nil, fmt.Errorf("tệp tin Excel không có dữ liệu để nhập (chỉ có dòng tiêu đề hoặc trống)")
	}

	var leads []*domain.Lead
	var validationErrors []string

	for i := 1; i < len(rows); i++ {
		rowNum := i + 1
		row := rows[i]

		if len(row) == 0 || isEmptyRow(row) {
			continue
		}

		getCol := func(idx int) string {
			if idx < len(row) {
				return strings.TrimSpace(row[idx])
			}
			return ""
		}

		idStr := getCol(0)
		studentName := getCol(1)
		parentName := getCol(2)
		phoneNumber := getCol(3)
		gradeStr := getCol(4)
		model := strings.ToLower(getCol(5))
		classType := strings.ToLower(getCol(6))
		perf := strings.ToLower(getCol(7))
		status := strings.ToLower(getCol(8))
		consultedBy := getCol(9)

		// 1. Kiểm tra Họ Tên
		if studentName == "" {
			validationErrors = append(validationErrors, fmt.Sprintf("Dòng %d: Họ tên Học sinh không được để trống", rowNum))
			continue
		}
		if parentName == "" {
			validationErrors = append(validationErrors, fmt.Sprintf("Dòng %d (Học sinh '%s'): Họ tên Phụ huynh không được để trống", rowNum, studentName))
			continue
		}

		// 2. Kiểm tra Số Điện Thoại Việt Nam
		phoneNumber = strings.ReplaceAll(phoneNumber, " ", "")
		phoneNumber = strings.ReplaceAll(phoneNumber, ".", "")
		phoneNumber = strings.ReplaceAll(phoneNumber, "-", "")
		if !phoneRegex.MatchString(phoneNumber) {
			validationErrors = append(validationErrors, fmt.Sprintf("Dòng %d (Học sinh '%s'): Số điện thoại '%s' không đúng định dạng Việt Nam (phải gồm 10 số bắt đầu bằng 03, 05, 07, 08, 09)", rowNum, studentName, phoneNumber))
			continue
		}

		// 3. Khối lớp
		grade, err := strconv.Atoi(gradeStr)
		if err != nil || grade < 1 || grade > 9 {
			validationErrors = append(validationErrors, fmt.Sprintf("Dòng %d (Học sinh '%s'): Khối lớp phải là số từ 1 đến 9 (nhận được: '%s')", rowNum, studentName, gradeStr))
			continue
		}

		// 4. Hình thức
		if model != "online" && model != "offline" {
			if strings.Contains(model, "on") {
				model = "online"
			} else {
				model = "offline"
			}
		}

		// 5. Loại lớp
		if classType != "basic" && classType != "advanced" && classType != "high_quality" {
			if strings.Contains(classType, "clc") || strings.Contains(classType, "chất lượng cao") {
				classType = "high_quality"
			} else if strings.Contains(classType, "nâng cao") || strings.Contains(classType, "nang cao") {
				classType = "advanced"
			} else {
				classType = "basic"
			}
		}

		// 6. Học lực
		if perf != "excellent" && perf != "good" && perf != "average" {
			if strings.Contains(perf, "giỏi") || strings.Contains(perf, "xuất sắc") {
				perf = "excellent"
			} else if strings.Contains(perf, "trung bình") {
				perf = "average"
			} else {
				perf = "good"
			}
		}

		// 7. Trạng thái
		if status != "pending" && status != "contacted" {
			status = "pending"
		}

		var id int64
		if idStr != "" {
			if parsedID, err := strconv.ParseInt(idStr, 10, 64); err == nil {
				id = parsedID
			}
		}

		leads = append(leads, &domain.Lead{
			ID:                  id,
			StudentName:         studentName,
			ParentName:          parentName,
			PhoneNumber:         phoneNumber,
			Grade:               grade,
			LearningModel:       model,
			ClassType:           classType,
			AcademicPerformance: perf,
			Status:              status,
			ConsultedBy:         consultedBy,
			CreatedAt:           time.Now(),
		})
	}

	return leads, validationErrors, nil
}

func isEmptyRow(row []string) bool {
	for _, cell := range row {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}
	return true
}
