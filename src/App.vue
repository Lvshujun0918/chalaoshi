<script setup>
import { ref, onMounted, watch, computed } from 'vue'
import { useTeacherSearch } from './composables/useTeacherSearch.js'
import { useCourseSearch } from './composables/useCourseSearch.js'
import SearchBar from './components/SearchBar.vue'
import TeacherCard from './components/TeacherCard.vue'
import CourseCard from './components/CourseCard.vue'
import CommentModal from './components/CommentModal.vue'
import AppIcon from './components/AppIcon.vue'

// 构建时由 vite define 注入的 git commit 短哈希
const gitHash = __GIT_HASH__

// 数据版本（从 ver.json 读取）
const dataVersion = ref('加载中...')

const currentMode = ref('teachers') // 'teachers' | 'courses'

// ── 教师搜索 ──
const teacherSearch = useTeacherSearch()

// ── 课程搜索 ──
const courseSearch = useCourseSearch()

// 根据当前模式，动态选择用哪个 composable
const query = computed({
  get: () => currentMode.value === 'teachers' ? teacherSearch.query.value : courseSearch.query.value,
  set: (val) => {
    if (currentMode.value === 'teachers') teacherSearch.query.value = val
    else courseSearch.query.value = val
  },
})

const loading = computed(() =>
  currentMode.value === 'teachers' ? teacherSearch.loading.value : courseSearch.loading.value
)
const displayResults = computed(() =>
  currentMode.value === 'teachers' ? teacherSearch.displayResults.value : courseSearch.displayResults.value
)
const totalResults = computed(() =>
  currentMode.value === 'teachers' ? teacherSearch.totalResults.value : courseSearch.totalResults.value
)
const page = computed(() =>
  currentMode.value === 'teachers' ? teacherSearch.page.value : courseSearch.page.value
)
const totalPages = computed(() =>
  currentMode.value === 'teachers' ? teacherSearch.totalPages.value : courseSearch.totalPages.value
)

// ── 分页器：生成展示页码（含省略号） ──
const pageNumbers = computed(() => {
  const total = totalPages.value
  const current = page.value
  const pages = []
  if (total <= 7) {
    for (let i = 1; i <= total; i++) pages.push(i)
    return pages
  }
  // 首页
  pages.push(1)
  if (current > 3) pages.push('...')
  // 当前页附近
  const start = Math.max(2, current - 1)
  const end = Math.min(total - 1, current + 1)
  for (let i = start; i <= end; i++) pages.push(i)
  if (current < total - 2) pages.push('...')
  // 末页
  pages.push(total)
  return pages
})

// 跳转指定页
const jumpPageInput = ref('')
function goToPage() {
  const p = parseInt(jumpPageInput.value, 10)
  if (isNaN(p) || p < 1 || p > totalPages.value) return
  jumpPageInput.value = ''
  if (currentMode.value === 'teachers') teacherSearch.setPage(p)
  else courseSearch.setPage(p)
}

const currentSortBy = computed(() =>
  currentMode.value === 'teachers' ? teacherSearch.sortBy.value : courseSearch.sortBy.value
)
const currentSortOrder = computed(() =>
  currentMode.value === 'teachers' ? teacherSearch.sortOrder.value : courseSearch.sortOrder.value
)

const { totalTeachers, totalColleges, departments, currentDepartment, getCollegeName } = teacherSearch

const selectedTeacher = ref(null)

function switchMode(mode) {
  if (currentMode.value === mode) return
  currentMode.value = mode
  if (mode === 'teachers') {
    teacherSearch.search(1)
  } else {
    courseSearch.search(1)
  }
}

function openComments(teacher) {
  selectedTeacher.value = teacher
}

function closeComments() {
  selectedTeacher.value = null
}

function handleClear() {
  if (currentMode.value === 'teachers') {
    teacherSearch.clearQuery()
    teacherSearch.search(1)
  } else {
    courseSearch.clearQuery()
    courseSearch.search(1)
  }
}

// 防抖搜索：用户输入时自动触发
let debounceTimer = null
watch(query, () => {
  if (debounceTimer) clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    if (currentMode.value === 'teachers') {
      teacherSearch.search(1)
    } else {
      courseSearch.search(1)
    }
  }, 350)
})

