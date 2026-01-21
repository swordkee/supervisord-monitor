<template>
  <div class="server-card card shadow-sm h-100">
    <div class="card-header bg-white border-0 py-3">
      <div class="d-flex justify-content-between align-items-center">
        <div>
          <h5 class="mb-0 fw-bold">
            <a :href="server.url" class="text-decoration-none text-dark">{{ server.name }}</a>
            <small v-if="showHost" class="text-muted ms-2">({{ hostname }})</small>
            <i v-if="server.has_auth" class="bi bi-shield-lock text-primary ms-2" title="Authenticated server connection"></i>
            <small class="text-muted ms-2">{{ server.version }}</small>
          </h5>
        </div>
        <div v-if="!server.error" class="btn-group btn-group-sm" role="group">
          <button @click.prevent="stopAll" class="btn btn-outline-danger" type="button" title="Stop all">
            <i class="bi bi-stop-circle-fill"></i>
          </button>
          <button @click.prevent="startAll" class="btn btn-outline-success" type="button" title="Start all">
            <i class="bi bi-play-circle-fill"></i>
          </button>
          <button @click.prevent="restartAll" class="btn btn-outline-primary" type="button" title="Restart all">
            <i class="bi bi-arrow-clockwise"></i>
          </button>
        </div>
      </div>
    </div>

    <div class="card-body p-0">
      <div v-if="server.error" class="alert alert-danger m-3 rounded-3">
        <i class="bi bi-exclamation-triangle me-2"></i>{{ server.error }}
      </div>

      <table class="table table-hover table-sm mb-0">
        <tbody>
          <tr v-for="process in server.processes" :key="getProcessName(process)" class="align-middle">
            <td>
              <span class="fw-medium">{{ getProcessName(process) }}</span>
              <span v-if="process.has_error" class="float-end">
                <button
                  :id="`${server.name}_${getProcessName(process)}`"
                  @click.prevent="showError(process)"
                  class="btn btn-sm btn-danger rounded-pill"
                  title="View error">
                  <i class="bi bi-exclamation-triangle"></i>
                </button>
              </span>
            </td>
            <td style="width: 100px;">
              <span :class="['badge', `bg-${getStatusClass(process.statename)}`]">
                {{ process.statename.toUpperCase() }}
              </span>
            </td>
            <td style="width: 120px; text-align:right;">
              <small class="text-muted">{{ getUptime(process.description) }}</small>
            </td>
            <td style="width: 100px;">
              <div class="btn-group btn-group-sm" role="group">
                <template v-if="process.statename === 'RUNNING' || process.statename === 'Running'">
                  <button @click.prevent="stopProcess(process)" class="btn btn-outline-dark" type="button" title="Stop">
                    <i class="bi bi-stop-circle"></i>
                  </button>
                  <button @click.prevent="restartProcess(process)" class="btn btn-outline-dark" type="button" title="Restart">
                    <i class="bi bi-arrow-clockwise"></i>
                  </button>
                </template>
                <template v-else-if="['STOPPED', 'STOPPED', 'EXITED', 'Exited', 'FATAL', 'Fatal'].includes(process.statename)">
                  <button @click.prevent="startProcess(process)" class="btn btn-outline-success" type="button" title="Start">
                    <i class="bi bi-play-circle"></i>
                  </button>
                </template>
                <template v-else>
                  <button @click.prevent="startProcess(process)" class="btn btn-outline-success" type="button" title="Start">
                    <i class="bi bi-play-circle"></i>
                  </button>
                </template>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="showModal" class="modal fade show d-block" tabindex="-1" aria-labelledby="errorModalLabel" aria-hidden="true">
      <div class="modal-dialog modal-lg">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="errorModalLabel">
              <i class="bi bi-exclamation-triangle text-danger me-2"></i>{{ errorTitle }}
            </h5>
            <button type="button" class="btn-close" @click="closeErrorModal"></button>
          </div>
          <div class="modal-body">
            <div class="bg-light p-3 rounded-3">
              <pre class="mb-0 text-wrap" style="white-space:pre-wrap;word-wrap:break-word;">{{ errorMessage }}</pre>
            </div>
          </div>
          <div class="modal-footer">
            <button class="btn btn-secondary" @click="closeErrorModal">Close</button>
            <button class="btn btn-danger" @click="clearErrorLog">
              <i class="bi bi-trash me-1"></i>Clear Log
            </button>
          </div>
        </div>
      </div>
    </div>
    <div v-if="showModal" class="modal-backdrop fade show"></div>
  </div>
</template>

<script>
import { ref, computed } from 'vue'
import axios from 'axios'

