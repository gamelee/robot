function ObjectProxy(obj, saver, every_time = false) {
    this.saver = saver
    this.obj = obj
    this.dirty = false
    every_time || (this.timer = setInterval(() => {
        if (!this.dirty) {
            return
        }
        this.dirty = false
        this.saver(this.obj)
    }, 1000))

    this.save = () => {
        if (every_time) {
            this.saver(this.obj)
        } else {
            this.dirty = true
        }
    }

    this.proxy = (obj) => {
        return new Proxy(obj, {
            set: (target, p, newVal, receiver) => {
                Reflect.set(target, p, newVal, receiver)
                this.save()
                return true
            },
            deleteProperty: (target, p) => {
                let v= Reflect.deleteProperty(target, p)
                this.save()
                return v
            },
            get: (target, p) => {
                if (target[p] instanceof Object && p !== 'prototype') {
                    return this.proxy(target[p])
                }
                return Reflect.get(target, p)
            }
        })
    }
    return this.proxy(this.obj)
}

let _parse = JSON.parse
JSON.parse = (str) => {
    let rs = {};
    try {
        rs = _parse(str)
    } catch (e) {
        return rs
    }
    return rs
}

function Each(obj, callback) {
    let length, i = 0;

    if (Array.isArray(obj)) {
        let out = [];
        length = obj.length;
        for (; i < length; i++) {
            out.push(callback(obj[i]))
        }
        return out
    } else {
        let out = {}
        for (i in obj) {
            out[i] = callback(obj[i], i)
        }
        return out
    }

}


function sum(list, field = false) {
    let ret = 0;

    (list ?? []).map((v) => {
        ret += field === false ? v : (v[field] ?? 0)
    })
    return ret
}

function merge(list, by, field) {
    let ret = {};
    (list ?? []).map((v) => {
        let id = v[by]
        ret[id] = (ret[id] ?? 0) + (v[field] ?? 0)
    })
    return ret
}

Array.prototype.to_object = function (fn, group = false) {
    let out = {}
    if (typeof fn === 'string') {
        let field = fn
        fn = (v, id) => {
            return [v[field], v]
        }
    }
    for (let i = 0; i < (this??[]).length; i++) {

        let [id, data] = fn(this[i], i)
        if (group) {
            id in out || (out[id] = [])
            out[id].push(data)
        } else {
            out[id] = data
        }
    }
    return out
}

to_slice = function (obj, fn) {
    return Object.entries(obj).map(([k, v]) => {
        return fn(k, v)
    })
}