function handleDeptFilter(dept) {
  teacherSearch.filterByDepartment(dept)
}

// ── 课程点击 → 查找对应教师并打开评分详情 ──
async function handleCourseClick(course) {
  const teacherName = course.teacher_name
  if (!teacherName) return
  try {
    const res = await fetch(`/api/search?q=${encodeURIComponent(teacherName)}`)
    const teachers = await res.json()
    if (teachers && teachers.length > 0) {
      // 精确匹配教师姓名
      const match = teachers.find(t => t.name === teacherName) || teachers[0]
      selectedTeacher.value = match
    }
  } catch (e) {
    console.error('查找教师失败:', e)
  }
}

// ── 排序 ──
const teacherSortOptions = [
  { value: 'rating', label: '评分' },
  { value: 'hotness', label: '热度' },
  { value: 'rating_count', label: '评论数' },
  { value: 'name', label: '姓名' },
]

const courseSortOptions = [
  { value: 'gpa', label: 'GPA' },
  { value: 'count', label: '选课人数' },
  { value: 'course_name', label: '课程名' },
  { value: 'std_dev', label: '标准差' },
]

const currentSortOptions = computed(() =>
  currentMode.value === 'teachers' ? teacherSortOptions : courseSortOptions
)

function handleSortChange(sortValue) {
  const isTeacher = currentMode.value === 'teachers'
  const current = isTeacher ? teacherSearch : courseSearch
  if (currentSortBy.value === sortValue) {
    // 点击同一项 → 切换正序/倒序
    const newOrder = currentSortOrder.value === 'desc' ? 'asc' : 'desc'
    current.setSortOrder(newOrder)
  } else {
    // 切换排序字段 → 默认倒序
    current.setSort(sortValue)
  }
}

// 初始加载热门教师 + 数据版本
onMounted(async () => {
  teacherSearch.search(1)
  try {
    const res = await fetch('/api/version')
    const v = await res.json()
    dataVersion.value = v.release_date || v.version || '未知'
  } catch {
    dataVersion.value = '未知'
  }
})
</script>

