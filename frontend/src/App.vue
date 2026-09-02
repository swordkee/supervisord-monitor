<template>
  <div id="app">
    <nav class="navbar navbar-expand-lg navbar-dark bg-dark shadow-lg fixed-top">
      <div class="container">
        <a class="navbar-brand fw-bold" href="/">
          <i class="bi bi-speedometer2 me-2"></i>Supervisord Monitor
        </a>
        <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav">
          <span class="navbar-toggler-icon"></span>
        </button>
        <div class="collapse navbar-collapse" id="navbarNav">
          <ul class="navbar-nav ms-auto">
            <li class="nav-item">
              <a class="nav-link active" href="/">Home</a>
            </li>
            <li class="nav-item" v-if="currentUser">
              <span class="nav-link">
                <i class="bi bi-person-circle me-1"></i>{{ currentUser }}
              </span>
            </li>
            <li class="nav-item">
              <a class="nav-link" href="#" @click="toggleMute">
                <i :class="['bi', muted ? 'bi-volume-mute' : 'bi-volume-up']"></i>
                {{ muted ? 'Unmute' : 'Mute' }}
              </a>
            </li>
            <li class="nav-item">
              <a class="nav-link" href="#" @click="refresh">
                <i class="bi bi-arrow-clockwise me-1"></i>
                Refresh
                <span class="badge bg-light text-primary ms-1">{{ refreshTimer }}</span>
              </a>
            </li>
            <li class="nav-item">
              <a class="nav-link" href="mailto:martin@lazarov.bg">Contact</a>
            </li>
          </ul>
        </div>
      </div>
    </nav>

    <div class="container py-4">
      <div v-if="muted" class="alert alert-warning d-flex justify-content-between align-items-center shadow-sm" style="margin-bottom: 20px;">
        <span><i class="bi bi-volume-mute me-2"></i>Sound muted</span>
        <a href="#" @click="toggleMute" class="btn btn-sm btn-warning">
          <i class="bi bi-volume-up me-1"></i> Unmute
        </a>
      </div>

      <div v-if="currentUser && servers.length === 0" class="alert alert-info shadow-sm">
        <i class="bi bi-info-circle me-2"></i>当前账号没有可访问的监控服务器，请联系管理员分配监控组。
      </div>

      <div class="row g-4">
        <div v-for="server in servers" :key="server.name"
             :class="supervisorCols === 2 ? 'col-lg-6 col-md-12' : 'col-lg-4 col-md-6'">
          <ServerCard :server="server" :show-host="showHost" @refresh="refresh" />
        </div>
      </div>
    </div>

    <div class="footer">
      <p>Powered by <a href="https://github.com/swordkee/supervisord-monitor" target="_blank">Supervisord Monitor</a></p>
    </div>

    <audio ref="alertSound" src="/sounds/alert.mp3"></audio>
  </div>
</template>

<script>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import axios from 'axios'
import ServerCard from './components/ServerCard.vue'

