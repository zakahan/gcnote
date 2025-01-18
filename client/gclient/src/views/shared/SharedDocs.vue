<template>
  <div class="shared-container">
    <div class="header">
      <h2>共享文档</h2>
    </div>

    <div class="file-list">
      <div 
        v-for="file in sharedFiles" 
        :key="file.shareFileId"
        class="file-item"
      >
        <div class="file-info" @click="handleFileClick(file)">
          <el-icon><Document /></el-icon>
          <span class="file-name">{{ file.fileName }}</span>
          <span class="share-time">{{ formatTime(file.createdAt) }}</span>
          <div class="password-section" @click.stop>
            <span class="label">密码：</span>
            <span class="password">{{ showPassword[file.shareFileId] ? file.password : '••••••••' }}</span>
            <el-icon 
              class="action-icon" 
              @click="togglePassword(file.shareFileId)"
            >
              <component :is="showPassword[file.shareFileId] ? 'View' : 'Hide'" />
            </el-icon>
            <el-icon 
              class="action-icon delete-icon" 
              @click="handleDelete(file)"
            >
              <Delete />
            </el-icon>
          </div>
        </div>
      </div>
      <el-empty v-if="sharedFiles.length === 0" description="暂无共享文档" />
    </div>

    <!-- 删除确认对话框 -->
    <el-dialog
      v-model="showDeleteDialog"
      title="取消分享"
      width="30%"
      :close-on-click-modal="false"
    >
      <div class="dialog-content">
        <el-icon class="dialog-icon warning-icon" :size="50"><Warning /></el-icon>
        <div class="dialog-title">确定要取消分享此文档吗？</div>
        <div class="dialog-desc">取消分享后，其他用户将无法通过此链接访问该文档</div>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showDeleteDialog = false">取消</el-button>
          <el-button type="danger" @click="confirmDelete">确定取消</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Document, Delete, Warning } from '@element-plus/icons-vue'
import axios from 'axios'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'

dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

const router = useRouter()
const sharedFiles = ref([])
const showPassword = ref({})
const showDeleteDialog = ref(false)
const fileToDelete = ref(null)

// 格式化时间
const formatTime = (time) => {
  return dayjs(time).fromNow()
}

// 切换密码显示状态
const togglePassword = (fileId) => {
  showPassword.value[fileId] = !showPassword.value[fileId]
}

// 处理删除
const handleDelete = (file) => {
  fileToDelete.value = file
  showDeleteDialog.value = true
}

// 确认删除
const confirmDelete = async () => {
  try {
    const token = localStorage.getItem('token')
    const response = await axios.post(
      'http://localhost:8086/share/delete',
      {
        share_file_id: fileToDelete.value.shareFileId
      },
      {
        headers: {
          'token': token,
          'Content-Type': 'application/json'
        }
      }
    )

    if (response.data.Code === 0) {
      ElMessage.success('取消分享成功')
      showDeleteDialog.value = false
      await fetchSharedFiles() // 刷新列表
    } else {
      ElMessage.error(response.data.Msg || '取消分享失败')
    }
  } catch (error) {
    console.error('取消分享失败:', error)
    ElMessage.error('取消分享失败')
  }
}

// 获取共享文档列表
const fetchSharedFiles = async () => {
  try {
    const token = localStorage.getItem('token')
    if (!token) {
      router.push('/login')
      return
    }

    const response = await axios.get('http://localhost:8086/share/info', {
      headers: { token }
    })

    if (response.data.Code === 0) {
      sharedFiles.value = response.data.Data.list
      // 初始化密码显示状态
      sharedFiles.value.forEach(file => {
        showPassword.value[file.shareFileId] = false
      })
    } else {
      ElMessage.error(response.data.Msg || '获取共享文档失败')
    }
  } catch (error) {
    console.error('获取共享文档失败:', error)
    ElMessage.error('获取共享文档失败')
  }
}

// 处理文件点击
const handleFileClick = (file) => {
  const url = `/shared-doc/${file.shareFileId}`
  window.open(url, '_blank')
}

onMounted(() => {
  fetchSharedFiles()
})
</script>

<style lang="less" scoped>
.shared-container {
  padding: 20px;

  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;

    h2 {
      margin: 0;
      color: #303133;
    }
  }

  .file-list {
    background: #fff;
    border-radius: 8px;
    padding: 20px;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);

    .file-item {
      display: flex;
      align-items: center;
      padding: 12px;
      border-radius: 4px;
      transition: background-color 0.3s;

      &:hover {
        background-color: #f5f7fa;
      }

      .file-info {
        display: flex;
        align-items: center;
        gap: 16px;
        flex: 1;
        cursor: pointer;

        &:hover .file-name {
          color: #409EFF;
        }

        .file-name {
          color: #303133;
          flex: 1;
        }

        .share-time {
          color: #909399;
          font-size: 13px;
          min-width: 100px;
        }

        .password-section {
          display: flex;
          align-items: center;
          gap: 8px;
          min-width: 200px;

          .label {
            color: #909399;
          }

          .password {
            font-family: monospace;
            color: #606266;
            width: 128px;  // 固定宽度，16个字符
            display: inline-block;
            text-align: left;
          }

          .action-icon {
            cursor: pointer;
            color: #606266;
            font-size: 16px;
            display: flex;
            align-items: center;
            justify-content: center;
            width: 24px;
            height: 24px;
            border-radius: 4px;
            transition: all 0.3s;

            &:hover {
              background-color: #f0f0f0;
              color: #409EFF;
            }

            &.delete-icon {
              &:hover {
                color: #F56C6C;
              }
            }
          }
        }
      }
    }
  }
}

.dialog-content {
  padding: 20px 0;
  text-align: center;
  
  .dialog-icon {
    color: #409EFF;
    margin-bottom: 20px;

    &.warning-icon {
      color: #E6A23C;
    }
  }
  
  .dialog-title {
    font-size: 16px;
    font-weight: 500;
    color: #303133;
    margin-bottom: 10px;
  }
  
  .dialog-desc {
    font-size: 14px;
    color: #909399;
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

:deep(.el-dialog) {
  border-radius: 8px;
  
  .el-dialog__header {
    margin: 0;
    padding: 20px;
    border-bottom: 1px solid #DCDFE6;
  }
  
  .el-dialog__headerbtn {
    top: 20px;
  }
  
  .el-dialog__title {
    font-size: 16px;
    font-weight: 500;
  }
  
  .el-dialog__body {
    padding: 20px;
  }
  
  .el-dialog__footer {
    padding: 20px;
    border-top: 1px solid #DCDFE6;
  }
}
</style> 