<template>
  <div class="app-shell">
    <!-- Header -->
    <header class="app-header">
      <div class="header-content">
        <h1 class="app-title">
          <AppIcon name="book" :size="32" class="title-icon" />
          查老师
        </h1>
        <p class="app-subtitle">教师评分 &amp; 课程GPA查询系统</p>
        <div class="header-stats">
          <div class="stat">
            <AppIcon name="users" :size="18" class="stat-icon" />
            <span class="stat-num">{{ totalTeachers }}</span>
            <span class="stat-label">位教师</span>
          </div>
          <div class="stat-divider" />
          <div class="stat">
            <AppIcon name="building" :size="18" class="stat-icon" />
            <span class="stat-num">{{ totalColleges }}</span>
            <span class="stat-label">个学院</span>
          </div>
        </div>
      </div>
    </header>

    <!-- Mode Switcher -->
    <div class="mode-switcher">
      <button
        :class="['mode-btn', { active: currentMode === 'teachers' }]"
        @click="switchMode('teachers')"
      >
        <AppIcon name="users" :size="18" class="mode-icon" />
        查老师
      </button>
      <button
        :class="['mode-btn', { active: currentMode === 'courses' }]"
        @click="switchMode('courses')"
      >
        <AppIcon name="graduation" :size="18" class="mode-icon" />
        查课程
      </button>
    </div>

    <!-- Search -->
    <SearchBar
      v-model="query"
      :placeholder="currentMode === 'teachers' ? '输入教师姓名、拼音或缩写…' : '输入课程名称或授课教师…'"
      @clear="handleClear()"
    />

    <!-- Sort Controls -->
    <div class="sort-bar">
      <span class="sort-label">
        <AppIcon name="sort" :size="14" class="sort-label-icon" />
        排序：
      </span>
      <div class="sort-options">
        <button
          v-for="opt in currentSortOptions"
          :key="opt.value"
          :class="['sort-chip', { active: currentSortBy === opt.value }]"
          @click="handleSortChange(opt.value)"
        >
          {{ opt.label }}
          <span v-if="currentSortBy === opt.value" class="sort-arrow">
            {{ currentSortOrder === 'asc' ? '↑' : '↓' }}
          </span>
        </button>
      </div>
    </div>

    <!-- Department Filter (仅教师模式) -->
    <div v-if="currentMode === 'teachers' && departments.length > 0" class="dept-filter">
      <button
        :class="['dept-chip', { active: !currentDepartment }]"
        @click="handleDeptFilter('')"
      >全部</button>
      <button
        v-for="d in departments.slice(0, 20)"
        :key="d.name"
        :class="['dept-chip', { active: currentDepartment === d.name }]"
        @click="handleDeptFilter(d.name)"
        :title="d.name + ' (' + d.count + ')'"
      >{{ d.name }}<span class="dept-count">{{ d.count }}</span></button>
    </div>

    <!-- Main Content -->
    <main class="app-main">
      <!-- Loading -->
      <div v-if="loading" class="loading-state">搜索中...</div>

      <!-- Has Results: Summary + List -->
      <template v-else-if="displayResults.length > 0">
        <!-- Result Summary -->
        <div class="result-summary">
          <span>找到 <strong>{{ totalResults }}</strong> 条结果</span>
          <span v-if="page > 1">第 {{ page }}/{{ totalPages }} 页</span>
        </div>

        <!-- Teacher List -->
        <TransitionGroup
          v-if="currentMode === 'teachers'"
          name="card-list"
          tag="div"
          class="teacher-list"
        >
          <TeacherCard
            v-for="t in displayResults"
            :key="t.id"
            :teacher="t"
            :college-name="getCollegeName(t.department)"
            @click="openComments"
          />
        </TransitionGroup>

        <!-- Course List -->
        <TransitionGroup
          v-if="currentMode === 'courses'"
          name="card-list"
          tag="div"
          class="teacher-list"
        >
          <CourseCard
            v-for="c in displayResults"
            :key="c.id"
            :course="c"
            @click="handleCourseClick"
          />
        </TransitionGroup>
      </template>

      <!-- No Results (initial welcome state) -->
      <div v-else-if="!query && !currentDepartment" class="empty-state welcome">
        <AppIcon name="pointer" :size="48" class="empty-icon" />
        <p v-if="currentMode === 'teachers'">输入教师姓名、拼音或缩写开始查询</p>
        <p v-else>输入课程名称或授课教师开始查询</p>
        <p class="empty-hint" v-if="currentMode === 'teachers'">支持教师姓名、拼音全拼和拼音首字母搜索</p>
        <p class="empty-hint" v-else>支持按课程名称、授课教师搜索，可按GPA/选课人数排序</p>
      </div>

      <!-- No Results (search returned empty) -->
      <div v-else class="empty-state">
        <AppIcon name="search-lg" :size="48" class="empty-icon" />
        <p>没有找到匹配的结果</p>
        <p class="empty-hint">试试其他关键词吧</p>
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="pagination">
        <button :disabled="page <= 1" @click="currentMode === 'teachers' ? teacherSearch.setPage(page - 1) : courseSearch.setPage(page - 1)">‹ 上一页</button>
        <template v-for="p in pageNumbers" :key="p">
          <span v-if="p === '...'" class="ellipsis">…</span>
          <button
            v-else
            :class="{ current: p === page }"
            @click="currentMode === 'teachers' ? teacherSearch.setPage(p) : courseSearch.setPage(p)"
          >{{ p }}</button>
        </template>
        <button :disabled="page >= totalPages" @click="currentMode === 'teachers' ? teacherSearch.setPage(page + 1) : courseSearch.setPage(page + 1)">下一页 ›</button>
        <span class="jump-box">
          跳至<input
            v-model="jumpPageInput"
            type="number"
            :min="1"
            :max="totalPages"
            placeholder="页"
            class="jump-input"
            @keyup.enter="goToPage"
          />页
          <button class="jump-btn" @click="goToPage">GO</button>
        </span>
      </div>
    </main>

    <!-- Footer -->
    <footer class="app-footer">
      <p>
        Made by lsj with
        <AppIcon name="heart" :size="12" class="footer-heart" />
        · 数据来自査老师离线版，仅供参考
      </p>
      <p class="footer-version">
        数据版本 {{ dataVersion }} · 构建版本 {{ gitHash }}
      </p>
    </footer>

    <!-- Comment Modal -->
    <CommentModal
      v-if="selectedTeacher"
      :teacher="selectedTeacher"
      @close="closeComments"
    />
  </div>
