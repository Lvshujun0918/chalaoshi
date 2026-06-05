<script setup>
import { ref, watch, onMounted, onUnmounted, computed } from 'vue'
import AppIcon from './AppIcon.vue'

const props = defineProps({
  teacher: { type: Object, required: true },
})

const emit = defineEmits(['close'])

const comments = ref([])
const courses = ref([])
const loading = ref(true)
const error = ref('')
const activeTab = ref('comments')
const page = ref(1)
const totalPages = ref(1)
const totalComments = ref(0)
const sortBy = ref('time')

const API_BASE = '/api'

async function loadComments(pageNum = 1) {
  loading.value = true
  error.value = ''
  try {
    const params = new URLSearchParams({
      page: String(pageNum),
      page_size: '20',
      sort_by: sortBy.value,
    })
    const res = await fetch(`${API_BASE}/teachers/${props.teacher.id}/comments?${params}`)
    if (!res.ok) throw new Error('加载失败')
    const data = await res.json()
    comments.value = data.comments
    page.value = data.page
    totalPages.value = data.total_pages
    totalComments.value = data.total
  } catch (e) {
    error.value = '加载评论失败: ' + e.message
  } finally {
    loading.value = false
  }
}

async function loadCourses() {
  try {
    const res = await fetch(`${API_BASE}/teachers/${props.teacher.id}`)
    if (!res.ok) return
    const data = await res.json()
    courses.value = data.courses || []
  } catch (e) {
    // ignore
  }
}

function formatTime(t) {
  if (!t) return ''
  const d = new Date(t)
  const now = new Date()
  const diff = now - d
  if (diff < 86400000) return '今天 ' + d.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  if (diff < 172800000) return '昨天'
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  if (y === now.getFullYear()) return `${m}-${day}`
  return `${y}-${m}-${day}`
}

function goPage(p) {
  if (p < 1 || p > totalPages.value) return
  loadComments(p)
}

function changeSort(s) {
  sortBy.value = s
  loadComments(1)
}

function handleEsc(e) {
  if (e.key === 'Escape') emit('close')
}

const initials = computed(() => props.teacher.name ? props.teacher.name.charAt(0) : '?')

const AVATAR_GRADIENTS = [
  'linear-gradient(135deg, #667eea, #764ba2)',
  'linear-gradient(135deg, #f093fb, #f5576c)',
  'linear-gradient(135deg, #4facfe, #00f2fe)',
  'linear-gradient(135deg, #43e97b, #38f9d7)',
  'linear-gradient(135deg, #fa709a, #fee140)',
  'linear-gradient(135deg, #a18cd1, #fbc2eb)',
  'linear-gradient(135deg, #fccb90, #d57eeb)',
  'linear-gradient(135deg, #e0c3fc, #8ec5fc)',
  'linear-gradient(135deg, #f5576c, #ff9a9e)',
  'linear-gradient(135deg, #667eea, #60a5fa)',
]

const avatarStyle = computed(() => ({
  background: AVATAR_GRADIENTS[props.teacher.id % AVATAR_GRADIENTS.length],
}))

const ratingDisplay = computed(() => {
  const r = props.teacher.rating
  if (!r) return '暂无评分'
  return r.toFixed(1) + ' 分'
})

const pageNumbers = computed(() => {
  const pages = []
  const total = totalPages.value
  const current = page.value
  let start = Math.max(1, current - 2)
  let end = Math.min(total, current + 2)
  if (end - start < 4) {
    if (start === 1) end = Math.min(total, start + 4)
    else start = Math.max(1, end - 4)
  }
  for (let i = start; i <= end; i++) pages.push(i)
  return pages
})

watch(() => props.teacher, (t) => {
  if (t) {
    loadComments(1)
    loadCourses()
  }
}, { immediate: true })

onMounted(() => {
  document.addEventListener('keydown', handleEsc)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleEsc)
})
</script>

