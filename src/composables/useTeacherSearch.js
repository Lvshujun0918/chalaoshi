import { ref, computed } from 'vue'
import data from '../../teacher.json'

const colleges = data.colleges || []
const allTeachers = (data.teachers || []).sort((a, b) => {
  const aRate = a.rate === 'N/A' ? -1 : parseFloat(a.rate)
  const bRate = b.rate === 'N/A' ? -1 : parseFloat(b.rate)
  if (bRate !== aRate) return bRate - aRate
  return b.hot - a.hot
})

const collegeMap = new Map(colleges.map(c => [c.id, c.name]))

export function useTeacherSearch() {
  const query = ref('')
  const loading = ref(false)

  const results = computed(() => {
    const q = query.value.trim().toLowerCase()
    if (!q) return []

    return allTeachers.filter(t =>
      t.name.includes(q) ||
      t.py.includes(q) ||
      t.sx.includes(q)
    )
  })

  const displayResults = computed(() => results.value.slice(0, 20))

  const totalTeachers = allTeachers.length
  const totalColleges = colleges.length

  function getCollegeName(xy) {
    return collegeMap.get(xy) || '未知学院'
  }

  function clearQuery() {
    query.value = ''
  }

  return {
    query,
    loading,
    results,
    displayResults,
    totalTeachers,
    totalColleges,
    getCollegeName,
    clearQuery,
  }
}
