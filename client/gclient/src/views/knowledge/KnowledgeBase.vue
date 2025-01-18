<template>
  <div class="kb-container">
    <div class="kb-content">
      <!-- 左侧文件列表 -->
      <div class="file-list">
        <!-- 知识库标题和新建按钮 -->
        <div class="kb-header">
          <div class="kb-title">
            <el-icon><FolderOpened /></el-icon>
            <span>{{ currentKbName }}</span>
          </div>
          <el-dropdown 
            trigger="click"
            @command="handleCreateCommand"
            class="create-btn"
          >
            <el-button type="primary" :icon="Plus" size="small">
              新建
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="new">新建文档</el-dropdown-item>
                <el-dropdown-item command="import">导入文档</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>

        <!-- 搜索框 -->
        <div class="search-box">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索文档..."
            :prefix-icon="Search"
            @click="showSearchDialog = true"
            readonly
          />
        </div>

        <el-menu class="file-menu">
          <el-menu-item
            v-for="file in kbFiles"
            :key="file.KBFileId"
            :index="file.KBFileId"
            @click="handleFileClick(file)"
          >
            <div class="file-item">
              <div class="file-info">
                <el-icon><Document /></el-icon>
                <span>{{ file.KBFileName }}</span>
              </div>
              <el-dropdown 
                trigger="hover" 
                @command="(command) => handleCommand(command, file)"
                class="file-actions"
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
          </el-menu-item>
        </el-menu>
        <el-empty v-if="kbFiles.length === 0" description="无文档" />
      </div>

      <!-- 右侧内容区域 -->
      <div class="content-area">
        <div v-if="currentFile" class="editor-container">
          <div class="editor-header">
            <div class="file-info">
              <el-icon><Document /></el-icon>
              <span class="file-name">{{ currentFile.KBFileName }}</span>
            </div>
            <div class="editor-actions">

              <el-button 
                type="primary" 
                :icon="Upload"
                size="small"
                :loading="updating"
                @click="handleUpdate"
              >
                更新文档
              </el-button>
            </div>
          </div>
          <div id="vditor" class="vditor" />
        </div>
        <el-empty v-else description="请选择要查看的文档" />
      </div>

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

      <!-- 新建文档对话框 -->
      <el-dialog
        v-model="createDialogVisible"
        :title="createMode === 'new' ? '新建文档' : '导入文档'"
        width="30%"
      >
        <el-form v-if="createMode === 'new'">
          <el-form-item label="文档名称">
            <el-input v-model="newDocName" placeholder="请输入文档名称" />
          </el-form-item>
        </el-form>
        <el-upload
          v-else
          class="upload-demo"
          drag
          action="http://localhost:8086/index/add_file"
          :headers="uploadHeaders"
          :data="uploadData"
          :on-success="handleUploadSuccess"
          :on-error="handleUploadError"
        >
          <el-icon class="el-icon--upload"><upload-filled /></el-icon>
          <div class="el-upload__text">
            拖拽文件到此处或 <em>点击上传</em>
          </div>
        </el-upload>
        <template #footer>
          <span class="dialog-footer">
            <el-button @click="createDialogVisible = false">取消</el-button>
            <el-button v-if="createMode === 'new'" type="primary" @click="handleCreate">确认</el-button>
          </span>
        </template>
      </el-dialog>

      <!-- 搜索对话框 -->
      <el-dialog
        v-model="showSearchDialog"
        title="细粒度搜索"
        width="50%"
        :close-on-click-modal="false"
      >
        <div class="search-dialog-content">
          <div class="search-input-area">
            <el-input
              v-model="searchKeyword"
              placeholder="输入关键词搜索..."
              :prefix-icon="Search"
              @keyup.enter="handleSearch"
            >
              <template #append>
                <el-button @click="handleSearch">搜索</el-button>
              </template>
            </el-input>
            <div class="search-options">
              <el-form :inline="true">
                <el-form-item label="检索方式">
                  <el-select v-model="searchMethod" placeholder="选择检索方式">
                    <el-option label="全文检索" value="full_text_search" />
                    <el-option label="向量检索" value="vector_search" />
                    <el-option label="关键词检索" value="keyword_search" />
                  </el-select>
                </el-form-item>
                <el-form-item label="返回数量">
                  <el-input-number 
                    v-model="topK" 
                    :min="1" 
                    :max="20" 
                    :step="1"
                    size="small"
                  />
                </el-form-item>
              </el-form>
            </div>
          </div>

          <div class="search-results" v-loading="searching">
            <template v-if="searchResults.length > 0">
              <div 
                v-for="result in searchResults" 
                :key="result.metadata?.kb_file_id"
                class="search-result-item"
                @click="handleResultClick(result)"
              >
                <div class="result-title">
                  <el-icon><Document /></el-icon>
                  <span>{{ getFileName(result.metadata?.kb_file_id) }}</span>
                </div>
                <div class="result-content">
                  {{ result.page_content }}
                </div>
              </div>
            </template>
            <el-empty v-else-if="!searching" description="暂无搜索结果" />
          </div>
        </div>
      </el-dialog>


      <!-- 删除确认对话框 -->
      <el-dialog
        v-model="showDeleteDialog"
        title="删除文件"
        width="30%"
        :close-on-click-modal="false"
      >
        <div class="share-dialog-content">
          <el-icon class="share-icon warning-icon" :size="50"><Warning /></el-icon>
          <div class="share-title">确定要删除此文件吗？</div>
          <div class="share-desc">删除后可在回收站中找回</div>
        </div>
        <template #footer>
          <span class="dialog-footer">
            <el-button @click="showDeleteDialog = false">取消</el-button>
            <el-button type="danger" @click="confirmDelete">确定删除</el-button>
          </span>
        </template>
      </el-dialog>


    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, nextTick, onBeforeUnmount, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Document, MoreFilled, FolderOpened, Plus, UploadFilled, Search, Upload, Warning } from '@element-plus/icons-vue'
