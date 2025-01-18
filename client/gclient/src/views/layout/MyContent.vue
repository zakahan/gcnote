/* eslint-disable */
<template>
  <div class="content-container">
    <!-- 操作按钮 -->
    <div class="action-buttons">
      <el-button type="primary" @click="handleCreateKb">
        <el-icon><Plus /></el-icon>新建知识库
      </el-button>
      <el-button type="primary" @click="handleCreateDoc">
        <el-icon><Document /></el-icon>新建文档
      </el-button>
      <el-button type="primary" @click="handleImportDoc">
        <el-icon><Upload /></el-icon>导入文档
      </el-button>
    </div>

    <!-- 最近文档 -->
    <div class="recent-docs">
      <div class="section">
        <h3>最近创建</h3>
        <div class="doc-list">
          <div 
            v-for="doc in recentCreated" 
            :key="doc.file_id"
            class="doc-item"
          >
            <div class="doc-info" @click="openDocument(doc)">
              <el-icon><Document /></el-icon>
              <span class="doc-name">{{ doc.file_name }}</span>
              <span class="doc-time">{{ formatTime(doc.created_at) }}</span>
            </div>
            <el-dropdown 
              trigger="hover" 
              @command="(command) => handleCommand(command, doc)"
              class="doc-actions"
            >
              <span class="el-dropdown-link">
                <el-icon><MoreFilled /></el-icon>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="rename">重命名</el-dropdown-item>
                  <el-dropdown-item command="delete" style="color: red">删除</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
          <el-empty v-if="recentCreated.length === 0" description="暂无文档" />
        </div>
      </div>

      <div class="section">
        <h3>最近修改</h3>
        <div class="doc-list">
          <div 
            v-for="doc in recentModified" 
            :key="doc.file_id"
            class="doc-item"
          >
            <div class="doc-info" @click="openDocument(doc)">
              <el-icon><Document /></el-icon>
              <span class="doc-name">{{ doc.file_name }}</span>
              <span class="doc-time">{{ formatTime(doc.modified_at) }}</span>
            </div>
            <el-dropdown 
              trigger="hover" 
              @command="(command) => handleCommand(command, doc)"
              class="doc-actions"
            >
              <span class="el-dropdown-link">
                <el-icon><MoreFilled /></el-icon>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="rename">重命名</el-dropdown-item>
                  <el-dropdown-item command="delete" style="color: red">删除</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
          <el-empty v-if="recentModified.length === 0" description="暂无文档" />
        </div>
      </div>
    </div>

    <!-- 新建知识库对话框 -->
    <el-dialog v-model="createKbDialog" title="新建知识库" width="30%">
      <el-form :model="newKbForm" label-width="80px">
        <el-form-item label="名称">
          <el-input v-model="newKbForm.name" placeholder="请输入知识库名称" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="createKbDialog = false">取消</el-button>
          <el-button type="primary" @click="submitCreateKb">确定</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 新建文档对话框 -->
    <el-dialog v-model="createDocDialog" title="新建文档" width="30%">
      <el-form :model="newDocForm" label-width="80px">
        <el-form-item label="标题">
          <el-input v-model="newDocForm.title" placeholder="请输入文档标题" />
        </el-form-item>
        <el-form-item label="知识库">
          <el-select v-model="newDocForm.kbId" placeholder="请选择知识库">
            <el-option
              v-for="item in kbList"
              :key="item.IndexId"
              :label="item.IndexName"
              :value="item.IndexId"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="createDocDialog = false">取消</el-button>
          <el-button type="primary" @click="submitCreateDoc">确定</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 导入文档对话框 -->
    <el-dialog v-model="importDocDialog" title="导入文档" width="30%">
      <el-form :model="importDocForm" label-width="80px">
        <el-form-item label="知识库">
          <el-select v-model="importDocForm.kbId" placeholder="请选择知识库">
            <el-option
              v-for="item in kbList"
              :key="item.IndexId"
              :label="item.IndexName"
              :value="item.IndexId"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="文件">
          <el-upload
            class="upload-demo"
            action="#"
            :auto-upload="false"
            :on-change="handleFileChange"
            :limit="1"
          >
            <template #trigger>
              <el-button type="primary">选择文件</el-button>
            </template>
            <template #tip>
              <div class="el-upload__tip">请选择要导入的文档文件</div>
            </template>
          </el-upload>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="importDocDialog = false">取消</el-button>
          <el-button type="primary" @click="submitImportDoc">确定</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 重命名对话框 -->
    <el-dialog
      v-model="renameDialogVisible"
      title="重命名文件"
      width="30%"
    >
      <el-input v-model="newFileName" placeholder="请输入新的文件名" />
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="renameDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleRename">确认</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, defineEmits } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Document, MoreFilled, Plus, Upload } from '@element-plus/icons-vue'
import axios from 'axios'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'

dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

const router = useRouter()
const emit = defineEmits(['kb-created'])

const recentCreated = ref([])
const recentModified = ref([])
const renameDialogVisible = ref(false)
const newFileName = ref('')
const currentDoc = ref(null)

// 对话框状态
const createKbDialog = ref(false)
const createDocDialog = ref(false)
const importDocDialog = ref(false)
const newKbForm = ref({ name: '' })
const newDocForm = ref({ title: '', kbId: '' })
const importDocForm = ref({ kbId: '', file: null })
const kbList = ref([])

// 获取最近文档
const fetchRecentDocs = async (mode) => {
  try {
    const token = localStorage.getItem('token')
    if (!token) {
      router.push('/login')
      return
    }

    const response = await axios.post(
      'http://localhost:8086/index/recent_docs',
      { mode },
      {
        headers: {
          'token': token,
          'Content-Type': 'application/json'
        }
      }
    )

    if (response.data.Code === 0) {
      if (mode === 'created') {
        recentCreated.value = response.data.Data
      } else {
        recentModified.value = response.data.Data
      }
    }
  } catch (error) {
    console.error('获取最近文档失败:', error)
    ElMessage.error('获取最近文档失败')
  }
}

// 获取知识库列表
const fetchKbList = async () => {
  try {
    const token = localStorage.getItem('token')
    const response = await axios.get('http://localhost:8086/index/show_indexes', {
      headers: { token }
    })
    if (response.data.Code === 0) {
      kbList.value = response.data.Data
    }
  } catch (error) {
    console.error('获取知识库列表失败:', error)
  }
}

// 格式化时间
const formatTime = (time) => {
  return dayjs(time).fromNow()
}

// 打开文档
const openDocument = (doc) => {
  router.push({
    name: 'knowledge-base',
    params: {
      index_id: doc.index_id,
      file_id: doc.file_id,
      file_name: doc.file_name
    }
  })
}

// 处理菜单命令
const handleCommand = (command, doc) => {
  if (command === 'rename') {
    currentDoc.value = doc
    newFileName.value = doc.file_name
    renameDialogVisible.value = true
  } else if (command === 'delete') {
    handleDelete(doc)
  }
}

// 处理重命名
const handleRename = async () => {
  try {
    const token = localStorage.getItem('token')
    const response = await axios.post(
      'http://localhost:8086/index/rename_file',
      {
        index_id: currentDoc.value.index_id,
        kb_file_id: currentDoc.value.file_id,
        kb_file_name: currentDoc.value.file_name,
        dest_kb_file_name: newFileName.value
      },
      {
        headers: {
          'token': token,
          'Content-Type': 'application/json'
        }
      }
    )

    if (response.data.Code === 0) {
      ElMessage.success('重命名成功')
      renameDialogVisible.value = false
      // 刷新列表
      fetchRecentDocs('created')
      fetchRecentDocs('modified')
    } else {
      ElMessage.error(response.data.Msg || '重命名失败')
    }
  } catch (error) {
    console.error('重命名失败:', error)
    ElMessage.error('重命名失败')
  }
}

// 处理删除
const handleDelete = (doc) => {
  ElMessageBox.confirm(
    '确定要删除该文件吗？',
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      const token = localStorage.getItem('token')
      const response = await axios.post(
        'http://localhost:8086/index/recycle_file',
        {
          index_id: doc.index_id,
          kb_file_id: doc.file_id,
          kb_file_name: doc.file_name
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
        // 刷新列表
        fetchRecentDocs('created')
        fetchRecentDocs('modified')
      } else {
        ElMessage.error(response.data.Msg || '删除失败')
      }
    } catch (error) {
      console.error('删除失败:', error)
      ElMessage.error('删除失败')
    }
  }).catch(() => {})
}