<template>
  <Teleport to="body">
    <div class="modal-overlay" @click.self="emit('close')">
      <div class="modal-panel">
        <!-- Header -->
        <div class="modal-header">
          <div class="modal-teacher-info">
            <div class="avatar" :style="avatarStyle">{{ initials }}</div>
            <div>
              <h2 class="modal-teacher-name">{{ teacher.name }}</h2>
              <p class="modal-teacher-meta">
                {{ teacher.department }} · {{ ratingDisplay }} ·
                {{ teacher.rating_count }} 条评论
              </p>
            </div>
          </div>
          <button class="modal-close" @click="emit('close')" aria-label="关闭">
            <AppIcon name="close" :size="20" />
          </button>
        </div>

        <!-- Tabs -->
        <div class="modal-tabs">
          <button
            :class="['tab', { active: activeTab === 'comments' }]"
            @click="activeTab = 'comments'"
          >
            评论 ({{ totalComments }})
          </button>
          <button
            v-if="courses.length > 0"
            :class="['tab', { active: activeTab === 'courses' }]"
            @click="activeTab = 'courses'"
          >
            课程 ({{ courses.length }})
          </button>
        </div>

        <!-- Comments Tab -->
        <div v-if="activeTab === 'comments'" class="modal-body">
          <!-- Sort -->
          <div class="comment-sort">
            <span>排序：</span>
            <button :class="{ active: sortBy === 'time' }" @click="changeSort('time')">最新</button>
            <button :class="{ active: sortBy === 'likes' }" @click="changeSort('likes')">最多赞</button>
            <button :class="{ active: sortBy === 'net_likes' }" @click="changeSort('net_likes')">最热</button>
          </div>

          <!-- Loading -->
          <div v-if="loading" class="modal-loading">加载中...</div>

          <!-- Error -->
          <div v-else-if="error" class="modal-error">{{ error }}</div>

          <!-- Empty -->
          <div v-else-if="comments.length === 0" class="modal-empty">
            <AppIcon name="comment" :size="48" class="empty-icon" />
            <p>暂无评论</p>
          </div>

          <!-- Comment List -->
          <div v-else class="comment-list">
            <div v-for="c in comments" :key="c.id" class="comment-item">
              <div class="comment-meta">
                <span class="comment-time">{{ formatTime(c.publish_time) }}</span>
                <span class="comment-likes" :class="{ positive: c.net_likes > 0, negative: c.net_likes < 0 }">
                  <AppIcon name="like" :size="14" />
                  {{ c.net_likes > 0 ? '+' : '' }}{{ c.net_likes }}
                </span>
              </div>
              <p class="comment-content">{{ c.content }}</p>
              <div class="comment-votes">
                <span class="votes-up">👍 {{ c.likes }}</span>
                <span class="votes-down">👎 {{ c.dislikes }}</span>
              </div>
            </div>
          </div>

          <!-- Pagination -->
          <div v-if="totalPages > 1" class="modal-pagination">
            <button :disabled="page <= 1" @click="goPage(page - 1)">‹ 上一页</button>
            <button
              v-for="p in pageNumbers"
              :key="p"
              :class="{ current: p === page }"
              @click="goPage(p)"
            >{{ p }}</button>
            <button :disabled="page >= totalPages" @click="goPage(page + 1)">下一页 ›</button>
          </div>
        </div>

        <!-- Courses Tab -->
        <div v-if="activeTab === 'courses'" class="modal-body">
          <div class="course-list">
            <div v-for="c in courses" :key="c.course_name" class="course-item">
              <div class="course-name">{{ c.course_name }}</div>
              <div class="course-stats">
                <span class="course-gpa" :class="{ 'gpa-high': c.gpa >= 4, 'gpa-low': c.gpa < 3 }">
                  GPA {{ c.gpa?.toFixed(2) }}
                </span>
                <span class="course-count">{{ c.count }}人</span>
                <span class="course-std">σ {{ c.std_dev?.toFixed(2) }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.modal-panel {
  background: var(--color-surface);
  border-radius: 16px;
  width: 100%;
  max-width: 680px;
  max-height: 85vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 20px 60px rgba(0,0,0,0.3);
  animation: slideUp 0.3s ease;
}

@keyframes slideUp {
  from { transform: translateY(30px); opacity: 0; }
  to { transform: translateY(0); opacity: 1; }
}

/* Header */
.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24px 24px 16px;
  border-bottom: 1px solid var(--color-border-light);
}

.modal-teacher-info {
  display: flex;
  align-items: center;
  gap: 14px;
}

.modal-teacher-info .avatar {
  width: 48px;
  height: 48px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  font-weight: 700;
  color: #fff;
  flex-shrink: 0;
}

