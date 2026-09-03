<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { cn } from '@/lib/utils'

const checked = defineModel<boolean>('checked', { default: false })

const props = defineProps<{
  id?: string
  class?: HTMLAttributes['class']
  disabled?: boolean
}>()

function toggle() {
  if (props.disabled) return
  checked.value = !checked.value
}
</script>

<template>
  <button
    :id="id"
    type="button"
    role="switch"
    :aria-checked="checked"
    :disabled="disabled"
    :class="cn(
      'relative inline-flex h-10 w-12 shrink-0 items-center justify-center rounded bg-transparent outline-none focus-visible:ring-2 focus-visible:ring-ring/50 disabled:opacity-50',
      props.class,
    )"
    @click="toggle"
  >
    <span :class="cn('switch-track', checked ? 'bg-primary' : 'bg-muted')">
      <span :class="cn('switch-thumb', checked ? 'translate-x-4' : 'translate-x-0')" />
    </span>
  </button>
</template>

<style scoped>
.switch-track {
  display: flex;
  width: 2.25rem;
  height: 1.25rem;
  align-items: center;
  border: 1px solid var(--border);
  border-radius: 999px;
  padding: 1px;
  transition: background-color 150ms ease;
}

.switch-thumb {
  display: block;
  width: 1rem;
  height: 1rem;
  border-radius: 50%;
  background: var(--background);
  box-shadow: 0 1px 2px rgb(0 0 0 / 20%);
  transition: transform 150ms ease;
}
</style>
