<script setup lang="ts">
import type { TimelineEntry } from '@/types';
import { isoToDateAndTime } from '@/utils/user';


const props = defineProps<{
    entry: TimelineEntry;
}>()
</script>

<template>
  <li class="entry">
    <div class="entry-spine" :class="'spine-' + entry.type"></div>
    <div class="entry-body">
      <div class="entry-head">
        <span class="entry-type" :class="'type-' + entry.type">{{ entry.type }}</span>
        <span class="entry-author mono">{{ entry.author }}</span>
        <span class="spacer"></span>
        <span class="entry-time mono dim">{{ isoToDateAndTime(entry.created_at) }}</span>
      </div>
      <p class="entry-text">{{ entry.text }}</p>
    </div>
  </li>
</template>

<style scoped>
.entry {
  display: flex;
  gap: 14px;
  list-style: none;
  padding: 4px 0;
}

.entry-spine {
  background-color: var(--color-border-strong);
  border-radius: 2px;
  flex-shrink: 0;
  width: 3px;
}

.spine-observation {
  background-color: #4a6cf7;
}
.spine-action {
  background-color: #0cce6b;
}
.spine-open_question {
  background-color: #ff3b5c;
}
.spine-discovery {
  background-color: #f5a623;
}
.spine-state_change {
  background-color: #8b5cf6;
}


.entry-body {
  flex: 1;
  padding-bottom: 16px;
}

.entry-head {
  align-items: center;
  display: flex;
  gap: 10px;
  margin-bottom: 4px;
}

.entry-type {
  color: var(--color-text-dim);
  font-family: var(--font-mono);
  font-size: 11px;
  letter-spacing: 1px;
  text-transform: uppercase;
}

.type-observation {
  color: #4a6cf7;
}
.type-action {
  color: #0cce6b;
}
.type-open_question {
  color: #ff3b5c;
}
.type-discovery {
  color: #f5a623;
}
.type-state_change {
  color: #8b5cf6;
}

.entry-author {
  color: var(--color-text-bright);
  font-size: 13px;
}

.entry-time {
  font-size: 12px;
}

.entry-text {
  color: var(--color-text);
}
</style>
