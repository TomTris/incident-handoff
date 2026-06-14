<script setup lang="ts">
import { whoAmI } from '@/api';
import type { Incident, IncidentStatus, Severity, UserContext } from '@/types';
import { makeEmptyUserContext } from '@/utils/user';
import { computed, onMounted, ref, watch } from 'vue';

const props = defineProps<{
    inc: Incident | undefined;
    errorUpdateIncidentMsg: string;
}>()

const emit = defineEmits<{
    'updateIncident' : [payload : {severity: Severity, status: IncidentStatus, on_call: string}]
}>()

const identity = ref<UserContext>(makeEmptyUserContext())
const updateOnCall = ref('')
const updateSeverity = ref<Severity>("SEV1")
const updateStatus = ref<IncidentStatus>("triggered")

watch(
  () => props.inc,
  (inc) => {
    if (!inc) return
    updateOnCall.value = inc.on_call
    updateSeverity.value = inc.severity
    updateStatus.value = inc.status
  },
  { immediate: true }
)

onMounted( async() => {
  try {
    identity.value = await whoAmI()
  } catch (e) {
      alert('Something is wrong.\n' + (e as Error).message)
      window.location.href = "/"
  }
})


const totalEntry = computed(() => props.inc?.entries.length || 0)
const takenActions = computed(() => props.inc?.entries.filter((e) => e.type == "action").length || 0)
const openQuestion = computed(() => props.inc?.entries.filter((e) => e.type == "open_question").length || 0)
const handoffCount = computed(() => {
  const entries = props.inc?.entries
  if (!entries?.length) return 0

  const ordered = [...entries].sort((a, b) => a.created_at.localeCompare(b.created_at))

  let handoffs = -1
  let current = 'not-exists--'
  for (const e of ordered) {
    if (e.author !== current) {
      handoffs++
      current = e.author
    }
  }
  return handoffs
})

function onIncidentUpdate() {
    emit('updateIncident', {severity: updateSeverity.value, status: updateStatus.value, on_call: updateOnCall.value})
}

function isIncidentOwner() : boolean {
  if (identity.value.role == "admin") {
    return true
  }
  if (identity.value.role == "engineer" && identity.value.username == props.inc?.on_call) {
    return true
  }
  return false
}

</script>

<template>
    <div class="detail-side">
        <div class="panel">
        <h3 class="panel-title">Handoff brief</h3>
        <div class="stat-grid">
            <div class="stat">
            <span class="stat-value mono">{{ totalEntry }}</span>
            <span class="stat-label">entries</span>
            </div>
            <div class="stat">
            <span class="stat-value mono">{{ takenActions }}</span>
            <span class="stat-label">actions</span>
            </div>
            <div class="stat">
            <span class="stat-value mono accent">{{ openQuestion }}</span>
            <span class="stat-label">open questions</span>
            </div>
            <div class="stat">
            <span class="stat-value mono">{{ handoffCount }}</span>
            <span class="stat-label">handoffs</span>
            </div>
        </div>
        </div>

        <div class="panel">
        <h3 class="panel-title">Update incident</h3>
        <p class="admin-note mono">ADMIN / ON-CALL ONLY</p>
        <div class="field">
            <label class="field-label">Status</label>
            <select class="select" :disabled="!isIncidentOwner()" v-model="updateStatus">
            <option value="triggered">Triggered</option>
            <option value="acknowledged">Acknowledged</option>
            <option value="investigating">Investigating</option>
            <option value="mitigated">Mitigated</option>
            <option value="resolved">Resolved</option>
            </select>
        </div>
        <div class="field">
            <label class="field-label">Severity</label>
            <select class="select" :disabled="!isIncidentOwner()" v-model="updateSeverity">
            <option value="SEV1">SEV1</option>
            <option value="SEV2">SEV2</option>
            <option value="SEV3">SEV3</option>
            </select>
        </div>
        <div class="field">
            <label class="field-label">On-call</label>
            <input
            class="input"
            type="text"
            :disabled="!isIncidentOwner()" 
            v-model="updateOnCall"
            placeholder="username to hand the pager to"
            />
        </div>
        <button v-if="isIncidentOwner()" class="btn btn-primary btn-block" @click.prevent="onIncidentUpdate">Apply update</button>
        <button
            v-else
            class="btn btn-block"
            @click.prevent="onIncidentUpdate"
            :disabled="!isIncidentOwner()"
            title='You are not admin nor on-call engineer to update this incident'
          >
          Apply update
        </button>
        <p class="error">{{ errorUpdateIncidentMsg }}</p>
        </div>
    </div>
</template>

<style scope>

.detail-side {
  display: flex;
  flex-direction: column;
  gap: 20px;
  width: 300px;
}

/* Handoff brief stats */
.stat-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 14px;
}

.stat {
  background-color: var(--color-input);
  border: 1px solid var(--color-border);
  border-radius: 5px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 14px;
  width: 116px;
}

.stat-value {
  color: var(--color-text-bright);
  font-size: 26px;
  font-weight: 700;
}

.stat-value.accent {
  color: var(--color-accent);
}

.stat-label {
  color: var(--color-text-dim);
  font-size: 11px;
  letter-spacing: 1px;
  text-transform: uppercase;
}

.admin-note {
  color: var(--color-text-dim);
  font-size: 10px;
  letter-spacing: 2px;
  margin-bottom: 14px;
}

@media (max-width: 768px) {
  .stat {
    flex: 1;
    min-width: 116px;
  }
}
</style>