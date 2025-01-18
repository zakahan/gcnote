import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import "normalize.css"
import "./assets/css/base.css"
// import store from './store'


const app = createApp(App);
app.use(router).mount("#app")
// createApp(App).use(store).use(router).mount('#app')