</template>

<style scoped>
.app-shell {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

/* ── Header ── */
.app-header {
  background: linear-gradient(135deg, #5b5ef7 0%, #7c5cfc 40%, #a78bfa 100%);
  padding: 40px 24px 56px;
  text-align: center;
  color: #fff;
  position: relative;
  overflow: hidden;
}

.app-header::before {
  content: '';
  position: absolute;
  inset: 0;
  background-image: url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23ffffff' fill-opacity='0.06'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E");
  pointer-events: none;
}

.header-content {
  position: relative;
  z-index: 1;
}

.app-title {
  font-size: 30px;
  font-weight: 800;
  letter-spacing: -0.5px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
}

.title-icon {
  flex-shrink: 0;
  color: rgba(255,255,255,.9);
}

.app-subtitle {
  margin-top: 4px;
  font-size: 15px;
  opacity: 0.85;
  font-weight: 400;
}

.header-stats {
  margin-top: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
}

/* ── Mode Switcher ── */
.mode-switcher {
  display: flex;
  justify-content: center;
  gap: 8px;
  margin: -20px auto 0;
  padding: 0 20px;
  position: relative;
  z-index: 10;
  max-width: 600px;
}

.mode-btn {
  flex: 1;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 10px 20px;
  border-radius: var(--radius-md);
  border: 2px solid var(--color-border-light);
  background: var(--color-surface);
  cursor: pointer;
  font-size: 15px;
  font-weight: 600;
  color: var(--color-text-secondary);
  transition: all var(--transition-fast);
  font-family: inherit;
  box-shadow: var(--shadow-sm);
}

.mode-btn:hover {
  border-color: var(--color-primary);
  color: var(--color-primary);
  box-shadow: var(--shadow-md);
}

.mode-btn.active {
  background: var(--color-primary);
  color: #fff;
  border-color: var(--color-primary);
  box-shadow: 0 4px 14px rgba(91, 94, 247, 0.35);
}

.mode-icon {
  flex-shrink: 0;
}

/* ── Sort Bar ── */
.sort-bar {
  max-width: 600px;
  margin: 12px auto 0;
  padding: 0 20px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.sort-label {
  font-size: 13px;
  color: var(--color-text-muted);
  display: inline-flex;
  align-items: center;
  gap: 4px;
  flex-shrink: 0;
}

.sort-label-icon {
  flex-shrink: 0;
}

.sort-options {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.sort-chip {
  padding: 4px 12px;
  border-radius: 16px;
  border: 1px solid var(--color-border-light);
  background: var(--color-surface);
  cursor: pointer;
  font-size: 12px;
  color: var(--color-text-secondary);
  transition: all var(--transition-fast);
  white-space: nowrap;
  font-family: inherit;
}

.sort-chip:hover {
  border-color: var(--color-primary);
  color: var(--color-primary);
}

.sort-chip.active {
  background: var(--color-primary);
  color: #fff;
  border-color: var(--color-primary);
}

.sort-arrow {
  margin-left: 2px;
  font-size: 11px;
  font-weight: 700;
}

.stat {
  display: flex;
  align-items: center;
  gap: 6px;
}

.stat-icon {
  opacity: 0.8;
  flex-shrink: 0;
}

.stat-num {
  font-size: 24px;
  font-weight: 700;
}

.stat-label {
  font-size: 13px;
  opacity: 0.8;
}

.stat-divider {
  width: 1px;
  height: 24px;
  background: rgba(255,255,255,.3);
}

/* ── Main ── */
.app-main {
  flex: 1;
  max-width: 760px;
  width: 100%;
  margin: 0 auto;
  padding: 28px 20px 48px;
}

/* ── Result Summary ── */
.result-summary {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 14px;
  font-size: 14px;
  color: var(--color-text-secondary);
}

.result-summary strong {
  color: var(--color-text);
}

/* ── Department Filter ── */
.dept-filter {
  max-width: 760px;
  width: 100%;
  margin: 0 auto;
  padding: 20px 8px;
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.dept-chip {
  padding: 4px 12px;
  border-radius: 16px;
  border: 1px solid var(--color-border-light);
  background: var(--color-surface);
  cursor: pointer;
  font-size: 12px;
  color: var(--color-text-secondary);
  transition: all 0.2s;
  white-space: nowrap;
}

.dept-chip:hover {
  border-color: var(--color-primary);
  color: var(--color-primary);
}

.dept-chip.active {
  background: var(--color-primary);
  color: #fff;
  border-color: var(--color-primary);
}

.dept-count {
  margin-left: 4px;
  font-size: 10px;
  opacity: 0.7;
}

/* ── Loading ── */
.loading-state {
  text-align: center;
  padding: 40px;
  color: var(--color-text-secondary);
  font-size: 14px;
}

/* ── Pagination ── */
.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  margin-top: 24px;
  padding-top: 20px;
}

.pagination button {
  min-width: 32px;
  height: 32px;
  padding: 0 10px;
  border-radius: 8px;
  border: 1px solid var(--color-border-light);
  background: var(--color-surface);
  cursor: pointer;
  font-size: 13px;
  color: var(--color-text-secondary);
  transition: all 0.15s;
}

.pagination button:hover:not(:disabled) {
  border-color: var(--color-primary);
  color: var(--color-primary);
}

.pagination button.current {
  background: var(--color-primary);
  color: #fff;
  border-color: var(--color-primary);
}

.pagination button:disabled {
  opacity: 0.4;
  cursor: default;
}

.pagination .ellipsis {
  padding: 0 2px;
  color: var(--color-text-muted);
}

.pagination .jump-box {
  margin-left: 12px;
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: var(--color-text-secondary);
}

.pagination .jump-input {
  width: 44px;
  height: 30px;
  padding: 0 6px;
  border-radius: 6px;
  border: 1px solid var(--color-border-light);
  background: var(--color-surface);
  font-size: 13px;
  text-align: center;
  color: var(--color-text);
  outline: none;
  transition: border-color 0.15s;
  -moz-appearance: textfield;
  appearance: textfield;
}

.pagination .jump-input::-webkit-outer-spin-button,
.pagination .jump-input::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}

.pagination .jump-input:focus {
  border-color: var(--color-primary);
}

.pagination .jump-btn {
  min-width: unset;
  width: 32px;
  height: 30px;
  padding: 0;
  font-size: 12px;
  font-weight: 600;
  border-radius: 6px;
}

/* ── Teacher List ── */
.teacher-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

/* ── Transition ── */
.card-list-enter-active {
  transition: all 0.35s ease;
}
.card-list-leave-active {
  transition: all 0.2s ease;
}
.card-list-enter-from {
  opacity: 0;
  transform: translateY(12px);
}
.card-list-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}

/* ── Empty State ── */
.empty-state {
  text-align: center;
  padding: 64px 20px;
  color: var(--color-text-muted);
}

.empty-icon {
  color: var(--color-text-muted);
  margin-bottom: 14px;
  opacity: 0.5;
}

.empty-state p {
  font-size: 15px;
  line-height: 1.6;
}

.empty-hint {
  font-size: 13px !important;
  margin-top: 4px;
}

/* ── Footer ── */
.app-footer {
  text-align: center;
  padding: 20px;
  font-size: 12px;
  color: var(--color-text-muted);
}

.footer-heart {
  color: #ef4444;
  vertical-align: -2px;
  display: inline;
}

.footer-version {
  margin-top: 4px;
  font-size: 11px;
  opacity: 0.55;
}

.footer-version code {
  font-family: 'SF Mono', 'Fira Code', 'Consolas', monospace;
  background: var(--color-border-light);
  padding: 1px 6px;
  border-radius: 4px;
  font-size: 11px;
}

/* ── Responsive ── */
@media (max-width: 768px) {
  .app-title {
    font-size: 24px;
    gap: 8px;
  }
  .app-header {
    padding: 28px 16px 48px;
  }
  .stat-num {
    font-size: 20px;
  }
  .stat-label {
    font-size: 12px;
  }
  .app-main {
    padding: 20px 12px 36px;
  }
}

@media (max-width: 480px) {
  .app-title {
    font-size: 20px;
    gap: 6px;
  }
  .app-header {
    padding: 24px 12px 44px;
  }
  .header-stats {
    gap: 12px;
  }
  .stat-num {
    font-size: 18px;
  }
  .stat {
    gap: 4px;
  }
  .empty-state {
    padding: 40px 16px;
  }
}
</style>
