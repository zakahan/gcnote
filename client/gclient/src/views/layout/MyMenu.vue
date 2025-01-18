<template>
  <div class="menu-container">
    <el-menu
      :default-active="activeMenu"
      class="el-menu-vertical"
      :collapse="isCollapse"
      @select="handleSelect"
    >
      <el-menu-item index="home">
        <el-icon><HomeFilled /></el-icon>
        <span>开始</span>
      </el-menu-item>

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

      <el-sub-menu index="kb">
        <template #title>
          <el-icon><Folder /></el-icon>
          <span>知识库</span>
        </template>
        <el-menu-item-group v-if="menuItems.length > 0">
          <div
            v-for="item in menuItems"
            :key="item.IndexId"
            class="kb-item-wrapper"
            @mouseenter="hoveredItem = item.IndexId"
            @mouseleave="hoveredItem = null"
          >
            <el-menu-item 
              :index="item.IndexId"
              @click="handleKbClick(item)"
            >
              <el-icon><Document /></el-icon>
              <span>{{ item.IndexName }}</span>
            </el-menu-item>
            <el-dropdown
              v-show="hoveredItem === item.IndexId"
              trigger="click"
              class="kb-actions"
              placement="bottom-end"
              :teleported="false"
              @command="(command) => handleCommand(command, item)"
            >
              <el-icon class="more-icon"><MoreFilled /></el-icon>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="rename">
                    <el-icon><EditPen /></el-icon>重命名
                  </el-dropdown-item>
                  <el-dropdown-item command="delete" style="color: #f56c6c;">
                    <el-icon><Delete /></el-icon>删除
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </el-menu-item-group>
        <el-empty v-else description="暂无知识库" />
      </el-sub-menu>

      <el-menu-item index="recycle" @click="router.push('/recycle')">
        <el-icon><Delete /></el-icon>
        <span>回收站</span>
      </el-menu-item>

      <!-- 搜索对话框 -->
      <el-dialog
        v-model="showSearchDialog"
        title="搜索文档"
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
          </div>

          <div class="search-results" v-loading="searching">
            <template v-if="searchResults.length > 0">
              <div 
                v-for="result in searchResults" 
                :key="result.KBFileId"
                class="search-result-item"
                @click="handleResultClick(result)"
              >
                <div class="result-title">
                  <el-icon><Document /></el-icon>
                  <span>{{ result.KBFileName }}</span>
                </div>
              </div>
            </template>
            <el-empty v-else-if="!searching" description="暂无搜索结果" />
          </div>
        </div>
      </el-dialog>

      <!-- 删除知识库确认对话框 -->
      <el-dialog
        v-model="showDeleteDialog"
        title="删除知识库"
        width="30%"
        :close-on-click-modal="false"
      >
        <div class="dialog-content">
          <el-icon class="dialog-icon warning-icon" :size="50"><Warning /></el-icon>
          <div class="dialog-title">确定要删除此知识库吗？</div>
          <div class="dialog-desc">删除后，该知识库下的所有文件将移至回收站</div>
        </div>
        <template #footer>
          <span class="dialog-footer">
            <el-button @click="showDeleteDialog = false">取消</el-button>
            <el-button type="danger" @click="confirmDelete">确定删除</el-button>
          </span>
        </template>
      </el-dialog>
    </el-menu>

    <!-- 重命名对话框 -->
    <el-dialog
      v-model="renameDialog"
      title="重命名知识库"
      width="30%"
      :close-on-click-modal="false"
      destroy-on-close
    >
      <el-form :model="renameForm" label-width="80px">
        <el-form-item label="新名称">
          <el-input v-model="renameForm.name" placeholder="请输入新的知识库名称" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="renameDialog = false">取消</el-button>
          <el-button type="primary" @click="submitRename">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, defineExpose } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { HomeFilled, Folder, Document, MoreFilled, EditPen, Delete, Search, Warning } from '@element-plus/icons-vue'
import axios from 'axios'

const router = useRouter()
const route = useRoute()
const menuItems = ref([])
const isCollapse = ref(false)
const activeMenu = ref('home')
const hoveredItem = ref(null)
const searchKeyword = ref('')
const showSearchDialog = ref(false)
const searchResults = ref([])
const searching = ref(false)

// 重命名相关
const renameDialog = ref(false)
const renameForm = ref({ name: '', id: '' })

