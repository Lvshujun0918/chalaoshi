<script setup>
import { computed } from 'vue'
import AppIcon from './AppIcon.vue'

const props = defineProps({
  teacher: { type: Object, required: true },
  collegeName: { type: String, default: '' },
})

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

const rateClass = computed(() => {
  if (props.teacher.rate === 'N/A') return 'rate-na'
  const n = parseFloat(props.teacher.rate)
  if (n >= 8) return 'rate-high'
  if (n >= 6) return 'rate-mid'
  return 'rate-low'
})

const hotClass = computed(() => {
  if (props.teacher.hot >= 80) return 'hot-high'
  if (props.teacher.hot >= 50) return 'hot-mid'
  return 'hot-low'
})

const initials = computed(() => props.teacher.name.charAt(0))
</script>

<template>
  <article class="teacher-card">
    <div class="avatar" :style="avatarStyle" aria-hidden="true">
      {{ initials }}
    </div>

    <div class="info">
      <div class="name-row">
        <h3 class="name">{{ teacher.name }}</h3>
        <span class="college-tag" :title="collegeName">{{ collegeName }}</span>
      </div>
      <div class="meta">
        <code>{{ teacher.py }}</code>
        <code>{{ teacher.sx }}</code>
      </div>
    </div>

    <div class="badges">
      <span :class="['badge', 'badge-rate', rateClass]" :aria-label="`评分 ${teacher.rate}`">
        <AppIcon name="star" :size="14" class="badge-icon" />
        {{ teacher.rate === 'N/A' ? '暂无' : teacher.rate }}
      </span>

      <span :class="['badge', 'badge-hot', hotClass]" :aria-label="`热度 ${teacher.hot}`">
        <AppIcon name="flame" :size="14" class="badge-icon" />
        {{ teacher.hot }}
      </span>
    </div>
  </article>
</template>

<style scoped>
.teacher-card {
  display: flex;
  align-items: center;
  gap: 16px;
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  padding: 16px 20px;
  box-shadow: var(--shadow-sm);
  border: 1px solid transparent;
  transition: transform var(--transition-fast), box-shadow var(--transition-fast), border-color var(--transition-fast);
  cursor: default;
}

.teacher-card:hover {
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
}

.college-tag {
  font-size: 11px;
  padding: 1px 8px;
  border-radius: 20px;
  background: var(--color-primary-light);
  color: var(--color-primary);
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 180px;
  flex-shrink: 1;
}

.meta {
  display: flex;
  gap: 8px;
  font-size: 12px;
}

.meta code {
  background: var(--color-bg);
  padding: 1px 6px;
  border-radius: 4px;
  font-size: 11px;
  color: var(--color-text-secondary);
  font-family: var(--font-mono);
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

/* Rate */
.badge-rate.rate-high {
  background: var(--color-rate-high-bg);
  color: var(--color-rate-high);
}
.badge-rate.rate-mid {
  background: var(--color-rate-mid-bg);
  color: var(--color-rate-mid);
}
.badge-rate.rate-low {
  background: var(--color-rate-low-bg);
  color: var(--color-rate-low);
}
.badge-rate.rate-na {
  background: var(--color-rate-na-bg);
  color: var(--color-rate-na);
}

/* Hot */
.badge-hot.hot-high { color: var(--color-hot-high); }
.badge-hot.hot-mid  { color: var(--color-hot-mid); }
.badge-hot.hot-low  { color: var(--color-hot-low); }
.badge-hot { background: transparent; padding-left: 0; padding-right: 0; }

/* Icons */
.badge-icon { flex-shrink: 0; }

/* Responsive */
@media (max-width: 768px) {
  .teacher-card {
    gap: 12px;
    padding: 14px 16px;
  }
  .college-tag {
    max-width: 140px;
  }
}

@media (max-width: 520px) {
  .teacher-card {
    flex-wrap: wrap;
    gap: 10px;
    padding: 12px 14px;
  }
  .badges {
    width: 100%;
    justify-content: flex-end;
  }
  .college-tag {
    max-width: 100px;
    font-size: 10px;
  }
  .name {
    font-size: 15px;
  }
}
</style>
