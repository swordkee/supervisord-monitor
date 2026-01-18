<template>
  <div class="server-card">
    <table class="table table-bordered table-condensed table-striped">
      <tr>
        <th colspan="4">
          <a :href="server.url">{{ server.name }}</a>
          <i v-if="showHost">{{ hostname }}</i>
          <i v-if="server.has_auth" class="icon-lock icon-green" style="color:blue" title="Authenticated server connection"></i>
          &nbsp;<i>{{ server.version }}</i>
          <span v-if="!server.error" class="server-btns pull-right">
            <a href="#" @click.prevent="stopAll" class="btn btn-mini btn-inverse" type="button">
              <i class="icon-stop icon-white"></i> Stop all
            </a>
            <a href="#" @click.prevent="startAll" class="btn btn-mini btn-success" type="button">
              <i class="icon-play icon-white"></i> Start all
            </a>
            <a href="#" @click.prevent="restartAll" class="btn btn-mini btn-primary" type="button">
              <i class="icon icon-refresh icon-white"></i> Restart all
            </a>
          </span>
        </th>
      </tr>
      <tr v-if="server.error">
        <td colspan="4">{{ server.error }}</td>
      </tr>
      <tr v-for="process in server.processes" :key="getProcessName(process)">
        <td>
          {{ getProcessName(process) }}
          <span v-if="process.has_error" class="pull-right">
            <a href="#"
               :id="`${server.name}_${getProcessName(process)}`"
               @click.prevent="showError(process)"
               class="pop btn btn-mini btn-danger">
              <img src="/img/alert_icon.png" />
            </a>
          </span>
        </td>
        <td width="10">
          <span :class="['label', `label-${getStatusClass(process.statename)}`]">
            {{ process.statename.toUpperCase() }}
          </span>
        </td>
        <td width="80" style="text-align:right">
          {{ getUptime(process.description) }}
        </td>
        <td style="width:1%">
          <div class="actions">
            <template v-if="process.statename === 'RUNNING' || process.statename === 'Running'">
              <a href="#" @click.prevent="stopProcess(process)" class="btn btn-mini btn-inverse" type="button" title="Stop">
                <i class="icon-stop icon-white"></i>
              </a>
              <a href="#" @click.prevent="restartProcess(process)" class="btn btn-mini btn-inverse" type="button" title="Restart">
                <i class="icon-refresh icon-white"></i>
              </a>
            </template>
            <template v-else-if="['STOPPED', 'STOPPED', 'EXITED', 'Exited', 'FATAL', 'Fatal'].includes(process.statename)">
              <a href="#" @click.prevent="startProcess(process)" class="btn btn-mini btn-success" type="button" title="Start">
                <i class="icon-play icon-white"></i>
              </a>
            </template>
            <template v-else>
              <a href="#" @click.prevent="startProcess(process)" class="btn btn-mini btn-success" type="button" title="Start">
                <i class="icon-play icon-white"></i>
              </a>
            </template>
          </div>
        </td>
      </tr>
    </table>

    <div class="modal hide fade" id="errorModal" tabindex="-1" role="dialog" aria-labelledby="errorModalLabel" aria-hidden="true" ref="errorModal">
      <div class="modal-header">
        <button type="button" class="close" data-dismiss="modal" aria-hidden="true" @click="closeErrorModal">×</button>
        <h3 id="errorModalLabel">{{ errorTitle }}</h3>
      </div>
      <div class="modal-body">
        <div class="well" style="padding:20px;">
          <pre style="white-space:pre-wrap;word-wrap:break-word;">{{ errorMessage }}</pre>
        </div>
      </div>
      <div class="modal-footer">
        <button class="btn btn-primary" @click="closeErrorModal">ok</button>
        <button class="btn btn-danger" @click="clearErrorLog">Clear</button>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, onUnmounted } from 'vue'
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

    const errorModal = ref(null)
    const errorTitle = ref('')
    const errorMessage = ref('')
    const currentProcess = ref(null)
    let $modal = null

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
      currentProcess.value = process
      const processName = getProcessName(process)
      errorTitle.value = `${processName}@${props.server.name}`
      errorMessage.value = process.log || ''
      if ($modal) {
        $modal.modal('show')
      }
    }

    const closeErrorModal = () => {
      if ($modal) {
        $modal.modal('hide')
      }
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

    onMounted(() => {
      if (errorModal.value) {
        $modal = window.$(errorModal.value)
        $modal.modal({ show: false })
      }
    })

    onUnmounted(() => {
      if ($modal) {
        $modal.modal('hide')
        $modal.remove()
      }
    })

    return {
      errorModal,
      errorTitle,
      errorMessage,
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
.server-btns {
  margin-top: -5px;
}

.server-btns a:not(:last-child) {
  margin-right: 3px;
}

.table-condensed {
  font-size: 13px;
  font-weight: normal;
}

.table-condensed td, .table-condensed th {
  font-weight: normal;
}

.btn-mini {
  padding: 4px 8px;
  font-size: 11px;
  line-height: 14px;
  font-weight: bold;
  min-width: 20px;
  display: inline-block;
}

.icon-white {
  color: white !important;
}

.label {
  font-size: 11px;
  font-weight: bold;
  padding: 3px 6px;
}

/* Bootstrap 2.3 icon fallback */
[class^="icon-"],
[class*=" icon-"] {
  display: inline-block;
  width: 14px;
  height: 14px;
  line-height: 14px;
  vertical-align: text-top;
  background-image: url("https://cdn.bootcdn.net/ajax/libs/twitter-bootstrap/2.3.2/img/glyphicons-halflings.png");
  background-position: 14px 14px;
  background-repeat: no-repeat;
}

.icon-white {
  background-image: url("https://cdn.bootcdn.net/ajax/libs/twitter-bootstrap/2.3.2/img/glyphicons-halflings-white.png");
}

.icon-stop {
  background-position: -312px 0;
}

.icon-play {
  background-position: -264px 0;
}

.icon-refresh {
  background-position: -240px 0;
}

.icon-lock {
  background-position: -24px 0;
}
</style>