import axios from 'axios'
import Vditor from 'vditor'
import 'vditor/dist/index.css'

const props = defineProps({
  index_id: {
    type: String,
    required: true
  },
  file_id: {
    type: String,
    default: ''
  },
  file_name: {
    type: String,
    default: ''
  }
})

const route = useRoute()
const router = useRouter()
const kbFiles = ref([])
const token = ref(null)
const currentFile = ref(null)
const vditor = ref(null)
const renameDialogVisible = ref(false)
const newFileName = ref('')
const fileToRename = ref(null)
const createDialogVisible = ref(false)
const createMode = ref('new')
const newDocName = ref('')
const currentKbName = ref('')
const searchKeyword = ref('')
const showSearchDialog = ref(false)
const searchMethod = ref('full_text_search')
const topK = ref(5)
const searchResults = ref([])
const searching = ref(false)
const updating = ref(false)
const showDeleteDialog = ref(false)
const fileToDelete = ref(null)

// 上传相关的数据
const uploadHeaders = computed(() => ({
  'token': token.value
}))

const uploadData = computed(() => ({
  index_id: route.params.index_id
}))

// 获取知识库文件列表
const fetchKbFiles = async () => {
  try {
    token.value = localStorage.getItem('token')
    if (!token.value) {
      ElMessage.error('未登录或登录已过期')
      await router.push('/login')
      return
    }

    // 获取知识库名称
    const kbResponse = await axios.get('http://localhost:8086/index/show_indexes', {
      headers: { 
        'token': token.value
      }
    })

    if (kbResponse.data.Code === 0) {
      const currentKb = kbResponse.data.Data.find(kb => kb.IndexId === route.params.index_id)
      if (currentKb) {
        currentKbName.value = currentKb.IndexName
      }
    }

    const response = await axios.post(
      'http://localhost:8086/index/show_files',
      {
        index_id: route.params.index_id
      },
      {
        headers: { 
          'token': token.value,
          'Content-Type': 'application/json'
        }
      }
    )
    
    if (response.data.Code === 0) {
      kbFiles.value = response.data.Data
    } else {
      ElMessage.error(response.data.Msg || '获取文件列表失败')
    }
  } catch (error) {
    console.error('获取文件列表失败:', error)
    ElMessage.error('获取文件列表失败')
  }
}

