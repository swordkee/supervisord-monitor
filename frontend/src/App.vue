<template>
  <div id="app">
    <div class="navbar navbar-inverse navbar-fixed-top">
      <div class="navbar-inner">
        <div class="container">
          <button type="button" class="btn btn-navbar" data-toggle="collapse" data-target=".nav-collapse">
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
          </button>
          <a class="brand" href="/">Support Center</a>
          <div class="nav-collapse collapse">
            <ul class="nav">
              <li class="active"><a href="/">Home</a></li>
              <li><a href="#" @click="toggleMute">
                <span :class="['icon', 'icon-music', 'icon-white']"></span>
                &nbsp;{{ muted ? 'Unmute' : 'Mute' }}
              </a></li>
              <li>
                <a href="#" @click="refresh">
                  Refresh
                  <b :id="'refresh'">({{ refreshTimer }})</b>
                </a>
              </li>
              <li><a href="mailto:martin@lazarov.bg">Contact</a></li>
            </ul>
          </div>
        </div>
      </div>
    </div>

    <div class="container">
      <div v-if="muted" class="row">
        <div class="span4 offset4 label label-important" style="padding:10px;margin-bottom:20px;text-align:center;">
          Sound muted
          <span class="pull-right">
            <a href="#" @click="toggleMute" style="color:white;">
              <span class="icon icon-music icon-white"></span> Unmute
            </a>
          </span>
        </div>
      </div>

      <div class="row">
        <div v-for="server in servers" :key="server.name"
             :class="['span', supervisorCols === 2 ? '6' : '4']">
          <ServerCard :server="server" :show-host="showHost" @refresh="refresh" />
        </div>
      </div>
    </div>

    <div class="footer">
      <p>Powered by <a href="https://github.com/mlazarov/supervisord-monitor" target="_blank">Supervisord Monitor</a> | Page rendered in <strong>0.05</strong> seconds</p>
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
        refreshIntervalSeconds.value = response.data.refresh
        enableAlarm.value = response.data.enable_alarm
        refreshTimer.value = refreshIntervalSeconds.value
        muted.value = response.data.muted

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
      startRefreshTimer()
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
      refresh,
      toggleMute
    }
  }
}
</script>

<style scoped>
</style>
