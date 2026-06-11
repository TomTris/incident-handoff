<script setup lang="ts">
import { createIncident, isAuthenticated, loadIncidents, logout } from '@/api';
import AddIncidentForm from '@/components/AddIncidentForm.vue';
import IncidentList from '@/components/IncidentList.vue';
import type { Severity, CreateIncidentRequest, Incident } from '@/types';
import { ref } from 'vue';

const createIncidentFormError = ref('')
const incidents = ref<Incident[]>([])

async function handleLogout() {
  await logout()
  window.location.href = "/"
}


async function handleCreateIncident({title, service, severity}: {title: string, service: string, severity: string}) {
  function isSeverity(v: string) : v is Severity {
    return v === "SEV1" || v === "SEV2" || v === "SEV3"
  }

  if (!isSeverity(severity)) {
    createIncidentFormError.value = "Severity must be SEV1, SEV2 or SEV3"
    return
  }

  const incReq : CreateIncidentRequest = {
    title: title,
    service: service,
    severity: severity
  }

  try {
    const res = await createIncident(incReq)
    createIncidentFormError.value = ''
    incidents.value.push(res)
  } catch (e) {
    createIncidentFormError.value = (e as Error).message
  }
}

async function init() {
  if (await isAuthenticated() == false) {
    window.location.href = "/"
  }
  incidents.value = await loadIncidents()
}

init()
</script>

<template>
  <main>
    <button type="button" @click="handleLogout">Logout</button>
    <AddIncidentForm :error="createIncidentFormError" @submit="handleCreateIncident"/>
    <IncidentList :incidents="incidents"/>
  </main>
</template>
