<script setup lang="ts">
import { ref } from 'vue';
import LoginForm from '../components/LoginForm.vue'
import { login } from '@/api.ts';

const error = ref('')
async function handleLogin({username, password} : {username:string; password:string})  {
    error.value = ''
    try {
      await login(username, password)
      window.location.href = "/incident-list"
    } catch (e) {
      error.value = (e as Error).message ?? 'Authentication failed'
    }
}
</script>

<template>
  <main>
    <LoginForm :error="error" @submit="handleLogin"/>
  </main>
</template>
