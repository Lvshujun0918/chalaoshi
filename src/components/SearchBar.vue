<script setup>
import AppIcon from './AppIcon.vue'

const props = defineProps({
  modelValue: { type: String, default: '' },
  placeholder: { type: String, default: '输入关键词搜索…' },
})

const emit = defineEmits(['update:modelValue', 'clear'])

function onInput(e) {
  emit('update:modelValue', e.target.value)
}

function onClear() {
  emit('update:modelValue', '')
  emit('clear')
}
</script>

<template>
  <div class="search-bar">
    <div class="search-input-wrap">
      <AppIcon name="search" :size="20" class="search-icon" />

      <input
        :value="modelValue"
        @input="onInput"
        type="text"
        :placeholder="placeholder"
        autofocus
        class="search-input"
        aria-label="搜索教师"
      />

      <button
        v-if="modelValue"
        @click="onClear"
        class="clear-btn"
        aria-label="清除搜索"
        title="清除"
      >
        <AppIcon name="x" :size="16" />
      </button>
    </div>

    <div class="search-hints">
      <span class="hint-tag">
        <AppIcon name="type" :size="13" class="hint-icon" />
        拼音 zhangsan
      </span>
      <span class="hint-tag">
        <AppIcon name="scissors" :size="12" class="hint-icon" />
        缩写 zs
      </span>
      <span class="hint-tag">
        <AppIcon name="languages" :size="13" class="hint-icon" />
        汉字 张三
      </span>
    </div>
  </div>
</template>

<style scoped>
.search-bar {
  max-width: 600px;
  margin: -24px auto 0;
  padding: 0 20px;
  position: relative;
  z-index: 10;
}

.search-input-wrap {
  display: flex;
  align-items: center;
  background: var(--color-surface);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-lg);
  padding: 4px 4px 4px 18px;
  transition: box-shadow var(--transition-fast);
}

.search-input-wrap:focus-within {
  box-shadow: 0 0 0 3px rgba(91,94,247,.22), var(--shadow-lg);
}

.search-icon {
  color: var(--color-text-muted);
  flex-shrink: 0;
  margin-right: 10px;
}

.search-input {
  flex: 1;
  border: none;
  outline: none;
  font-size: 16px;
  font-family: inherit;
  background: transparent;
  color: var(--color-text);
  padding: 12px 0;
  min-width: 0;
}

.search-input::placeholder {
  color: var(--color-text-muted);
}

.clear-btn {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: none;
  border-radius: 50%;
  background: transparent;
  color: var(--color-text-muted);
  cursor: pointer;
  transition: background var(--transition-fast), color var(--transition-fast);
}

.clear-btn:hover {
  background: var(--color-bg);
  color: var(--color-text-secondary);
}

.search-hints {
  display: flex;
  justify-content: center;
  gap: 10px;
  margin-top: 10px;
  flex-wrap: wrap;
}

.hint-tag {
  font-size: 12px;
  color: var(--color-text-muted);
  background: var(--color-surface);
  padding: 3px 10px;
  border-radius: 20px;
  border: 1px solid var(--color-border);
  white-space: nowrap;
  display: inline-flex;
  align-items: center;
  gap: 5px;
}

.hint-icon {
  flex-shrink: 0;
}

@media (max-width: 768px) {
  .search-bar {
    padding: 0 16px;
  }
}

@media (max-width: 480px) {
  .search-input-wrap {
    padding: 2px 2px 2px 12px;
  }
  .search-input {
    font-size: 15px;
  }
  .hint-tag {
    font-size: 11px;
    padding: 2px 8px;
  }
}
</style>
