<template>
  <div class="shared-doc-view">
    <div v-if="isVerified" class="doc-container">
      <div class="header">
        <div class="title">
          <el-icon><Document /></el-icon>
          <h2>{{fileName}}</h2>
        </div>
        <div class="collaborators">
          <div v-for="client in activeClients" :key="client.id" class="collaborator">
            <el-tooltip :content="client.name" placement="top">
              <div class="avatar" :style="{ backgroundColor: getClientColor(client.id) }">
                {{ client.name[0] }}
              </div>
            </el-tooltip>
          </div>
        </div>
        <div class="connected-users">
          <el-tooltip content="Connected Users" placement="top">
            <span>{{ connectedUserCount }} Users Connected</span>
          </el-tooltip>
        </div>
      </div>
      <div class="content">
        <div id="vditor" class="vditor"></div>
      </div>
    </div>

    <!-- 密码验证对话框 -->
    <el-dialog
      v-model="showPasswordDialog"
      title="访问验证"
      width="400px"
      :close-on-click-modal="false"
      :show-close="false"
      :close-on-press-escape="false"
    >
      <div class="password-dialog-content">
        <el-icon class="dialog-icon"><Lock /></el-icon>
        <div class="dialog-title">请输入访问密码</div>
        <el-input
          v-model="inputPassword"
          type="password"
          placeholder="请输入密码"
          show-password
          @keyup.enter="verifyPassword"
        />
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button type="primary" @click="verifyPassword" :loading="verifying">确认</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onBeforeUnmount, computed, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Document, Lock } from '@element-plus/icons-vue'
import Vditor from 'vditor'
import 'vditor/dist/index.css'
import * as Y from 'yjs'
import { WebsocketProvider } from 'y-websocket'

const route = useRoute()
const fileName = ref('未命名文档')
const vditor = ref(null)
const showPasswordDialog = ref(true)
const inputPassword = ref('')
const isVerified = ref(false)
const verifying = ref(false)

// Yjs 相关变量
const ydoc = ref(null)
const provider = ref(null)
const ytext = ref(null)

// 计算活跃的客户端（排除自己）
const activeClients = computed(() => {
  if (!provider.value?.awareness?.getStates()) return []
  const states = Array.from(provider.value.awareness.getStates().entries())
  const currentClientId = ydoc.value?.clientID
  return states
    .filter(([clientId]) => clientId !== currentClientId)
    .map(([clientId, state]) => ({
      id: clientId,
      name: state.user?.name || '匿名用户'
    }))
})

// 生成客户端颜色
const getClientColor = (id) => {
  const colors = ['#f56c6c', '#e6a23c', '#67c23a', '#409eff', '#909399']
  const index = Math.abs(hashCode(id)) % colors.length
  return colors[index]
}

// 简单的字符串哈希函数
const hashCode = (str) => {
  let hash = 0
  for (let i = 0; i < str.length; i++) {
    hash = ((hash << 5) - hash) + str.charCodeAt(i)
    hash = hash & hash
  }
  return hash
}
const initYjs = () => {
  return new Promise((resolve) => {
    ydoc.value = new Y.Doc()
    ytext.value = ydoc.value.getText('shared-text')
    
    const roomName = route.params.id
    console.log('Connecting to room:', roomName)

    provider.value = new WebsocketProvider(
      `ws://localhost:1234`, // 使用正确的 WebSocket URL
      roomName,
      ydoc.value,
      {
        WebSocketPolyfill: class extends WebSocket {
          constructor(...args) {
            super(...args)
            this.addEventListener('message', (event) => {
              console.log('WebSocket received message:', new Uint8Array(event.data))
            })
            this.addEventListener('open', () => {
              console.log('WebSocket connection opened')
            })
          }
        }
      }
    )

    const username = localStorage.getItem('username') || '匿名用户'
    provider.value.awareness.setLocalStateField('user', {
      name: username,
      color: getClientColor(ydoc.value.clientID),
      docId: route.params.id
    })

    provider.value.on('status', ({ status }) => {
      console.log('Connection status:', status)
      if (status === 'connected') {
        console.log('Connected to room:', roomName, 'with client ID:', ydoc.value.clientID)
        resolve() // 确保连接成功后再进行后续操作
      }
    })

    // 监听同步状态
    provider.value.on('synced', (isSynced) => {
      console.log('文档同步状态:', isSynced)
      if (isSynced) {
        // 增加延时，确保接收到远程内容
        setTimeout(() => {
          console.log('远程同步内容:', ytext.value.toString())
        }, 1000)  // 增加延时，确保同步完成
      }
    })
  })
}


