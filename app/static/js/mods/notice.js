
layui.define2(function ($) {

    let $container;
    let listener;
    let toastId = 0;
    let toastType = {
        error: 'error',
        info: 'info',
        success: 'success',
        warning: 'warning'
    };

    layui.link("/css/notice.css", "notice")
    let cssStyle = $('<style type="text/css"></style>');
    $("head").append(cssStyle)
    let mod = {
        name: 'notice',
        options: {},
        version: '2.1.4',
    };

    let previousToast;


    ////////////////

    mod.error = function (message, title, optionsOverride) {
        return mod.notify({
            type: toastType.error,
            iconClass: getOptions().iconClasses.error,
            message: message,
            optionsOverride: optionsOverride,
            title: title
        });
    }

    function getContainer(options, create) {
        if (!options) {
            options = getOptions();
        }
        let $container = $('#' + options.containerId);
        if ($container.length) {
            return $container;
        }
        if (create) {
            $container = createContainer(options);
        }
        return $container;
    }

    mod.info = function (message, title, optionsOverride) {
        return mod.notify({
            type: toastType.info,
            iconClass: getOptions().iconClasses.info,
            message: message,
            optionsOverride: optionsOverride,
            title: title
        })
    }


    mod.success = function (message, title, optionsOverride) {
        return mod.notify({
            type: toastType.success,
            iconClass: getOptions().iconClasses.success,
            message: message,
            optionsOverride: optionsOverride,
            title: title
        })
    }

    mod.warning = function (message, title, optionsOverride) {
        return mod.notify({
            type: toastType.warning,
            iconClass: getOptions().iconClasses.warning,
            message: message,
            optionsOverride: optionsOverride,
            title: title
        })
    }

    mod.clear = function ($toastElement, clearOptions) {
        let options = getOptions();
        if (!$container) {
            getContainer(options);
        }
        if (!clearToast($toastElement, options, clearOptions)) {
            clearContainer(options);
        }
    }

    mod.remove = function ($toastElement) {
        let options = getOptions();
        if ($toastElement && $(':focus', $toastElement).length === 0) {
            removeToast($toastElement);
            return;
        }
        if (mod.$container.children().length) {
            mod.$container.remove();
        }
    }

    // internal functions

    function clearContainer(options) {
        let toastsToClear = mod.$container.children();
        for (let i = toastsToClear.length - 1; i >= 0; i--) {
            clearToast($(toastsToClear[i]), options);
        }
    }

    function clearToast($toastElement, options, clearOptions) {
        let force = clearOptions && clearOptions.force ? clearOptions.force : false;
        if ($toastElement && (force || $(':focus', $toastElement).length === 0)) {
            $toastElement[options.hideMethod]({
                duration: options.hideDuration,
                easing: options.hideEasing,
                complete: function () {
                    removeToast($toastElement);
                }
            });
            return true;
        }
        return false;
    }

    function createContainer(options) {
        let $container = $('<div/>')
            .attr('id', options.containerId)
            .addClass(options.positionClass);

        $container.appendTo($(options.target));
        return $container;
    }

    function getDefaults() {
        return {
            tapToDismiss: true,
            toastClass: 'toast',
            containerId: 'toast-container',
            debug: false,

            showMethod: 'fadeIn', //fadeIn, slideDown, and show are built into jQuery
            showDuration: 300,
            showEasing: 'swing', //swing and linear are built into jQuery
            onShown: undefined,
            hideMethod: 'fadeOut',
            hideDuration: 1000,
            hideEasing: 'swing',
            onHidden: undefined,
            closeMethod: false,
            closeDuration: false,
            closeEasing: false,
            closeOnHover: true,

            extendedTimeOut: 1000,
            iconClasses: {
                error: 'toast-error',
                info: 'toast-info',
                success: 'toast-success',
                warning: 'toast-warning'
            },
            iconClass: 'toast-info',
            positionClass: 'toast-top-mid',
            timeOut: 5000, // Set timeOut and extendedTimeOut to 0 to make it sticky
            titleClass: 'toast-title',
            messageClass: 'toast-message',
            escapeHtml: false,
            target: 'body',
            closeHtml: '<button type="button">&times;</button>',
            closeClass: 'toast-close-button',
            newestOnTop: true,
            preventDuplicates: false,
            progressBar: false,
            progressClass: 'toast-progress',
            rtl: false
        };
    }

    mod.notify = function (map) {
        let options = getOptions();
        let iconClass = map.iconClass || options.iconClass;

        if (typeof (map.optionsOverride) !== 'undefined') {
            options = $.extend(options, map.optionsOverride);
            iconClass = map.optionsOverride.iconClass || iconClass;
        }

        if (shouldExit(options, map)) {
            return;
        }

        toastId++;

        mod.$container = getContainer(options, true);

        let intervalId = null;
        let $toastElement = $('<div/>');
        let $titleElement = $('<div/>');
        let $messageElement = $('<div/>');
        let $progressElement = $('<div/>');
        let $closeElement = $(options.closeHtml);
        let progressBar = {
            intervalId: null,
            hideEta: null,
            maxHideTime: null
        };
        let response = {
            toastId: toastId,
            state: 'visible',
            startTime: new Date(),
            options: options,
            map: map
        };

        personalizeToast();

        displayToast();

        handleEvents();

        function escapeHtml(source) {
            if (source == null) {
                source = '';
            }

            return source
                .replace(/&/g, '&amp;')
                .replace(/"/g, '&quot;')
                .replace(/'/g, '&#39;')
                .replace(/</g, '&lt;')
                .replace(/>/g, '&gt;');
        }

        function personalizeToast() {
            setIcon();
            setTitle();
            setMessage();
            setCloseButton();
            setProgressBar();
            setRTL();
            setSequence();
            setAria();
        }

        function setAria() {
            var ariaValue = '';
            switch (map.iconClass) {
                case 'toast-success':
                case 'toast-info':
                    ariaValue = 'polite';
                    break;
                default:
                    ariaValue = 'assertive';
            }
            $toastElement.attr('aria-live', ariaValue);
        }

        function handleEvents() {
            if (options.closeOnHover) {
                $toastElement.hover(stickAround, delayedHideToast);
            }

            if (!options.onclick && options.tapToDismiss) {
                $toastElement.click(hideToast);
            }

            if (options.closeButton && $closeElement) {
                $closeElement.click(function (event) {
                    if (event.stopPropagation) {
                        event.stopPropagation();
                    } else if (event.cancelBubble !== undefined && event.cancelBubble !== true) {
                        event.cancelBubble = true;
                    }

                    if (options.onCloseClick) {
                        options.onCloseClick(event);
                    }

                    hideToast(true);
                });
            }

            if (options.onclick) {
                $toastElement.click(function (event) {
                    options.onclick(event);
                    hideToast();
                });
            }
        }

        function displayToast() {
            $toastElement.hide();

            $toastElement[options.showMethod](
                {duration: options.showDuration, easing: options.showEasing, complete: options.onShown}
            );

            if (options.timeOut > 0) {
                intervalId = setTimeout(hideToast, options.timeOut);
                progressBar.maxHideTime = parseFloat(options.timeOut);
                progressBar.hideEta = new Date().getTime() + progressBar.maxHideTime;
                if (options.progressBar) {
                    progressBar.intervalId = setInterval(updateProgress, 10);
                }
            }
        }

        function setIcon() {
            if (map.iconClass) {
                $toastElement.addClass(options.toastClass).addClass(iconClass);
            }
        }

        function setSequence() {
            if (options.newestOnTop) {
                mod.$container.prepend($toastElement);
            } else {
                mod.$container.append($toastElement);
            }
        }

        function setTitle() {
            if (map.title) {
                let suffix = map.title;
                if (options.escapeHtml) {
                    suffix = escapeHtml(map.title);
                }
                $titleElement.append(suffix).addClass(options.titleClass);
                $toastElement.append($titleElement);
            }
        }

        function setMessage() {
            if (map.message) {
                var suffix = map.message;
                if (options.escapeHtml) {
                    suffix = escapeHtml(map.message);
                }
                $messageElement.append(suffix).addClass(options.messageClass);
                $toastElement.append($messageElement);
            }
        }

        function setCloseButton() {
            if (options.closeButton) {
                $closeElement.addClass(options.closeClass).attr('role', 'button');
                $toastElement.prepend($closeElement);
            }
        }

        function setProgressBar() {
            if (options.progressBar) {
                $progressElement.addClass(options.progressClass);
                $toastElement.prepend($progressElement);
            }
        }

        function setRTL() {
            if (options.rtl) {
                $toastElement.addClass('rtl');
            }
        }

        function shouldExit(options, map) {
            if (options.preventDuplicates) {
                if (map.message === previousToast) {
                    return true;
                } else {
                    previousToast = map.message;
                }
            }
            return false;
        }

        function hideToast(override) {
            let method = override && options.closeMethod !== false ? options.closeMethod : options.hideMethod;
            let duration = override && options.closeDuration !== false ?
                options.closeDuration : options.hideDuration;
            let easing = override && options.closeEasing !== false ? options.closeEasing : options.hideEasing;
            if ($(':focus', $toastElement).length && !override) {
                return;
            }
            clearTimeout(progressBar.intervalId);
            return $toastElement[method]({
                duration: duration,
                easing: easing,
                complete: function () {
                    removeToast($toastElement);
                    clearTimeout(intervalId);
                    if (options.onHidden && response.state !== 'hidden') {
                        options.onHidden();
                    }
                    response.state = 'hidden';
                    response.endTime = new Date();
                }
            });
        }

        function delayedHideToast() {
            if (options.timeOut > 0 || options.extendedTimeOut > 0) {
                intervalId = setTimeout(hideToast, options.extendedTimeOut);
                progressBar.maxHideTime = parseFloat(options.extendedTimeOut);
                progressBar.hideEta = new Date().getTime() + progressBar.maxHideTime;
            }
        }

        function stickAround() {
            clearTimeout(intervalId);
            progressBar.hideEta = 0;
            $toastElement.stop(true, true)[options.showMethod](
                {duration: options.showDuration, easing: options.showEasing}
            );
        }

        function updateProgress() {
            let percentage = ((progressBar.hideEta - (new Date().getTime())) / progressBar.maxHideTime) * 100;
            $progressElement.width(percentage + '%');
        }

        return $toastElement;
    }

    function getOptions() {
        return $.extend({}, getDefaults(), mod.options);
    }

    function removeToast($toastElement) {
        if (!$container) {
            $container = getContainer();
        }
        if ($toastElement.is(':visible')) {
            return;
        }
        $toastElement.remove();
        $toastElement = null;
        if ($container.children().length === 0) {
            $container.remove();
            previousToast = undefined;
        }
    }
    return mod
})