// 处理新建/导入命令
const handleCreateCommand = (command) => {
  createMode.value = command
  createDialogVisible.value = true
}

// 处理上传成功
const handleUploadSuccess = (response) => {
  if (response.Code === 0) {
    ElMessage.success('上传成功')
    createDialogVisible.value = false
    fetchKbFiles()
  } else {
    ElMessage.error(response.Msg || '上传失败')
  }
}

// 处理上传失败
const handleUploadError = () => {
  ElMessage.error('上传失败')
}

// 读取文件内容
const fetchFileContent = async (file) => {
  try {
    const response = await axios.post(
      'http://localhost:8086/index/read_file',
      {
        kb_file_name: file.KBFileName,
        index_id: route.params.index_id,
        kb_file_id: file.KBFileId
      },
      {
        headers: {
          'token': token.value,
          'Content-Type': 'application/json'
        }
      }
    )

    if (response.data.Code === 0) {
      return response.data.Data
    } else {
      ElMessage.error(response.data.Msg || '读取文件失败')
      return null
    }
  } catch (error) {
    console.error('读取文件失败:', error)
    ElMessage.error('读取文件失败')
    return null
  }
}

// 初始化编辑器
const initVditor = async () => {
  // 如果已经存在实例，先销毁
  if (vditor.value) {
    vditor.value.destroy()
    vditor.value = null
  }

  await nextTick()
  vditor.value = new Vditor('vditor', {
    height: '100%',
    mode: 'ir',
    preview: {
      theme: {
        current: 'light'
      }
    },
    upload: {
      url: 'http://localhost:8086/images/upload',
      fieldName: 'file',
      headers: {
        'token': token.value
      },
      success: (_, response) => {
        try {
          const res = JSON.parse(response)
          if (res.Code === 0) {
            vditor.value.insertValue(`![image](${res.Data.url})`)
          } else {
            ElMessage.error(res.Msg || '上传失败')
          }
        } catch (error) {
          console.error('解析响应失败:', error)
          ElMessage.error('上传失败')
        }
      },
      error: () => {
        ElMessage.error('图片上传失败')
      }
    },
    after: () => {
      // 编辑器初始化完成后，如果有当前文件，设置内容
      if (currentFile.value) {
        fetchFileContent(currentFile.value).then(content => {
          if (content !== null && vditor.value) {
            setTimeout(() => {
              vditor.value.setValue(content)
            }, 100)
          }
        })
      }
    }
  })
}

// 监听 currentFile 变化，重新初始化编辑器以更新上传参数
watch(() => currentFile.value, () => {
  if (currentFile.value) {
    initVditor()
  }
})

// 处理文件点击
const handleFileClick = async (file) => {
  currentFile.value = file
  // 如果编辑器还没初始化，先初始化
  if (!vditor.value) {
    await initVditor()
  } else {
    // 如果编辑器已经初始化，直接获取并设置内容
    const content = await fetchFileContent(file)
    if (content !== null) {
      // 使用 setTimeout 确保编辑器完全准备好
      setTimeout(() => {
        vditor.value.setValue(content)
      }, 100)
    }
  }
}

// 处理下拉菜单命令
const handleCommand = (command, file) => {
  if (command === 'rename') {
    fileToRename.value = file
    newFileName.value = file.KBFileName
    renameDialogVisible.value = true
  } else if (command === 'delete') {
    handleDelete(file)
  }
}

