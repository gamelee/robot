layui.define2(function (layer, notice) {
    const code_success = 1
    const code_failed = 2

    let mod = {
        name: 'go',
        seq: 0,
        reg: {},
        files:{},
    }

    window.GO = window.GO || false
    if (!GO) {
        layer.alert(mod.name + " 进程初始化失败")
        return
    }

    mod.on = function (action, fn) {
        if (mod.reg[action]) {
            return
        }
        mod.reg[action] = fn
    }
    // 调用 golang
    mod.call = async function (act, arg, msg = "") {
        let req = {seq: this.seq++, action: act, msg: msg, arg: arg,}

        return GO(req).then(rsp => {
            let info = typeof arg === 'string' ? arg : arg[0]
            info = !info ? 'no_arg' : info
            console.debug(req.action, info, {arg, req, rsp})
            if (rsp.seq !== req.seq) {
                return new Promise((resolve, reject) => {
                    reject("seq 不相等")
                })
            }
            if (rsp.code !== code_success) {
                return new Promise((resolve, reject) => {
                    reject(rsp.msg);
                })
            }
            return rsp.data
        })
    }
    // 此方法 golang 会调用
    window.CallJS = function ({action, seq, arg}) {

        if (action === 'event' && (arg.id === 1 || arg.id === 1001)) {
            return rsp
        }

        let rsp = {seq: seq, code: code_success}
        let fn = mod.reg[action]
        if (!fn) {
            rsp.code = code_failed
            rsp.msg = "unknown action " + action
            console.error('js', action, "未注册处理函数", arg)
            return rsp
        }
        rsp.code = code_success

        rsp.data = fn(arg)
        // 返回 true golang 会特殊处理
        console.debug(arg.type, arg.from, {arg, rsp})
        if (action === 'event' && (arg.type === "rsp" || arg.type === 'ntf')) {
            rsp.data = true
        }

        return rsp;
    }


    mod.reqs = mod.call("req.all", "", "获取所有CS请求详情").then(x => x)

    mod.send = async (req) => mod.call("req.send", req, "发送请求")

    mod.go_version = async (url) => mod.call("golang.version", url, "golang 版本")

    mod.nodes = async () => mod._nodes || mod.call("nodes", '', "节点").then(data => mod._nodes = data)

    mod.run = async (id, data) => {
        return mod.call("robot.run", [id, data], "运行行为树").then((rst) => {
            notice.success("运行成功")
        }).catch(e => {
            notice.error("运行失败")
        })
    }

    mod.read = async (file) => mod.call("file.read", file, "读取文件")
    mod.read_json = async (file) => {
        let buf = await mod.read(file)
        return JSON.parse(buf)
    }

    mod.write = async (file, content) => mod.call("file.write", [file, content], "写入文件")
    mod.write_json = async (name, obj) => mod.write(name, JSON.stringify(obj, undefined, "    "))

    let loading = Symbol('loading')
    mod.open_json = async function (file, dft) {
        let i = 0
        if (mod.files[file]) {
            return new Promise((resolve, reject) => {
                let fn = function () {
                    i++
                    if (mod.files[file] !== loading) {
                        resolve(mod.files[file])
                    }
                }
                if (i > 100) {
                    console.error("timeout", 'file.read wait', file)
                    reject("timeout", 'file.read wait', file)
                }
                if (i === 0) {
                    fn()
                }
                setTimeout(fn, 10)
            })
        }
        mod.files[file] = loading
        return mod.read_json(file).then(obj => {
            mod.files[file] = new ObjectProxy(obj, (dat) => mod.write_json(file, dat))
            return mod.files[file]
        }).catch(e => {
            mod.files[file] = new ObjectProxy(dft, (dat) => mod.write_json(file, dat))
            return mod.write_json(file, dft).then(()=> mod.files[file])
        })
    }

    return mod
});
