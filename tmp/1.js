const WebSocket = require("ws");

// 配置参数
const WS_SERVER_URL =
	"ws://mainnet.helius-rpc.com/?api-key=75d6bbc9-fa1b-46e2-9cea-1248f1235d6f"; // 替换为实际WS地址
const JSON_RPC_REQUEST = {
	jsonrpc: "2.0",
	id: 1,
	method: "logsSubscribe",
	params: [
		{
			mentions: ["DfMxre4cKmvogbLrPigxmibVTTQDuzjdXojWzjCXXhzj"],
		},
		{
			commitment: "confirmed",
		},
	],
};

// 创建WebSocket连接
const socket = new WebSocket(WS_SERVER_URL);

// 添加监听器
socket.on("open", () => {
	console.log("✅ WebSocket连接已建立");

	// 发送JSON-RPC请求
	const requestJson = JSON.stringify(JSON_RPC_REQUEST);
	socket.send(requestJson);
	console.log("📤 已发送请求:", requestJson);
});

socket.on("message", (data) => {
	try {
		const response = JSON.parse(data);

		// 检查错误响应
		if (response.error) {
			console.error("❌ 收到错误响应:");
			console.error(`  错误码: ${response.error.code}`);
			console.error(`  错误信息: ${response.error.message}`);
			return;
		}

		// 处理不同的响应类型
		if (response.id === JSON_RPC_REQUEST.id) {
			// 直接请求响应
			console.log("📨 收到订阅响应:");
			console.log(`  订阅ID: ${response.result}`);
			console.log("⚡ 准备接收日志通知...");
		} else if (response.method === "logsNotification") {
			// 订阅通知处理
			console.log("🔔 收到日志通知:");
			console.log("  订阅ID:", response.params.subscription);
			console.log("  日志内容:", response.params.result.value);
		} else {
			console.log("ℹ️ 收到其他消息:", JSON.stringify(response, null, 2));
		}
	} catch (e) {
		console.error("⚠️ 解析响应失败:", e.message);
		console.log("原始数据:", data);
	}
});

socket.on("close", (code, reason) => {
	console.log(`🚪 连接已关闭 - 状态码: ${code}, 原因: ${reason}`);
});

socket.on("error", (error) => {
	console.error("‼️ WebSocket错误:", error.message);
});

// 添加心跳保持连接
setInterval(() => {
	if (socket.readyState === WebSocket.OPEN) {
		socket.ping();
	}
}, 30000); // 每30秒发送一次心跳
