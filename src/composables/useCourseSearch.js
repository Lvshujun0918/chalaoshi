import { ref, computed } from 'vue'

const API_BASE = '/api'

export function useCourseSearch() {
  const query = ref('')
  const loading = ref(false)
  const results = ref([])
  const totalResults = ref(0)
  const page = ref(1)
  const totalPages = ref(1)
  const sortBy = ref('gpa')
  const sortOrder = ref('desc')

  // 搜索
  async function search(pageNum = 1) {
    loading.value = true
    try {
      const params = new URLSearchParams({
        page: String(pageNum),
        page_size: '20',
        sort_by: sortBy.value,
        sort_order: sortOrder.value,
      })
      const q = query.value.trim()
      if (q) params.set('q', q)
      const res = await fetch(`${API_BASE}/courses?${params}`)
      const data = await res.json()
      results.value = data.courses
      totalResults.value = data.total
      page.value = data.page
      totalPages.value = data.total_pages
    } catch (e) {
      console.error('课程搜索失败:', e)
    } finally {
      loading.value = false
    }
  }

  // 显示结果（已分页）
  const displayResults = computed(() => results.value)

  function clearQuery() {
    query.value = ''
    results.value = []
    totalResults.value = 0
  }

  function setPage(pageNum) {
    search(pageNum)
  }

  function setSort(sort) {
    sortBy.value = sort
    sortOrder.value = 'desc'
    search(1)
  }

  function setSortOrder(order) {
    sortOrder.value = order
    search(1)
  }

  return {
    query,
    loading,
    results,
    displayResults,
    totalResults,
    page,
    totalPages,
    sortBy,
    sortOrder,
    clearQuery,
    search,
    setPage,
    setSort,
    setSortOrder,
  }
}