// 处理新建知识库
const handleCreateKb = () => {
  createKbDialog.value = true
}

const submitCreateKb = async () => {
  try {
    const token = localStorage.getItem('token')
    const response = await axios.post('http://localhost:8086/index/create_index', {
      index_name: newKbForm.value.name
    }, {
      headers: { token }
    })
    
    if (response.data.Code === 0) {
      ElMessage.success('创建知识库成功')
      createKbDialog.value = false
      newKbForm.value.name = ''
      await fetchKbList()
      // 通知父组件刷新菜单
      emit('kb-created')
    } else {
      ElMessage.error(response.data.Msg || '创建知识库失败')
    }
  } catch (error) {
    console.error('创建知识库失败:', error)
    ElMessage.error('创建知识库失败')
  }
}

// 处理新建文档
const handleCreateDoc = () => {
  createDocDialog.value = true
}

const submitCreateDoc = async () => {
  try {
    const token = localStorage.getItem('token')
    const response = await axios.post('http://localhost:8086/index/create_file', {
      index_id: newDocForm.value.kbId,
      kb_file_name: newDocForm.value.title
    }, {
      headers: { token }
    })
    
    if (response.data.Code === 0) {
      ElMessage.success('创建文档成功')
      createDocDialog.value = false
      newDocForm.value = { title: '', kbId: '' }
      // 刷新列表
      fetchRecentDocs('created')
      fetchRecentDocs('modified')
    } else {
      ElMessage.error(response.data.Msg || '创建文档失败')
    }
  } catch (error) {
    console.error('创建文档失败:', error)
    ElMessage.error('创建文档失败')
  }
}

// 处理导入文档
const handleImportDoc = () => {
  importDocDialog.value = true
}

const handleFileChange = (file) => {
  importDocForm.value.file = file.raw
}

const submitImportDoc = async () => {
  try {
    if (!importDocForm.value.file) {
      ElMessage.warning('请选择要导入的文件')
      return
    }

    const token = localStorage.getItem('token')
    if (!importDocForm.value.kbId) {
      ElMessage.warning('请选择知识库')
      return
    }

    const formData = new FormData()
    formData.append('index_id', importDocForm.value.kbId)
    formData.append('file', importDocForm.value.file)

    const response = await axios.post('http://localhost:8086/index/add_file', formData, {
      headers: {
        token,
        'Content-Type': 'multipart/form-data'
      }
    })

    if (response.data.Code === 0) {
      ElMessage.success('导入文档成功')
      importDocDialog.value = false
      importDocForm.value = { kbId: '', file: null }
      // 刷新列表
      fetchRecentDocs('created')
      fetchRecentDocs('modified')
    } else {
      ElMessage.error(response.data.Msg || '导入文档失败')
    }
  } catch (error) {
    console.error('导入文档失败:', error)
    ElMessage.error('导入文档失败')
  }
}

onMounted(async () => {
  await fetchKbList()
  fetchRecentDocs('created')
  fetchRecentDocs('modified')
})
</script>

<style lang="less" scoped>
.content-container {
  padding: 20px;

  .action-buttons {
    margin-bottom: 30px;

  }

  .recent-docs {
    display: flex;
    gap: 40px;

    .section {
      flex: 1;
      background: #fff;
      border-radius: 8px;
      padding: 20px;
      box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);

      h3 {
        color: #606266;
        font-weight: 500;
        margin-bottom: 16px;
      }

      .doc-list {
        .doc-item {
          display: flex;
          align-items: center;
          justify-content: space-between;
          padding: 12px;
          border-radius: 4px;
          cursor: pointer;
          transition: background-color 0.3s;

          &:hover {
            background-color: #f5f7fa;

            .doc-actions {
              opacity: 1;
            }
          }

          .doc-info {
            display: flex;
            align-items: center;
            gap: 8px;
            flex: 1;

            .doc-name {
              color: #303133;
              flex: 1;
            }

            .doc-time {
              color: #909399;
              font-size: 13px;
              margin-left: 16px;
            }
          }

          .doc-actions {
            opacity: 0;
            transition: opacity 0.3s;
            color: #909399;
          }
        }
      }
    }
  }
}
</style>


/* eslint-enable */