.modal-teacher-name {
  font-size: 20px;
  font-weight: 700;
  color: var(--color-text);
  margin: 0;
}

.modal-teacher-meta {
  font-size: 13px;
  color: var(--color-text-secondary);
  margin: 4px 0 0;
}

.modal-close {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  border: none;
  background: var(--color-bg);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.2s;
  color: var(--color-text);
}
.modal-close:hover {
  background: var(--color-border-light);
}

/* Tabs */
.modal-tabs {
  display: flex;
  padding: 0 24px;
  border-bottom: 1px solid var(--color-border-light);
  gap: 0;
}

.tab {
  padding: 12px 20px;
  border: none;
  background: none;
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text-secondary);
  cursor: pointer;
  border-bottom: 2px solid transparent;
  transition: all 0.2s;
}
.tab:hover {
  color: var(--color-text);
}
.tab.active {
  color: var(--color-primary);
  border-bottom-color: var(--color-primary);
}

/* Body */
.modal-body {
  flex: 1;
  overflow-y: auto;
  padding: 20px 24px;
}

/* Sort */
.comment-sort {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 16px;
  font-size: 13px;
  color: var(--color-text-secondary);
}

.comment-sort button {
  padding: 4px 12px;
  border-radius: 14px;
  border: 1px solid var(--color-border-light);
  background: var(--color-bg);
  cursor: pointer;
  font-size: 12px;
  color: var(--color-text-secondary);
  transition: all 0.2s;
}
.comment-sort button.active {
  background: var(--color-primary);
  color: #fff;
  border-color: var(--color-primary);
}

/* Loading/Error/Empty */
.modal-loading,
.modal-error,
.modal-empty {
  text-align: center;
  padding: 40px 20px;
  color: var(--color-text-secondary);
}

.modal-error { color: var(--color-rate-low); }

.empty-icon { opacity: 0.3; margin-bottom: 12px; }

/* Comment List */
.comment-item {
  padding: 16px 0;
  border-bottom: 1px solid var(--color-border-light);
}
.comment-item:last-child {
  border-bottom: none;
}

.comment-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.comment-time {
  font-size: 12px;
  color: var(--color-text-tertiary);
}

.comment-likes {
  font-size: 13px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 3px;
}
.comment-likes.positive { color: #22c55e; }
.comment-likes.negative { color: #ef4444; }

.comment-content {
  font-size: 14px;
  line-height: 1.7;
  color: var(--color-text);
  white-space: pre-wrap;
  word-break: break-word;
  margin: 0;
}

.comment-votes {
  margin-top: 8px;
  font-size: 12px;
  color: var(--color-text-tertiary);
  display: flex;
  gap: 12px;
}

/* Pagination */
.modal-pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  margin-top: 20px;
  padding-top: 16px;
  border-top: 1px solid var(--color-border-light);
}

.modal-pagination button {
  min-width: 32px;
  height: 32px;
  padding: 0 8px;
  border-radius: 8px;
  border: 1px solid var(--color-border-light);
  background: var(--color-surface);
  cursor: pointer;
  font-size: 13px;
  color: var(--color-text-secondary);
  transition: all 0.15s;
}
.modal-pagination button:hover:not(:disabled) {
  border-color: var(--color-primary);
  color: var(--color-primary);
}
.modal-pagination button.current {
  background: var(--color-primary);
  color: #fff;
  border-color: var(--color-primary);
}
.modal-pagination button:disabled {
  opacity: 0.4;
  cursor: default;
}

/* Course List */
.course-item {
  padding: 14px 0;
  border-bottom: 1px solid var(--color-border-light);
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
}
.course-item:last-child {
  border-bottom: none;
}

.course-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--color-text);
  flex: 1;
}

.course-stats {
  display: flex;
  gap: 10px;
  font-size: 12px;
  color: var(--color-text-secondary);
  flex-shrink: 0;
}

.course-gpa {
  font-weight: 600;
}
.gpa-high { color: #22c55e; }
.gpa-low { color: #ef4444; }

/* Responsive */
@media (max-width: 600px) {
  .modal-panel {
    max-height: 95vh;
    border-radius: 12px 12px 0 0;
    margin-top: auto;
  }
  .modal-header {
    padding: 18px 16px 12px;
  }
  .modal-body {
    padding: 16px;
  }
  .modal-teacher-name { font-size: 17px; }
}
</style>
