import re

print("Running HTML upgrade on playground.html...")

with open("playground.html", "r", encoding="utf-8") as f:
    html = f.read()

# 1. Update Navigation Links
old_nav = """            <ul class="nav-links">
                <li><button id="btn-homepage" class="active" onclick="switchTab('tab-homepage')">Trang Chủ</button></li>
                <li><button id="btn-about" onclick="switchTab('tab-about')">Giới Thiệu</button></li>
                <li><button id="btn-honor" onclick="switchTab('tab-honor')">Vinh Danh</button></li>
                <li><button id="btn-sales" onclick="switchTab('tab-sales')">Quản Trị Sales</button></li>
                <li><button id="btn-admin" onclick="switchTab('tab-admin')">Bảng Admin</button></li>
            </ul>"""

new_nav = """            <ul class="nav-links">
                <li><button id="btn-homepage" class="active" onclick="switchTab('tab-homepage')">Trang Chủ</button></li>
                <li><button id="btn-about" onclick="switchTab('tab-about')">Giới Thiệu</button></li>
                <li><button id="btn-honor" onclick="switchTab('tab-honor')">Vinh Danh</button></li>
                <li id="li-sales-dashboard" style="display: none;"><button id="btn-sales" onclick="switchTab('tab-sales')">Quản Trị Sales</button></li>
                <li id="li-admin-dashboard" style="display: none;"><button id="btn-admin" onclick="switchTab('tab-admin')">Bảng Admin</button></li>
                <li id="li-login"><button id="btn-login" onclick="openLoginModal()" style="background: rgba(16, 185, 129, 0.1); color: var(--accent-green); border: var(--glass-border); padding: 8px 16px; border-radius: 20px; font-weight: 600; cursor: pointer;">Đăng Nhập</button></li>
                <li id="li-user-info" style="display: none; align-items: center; gap: 10px;">
                    <span id="nav-username-badge" class="user-badge" style="font-weight: 600; font-size: 13px;"></span>
                    <button onclick="handleLogout()" class="btn-logout" style="padding: 6px 12px; font-size: 12px; background: rgba(239, 68, 68, 0.1); color: #EF4444; border: none; border-radius: 8px; cursor: pointer;">Đăng Xuất</button>
                </li>
            </ul>"""

if old_nav in html:
    html = html.replace(old_nav, new_nav)
    print("SUCCESS: Navigation Links updated successfully!")
else:
    print("Warning: old nav block not found directly.")

# 2. Update Tab Sales Lock Screen
old_sales_lock = """            <div id="sales-login-section" class="login-overlay-container">
                <h3>Đăng Nhập Đội Ngũ Sales</h3>
                <form id="sales-login-form">
                    <div class="form-group">
                        <label for="sales-user">Tên đăng nhập</label>
                        <input type="text" id="sales-user" class="form-control" placeholder="Gợi ý: teacher1" required>
                    </div>
                    <div class="form-group">
                        <label for="sales-pass">Mật khẩu</label>
                        <input type="password" id="sales-pass" class="form-control" placeholder="Mật khẩu (password123)" required>
                    </div>
                    <button type="submit" class="btn-submit">Đăng Nhập Hệ Thống</button>
                </form>
            </div>"""

new_sales_lock = """            <!-- Lock screen if not logged in -->
            <div id="sales-login-section" class="login-overlay-container" style="text-align: center; padding: 40px 20px;">
                <div style="font-size: 48px; margin-bottom: 15px;">🔒</div>
                <h3>Yêu Cầu Đăng Nhập Hệ Thống</h3>
                <p style="color: var(--text-muted); margin-bottom: 20px; max-width: 400px; margin-left: auto; margin-right: auto;">Vui lòng đăng nhập bằng tài khoản Sales hoặc Giáo Viên của bạn để sử dụng phân hệ này.</p>
                <button onclick="openLoginModal()" class="btn-submit" style="width: auto; padding: 10px 25px;">Đăng Nhập Ngay</button>
            </div>"""

