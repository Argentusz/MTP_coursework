import { createApp } from 'vue';
import App from './App.vue';
import { InitMtpEditor } from "./js/mtpLang.js";

const app = createApp(App)
InitMtpEditor()

app.mount('#app');