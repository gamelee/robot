class DDChessDB {

    constructor(projectName = 'ddchess-db') {
        this.projectName = projectName;
        this.ready(projectName)
    }

    ready(projectName = this.projectName) {
        return new Promise((resolve, reject) => {
            if (this.db) {
                resolve(this);
            } else {
                this.objectStoreName = 'ddchess-table';
                const request = window.indexedDB.open(projectName, 1);
                request.onsuccess = (event) => {
                    this.db = event.target.result;
                    resolve(this);
                };
                request.onupgradeneeded = (event) => {
                    const db = event.target.result;
                    if (!db.objectStoreNames.contains(this.objectStoreName)) {
                        db.createObjectStore(this.objectStoreName);
                    }
                };
                request.onerror = (event) => {
                    reject(event);
                }
            }
        })
    }

    init(request, callback) {
        return new Promise((resolve, reject) => {
            const success = value => {
                if (callback && typeof callback === 'function') {
                    callback(false, value);
                }
                resolve(value);
            }
            const error = event => {
                if (callback && typeof callback === 'function') {
                    callback(event);
                }
                reject(event);
            }
            return this.ready().then(() => {
                request(success, error);
            }).catch(error)
        })
    }

    eventList = []

    // 触发事件
    triggerEvent() {
        document.dispatchEvent(new CustomEvent('ddchess-storage', {
            detail: [...arguments]
        }))
        this.eventList.forEach(fn => {
            fn(...arguments)
        })
    }

    // 自定义事件
    addEventListener(callback) {
        this.eventList.push(callback);
    }

    async setItem(key, value, callback) {
        return this.init((success, error) => {
            const request = this.db.transaction(this.objectStoreName, 'readwrite').objectStore(this.objectStoreName).put(value, key);
            request.onsuccess = () => {
                success(value);
                this.triggerEvent('setItem', {key, value}, this.projectName);
            };
            request.onerror = error;
            // this.dispatchEvent(new CustomEvent('update'))
        }, callback)
    }

    async getItem(key, callback) {
        return this.init((success, error) => {
            const request = this.db.transaction(this.objectStoreName).objectStore(this.objectStoreName).get(key);
            request.onsuccess = () => success(request.result);
            request.onerror = error;
        }, callback)
    }

    async removeItem(key, callback) {
        return this.init((success, error) => {
            const request = this.db.transaction(this.objectStoreName, 'readwrite').objectStore(this.objectStoreName).delete(key);
            request.onsuccess = () => {
                success(key);
                this.triggerEvent('removeItem', {key}, this.projectName);
            };
            request.onerror = error;
        }, callback)
    }

    key(index, callback) {
        return this.init((success, error) => {
            const request = this.db.transaction(this.objectStoreName).objectStore(this.objectStoreName).getAllKeys();
            request.onsuccess = () => success(request.result[index]);
            request.onerror = error;
        }, callback)
    }

    keys(callback) {
        return this.init((success, error) => {
            const request = this.db.transaction(this.objectStoreName).objectStore(this.objectStoreName).getAllKeys();
            request.onsuccess = () => success(request.result);
            request.onerror = error;
        }, callback)
    }

    clear(callback) {
        return this.init((success, error) => {
            const request = this.db.transaction(this.objectStoreName, 'readwrite').objectStore(this.objectStoreName).clear();
            request.onsuccess = () => {
                success(null);
                this.triggerEvent('clear', {}, this.projectName);
            };
            request.onerror = error;
        }, callback)
    }
}

window.ddchess.db = new DDChessDB()
window.ddchess.dao = async (key, dft = {}) => {
    let obj = (await window.ddchess.db.getItem(key)) || dft
    return new ObjectProxy(obj, dat => window.ddchess.db.setItem(key, dat), false)
}