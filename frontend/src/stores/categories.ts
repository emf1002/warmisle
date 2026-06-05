import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getCategories } from '@/api/category'

export const useCategoriesStore = defineStore('categories', () => {
  const categories = ref<any[]>([])
  const loaded = ref(false)

  async function fetchCategories() {
    if (loaded.value) return categories.value
    const res: any = await getCategories()
    categories.value = res.data || []
    loaded.value = true
    return categories.value
  }

  function reset() {
    categories.value = []
    loaded.value = false
  }

  return { categories, loaded, fetchCategories, reset }
})