// 处理重命名
const handleRename = async () => {
  try {
    const response = await axios.post(
      'http://localhost:8086/index/rename_file',
      {
        index_id: route.params.index_id,
        kb_file_id: fileToRename.value.KBFileId,
        kb_file_name: fileToRename.value.KBFileName,
        dest_kb_file_name: newFileName.value
      },
      {
        headers: {
          'token': token.value,
          'Content-Type': 'application/json'
        }
      }
    )

    if (response.data.Code === 0) {
      ElMessage.success('重命名成功')
      renameDialogVisible.value = false
      await fetchKbFiles() // 刷新文件列表
    } else {
      ElMessage.error(response.data.Msg || '重命名失败')
    }
  } catch (error) {
    console.error('重命名失败:', error)
    ElMessage.error('重命名失败')
  }
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
      'http://localhost:8086/index/delete_file',
      {
        index_id: route.params.index_id,
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
      await fetchKbFiles()
      if (currentFile.value?.KBFileId === fileToDelete.value.KBFileId) {
        currentFile.value = null
      }
    } else {
      ElMessage.error(response.data.Msg || '删除失败')
    }
  } catch (error) {
    console.error('删除失败:', error)
    ElMessage.error('删除失败')
  }
}

// 监听路由参数变化
watch(
  () => route.params.index_id,
  async (newId, oldId) => {
    if (newId && newId !== oldId) {
      // 清理当前状态
      currentFile.value = null
      if (vditor.value) {
        vditor.value.destroy()
        vditor.value = null
      }
      await fetchKbFiles()
      
      // 如果有文件参数，自动打开对应文件
      if (route.params.file_id && route.params.file_name) {
        const file = {
          KBFileId: route.params.file_id,
          KBFileName: route.params.file_name
        }
        await handleFileClick(file)
      }
    }
  }
)

// 组件卸载时清理
onBeforeUnmount(() => {
  if (vditor.value) {
    vditor.value.destroy()
    vditor.value = null
  }
})

onMounted(async () => {
  await fetchKbFiles()
  // 如果有文件参数，自动打开对应文件
  if (props.file_id && props.file_name) {
    const file = {
      KBFileId: props.file_id,
      KBFileName: props.file_name
    }
    await handleFileClick(file)
  }
})

// 处理搜索
const handleSearch = async () => {
  if (!searchKeyword.value.trim()) {
    ElMessage.warning('请输入搜索关键词')
    return
  }

  searching.value = true
  searchResults.value = []

  try {
    const response = await axios.post(
      'http://localhost:8086/index/retrieval',
      {
        index_id: route.params.index_id,
        query: searchKeyword.value,
        retrieval_method: searchMethod.value,
        top_k: topK.value
      },
      {
        headers: {
          'token': token.value,
          'Content-Type': 'application/json'
        }
      }
    )

    if (response.data.Code === 0) {
      searchResults.value = response.data.Data
    } else {
      ElMessage.error(response.data.Msg || '搜索失败')
    }
  } catch (error) {
    console.error('搜索失败:', error)
    ElMessage.error('搜索失败')
  } finally {
    searching.value = false
  }
}

// 处理搜索结果点击
const handleResultClick = async (result) => {
  if (result.metadata?.kb_file_id) {
    const file = {
      KBFileId: result.metadata.kb_file_id,
      KBFileName: result.metadata.kb_file_id // 使用文件ID作为文件名
    }
    await handleFileClick(file)
    showSearchDialog.value = false
  }
}

// 获取文件名
const getFileName = (fileId) => {
  const file = kbFiles.value.find(f => f.KBFileId === fileId)
  return file ? file.KBFileName : fileId
}

// 处理新建文档
const handleCreate = async () => {
  try {
    const response = await axios.post(
      'http://localhost:8086/index/create_file',
      {
        index_id: route.params.index_id,
        kb_file_name: newDocName.value
      },
      {
        headers: {
          'token': token.value,
          'Content-Type': 'application/json'
        }
      }
    )

    if (response.data.Code === 0) {
      ElMessage.success('创建成功')
      createDialogVisible.value = false
      newDocName.value = ''
      await fetchKbFiles()
    } else {
      ElMessage.error(response.data.Msg || '创建失败')
    }
  } catch (error) {
    console.error('创建失败:', error)
    ElMessage.error('创建失败')
  }
}

