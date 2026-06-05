<script setup>
import { computed } from 'vue'
import AppIcon from './AppIcon.vue'

const props = defineProps({
  teacher: { type: Object, required: true },
  collegeName: { type: String, default: '' },
})

const emit = defineEmits(['click'])

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

const rateDisplay = computed(() => {
  const r = props.teacher.rating
  if (r == null || r === 0) return 'N/A'
  return r.toFixed(1)
})

const rateClass = computed(() => {
  if (rateDisplay.value === 'N/A') return 'rate-na'
  const n = parseFloat(rateDisplay.value)
  if (n >= 8) return 'rate-high'
  if (n >= 6) return 'rate-mid'
  return 'rate-low'
})

const hotClass = computed(() => {
  const h = props.teacher.hotness
  if (h >= 80) return 'hot-high'
  if (h >= 50) return 'hot-mid'
  return 'hot-low'
})

const hotDisplay = computed(() => {
  const h = props.teacher.hotness
  if (h >= 10000) return (h / 10000).toFixed(1) + 'w'
  if (h >= 1000) return (h / 1000).toFixed(1) + 'k'
  return String(h)
})

const initials = computed(() => props.teacher.name ? props.teacher.name.charAt(0) : '?')
</script>

<template>
  <article class="teacher-card" @click="emit('click', teacher)">
    <div class="avatar" :style="avatarStyle" aria-hidden="true">
      {{ initials }}
    </div>

    <div class="info">
      <div class="name-row">
        <h3 class="name">{{ teacher.name }}</h3>
        <span class="college-tag" :title="collegeName">{{ collegeName }}</span>
      </div>
      <div class="meta">
        <code>{{ teacher.pinyin }}</code>
        <code>{{ teacher.pinyin_abbr }}</code>
      </div>
    </div>

    <div class="badges">
      <span :class="['badge', 'badge-rate', rateClass]" :aria-label="`评分 ${rateDisplay}`">
        <AppIcon name="star" :size="14" class="badge-icon" />
        {{ rateDisplay === 'N/A' ? '暂无' : rateDisplay }}
      </span>

      <span :class="['badge', 'badge-hot', hotClass]" :aria-label="`热度 ${hotDisplay}`">
        <AppIcon name="flame" :size="14" class="badge-icon" />
        {{ hotDisplay }}
      </span>

      <span class="badge badge-comments">
        <AppIcon name="comment" :size="14" class="badge-icon" />
        {{ teacher.rating_count }}
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
  cursor: pointer;
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

/* Comments count */
.badge-comments {
  background: var(--color-bg);
  color: var(--color-text-secondary);
  font-size: 12px;
}

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
