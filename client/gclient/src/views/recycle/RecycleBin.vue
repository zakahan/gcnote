<template>
  <div class="recycle-container">
    <div class="header">
      <h2>回收站</h2>
      <el-button 
        type="danger" 
        :icon="Delete"
        @click="handleClearAll"
      >
        清空回收站
      </el-button>
    </div>

    <div class="file-list">
      <div 
        v-for="file in recycleFiles" 
        :key="file.ID"
        class="file-item"
      >
        <div class="file-info">
          <el-icon><Document /></el-icon>
          <span class="file-name">{{ file.KBFileName }}</span>
          <span class="delete-time">{{ formatTime(file.CreatedAt) }}</span>
        </div>
        <el-dropdown 
          trigger="hover" 
          @command="(command) => handleCommand(command, file)"
        >
          <el-icon class="more-icon"><MoreFilled /></el-icon>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="restore">
                <el-icon><RefreshRight /></el-icon>还原
              </el-dropdown-item>
              <el-dropdown-item command="delete" style="color: #f56c6c;">
                <el-icon><Delete /></el-icon>彻底删除
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
      <el-empty v-if="recycleFiles.length === 0" description="回收站为空" />
    </div>

    <!-- 彻底删除确认对话框 -->
    <el-dialog
      v-model="showDeleteDialog"
      title="彻底删除"
      width="30%"
      :close-on-click-modal="false"
    >
      <div class="dialog-content">
        <el-icon class="dialog-icon warning-icon" :size="50"><Warning /></el-icon>
        <div class="dialog-title">确定要彻底删除此文件吗？</div>
        <div class="dialog-desc">此操作将永久删除该文件，无法恢复！</div>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showDeleteDialog = false">取消</el-button>
          <el-button type="danger" @click="confirmDelete">确定删除</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 清空回收站确认对话框 -->
    <el-dialog
      v-model="showClearDialog"
      title="清空回收站"
      width="30%"
      :close-on-click-modal="false"
    >
      <div class="dialog-content">
        <el-icon class="dialog-icon warning-icon" :size="50"><Warning /></el-icon>
        <div class="dialog-title">确定要清空回收站吗？</div>
        <div class="dialog-desc">此操作将永久删除回收站中的所有文件，无法恢复！</div>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showClearDialog = false">取消</el-button>
          <el-button type="danger" @click="confirmClear">确定清空</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Document, Delete, MoreFilled, RefreshRight, Warning } from '@element-plus/icons-vue'
import axios from 'axios'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'

dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

const router = useRouter()
const recycleFiles = ref([])
const showDeleteDialog = ref(false)
const showClearDialog = ref(false)
const fileToDelete = ref(null)

// 格式化时间
const formatTime = (time) => {
  return dayjs(time).fromNow()
}

// 获取回收站文件列表
const fetchRecycleFiles = async () => {
  try {
    const token = localStorage.getItem('token')
    if (!token) {
      router.push('/login')
      return
    }

    const response = await axios.get('http://localhost:8086/recycle/show_files', {
      headers: { token }
    })

    if (response.data.Code === 0) {
      recycleFiles.value = response.data.Data
    } else {
      ElMessage.error(response.data.Msg || '获取回收站文件失败')
    }
  } catch (error) {
    console.error('获取回收站文件失败:', error)
    ElMessage.error('获取回收站文件失败')
  }
}

// 处理命令
const handleCommand = (command, file) => {
  if (command === 'restore') {
    handleRestore(file)
  } else if (command === 'delete') {
    fileToDelete.value = file
    showDeleteDialog.value = true
  }
}

// 处理还原
const handleRestore = async (file) => {
  try {
    const token = localStorage.getItem('token')
    const response = await axios.post(
      'http://localhost:8086/recycle/restore',
      {
        index_id: file.SourceIndexId,
        kb_file_id: file.KBFileId
      },
      {
        headers: {
          'token': token,
          'Content-Type': 'application/json'
        }
      }
    )

    if (response.data.Code === 0) {
      ElMessage.success('还原成功')
      await fetchRecycleFiles()
    } else {
      ElMessage.error(response.data.Msg || '还原失败')
    }
  } catch (error) {
    console.error('还原失败:', error)
    ElMessage.error('还原失败')
  }
}

// 处理清空回收站
const handleClearAll = () => {
  showClearDialog.value = true
}

// 确认清空回收站
const confirmClear = async () => {
  try {
    const token = localStorage.getItem('token')
    const response = await axios.get('http://localhost:8086/recycle/clear', {
      headers: { token }
    })

    if (response.data.Code === 0) {
      ElMessage.success('清空成功')
      showClearDialog.value = false
      await fetchRecycleFiles()
    } else {
      ElMessage.error(response.data.Msg || '清空失败')
    }
  } catch (error) {
    console.error('清空失败:', error)
    ElMessage.error('清空失败')
  }
}

// 确认删除
const confirmDelete = async () => {
  try {
    const token = localStorage.getItem('token')
    const response = await axios.post(
      'http://localhost:8086/recycle/delete_file',
      {
        index_id: fileToDelete.value.SourceIndexId,
        kb_file_id: fileToDelete.value.KBFileId
      },
      {
        headers: {
          'token': token,
          'Content-Type': 'application/json'
        }
      }
    )

    if (response.data.Code === 0) {
      ElMessage.success('删除成功')
      showDeleteDialog.value = false
      await fetchRecycleFiles()
    } else {
      ElMessage.error(response.data.Msg || '删除失败')
    }
  } catch (error) {
    console.error('删除失败:', error)
    ElMessage.error('删除失败')
  }
}

onMounted(() => {
  fetchRecycleFiles()
})
</script>

<style lang="less" scoped>
.recycle-container {
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
      justify-content: space-between;
      padding: 12px;
      border-radius: 4px;
      transition: background-color 0.3s;

      &:hover {
        background-color: #f5f7fa;
      }

      .file-info {
        display: flex;
        align-items: center;
        gap: 8px;
        flex: 1;

        .file-name {
          color: #303133;
          flex: 1;
        }

        .delete-time {
          color: #909399;
          font-size: 13px;
          margin-left: 16px;
        }
      }

      .more-icon {
        font-size: 16px;
        color: #909399;
        cursor: pointer;

        &:hover {
          color: #409EFF;
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