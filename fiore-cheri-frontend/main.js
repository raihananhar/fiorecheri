const { app, BrowserWindow } = require('electron');

function createWindow() {
    const win = new BrowserWindow({
        width: 800,
        height: 600,
        webPreferences: {
            nodeIntegration: true, // Fix require() issue
            contextIsolation: false // Tambahkan ini biar axios bisa jalan
        }
    });

    win.loadFile('index.html');
}

app.whenReady().then(createWindow);
