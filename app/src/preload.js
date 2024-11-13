// See the Electron documentation for details on how to use preload scripts:
// https://www.electronjs.org/docs/latest/tutorial/process-model#preload-scripts
console.log("PRELOAD.JS")

import { ipcRenderer } from "electron"
ipcRenderer.send("request", 1)

ipcRenderer.on("response", (e, a) => {
    console.log("ipcRenderer got response", a)
})