if old_sales_lock in html:
    html = html.replace(old_sales_lock, new_sales_lock)
    print("SUCCESS: Tab Sales lock screen updated successfully!")
else:
    print("Warning: old sales lock block not found directly.")

# 3. Update Tab Admin Lock Screen
old_admin_lock = """            <div id="admin-login-section" class="login-overlay-container">
                <h3>Đăng Nhập Tài Khoản Admin</h3>
                <form id="admin-login-form">
                    <div class="form-group">
                        <label for="admin-user">Tên đăng nhập</label>
                        <input type="text" id="admin-user" class="form-control" placeholder="Gợi ý: admin1" required>
                    </div>
                    <div class="form-group">
                        <label for="admin-pass">Mật khẩu</label>
                        <input type="password" id="admin-pass" class="form-control" placeholder="Mật khẩu (admin123)" required>
                    </div>
                    <button type="submit" class="btn-submit">Đăng Nhập Admin</button>
                </form>
            </div>"""

new_admin_lock = """            <!-- Lock screen if not logged in -->
            <div id="admin-login-section" class="login-overlay-container" style="text-align: center; padding: 40px 20px;">
                <div style="font-size: 48px; margin-bottom: 15px;">🛡️</div>
                <h3>Yêu Cầu Quyền Quản Trị Viên</h3>
                <p style="color: var(--text-muted); margin-bottom: 20px; max-width: 400px; margin-left: auto; margin-right: auto;">Mục này chỉ dành cho Admin tối cao của Toán Cô Trà. Vui lòng đăng nhập tài khoản Admin của bạn.</p>
                <button onclick="openLoginModal()" class="btn-submit" style="width: auto; padding: 10px 25px;">Đăng Nhập Quyền Admin</button>
            </div>"""

if old_admin_lock in html:
    html = html.replace(old_admin_lock, new_admin_lock)
    print("SUCCESS: Tab Admin lock screen updated successfully!")
else:
    print("Warning: old admin lock block not found directly.")

# 4. Update Price Editor Card
old_price_editor = """                    <!-- 1. Tùy chỉnh học phí lớp học -->
                    <div class="editor-card">
                        <h3>Điều Chỉnh Học Phí Lớp Học</h3>
                        <form id="admin-form-prices">
                            <div class="class-price-row">
                                <span>Lớp 5 Cơ Bản:</span>
                                <input type="text" id="price-basic" class="form-control" style="width: 160px;" required>
                            </div>
                            <div class="class-price-row">
                                <span>Lớp 5 Nâng Cao:</span>
                                <input type="text" id="price-advanced" class="form-control" style="width: 160px;" required>
                            </div>
                            <div class="class-price-row">
                                <span>CLC Ôn Luyện Chuyên:</span>
                                <input type="text" id="price-clc" class="form-control" style="width: 160px;" required>
                            </div>
                            <button type="submit" class="btn-action-small" style="width: 100%; margin-top: 15px;">Lưu Thay Đổi Học Phí</button>
                        </form>
                    </div>"""

