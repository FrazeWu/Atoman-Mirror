import { readFileSync } from 'node:fs'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'
import test from 'node:test'
import assert from 'node:assert/strict'

const root = resolve(dirname(fileURLToPath(import.meta.url)), '..')
const mirrorRoot = resolve(root, '..')
const readWeb = (path) => readFileSync(resolve(root, path), 'utf8')
const readMirror = (path) => readFileSync(resolve(mirrorRoot, path), 'utf8')

test('product controls keep a 40px minimum target and visible focus', () => {
  const button = readWeb('src/components/ui/Button.vue')
  const switchControl = readWeb('src/components/ui/Switch.vue')

  assert.match(button, /sm:\s*'h-10/)
  assert.match(button, /focus-visible:ring-2/)
  assert.match(switchControl, /h-10 w-12/)
  assert.match(switchControl, /switch-track/)
})

test('shell exposes service state and mobile navigation remains usable', () => {
  const shell = readWeb('src/components/AppShell.vue')
  const router = readWeb('src/router/index.ts')
  const docsConfig = readMirror('docs/astro.config.mjs')

  assert.match(shell, /aria-live="polite"/)
  assert.match(shell, /serviceStatus/)
  assert.match(shell, /overflow-x-auto/)
  assert.match(router, /Atoman Mirror/)
  assert.match(docsConfig, /title: 'Atoman Mirror'/)
})

test('async product feedback is announced to assistive technology', () => {
  const images = readWeb('src/pages/ImagesPage.vue')
  const search = readWeb('src/pages/SearchPage.vue')

  assert.match(images, /role="alert"/)
  assert.match(images, /role="status"/)
  assert.match(search, /aria-label="上一页"/)
  assert.match(search, /aria-label="下一页"/)
  assert.match(search, /aria-live="polite"/)
})

test('documentation home is compact and does not reserve the first viewport for artwork', () => {
  const styles = readMirror('docs/src/styles/custom.css')
  const zhHome = readMirror('docs/src/content/docs/index.mdx')
  const enHome = readMirror('docs/src/content/docs/en/index.mdx')

  assert.doesNotMatch(zhHome, /\n\s*image:\s*\n/)
  assert.doesNotMatch(enHome, /\n\s*image:\s*\n/)
  assert.match(styles, /\.hero h1\s*\{/)
  assert.match(styles, /max-width:\s*52rem/)
  assert.match(styles, /border-radius:\s*4px/)
})
