// See the Electron documentation for details on how to use preload scripts:
// https://www.electronjs.org/docs/latest/tutorial/process-model#preload-scripts

import { contextBridge, ipcRenderer } from "electron"

ipcRenderer.send("connect")

contextBridge.exposeInMainWorld("mtpAPI", {
    request: (command) => ipcRenderer.send("request", command),
    onUpdate: (callback) => ipcRenderer.on("update", (e, a) => callback(a)),
    onError: (callback) => ipcRenderer.on("error", (e, a) => callback(a)),
})

contextBridge.exposeInMainWorld("fileAPI", {
    open: () => ipcRenderer.send("f-open"),
    save: (filePath, data) => ipcRenderer.send("f-save", { filePath, data }),
    onOpened: (callback) => ipcRenderer.on("f-opened", (e, a) => callback(a)),
    onSaved: (callback) => ipcRenderer.on("f-saved", (e, a) => callback(a))
})