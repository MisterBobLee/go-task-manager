<template>
  <div>
    <header class="header">
      <div class="header-left">
        {{ currentTime }}
      </div>

      <div class="header-right">
        <span class="username">{{ username }}</span>
        <button @click="logout" class="logout-button">
          登出
        </button>
      </div>
    </header>

    <div class="main">
      <div class="container">
        <h1 class="title">任務清單</h1>

        <div class="toolbar">
          <div class="input-wrapper">
            <template v-if="isCreatingTask">
              <input
                v-model="newTaskTitle"
                type="text"
                placeholder="輸入任務標題..."
                class="task-input"
              />
            </template>
          </div>

          <div class="button-group">
            <template v-if="isCreatingTask">
              <button @click="confirmCreateTask" class="btn create-btn">
                確認
              </button>
              <button @click="cancelCreateTask" class="btn cancel-btn">
                取消
              </button>
            </template>
            <template v-else>
              <button @click="startCreateTask" class="btn create-btn">
                新增任務
              </button>
              <button @click="toggleEditMode" class="btn edit-btn">
                {{ isEditMode ? '取消' : '編輯' }}
              </button>
              <button
                v-if="isEditMode"
                @click="deleteSelectedTasks"
                class="btn delete-btn"
              >
                刪除已選任務
              </button>
            </template>
          </div>
        </div>

        <div class="page-size-selector">
          <label class="item-numbers">每頁顯示筆數：</label>
          <select v-model.number="size" @change="updateSize" class="page-select">
            <option :value="10">10</option>
            <option :value="20">20</option>
            <option :value="50">50</option>
            <option :value="100">100</option>
          </select>
        </div>

        <div v-if="error" class="error">{{ error }}</div>
        <div v-else-if="tasks.length === 0" class="no-task">目前沒有任務</div>

        <ul class="task-list">
          <li
            v-for="task in tasks"
            :key="task.id"
            class="task-card"
          >
            <div class="task-content">
              <input
                type="checkbox"
                :checked="isEditMode ? selectedTaskIds.has(task.id) : task.completed"
                @change="isEditMode ? toggleSelectTask(task.id) : toggleCompleted(task)"
                :class="['task-checkbox', { 'edit-mode-checkbox': isEditMode }]"
              />
              <div>
                <h2 class="task-title" :class="{ 'task-completed': task.completed }">
                  {{ task.title }}
                </h2>
                <span
                  class="task-status"
                  :class="task.completed ? 'completed' : 'pending'"
                >
                  {{ task.completed ? '已完成' : '未完成' }}
                </span>
              </div>
            </div>
          </li>
        </ul>
        <div class="pagination">
          <button @click="goToPage(1)" class="page-button">&laquo;</button>
          <button @click="goToPage(page - 1)" :disabled="page === 1" class="page-button">
            &lt;
          </button>
          <button
            v-for="p in visiblePages"
            :key="p"
            @click="goToPage(p)"
            :class="['page-button', { active: page === p }]"
          >
            {{ p }}
          </button>
          <span v-if="page + 2 < totalPages" class="ellipsis">...</span>
          <button
            v-if="page + 2 < totalPages"
            @click="goToPage(totalPages)"
            class="page-button"
          >
            {{ totalPages }}
          </button>
          <button @click="goToPage(page + 1)" :disabled="page === totalPages" class="page-button">
            &gt;
          </button>
          <button @click="goToPage(totalPages)" class="page-button">&raquo;</button>
          <div class="jump-to-page">
            <input
              v-model.number="inputPage"
              type="number"
              placeholder="跳至頁碼"
              class="jump-input"
              min="1"
              :max="totalPages"
            />
            <button @click="jumpToInputPage" class="btn jump-btn">跳轉</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watchEffect, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'

interface Task {
  id: number
  title: string
  completed: boolean
}

const tasks = ref<Task[]>([])
const error = ref('')
const route = useRoute()
const router = useRouter()

const page = ref(parseInt(route.query.page as string) || 1)
const size = ref(parseInt(route.query.size as string) || 10)
const totalPages = ref(1)
const hasMorePages = ref(false)
const newTaskTitle = ref('')
const isCreatingTask = ref(false)
const isEditMode = ref(false)
const selectedTaskIds = ref<Set<number>>(new Set())
const currentTime = ref('')
const username = ref('')

