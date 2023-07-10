layui.define2(function (colorpicker, server, $, log_config, behavior_tree, reqs) {

    let mod = {
        name: 'menu',
        css_vars: ["--color-0", "--color-1", "--color-2", "--color-4", "--color-5"],
    }


    mod.init = async function () {
        mod.$dom = $(`<div class="h100p menu" style="display: none">
    <div class="servers layui-form"><select title="servers" lay-search="" lay-filter="servers"></select></div>
    <div class="panel layui-btn-group"></div>
    <div class="behavior-tree"></div>
    <div class="reqs"></div>
</div>`)

        mod.$theme = $(`<div id="theme" style="display: none"></div>`)
        $("body").append(mod.$dom, mod.$theme)

        layer.open_once('菜单', mod.$dom, {area: ["240px", "100%"]})
        mod.$panel = mod.$dom.find(".panel")

        server.init(mod.$dom.find(".servers select"), mod.$panel)
        log_config.init(mod.$panel)
        behavior_tree.init(mod.$dom.find(".behavior-tree"), mod.$panel)
        reqs.init(mod.$dom.find(".reqs"), mod.$panel)
        mod.init_theme()

    }

    mod.init_theme = function () {
        mod.css_vars.to_object((v, i) => {
            let id = "vars-css-" + v
            let cv = localStorage.getItem(id)
            if (cv) {
                document.documentElement.style.setProperty(v, cv)
            } else {
                cv = getComputedStyle(document.documentElement).getPropertyValue(v)
            }
            mod.$theme.append($(`<div>`, {text: v}).append($(`<div id="${id}">`)))
            colorpicker.render({
                elem: '#' + id,
                color: cv,
                done: (color) => {
                    localStorage.setItem(id, color)
                    document.documentElement.style.setProperty(v, color)
                },
                change: (color) => {
                    document.documentElement.style.setProperty(v, color)
                },
            });
            return [v, cv]
        })
        mod.$panel.append_btn("🌈", "主题设置", function () {
            layer.open_once("主题设置", mod.$theme, {btn: false})
        })
    }

    return mod
})
