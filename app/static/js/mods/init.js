layui.define2(function (layer, $, notice) {

    let mod = {name: 'init'}

    mod.layer_config = {
        skin: 'class-layer-ddchess',
        type: 1,
        shadeClose: true,
    }

    layer.open_left = function (title, $content, yes, width = '60%', opt = {}) {
        let tmp = $.extend({}, mod.layer_config, {
            offset: 'l', anim: 'slideRight',
            btn: ["保存", "取消"],
            area: [width, '100%'],
            title: title, content: $content,
            yes: yes, zIndex: layer.zIndex,
            success(layero, index) {},
        }, opt)
        tmp.area[0] = width
        return layer.open(tmp)
    }


    layer.open_right = function (title, $content, yes, width = '60%', opt = {}) {
        let tmp = $.extend({}, mod.layer_config, {
            offset: 'r', anim: 'slideLeft',
            btn: ["保存", "取消"],
            area: [width, '100%'],
            title: title, content: $content,
            yes: yes, zIndex: layer.zIndex,
            success(layero, index) {},
        }, opt)
        tmp.area[0] = width
        return layer.open(tmp)
    }

    layer.choose = function (msg, title = "信息", btn = ["确认", "取消"]) {
        return new Promise((resolve, reject) => {
            let fns = btn.to_object((v, id) => ['btn' + (id + 1), (idx) => {
                cancel = false
                layer.close(idx)
                resolve(id)
            }])
            let opt = {btn, ...fns}
            opt.skin = mod.layer_config.skin
            opt.title = title
            let cancel = true
            opt.end = () => {
                cancel && resolve(-1)
            }
            layer.confirm(msg, opt)
        })
    }

    layer.input = function (title) {
        return new Promise((resolve, reject) => {
            layer.prompt({
                title: title, end: () => {
                    resolve(false)
                },
                skin: mod.layer_config.skin,
            }, function (name, idx, elem) {
                layer.close(idx)
                resolve(name)
            })
        })
    }

    window.ddchess.once = window.ddchess.once || {}

    layer.open_once = function (title, content, opt) {

        let open = function () {
            let tmp = $.extend({}, mod.layer_config, {
                id: title, type: 1, shade: 0, maxmin: true,
                resize: true, fixed: true, title: title,
                closeBtn: false, content: content, btn: false,
                maxHeight: '50%', zIndex: layer.zIndex,
                offset: 'rb', success(layero, index) {
                    layer.setTop(layero)
                },
                end: () => {
                    delete window.ddchess.once[title]
                },
            }, opt)
            window.ddchess.once[title] = layer.open(tmp)
        }
        if (window.ddchess.once[title]) {
            layer.close(window.ddchess.once[title], () => {
                open()
            })
        } else {
            open()
        }

    }

    layer.get_index = function (id) {
        return window.ddchess.once[id] ?? false
    }

    $.fn.extend({
        append_btn: function (content, tips = '', click = undefined, cls = "layui-btn-primary layui-btn-xs") {
            let html = `<button type="button" class="layui-btn ${cls}" style="border-color: white"`
            if (tips) {
                html += ` onmouseenter="this.layer_idx = layer.tips('${tips}', this,{tips:3})"`
                html += ` onmouseout="layer.close(this.layer_idx)"`
            }
            html += `>${content}</button>`
            let $dom = $(html)
            this.append($dom)
            click && $dom.on('click', click)
            return this
        },
        htmlJSON: function (obj) {
            this.html(JSON.stringify(obj))
        }
    })

    mod.$json_editor = $(`<div>`, {class: 'json-editor', style: "display: none", id: "init-json-editor"})
    $('body').append(mod.$json_editor, mod.$js_editor)
    mod.json_editor = new JSONEditor(mod.$json_editor[0], {mode: 'tree', mainMenuBar: true})

    $.extend({
        json_editor: async function (obj, title, left = true) {
            return new Promise((resolve, reject) => {
                let cancel = true
                mod.json_editor.set(obj);
                mod.json_editor.setMode('tree');
                (left ? layer.open_left : layer.open_right)?.(title, mod.$json_editor, (idx) => {
                    resolve(mod.json_editor.get())
                    cancel = false
                    notice.success("保存成功")
                    layer.close(idx)
                }, '65%', {
                    end: (idx) => {
                        cancel && resolve(false)
                    }
                })
            })
        },
        json_viewer: function (obj, title = "", left = true) {
            mod.json_editor.set(obj);
            mod.json_editor.setMode('view');
            (left ? layer.open_left : layer.open_right)?.(title, mod.$json_editor, async (idx) => {
                await navigator?.clipboard?.writeText?.(JSON.stringify(obj))
                notice.success("复制成功")
                layer.close(idx)
            }, '65%', {btn: ["复制", "关闭"]});
        },
    })

    function stopF5Refresh() {
        document.onkeydown = function (oEvent) {
            oEvent = oEvent || window.oEvent;
            //获取键盘的keyCode值
            let nKeyCode = oEvent.keyCode || oEvent.which || oEvent.charCode;
            //获取ctrl 键对应的事件属性
            let bCtrlKeyCode = oEvent.ctrlKey || oEvent.metaKey;
            if (oEvent?.nKeyCode === 83 && bCtrlKeyCode) {
                alert('save');
                //doSomeThing...
            }
        }

        document.onkeydown = function (e) {
            let evt = window.event || e
            let code = evt.keyCode || evt.which;
            //屏蔽F1---F12
            if (code > 111 && code < 123 || (e.ctrlKey && e.key.toUpperCase() === 'R')) {
                if (evt.preventDefault) {
                    evt.preventDefault();
                } else {
                    evt.keyCode = 0;
                    evt.returnValue = false;
                }
            }
        }
        //禁止鼠标右键菜单
        document.oncontextmenu = function (e) {
            return false;
        }
        //阻止后退的所有动作，包括 键盘、鼠标手势等产生的后退动作。
        history.pushState(null, null, window.location.href);
        window.addEventListener("popstate", function () {
            history.pushState(null, null, window.location.href);
        })
    }

    stopF5Refresh()
    return mod
})

