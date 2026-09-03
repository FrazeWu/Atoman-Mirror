<script setup lang="ts">
defineProps<{
  eyebrow: string
  title: string
  subtitle: string
}>()
</script>

<template>
  <header class="page-hero">
    <p class="eyebrow">{{ eyebrow }}</p>
    <h1 class="display-title">{{ title }}</h1>
    <p class="max-w-2xl text-base leading-7 text-muted-foreground sm:text-lg">
      {{ subtitle }}
    </p>
    <slot />
  </header>
</template>
