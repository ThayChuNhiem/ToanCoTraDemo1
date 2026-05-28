import re

print("Running DOM & JS upgrades on playground.html...")

with open("playground.html", "r", encoding="utf-8") as f:
    html = f.read()

# 1. Replace the static Class Price Editor form with dynamic Class CRUD Card
old_price_card = """                    <!-- 1. Tùy chỉnh học phí lớp học -->
                    <div class="editor-card">
                        <h3>Điều Chỉnh Học Phí Lớp Học</h3>
                        <p style="font-size: 11px; color: var(--text-muted); margin-bottom: 12px;">Cập nhật học phí trực quan cho toàn bộ các lớp học hoạt động trên hệ thống</p>
                        <form id="admin-form-prices">
                            <div id="admin-classes-price-list" style="display: flex; flex-direction: column; gap: 8px;">
                                <!-- Load dong danh sach lop hoc -->
                            </div>
                            <button type="submit" class="btn-action-small" style="width: 100%; margin-top: 15px;">Lưu Thay Đổi Học Phí</button>
                        </form>
                    </div>"""

new_price_card = """                    <!-- 1. Quản Lý Khóa Học & Chương Trình Mũi Nhọn -->
                    <div class="editor-card" style="grid-column: span 2;">
                        <h3>Quản Lý Danh Mục Khóa Học & Chương Trình Mũi Nhọn</h3>
                        <p style="font-size: 11px; color: var(--text-muted); margin-bottom: 12px;">Thêm mới lớp học, cập nhật thông tin học phí, hình thức học, và đánh dấu Chương Trình Đào Tạo Mũi Nhọn hiển thị nổi bật ở Trang Chủ.</p>
                        
                        <div style="background-color: var(--bg-snowy); padding: 15px; border-radius: 12px; margin-bottom: 15px; border: var(--glass-border);">
                            <h4 id="class-form-title" style="font-size: 14px; margin-bottom: 10px; font-weight: 600;">Thêm Lớp Học Mới</h4>
                            <form id="admin-form-manage-class" style="margin-top: 10px;">
                                <input type="hidden" id="class-edit-id" value="">
                                <div style="display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 12px; margin-bottom: 10px;">
                                    <div class="form-group" style="margin-bottom: 0;">
                                        <label style="font-size: 11px;">Mã định danh lớp (ID)</label>
                                        <input type="text" id="class-id-input" class="form-control" style="padding: 6px 10px; font-size: 12px;" placeholder="Ví dụ: class_4_adv" required>
                                    </div>
                                    <div class="form-group" style="margin-bottom: 0;">
                                        <label style="font-size: 11px;">Tên lớp học</label>
                                        <input type="text" id="class-name-input" class="form-control" style="padding: 6px 10px; font-size: 12px;" placeholder="Toán Tư Duy Lớp 4 Nâng Cao" required>
                                    </div>
                                    <div class="form-group" style="margin-bottom: 0;">
                                        <label style="font-size: 11px;">Học phí niêm yết</label>
                                        <input type="text" id="class-price-input" class="form-control" style="padding: 6px 10px; font-size: 12px;" placeholder="Ví dụ: 1.200.000đ/tháng" required>
                                    </div>
                                </div>
                                <div style="display: grid; grid-template-columns: 1fr 1fr 1fr 1fr; gap: 12px; margin-bottom: 10px;">
                                    <div class="form-group" style="margin-bottom: 0;">
                                        <label style="font-size: 11px;">Khối lớp (1-9)</label>
                                        <input type="number" id="class-grade-input" class="form-control" style="padding: 6px 10px; font-size: 12px;" min="1" max="9" value="5" required>
                                    </div>
                                    <div class="form-group" style="margin-bottom: 0;">
                                        <label style="font-size: 11px;">Hình thức học</label>
                                        <select id="class-model-input" class="form-control" style="padding: 6px 10px; font-size: 12px; height: 34px;">
                                            <option value="online">Online (Trực tuyến)</option>
                                            <option value="offline">Offline (Tại lớp)</option>
                                        </select>
                                    </div>
                                    <div class="form-group" style="margin-bottom: 0;">
                                        <label style="font-size: 11px;">Loại lớp</label>
                                        <select id="class-type-input" class="form-control" style="padding: 6px 10px; font-size: 12px; height: 34px;">
                                            <option value="basic">Cơ Bản (basic)</option>
                                            <option value="advanced">Nâng Cao (advanced)</option>
                                            <option value="high_quality">Chất Lượng Cao (high_quality)</option>
                                        </select>
                                    </div>
                                    <div class="form-group" style="margin-bottom: 0; display: flex; align-items: center; gap: 8px; padding-top: 15px;">
                                        <input type="checkbox" id="class-popular-input" style="width: 18px; height: 18px; cursor: pointer;">
                                        <label for="class-popular-input" style="font-size: 11px; font-weight: 600; cursor: pointer; color: var(--accent-green);">Chương trình Mũi nhọn</label>
                                    </div>
                                </div>
                                <div class="form-group" style="margin-bottom: 12px;">
                                    <label style="font-size: 11px;">Mô tả ngắn gọn lớp học</label>
                                    <input type="text" id="class-desc-input" class="form-control" style="padding: 6px 10px; font-size: 12px;" placeholder="Mô tả tóm tắt khóa học..." required>
                                </div>
                                <div style="display: flex; gap: 10px;">
                                    <button type="submit" id="btn-save-class" class="btn-action-small" style="padding: 8px 15px; background: var(--accent-green); color: white;">Lưu Lớp Học</button>
                                    <button type="button" id="btn-cancel-class-edit" class="btn-action-small" style="padding: 8px 15px; background: var(--text-muted); color: white; display: none;" onclick="resetClassForm()">Hủy Sửa</button>
                                </div>
                            </form>
                        </div>

                        <h4 style="font-size: 14px; margin-bottom: 8px; font-weight: 600;">Danh Sách Lớp Học Hệ Thống</h4>
                        <div style="max-height: 250px; overflow-y: auto; overflow-x: auto; border: var(--glass-border); border-radius: 8px; background: white;">
                            <table class="admin-list-table" style="margin-top: 0; font-size: 12px;">
                                <thead>
                                    <tr>
                                        <th>Mã Lớp</th>
                                        <th>Tên Lớp</th>
                                        <th>Khối</th>
                                        <th>Hình Thức</th>
                                        <th>Loại Lớp</th>
                                        <th>Học Phí</th>
                                        <th>Mũi Nhọn</th>
                                        <th>Thao Tác</th>
                                    </tr>
                                </thead>
                                <tbody id="admin-table-classes">
                                    <!-- Load dong danh sach lop hoc -->
                                </tbody>
                            </table>
                        </div>
                    </div>"""