// 搜索相关方法
const handleSearch = async () => {
  if (!searchKeyword.value.trim()) {
    ElMessage.warning('请输入搜索关键词')
    return
  }

  searching.value = true
  searchResults.value = []

  try {
    const token = localStorage.getItem('token')
    const response = await axios.post(
      'http://localhost:8086/index/search_file',
      {
        kb_file_name: searchKeyword.value,
        is_fuzzy_search: true
      },
      {
        headers: {
          'token': token,
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



const handleResultClick = (result) => {
  if (result.IndexId && result.KBFileId) {
    showSearchDialog.value = false
    router.push({
      name: 'knowledge-base',
      params: { 
        index_id: result.IndexId,
        file_id: result.KBFileId,
        file_name: result.KBFileName
      }
    })
  }
}

// 原有的方法保持不变
const handleCommand = (command, item) => {
  if (command === 'rename') {
    renameForm.value = {
      name: item.IndexName,
      id: item.IndexId
    }
    renameDialog.value = true
  } else if (command === 'delete') {
    itemToDelete.value = item
    showDeleteDialog.value = true
  }
}

const submitRename = async () => {
  try {
    const token = localStorage.getItem('token')
    const response = await axios.post('http://localhost:8086/index/rename_index', {
      dest_index_name: renameForm.value.name,
      index_id: renameForm.value.id
    }, {
      headers: { token }
    })

    if (response.data.Code === 0) {
      ElMessage.success('重命名成功')
      renameDialog.value = false
      await fetchMenuItems()
    } else {
      ElMessage.error(response.data.Msg || '重命名失败')
    }
  } catch (error) {
    console.error('重命名失败:', error)
    ElMessage.error('重命名失败')
  }
}

const showDeleteDialog = ref(false)
const itemToDelete = ref(null)

// 确认删除
const confirmDelete = async () => {
  try {
    const token = localStorage.getItem('token')
    const response = await axios.post(
      'http://localhost:8086/index/delete',
      {
        index_id: itemToDelete.value.IndexId
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
      await fetchMenuItems()
    } else {
      ElMessage.error(response.data.Msg || '删除失败')
    }
  } catch (error) {
    console.error('删除失败:', error)
    ElMessage.error('删除失败')
  }
}

const handleKbClick = (item) => {
  console.log('点击知识库:', item)
  router.push({
    name: 'knowledge-base',
    params: { index_id: item.IndexId }
  })
}

const handleSelect = (index) => {
  if (index === 'home') {
    router.push('/home')
  } else {
    // 直接使用index作为index_id
    router.push({
      name: 'knowledge-base',
      params: { index_id: index }
    })
  }
}

const fetchMenuItems = async () => {
  try {
    const token = localStorage.getItem('token')
    if (!token) {
      router.push('/login')
      return
    }

    const response = await axios.get('http://localhost:8086/index/show_indexes', {
      headers: {
        token: token
      }
    })

    if (response.data && response.data.Code === 0) {
      menuItems.value = response.data.Data
    } else {
      ElMessage.error(response.data.Msg || '获取知识库列表失败')
    }
  } catch (error) {
    console.error('获取知识库列表失败:', error)
    ElMessage.error('获取知识库列表失败')
  }
}

// 暴露刷新方法给父组件
defineExpose({
  refresh: fetchMenuItems
})

onMounted(() => {
  fetchMenuItems()
  // 根据当前路由设置活动菜单项
  const currentPath = route.path
  if (currentPath.includes('/kb/')) {
    activeMenu.value = 'kb-' + route.params.id
  } else {
    activeMenu.value = 'home'
  }
})
</script>

<style lang="less" scoped>
.el-menu-vertical {
  height: 100%;
  border-right: none;

  .el-menu-item [class^="el-icon-"] {
    margin-right: 5px;
    width: 24px;
    text-align: center;
    font-size: 18px;
  }

}

.search-box {
  padding: 10px 20px;
  border-bottom: 1px solid #e6e6e6;

}

.kb-item-wrapper {
  position: relative;
  display: flex;
  align-items: center;

  .kb-actions {
    position: absolute;
    right: 20px;
    cursor: pointer;
    z-index: 10;
    
    .more-icon {
      font-size: 16px;
      color: #909399;
      
      &:hover {
        color: #409EFF;
      }
    }

  }
}

.search-dialog-content {
  display: flex;
  flex-direction: column;
  gap: 20px;

  .search-input-area {
    display: flex;
    flex-direction: column;
    gap: 10px;

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

</style>