<template>
  <div class="login-container">
    <div class="login-box">
      <div class="title">
        <h2>{{ isLogin ? '欢迎登录' : '用户注册' }}</h2>
      </div>
      <el-form :model="formData" :rules="rules" ref="formRef">
        <el-form-item prop="username">
          <el-input v-model="formData.username" placeholder="用户名" prefix-icon="User" />
        </el-form-item>
        <el-form-item v-if="!isLogin" prop="email">
          <el-input v-model="formData.email" placeholder="邮箱" prefix-icon="Message" />
        </el-form-item>
        <el-form-item prop="password">
          <el-input v-model="formData.password" type="password" placeholder="密码" prefix-icon="Lock" show-password />
        </el-form-item>
        <el-form-item v-if="!isLogin" prop="confirmPassword">
          <el-input v-model="formData.confirmPassword" type="password" placeholder="确认密码" prefix-icon="Lock" show-password />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" class="submit-btn" @click="handleSubmit">
            {{ isLogin ? '登录' : '注册' }}
          </el-button>
        </el-form-item>
      </el-form>
      <div class="switch-mode">
        <span @click="isLogin = !isLogin">
          {{ isLogin ? '没有账号？去注册' : '已有账号？去登录' }}
        </span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import axios from 'axios'
// import { User, Lock, Message } from '@element-plus/icons-vue'

const router = useRouter()
const isLogin = ref(true)
const formRef = ref(null)

const formData = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: ''
})

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '长度在 6 到 20 个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { 
      validator: (rule, value, callback) => {
        if (!isLogin.value && value !== formData.password) {
          callback(new Error('两次输入密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    const url = `http://localhost:8086/user/${isLogin.value ? 'login' : 'register'}`
    const requestData = isLogin.value ? {
      user_name: formData.username,
      password: formData.password
    } : {
      user_name: formData.username,
      email: formData.email,
      password: formData.password
    }
    
    const response = await axios.post(url, requestData)
    
    if (response.data.Code === 0) {
      ElMessage.success(isLogin.value ? '登录成功' : '注册成功')
      if (isLogin.value) {
        localStorage.setItem('token', response.data.Data)
        localStorage.setItem('username', formData.username)
        await router.push('/home')
      } else {
        isLogin.value = true
      }
    } else {
      ElMessage.error(response.data.Msg || (isLogin.value ? '登录失败' : '注册失败'))
    }
  } catch (error) {
    console.error(error)
    ElMessage.error(error.message || '操作失败')
  }
}
</script>

<style scoped>
.login-container {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f5f5f5;
}

.login-box {
  width: 400px;
  padding: 40px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.title {
  text-align: center;
  margin-bottom: 30px;
}

.title h2 {
  color: #409EFF;
  font-weight: 500;
}

.submit-btn {
  width: 100%;
}

.switch-mode {
  text-align: center;
  margin-top: 20px;
  color: #409EFF;
  cursor: pointer;
}

.switch-mode span:hover {
  text-decoration: underline;
}
</style>