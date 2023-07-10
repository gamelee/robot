layui.define2(function ($, layer, go, log_config, tree, notice, menu) {

    let mod = {name: 'log_view',}

    go.on("event", function (ev) {
        if (ev.type === "sys") {
            console.log("server", ev)
            return
        }
        mod.add_log(ev)
        if (ev.type === "rsp" || ev.type === 'ntf') {
            // 可以在 golang中访问
            window.vars.msg = Object.assign(window.vars.msg, ev.data)
        }
    })

    mod.js_error = ""

    mod.filter = async function (id, data, $dom= false) {
        let err = log_config.filter(id, data, $dom)
        if (err === false) return true
        if(mod.js_error !== "") return true // 有错误未处理
        mod.js_error = err

        let flag = await layer.choose("js 代码有报错,是否修改")
        if (flag !== 0) return true

        let code = log_config.get_filter(id) || log_config.default_filter(data)
        code = await $.js_editor(code, 'js 代码编辑器: 消息输出格式 - ' + id)
        if (code === false) {
            return true
        }
        if (code === true) {
            let del = await layer.choose("确认删除 " + id + " 的配置么？")
            del || log_config.del_filter(id)
            notice.success("删除成功")
            return true
        }
        log_config.set_filter(id, code)

        err = await mod.filter(id, data, $dom)
        if ( err !== false) {
            mod.js_error = err
            notice.error("代码有报错")
            return false
        }
        mod.js_error = ""
        notice.success("保存成功")
        return true
    }


    mod.add_log = function (info) {
        info.id_pretty = info.id_pretty || info.id
        let id = info.from + "_" + info.id_pretty.toLowerCase()
        if (log_config.is_skip(id)) {
            mod.filter(id, info.data, false)
            return
        }


        let $dom = $(`<div id="${id}" class="log-item ${info.id} ${info.type} ${info.from}">
    <div class="layui-btn-group head layui-anim-down">
        <button type="button" class="layui-btn layui-btn-xs server">${info.from}</button>
        <button type="button" class="layui-btn layui-btn-xs layui-btn-primary id">${info.id_pretty}</button>
        <button type="button" class="layui-btn layui-btn-xs layui-btn-primary skip">🗑️</button>
        <button type="button" class="layui-btn layui-btn-xs layui-btn-primary edit">🖋️</button>
    </div>
    <div class="content"></div>
</div>`)
        $dom.data(info)
        $dom.find(".id").on("click", mod.on_view_log)
        $dom.find(".server").text(info.from)
        $dom.find(".skip").on('click', mod.log_skip)
        $dom.find(".edit").on('click', mod.log_edit)
        mod.$dom.prepend($dom)
        if (info?.data?.ErrMsg) {
            $dom.find(".content").html(info?.data?.ErrMsg)
        } else{
            mod.filter(id, info.data, $dom)
        }
    }

    mod.on_view_log = function () {
        $.json_viewer($(this).parent().parent().data(), "详情", false)
    }

    mod.init = function () {
        mod.$dom = $(`<div class="log"></div>`)
        $('body').append(mod.$dom)

        menu.$panel.append_btn(`🧹`, '清除日志', mod.clean_log)
        menu.$panel.append_btn(`🗂️️`, '已过滤日志', async () => {
            let config = await $.json_editor(log_config.log_skip, '编辑日志过滤列表')
            log_config.replace_skip(config)
        })
    }

    mod.log_tree = function ($dom, id, title, items) {
        $dom.append(`<div id="log-tree-${id}" style="display: inline-block"></div>`)
        tree.render({
            elem: '#log-tree-' + id,
            data: [{
                id: id,
                title: title,
                spread: false,
                children: items
            }],
            edit: [],
        })
    }


    mod.log_skip = function ($dom) {
        let id = $(this).parent().parent().attr("id")
        log_config.add_skip(id)
        mod.$dom.find(".log-item[id=" + id + "]").remove()
    }


    mod.log_edit = async function ($dom) {
        let $log = $(this).parent().parent()
        let info = $log.data()
        let id = $log.attr("id")

        let code = log_config.get_filter(id) || log_config.default_filter(info.data)
        code = await $.js_editor(code, 'js编辑器: 消息输出格式 - ' + info.id_pretty)
        if (code === false) {
            return
        }
        if (code === true) {
            let del = await layer.choose("确认删除 " + id + " 的配置么？")
            del || log_config.del_filter(id)
            notice.success("删除成功")
            return
        }

        log_config.set_filter(id, code)
        if (log_config.filter(id, info.data, $log) !== false) {
            notice.error("代码有报错")
            return
        }
        notice.success("保存成功")
    }


    mod.clean_log = function () {
        mod.$dom.empty()
    }
    return mod
})