const initVditor = async () => {
  try {
    console.log('Initializing Vditor...')
    return new Promise((resolve) => {
      vditor.value = new Vditor('vditor', {
        height: '100%',
        mode: 'ir',
        preview: {
          theme: {
            current: 'light'
          },
          hljs: {
            enable: true,
            style: 'github'
          }
        },
        toolbar: [
          'emoji',
          'headings',
          'bold',
          'italic',
          'strike',
          'link',
          '|',
          'list',
          'ordered-list',
          'check',
          'outdent',
          'indent',
          '|',
          'quote',
          'line',
          'code',
          'inline-code',
          'insert-before',
          'insert-after',
          '|',
          'upload',
          'table',
          '|',
          'undo',
          'redo',
          '|',
          'outline',
          'preview',
          'fullscreen',
          'content-theme',
          'code-theme',
          'export'
        ],
        counter: {
          enable: true
        },
        cache: {
          enable: false
        },
        after: () => {
          console.log('Vditor initialized successfully')
          
          // 延时等待，确保 Yjs 内容加载完成
          setTimeout(() => {
            const initialContent = ytext.value.toString()
            console.log('Setting initial content:', initialContent)
            if (initialContent) {
              vditor.value.setValue(initialContent)
            }
          }, 1000)

          // 监听 Yjs 文本变化
          ytext.value.observe(() => {
            const newContent = ytext.value.toString()
            console.log('Yjs content updated:', newContent)
            if (newContent !== vditor.value.getValue()) {
              console.log('Updating Vditor content:', newContent)
              vditor.value.setValue(newContent)
            }
          })

          resolve()
        },
        input: (value) => {
          if (!ytext.value || value === ytext.value.toString()) return
          ytext.value.delete(0, ytext.value.length)
          ytext.value.insert(0, value)
        }
      })
    })
  } catch (error) {
    console.error('初始化编辑器失败:', error)
    ElMessage.error('初始化编辑器失败')
  }
}


// 验证密码
const verifyPassword = async () => {
  if (!inputPassword.value.trim()) {
    ElMessage.warning('密码不能为空')
    return
  }

  showPasswordDialog.value = false
  isVerified.value = true
  
  // 等待 DOM 更新
  await nextTick()
  
  // 初始化 Yjs 并等待连接成功
  await initYjs()
  console.log('Yjs initialized and connected, current content:', ytext.value.toString())
  
  // 初始化编辑器
  await initVditor()
}

// 清理
onBeforeUnmount(() => {
  if (provider.value) {
    provider.value.destroy()
  }
  if (ydoc.value) {
    ydoc.value.destroy()
  }
})

// Add a computed property to calculate the number of connected users
const connectedUserCount = computed(() => {
  if (!provider.value?.awareness?.getStates()) return 0
  return provider.value.awareness.getStates().size
})
</script>

<style lang="less" scoped>
.shared-doc-view {
  min-height: 100vh;
  background-color: #f5f7fa;
  padding: 20px;
  box-sizing: border-box;

  .doc-container {
    .header {
      background: #fff;
      padding: 16px 24px;
      border-radius: 8px;
      box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
      margin-bottom: 20px;
      display: flex;
      justify-content: space-between;
      align-items: center;

      .title {
        display: flex;
        align-items: center;
        gap: 12px;

        h2 {
          margin: 0;
          color: #303133;
          font-size: 18px;
          font-weight: 500;
        }
      }

      .collaborators {
        display: flex;
        gap: 8px;

        .collaborator {
          .avatar {
            width: 32px;
            height: 32px;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            color: #fff;
            font-weight: 500;
            font-size: 14px;
            cursor: default;
          }
        }
      }

      .connected-users {
        font-size: 14px;
        color: #303133;
        margin-left: 20px;
      }
    }

    .content {
      background: #fff;
      padding: 24px;
      border-radius: 8px;
      box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
      height: calc(100vh - 120px);

      .vditor {
        height: 100%;
      }
    }
  }
}

.password-dialog-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px 0;

  .dialog-icon {
    font-size: 48px;
    color: #409EFF;
    margin-bottom: 20px;
  }

  .dialog-title {
    font-size: 16px;
    font-weight: 500;
    color: #303133;
    margin-bottom: 20px;
  }

  .el-input {
    width: 100%;
  }
}

:deep(.el-dialog) {
  border-radius: 8px;
  
  .el-dialog__header {
    margin: 0;
    padding: 20px;
    border-bottom: 1px solid #DCDFE6;
  }
  
  .el-dialog__body {
    padding: 20px;
  }
  
  .el-dialog__footer {
    padding: 20px;
    border-top: 1px solid #DCDFE6;
    display: flex;
    justify-content: center;
  }
}

:deep(.remote-cursor) {
  position: absolute;
  width: 2px;
  height: 20px;
  background-color: #409EFF;
  pointer-events: none;
  z-index: 1000;
  transition: transform 0.1s ease;

  &::after {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    border-left: 4px solid transparent;
    border-right: 4px solid transparent;
    border-bottom: 4px solid currentColor;
    transform: translateX(-3px) translateY(-4px);
  }
}
</style> 