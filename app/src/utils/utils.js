export default class Connector {
    constructor(ipc) {
        if (Connector._instance) {
            return Connector._instance
        }

        console.log("Initializing Connector")

        Connector._instance = this

        this.state = ""
        this.error = ""
        this._ipc = ipc
        this._ipc.send("connect")
        this._ipc.on("update", (e, a) => {
            this.state = a
            console.log("update")
            console.log(a)
        })
        this._ipc.on("error", (e, a) => {
            this.error = a
            console.log("error")
            console.log(a)
        })
        return this
    }

    request(command) {
        this._ipc.send("request", command)
    }
}