// 处理文档更新
const handleUpdate = async () => {
  if (!currentFile.value || !vditor.value) {
    ElMessage.warning('请先选择文件')
    return
  }

  updating.value = true
  try {
    // 获取编辑器内容并转换为文件
    const content = vditor.value.getValue()
    const file = new File([content], currentFile.value.KBFileName, {
      type: 'text/markdown'
    })

    // 创建 FormData
    const formData = new FormData()
    formData.append('index_id', route.params.index_id)
    formData.append('kb_file_id', currentFile.value.KBFileId)
    formData.append('file', file)

    // 发送更新请求
    const response = await axios.post(
      'http://localhost:8086/index/update_file',
      formData,
      {
        headers: {
          'token': token.value,
          'Content-Type': 'multipart/form-data'
        }
      }
    )

    if (response.data.Code === 0) {
      ElMessage.success('更新成功')
    } else {
      ElMessage.error(response.data.Msg || '更新失败')
    }
  } catch (error) {
    console.error('更新失败:', error)
    ElMessage.error('更新失败')
  } finally {
    updating.value = false
  }
}


</script>

<style lang="less" scoped>
.kb-container {
  height: 100%;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 20px;

  .kb-content {
    height: 100%;
    display: flex;
    gap: 20px;
    background: #fff;
    border-radius: 4px;
    padding: 20px;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);

    .file-list {
      width: 250px;
      border-right: 1px solid #e6e6e6;
      overflow: auto;

      .kb-header {
        padding: 16px 20px;
        border-bottom: 1px solid #e6e6e6;
        display: flex;
        justify-content: space-between;
        align-items: center;

        .kb-title {
          display: flex;
          align-items: center;
          gap: 8px;
          color: #606266;
          font-size: 14px;
          font-weight: 500;
        }

        .create-btn {
        }
      }

      .file-menu {
        border-right: none;

      }

      .search-box {
        padding: 10px 20px;
        border-bottom: 1px solid #e6e6e6;

      }
    }

    .content-area {
      flex: 1;
      display: flex;
      background: #f5f5f5;
      border-radius: 4px;
      overflow: hidden;

      .editor-container {
        width: 100%;
        height: 100%;
        display: flex;
        flex-direction: column;
        background: #fff;

        .editor-header {
          display: flex;
          justify-content: space-between;
          align-items: center;
          padding: 12px 20px;
          border-bottom: 1px solid #e6e6e6;
          background-color: #fff;

          .file-info {
            display: flex;
            align-items: center;
            gap: 8px;

            .file-name {
              font-size: 16px;
              color: #303133;
              font-weight: 500;
            }
          }

          .editor-actions {
            display: flex;
            gap: 12px;
          }
        }

        .vditor {
          flex: 1;
          height: 0; // 让 flex: 1 生效
        }
      }
    }
  }
}

:deep(.el-dropdown-link) {
  cursor: pointer;
  display: flex;
  align-items: center;
  color: #909399;
}

.search-dialog-content {
  display: flex;
  flex-direction: column;
  gap: 20px;

  .search-input-area {
    display: flex;
    flex-direction: column;
    gap: 10px;

    .search-options {
      padding: 10px 0;
      border-bottom: 1px solid #e6e6e6;
    }
  }

  .search-results {
    max-height: 400px;
    overflow-y: auto;
    padding: 10px 0;

    .search-result-item {
      padding: 12px;
      border-radius: 4px;
      cursor: pointer;
      transition: background-color 0.3s;

      &:hover {
        background-color: #f5f7fa;
      }

      .result-title {
        display: flex;
        align-items: center;
        gap: 8px;
        color: #409EFF;
        margin-bottom: 4px;
      }

      .result-content {
        color: #606266;
        font-size: 13px;
        margin-left: 24px;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
        text-overflow: ellipsis;
      }
    }
  }
}

.share-dialog-content {
  padding: 20px 0;
  text-align: center;
  
  .share-icon {
    color: #409EFF;
    margin-bottom: 20px;

    &.warning-icon {
      color: #E6A23C;
    }
  }
  
  .share-title {
    font-size: 16px;
    font-weight: 500;
    color: #303133;
    margin-bottom: 10px;
  }
  
  .share-desc {
    font-size: 14px;
    color: #909399;
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

// 添加全局遮罩样式
</style>