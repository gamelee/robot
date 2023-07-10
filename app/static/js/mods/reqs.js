layui.define2(function (go, $, form) {
    let mod = {
        name: 'reqs',
        modal_html: 'reqs_modal',
    }
    mod.on_req_select = function () {
        let select = mod.$req_select.val()
        mod.$req_name.val(select)
        mod.json_editor.set({
            title: select,
            name: "SendReq",
            category: "action",
            properties: {
                Server: "hall",
                ID: 0,
                Req: select,
                Body: mod.reqs[select],
            }
        })
    }

    mod.add_req = function (id, server = "hall") {
        let $btn = $(`<div class="layui-btn-group req" req="${id}">
    <button class="layui-btn layui-btn-sm server ${server}">${server}</button>
    <button class="layui-btn layui-btn-sm title">${id}</button>
    <button class="layui-btn layui-btn-sm edit"><span class="iconfont">&#xe60f;</span></button>
    <button class="layui-btn layui-btn-sm del "><span class="iconfont">&#xeafb;</span></button>
</div>`)
        $btn.find(".title").on('click', mod.on_send_req)
        $btn.find(".edit").on('click', mod.on_edit_req)
        $btn.find(".del").on('click', mod.on_del_req)
        mod.$reqs.prepend($btn)
    }

    mod.on_add_req = function () {
        layer.open_left('添加请求', mod.$modal, (idx) => {
            if (mod.$req_name.val() in mod.file) {
                layer.msg("保存失败，已存在同名请求，改个名字试试", {icon: 2})
                return
            }
            let id = $.trim(mod.$req_name.val())
            let req = mod.json_editor.get()
            req.id = id
            mod.file[id] = req
            mod.add_req(id, req.properties.Server)
            layer.msg("保存成功", {icon: 1})
            layer.close(idx)
        })
    }


    mod.on_send_req = function () {
        let id = $(this).parent().attr('req')
        mod.send_req(mod.file[id])
    }

    mod.send_req = function (req, body, server = "hall", reqID = 0, title = "发送请求") {
        if (typeof req !== 'object') {
            req = {
                title: title,
                name: "SendReq",
                category: "action",
                properties: {
                    Server: server,
                    Body: body,
                    Req: req,
                    ID: reqID,
                }
            }
        }
        return go.send(req).catch(e => {
            layer.msg(e, {icon: 2})
        })
    }

    mod.send_world_chat = function (content) {
        return mod.send_req('WorldChat', {
            "Content": content,
        })
    }

    mod.on_edit_req = async function () {

        let id = $(this).parent().attr('req')
        let config = await $.json_editor(mod.file[id], '修改参数: ' + id)
        if (config !== false) {
            (mod.file[id] = config)
            let $server = $($(this).parent().children()[0])
            let server = $server.text()
            $server.removeClass(server).addClass(mod.file[id].properties.Server).text(mod.file[id].properties.Server)
        }
    }


    mod.on_del_req = async function () {
        let id = $(this).parent().attr('req')
        let del = await layer.choose("确认删除 " + id + "？")
        if (del === 0) {
            $(this).parent().remove()
            delete mod.file[id]
        }

    }

    mod.$modal = $(`<div id="add-reqs" class="layui-form" lay-filter="reqs" style="display: none"></div>`, {})
    $('body').append(mod.$modal)

    mod.init = async function ($reqs, $panel) {
        mod.$reqs = $reqs

        $panel.append_btn(`📨`, '创建快捷请求', mod.on_add_req)

        const rsp = await fetch(ddchess.path.mod_html("reqs_modal"))
        const html = await rsp.text()
        mod.$modal.append($(html))

        mod.$req_select = mod.$modal.find('select')
        mod.$req_name = mod.$modal.find('input')
        mod.json_editor = new JSONEditor(mod.$modal.find(".json-editor")[0], {mode: 'tree', mainMenuBar: false})
        form.on('select(req)', mod.on_req_select)
        mod.reqs = await go.reqs
        mod.$req_select.children().remove()
        let keys = Object.keys(mod.reqs)
        keys.sort()
        keys.forEach(function (v) {
            mod.$req_select.append($('<option>', {text: v}))
        })
        form.render(null, 'reqs')
        mod.on_req_select()


        mod.file = await go.open_json(ddchess.path.mod_config(mod.name), {})
        Object.keys(mod.file).forEach(function (idx) {
            mod.add_req(idx, mod.file[idx].properties.Server)
        })
    }
    return mod
})

