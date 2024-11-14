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
    this._process = spawn("mtp/mtp_darwin", ["-intr", "-trace" ,"-sudo", "-marshall=JSON", "-base=./mtp/projects"])
    this._renderer = null

    this._process.stdout.on("data", (data) => {
      if (!this._renderer) {
        return
      }

      this._renderer.send("update", `${data}`)
    })

    this._process.stderr.on("data", (data) => {
      if (!this._renderer) {
        return
      }

      this._renderer.send("error", `${data}`)
    })
  }

  connectRenderer(renderer) {
    this._renderer = renderer
  }

  send(command) {
    if (!command.endsWith("\n")) {
      command += "\n"
    }

    this._process.stdin.write(command)
  }
}

const bridge = new MtpBridge()

const { ipcMain } = require("electron")
ipcMain.on("connect", (e, a) => {
  bridge.connectRenderer(e.sender)
})

ipcMain.on("request", (e, a) => {
  bridge.send(a)
})