export default {
  name: 'ServerCard',
  props: {
    server: {
      type: Object,
      required: true
    },
    showHost: {
      type: Boolean,
      default: false
    }
  },
  emits: ['refresh'],
  setup(props, { emit }) {
    const hostname = computed(() => {
      try {
        const url = new URL(props.server.url)
        return url.hostname
      } catch {
        return ''
      }
    })

    const errorTitle = ref('')
    const errorMessage = ref('')
    const currentProcess = ref(null)
    const showModal = ref(false)

    const getProcessName = (process) => {
      if (process.group !== process.name) {
        return `${process.group}:${process.name}`
      }
      return process.name
    }

    const getStatusClass = (status) => {
      switch (status) {
        case 'RUNNING':
        case 'Running':
          return 'success'
        case 'STARTING':
        case 'Starting':
          return 'info'
        case 'FATAL':
        case 'Fatal':
          return 'important'
        case 'STOPPED':
        case 'Stopped':
          return 'inverse'
        default:
          return 'error'
      }
    }

    const getUptime = (description) => {
      if (!description) return '&nbsp;'
      if (description === 'Fatal' || description === 'fatal' || description === 'STOPPED' || description === 'STOPPING') return '&nbsp;'
      const parts = description.split(',')
      if (parts.length >= 2) {
        const uptime = parts[1].trim()
        if (uptime.startsWith('uptime ')) {
          return uptime.replace('uptime ', '')
        }
        return uptime
      }
      return '&nbsp;'
    }

    const showError = (process) => {
      const processName = getProcessName(process)
      errorTitle.value = `${processName} @ ${props.server.name}`
      errorMessage.value = process.log || 'No error log available'
      currentProcess.value = process
      showModal.value = true
    }

    const closeErrorModal = () => {
      showModal.value = false
      currentProcess.value = null
    }

    const clearErrorLog = async () => {
      if (currentProcess.value) {
        const processName = getProcessName(currentProcess.value)
        try {
          await axios.post(`/api/clear/${props.server.name}/${processName}`)
          closeErrorModal()
          emit('refresh')
        } catch (error) {
          console.error('Clear log failed:', error)
          alert('Clear log failed: ' + (error.response?.data?.error || error.message))
        }
      }
    }

    const executeAction = async (action) => {
      try {
        await action()
        emit('refresh')
      } catch (error) {
        console.error('Action failed:', error)
        alert('Action failed: ' + (error.response?.data?.error || error.message))
      }
    }

    const startProcess = (process) => {
      const processName = getProcessName(process)
      executeAction(() => 
        axios.post(`/api/start/${props.server.name}/${processName}`)
      )
    }

    const stopProcess = (process) => {
      const processName = getProcessName(process)
      executeAction(() => 
        axios.post(`/api/stop/${props.server.name}/${processName}`)
      )
    }

    const restartProcess = (process) => {
      const processName = getProcessName(process)
      executeAction(() => 
        axios.post(`/api/restart/${props.server.name}/${processName}`)
      )
    }

    const clearLog = (process) => {
      const processName = getProcessName(process)
      executeAction(() => 
        axios.post(`/api/clear/${props.server.name}/${processName}`)
      )
    }

    const startAll = () => {
      executeAction(() => 
        axios.post(`/api/startall/${props.server.name}`)
      )
    }

    const stopAll = () => {
      executeAction(() => 
        axios.post(`/api/stopall/${props.server.name}`)
      )
    }

    const restartAll = () => {
      executeAction(() =>
        axios.post(`/api/restartall/${props.server.name}`)
      )
    }

    return {
      errorTitle,
      errorMessage,
      showModal,
      hostname,
      getProcessName,
      getStatusClass,
      getUptime,
      showError,
      closeErrorModal,
      clearErrorLog,
      startProcess,
      stopProcess,
      restartProcess,
      clearLog,
      startAll,
      stopAll,
      restartAll
    }
  }
}
</script>

<style scoped>
.server-card {
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 20px;
  overflow: hidden;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05), 0 10px 20px rgba(0, 0, 0, 0.08);
}

.server-card:hover {
  transform: translateY(-8px) scale(1.01);
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.15), 0 10px 20px rgba(0, 0, 0, 0.1) !important;
  border-color: rgba(99, 102, 241, 0.3);
}

.card-header {
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
  background: linear-gradient(180deg, rgba(255, 255, 255, 1) 0%, rgba(249, 250, 251, 1) 100%);
}

.card-header h5 a {
  transition: color 0.2s ease;
  color: #1f2937;
}

.card-header h5 a:hover {
  color: #6366f1 !important;
}

.btn-group .btn {
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  border-width: 1.5px;
}

.btn-group .btn:hover {
  transform: scale(1.1);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.btn-group .btn:active {
  transform: scale(0.95);
}

.table {
  margin-bottom: 0;
}

.table tbody tr {
  transition: all 0.2s ease;
  border-bottom: 1px solid rgba(0, 0, 0, 0.04);
}

.table tbody tr:last-child {
  border-bottom: none;
}

.table tbody tr:hover {
  background-color: rgba(99, 102, 241, 0.04);
  transform: translateX(4px);
}

.badge {
  font-size: 0.7rem;
  font-weight: 700;
  padding: 0.4em 0.7em;
  border-radius: 8px;
  letter-spacing: 0.5px;
  text-transform: uppercase;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.btn-sm {
  padding: 0.3rem 0.6rem;
  font-size: 0.8rem;
  border-radius: 10px;
  font-weight: 600;
}

.modal-content {
  border: none;
  border-radius: 20px;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.4);
  overflow: hidden;
}

.modal-header {
  border-bottom: 1px solid rgba(0, 0, 0, 0.08);
  padding: 1.5rem 2rem;
  background: linear-gradient(180deg, rgba(255, 255, 255, 1) 0%, rgba(249, 250, 251, 1) 100%);
}

.modal-footer {
  border-top: 1px solid rgba(0, 0, 0, 0.08);
  padding: 1.5rem 2rem;
  background: #f9fafb;
}

.modal-body {
  padding: 2rem;
}

.modal-body pre {
  font-size: 0.85rem;
  line-height: 1.7;
  color: #374151;
  background: #f3f4f6;
  padding: 1rem;
  border-radius: 8px;
  border: 1px solid rgba(0, 0, 0, 0.05);
}

.alert {
  border-radius: 14px;
  border: none;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

@media (max-width: 768px) {
  .server-card:hover {
    transform: none;
  }
}
</style>
