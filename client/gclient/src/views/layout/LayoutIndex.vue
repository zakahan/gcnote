<template>
  <el-container class="layout-container">
    <el-header class="header">
      <div class="logo">
        <h2>知识库系统</h2>
      </div>
      <div class="user-info">
        <el-dropdown @command="handleCommand">
          <span class="el-dropdown-link">
            {{ username }}<el-icon class="el-icon--right"><arrow-down /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="profile">个人信息</el-dropdown-item>
              <el-dropdown-item command="logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </el-header>
    <el-container>
      <el-aside width="200px" class="aside">
        <MyMenu ref="menuRef"></MyMenu>
      </el-aside>
      <el-main class="main">
        <router-view></router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ArrowDown } from '@element-plus/icons-vue'
import MyMenu from './MyMenu.vue'
// import MyContent from './MyContent.vue'

const router = useRouter()
const username = ref(localStorage.getItem('username') || '用户')
const menuRef = ref(null)

// const refreshMenu = () => {
//   menuRef.value?.refresh()
// }

const handleCommand = (command) => {
  if (command === 'logout') {
    localStorage.removeItem('token')
    localStorage.removeItem('username')
    router.push('/login')
    ElMessage.success('已退出登录')
  } else if (command === 'profile') {
    // TODO: 跳转到个人信息页面
    ElMessage.info('功能开发中')
  }
}
</script>

<style lang="less" scoped>
.layout-container {
  height: 100vh;
  
  .header {
    background-color: #fff;
    border-bottom: 1px solid #e6e6e6;
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0 20px;
    
    .logo {
      h2 {
        margin: 0;
        color: #409EFF;
      }
    }
    
    .user-info {
      .el-dropdown-link {
        cursor: pointer;
        display: flex;
        align-items: center;
        color: #606266;
        
        &:hover {
          color: #409EFF;
        }
      }
    }
  }
  
  .aside {
    background-color: #fff;
    border-right: 1px solid #e6e6e6;
  }
  
  .main {
    background-color: #f5f5f5;
    padding: 20px;
  }
}
</style>
 