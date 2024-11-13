import { createApp } from 'vue';
import App from './App.vue';

createApp(App).mount('#app');

// import ipcRenderer from 'electron'
//
// ipcRenderer.send("request", 1)
//
// ipcRenderer.on("response", (e, a) => {
//     console.log("ipcRenderer got response", a)
// })

console.log("RENDERER.JS")
// import ipcRenderer from 'electron'
// require("electron")