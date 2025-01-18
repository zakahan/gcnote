<template>
    <div class="test-container">
        <el-button type="primary" @click="testApi">测试后端接口</el-button>
        <div v-if="response" class="response-box">
            <pre>{{ response }}</pre>
        </div>
    </div>
</template>

<script>
import axios from 'axios';
import { ref } from 'vue';

export default {
    setup() {
        const response = ref(null);

        const testApi = async () => {
            try {
                console.log('开始请求...');
                const res = await axios({
                    method: 'get',
                    url: 'http://127.0.0.1:8086/index/show_indexes',
                    headers: {
                        'token': 'eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzc3OTM3NTksIm5hbWUiOiJzdXBlciIsInN1YiI6IjJlNzVjMzVkLWNjODEtNGQxYy1iN2MwLTYyZTkwMzA4MjFjZSJ9.kxn5MXstASeWXsKvWGoOS3fmolcGdZmMZAHtQYdaYoDL7TePYMXscD79WEkaW-6I03zjNiMZFxnnXKt2Efmv9VgDf1ZFeWc-yhwiTnELrfY7Tza3DJ6CqivrV7w4a8q68TZirPDNsQlCC3BnGVkdCuc-QEVSp4fQU5mHgamSHpl59b24AX0HEv2PNcJ-WQhtu1BBuT-mYIPIeAcetVnLOvzK4B6ZgUDDu2gjiEhmjyVe4y5XoAE033Y33tyuHTSN9EUwug63zwoJitsSqqvZjVIdkCJheZ3Y0WWONjiKlC4Ldv92JTKlEPImQVD-C3v43W8cCUdjPjK85iDDbXBSuQ'
                    }
                });
                console.log('请求完成，响应数据:', res);
                response.value = res.data;
            } catch (error) {
                console.log('请求配置:', error.config);
                console.log('错误信息:', error.message);
                if (error.response) {
                    console.log('响应状态:', error.response.status);
                    console.log('响应头:', error.response.headers);
                    console.log('响应数据:', error.response.data);
                }
                response.value = error.message;
            }
        };

        return {
            testApi,
            response
        };
    }
}
</script>

<style scoped>
.test-container {
    padding: 20px;
}
.response-box {
    margin-top: 20px;
    padding: 15px;
    background: #f5f5f5;
    border-radius: 4px;
    white-space: pre-wrap;
}
pre {
    margin: 0;
    font-family: monospace;
}
</style>