if old_price_card in html:
    html = html.replace(old_price_card, new_price_card)
    print("SUCCESS: Price static card upgraded to full Class CRUD table successfully!")
else:
    print("Warning: old price static card not found directly.")

# 2. Add Mobile Responsive wraps (overflow-x: auto) around all dashboard tables
# Wrap Honored Students table
old_honored_table = """                        <h4>Danh Sách Học Sinh Vinh Danh Hiện Tại</h4>
                        <table class="admin-list-table">"""

new_honored_table = """                        <h4>Danh Sách Học Sinh Vinh Danh Hiện Tại</h4>
                        <div style="overflow-x: auto; max-width: 100%; border: var(--glass-border); border-radius: 8px; background: white;">
                            <table class="admin-list-table" style="margin-top: 0;">"""

# Note: We need to replace the ending table tag as well or just wrap it cleanly. Wrapping the inner table directly is much cleaner!
if old_honored_table in html:
    html = html.replace(old_honored_table, new_honored_table)
    # Append closing div for honored table wrapper
    old_end_honored_table = """                            <tbody id="admin-table-students">
                                <!-- Load động danh sách học sinh kèm nút xóa -->
                            </tbody>
                        </table>"""
    new_end_honored_table = """                            <tbody id="admin-table-students">
                                <!-- Load động danh sách học sinh kèm nút xóa -->
                            </tbody>
                        </table>
                        </div>"""
    html = html.replace(old_end_honored_table, new_end_honored_table)
    print("SUCCESS: Honored students table mobile responsiveness updated!")
else:
    print("Warning: old honored table not found.")

# Update Teachers Table container to support horizontal scroll as well
old_teachers_container = """                        <h4 style="font-size: 14px; margin-bottom: 8px; font-weight: 600;">Danh Sách Giảng Viên Trang Giới Thiệu</h4>
                        <div style="max-height: 250px; overflow-y: auto; border: var(--glass-border); border-radius: 8px; background: white;">"""

new_teachers_container = """                        <h4 style="font-size: 14px; margin-bottom: 8px; font-weight: 600;">Danh Sách Giảng Viên Trang Giới Thiệu</h4>
                        <div style="max-height: 250px; overflow-y: auto; overflow-x: auto; border: var(--glass-border); border-radius: 8px; background: white;">"""

if old_teachers_container in html:
    html = html.replace(old_teachers_container, new_teachers_container)
    print("SUCCESS: Teachers table mobile responsiveness updated!")
else:
    print("Warning: old teachers container not found.")

# Update Sales Accounts Table container to support horizontal scroll as well
old_users_container = """                        <h4 style="font-size: 14px; margin-bottom: 8px; font-weight: 600;">Danh Sách Nhân Viên Sales Hiện Tại</h4>
                        <div style="max-height: 200px; overflow-y: auto; border: var(--glass-border); border-radius: 8px; background: white;">"""

new_users_container = """                        <h4 style="font-size: 14px; margin-bottom: 8px; font-weight: 600;">Danh Sách Nhân Viên Sales Hiện Tại</h4>
                        <div style="max-height: 200px; overflow-y: auto; overflow-x: auto; border: var(--glass-border); border-radius: 8px; background: white;">"""

if old_users_container in html:
    html = html.replace(old_users_container, new_users_container)
    print("SUCCESS: Sales accounts table mobile responsiveness updated!")
else:
    print("Warning: old users container not found.")

# Save modified html back
with open("playground.html", "w", encoding="utf-8") as f:
    f.write(html)

print("Upgrades complete on playground.html!")
