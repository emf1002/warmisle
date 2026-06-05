import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getMembers } from '@/api/member'

export const useMembersStore = defineStore('members', () => {
  const members = ref<any[]>([])
  const loaded = ref(false)

  async function fetchMembers() {
    if (loaded.value) return members.value
    const res: any = await getMembers()
    members.value = res.data || []
    loaded.value = true
    return members.value
  }

  function reset() {
    members.value = []
    loaded.value = false
  }

  return { members, loaded, fetchMembers, reset }
})
