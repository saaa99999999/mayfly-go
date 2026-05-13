import { createWebSocket } from '@/common/request';
import { ref, onBeforeUnmount, type Ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { ElMessage } from 'element-plus';

/**
 * AI Chat WebSocket 连接管理 Hook
 * 负责 WebSocket 连接、重连、消息收发等
 */
export function useAiChatWebSocket(
  onMessage: (data: any) => void,
  currentSessionId: Ref<string>,
  isNewSession: Ref<boolean>
) {
  const { t } = useI18n();
  
  const socket = ref<WebSocket | null>(null);
  const reconnectTimer = ref<any>(null);
  const reconnectAttempts = ref(0);
  const MAX_RECONNECT_ATTEMPTS = 5;
  const RECONNECT_DELAY = 3000;

  /**
   * 初始化 WebSocket 连接
   */
  const initSocket = async () => {
    try {
      console.log('init chat ws...');
      const ws = await createWebSocket(`/ai/chat`);
      socket.value = ws;
      
      ws.onmessage = (e) => {
        const data = JSON.parse(e.data);
        
        // 会话隔离：只处理属于当前激活会话的消息
        if (data.sessionId && data.sessionId !== currentSessionId.value) {
          // 新会话首次收到后端返回的真实 sessionId，更新并通知父组件
          if (isNewSession.value) {
            currentSessionId.value = data.sessionId;
          } else {
            console.log(`忽略不属于当前会话的消息: ${data.sessionId} !== ${currentSessionId.value}`);
            return;
          }
        }
        
        onMessage(data);
      };

      ws.onclose = (event) => {
        console.log('chat ws 连接关闭:', event.code, event.reason);
        if (!event.wasClean) {
          attemptReconnect();
        }
      };

      ws.onerror = (error) => {
        console.error('chat ws 错误:', error);
      };

      // 连接成功，重置重连计数
      reconnectAttempts.value = 0;
    } catch (e) {
      console.log('连接错误', e);
      // 直接显示错误提示，不传递到消息处理器
      ElMessage.error(t('ai.chat.connectionFailed'));
      attemptReconnect();
    }
  };

  /**
   * 尝试重连
   */
  const attemptReconnect = () => {
    if (reconnectAttempts.value >= MAX_RECONNECT_ATTEMPTS) {
      console.warn('达到最大重连次数，停止重连');
      ElMessage.error(t('ai.chat.connectionDisconnected'));
      return;
    }

    reconnectAttempts.value++;
    console.log(`尝试第 ${reconnectAttempts.value} 次重连...`);

    if (reconnectTimer.value) {
      clearTimeout(reconnectTimer.value);
    }

    reconnectTimer.value = setTimeout(() => {
      initSocket();
    }, RECONNECT_DELAY);
  };

  /**
   * 清理 WebSocket 连接
   */
  const cleanupSocket = () => {
    if (reconnectTimer.value) {
      clearTimeout(reconnectTimer.value);
      reconnectTimer.value = null;
    }
    if (socket.value) {
      socket.value.onclose = null;
      socket.value.onerror = null;
      socket.value.onmessage = null;
      if (socket.value.readyState === WebSocket.OPEN || socket.value.readyState === WebSocket.CONNECTING) {
        socket.value.close();
      }
      socket.value = null;
    }
    reconnectAttempts.value = 0;
  };

  /**
   * 发送消息
   */
  const sendMessage = (type: 'text' | 'interruptResume', content: string) => {
    // 检查 WebSocket 连接状态
    if (!socket.value || socket.value.readyState === WebSocket.CLOSED || socket.value.readyState === WebSocket.CLOSING) {
      console.warn('WebSocket 连接已关闭，尝试重连...');
      
      // 如果正在重连中，等待重连完成
      if (reconnectAttempts.value > 0 && reconnectAttempts.value < MAX_RECONNECT_ATTEMPTS) {
        ElMessage.warning(t('ai.chat.reconnecting'));
        attemptReconnect();
        return;
      }

      // 立即尝试重连
      attemptReconnect();
      ElMessage.error(t('ai.chat.connectionLost'));
      return;
    }

    socket.value.send(
      JSON.stringify({
        type,
        sessionId: currentSessionId.value,
        content,
      })
    );
  };

  /**
   * 获取当前连接状态
   */
  const isConnected = () => {
    return socket.value && socket.value.readyState === WebSocket.OPEN;
  };

  // 组件卸载时清理连接
  onBeforeUnmount(() => {
    cleanupSocket();
  });

  return {
    initSocket,
    sendMessage,
    reconnectAttempts,
    MAX_RECONNECT_ATTEMPTS,
  };
}
