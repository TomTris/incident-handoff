<script setup lang="ts">
import { addEntry, getIncident, updateIncident, whoAmI, } from '@/api';
import type { Incident, IncidentStatus, Severity, TimelineEntryType, UserContext } from '@/types';
import { useRoute } from 'vue-router';
import AppHeader from '@/components/AppHeader.vue';
import SeverityBadge from '@/components/SeverityBadge.vue';
import StatusBadge from '@/components/StatusBadge.vue';
import PanelCatchup from '@/components/PanelCatchup.vue';
import TimelineEntryCard from '@/components/TimelineEntryCard.vue';
import AddEntry from '@/components/AddEntry.vue';
import DetailSide from '@/components/DetailSide.vue';
import { onMounted, ref } from 'vue';
import { makeEmptyUserContext } from '@/utils/user';

const route = useRoute()
const inc = ref<Incident | undefined>()
const errIncidentLoadingMsg = ref('')
const errAddEntryMsg = ref('')
const errUpdateIncidentMsg = ref('')

onMounted(async() => {
    try {
        inc.value = await getIncident(route.params.id as string)
    } catch (e) {
        errIncidentLoadingMsg.value = (e as Error).message
        return
    }
})

async function handleAddEntry(payload: { type: TimelineEntryType; text: string }) {
    if (!inc.value) {
        return
    }
    errAddEntryMsg.value = ''
    
    if (payload.text.trim() == '') {
        errAddEntryMsg.value = "Text shouldn't be empty"
        return
    }

    try {
        const newEntry = await addEntry(inc.value.id, payload.type, payload.text.trim())
        inc.value.entries.push(newEntry)
    } catch (e) {
        errAddEntryMsg.value = (e as Error).message
    }
}

async function handleIncidentUpdate(payload: {severity: Severity, status: IncidentStatus, on_call: string}) {
    errUpdateIncidentMsg.value = ''
    
    if (!inc.value) {
        return
    }

    if (payload.on_call.trim() == '') {
        errUpdateIncidentMsg.value = "on-call shouldn't be empty"
        return
    }

    try {
        await updateIncident(inc.value.id, payload)
        inc.value.status = payload.status
        inc.value.severity = payload.severity
        inc.value.on_call = payload.on_call.trim()
    } catch (e) {
        errUpdateIncidentMsg.value = (e as Error).message
    }
}
</script>

<template>
    <div>
        <AppHeader></AppHeader>
        <p class="error">{{ errIncidentLoadingMsg }}</p>
        <div class="page">
            <RouterLink :to="{name:'incidents'}" class="back mono">← Back to incident</RouterLink>
            <p v-if="errIncidentLoadingMsg" class="error"> {{ errIncidentLoadingMsg }}</p>
            <template v-else-if="inc">
                <div class="detail-head">
                    <div class="head-id-row">
                        <span class="detail-id mono">{{ inc.id }}</span>
                        <SeverityBadge :severity="inc.severity"></SeverityBadge>
                        <StatusBadge :status="inc.status"></StatusBadge>
                    </div>
                    <h1 class="detail-title"> {{ inc.title }}</h1>
                    <div class="head-meta">
                        <span><span class="meta-key">service</span> <b class="mono">{{ inc.service }}</b></span>
                        <span><span class="meta-key">on-call</span> <b class="mono">{{ inc.on_call }}</b></span>
                        <span><span class="meta-key">opened by</span> <b class="mono">{{ inc.opened_by }}</b></span>
                        <!-- <span><span class="meta-key">elapsed</span> <b class="mono">{{ inc.elapsed }}</b></span> -->
                    </div>
                </div>
            </template>

            <div class="detail-grid">
                <div class="detail-main">
                    <PanelCatchup :inc="inc"></PanelCatchup>

                    <div class="panel">
                        <h3 class="panel-title">Timeline</h3>
                        <ul class="timeline">
                            <p v-if="!inc?.entries.length" class="catchup-empty">
                                No entry yet
                            </p>
                            <TimelineEntryCard v-for="e in inc?.entries" :key="e.id" :entry="e"></TimelineEntryCard>
                        </ul>
                    </div>

                    <div class="panel add-entry">
                        <AddEntry :errAddEntryMsg="errAddEntryMsg" @add-entry="handleAddEntry" ></AddEntry>
                    </div>

                </div>
                <aside class="detail-side">
                    <DetailSide :inc='inc' :errorUpdateIncidentMsg="errUpdateIncidentMsg" @update-incident="handleIncidentUpdate"></DetailSide>
                </aside>
            </div>
        
            <div class="detail-grid">
                <div class="detail-main"></div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.back {
  color: var(--color-text-dim);
  display: inline-block;
  font-size: 13px;
  margin-bottom: 20px;
}

.back:hover {
  color: var(--color-text-bright);
}

.detail-head {
  border-bottom: 1px solid var(--color-border);
  margin-bottom: 24px;
  padding-bottom: 20px;
}

.head-id-row {
  align-items: center;
  display: flex;
  gap: 12px;
  margin-bottom: 10px;
}

.detail-id {
  color: var(--color-text-dim);
  font-size: 14px;
}

.detail-title {
  color: var(--color-text-bright);
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 14px;
}

.head-meta {
  color: var(--color-text);
  display: flex;
  font-size: 14px;
  gap: 28px;
}

.meta-key {
  color: var(--color-text-dim);
  font-size: 11px;
  letter-spacing: 1px;
  margin-right: 4px;
  text-transform: uppercase;
}

/* Two-column layout: main timeline + sidebar */
.detail-grid {
  display: flex;
  gap: 20px;
}

.detail-main {
  display: flex;
  flex: 1;
  flex-direction: column;
  gap: 20px;
}

@media (max-width: 768px) {
  /* Stack the timeline column above the sidebar */
  .detail-grid {
    flex-direction: column;
  }

  .detail-side {
    width: 100%;
  }

  .head-meta {
    flex-wrap: wrap;
    gap: 14px;
  }

  .head-id-row {
    flex-wrap: wrap;
  }

  .detail-title {
    font-size: 20px;
  }
}
</style>