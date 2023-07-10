layui.define2(function (go, $, form, notice) {

    let mod = {
        name: 'server',
    }

    mod.init = async function ($select, $panel) {
        mod.$select = $select
        $panel.append_btn(`🗑️`, '删除服务器配置', mod.on_delete_config)
        $panel.append_btn(`✍️️`, '编辑服务器配置', mod.on_edit_config)
        $panel.append_btn(`➕`, '添加服务器配置', mod.on_add_config)
        form.on('select(servers)', () => mod.file.select = mod.$select.val())
        mod.file = await go.open_json(ddchess.path.mod_config(mod.name), ddchess.default.server)
        mod.render_select()
    }

    mod.render_select = function () {
        if (!mod.file.servers || Object.keys(mod.file.servers).length === 0) {
            notice.error`服务器配置为空, 请尝试删除配置文件`
            return
        }
        mod.$select.append(Object.keys(mod.file.servers).map(idx => $(`<option>`, {text: idx, value: idx})))

        if (!mod.file.select || !mod.file.servers[mod.file.select]) {
            mod.file.select = Object.keys(mod.file.servers)[0]
        }
        mod.$select.val(mod.file.select)
        form.render('select')
    }


    mod.on_edit_config = async () => {
        let config = await $.json_editor(mod.file.servers[mod.file.select], '编辑服务配置: ' + mod.file.select)
        config === false || (mod.file.servers[mod.file.select] = config)
    }

    mod.on_delete_config = async () => {
        if (Object.keys(mod.file.servers).length <= 1) {
            notice.warning`再删就一个都不剩了`
            return
        }

        let del = await layer.choose("确认删除服务器" + mod.file.select + " 的配置么？")
        if (del === 0) {
            mod.$select.find("option[value='" + mod.file.select + "']").remove()
            delete mod.file.servers[mod.file.select]
            mod.file.select = Object.keys(mod.file.servers)[0]
            mod.$select.val(mod.file.select)
            form.render('select')
            notice.success`删除成功`;
        }
    }

    mod.on_add_config = async () => {
        let name = await layer.input(`请输入服务器名`)
        if (name === false) return
        name = $.trim(name)
        if (name === "") {
            notice.warning`配置名不能为空`
            return
        }
        if (name in mod.file.servers) {
            notice.warning`配置已存在`
            return
        }

            mod.file.servers[name] = mod.file.servers[mod.file.select]
        mod.file.select = name
        mod.$select.append($('<option>', {text: name}))
        mod.$select.val(mod.file.select)
        form.render('select')
        notice.success`保存成功`
    }

    return mod
})

