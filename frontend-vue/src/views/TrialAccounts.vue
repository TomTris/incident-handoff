<script setup lang="ts">
import { ref } from 'vue';
import { login, whoAmI } from '@/api.ts';
import { useUserContextStore } from '@/stores/userIdentity';

const error = ref('')
const username = ref('')
const password = ref('')
async function handleLogin()  {
    error.value = ''
    try {
      await login(username.value, password.value)
      window.location.href = "/incidents"
    } catch (e) {
      error.value = (e as Error).message ?? 'Authentication failed'
    }
}
</script>

<template>
  <main>
    
    <div class="login-screen">
        <div class ="login-card">
            <div class="login-brand">
                <span class="login-wordmark">HANDOFF</span>
                <span class="login-mark">\\</span>
            </div>
            <p class="login-tag"> Trial accounts</p>
            <div class="accounts">
                <div class="account">
                    <p class="account-info"> username: anh </p>
                    <p class="account-info"> password: anh123 </p>
                    <p class="account-info"> role: engineer </p>
                </div>
                <div class="account">
                    <p class="account-info"> username: bernd </p>
                    <p class="account-info"> password: bernd123 </p>
                    <p class="account-info"> role: engineer </p>
                </div>
                <div class="account">
                    <p class="account-info"> username: admin </p>
                    <p class="account-info"> password: admin123 </p>
                    <p class="account-info"> role: admin </p>
                </div>
            </div>
            <p class="login-foot mono dim">ON-CALL ACCESS ONLY</p>
            <RouterLink :to="{name:'incidents'}" class="back mono">← Back to incident</RouterLink>
        </div>
    </div>
  </main>
</template>

<style scoped>
.back {
        padding-top: 30px;
    display: flex;
    justify-content: center;
    font-size: 18px;
}
.accounts {
    display: flex;
    flex-direction: column;
    gap: 20px;
    padding-left: 30px;
    margin-bottom: 40px;
}

.login-screen {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    min-height: 90vh;
    padding: 24px;
}

.login-card {
    background-color: var(--color-panel);
    border: 1px solid var(--color-border);
    border-top: 3px solid var(--color-accent);
    border-radius: 8px;
    padding: 40px;
    max-width: 380px;
    width: 100%;
}

.login-brand{
    display: flex;
    flex-direction: row;
    justify-content: center;
    align-items: center;
    gap: 5px;
}

.login-wordmark {
    color: var(--color-text-bright);
    font-family: var(--font-mono);
    font-size: 26px;
    font-weight: 700;
    letter-spacing: 4px;
}

.login-tag {
    color: var(--color-text-dim);
    font-size: 13px;
    margin: 12px 0 28px;
    text-align: center;
}

.login-mark {
    color: var(--color-accent);
    font-size: 26px;
}

.login-form {
    margin-bottom: 20px;
}

.login-foot{
    font-size: 10px;
    letter-spacing: 2px;
    text-align: center;
}
</style>