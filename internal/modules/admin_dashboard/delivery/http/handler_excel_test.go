package http

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"toan-co-tra-backend/internal/domain"
	"toan-co-tra-backend/internal/modules/admin_dashboard/usecase"
	salesHttp "toan-co-tra-backend/internal/modules/sales/delivery/http"
	salesUsecase "toan-co-tra-backend/internal/modules/sales/usecase"
	"toan-co-tra-backend/internal/repository"
	"toan-co-tra-backend/pkg/excel"
)

func setupTestServer(t *testing.T) (*AdminHandler, *salesHttp.SalesHandler, string, func()) {
	tmpDir, err := os.MkdirTemp("", "excel_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	dbPath := filepath.Join(tmpDir, "database_test.json")

	repo := repository.NewCentralRepository(dbPath)
	adminUC := usecase.NewAdminUsecase(repo)
	salesUC := salesUsecase.NewSalesUsecase(repo)

	salesHandler := salesHttp.NewSalesHandler(salesUC, repo, nil)
	adminHandler := NewAdminHandler(adminUC, salesHandler)

	// Đăng nhập tài khoản admin1 để lấy token hợp lệ
	loginPayload := map[string]string{
		"username": "admin1",
		"password": "admin123",
	}
	body, _ := json.Marshal(loginPayload)
	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	salesHandler.Login(w, req)

	var loginResp struct {
		Success bool `json:"success"`
		Data    struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &loginResp)

	if !loginResp.Success || loginResp.Data.Token == "" {
		t.Fatalf("Admin login failed in test setup: %s", w.Body.String())
	}

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return adminHandler, salesHandler, loginResp.Data.Token, cleanup
}

func TestClassesExportAndTemplateHTTP(t *testing.T) {
	adminHandler, _, token, cleanup := setupTestServer(t)
	defer cleanup()

	// 1. Test Export Classes
	reqExport := httptest.NewRequest("GET", "/api/v1/admin/classes/export", nil)
	reqExport.Header.Set("Authorization", "Bearer "+token)
	wExport := httptest.NewRecorder()
	adminHandler.ExportClassesExcel(wExport, reqExport)

	if wExport.Code != http.StatusOK {
		t.Fatalf("ExportClassesExcel failed with status %d: %s", wExport.Code, wExport.Body.String())
	}
	if wExport.Header().Get("Content-Type") != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		t.Errorf("Unexpected Content-Type: %s", wExport.Header().Get("Content-Type"))
	}

	// 2. Test Download Template
	reqTmpl := httptest.NewRequest("GET", "/api/v1/admin/classes/template", nil)
	reqTmpl.Header.Set("Authorization", "Bearer "+token)
	wTmpl := httptest.NewRecorder()
	adminHandler.DownloadClassesTemplate(wTmpl, reqTmpl)

	if wTmpl.Code != http.StatusOK {
		t.Fatalf("DownloadClassesTemplate failed with status %d: %s", wTmpl.Code, wTmpl.Body.String())
	}
}

func TestClassesImportHTTP(t *testing.T) {
	adminHandler, _, token, cleanup := setupTestServer(t)
	defer cleanup()

	// Tạo dữ liệu lớp học mới qua Excel
	newClasses := []*domain.Class{
		{
			ID:        "excel_class_1",
			Name:      "Toán Thí Nghiệm Excel 1",
			Grade:     5,
			Model:     "online",
			Type:      "high_quality",
			Price:     "1.800.000đ/tháng",
			IsPopular: true,
			Desc:      "Lớp học thử nghiệm từ Excel import",
		},
		{
			ID:        "excel_class_2",
			Name:      "Toán Thí Nghiệm Excel 2",
			Grade:     4,
			Model:     "offline",
			Type:      "basic",
			Price:     "900.000đ/tháng",
			IsPopular: false,
			Desc:      "Lớp cơ bản",
		},
	}

	excelBytes, err := excel.ExportClassesToExcel(newClasses)
	if err != nil {
		t.Fatalf("Failed to create excel bytes: %v", err)
	}

	// Gửi multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "classes.xlsx")
	if err != nil {
		t.Fatalf("CreateFormFile failed: %v", err)
	}
	_, _ = part.Write(excelBytes)
	_ = writer.Close()

	req := httptest.NewRequest("POST", "/api/v1/admin/classes/import", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	adminHandler.ImportClassesExcel(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("ImportClassesExcel failed with code %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Success  bool `json:"success"`
		Inserted int  `json:"inserted"`
		Updated  int  `json:"updated"`
		Total    int  `json:"total"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	if !resp.Success || resp.Total != 2 || resp.Inserted != 2 {
		t.Errorf("Unexpected import response: %+v", resp)
	}
}

func TestLeadsExportAndImportHTTP(t *testing.T) {
	adminHandler, _, token, cleanup := setupTestServer(t)
	defer cleanup()

	// 1. Test Leads Export
	reqExport := httptest.NewRequest("GET", "/api/v1/admin/leads/export", nil)
	reqExport.Header.Set("Authorization", "Bearer "+token)
	wExport := httptest.NewRecorder()
	adminHandler.ExportLeadsExcel(wExport, reqExport)

	if wExport.Code != http.StatusOK {
		t.Fatalf("ExportLeadsExcel failed with code %d: %s", wExport.Code, wExport.Body.String())
	}

	// 2. Test Leads Import
	newLeads := []*domain.Lead{
		{
			StudentName:         "Học Sinh Mới Excel",
			ParentName:          "Phụ Huynh Mới Excel",
			PhoneNumber:         "0983111222",
			Grade:               5,
			LearningModel:       "online",
			ClassType:           "basic",
			AcademicPerformance: "good",
		},
	}

	excelBytes, err := excel.ExportLeadsToExcel(newLeads)
	if err != nil {
		t.Fatalf("Failed to export leads: %v", err)
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "leads.xlsx")
	_, _ = part.Write(excelBytes)
	_ = writer.Close()

	reqImport := httptest.NewRequest("POST", "/api/v1/admin/leads/import", body)
	reqImport.Header.Set("Content-Type", writer.FormDataContentType())
	reqImport.Header.Set("Authorization", "Bearer "+token)
	wImport := httptest.NewRecorder()

	adminHandler.ImportLeadsExcel(wImport, reqImport)

	if wImport.Code != http.StatusOK {
		t.Fatalf("ImportLeadsExcel failed with code %d: %s", wImport.Code, wImport.Body.String())
	}

	var resp struct {
		Success  bool `json:"success"`
		Inserted int  `json:"inserted"`
	}
	if err := json.Unmarshal(wImport.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	if !resp.Success || resp.Inserted != 1 {
		t.Errorf("Unexpected lead import response: %+v", resp)
	}
}

