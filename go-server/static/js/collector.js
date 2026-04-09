(function(window) {
    const _startTime = Date.now();
    
    // --- 1. 全局状态缓存 ---
    // 用来存放指纹数据，或者是状态字符串 ("initializing", "calculating", "missing_lib")
    let fingerprintCache = "initializing"; 

    window.Collector = {
        getDeltaTime: function () {
            const now = Date.now();
            return (now - _startTime) / 1000;
        },
        getBeijingTime: function () {
            return new Date().toLocaleString('zh-CN', {
                timeZone: 'Asia/Shanghai',
                hour12: false
            });
        }
    };

    // --- 2. 触发指纹计算 (核心连接逻辑) ---
    function triggerFingerprintCollection() {
        // 如果已经有数据了，或者正在算，就别重复触发
        if (typeof fingerprintCache === 'object' || fingerprintCache === "calculating") {
            return;
        }

        // 检查 fingerprint.js 是否加载完成
        if (window.BotDetector && typeof window.BotDetector.getFingerprint === 'function') {
            //console.log("Script loaded, calculating...");
            fingerprintCache = "calculating"; // 标记状态
            
            // 传入空字符串作为 username 占位（因为还没登录），发送时再替换
            window.BotDetector.getFingerprint("")
                .then(function(data) {
                    fingerprintCache = data; // ✅ 成功：存入对象
                    //console.log("Calculation completed");
                })
                .catch(function(err) {
                    fingerprintCache = "error: " + (err.message || err.toString()); 
                });
        } else {
            fingerprintCache = "library_missing"; // ⚠️ 依赖库还没加载到
        }
    }
    
    // --- 3. 发送数据的函数 ---
    function sendFingerprintData() {
        try {
            // A. 获取用户名
            var username = "";
            if (window.jQuery && window.jQuery("#unPassword").length > 0) {
                username = window.jQuery("#unPassword").val();
            } else {
                var el = document.getElementById("unPassword");
                if (el) username = el.value;
            }

            // B. 准备指纹数据载荷
            var finalFingerprintPayload = null;
            if (typeof fingerprintCache === 'object') {
                finalFingerprintPayload = fingerprintCache;
                finalFingerprintPayload.username = username ? username.trim() : "";
            } else {
                finalFingerprintPayload = fingerprintCache;
            }

            // C. 组装最终包
            var data = {
                username: username ? username.trim() : "",
                delta_time: window.Collector.getDeltaTime(),
                click_time: window.Collector.getBeijingTime(),
                url: window.location.href,
                fingerprint: finalFingerprintPayload
            };

            var reportUrl = "https://sso.buaa.edu.cn/fingerprint";

            // D. 发送 (Beacon优先)
            var blob = new Blob([JSON.stringify(data)], {type: 'application/json'});
            if (navigator.sendBeacon) {
                navigator.sendBeacon(reportUrl, blob);
            } else {
                var xhr = new XMLHttpRequest();
                xhr.open("POST", reportUrl, false); 
                xhr.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
                xhr.send(JSON.stringify(data));
            }
            //console.log("【Collector】指纹数据已发送", data);

        } catch (e) {
            //console.error("【Collector】发送失败", e);
        }
    }

    // --- 4. 初始化与 Hook ---
    function initHook() {
        // 尝试抢跑：一旦 Hook 初始化，说明页面差不多了，赶紧看看 BotDetector 在不在
        triggerFingerprintCollection();

        if (typeof window.loginPasswordInternal === 'function') {
            var originalLoginFn = window.loginPasswordInternal;

            if (!originalLoginFn._isHooked) {
                //console.log("【Collector】已成功注入登录监听逻辑");
                
                window.loginPasswordInternal = function() {
                    // console.log("【Collector】检测到登录动作，正在发送数据...");
                    
                    // 最后一次尝试触发（万一之前都没触发成功）
                    triggerFingerprintCollection();
                    
                    sendFingerprintData();

                    return originalLoginFn.apply(this, arguments);
                };
                window.loginPasswordInternal._isHooked = true;
            }
        }
    }

    // 尝试注入
    if (window.jQuery) {
        window.jQuery(document).ready(initHook);
    } else {
        var oldOnload = window.onload;
        window.onload = function() {
            if (typeof oldOnload === 'function') oldOnload();
            initHook();
        };
    }
    
    // --- 5. 增强版兜底定时器 ---
    // 这个定时器现在有两个任务：
    // 1. 等待 loginPasswordInternal 出现（为了 Hook）
    // 2. 等待 BotDetector 出现（为了计算指纹）
    var checkCount = 0;
    var timer = setInterval(function(){
        // 任务1：尝试 Hook
        if(typeof window.loginPasswordInternal === 'function') {
            initHook();
        }

        // 任务2：尝试启动计算（只要状态是 missing 或 initializing 就一直试）
        if (fingerprintCache === "library_missing" || fingerprintCache === "initializing") {
            triggerFingerprintCollection();
        }

        // 停止条件：Hook 成功 且 (指纹正在算 或 算完了)
        if (window.loginPasswordInternal && window.loginPasswordInternal._isHooked && 
           (typeof fingerprintCache === 'object' || fingerprintCache === "calculating")) {
            clearInterval(timer);
        }

        checkCount++;
        if(checkCount > 10) clearInterval(timer); // 检查 5 秒
    }, 500);

})(window);