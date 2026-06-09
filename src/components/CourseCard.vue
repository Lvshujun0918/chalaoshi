<script setup>
import { computed } from 'vue'
import AppIcon from './AppIcon.vue'

const props = defineProps({
  course: { type: Object, required: true },
})

const emit = defineEmits(['click'])

const AVATAR_GRADIENTS = [
  'linear-gradient(135deg, #0ea5e9, #06b6d4)',
  'linear-gradient(135deg, #10b981, #34d399)',
  'linear-gradient(135deg, #14b8a6, #2dd4bf)',
  'linear-gradient(135deg, #06b6d4, #22d3ee)',
  'linear-gradient(135deg, #0891b2, #06b6d4)',
  'linear-gradient(135deg, #0d9488, #14b8a6)',
  'linear-gradient(135deg, #2563eb, #3b82f6)',
  'linear-gradient(135deg, #0284c7, #38bdf8)',
  'linear-gradient(135deg, #059669, #34d399)',
  'linear-gradient(135deg, #0f766e, #14b8a6)',
]

const avatarStyle = computed(() => ({
  background: AVATAR_GRADIENTS[props.course.id % AVATAR_GRADIENTS.length],
}))

const gpaDisplay = computed(() => {
  const g = props.course.gpa
  if (g == null || g === 0) return 'N/A'
  return g.toFixed(2)
})

const gpaClass = computed(() => {
  if (gpaDisplay.value === 'N/A') return 'gpa-na'
  const n = parseFloat(gpaDisplay.value)
  if (n >= 4.0) return 'gpa-high'
  if (n >= 3.5) return 'gpa-mid'
  return 'gpa-low'
})

const initials = computed(() => {
  const name = props.course.course_name
  return name ? name.charAt(0) : '?'
})

const countDisplay = computed(() => {
  const c = props.course.count
  if (c >= 10000) return (c / 10000).toFixed(1) + 'w'
  if (c >= 1000) return (c / 1000).toFixed(1) + 'k'
  return String(c || 0)
})
</script>

<template>
  <article class="course-card" @click="emit('click', course)">
    <div class="avatar" :style="avatarStyle" aria-hidden="true">
      {{ initials }}
    </div>

    <div class="info">
      <div class="name-row">
        <h3 class="name">{{ course.course_name }}</h3>
      </div>
      <div class="meta">
        <span class="teacher-tag">
          <AppIcon name="user" :size="12" class="meta-icon" />
          {{ course.teacher_name }}
        </span>
        <span v-if="course.std_dev" class="std-dev-tag">
          标准差 {{ course.std_dev.toFixed(2) }}
        </span>
      </div>
    </div>

    <div class="badges">
      <span :class="['badge', 'badge-gpa', gpaClass]" :aria-label="`GPA ${gpaDisplay}`">
        <AppIcon name="bookmark" :size="14" class="badge-icon" />
        {{ gpaDisplay === 'N/A' ? '暂无' : gpaDisplay }}
      </span>

      <span class="badge badge-count">
        <AppIcon name="users" :size="14" class="badge-icon" />
        {{ countDisplay }}
      </span>
    </div>
  </article>
</template>

<style scoped>
.course-card {
  display: flex;
  align-items: center;
  gap: 16px;
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  padding: 16px 20px;
  box-shadow: var(--shadow-sm);
  border: 1px solid transparent;
  transition: transform var(--transition-fast), box-shadow var(--transition-fast), border-color var(--transition-fast);
  cursor: pointer;
}

.course-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
  border-color: var(--color-border-hover);
}

/* Avatar */
.avatar {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: 700;
  color: #fff;
  user-select: none;
}

/* Info */
.info {
  flex: 1;
  min-width: 0;
}

.name-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.name {
  font-size: 16px;
  font-weight: 700;
  color: var(--color-text);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.meta {
  display: flex;
  gap: 8px;
  font-size: 12px;
  align-items: center;
}

.teacher-tag {
  font-size: 11px;
  padding: 1px 8px;
  border-radius: 20px;
  background: var(--color-primary-light);
  color: var(--color-primary);
  font-weight: 500;
  white-space: nowrap;
  display: inline-flex;
  align-items: center;
  gap: 3px;
}

.meta-icon {
  flex-shrink: 0;
}

.std-dev-tag {
  font-size: 11px;
  color: var(--color-text-muted);
}

/* Badges */
.badges {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-shrink: 0;
}

.badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-weight: 600;
  font-size: 13px;
  padding: 5px 10px;
  border-radius: 20px;
  white-space: nowrap;
}

.badge-icon {
  flex-shrink: 0;
}

/* GPA */
.badge-gpa.gpa-high {
  background: #ecfdf5;
  color: #059669;
}
.badge-gpa.gpa-mid {
  background: #fffbeb;
  color: #d97706;
}
.badge-gpa.gpa-low {
  background: #fef2f2;
  color: #dc2626;
}
.badge-gpa.gpa-na {
  background: #f9fafb;
  color: #9ca3af;
}

/* Count */
.badge-count {
  background: var(--color-bg);
  color: var(--color-text-secondary);
}

@media (max-width: 480px) {
  .course-card {
    padding: 12px 14px;
    gap: 10px;
  }
  .badges {
    gap: 6px;
  }
  .badge {
    font-size: 11px;
    padding: 4px 8px;
  }
}
</style>
