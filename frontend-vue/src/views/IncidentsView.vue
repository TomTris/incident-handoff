<script setup lang="ts">
import { createIncident, isAuthenticated, loadIncidents, logout } from '@/api';
import AddIncidentForm from '@/components/AddIncidentForm.vue';
import AppHeader from '@/components/AppHeader.vue';
import IncidentList from '@/components/IncidentList.vue';
import IncidentListItem from '@/components/IncidentListItem.vue';
import type { Severity, CreateIncidentRequest, Incident } from '@/types';
import { onMounted, ref } from 'vue';
const incidents = ref<Incident[]>([])

async function handleLogout() {
  await logout()
  window.location.href = "/"
}

onMounted(async() => {
  if (await isAuthenticated() == false) {
    window.location.href = "/"
  }
  incidents.value = await loadIncidents()
})

const filterStatus = ref('All')
const service = ref('')

</script>

<template>
  <main>
    <AppHeader @submit="handleLogout"></AppHeader>
    <div class="page">
      <div class="dash-head">
        <div>
          <p class="eyebrow">Active board</p>
          <h1 class="page-title">Incidents</h1>
        </div>
        <div class="spacer"></div>
        <RouterLink :to="{name: 'incidents-new'}" class="btn btn-primary">+ New Incident</RouterLink>
      </div>

      <div class="filters">
        <div class="filter">
          <label class="field-label">Status</label>
          <select class="select" v-model="filterStatus">
            <option value="">All</option>
            <option value="active">Active</option>
            <option value="triggered">Triggered</option>
            <option value="investigating">investigating</option>
            <option value="mitigated">Mitigated</option>
            <option value="resolved">Resolved</option>
          </select>
        </div>
        <div class="filter">
          <label class="field-label">Service</label>
          <input class="input" type="text" v-model="service" placeholder="e.g Payment">
        </div>
      </div>

      <ul class="incident-list">
        <li v-for="inc in incidents">
          <IncidentListItem :key="inc.id" :inc="inc"></IncidentListItem>
        </li>
      </ul>



    </div>
  </main>
</template>

<style scoped>
.dash-head {
  align-items: flex-end;
  display: flex;
  margin-bottom: 20px;
}

.filters {
  display: flex;
  gap: 16px;
  list-style: none;
}

.filter {
  width: 220px;
}
</style>