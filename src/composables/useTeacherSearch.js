import { ref, computed } from 'vue'

const API_BASE = '/api'

const collegeMap = new Map()
let departmentsCache = null

export function useTeacherSearch() {
  const query = ref('')
  const loading = ref(false)
  const results = ref([])
  const departments = ref([])
  const totalTeachers = ref(0)
  const totalColleges = ref(0)
  const currentDepartment = ref('')
  const totalResults = ref(0)
  const page = ref(1)
  const totalPages = ref(1)
  const sortBy = ref('rating')
  const sortOrder = ref('desc')

  // 加载统计数据
  async function loadStats() {
    try {
      const res = await fetch(`${API_BASE}/stats`)
      const data = await res.json()
      totalTeachers.value = data.total_teachers
      totalColleges.value = data.total_departments
    } catch (e) {
      console.error('加载统计失败:', e)
    }
  }

  // 加载院系列表
  async function loadDepartments() {
    if (departmentsCache) {
      departments.value = departmentsCache
      return
    }
    try {
      const res = await fetch(`${API_BASE}/departments`)
      const data = await res.json()
      departmentsCache = data
      departments.value = data
      data.forEach(d => {
        collegeMap.set(d.name, d.name)
      })
    } catch (e) {
      console.error('加载院系失败:', e)
    }
  }

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
      if (currentDepartment.value) {
        params.set('department', currentDepartment.value)
      }
      const res = await fetch(`${API_BASE}/teachers?${params}`)
      const data = await res.json()
      results.value = data.teachers
      totalResults.value = data.total
      page.value = data.page
      totalPages.value = data.total_pages
    } catch (e) {
      console.error('搜索失败:', e)
    } finally {
      loading.value = false
    }
  }

  // 显示结果（已分页）
  const displayResults = computed(() => results.value)

  function getCollegeName(name) {
    return name || '未知学院'
  }

  function clearQuery() {
    query.value = ''
    currentDepartment.value = ''
    results.value = []
    totalResults.value = 0
  }

  function filterByDepartment(dept) {
    currentDepartment.value = dept
    search(1)
  }

  function setPage(pageNum) {
    search(pageNum)
  }

  function setSort(sort) {
    sortBy.value = sort
    search(1)
  }

  function setSortOrder(order) {
    sortOrder.value = order
    search(1)
  }

  // 初始化
  loadStats()
  loadDepartments()

  return {
    query,
    loading,
    results,
    displayResults,
    totalTeachers,
    totalColleges,
    totalResults,
    page,
    totalPages,
    departments,
    currentDepartment,
    sortBy,
    sortOrder,
    getCollegeName,
    clearQuery,
    search,
    filterByDepartment,
    setPage,
    setSort,
    setSortOrder,
  }
}