new_price_editor = """                    <!-- 1. Tùy chỉnh học phí lớp học -->
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

if old_price_editor in html:
    html = html.replace(old_price_editor, new_price_editor)
    print("SUCCESS: Price editor card updated successfully!")
else:
    print("Warning: old price editor block not found directly.")

# 5. Replace User Creation Form with Sales Account Management Card
old_user_creator = """                    <!-- 2. Tạo tài khoản giáo viên/sales mới -->
                    <div class="editor-card">
                        <h3>Cấp Tài Khoản Sales / Giáo Viên Mới</h3>
                        <p style="font-size: 11px; color: var(--text-muted); margin-bottom: 12px;">Chỉ Admin mới có thẩm quyền thiết lập tài khoản tuyển sinh bổ sung</p>
                        <form id="admin-form-create-user">
                            <div class="form-group" style="margin-bottom: 10px;">
                                <label>Tên tài khoản mới</label>
                                <input type="text" id="new-user-name" class="form-control" placeholder="Ví dụ: teacher2" required>
                            </div>
                            <div class="form-group" style="margin-bottom: 10px;">
                                <label>Mật khẩu</label>
                                <input type="password" id="new-user-pass" class="form-control" placeholder="Mật khẩu" required>
                            </div>
                            <div class="form-group" style="margin-bottom: 15px;">
                                <label>Họ tên đầy đủ</label>
                                <input type="text" id="new-user-fullname" class="form-control" placeholder="Ví dụ: Thầy Hoàng Lâm" required>
                            </div>
                            <button type="submit" class="btn-action-small" style="width: 100%; background-color: var(--accent-green); color: white;">Tạo & Kích Hoạt Tài Khoản</button>
                        </form>
                    </div>"""

new_user_creator = """                    <!-- 2. Quan ly Tai Khoan Sales / Tuyen Sinh -->
                    <div class="editor-card" style="grid-column: span 2;">
                        <h3>Quản Lý Tài Khoản Sales & Đội Ngũ Tuyển Sinh</h3>
                        <p style="font-size: 11px; color: var(--text-muted); margin-bottom: 12px;">Cấp mới và thu hồi tài khoản chăm sóc phụ huynh cho nhân viên sales</p>
                        
                        <div style="background-color: var(--bg-snowy); padding: 15px; border-radius: 12px; margin-bottom: 15px; border: var(--glass-border);">
                            <h4 style="font-size: 14px; margin-bottom: 10px; font-weight: 600;">Cấp Tài Khoản Mới</h4>
                            <form id="admin-form-create-sales-user" style="display: flex; gap: 10px; align-items: flex-end;">
                                <div class="form-group" style="margin-bottom: 0; flex: 1;">
                                    <label style="font-size: 11px; margin-bottom: 3px;">Tên tài khoản</label>
                                    <input type="text" id="new-sales-username" class="form-control" style="padding: 6px 10px; font-size: 12px;" placeholder="teacher2" required>
                                </div>
                                <div class="form-group" style="margin-bottom: 0; flex: 1;">
                                    <label style="font-size: 11px; margin-bottom: 3px;">Mật khẩu</label>
                                    <input type="password" id="new-sales-password" class="form-control" style="padding: 6px 10px; font-size: 12px;" required>
                                </div>
                                <div class="form-group" style="margin-bottom: 0; flex: 150px;">
                                    <label style="font-size: 11px; margin-bottom: 3px;">Họ tên đầy đủ</label>
                                    <input type="text" id="new-sales-fullname" class="form-control" style="padding: 6px 10px; font-size: 12px;" placeholder="Cô Minh Châu" required>
                                </div>
                                <button type="submit" class="btn-action-small" style="padding: 8px 15px; height: 34px; background: var(--accent-green); color: white;">Tạo Tài Khoản</button>
                            </form>
                        </div>

                        <h4 style="font-size: 14px; margin-bottom: 8px; font-weight: 600;">Danh Sách Nhân Viên Sales Hiện Tại</h4>
                        <div style="max-height: 200px; overflow-y: auto; border: var(--glass-border); border-radius: 8px; background: white;">
                            <table class="admin-list-table" style="margin-top: 0; font-size: 12px;">
                                <thead>
                                    <tr>
                                        <th>Tên Tài Khoản</th>
                                        <th>Họ & Tên</th>
                                        <th>Vai Trò</th>
                                        <th>Thao Tác</th>
                                    </tr>
                                </thead>
                                <tbody id="admin-table-users">
                                    <!-- Load dong danh sach users -->
                                </tbody>
                            </table>
                        </div>
                    </div>"""

if old_user_creator in html:
    html = html.replace(old_user_creator, new_user_creator)
    print("SUCCESS: User creator card updated to full sales account management successfully!")
else:
    print("Warning: old user creator block not found directly.")

# 6. Add Teacher CRUD Card right before Honored Students Card
old_students_card = """                    <!-- 4. Quản lý Học sinh Vinh danh -->"""

new_teacher_crud_card = """                    <!-- 3. Quan ly Doi Ngu Giao Vien & Cong Su -->
                    <div class="editor-card" style="grid-column: span 2;">
                        <h3>Quản Lý Đội Ngũ Giảng Viên & Cộng Sự (Trang Giới Thiệu)</h3>
                        <p style="font-size: 11px; color: var(--text-muted); margin-bottom: 12px;">Thêm mới, chỉnh sửa thông tin tiểu sử và ảnh đại diện của giáo viên hiển thị trên trang Giới Thiệu</p>
                        
                        <div style="background-color: var(--bg-snowy); padding: 15px; border-radius: 12px; margin-bottom: 15px; border: var(--glass-border);">
                            <h4 id="teacher-form-title" style="font-size: 14px; margin-bottom: 10px; font-weight: 600;">Thêm Mới Giảng Viên</h4>
                            <form id="admin-form-manage-teacher" style="margin-top: 10px;">
                                <input type="hidden" id="teacher-edit-id" value="">
                                <div style="display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 12px; margin-bottom: 10px;">
                                    <div class="form-group" style="margin-bottom: 0;">
                                        <label style="font-size: 11px;">Họ tên giáo viên</label>
                                        <input type="text" id="teacher-name" class="form-control" style="padding: 6px 10px; font-size: 12px;" placeholder="Cô Phạm Minh Châu" required>
                                    </div>
                                    <div class="form-group" style="margin-bottom: 0;">
                                        <label style="font-size: 11px;">Vai trò giảng dạy</label>
                                        <input type="text" id="teacher-role" class="form-control" style="padding: 6px 10px; font-size: 12px;" placeholder="Giáo viên Toán tư duy" required>
                                    </div>
                                    <div class="form-group" style="margin-bottom: 0;">
                                        <label style="font-size: 11px;">Trình độ học vấn</label>
                                        <input type="text" id="teacher-education" class="form-control" style="padding: 6px 10px; font-size: 12px;" placeholder="Cử nhân ĐH Sư Phạm Hà Nội" required>
                                    </div>
                                </div>
                                <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 12px; margin-bottom: 10px;">
                                    <div class="form-group" style="margin-bottom: 0;">
                                        <label style="font-size: 11px;">Thứ tự hiển thị (Số nhỏ xếp trước)</label>
                                        <input type="number" id="teacher-order" class="form-control" style="padding: 6px 10px; font-size: 12px; width: 100px;" value="5" required>
                                    </div>
                                    <div class="form-group" style="margin-bottom: 0;">
                                        <label style="font-size: 11px;">Tải Ảnh Giáo Viên Lên (Upload ảnh trực tiếp)</label>
                                        <div style="display: flex; gap: 10px; align-items: center;">
                                            <input type="file" id="teacher-avatar-file" class="form-control" style="padding: 4px; font-size: 11px;" accept="image/*">
                                            <input type="hidden" id="teacher-avatar-url">
                                        </div>
                                    </div>
                                </div>
                                <div class="form-group" style="margin-bottom: 12px;">
                                    <label style="font-size: 11px;">Tiểu sử giảng dạy tóm tắt</label>
                                    <textarea id="teacher-bio" class="form-control" rows="3" style="padding: 8px; font-size: 12px; resize: vertical;" placeholder="Nhập kinh nghiệm, phương pháp giảng dạy..." required></textarea>
                                </div>
                                <div style="display: flex; gap: 10px;">
                                    <button type="submit" id="btn-save-teacher" class="btn-action-small" style="padding: 8px 15px; background: var(--accent-green); color: white;">Lưu Giảng Viên</button>
                                    <button type="button" id="btn-cancel-teacher-edit" class="btn-action-small" style="padding: 8px 15px; background: var(--text-muted); color: white; display: none;" onclick="resetTeacherForm()">Hủy Chỉnh Sửa</button>
                                </div>
                            </form>
                        </div>

                        <h4 style="font-size: 14px; margin-bottom: 8px; font-weight: 600;">Danh Sách Giảng Viên Trang Giới Thiệu</h4>
                        <div style="max-height: 250px; overflow-y: auto; border: var(--glass-border); border-radius: 8px; background: white;">
                            <table class="admin-list-table" style="margin-top: 0; font-size: 12px;">
                                <thead>
                                    <tr>
                                        <th>Ảnh</th>
                                        <th>Họ & Tên</th>
                                        <th>Vai Trò</th>
                                        <th>Học Vấn</th>
                                        <th>Thứ Tự</th>
                                        <th>Thao Tác</th>
                                    </tr>
                                </thead>
                                <tbody id="admin-table-teachers">
                                    <!-- Load dong danh sach giao vien -->
                                </tbody>
                            </table>
                        </div>
                    </div>

                    <!-- 4. Quản lý Học sinh Vinh danh -->"""

if old_students_card in html:
    html = html.replace(old_students_card, new_teacher_crud_card)
    print("SUCCESS: Teacher CRUD card inserted successfully!")
else:
    print("Warning: old students card comment not found directly.")

# 7. Add unified login modal before Chatbot Widget
old_chatbot = """    <!-- Chatbot Widget góc màn hình -->"""

new_modal_and_chatbot = """    <!-- Unified Login Modal -->
    <div id="login-modal" style="display: none; position: fixed; top: 0; left: 0; width: 100%; height: 100%; background: rgba(30, 41, 59, 0.6); z-index: 1000; align-items: center; justify-content: center; backdrop-filter: blur(5px);">
        <div style="background: var(--white); border-radius: 16px; padding: 30px; width: 400px; max-width: 95%; box-shadow: 0 20px 40px rgba(0,0,0,0.1); border: var(--glass-border); position: relative;">
            <button onclick="closeLoginModal()" style="position: absolute; top: 15px; right: 15px; background: none; border: none; font-size: 24px; cursor: pointer; color: var(--text-muted);">&times;</button>
            <h3 style="margin-bottom: 20px; text-align: center; color: var(--text-dark); font-weight: 700; font-size: 20px;">Đăng Nhập Hệ Thống</h3>
            <form id="unified-login-form">
                <div class="form-group" style="margin-bottom: 15px;">
                    <label style="display: block; margin-bottom: 5px; font-size: 13px; font-weight: 600;">Tên đăng nhập</label>
                    <input type="text" id="login-username" class="form-control" placeholder="Ví dụ: teacher1, admin1" required style="width:100%; padding:10px; border-radius:8px; border: var(--glass-border);">
                </div>
                <div class="form-group" style="margin-bottom: 20px;">
                    <label style="display: block; margin-bottom: 5px; font-size: 13px; font-weight: 600;">Mật khẩu</label>
                    <input type="password" id="login-password" class="form-control" placeholder="Mật khẩu" required style="width:100%; padding:10px; border-radius:8px; border: var(--glass-border);">
                </div>
                <button type="submit" class="btn-submit" style="width: 100%; padding: 12px; background: var(--accent-green); color: white; border: none; border-radius: 8px; font-weight: 600; cursor: pointer;">Đăng Nhập</button>
            </form>
        </div>
    </div>

    <!-- Chatbot Widget góc màn hình -->"""

if old_chatbot in html:
    html = html.replace(old_chatbot, new_modal_and_chatbot)
    print("SUCCESS: Unified login modal added successfully!")
else:
    print("Warning: old chatbot comment not found directly.")

with open("playground_updated.html", "w", encoding="utf-8") as f:
    f.write(html)

print("HTML upgrade completed on playground_updated.html!")
