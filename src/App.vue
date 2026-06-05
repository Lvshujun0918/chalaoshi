<script setup>
import { useTeacherSearch } from './composables/useTeacherSearch.js'
import SearchBar from './components/SearchBar.vue'
import TeacherCard from './components/TeacherCard.vue'
import AppIcon from './components/AppIcon.vue'

const {
  query,
  results,
  displayResults,
  totalTeachers,
  totalColleges,
  getCollegeName,
  clearQuery,
} = useTeacherSearch()
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
        <p class="app-subtitle">教师评分查询系统</p>
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

    <!-- Search -->
    <SearchBar v-model="query" @clear="clearQuery" />

    <!-- Main Content -->
    <main class="app-main">
      <!-- Result Summary -->
      <div v-if="query && results.length > 0" class="result-summary">
        <span>找到 <strong>{{ results.length }}</strong> 条结果</span>
        <span v-if="results.length > 20" class="result-note">仅显示前 20 条</span>
      </div>

      <!-- Teacher List -->
      <TransitionGroup
        v-if="displayResults.length > 0"
        name="card-list"
        tag="div"
        class="teacher-list"
      >
        <TeacherCard
          v-for="t in displayResults"
          :key="t.id"
          :teacher="t"
          :college-name="getCollegeName(t.xy)"
        />
      </TransitionGroup>

      <!-- No Results -->
      <div v-else-if="query" class="empty-state">
        <AppIcon name="search-lg" :size="48" class="empty-icon" />
        <p>没有找到匹配的教师</p>
        <p class="empty-hint">试试其他关键词吧</p>
      </div>

      <!-- Welcome -->
      <div v-else class="empty-state welcome">
        <AppIcon name="pointer" :size="48" class="empty-icon" />
        <p>输入教师姓名、拼音或缩写开始查询</p>
      </div>
    </main>

    <!-- Footer -->
    <footer class="app-footer">
      <p>
        Made with
        <AppIcon name="heart" :size="12" class="footer-heart" />
        · 数据仅供参考
      </p>
    </footer>
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

.result-note {
  font-size: 12px;
  color: var(--color-text-muted);
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
