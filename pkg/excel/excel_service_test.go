package excel

import (
	"bytes"
	"testing"
	"toan-co-tra-backend/internal/domain"

	"github.com/xuri/excelize/v2"
)

func TestExportAndParseClasses(t *testing.T) {
	sampleClasses := []*domain.Class{
		{
			ID:        "class_5_clc_test",
			Name:      "Toán 5 CLC Test",
			Grade:     5,
			Model:     "offline",
			Type:      "high_quality",
			Price:     "1.500.000đ/tháng",
			IsPopular: true,
			Desc:      "Lớp ôn thi chuyên",
		},
		{
			ID:        "class_4_adv_test",
			Name:      "Toán 4 Nâng Cao Test",
			Grade:     4,
			Model:     "online",
			Type:      "advanced",
			Price:     "1.200.000đ/tháng",
			IsPopular: false,
			Desc:      "Học trực tuyến",
		},
	}

	data, err := ExportClassesToExcel(sampleClasses)
	if err != nil {
		t.Fatalf("ExportClassesToExcel failed: %v", err)
	}

	if len(data) == 0 {
		t.Fatalf("Exported data is empty")
	}

	// Đọc lại từ buffer
	classes, errs, err := ParseClassesFromExcel(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("ParseClassesFromExcel failed: %v", err)
	}

	if len(errs) > 0 {
		t.Errorf("Expected 0 errors, got %v", errs)
	}

	if len(classes) != 2 {
		t.Fatalf("Expected 2 classes, got %d", len(classes))
	}

	if classes[0].ID != "class_5_clc_test" || classes[0].Name != "Toán 5 CLC Test" {
		t.Errorf("First class mismatch: %+v", classes[0])
	}
	if !classes[0].IsPopular {
		t.Errorf("Expected IsPopular to be true")
	}

	if classes[1].ID != "class_4_adv_test" || classes[1].Grade != 4 {
		t.Errorf("Second class mismatch: %+v", classes[1])
	}
}

func TestClassValidationErrors(t *testing.T) {
	// Tạo file Excel với một số dòng lỗi
	f := excelize.NewFile()
	defer f.Close()
	sheet := "Sheet1"

	headers := []string{"Mã Lớp", "Tên Lớp", "Khối Lớp", "Hình Thức", "Loại Lớp", "Học Phí", "Mũi Nhọn", "Mô Tả"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		_ = f.SetCellValue(sheet, cell, h)
	}

	// Dòng 2: Hợp lệ
	row2 := []interface{}{"c1", "Lớp 1 Cơ Bản", 1, "offline", "basic", "1.000.000đ", "Không", "Mô tả"}
	// Dòng 3: Thiếu Tên Lớp
	row3 := []interface{}{"c2", "", 2, "offline", "basic", "1.000.000đ", "Không", "Mô tả"}
	// Dòng 4: Khối lớp sai (khối 12)
	row4 := []interface{}{"c3", "Lớp Sai Khối", 12, "offline", "basic", "1.000.000đ", "Không", "Mô tả"}
	// Dòng 5: CLC nhưng cho Khối 3 (vi phạm CLC chỉ dành cho Khối 5)
	row5 := []interface{}{"c4", "Lớp CLC Khối 3", 3, "offline", "high_quality", "1.000.000đ", "Không", "Mô tả"}

	for col, val := range row2 {
		cell, _ := excelize.CoordinatesToCellName(col+1, 2)
		_ = f.SetCellValue(sheet, cell, val)
	}
	for col, val := range row3 {
		cell, _ := excelize.CoordinatesToCellName(col+1, 3)
		_ = f.SetCellValue(sheet, cell, val)
	}
	for col, val := range row4 {
		cell, _ := excelize.CoordinatesToCellName(col+1, 4)
		_ = f.SetCellValue(sheet, cell, val)
	}
	for col, val := range row5 {
		cell, _ := excelize.CoordinatesToCellName(col+1, 5)
		_ = f.SetCellValue(sheet, cell, val)
	}

	var buf bytes.Buffer
	_ = f.Write(&buf)

	classes, errs, err := ParseClassesFromExcel(&buf)
	if err != nil {
		t.Fatalf("ParseClassesFromExcel returned error: %v", err)
	}

	if len(classes) != 1 {
		t.Errorf("Expected 1 valid class, got %d", len(classes))
	}

	if len(errs) != 3 {
		t.Errorf("Expected 3 validation errors, got %d: %v", len(errs), errs)
	}
}

func TestExportAndParseLeads(t *testing.T) {
	sampleLeads := []*domain.Lead{
		{
			StudentName:         "Nguyễn Văn A",
			ParentName:          "Nguyễn Văn B",
			PhoneNumber:         "0983123456",
			Grade:               5,
			LearningModel:       "online",
			ClassType:           "basic",
			AcademicPerformance: "good",
		},
	}

	data, err := ExportLeadsToExcel(sampleLeads)
	if err != nil {
		t.Fatalf("ExportLeadsToExcel failed: %v", err)
	}

	leads, errs, err := ParseLeadsFromExcel(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("ParseLeadsFromExcel failed: %v", err)
	}

	if len(errs) > 0 {
		t.Errorf("Expected 0 errors, got %v", errs)
	}

	if len(leads) != 1 {
		t.Fatalf("Expected 1 lead, got %d", len(leads))
	}

	if leads[0].PhoneNumber != "0983123456" {
		t.Errorf("Phone mismatch: %s", leads[0].PhoneNumber)
	}
}

func TestLeadValidationErrors(t *testing.T) {
	f := excelize.NewFile()
	defer f.Close()
	sheet := "Sheet1"

	headers := []string{"ID", "Học Sinh", "Phụ Huynh", "SĐT", "Khối", "Hình Thức", "Loại", "Học Lực"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		_ = f.SetCellValue(sheet, cell, h)
	}

	// Dòng 2: Hợp lệ
	row2 := []interface{}{1, "Học Sinh A", "Phụ Huynh A", "0912345678", 5, "online", "basic", "good"}
	// Dòng 3: SĐT sai định dạng (9 số)
	row3 := []interface{}{2, "Học Sinh B", "Phụ Huynh B", "091234567", 5, "online", "basic", "good"}
	// Dòng 4: Thiếu tên học sinh
	row4 := []interface{}{3, "", "Phụ Huynh C", "0987654321", 5, "online", "basic", "good"}

	for col, val := range row2 {
		cell, _ := excelize.CoordinatesToCellName(col+1, 2)
		_ = f.SetCellValue(sheet, cell, val)
	}
	for col, val := range row3 {
		cell, _ := excelize.CoordinatesToCellName(col+1, 3)
		_ = f.SetCellValue(sheet, cell, val)
	}
	for col, val := range row4 {
		cell, _ := excelize.CoordinatesToCellName(col+1, 4)
		_ = f.SetCellValue(sheet, cell, val)
	}

	var buf bytes.Buffer
	_ = f.Write(&buf)

	leads, errs, err := ParseLeadsFromExcel(&buf)
	if err != nil {
		t.Fatalf("ParseLeadsFromExcel failed: %v", err)
	}

	if len(leads) != 1 {
		t.Errorf("Expected 1 valid lead, got %d", len(leads))
	}

	if len(errs) != 2 {
		t.Errorf("Expected 2 errors, got %d: %v", len(errs), errs)
	}
}
