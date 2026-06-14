<script setup lang="ts">
import type { Incident } from '@/types';
import { computed } from 'vue';
import { isoToDateAndTime } from '@/utils/user';
const props = defineProps<{
    inc: Incident | undefined;
}>()

const actionsTaken = computed(() => props.inc?.entries.filter((entry) => entry.type == "action"))
const openQuestions = computed(() => props.inc?.entries.filter((entry) => entry.type === "open_question") ?? [])
const lastOpenQuestion = computed(() => openQuestions.value.at(-1))

</script>

<template>
    <div class="catchup">
            <div class="wrapper">
            <h3 class="catchup-block">Where things stand</h3>
            <div class="catchup-block">
                <p class="catchup-label">Actions taken</p>
                <ul class="actions-list">
                    <li v-if="!actionsTaken?.length" class="catchup-empty">
                        No actions logged yet.
                    </li>
                    <li v-for="a in actionsTaken" class="action-row">
                        <span class="action-time mono dim" >{{ isoToDateAndTime(a.created_at) }}</span>
                        <span class="action-who mono dim" >{{ a.author }}</span>
                        <span class="action-text mono dim" >{{ a.text }}</span>
                    </li>
                </ul>
            </div>
        
            <div class="catchup-block">
                <p class="catchup-label">Last open questions</p>
                <ul class="open-question">
                    <li v-if="!lastOpenQuestion" class="catchup-empty">
                        No open questions
                    </li>
                    <li v-else class="open_question">
                        <span class="oq-mark">?</span>
                        <div>
                        <p class="oq-text">{{ lastOpenQuestion.text }}</p>
                        <p class="oq-meta mono dim">
                            {{ lastOpenQuestion.author }} · {{ lastOpenQuestion.created_at }}
                        </p>
                        </div>
                    </li>
                </ul>
            </div>
        </div>
    </div>
</template>

<style scope>
.wrapper {
    margin: 0 15px;
}
.catchup {
    border-left: 3px solid var(--color-accent);
}

.catchup-block {
  margin-bottom: 18px;
}

.catchup-block:last-child {
  margin-bottom: 0;
}

.catchup-label {
  color: var(--color-text-dim);
  font-family: var(--font-mono);
  font-size: 11px;
  letter-spacing: 1px;
  margin-bottom: 10px;
  text-transform: uppercase;
}

.actions-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  list-style: none;
}

.action-row {
  align-items: baseline;
  display: flex;
  gap: 10px;
}

.action-time {
  font-size: 12px;
  width: 48px;
}

.action-who {
  color: var(--color-accent);
  font-size: 13px;
  width: 65px;
}

.action-text {
  color: var(--color-text);
  flex: 1;
}

.open-question {
  background-color: rgba(210, 153, 34, 0.08);
  border: 1px solid var(--color-sev3);
  border-radius: 5px;
  display: flex;
  gap: 12px;
  padding: 14px;
}

.oq-mark {
  color: var(--color-sev3);
  font-family: var(--font-mono);
  font-size: 22px;
  font-weight: 700;
  line-height: 1;
}

.oq-text {
  color: var(--color-text-bright);
  margin-bottom: 4px;
}

.oq-meta {
  font-size: 12px;
}

.catchup-empty {
  color: var(--color-text-dim);
  font-size: 13px;
}

@media (max-width: 768px) {
  .action-row {
    flex-wrap: wrap;
    gap: 4px 10px;
  }

}
</style>