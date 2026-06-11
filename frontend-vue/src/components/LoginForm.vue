<script setup lang="ts">
import { ref } from 'vue';
defineProps({
    error: {type:String, default: ''}
})
const emit = defineEmits(['submit'])
const username = ref('')
const password = ref('')
function onSubmit() {
    emit('submit', {username: username.value, password: password.value})
}
</script>

<template>
    <form @submit.prevent="onSubmit">
        <p v-if="error" class="error" role="alert">{{ error }}</p>
        <label>
            Username
            <input v-model="username" type="test" required>
        </label>
        <label>
            Password
            <input v-model="password" type="password" required>
        </label>
        <button type="submit">Authenticate</button>
    </form>
</template>

<style scoped>
* {
  box-sizing: border-box;
}

body {
    margin: 0;
    padding: 1rem;
    font-family: system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
    display: flex;
    flex-direction: column;
    justify-content: center; /* Vertically */
    align-items: center; /* Horizontally */
    min-height: 100vh; /* percentage */
    background-color: #f5f5f5;
    /* color: #333; */
}

form {
    display: flex;
    flex-direction: column;
    gap: 1rem;  /* between elements inside */
    width: 100%;
    max-width: 30rem;
    padding: 2rem; /* from element to box border */
    border: 1px solid #ccc;
    border-radius: 0.5rem;
}

label {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    font-weight: 600;
}

input {
    padding: 0.5rem;
    border: 1px solid #ccc;
    border-radius: 0.25rem;
    font-size: 1rem;
}

button {
    padding: 0.75rem;
    background-color: #252525;
    color: white;
    border: none;
    border-radius: 0.25rem;
    font-size: 1rem;
    cursor: pointer;
}
button:hover {
    background-color: grey
}


h1 {
  margin-bottom: 2rem;
}
.error {color: red;}
</style>