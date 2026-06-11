<script setup lang="ts">
import { ref } from 'vue';

const props = defineProps<{error: string}>()
const emit = defineEmits(['submit'])

const title = ref('')
const service = ref('')
const severity = ref('')

function onSubmit() {
    emit("submit", {title: title.value, service: service.value, severity: severity.value})
}

</script>

<template>
    <form @submit.prevent="onSubmit">
        <p class="error">{{ props.error }}</p>
        <label>
            Title
            <input v-model="title" type="text">
        </label>
        <label>
            Service
            <input v-model="service" type="text">
        </label>
        <label>
            Severity
            <input v-model="severity" type="text">
        </label>
        <button type="submit">Create Incident</button>
    </form>
</template>

<style>
.error {color: red}
</style>