// See the Electron documentation for details on how to use preload scripts:
// https://www.electronjs.org/docs/latest/tutorial/process-model#preload-scripts

import { contextBridge, ipcRenderer } from "electron"

ipcRenderer.send("greet")

contextBridge.exposeInMainWorld("mtpAPI", {
    connect: (flags) => ipcRenderer.send("connect", flags),
    disconnect: () => ipcRenderer.send("disconnect"),
    request: (command) => ipcRenderer.send("request", command),
    ping: () => ipcRenderer.send("ping"),

    onUpdate: (callback) => ipcRenderer.on("update", (e, a) => callback(a)),
    onError: (callback) => ipcRenderer.on("error", (e, a) => callback(a)),
})

contextBridge.exposeInMainWorld("fileAPI", {
    new: () => ipcRenderer.send("f-new"),
    open: () => ipcRenderer.send("f-open"),
    save: (filePath, data) => ipcRenderer.sendSync("f-save", { filePath, data }),

    onOpened: (callback) => ipcRenderer.on("f-opened", (e, a) => callback(a)),
    onSaved: (callback) => ipcRenderer.on("f-saved", (e, a) => callback(a))
})