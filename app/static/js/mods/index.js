layui.define2(function ($, init, log_view, menu, server) {

    let vars = {
        server() {
            return server.file.servers[server.file.select]
        },
        msg:{},
    }

    window.vars = new Proxy({}, {
        get(target, p) {
            if (typeof vars[p] === "function") {
                return vars[p]()
            }
            return vars[p]
        }
    })

    let mod = {name: 'index',}
    menu.init()
    log_view.init()
    return mod
})

