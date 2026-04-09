(function (window) {
    // 版本号
    const VERSION = "v1.1.0"

    // 先加载这三个LLM特征的js
    // DOM API
    const DomMonitor = (function () {
        const stats = {
            qs: 0, qsAll: 0, layoutReads: 0, timestamps: [], ops: []
        };
        const CONFIG = { windowMs: 60000 };

        function init() {
            // Hook querySelector
            const originalQs = Document.prototype.querySelector;
            Document.prototype.querySelector = function (...args) {
                record('qs');
                return originalQs.apply(this, args);
            };

            // Hook querySelectorAll
            const originalQsAll = Document.prototype.querySelectorAll;
            Document.prototype.querySelectorAll = function (...args) {
                record('qsAll');
                return originalQsAll.apply(this, args);
            };

            // Hook offsetHeight (强制重排检测)
            const descriptor = Object.getOwnPropertyDescriptor(HTMLElement.prototype, 'offsetHeight');
            if (descriptor && descriptor.get) {
                Object.defineProperty(HTMLElement.prototype, 'offsetHeight', {
                    get: function () {
                        record('layoutReads');
                        return descriptor.get.apply(this);
                    }
                });
            }
        }

        function record(type) {
            const now = performance.now();
            stats[type]++;
            stats.timestamps.push(now);
            stats.ops.push(type === 'qs' ? 1 : type === 'layoutReads' ? 2 : 3);
        }

        function getStats() {
            // 简单的爆发计算逻辑
            let burstCount = 0;
            for (let i = 1; i < stats.timestamps.length; i++) {
                if (stats.timestamps[i] - stats.timestamps[i - 1] < 4) burstCount++;
            }
            return {
                qsCount: stats.qs,
                qsAllCount: stats.qsAll,
                layoutReads: stats.layoutReads,
                burstCount: burstCount
            };
        }

        return { init, getStats };
    })();

    // Mutation Monitor 
    const MutationMonitor = (function () {
        const stats = { totalMutations: 0, uniqueNodes: 0, records: [] };
        const accessedNodes = new Set();
        let observer = null;

        function init() {
            observer = new MutationObserver((mutations) => {
                mutations.forEach(m => {
                    if (m.type !== 'attributes') return;
                    accessedNodes.add(m.target);
                    stats.totalMutations++;
                    if (stats.records.length < 100) { // 限制日志数量
                        stats.records.push({
                            tag: m.target.tagName,
                            attr: m.attributeName,
                            ts: performance.now()
                        });
                    }
                });
                stats.uniqueNodes = accessedNodes.size;
            });

            observer.observe(document.documentElement, {
                attributes: true, subtree: true, attributeFilter: ['style', 'class', 'id', 'readonly']
            });
        }

        function getStats() {
            return {
                totalMutations: stats.totalMutations,
                uniqueNodes: stats.uniqueNodes,
                mutationRecords: stats.records
            };
        }

        return { init, getStats };
    })();

    // Honeypot 
    const Honeypot = (function () {
        const report = { triggered: false, triggers: [] };

        function init() {
            if (document.readyState === 'loading') {
                document.addEventListener('DOMContentLoaded', renderTrap);
            } else {
                renderTrap();
            }
        }

        function renderTrap() {
            const fakeBtn = document.createElement('button');
            fakeBtn.id = 'loginbutton'; // 诱惑性 ID
            fakeBtn.innerText = 'Login';
            Object.assign(fakeBtn.style, {
                position: 'absolute', opacity: '0', pointerEvents: 'auto', zIndex: '-1', top: '0', left: '0'
            });

            fakeBtn.addEventListener('click', (e) => {
                e.preventDefault();
                e.stopPropagation();
                report.triggered = true;
                report.triggers.push({ type: 'fake_btn_click', ts: performance.now() });
            });

            document.body.appendChild(fakeBtn);
        }

        function getStats() {
            return {
                triggered: report.triggered,
                triggerCount: report.triggers.length,
                triggers: report.triggers
            };
        }

        return { init, getStats };
    })();


    try {
        DomMonitor.init();
        MutationMonitor.init();
        Honeypot.init();
        // console.log("Active Defense Modules Loaded.");
    } catch (e) {
        // console.error("Init failed:", e);
    }

    // hash
    async function sha256(message) {
        if (!message) return "";
        try {
            const msgBuffer = new TextEncoder().encode(message);
            const hashBuffer = await crypto.subtle.digest('SHA-256', msgBuffer);
            const hashArray = Array.from(new Uint8Array(hashBuffer));
            return hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
        } catch (e) {
            // console.error("Hash error:", e);
            return "hash_not_supported";
        }
    }

    // 获取 Hash Cookie 
    async function getCookieFingerprint() {
        const rawCookie = document.cookie || "";
        if (!rawCookie) return sha256("cookie_not_supported");

        const cookies = rawCookie.split(';')
            .map(c => c.trim())
            .filter(Boolean)
            .sort();   // 关键：排序保证稳定顺序

        const normalized = cookies.join(';');
        return await sha256(normalized);
    }

    // FontsFingerprint
    function getFontsFingerprint() {
        return new Promise(function (resolve) {
            var testString = 'mmMwWLliI0O&1';
            var textSize = '48px';
            var baseFonts = ['monospace', 'sans-serif', 'serif'];
            var fontList = [
                'sans-serif-thin', 'ARNO PRO', 'Agency FB', 'Arabic Typesetting',
                'Arial Unicode MS', 'AvantGarde Bk BT', 'BankGothic Md BT', 'Batang',
                'Bitstream Vera Sans Mono', 'Calibri', 'Century', 'Century Gothic',
                'Clarendon', 'EUROSTILE', 'Franklin Gothic', 'Futura Bk BT', 'Futura Md BT',
                'GOTHAM', 'Gill Sans', 'HELV', 'Haettenschweiler', 'Helvetica Neue',
                'Humanst521 BT', 'Leelawadee', 'Letter Gothic', 'Levenim MT', 'Lucida Bright',
                'Lucida Sans', 'Menlo', 'MS Mincho', 'MS Outlook', 'MS Reference Specialty',
                'MS UI Gothic', 'MT Extra', 'MYRIAD PRO', 'Marlett', 'Meiryo UI',
                'Microsoft Uighur', 'Minion Pro', 'Monotype Corsiva', 'PMingLiU',
                'Pristina', 'SCRIPTINA', 'Segoe UI Light', 'Serifa', 'SimHei',
                'Small Fonts', 'Staccato222 BT', 'TRAJAN PRO', 'Univers CE 55 Medium',
                'Vrinda', 'ZWAdobeF'
            ];
            var defaultWidth = {};
            var defaultHeight = {};
            var spansContainer = document.createElement('div');
            spansContainer.style.visibility = 'hidden';
            spansContainer.style.position = 'absolute';
            spansContainer.style.top = '0';
            spansContainer.style.left = '0';
            spansContainer.style.fontSize = textSize;
            var createSpan = function (fontFamily) {
                var span = document.createElement('span');
                span.style.fontFamily = fontFamily;
                span.textContent = testString;
                spansContainer.appendChild(span);
                return span;
            };
            var createSpanWithFallback = function (font, baseFont) {
                return createSpan("'".concat(font, "',").concat(baseFont));
            };
            // Create spans for base fonts and record their default dimensions
            var baseFontSpans = baseFonts.map(createSpan);
            baseFonts.forEach(function (font, i) {
                defaultWidth[font] = baseFontSpans[i].offsetWidth;
                defaultHeight[font] = baseFontSpans[i].offsetHeight;
            });
            // Create spans for all test fonts against each base font
            var fontSpans = {};
            var _loop_1 = function (font) {
                fontSpans[font] = baseFonts.map(function (base) { return createSpanWithFallback(font, base); });
            };
            for (var _i = 0, fontList_1 = fontList; _i < fontList_1.length; _i++) {
                var font = fontList_1[_i];
                _loop_1(font);
            }
            document.body.appendChild(spansContainer);
            var availableFonts = fontList.filter(function (font) {
                return fontSpans[font].some(function (span, i) {
                    return span.offsetWidth !== defaultWidth[baseFonts[i]] ||
                        span.offsetHeight !== defaultHeight[baseFonts[i]];
                });
            });
            document.body.removeChild(spansContainer);
            resolve(availableFonts);
        });
    }

    // CanvasFingerprint
    function getCanvasFingerprint() {
        var winding = false;
        var geometry;
        var text;
        var canvas = document.createElement('canvas');
        canvas.width = 1;
        canvas.height = 1;
        var context = canvas.getContext('2d');
        if (!context || !canvas.toDataURL) {
            geometry = text = 'unsupported';
        }
        else {
            winding = doesSupportWinding(context);
            var _a = renderImages(canvas, context), geo = _a[0], txt = _a[1];
            geometry = geo;
            text = txt;
        }
        return { winding: winding, geometry: geometry, text: text };
    }
    function doesSupportWinding(context) {
        context.rect(0, 0, 10, 10);
        context.rect(2, 2, 6, 6);
        return !context.isPointInPath(5, 5, 'evenodd');
    }
    function renderImages(canvas, context) {
        renderTextImage(canvas, context);
        var textImage1 = canvas.toDataURL();
        var textImage2 = canvas.toDataURL();
        if (textImage1 !== textImage2) {
            return ['unstable', 'unstable'];
        }
        renderGeometryImage(canvas, context);
        var geometryImage = canvas.toDataURL();
        return [geometryImage, textImage1];
    }
    function renderTextImage(canvas, context) {
        canvas.width = 240;
        canvas.height = 60;
        context.textBaseline = 'alphabetic';
        context.fillStyle = '#f60';
        context.fillRect(100, 1, 62, 20);
        context.fillStyle = '#069';
        context.font = '11pt "Times New Roman"';
        var printedText = "Cwm fjordbank gly ".concat(String.fromCharCode(55357, 56835) /* 😃 */);
        context.fillText(printedText, 2, 15);
        context.fillStyle = 'rgba(102, 204, 0, 0.2)';
        context.font = '18pt Arial';
        context.fillText(printedText, 4, 45);
    }
    function renderGeometryImage(canvas, context) {
        canvas.width = 122;
        canvas.height = 110;
        context.globalCompositeOperation = 'multiply';
        var circles = [
            ['#f2f', 40, 40],
            ['#2ff', 80, 40],
            ['#ff2', 60, 80],
        ];
        for (var _i = 0, circles_1 = circles; _i < circles_1.length; _i++) {
            var _a = circles_1[_i], color = _a[0], x = _a[1], y = _a[2];
            context.fillStyle = color;
            context.beginPath();
            context.arc(x, y, 40, 0, Math.PI * 2, true);
            context.closePath();
            context.fill();
        }
        context.fillStyle = '#f9c';
        context.arc(60, 60, 60, 0, Math.PI * 2, true);
        context.arc(60, 60, 20, 0, Math.PI * 2, true);
        context.fill('evenodd');
    }

    // WebGLFingerprint
    function safeGet(gl, pname) {
        try {
            if (typeof pname === "number" || (pname in gl)) {
                return gl.getParameter(pname);
            }
        } catch (e) { }
        return null;
    }
    function getWebGLContext() {
        const canvas = document.createElement("canvas");
        return (
            canvas.getContext("webgl") ||
            canvas.getContext("experimental-webgl") ||
            null
        );
    }
    async function getWebGLFingerprint() {
        try {
            const gl = getWebGLContext();
            if (!gl) return { supported: false };

            // --- 基础信息 ---
            const basicVendor = safeGet(gl, gl.VENDOR) || "";
            const basicRenderer = safeGet(gl, gl.RENDERER) || "";
            const version = safeGet(gl, gl.VERSION) || "";
            const shadingLang = safeGet(gl, gl.SHADING_LANGUAGE_VERSION) || "";

            // --- Chrome/Edge 扩展 ---
            let vendorUnmasked = "";
            let rendererUnmasked = "";
            const dbgRendererInfo = gl.getExtension("WEBGL_debug_renderer_info");
            if (dbgRendererInfo) {
                // ⚠️ Firefox 已经 deprecated，会 fallback 到 basic
                vendorUnmasked =
                    safeGet(gl, dbgRendererInfo.UNMASKED_VENDOR_WEBGL) || basicVendor;
                rendererUnmasked =
                    safeGet(gl, dbgRendererInfo.UNMASKED_RENDERER_WEBGL) || basicRenderer;
            }

            const basicInfo = {
                supported: true,
                version,
                shadingLanguageVersion: shadingLang,
                vendor: basicVendor,
                renderer: basicRenderer,
                vendorUnmasked,
                rendererUnmasked,
            };

            // --- 常用参数 ---
            const parameters = {
                aliasedLineWidthRange: safeGet(gl, gl.ALIASED_LINE_WIDTH_RANGE),
                aliasedPointSizeRange: safeGet(gl, gl.ALIASED_POINT_SIZE_RANGE),
                alphaBits: safeGet(gl, gl.ALPHA_BITS),
                depthBits: safeGet(gl, gl.DEPTH_BITS),
                stencilBits: safeGet(gl, gl.STENCIL_BITS),
                maxRenderbufferSize: safeGet(gl, gl.MAX_RENDERBUFFER_SIZE),
                maxTextureSize: safeGet(gl, gl.MAX_TEXTURE_SIZE),
            };

            // --- 扩展列表 ---
            const extensions = gl.getSupportedExtensions() || [];

            // --- 着色器精度 ---
            const shaderPrecisions = [];
            const shaderTypes = ["FRAGMENT_SHADER", "VERTEX_SHADER"];
            const precisionTypes = [
                "LOW_FLOAT",
                "MEDIUM_FLOAT",
                "HIGH_FLOAT",
                "LOW_INT",
                "MEDIUM_INT",
                "HIGH_INT",
            ];

            for (const shader of shaderTypes) {
                for (const precision of precisionTypes) {
                    const shaderTypeConst = gl[shader];
                    const precisionConst = gl[precision];
                    if (
                        typeof shaderTypeConst === "number" &&
                        typeof precisionConst === "number"
                    ) {
                        try {
                            const res = gl.getShaderPrecisionFormat(
                                shaderTypeConst,
                                precisionConst
                            );
                            if (res) {
                                shaderPrecisions.push(
                                    `${shader}.${precision}=${res.rangeMin},${res.rangeMax},${res.precision}`
                                );
                            }
                        } catch (e) {
                            // 某些设备/浏览器不支持，跳过
                        }
                    }
                }
            }

            return {
                ...basicInfo,
                parameters,
                extensions,
                shaderPrecisions,
            };
        } catch (e) {
            return { supported: false, error: e.message || "webgl_error" };
        }
    }

    // Level 1
    function getLevelOneSignals() {
        const signals = {};

        try {
            if ('webdriver' in navigator) {
                signals.webdriver = String(navigator.webdriver);
            } else {
                signals.webdriver = 'not_supported';
            }
            // signals.pluginsLength = navigator.plugins ? navigator.plugins.length : 0;
            (function () {
                try {
                    if (!('plugins' in navigator)) {
                        signals.pluginsList = null; // 浏览器不暴露
                        return;
                    }

                    const list = [];
                    for (let i = 0; i < navigator.plugins.length; i++) {
                        const p = navigator.plugins[i];
                        list.push({
                            name: p.name,
                            filename: p.filename,
                            description: p.description,
                            mimeTypes: Array.from(p).map(m => m.type)
                        });
                    }

                    signals.pluginsList = list;
                    signals.pluginsLength = navigator.plugins.length;
                } catch (e) {
                    signals.pluginsList = 'error';
                    signals.pluginsLength = '-1';
                }
            })();
            signals.languages = navigator.languages || [];
            signals.hasChrome = !!window.chrome;
            signals.hasChromeRuntime = !!(window.chrome && window.chrome.runtime);
            signals.mimeTypesLength = navigator.mimeTypes ? navigator.mimeTypes.length : 0;
            signals.hardwareConcurrency = navigator.hardwareConcurrency || 0;
            signals.outerVsScreenWidth = window.outerWidth - screen.width;
            signals.userAgent = navigator.userAgent || '';
        } catch (e) {
            // console.error("Bot signal collection error:", e);
        }

        return signals;
    }

    // Level 2
    async function getLevel2Signals() {
        const signals = {};

        // 基础属性
        signals.userAgent = navigator.userAgent || '';
        signals.platform = navigator.platform || '';
        signals.languages = navigator.languages || [];

        // UA 和 platform 一致性
        signals.uaPlatformMismatch = (signals.userAgent.includes("windows") && signals.platform.toLowerCase().includes("linux")) ||
            (signals.userAgent.includes("mac") && signals.platform.toLowerCase().includes("win")) ||
            (signals.userAgent.includes("linux") && signals.platform.toLocaleLowerCase().includes("mac"));

        // ✅ WebGL
        try {
            const webglData = await getWebGLFingerprint();
            signals.gpuVendor = webglData.vendorUnmasked || webglData.vendor || "unknown";
            signals.gpuRenderer = webglData.rendererUnmasked || webglData.renderer || "unknown";
        } catch (e) {
            signals.gpuVendor = "error";
            signals.gpuRenderer = "error";
        }

        // Permissions API
        try {
            const perm = await navigator.permissions.query({ name: 'notifications' });
            signals.notificationPermission = perm.state;
        } catch {
            signals.notificationPermission = 'error';
        }

        // window.chrome 内部结构检测
        signals.hasChromeApp = !!(window.chrome && window.chrome.app);
        signals.hasChromeRuntime = !!(window.chrome && window.chrome.runtime);

        // 屏幕与窗口大小
        signals.screenVsWindowMismatch = (screen.width !== window.innerWidth) || (screen.height !== window.innerHeight);

        // deviceMemory 与 hardwareConcurrency
        signals.deviceMemory = navigator.deviceMemory || 0;
        signals.hardwareConcurrency = navigator.hardwareConcurrency || 0;

        return signals;
    }

    // Level 3
    async function getAudioFingerprint() {
        try {
            // 使用 OfflineAudioContext 生成离线音频数据
            const OfflineCtx = window.OfflineAudioContext || window.AudioContext;
            if (!OfflineCtx) {
                return { error: "OfflineAudioContext not supported" };
            }

            // 参数: 1 个声道, 5000 采样点, 采样率 44100Hz
            const context = new OfflineCtx(1, 5000, 44100);

            // 创建振荡器 (OscillatorNode) 作为信号源
            const oscillator = context.createOscillator();
            oscillator.type = "sine";   // 正弦波
            oscillator.frequency.value = 1000; // 频率 1kHz

            // 创建滤波器 (BiquadFilterNode)
            const filter = context.createBiquadFilter();
            filter.type = "lowpass"; // 低通滤波器
            filter.frequency.value = 1500;

            // 链接音频图：oscillator -> filter -> destination
            oscillator.connect(filter);
            filter.connect(context.destination);

            oscillator.start(0);

            // 渲染离线音频
            const buffer = await context.startRendering();

            // 取出前 10 个采样点作为 fingerprint
            const channelData = buffer.getChannelData(0);
            const sample = Array.from(channelData.slice(0, 10));

            // 简单计算「抖动特征」: 相邻差值方差
            let diffs = [];
            for (let i = 1; i < sample.length; i++) {
                diffs.push(sample[i] - sample[i - 1]);
            }
            const mean = diffs.reduce((a, b) => a + b, 0) / diffs.length;
            const variance = diffs.reduce((a, b) => a + Math.pow(b - mean, 2), 0) / diffs.length;

            return {
                sample: sample.map(v => Number(v.toFixed(6))), // 固定精度
                jitterVariance: Number(variance.toFixed(8)),   // 抖动方差
            };
        } catch (e) {
            return { error: e.toString() };
        }
    }
    async function getRealtimeAudioFingerprint() {
        const AudioCtx = window.AudioContext || window.webkitAudioContext;
        if (!AudioCtx) {
            return { sample: null, jitterVar: null, error: "unsupported" };
        }

        // 定义一个 300ms 的硬超时，保证函数必返回
        const timeout = new Promise(resolve => {
            setTimeout(() => resolve({ sample: null, jitterVar: null, error: "timeout" }), 300);
        });

        const attempt = (async () => {
            try {
                const ctx = new AudioCtx();

                // 某些浏览器创建就失败
                if (!ctx || ctx.state === "closed") {
                    return { sample: null, jitterVar: null, error: "creation_failed" };
                }

                if (ctx.state === "suspended") return { sample: null, jitterVar: null, error: "blocked_by_autoplay_policy" }

                const oscillator = ctx.createOscillator();
                const analyser = ctx.createAnalyser();

                analyser.fftSize = 2048;
                oscillator.type = "sine";
                oscillator.frequency.value = 440;

                oscillator.connect(analyser);
                analyser.connect(ctx.destination);
                oscillator.start();

                await new Promise(r => setTimeout(r, 150));

                const array = new Float32Array(analyser.frequencyBinCount);
                analyser.getFloatFrequencyData(array);

                const sample = Array.from(array.slice(0, 5));
                const diffs = [];
                for (let i = 1; i < sample.length; i++) diffs.push(sample[i] - sample[i - 1]);
                const mean = diffs.reduce((a, b) => a + b, 0) / diffs.length;
                const variance = diffs.reduce((a, b) => a + Math.pow(b - mean, 2), 0) / diffs.length;

                oscillator.stop();
                ctx.close();

                return { sample, jitterVar: variance };
            } catch (e) {
                return { sample: null, jitterVar: null, error: e.message || "exception" };
            }
        })();

        return Promise.race([attempt, timeout]);
    }
    async function getLevel3Signals() {
        const signals = {};

        // 1. 环境完整性
        try {
            const availableFonts = await getFontsFingerprint();
            signals.fontsCount = availableFonts.length;
            //signals.fontsList = availableFonts; // 如果要看具体字体，可以加上
        } catch (e) {
            signals.fontsError = e.toString();
        }

        try {
            signals.hasMediaDevices = !!(navigator.mediaDevices && navigator.mediaDevices.enumerateDevices);
            signals.hasSpeechSynthesis = typeof speechSynthesis !== 'undefined';
            signals.intlTimeZone = Intl.DateTimeFormat().resolvedOptions().timeZone;
            signals.dateTimeZoneOffset = new Date().getTimezoneOffset();
        } catch (e) {
            signals.envError = e.toString();
        }

        // 2. 图形与音频
        try {
            const webglData = await getWebGLFingerprint();
            signals.gpuVendor = webglData.vendorUnmasked || webglData.vendor || "unknown";
            signals.gpuRenderer = webglData.rendererUnmasked || webglData.renderer || "unknown";
        } catch (e) {
            signals.gpuVendor = "error";
            signals.gpuRenderer = "error";
        }

        try {
            const audioData = await getAudioFingerprint();
            signals.audioSample = audioData.sample || [];
            signals.audioJitterVar = audioData.jitterVariance || 0;
        } catch (e) {
            signals.audioError = e.toString();
        }

        // try {
        //     const audioResult = await getRealtimeAudioFingerprint();
        //     signals.realtimeAudioSample = audioResult.sample;
        //     signals.realtimeAudioJitterVar = audioResult.jitterVar;
        //     signals.realtimeAudioError = audioResult.error;
        // } catch (e) {
        //     signals.realtimeAudioSample = null;
        //     signals.realtimeAudioJitterVar = null;
        //     signals.realtimeAudioError = e.toString();
        // }

        // 3. 执行上下文
        signals.requestIdleCallbackSupported = typeof requestIdleCallback === 'function';
        signals.queueMicrotaskSupported = typeof queueMicrotask === 'function';
        signals.touchEventSupported = typeof TouchEvent !== 'undefined';
        signals.pointerEventSupported = typeof PointerEvent !== 'undefined';

        // 简单测试 requestIdleCallback 是否能执行
        signals.idleCallbackExecuted = false;
        if (signals.requestIdleCallbackSupported) {
            await new Promise(resolve => {
                requestIdleCallback(() => {
                    signals.idleCallbackExecuted = true;
                    resolve();
                }, { timeout: 100 });
            });
        }

        // 4. Chrome DevTools 残留
        signals.hasReactDevTools = !!window.__REACT_DEVTOOLS_GLOBAL_HOOK__;
        signals.hasDevtools = !!window.devtools;

        return signals;
    }

    // Mousemove
    const CONFIG = {
        WINDOW_MS: 15000,       // 扩大到 15000ms，降低短流程场景下 eventCount=0 概率
        MIN_PUSH_INTERVAL: 8,   // 针对高频事件(mousemove/scroll)的节流间隔
        MAX_QUEUE_SIZE: 1000    // 防止内存溢出的安全上限
    };

    // 统一的事件队列，存储所有类型的交互
    let eventQueue = [];
    let lastPushTs = 0;

    function recordEvent(e) {
        const now = Date.now();

        // 1. 提取空间坐标 (Spatial Data)
        // 并不是所有事件都有坐标，键盘事件通常为 null，但在分析时这本身就是特征
        let x = null;
        let y = null;

        if (e.type.startsWith('mouse') || e.type === 'click' || e.type === 'contextmenu') {
            x = e.clientX || 0;
            y = e.clientY || 0;
        } else if (e.type.startsWith('touch') && e.touches.length > 0) {
            x = e.touches[0].clientX || 0;
            y = e.touches[0].clientY || 0;
        }

        // 2. 提取目标元素信息 (Contextual Data)
        // 论文提到记录 DOM Object IDs
        let targetTag = 'document';
        let targetId = null;
        if (e.target) {
            targetTag = e.target.tagName;
            targetId = e.target.id || null;
        }

        // 3. 构造数据包
        const eventData = {
            t: now,             // Temporal: 时间戳
            e: e.type,          // Type: 事件类型
            x: x,               // Spatial: X坐标
            y: y,               // Spatial: Y坐标
            tag: targetTag,     // Context: 标签名
            id: targetId,       // Context: 元素ID (如果有)
            // 对于键盘事件，可以选记录 key (慎用，涉及隐私，论文中主要用作频率分析)
            k: e.type.startsWith('key') ? e.code : undefined
        };

        // 4. 节流处理 (针对 mousemove 和 scroll)
        // 如果是高频事件，且距离上次记录时间太短，则跳过（除非是不同的事件类型）
        const isHighFreq = (e.type === 'mousemove' || e.type === 'scroll' || e.type === 'wheel');
        if (isHighFreq) {
            if (now - lastPushTs < CONFIG.MIN_PUSH_INTERVAL) {
                // 可选：更新队列中最后一个同类型事件的坐标和时间，保持最新状态
                const last = eventQueue[eventQueue.length - 1];
                if (last && last.e === e.type) {
                    last.x = x;
                    last.y = y;
                    last.t = now;
                    cleanup(now);
                    return;
                }
            }
        }

        // 入队
        eventQueue.push(eventData);
        lastPushTs = now;

        // 清理过期数据
        cleanup(now);
    }

    // 清理函数：移除 3000ms 之前的数据
    function cleanup(now = Date.now()) {
        const cutoff = now - CONFIG.WINDOW_MS;

        // 移除过期数据
        while (eventQueue.length && eventQueue[0].t < cutoff) {
            eventQueue.shift();
        }

        // 安全截断：防止极端情况下队列过大
        if (eventQueue.length > CONFIG.MAX_QUEUE_SIZE) {
            eventQueue.splice(0, eventQueue.length - CONFIG.MAX_QUEUE_SIZE);
        }
    }

    function initWebGuard() {
        const options = { passive: true, capture: true };

        // 鼠标类 (用于捕捉轨迹和点击意图)
        document.addEventListener('mousemove', recordEvent, options);
        document.addEventListener('mousedown', recordEvent, options);
        document.addEventListener('mouseup', recordEvent, options);
        document.addEventListener('click', recordEvent, options);
        document.addEventListener('wheel', recordEvent, options);

        // 键盘类 (LLM Agent 填写表单时会有极快或极规律的输入)
        document.addEventListener('keydown', recordEvent, options);
        document.addEventListener('keyup', recordEvent, options);

        // 页面交互类 (Agent 可能会触发 focus/blur 但没有鼠标移动)
        window.addEventListener('scroll', recordEvent, options);
        window.addEventListener('focus', recordEvent, options);
        window.addEventListener('blur', recordEvent, options);

        // 移动端支持 (防止 Agent 伪装成移动设备)
        document.addEventListener('touchstart', recordEvent, options);
        document.addEventListener('touchend', recordEvent, options);
    }

    // 导出数据接口
    function getInteractionTrace() {
        const now = Date.now();
        cleanup(now);

        // 返回深拷贝，防止外部修改
        return {
            trace: JSON.parse(JSON.stringify(eventQueue)),
            meta: {
                timestamp: now,
                windowMs: CONFIG.WINDOW_MS,
                eventCount: eventQueue.length
            }
        };
    }
    initWebGuard();

    // KeyBoard
    const keystrokes = [];
    let lastKeyTime = null;

    // 最大保留按键数量（避免日志无限增长，可根据需求调大或设为 null）
    const MAX_KEYS = 500;

    // 键盘监听
    document.addEventListener('keydown', (e) => {
        const now = Date.now();
        const interval = lastKeyTime ? now - lastKeyTime : 0;
        lastKeyTime = now;

        keystrokes.push({
            key: e.key,
            code: e.code,
            timestamp: now,
            interval: interval
        });


        if (MAX_KEYS && keystrokes.length > MAX_KEYS) {
            keystrokes.shift();
        }
    }, { passive: true }); // 不阻止默认事件


    function getKeyboardData() {
        return {
            strokes: keystrokes.slice(),
            total: keystrokes.length,
            timestamp: Date.now()
        };
    }

    // 辅助工具：给 Promise 加超时，超时就返回默认值
    const withTimeout = (promise, ms, defaultValue) => {
        const timeout = new Promise(resolve => setTimeout(() => resolve(defaultValue), ms));
        return Promise.race([promise, timeout]);
    };

    async function collectAllData(username) {
        const TIMEOUT = 1000;

        const canvasData = getCanvasFingerprint() || "missing";
        const level1Signals = getLevelOneSignals() || "missing";
        const interactionData = getInteractionTrace() || "missing";
        const keyboardData = getKeyboardData() || "missing";

        const [cookieData, fontsData, webglData, level2Signals, level3Signals] = await Promise.all([
            withTimeout(getCookieFingerprint(), TIMEOUT, "calculating"),
            withTimeout(getFontsFingerprint(), TIMEOUT, "calculating"),
            withTimeout(getWebGLFingerprint(), TIMEOUT, "calculating"),
            withTimeout(getLevel2Signals(), TIMEOUT, "calculating"),
            withTimeout(getLevel3Signals(), TIMEOUT, "calculating")
        ]);

        const llmData = {
            domAnomaly: DomMonitor.getStats(),
            mutation: MutationMonitor.getStats(),
            honeypot: Honeypot.getStats()
        };
        return {
            username: username,
            version: VERSION,
            cookie: cookieData,
            fonts: fontsData,
            canvas: canvasData,
            webgl: webglData,
            leve1: level1Signals,
            leve2: level2Signals,
            level3: level3Signals,
            mousemove: interactionData,
            keyboard: keyboardData,
            llmNature: llmData
        };
    }

    window.BotDetector = {
        getFingerprint: collectAllData
    };

})(window);