layui.define2(function (go) {

    let mod = {
        name: 'log_config',
    }

    mod.config_file = ddchess.path.mod_config(mod.name)


    mod.init = async function ($panel) {
        mod.$dom = $panel
        mod.file = await go.open_json(mod.config_file, ddchess.default.log)
        mod.log_skip = await ddchess.dao("log_skip", {})
    }

    mod.add_skip = id => mod.log_skip[id] = true
    mod.is_skip = id => id in (mod.log_skip ?? {})
    mod.del_skip = id => delete mod.log_skip[id]
    mod.replace_skip = function (skip) {
        if (skip === false) {
            return
        }
        for (let i = 0, ks = Object.keys(mod.log_skip), l = ks.length; i < l; i++) {
            mod.del_skip(ks[i])
        }
        Object.assign(mod.log_skip, skip)
    }

    mod.get_filter = id => mod.file[id]
    mod.has_filter = id => id in mod.file
    mod.set_filter = (id, code) => mod.file[id] = code
    mod.del_filter = id => delete mod.file[id]


    let extract_msg = (msg) => {
        for (const k in msg) {
            if (k in {Seq: 1, ErrCode: 1, ErrMsg: 1,}) continue
            msg = msg[k]
            break
        }
        return msg
    }

    mod.filter = function (id, msg, $dom = false) {
        if (!mod.file?.[id]) return false
        if ($dom !== false) {
            $dom = $dom.find(".content")
        }
        let data = msg // eval 中使用
        if (data.ErrMsg) {
            $dom ? $dom.text(data.ErrMsg).css("background-color", "var(--color-0)").css("color", "var(--color-4)") :
                console.error(data.ErrMsg)
            return
        }
        try {
            eval(mod.file[id])
        } catch (e) {
            console.error(id, e)
            return e
        }
        return false
    }

    function format(val) {
        if (Array.isArray(val)) {
            if (val.length > 0) {
                return [format(val[0])]
            }
            return 'array'
        } else if (typeof val === 'object') {
            let out = {}
            for (const k in val) {
                out[k] = format(val[k])
            }
            return out
        }
        return val
    }


    mod.default_filter = function (msg) {
        return `// 处理函数
data || (data = {}) // 原始消息数据
$dom || ($dom = false)
let msg = extract_msg(data) // 提取出内层数据 
// 开始处理
$dom.htmlJSON(msg)
`
    }

    return mod
})

