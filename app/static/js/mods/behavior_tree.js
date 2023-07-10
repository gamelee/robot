layui.define2(function (element, go, $, server, notice, form, tree, reqs) {
    let mod = {name: 'behavior_tree',}

    mod.init = async function ($behavior_tree, $panel) {
        mod.$behavior_tree = $behavior_tree

        $panel.append_btn(`🌳`, '创建行为树', mod.on_add_tree)

        mod.file = await go.open_json(ddchess.path.mod_config(mod.name), ddchess.default.behavior_tree)
        Object.values(mod.file).forEach(mod.$add_tree)
        await mod.init_editor()
    }

    mod.gen_uuid = function () {
        let s = [];
        let hexDigits = "0123456789abcdef";
        for (let i = 0; i < 36; i++) {
            s[i] = hexDigits[Math.floor(Math.random() * 0x10)]
        }
        s[14] = "4";
        s[19] = hexDigits[(s[19] & 0x3) | 0x8]
        s[8] = s[13] = s[18] = s[23] = "-";

        return s.join("")
    }

    mod.$add_tree = function (tree) {
        let $group = $(`<div id="${tree.id}" class="layui-btn-group tree">
    <button class="layui-btn layui-btn-sm run"><span class="iconfont">&#xe6c1;</span></button>
    <button class="layui-btn layui-btn-sm name">${tree.title}</button>
    <button class="layui-btn layui-btn-sm del"><span class="iconfont">&#xeafb;</span></button>
</div>`)
        $group.find('button.run').on('click', mod.run_tree)
        $group.find('button.del').on('click', mod.on_del_tree)
        $group.find('button.name').on("click", mod.edit_tree)
        mod.$behavior_tree.append($group)
    }

    mod.run_tree = function () {
        go.run($(this).parent().attr("id"), server.file.servers[server.file.select])
    }

    mod.edit_tree = function () {
        mod.editor($(this).parent().attr("id"), mod.file)
    }

    mod.new_tree = function (name) {
        let tree = {
            id: mod.gen_uuid(),
            title: name,
            root: mod.gen_uuid(),
            properties: {},
            nodes: {},
        }
        let root = {
            id: tree.root,
            name: "Sequence",
            category: "composite",
            title: "根节点",
            properties: {},
            children: [],
        }

        tree.nodes[root.id] = root
        return tree
    }

    mod.on_add_tree = function () {
        layer.prompt("请输入行为树名称", (name, idx, elem) => {
            layer.close(idx);
            name = $.trim(name)
            if (name === "") {
                notice.error`名称不能为空`
                return
            }

            if (Object.values(mod.file).to_object(v => [v.title, true])[name]) {
                notice.error`配置已存在`
                return
            }

            let tree = mod.new_tree(name)
            mod.file[tree.id] = tree
            mod.$add_tree(tree)
            notice.success`保存成功`
        })
    }

    mod.on_del_tree = async function () {
        let $dom = $(this).parent()
        let tid = $dom.attr("id")


        let del = await layer.choose("确认要删除行为树 " + mod.file[tid].title + '"么?')
        if (del === 0) {
            delete mod.file[tid]
            $dom.remove()
        }
    }


    mod.$editor = $(`<div id="behavior-tree-editor" style="display: none"></div>`, {})
    mod.$node_editor = $(`<div id="behavior-tree-editor-node" class="layui-form" lay-filter="behavior-tree" style="display: none">`)
    $('body').append(mod.$editor, mod.$node_editor)


    mod.init_editor = async function () {

        const rsp = await fetch(ddchess.path.mod_html('behavior_tree_node'))
        const html = await rsp.text()
        mod.$node_editor.append($(html))

        form.render(null, "behaviro-tree")

        mod.$node_name = mod.$node_editor.find("#name")
        mod.$node_category = mod.$node_editor.find("#category")
        mod.$node_id = mod.$node_editor.find("#id")
        mod.$node_title = mod.$node_editor.find("#title")
        mod.$node_props = mod.$node_editor.find("#properties")
        mod.json_editor = new JSONEditor(mod.$node_props[0], {mode: 'tree', mainMenuBar: false})

        form.on('select(category)', mod.select_node_cate)
        form.on('select(name)', mod.select_node_name)

        let nodes = await go.nodes()

        Object.values(nodes).forEach(function (node) {
            mod.$node_name.append($(`<option/>`, {
                text: node.title,
                value: node.name,
                category: node.category,
            }).data(node.properties ?? {}))
        })

        let reqs = await go.reqs

        let keys = Object.keys(reqs)
        keys.sort()
        keys.forEach(function (idx) {
            mod.$node_name.append($(`<option/>`, {text: idx, value: idx, category: 'custom'}).data({
                Server: "hall",
                ID: 0,
                Req: idx,
                Body: reqs[idx],
            }))
        })
    }


    mod.prepare_nodes = (tree) => [mod.prepare_node(tree.nodes, tree.nodes[tree.root])]

    mod.render_node_title = function (node, parent, idx) {
        let title = [
            mod.render_node_title_pad(node),
            mod.render_node_title_info(node),
            mod.render_node_title_parent(parent, idx),
        ]
        if ((node.category === "composite" || node.category === "decorator") && parent && parent.name === "IF") {
            title = [title[0], title[2], title[1]]
        }
        return title.join('')
    }

    mod.render_node_title_pad = node => (!node.children || node.children.length === 0) ? "&ensp;&ensp;" : ''

    mod.render_node_title_parent = function (parent, idx) {
        let prefix = ""
        if (parent && parent.name === 'IF') {
            prefix += ((idx === 0 || idx === '0') ? `<span class="layui-badge layui-bg-orange">Y</span>&ensp;` :
                `<span class="layui-badge layui-bg-orange">N</span>&ensp;`)
        } else if (parent && parent.name === 'Sequence') {
            prefix += (`<span class="layui-badge layui-bg-cyan">` + (idx - 0 + 1) + `</span>&ensp;`)
        }
        return prefix
    }
    mod.render_node_title_info = function (node) {
        let infos = {
            Sequence: `<span class="layui-badge layui-bg-cyan">⬇`,
            SubTree: `<span class="layui-badge layui-bg-gray">🌳`,
            WaitMsg: `<span class="layui-badge layui-bg-blue">→️`,
            IF: `<span class="layui-badge layui-bg-orange">❓`,
            RepeaterTimes: `<span class="layui-badge layui-bg-blue">🔁`,
        }

        let title = infos[node.name] ?? ''
        switch (node.name) {
            case 'IF':
                title += node.properties.Cond.substring(2, node.properties.Cond.length - 2);
                break;
            case 'SubTree':
                title += node.properties.Tree;
                break;
            case 'WaitMsg':
                title += node.properties.MsgID;
                break;
            case 'RepeaterTimes':
                title += node.properties.Times;
                break;

        }
        title += `</span>&ensp;`
        return title
    }

    mod.prepare_node = function (nodes, node, parent = undefined, child_idx = 0) {
        let copy_node = Object.assign({}, node)

        copy_node.head = mod.render_node_title(node, parent, child_idx)
        copy_node.edit = ['update']
        copy_node.spread = true

        if (parent !== undefined) {
            copy_node.pid = parent.id
            copy_node.edit.push("del")
        }

        child_idx !== 0 && copy_node.edit.push("move_up");

        node.category === 'decorator' && !(node?.children?.length) && copy_node.edit.push("add")
        node.category !== "action" && node.category !== 'decorator' && copy_node.edit.push("add")

        if (!(node?.children?.length)) {
            return copy_node
        }
        copy_node.children = []
        node.children.forEach((id, idx) => copy_node.children[idx] = mod.prepare_node(nodes, nodes[id], node, idx))

        return copy_node
    }

    mod.format_req_nodes = function () {

        mod.$node_name.find(`option[category=reqs]`).remove()
        Object.values(reqs.file).forEach(function (req, idx) {
            mod.$node_name.append($(`<option/>`, {text: req.id, value: req.id, category: 'reqs'}).data('properties', {
                Server: "hall",
                ID: 0,
                Req: req.req,
                Body: req.params,
            }))
        })
        form.render("select")
    }


    mod.select_node_name = function () {
        let $opt = mod.$node_name.find('option:selected')
        mod.$node_title.val($opt.text())
        mod.json_editor.set($opt.data() ?? {})
        form.render("select")
    }

    mod.select_node_cate = function () {
        let cate = mod.$node_category.val()
        let opts = mod.$node_name.find(`option[category=${cate}]`)
        opts.attr("disabled", false)
        mod.$node_name.find(`option[category!=${cate}]`).attr("disabled", true)

        if (cate !== mod.$node_name.find('option:selected').attr("category")) {
            mod.$node_title.val($(opts[0]).text())
            mod.$node_name.val($(opts[0]).val())
        }
        mod.select_node_name()
    }

    mod.render_node_editor = function (node) {
        mod.$node_id.val(node.id)
        node.category && mod.$node_category.val(node.category)
        node.title && mod.$node_title.val(node.title)
        node.name && mod.$node_name.val(node.name)
        mod.select_node_cate()
        Object.keys(node.properties ?? {}).length && mod.json_editor.set(node.properties)
    }

    mod.on_node_operated_add = function (parent) {
        let node = {id: mod.gen_uuid()}
        mod.$node_category.find("option").attr("disabled", false)
        mod.render_node_editor(node)
        layer.open_right('新建节点', mod.$node_editor, (idx) => {
            node.title = mod.$node_title.val()
            node.name = mod.$node_name.val()
            node.category = mod.$node_category.val()
            if (node.category === 'reqs' || node.category === 'custom') {
                node.category = 'action'
                node.name = 'SendReq'
            }

            if (node.category === "subtree") {
                node.category = 'action'
                node.name = 'SubTree'
            }

            node.properties = mod.json_editor.get()
            mod.current_tree.nodes[node.id] = node
            if (node.category === 'composite' || node.category === 'decorator') node.children = []
            mod.current_tree.nodes[parent.id].children.push(node.id)
            mod.render_tree()
            layer.close(idx)
        }, '65%')
    }


    mod.render_tree = function () {

        let data = mod.prepare_nodes(mod.current_tree)
        // 渲染
        tree.render({
            elem: '#behavior-tree-editor',
            data: data,
            edit: [],
            operate: ev => mod['on_node_operated_' + ev.type]?.(ev.data),
        })
    }

    mod.delete_node = function (node, del_parent = true) {
        delete mod.current_tree.nodes[node.id]

        if (del_parent && node.pid && mod.current_tree.nodes[node.pid]) {
            let parent = mod.current_tree.nodes[node.pid]
            let idx = parent.children.indexOf(node.id)
            if (idx === -1) return
            parent.children.splice(idx, 1)
        }

        if (node.children) {
            for (const child of node.children) {
                mod.delete_node(child, false)
            }
        }
        mod.render_tree()
    }

    mod.on_node_operated_del = async function (node) {
        let del = await layer.choose("确认要删除么？")
        if (del === 0) {
            mod.delete_node(node)
        }
    }


    mod.node_is_leaf = function (node) {
        return node.category === 'action' || node.category === 'subtree'
    }

    mod.on_node_operated_update = function (node) {
        let type = mod.$node_category.find(`option[value=${node.category}]`).attr('type')
        mod.$node_category.find(`option[type=${type}]`).attr("disabled", false)
        mod.$node_category.find(`option[type!=${type}]`).attr("disabled", true)

        mod.render_node_editor(node)
        layer.open_right("修改: " + node.title, mod.$node_editor, (idx) => {
            let oriNode = mod.current_tree.nodes[node.id]
            oriNode.title = mod.$node_title.val()
            oriNode.name = mod.$node_name.val()
            oriNode.category = mod.$node_category.val()
            oriNode.properties = mod.json_editor.get()
            mod.current_tree.nodes[node.id] = oriNode
            mod.render_tree()
            layer.close(idx)
        }, '65%')
    }

    mod.on_node_operated_move_up = function (node) {

        if (!node.pid || !mod.current_tree.nodes[node.pid]) {
            return
        }
        let parent = mod.current_tree.nodes[node.pid]
        if (!parent.children) {
            return
        }
        let idx = parent.children.findIndex((e) => node.id === e)
        if (idx === 0) {
            return
        }
        let tmp = parent.children[idx - 1]
        parent.children[idx - 1] = node.id
        parent.children[idx] = tmp
        mod.render_tree()
    }


    mod.editor = function (id) {
        if (!mod.file[id]) {
            notice.warning("未找到树:" + id)
            return
        }
        mod.current_tree = mod.file[id]

        mod.render_tree()

        mod.format_req_nodes()
        // 模拟数据
        layer.open_right(mod.file[id].title, mod.$editor, idx => {
            mod.file[id] = mod.file[id]
            layer.close(idx)
        }, '85%', {btn:false})
    }

    return mod
})

