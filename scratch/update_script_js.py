import re

print("Running Javascript upgrade on playground_updated.html...")

with open("playground_updated.html", "r", encoding="utf-8") as f:
    html = f.read()

# Locate the <script> block and replace it
script_start = html.find("<script>")
script_end = html.rfind("</script>")

if script_start != -1 and script_end != -1:
    new_script_content = """<script>
        // Global State Variables
        let currentGrade = 5;
        let currentModel = 'online';
        let currentClassType = 'basic';
        let currentPerformance = 'good';
        
        // Unified Authentication State Variables
        let authToken = localStorage.getItem('auth_token') || '';
        let authUsername = localStorage.getItem('auth_username') || '';
        let authRole = localStorage.getItem('auth_role') || '';
        let authName = localStorage.getItem('auth_name') || '';

        // Update navigation and login status on UI
        function updateNavbarAuthUI() {
            const liSales = document.getElementById('li-sales-dashboard');
            const liAdmin = document.getElementById('li-admin-dashboard');
            const liLogin = document.getElementById('li-login');
            const liUser = document.getElementById('li-user-info');
            const badge = document.getElementById('nav-username-badge');

            if (authToken) {
                liLogin.style.display = 'none';
                liUser.style.display = 'flex';
                badge.innerText = authName + ' (' + (authRole === 'admin' ? 'Admin' : 'Sales') + ')';
                
                if (authRole === 'admin') {
                    liAdmin.style.display = 'inline-block';
                    liSales.style.display = 'none';
                } else if (authRole === 'teacher') {
                    liSales.style.display = 'inline-block';
                    liAdmin.style.display = 'none';
                }
            } else {
                liLogin.style.display = 'inline-block';
                liUser.style.display = 'none';
                liSales.style.display = 'none';
                liAdmin.style.display = 'none';
            }
        }

        // SPA Navigation tabs switcher
        function switchTab(tabId) {
            document.querySelectorAll('.tab-panel').forEach(panel => {
                panel.classList.remove('active');
            });
            document.querySelectorAll('.nav-links button').forEach(btn => {
                btn.classList.remove('active');
            });

            document.getElementById(tabId).classList.add('active');
            
            // Map button active
            if (tabId === 'tab-homepage') {
                document.getElementById('btn-homepage').classList.add('active');
                loadHomepageClasses();
                selectGrade(currentGrade);
            } else if (tabId === 'tab-about') {
                document.getElementById('btn-about').classList.add('active');
                loadAboutTeachers();
            } else if (tabId === 'tab-honor') {
                document.getElementById('btn-honor').classList.add('active');
                loadHonorStudents();
            } else if (tabId === 'tab-sales') {
                document.getElementById('btn-sales').classList.add('active');
                checkSalesAuth();
            } else if (tabId === 'tab-admin') {
                document.getElementById('btn-admin').classList.add('active');
                checkAdminAuth();
            }
        }

        // Open and close login modal
        function openLoginModal() {
            document.getElementById('login-modal').style.display = 'flex';
        }

        function closeLoginModal() {
            document.getElementById('login-modal').style.display = 'none';
            document.getElementById('unified-login-form').reset();
        }

        // Unified login form submission
        document.getElementById('unified-login-form').onsubmit = function(e) {
            e.preventDefault();
            const user = document.getElementById('login-username').value.trim();
            const pass = document.getElementById('login-password').value;

            fetch('/api/v1/auth/login', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ username: user, password: pass })
            })
            .then(res => res.json())
            .then(result => {
                if (result.success) {
                    authToken = result.data.token;
                    authUsername = result.data.username;
                    authRole = result.data.role;
                    authName = result.data.name;

                    localStorage.setItem('auth_token', authToken);
                    localStorage.setItem('auth_username', authUsername);
                    localStorage.setItem('auth_role', authRole);
                    localStorage.setItem('auth_name', authName);

                    updateNavbarAuthUI();
                    closeLoginModal();

                    // Navigate to appropriate tab based on role
                    if (authRole === 'admin') {
                        switchTab('tab-admin');
                    } else if (authRole === 'teacher') {
                        switchTab('tab-sales');
                    }
                } else {
                    alert(result.error);
                }
            });
        };

        // Unified Logout
        function handleLogout() {
            fetch('/api/v1/auth/logout', {
                method: 'POST',
                headers: { 'Authorization': 'Bearer ' + authToken }
            }).finally(() => {
                authToken = '';
                authUsername = '';
                authRole = '';
                authName = '';

                localStorage.removeItem('auth_token');
                localStorage.removeItem('auth_username');
                localStorage.removeItem('auth_role');
                localStorage.removeItem('auth_name');

                updateNavbarAuthUI();
                switchTab('tab-homepage');
            });
        }

        // Dashboard authentications checks
        function checkSalesAuth() {
            if (authToken && (authRole === 'teacher' || authRole === 'admin')) {
                document.getElementById('sales-login-section').style.display = 'none';
                document.getElementById('sales-monitor').style.display = 'block';
                document.getElementById('sales-active-user').innerText = '✓ Sales: ' + authName;
                loadLeadsSalesPortal();
                connectSSEStream();
            } else {
                document.getElementById('sales-login-section').style.display = 'block';
                document.getElementById('sales-monitor').style.display = 'none';
            }
        }

        function checkAdminAuth() {
            if (authToken && authRole === 'admin') {
                document.getElementById('admin-login-section').style.display = 'none';
                document.getElementById('admin-monitor').style.display = 'block';
                loadAdminDashboardData();
            } else {
                document.getElementById('admin-login-section').style.display = 'block';
                document.getElementById('admin-monitor').style.display = 'none';
            }
        }

        function logoutSession(portal) {
            handleLogout();
        }

        // Selection helpers for home page registration
        function selectGrade(grade) {
            currentGrade = grade;
            document.querySelectorAll('#grade-selector .select-btn').forEach((btn, idx) => {
                btn.classList.toggle('active', (idx + 1) === grade);
            });
            
            const hqBtn = document.getElementById('type-hq');
            if (grade === 5) {
                hqBtn.style.display = 'inline-block';
            } else {
                hqBtn.style.display = 'none';
                if (currentClassType === 'high_quality') {
                    selectType('basic');
                }
            }
            checkClassTypeConstraints();
        }

        function selectModel(model) {
            currentModel = model;
            document.getElementById('model-online').classList.toggle('active', model === 'online');
            document.getElementById('model-offline').classList.toggle('active', model === 'offline');
        }

        function selectType(type) {
            currentClassType = type;
            document.getElementById('type-basic').classList.toggle('active', type === 'basic');
            document.getElementById('type-advanced').classList.toggle('active', type === 'advanced');
            document.getElementById('type-hq').classList.toggle('active', type === 'high_quality');
            checkClassTypeConstraints();
        }

        function checkClassTypeConstraints() {
            if (currentClassType === 'high_quality' && currentGrade !== 5) {
                selectType('basic');
                alert('Chú ý: Lớp Chất Lượng Cao (CLC Ôn Luyện Chuyên) chỉ dành riêng cho khối Lớp 5. Hệ thống đã tự động đưa về lớp Cơ Bản.');
            }
        }

        function selectPerformance(perf) {
            currentPerformance = perf;
            document.getElementById('perf-average').classList.toggle('active', perf === 'average');
            document.getElementById('perf-good').classList.toggle('active', perf === 'good');
            document.getElementById('perf-excellent').classList.toggle('active', perf === 'excellent');
        }

        // Load dynamic home page classes
        function loadHomepageClasses() {
            fetch('/api/v1/homepage/classes')
                .then(res => res.json())
                .then(result => {
                    if (result.success) {
                        const grid = document.getElementById('homepage-classes-grid');
                        grid.innerHTML = '';
                        result.data.forEach(c => {
                            const isPop = c.is_popular ? '<div class="popular-badge">Khuyên dùng</div>' : '';
                            const popClass = c.is_popular ? 'popular' : '';
                            
                            grid.innerHTML += '<div class="class-card ' + popClass + '">' + 
                                isPop +
                                '<div class="class-header">' +
                                    '<h3>' + c.name + '</h3>' +
                                    '<p class="class-desc">' + c.desc + '</p>' +
                                '</div>' +
                                '<div class="class-footer">' +
                                    '<div class="price-label">Học phí niêm yết:</div>' +
                                    '<div class="price-value">' + c.price + '</div>' +
                                '</div>' +
                            '</div>';
                        });
                    }
                });
        }

        // Load About page teachers list
        function loadAboutTeachers() {
            fetch('/api/v1/about/teachers')
                .then(res => res.json())
                .then(result => {
                    if (result.success) {
                        const cotra = result.data.find(t => t.id === 1);
                        if (cotra) {
                            document.getElementById('about-cotra-img').src = cotra.avatar || '/assets/cotra.jpg';
                            document.getElementById('about-cotra-role').innerText = cotra.role;
                            document.getElementById('about-cotra-bio').innerText = cotra.bio;
                        }

                        const colleagues = result.data.filter(t => t.id !== 1);
                        const grid = document.getElementById('about-teachers-grid');
                        grid.innerHTML = '';
                        colleagues.forEach(t => {
                            const avatar = t.avatar || 'data:image/svg+xml;utf8,<svg xmlns="http://www.w3.org/2000/svg" width="90" height="90" viewBox="0 0 100 100"><circle cx="50" cy="50" r="45" fill="%23E2E8F0"/><text x="50" y="55" font-family="sans-serif" font-size="24" text-anchor="middle" fill="%2364748B">' + t.name.charAt(0) + '</text></svg>';
                            grid.innerHTML += '<div class="teacher-card">' +
                                '<div class="teacher-avatar-circle">' +
                                    '<img src="' + avatar + '" alt="' + t.name + '">' +
                                '</div>' +
                                '<h3>' + t.name + '</h3>' +
                                '<h5>' + t.role + '</h5>' +
                                '<div class="teacher-edu">' + t.education + '</div>' +
                                '<p class="teacher-bio">' + t.bio + '</p>' +
                            '</div>';
                        });
                    }
                });
        }

        // Load Honor page students list
        function loadHonorStudents() {
            fetch('/api/v1/honor/students')
                .then(res => res.json())
                .then(result => {
                    if (result.success) {
                        const grid = document.getElementById('honor-students-grid');
                        grid.innerHTML = '';
                        result.data.forEach(s => {
                            const avatar = s.avatar || 'data:image/svg+xml;utf8,<svg xmlns="http://www.w3.org/2000/svg" width="100" height="100" viewBox="0 0 100 100"><circle cx="50" cy="50" r="45" fill="%23E2E8F0"/><text x="50" y="55" font-family="sans-serif" font-size="24" text-anchor="middle" fill="%2364748B">' + s.name.charAt(0) + '</text></svg>';
                            grid.innerHTML += '<div class="student-card">' +
                                '<div class="student-badge-crown">Tuyên Dương ✓</div>' +
                                '<div class="student-avatar">' +
                                    '<img src="' + avatar + '">' +
                                '</div>' +
                                '<h3>' + s.name + '</h3>' +
                                '<div class="student-class">' + s.class + ' | ' + s.year + '</div>' +
                                '<div class="student-achievement">' + s.achievement + '</div>' +
                            '</div>';
                        });
                    }
                });
        }

        // Submit study consultation registration form
        document.getElementById('lead-register-form').onsubmit = function(e) {
            e.preventDefault();
            const parentName = document.getElementById('reg-parent').value;
            const phone = document.getElementById('reg-phone').value;
            const student = document.getElementById('reg-student').value;

            fetch('/api/v1/registrations', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    parent_name: parentName,
                    phone_number: phone,
                    student_name: student,
                    grade: currentGrade,
                    learning_model: currentModel,
                    class_type: currentClassType,
                    academic_performance: currentPerformance
                })
            })
            .then(res => res.json())
            .then(result => {
                if (result.success) {
                    alert('Đăng ký tư vấn học thử thành công! Cô Trà và đội ngũ tuyển sinh sẽ sớm liên hệ trực tiếp với Phụ huynh.');
                    document.getElementById('lead-register-form').reset();
                    selectGrade(5);
                    selectModel('online');
                    selectType('basic');
                    selectPerformance('good');
                } else {
                    alert('Lỗi đăng ký: ' + result.error);
                }
            });
        };

        // Load and render Leads on Sales Portal
        function loadLeadsSalesPortal() {
            fetch('/api/v1/registrations', {
                headers: { 'Authorization': 'Bearer ' + authToken }
            })
            .then(res => res.json())
            .then(result => {
                if (result.success) {
                    renderLeads(result.data);
                }
            });
        }

        function renderLeads(leads) {
            const container = document.getElementById('sales-leads-list');
            container.innerHTML = '';
            
            if (leads.length === 0) {
                container.innerHTML = '<div style="text-align: center; color: var(--text-muted); padding: 30px;">Hiện tại chưa có lượt đăng ký tư vấn học thử nào.</div>';
                return;
            }

            leads.sort((a,b) => new Date(b.created_at) - new Date(a.created_at));

            leads.forEach(l => {
                const isContacted = l.status === 'contacted';
                const contactedClass = isContacted ? 'contacted' : '';
                const toggleActiveClass = isContacted ? 'active' : '';
                const consultText = isContacted ? 'Đã Tư Vấn ✓' : 'Chờ Tư Vấn';
                const learningText = l.learning_model === 'online' ? '💻 Online' : '🏫 Tại Lớp';
                const classText = l.class_type === 'basic' ? 'Cơ Bản' : l.class_type === 'advanced' ? 'Nâng Cao' : 'CLC Luyện Thi';
                const perfText = l.academic_performance === 'excellent' ? '🏆 Giỏi/Xuất sắc' : l.academic_performance === 'good' ? '✨ Học lực Khá' : '📚 Học lực TB';
                
                // Audit tag if consulted
                const auditText = (isContacted && l.consulted_by) ? '<span style="background: rgba(16, 185, 129, 0.1); color: var(--accent-green); font-weight: 600; padding: 2px 8px; border-radius: 4px; font-size: 11px; margin-left: 8px;">Đã tư vấn bởi: ' + l.consulted_by + '</span>' : '';

                container.innerHTML += '<div class="lead-card ' + contactedClass + '">' +
                    '<div class="lead-info-left">' +
                        '<h4>' + l.student_name + ' (Học sinh) - Phụ huynh: ' + l.parent_name + '</h4>' +
                        '<div class="lead-meta">' +
                            '<span>📞 ' + l.phone_number + '</span>' +
                            '<span>Khối Lớp: ' + l.grade + '</span>' +
                            '<span>' + learningText + '</span>' +
                            '<span>' + classText + '</span>' +
                            '<span>' + perfText + '</span>' +
                            '<span>Thời gian: ' + new Date(l.created_at).toLocaleTimeString() + '</span>' +
                            auditText +
                        '</div>' +
                    '</div>' +
                    '<div class="lead-actions-right">' +
                        '<div class="ios-toggle-container ' + toggleActiveClass + '" onclick="toggleLeadConsult(' + l.id + ')">' +
                            '<span>' + consultText + '</span>' +
                            '<div class="ios-switch"></div>' +
                        '</div>' +
                    '</div>' +
                '</div>';
            });
        }

        function toggleLeadConsult(leadId) {
            fetch('/api/v1/registrations/toggle-consulted', {
                method: 'POST',
                headers: { 
                    'Content-Type': 'application/json',
                    'Authorization': 'Bearer ' + authToken
                },
                body: JSON.stringify({ id: leadId })
            })
            .then(res => res.json())
            .then(result => {
                if (result.success) {
                    loadLeadsSalesPortal();
                } else {
                    alert(result.error);
                }
            });
        }

        // Live SSE stream connection
        let sseConnected = false;
        function connectSSEStream() {
            if (sseConnected) return;

            const eventSource = new EventSource('/api/v1/registrations/stream?token=' + authToken);
            
            eventSource.addEventListener('registration_created', function(e) {
                playNotificationBeep();
                loadLeadsSalesPortal();
            });

            eventSource.addEventListener('registration_updated', function(e) {
                loadLeadsSalesPortal();
            });

            eventSource.addEventListener('connected', function(e) {
                sseConnected = true;
                document.getElementById('sales-active-user').innerText = '✓ Sales trực tuyến: ' + authName + ' (Đang kết nối)';
            });

            eventSource.onerror = function() {
                sseConnected = false;
                document.getElementById('sales-active-user').innerText = '⚠ Mất kết nối SSE Stream';
            };
        }

        function playNotificationBeep() {
            try {
                const context = new (window.AudioContext || window.webkitAudioContext)();
                const osc = context.createOscillator();
                osc.type = 'sine';
                osc.frequency.setValueAtTime(800, context.currentTime);
                osc.connect(context.destination);
                osc.start();
                osc.stop(context.currentTime + 0.15);
            } catch (err) {}
        }

        // ----------------------------------------------------
        // ADMIN DASHBOARD DATA & CRUD LOADING
        // ----------------------------------------------------

        function loadAdminDashboardData() {
            // 1. Tải báo cáo widgets
            fetch('/api/v1/admin/reports', {
                headers: { 'Authorization': 'Bearer ' + authToken }
            })
            .then(res => res.json())
            .then(result => {
                if (result.success) {
                    const r = result.data;
                    document.getElementById('wd-total').innerText = r.total_leads;
                    document.getElementById('wd-contacted').innerText = r.contacted_leads;
                    document.getElementById('wd-pending').innerText = r.pending_leads;
                    document.getElementById('wd-ratio').innerText = Math.round(r.online_percentage) + '% / ' + Math.round(r.offline_percentage) + '%';
                }
            });

            // 2. Tải cấu hình mức giá lớp học động
            loadAdminClassesPrices();

            // 3. Tải danh sách giảng viên
            loadAdminTeachers();

            // 4. Tải danh sách vinh danh học sinh
            loadAdminHonoredStudents();

            // 5. Tải danh sách tài khoản sales
            loadAdminUsers();
        }

        // Load all classes and render price inputs under Admin Panel
        function loadAdminClassesPrices() {
            fetch('/api/v1/homepage/classes')
                .then(res => res.json())
                .then(result => {
                    if (result.success) {
                        const container = document.getElementById('admin-classes-price-list');
                        container.innerHTML = '';
                        result.data.forEach(c => {
                            container.innerHTML += '<div class="class-price-row" style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px;">' +
                                '<span style="font-weight: 500; font-size: 13px;">' + c.name + ':</span>' +
                                '<input type="text" data-class-id="' + c.id + '" class="form-control class-price-input" value="' + c.price + '" style="width: 160px; padding: 6px 10px; font-size: 12px;" required>' +
                            '</div>';
                        });
                    }
                });
        }

        // Save class prices form
        document.getElementById('admin-form-prices').onsubmit = function(e) {
            e.preventDefault();
            const inputs = document.querySelectorAll('.class-price-input');
            const promises = [];
            
            inputs.forEach(input => {
                const id = input.getAttribute('data-class-id');
                const val = input.value;
                promises.push(
                    fetch('/api/v1/admin/classes/update', {
                        method: 'POST',
                        headers: { 
                            'Content-Type': 'application/json',
                            'Authorization': 'Bearer ' + authToken
                        },
                        body: JSON.stringify({ id: id, price: val })
                    })
                );
            });

            Promise.all(promises).then(() => {
                alert('Tùy chỉnh học phí hoàn thành! Mức giá mới đã được cập nhật ngay lên Trang Chủ.');
                loadHomepageClasses();
                loadAdminDashboardData();
            });
        };

        // Load and render active users in Admin Panel
        function loadAdminUsers() {
            fetch('/api/v1/admin/users', {
                headers: { 'Authorization': 'Bearer ' + authToken }
            })
            .then(res => res.json())
            .then(result => {
                if (result.success) {
                    const tbody = document.getElementById('admin-table-users');
                    tbody.innerHTML = '';
                    result.data.forEach(u => {
                        const isSuperAdmin = u.username === 'admin1';
                        const deleteBtn = isSuperAdmin ? '' : '<button class="btn-delete-small" onclick="deleteSalesUser(\'' + u.username + '\')" style="background: #EF4444; color: white; padding: 4px 8px; font-size: 11px;">Xóa tài khoản</button>';
                        const roleBadge = u.role === 'admin' ? '<span class="user-badge" style="background: rgba(255, 184, 0, 0.1); color: #B45309; padding: 2px 6px; border-radius: 4px; font-size: 10px;">Admin</span>' : '<span class="user-badge" style="background: rgba(16, 185, 129, 0.1); color: var(--accent-green); padding: 2px 6px; border-radius: 4px; font-size: 10px;">Sales / Teacher</span>';
                        
                        tbody.innerHTML += '<tr>' +
                            '<td style="font-weight: 700;">' + u.username + '</td>' +
                            '<td>' + (u.name || 'Chưa cập nhật') + '</td>' +
                            '<td>' + roleBadge + '</td>' +
                            '<td>' + deleteBtn + '</td>' +
                        '</tr>';
                    });
                }
            });
        }

        // Create new sales account
        document.getElementById('admin-form-create-sales-user').onsubmit = function(e) {
            e.preventDefault();
            const u = document.getElementById('new-sales-username').value.trim();
            const p = document.getElementById('new-sales-password').value;
            const n = document.getElementById('new-sales-fullname').value.trim();

            fetch('/api/v1/admin/users/create', {
                method: 'POST',
                headers: { 
                    'Content-Type': 'application/json',
                    'Authorization': 'Bearer ' + authToken
                },
                body: JSON.stringify({ username: u, password: p, name: n })
            })
            .then(res => res.json())
            .then(result => {
                if (result.success) {
                    alert(result.message);
                    document.getElementById('admin-form-create-sales-user').reset();
                    loadAdminUsers();
                } else {
                    alert('Lỗi tạo tài khoản: ' + result.error);
                }
            });
        };

        // Delete sales account
        function deleteSalesUser(username) {
            if (!confirm('Bạn có chắc chắn muốn xóa vĩnh viễn tài khoản sales \"' + username + '\" này không?')) return;

            fetch('/api/v1/admin/users/delete', {
                method: 'POST',
                headers: { 
                    'Content-Type': 'application/json',
                    'Authorization': 'Bearer ' + authToken
                },
                body: JSON.stringify({ username: username })
            })
            .then(res => res.json())
            .then(result => {
                if (result.success) {
                    alert(result.message);
                    loadAdminUsers();
                } else {
                    alert(result.error);
                }
            });
        }

        // Load and render teachers list under Admin Panel
        function loadAdminTeachers() {
            fetch('/api/v1/about/teachers')
                .then(res => res.json())
                .then(result => {
                    if (result.success) {
                        const cotra = result.data.find(t => t.id === 1);
                        if (cotra) {
                            document.getElementById('admin-cotra-name').value = cotra.name;
                            document.getElementById('admin-cotra-role').value = cotra.role;
                            document.getElementById('admin-cotra-bio').value = cotra.bio;
                            document.getElementById('admin-cotra-avatar-url').value = cotra.avatar;
                            document.getElementById('admin-cotra-avatar-preview').src = cotra.avatar || '/assets/cotra.jpg';
                        }

                        const tbody = document.getElementById('admin-table-teachers');
                        tbody.innerHTML = '';
                        result.data.forEach(t => {
                            const img = t.avatar || '/assets/cotra.jpg';
                            const isFounder = t.id === 1;
                            const deleteBtn = isFounder ? '' : '<button class="btn-delete-small" onclick="deleteTeacher(' + t.id + ')" style="background: #EF4444; color: white; padding: 4px 8px; font-size: 11px;">Xóa</button>';
                            const editBtn = '<button class="btn-action-small" onclick="editTeacher(' + JSON.stringify(t).replace(/"/g, '&quot;') + ')" style="background: var(--accent-amber); color: var(--text-dark); margin-right: 5px; padding: 4px 8px; font-size: 11px;">Sửa</button>';
                            
                            tbody.innerHTML += '<tr>' +
                                '<td><img src=\"' + img + '\" style=\"width: 30px; height: 30px; border-radius: 50%; object-fit: cover;\"></td>' +
                                '<td style=\"font-weight: 700;\">' + t.name + '</td>' +
                                '<td>' + t.role + '</td>' +
                                '<td>' + t.education + '</td>' +
                                '<td>' + t.order + '</td>' +
                                '<td>' + editBtn + deleteBtn + '</td>' +
                            '</tr>';
                        });
                    }
                });
        }

        function editTeacher(t) {
            document.getElementById('teacher-form-title').innerText = 'Chỉnh Sửa Giảng Viên: ' + t.name;
            document.getElementById('teacher-edit-id').value = t.id;
            document.getElementById('teacher-name').value = t.name;
            document.getElementById('teacher-role').value = t.role;
            document.getElementById('teacher-education').value = t.education;
            document.getElementById('teacher-avatar-url').value = t.avatar;
            document.getElementById('teacher-order').value = t.order;
            document.getElementById('teacher-bio').value = t.bio;
            document.getElementById('btn-cancel-teacher-edit').style.display = 'inline-block';
        }

        function resetTeacherForm() {
            document.getElementById('teacher-form-title').innerText = 'Thêm Mới Giảng Viên';
            document.getElementById('teacher-edit-id').value = '';
            document.getElementById('admin-form-manage-teacher').reset();
            document.getElementById('teacher-avatar-url').value = '';
            document.getElementById('btn-cancel-teacher-edit').style.display = 'none';
        }

        // Upload avatar teacher
        document.getElementById('teacher-avatar-file').onchange = function(e) {
            const file = e.target.files[0];
            if (!file) return;

            const formData = new FormData();
            formData.append('file', file);

            fetch('/api/v1/admin/upload', {
                method: 'POST',
                headers: { 'Authorization': 'Bearer ' + authToken },
                body: formData
            })
            .then(res => res.json())
            .then(result => {
                if (result.success) {
                    document.getElementById('teacher-avatar-url').value = result.url;
                    alert('Tải ảnh đại diện giáo viên thành công!');
                } else {
                    alert(result.error);
                }
            });
        };

        // Manage teacher save form
        document.getElementById('admin-form-manage-teacher').onsubmit = function(e) {
            e.preventDefault();
            const editId = document.getElementById('teacher-edit-id').value;
            const name = document.getElementById('teacher-name').value.trim();
            const role = document.getElementById('teacher-role').value.trim();
            const education = document.getElementById('teacher-education').value.trim();
            const order = parseInt(document.getElementById('teacher-order').value);
            const avatar = document.getElementById('teacher-avatar-url').value;
            const bio = document.getElementById('teacher-bio').value.trim();

            if (editId) {
                fetch('/api/v1/admin/teachers/update', {
                    method: 'POST',
                    headers: { 
                        'Content-Type': 'application/json',
                        'Authorization': 'Bearer ' + authToken
                    },
                    body: JSON.stringify({ id: parseInt(editId), name: name, role: role, education: education, order: order, avatar: avatar, bio: bio })
                })
                .then(res => res.json())
                .then(result => {
                    if (result.success) {
                        alert('Cập nhật giảng viên thành công!');
                        resetTeacherForm();
                        loadAdminTeachers();
                        loadAboutTeachers();
                    } else {
                        alert(result.error);
                    }
                });
            } else {
                fetch('/api/v1/admin/teachers/manage', {
                    method: 'POST',
                    headers: { 
                        'Content-Type': 'application/json',
                        'Authorization': 'Bearer ' + authToken
                    },
                    body: JSON.stringify({
                        action: 'add',
                        teacher_data: { name: name, role: role, education: education, order: order, avatar: avatar, bio: bio }
                    })
                })
                .then(res => res.json())
                .then(result => {
                    if (result.success) {
                        alert(result.message);
                        resetTeacherForm();
                        loadAdminTeachers();
                        loadAboutTeachers();
                    } else {
                        alert(result.error);
                    }
                });
            }
        };

        function deleteTeacher(teacherId) {
            if (!confirm('Bạn có chắc chắn muốn xóa giảng viên này khỏi danh sách giới thiệu không?')) return;

            fetch('/api/v1/admin/teachers/manage', {
                method: 'POST',
                headers: { 
                    'Content-Type': 'application/json',
                    'Authorization': 'Bearer ' + authToken
                },
                body: JSON.stringify({ action: 'delete', id: teacherId })
            })
            .then(res => res.json())
            .then(result => {
                if (result.success) {
                    alert(result.message);
                    loadAdminTeachers();
                    loadAboutTeachers();
                } else {
                    alert(result.error);
                }
            });
        }

        // Upload avatar Co Tra directly
        document.getElementById('admin-cotra-avatar-file').onchange = function(e) {
            const file = e.target.files[0];
            if (!file) return;

            const formData = new FormData();
            formData.append('file', file);

            fetch('/api/v1/admin/upload', {
                method: 'POST',
                headers: { 'Authorization': 'Bearer ' + authToken },
                body: formData
            })
            .then(res => res.json())
            .then(result => {
                if (result.success) {
                    document.getElementById('admin-cotra-avatar-url').value = result.url;
                    document.getElementById('admin-cotra-avatar-preview').src = result.url;
                    alert('Tải ảnh đại diện Cô Trà hoàn tất!');
                } else {
                    alert('Lỗi upload ảnh: ' + result.error);
                }
            });
        };

        // Update Co Tra bio
        document.getElementById('admin-form-cotra-bio').onsubmit = function(e) {
            e.preventDefault();
            const name = document.getElementById('admin-cotra-name').value;
            const role = document.getElementById('admin-cotra-role').value;
            const bio = document.getElementById('admin-cotra-bio').value;
            const avatar = document.getElementById('admin-cotra-avatar-url').value;

            fetch('/api/v1/admin/teachers/update', {
                method: 'POST',
                headers: { 
                    'Content-Type': 'application/json',
                    'Authorization': 'Bearer ' + authToken
                },
                body: JSON.stringify({ id: 1, name: name, role: role, bio: bio, avatar: avatar })
            })
            .then(res => res.json())
            .then(result => {
                if (result.success) {
                    alert('Cập nhật tiểu sử Cô Trà hoàn tất!');
                    loadAdminDashboardData();
                    loadAboutTeachers();
                } else {
                    alert(result.error);
                }
            });
        };

        // Load admin honored students list
        function loadAdminHonoredStudents() {
            fetch('/api/v1/honor/students')
                .then(res => res.json())
                .then(result => {
                    if (result.success) {
                        const tbody = document.getElementById('admin-table-students');
                        tbody.innerHTML = '';
                        result.data.forEach(s => {
                            const img = s.avatar || 'data:image/svg+xml;utf8,<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"30\" height=\"30\" viewBox=\"0 0 100 100\"><circle cx=\"50\" cy=\"50\" r=\"45\" fill=\"%23E2E8F0\"/></svg>';
                            tbody.innerHTML += 'tr' +
                                'td' + 'img src=\"' + img + '\" style=\"width: 30px; height: 30px; border-radius: 50%; object-fit: cover;\"' + '/td' +
                                'td style=\"font-weight: 700;\"' + s.name + '/td' +
                                'td' + s.class + '/td' +
                                'td' + s.year + '/td' +
                                'td' + s.achievement + '/td' +
                                'td' + 'button class=\"btn-delete-small\" onclick=\"deleteHonoredStudent(' + s.id + ')\"' + 'Xóa' + '/button' + '/td' +
                            '/tr';
                        });
                        tbody.innerHTML = tbody.innerHTML
                            .replace(/tr/g, '<tr>').replace(/\/tr/g, '</tr>')
                            .replace(/td/g, '<td>').replace(/\/td/g, '</td>')
                            .replace(/<td><td>/g, '<td>').replace(/<\/td><\/td>/g, '</td>');
                    }
                });
        }

        // Upload honored student photo
        document.getElementById('new-stu-avatar-file').onchange = function(e) {
            const file = e.target.files[0];
            if (!file) return;

            const formData = new FormData();
            formData.append('file', file);

            fetch('/api/v1/admin/upload', {
                method: 'POST',
                headers: { 'Authorization': 'Bearer ' + authToken },
                body: formData
            })
            .then(res => res.json())
            .then(result => {
                if (result.success) {
                    document.getElementById('new-stu-avatar-url').value = result.url;
                    alert('Tải ảnh chân dung học sinh thành công!');
                } else {
                    alert(result.error);
                }
            });
        };

        // Add honored student
        document.getElementById('admin-form-add-student').onsubmit = function(e) {
            e.preventDefault();
            const name = document.getElementById('new-stu-name').value;
            const sClass = document.getElementById('new-stu-class').value;
            const year = document.getElementById('new-stu-year').value;
            const achievement = document.getElementById('new-stu-achievement').value;
            const avatar = document.getElementById('new-stu-avatar-url').value;

            fetch('/api/v1/admin/students/manage', {
                method: 'POST',
                headers: { 
                    'Content-Type': 'application/json',
                    'Authorization': 'Bearer ' + authToken
                },
                body: JSON.stringify({
                    action: 'add',
                    student_data: { name: name, class: sClass, year: year, achievement: achievement, avatar: avatar }
                })
            })
            .then(res => res.json())
            .then(result => {
                if (result.success) {
                    alert(result.message);
                    document.getElementById('admin-form-add-student').reset();
                    document.getElementById('new-stu-avatar-url').value = '';
                    loadAdminDashboardData();
                    loadHonorStudents();
                } else {
                    alert(result.error);
                }
            });
        };

        // Delete honored student
        function deleteHonoredStudent(stuId) {
            if (!confirm('Bạn có chắc chắn muốn xóa học sinh này khỏi Bảng Vàng Vinh Danh hay không?')) return;

            fetch('/api/v1/admin/students/manage', {
                method: 'POST',
                headers: { 
                    'Content-Type': 'application/json',
                    'Authorization': 'Bearer ' + authToken
                },
                body: JSON.stringify({ action: 'delete', id: stuId })
            })
            .then(res => res.json())
            .then(result => {
                if (result.success) {
                    alert(result.message);
                    loadAdminDashboardData();
                    loadHonorStudents();
                } else {
                    alert(result.error);
                }
            });
        }

        // Chatbot Widget Simulation
        function toggleChat() {
            document.getElementById('chat-win').classList.toggle('active');
        }

        function handleChatKey(e) {
            if (e.key === 'Enter') {
                sendChatMessage();
            }
        }

        function sendChatMessage() {
            const input = document.getElementById('chat-txt-input');
            const txt = input.value.trim();
            if (!txt) return;

            const list = document.getElementById('chat-messages-list');
            list.innerHTML += '<div class="chat-msg user">' + txt + '</div>';
            input.value = '';
            list.scrollTop = list.scrollHeight;

            fetch('/api/v1/chatbot/webhook', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    message: { text: txt }
                })
            })
            .then(res => res.json())
            .then(result => {
                let reply = 'Cảm ơn phụ huynh đã liên hệ! Các giáo viên sẽ lập tức hỗ trợ trực tiếp giải đáp mọi thắc mắc qua Zalo số 098 345 93 93.';
                if (result && result.reply) {
                    reply = result.reply;
                }
                
                reply = reply.replace(/Gia sư AI 24\/7/g, 'Các giáo viên sát sao hỗ trợ nhiệt tình mọi lúc');

                setTimeout(() => {
                    list.innerHTML += '<div class="chat-msg bot">' + reply + '</div>';
                    list.scrollTop = list.scrollHeight;
                }, 800);
            });
        }

        // Initialize default tabs and configs on load
        window.onload = function() {
            updateNavbarAuthUI();
            switchTab('tab-homepage');
        };
    </script>"""

    # Do replacement of script block
    html = html[:script_start] + new_script_content + html[script_end + len("</script>"):]
    with open("playground_updated.html", "w", encoding="utf-8") as f:
        f.write(html)
    print("SUCCESS: Javascript updated successfully in playground_updated.html!")
else:
    print("ERROR: Could not find script start and end tags in playground_updated.html!")