export default {
  name: 'App',
  components: { ServerCard },
  setup() {
    const servers = ref([])
    const muted = ref(false)
    const refreshTimer = ref(10)
    const refreshInterval = ref(null)
    const alertSound = ref(null)
    const supervisorCols = ref(2)
    const showHost = ref(false)
    const refreshIntervalSeconds = ref(10)
    const enableAlarm = ref(true)
    const hasUserInteracted = ref(false)
    const currentUser = ref('')

    const hasAlert = computed(() => {
      return servers.value.some(server =>
        server.processes?.some(proc => proc.has_error || proc.statename === 'FATAL')
      )
    })

    const fetchDashboard = async () => {
      try {
        const response = await axios.get('/api/dashboard')
        servers.value = response.data.servers
        supervisorCols.value = response.data.supervisor_cols
        showHost.value = response.data.show_host
        enableAlarm.value = response.data.enable_alarm
        currentUser.value = response.data.user || ''
        refreshTimer.value = 0  // 不自动刷新

        if (hasAlert.value && !muted.value && enableAlarm.value && hasUserInteracted.value) {
          playAlert()
        }

        if (hasAlert.value) {
          document.title = '!!! WARNING !!!'
        } else {
          document.title = 'Supervisord Monitoring'
        }
      } catch (error) {
        console.error('Failed to fetch dashboard:', error)
      }
    }

    const playAlert = () => {
      if (alertSound.value && hasUserInteracted.value) {
        alertSound.value.play().catch(e => console.log('Audio play failed:', e))
      }
    }

    const toggleMute = () => {
      muted.value = !muted.value
      localStorage.setItem('muted', muted.value ? 'true' : 'false')
    }

    const startRefreshTimer = () => {
      if (refreshInterval.value) {
        clearInterval(refreshInterval.value)
      }
      if (refreshIntervalSeconds.value > 0) {
        refreshInterval.value = setInterval(() => {
          refreshTimer.value--
          if (refreshTimer.value <= 0) {
            fetchDashboard()
          }
        }, 1000)
      }
    }

    const refresh = () => {
      fetchDashboard()
    }

    const handleUserInteraction = () => {
      if (!hasUserInteracted.value) {
        hasUserInteracted.value = true
      }
    }

    onMounted(() => {
      const savedMuted = localStorage.getItem('muted')
      if (savedMuted === 'true') {
        muted.value = true
      }

      document.addEventListener('click', handleUserInteraction, { once: true })
      document.addEventListener('keydown', handleUserInteraction, { once: true })
      document.addEventListener('touchstart', handleUserInteraction, { once: true })

      fetchDashboard()
      // 不自动刷新
    })

    onUnmounted(() => {
      if (refreshInterval.value) {
        clearInterval(refreshInterval.value)
      }
      document.removeEventListener('click', handleUserInteraction)
      document.removeEventListener('keydown', handleUserInteraction)
      document.removeEventListener('touchstart', handleUserInteraction)
    })

    return {
      servers,
      muted,
      refreshTimer,
      alertSound,
      supervisorCols,
      showHost,
      currentUser,
      refresh,
      toggleMute
    }
  }
}
</script>

<style scoped>
#app {
  min-height: 100vh;
}

.navbar {
  background: linear-gradient(135deg, #1f2937 0%, #111827 100%) !important;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.navbar-brand {
  font-size: 1.5rem;
  letter-spacing: 0.5px;
  font-weight: 700;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.navbar-brand:hover {
  background: linear-gradient(135deg, #818cf8 0%, #a78bfa 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.nav-link {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  font-weight: 500;
  padding: 0.5rem 1rem !important;
  border-radius: 8px;
  margin: 0 0.2rem;
}

.nav-link:hover {
  transform: translateY(-2px);
  background: rgba(99, 102, 241, 0.15);
}

.nav-link.active {
  background: rgba(99, 102, 241, 0.2);
  color: #818cf8 !important;
}

.navbar-toggler {
  border: 2px solid rgba(255, 255, 255, 0.2);
  padding: 0.5rem 0.75rem;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.navbar-toggler:hover {
  border-color: rgba(99, 102, 241, 0.5);
  background: rgba(99, 102, 241, 0.1);
}

.container {
  max-width: 1600px;
}

.alert {
  border-radius: 14px;
  border: none;
  animation: slideDown 0.3s ease-out;
  background: rgba(255, 193, 7, 0.15);
  color: #ffc107;
  backdrop-filter: blur(10px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.badge {
  font-weight: 600;
  padding: 0.35em 0.65em;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.footer {
  margin-top: 40px;
  padding: 30px 0;
  background: rgba(0, 0, 0, 0.4);
  backdrop-filter: blur(20px);
  border-radius: 16px;
  margin-bottom: 20px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.2);
}

.footer p {
  margin: 0;
  font-size: 0.9rem;
  color: rgba(255, 255, 255, 0.7);
}

.footer a {
  font-weight: 500;
}

.footer a:hover {
  color: #818cf8 !important;
}
</style>
