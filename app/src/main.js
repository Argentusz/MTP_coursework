import * as child_process from "child_process";

const { app, BrowserWindow } = require('electron');
const path = require('node:path');

import Squirrel from "electron-squirrel-startup"

// Handle creating/removing shortcuts on Windows when installing/uninstalling.
if (Squirrel) {
  app.quit();
}

const createWindow = () => {
  // Create the browser window.
  const win = new BrowserWindow({
    webPreferences: {
      preload: path.join(__dirname, 'preload.js'),
    },
  });

  win.maximize()

  // and load the index.html of the app.
  if (MAIN_WINDOW_VITE_DEV_SERVER_URL) {
    win.loadURL(MAIN_WINDOW_VITE_DEV_SERVER_URL);
  } else {
    win.loadFile(path.join(__dirname, `../renderer/${MAIN_WINDOW_VITE_NAME}/index.html`));
  }

  // Open the DevTools.
  // win.webContents.openDevTools();
};

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.whenReady().then(() => {
  createWindow()

  // On OS X it's common to re-create a window in the app when the
  // dock icon is clicked and there are no other windows open.
  app.on('activate', () => {
    if (BrowserWindow.getAllWindows().length === 0) {
      createWindow();
    }
  });
});

// Quit when all windows are closed, except on macOS. There, it's common
// for applications and their menu bar to stay active until the user quits
// explicitly with Cmd + Q.
app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit();
  }
});

// In this file you can include the rest of your app's specific main process
// code. You can also put them in separate files and import them here.
import { spawn } from "node:child_process";
class MtpBridge {
  constructor() {
    this._renderer = null
    this._server = null
    this._stdoutbuff = null
    this._stdoutreading = false
    this._filePathBuff = null
    this._fileDataBuff = null
  }
  connectRenderer(renderer) {
    this._renderer = renderer
  }
  connectServer(path, flags) {
    this._server = spawn(path, [...flags, "-marshall=JSON", "-base="])
    this._server.stderr.on("data", data => this._renderer?.send("error", `${data}`))
    this._server.stdout.on("data", data => {
      const output = `${data}`

      this._stdoutreading
        ? this._stdoutbuff += output
        : this._stdoutbuff = output

      this._stdoutreading = !output.endsWith("\x00\n")

      if (!this._stdoutreading) {
        this._stdoutbuff = this._stdoutbuff.substring(0, this._stdoutbuff.length - 2)
        this._renderer?.send("update", this._stdoutbuff)
      }
    })
  }
  setFileBuffer(filePath, data) {
    this._filePathBuff = filePath
    this._fileDataBuff = data
  }
  disconnectServer() {
    this._server?.kill()
    this._stdoutbuff = null
    this._renderer?.send("update", null)
  }
  send(command) {
    this._server?.stdin.write(command)
  }
  ping() {
    this._renderer?.send("update", this._stdoutbuff)
    if (this._filePathBuff) {
      this._renderer?.send("f-opened", { filePath: this._filePathBuff, data: this._fileDataBuff })
    }
  }
}

const bridge = new MtpBridge()

const { ipcMain, dialog } = require("electron")
const fs = require("fs")

ipcMain.on("greet", (e, a) => bridge.connectRenderer(e.sender))
ipcMain.on("connect", (e, a) => {
  dialog.showOpenDialog({ properties: ['openFile', 'multiSelections'] }).then(r => {
    if (r.canceled || r.filePaths.length === 0) return
    bridge.connectServer(r.filePaths[0], a)
  })
})
ipcMain.on("disconnect", (e, a) => bridge.disconnectServer())
ipcMain.on("request", (e, a) => {
  if (!a.endsWith("\n")) a += "\n"
  bridge.send(a)
})
ipcMain.on("ping", (e, a) => bridge.ping())

ipcMain.on("f-open", (e, a) => {
  dialog.showOpenDialog({ properties: ['openFile'] })
      .then(r => {
        if (r.canceled || r.filePaths.length === 0) return
        const filePath = r.filePaths[0]
        openFile(e, filePath)
      })
})
ipcMain.on("f-save", (e, a) => {
  if (a.filePath) {
    writeFile(e, a)
    return
  }

  newFile(e, a).then(filePath => {
    if (!filePath) {
      return
    }

    writeFile(e, { filePath, data: a.data }, (err => {
      if (err) {
        e.returnValue = false
        return
      }

      openFile(e, filePath)
      e.sender.send("f-saved")
      e.returnValue = true
    }))
  })
})
ipcMain.on("f-new", (e, a) => {
  newFile(e, a).then(filePath => {
    if (filePath) {
      openFile(e, filePath)
    }
  })
})

const newFile = (e, a) => {
  return dialog.showSaveDialog({}).then(file => {
    if (file.canceled) return null
    fs.open(file.filePath, "w", () => {})
    return file.filePath
  })
}

const writeFile = (e, a, callback = null) => {
  callback = callback || (err => {
    if (!err) e.sender.send("f-saved")
    e.returnValue = !err
  })
  bridge.setFileBuffer(a.filePath, a.data)
  fs.writeFile(a.filePath, a.data, callback)
}

const openFile = (e, filePath) => {
  fs.readFile(filePath, "utf-8", (err, data) => {
    if (err) e.sender.send("error", err)
    bridge.setFileBuffer(filePath, data)
    e.sender.send("f-opened", { filePath, data })
  })
}