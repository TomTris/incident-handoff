import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import type { UserContext } from '@/types'

// export const useUserContextStore = defineStore('userContext', () => {
//   const makeEmpty =() => {
//     return {id: '', username : '', role : ''}
//   }
//   const identity = ref<UserContext>(makeEmpty())
//   function assign(data : UserContext) {
//     identity.value = data
//   }
//   function clean() {
//     identity.value = makeEmpty()
//   }
//   return { identity, assign, clean }
// })
