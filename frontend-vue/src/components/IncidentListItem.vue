<script setup lang="ts">
import type { Incident } from '@/types';

const prop = defineProps<{inc : Incident}>()
</script>

<template>
    <h3>title: {{ inc.title }}</h3>
    <p>service: {{ inc.service }}</p>
    <p>severity: {{ inc.severity }}</p>
    <p>status: {{ inc.status }}</p>
    <p>opened_by: {{ inc.opened_by }}</p>
    <p>on_call: {{ inc.on_call }}</p>
    <p>created_at: {{ inc.created_at }}</p>
    <p>updated_at: {{ inc.updated_at }}</p>
    <p>version: {{ inc.version }}</p>
    <p>.</p>
    <p>.</p>
</template>