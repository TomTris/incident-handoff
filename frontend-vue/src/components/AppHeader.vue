<script setup lang="ts">
import { logout, whoAmI } from '@/api';
import type { UserContext } from '@/types';
import { makeEmptyUserContext } from '@/utils/user';
import { onMounted, onUnmounted, ref } from 'vue';

async function handleLogout() {
    await logout()
    window.location.href = "/"
}

const identity = ref<UserContext>(makeEmptyUserContext())

const currentTime = ref('')
let timer: ReturnType<typeof setInterval> | undefined

function updateTime() {
    currentTime.value = new Date().toLocaleTimeString('de-DE', {
        timeZone: 'Europe/Berlin',
        hour: '2-digit',
        minute: '2-digit',
    })
}

onMounted(async () => {
    updateTime()
    // align to the next minute boundary, then tick every minute
    const msToNextMinute = 60000 - (Date.now() % 60000)
    timer = setTimeout(() => {
        updateTime()
        timer = setInterval(updateTime, 60000)
    }, msToNextMinute)

    try {
        identity.value = await whoAmI()
    } catch (e) {
        alert('Something is wrong.\n' + (e as Error).message)
        window.location.href = "/"
    }
})

onUnmounted(() => {
    clearTimeout(timer as ReturnType<typeof setTimeout>)
    clearInterval(timer as ReturnType<typeof setInterval>)
})

</script>

<template>
    <header class="appbar">
        <div class="appbar-inner">
            <RouterLink :to="{name: 'entry'}" class="brand">
                <span class="brand-name">HANDOFF</span>
                <span class="brand-mark">//</span>
            </RouterLink>
            <nav class="appbar-nav">
                <RouterLink :to="{name: 'incidents'}" class="nav-link">Incidents</RouterLink>
            </nav>


            <div class="spacer"></div>

            <span class="current-time mono">🕞 {{ currentTime }}</span>
            <div class="appbar-user">
                <span class="user-name mono"> {{ identity.username }}</span>
                <span class="user-role">{{identity.role}}</span>
            </div>
            <button class="btn appbar-logout" @click="handleLogout">Log out</button>
        </div>
    </header>
</template>

<style scoped>
.appbar {
  background-color: var(--color-panel);
  border-bottom: 1px solid var(--color-border);
}

.appbar-inner {
  align-items: center;
  display: flex;
  gap: 24px;
  margin: 0 auto;
  max-width: 1100px;
  padding: 0 24px;
}

.brand {
  align-items: center;
  color: var(--color-text-bright);
  display: flex;
  gap: 5px;
  height: 56px;
}

.brand-mark {
  color: var(--color-accent);
  font-size: 18px;
}

.brand-name {
  color: var(--color-text-bright);
  font-family: var(--font-mono);
  font-size: 16px;
  font-weight: 700;
  letter-spacing: 4px;
}

.appbar-nav {
  display: flex;
  gap: 18px;
}

.nav-link {
  color: var(--color-text-dim);
  font-size: 14px;
  font-weight: 600;
  transition: color 0.15s;
}

.nav-link:hover {
  color: var(--color-text-bright);
}

.router-link-exact-active.nav-link {
  color: var(--color-accent);
}

.current-time {
  align-items: center;
  color: var(--color-text-dim);
  display: flex;
  font-size: 13px;
  gap: 6px;
}

.clock-icon {
  color: var(--color-accent);
  flex-shrink: 0;
}

.appbar-user {
  display: flex;
  flex-direction: column;
  line-height: 1.2;
  text-align: right;
}

.user-name {
  color: var(--color-text-bright);
  font-size: 13px;
}

.user-role {
  color: var(--color-text-dim);
  font-size: 11px;
  letter-spacing: 1px;
  text-transform: uppercase;
}

.appbar-logout {
  padding: 7px 14px;
}

.router-link-exact-active.nav-link {
  color: var(--color-accent);
}

@media (max-width: 768px) {
  .appbar-inner {
    gap: 12px;
    padding: 0 16px;
  }

  .appbar-nav {
    display: none;
  }
}

</style>