const fetchTasks = async () => {
  const token = localStorage.getItem('token')
  if (!token) {
    router.push('/login')
    return
  }

  try {
    const res = await axios.get('http://localhost:8080/api/tasks', {
      headers: {
        Authorization: `Bearer ${token}`
      },
      params: {
        page: page.value,
        size: size.value
      }
    })

    tasks.value = res.data.tasks
    totalPages.value = Math.ceil(res.data.total / size.value) || 1
    hasMorePages.value = page.value + 2 < totalPages.value
  } catch (err: any) {
    if (err.response?.status === 401) {
      localStorage.removeItem('token')
      router.replace('/login')
    } else {
      error.value = err.response?.data?.error || '任務讀取失敗'
    }
  }
}

watchEffect(() => {
  fetchTasks()
})

const updateTime = () => {
  const now = new Date()
  currentTime.value = now.toLocaleString()
}

onMounted(() => {
  updateTime()
  setInterval(updateTime, 1000)
})

username.value = localStorage.getItem('username') || '未登入使用者'

const toggleEditMode = () => {
  isEditMode.value = !isEditMode.value
  selectedTaskIds.value.clear()
}

const toggleSelectTask = (id: number) => {
  if (selectedTaskIds.value.has(id)) {
    selectedTaskIds.value.delete(id)
  } else {
    selectedTaskIds.value.add(id)
  }
}

const deleteSelectedTasks = async () => {
  if (!selectedTaskIds.value.size) return
  if (!confirm('確定要刪除選取的任務嗎？')) return

  const token = localStorage.getItem('token')
  if (!token) {
    router.push('/login')
    return
  }

  try {
    const ids = Array.from(selectedTaskIds.value)
    await Promise.all(ids.map(id =>
      axios.delete(`http://localhost:8080/api/tasks/${id}`, {
        headers: { Authorization: `Bearer ${token}` }
      })
    ))

    tasks.value = tasks.value.filter(t => !selectedTaskIds.value.has(t.id))
    selectedTaskIds.value.clear()
    isEditMode.value = false

  } catch (err: any) {
    alert(err.response?.data?.error || '刪除任務失敗')
  }
}

const startCreateTask = () => {
  isCreatingTask.value = true
  newTaskTitle.value = ''
}

const confirmCreateTask = async () => {
  if (!newTaskTitle.value.trim()) {
    alert('請輸入任務標題')
    return
  }

  const token = localStorage.getItem('token')
  if (!token) {
    router.push('/login')
    return
  }

  try {
    await axios.post(
      'http://localhost:8080/api/tasks',
      { title: newTaskTitle.value },
      {
        headers: {
          Authorization: `Bearer ${token}`
        }
      }
    )
    newTaskTitle.value = ''
    isCreatingTask.value = false
    fetchTasks()
  } catch (err: any) {
    alert(err.response?.data?.error || '新增任務失敗')
  }
}

const cancelCreateTask = () => {
  isCreatingTask.value = false
  newTaskTitle.value = ''
}

const toggleCompleted = async (task: Task) => {
  const token = localStorage.getItem('token')
  if (!token) {
    router.push('/login')
    return
  }

  try {
    await axios.put(
      `http://localhost:8080/api/tasks/${task.id}`,
      {
        completed: !task.completed
      },
      {
        headers: {
          Authorization: `Bearer ${token}`
        }
      }
    )

    const target = tasks.value.find(t => t.id === task.id)
    if (target) target.completed = !target.completed
  } catch (err: any) {
    alert(err.response?.data?.error || '更新任務狀態失敗')
  }
}

const logout = () => {
  if (confirm('確定要登出嗎？')) {
    localStorage.removeItem('token')
    router.replace('/login')
  }
}

const goToPage = (p: number) => {
  if (p < 1 || p > totalPages.value) return
  router.push({ path: '/tasks', query: { page: p.toString(), size: size.value.toString() } })
  page.value = p
}

const updateSize = () => {
  page.value = 1
  router.replace({
    path: '/tasks',
    query: { page: '1', size: size.value.toString() }
  })
}

const visiblePages = computed(() => {
  const pages: number[] = []
  const start = Math.max(1, page.value - 2)
  const end = Math.min(totalPages.value, page.value + 2)
  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  return pages
})

const inputPage = ref<number | null>(null)

const jumpToInputPage = () => {
  if (!inputPage.value || inputPage.value < 1 || inputPage.value > totalPages.value) {
    alert(`請輸入 1 到 ${totalPages.value} 之間的頁碼`)
    return
  }

  goToPage(inputPage.value)
  inputPage.value = null
}
</script>
