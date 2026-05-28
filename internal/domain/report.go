package domain

// SalesReport đại diện cho thông tin báo cáo thống kê hiển thị trên Dashboard Admin
type SalesReport struct {
	TotalLeads        int64            `json:"total_leads"`
	ContactedLeads    int64            `json:"contacted_leads"`
	PendingLeads      int64            `json:"pending_leads"`
	OnlinePercentage  float64          `json:"online_percentage"`
	OfflinePercentage float64          `json:"offline_percentage"`
	GradeDistribution map[int]int      `json:"grade_distribution"` // Phân bố số lượng lead theo khối lớp 1-9
}
