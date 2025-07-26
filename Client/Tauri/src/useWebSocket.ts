import { ref, computed, onUnmounted } from 'vue';

export interface UseWebSocketOptions {
  url: string | (() => string);
  password?: string | (() => string);
  onMessage?: (data: any) => void;
  onOpen?: () => void;
  onClose?: () => void;
  onError?: (err: Event) => void;
  onConnecting?: () => void;
  onConnectFailed?: (error: string) => void;
}

export type ConnectionStatus = 'disconnected' | 'connecting' | 'connected' | 'failed';

export function useWebSocket(options: UseWebSocketOptions) {
  const ws = ref<WebSocket | null>(null);
  const connectionStatus = ref<ConnectionStatus>('disconnected');
  const lastError = ref<string | null>(null);
  const isConnecting = ref(false);
  let onMessage: ((data: any) => void) | undefined = options.onMessage;

  function connect() {
    console.log('WebSocket connect called, current status:', connectionStatus.value, 'isConnecting:', isConnecting.value);
    
    if (isConnecting.value) {
      console.log('WebSocket already connecting, ignoring connect call');
      return; // 防止重复连接
    }
    
    if (ws.value) {
      console.log('Closing existing WebSocket connection');
      // 临时设置状态为connecting，避免触发重连逻辑
      const oldStatus = connectionStatus.value;
      connectionStatus.value = 'connecting';
      ws.value.close();
      ws.value = null;
      // 恢复原状态，让新的连接流程继续
      connectionStatus.value = oldStatus;
    }
    
    let wsUrl = typeof options.url === 'function' ? options.url() : options.url;
    if (!wsUrl) {
      lastError.value = 'WebSocket URL 未设置';
      connectionStatus.value = 'failed';
      options.onConnectFailed && options.onConnectFailed('WebSocket URL 未设置');
      return;
    }
    
    isConnecting.value = true;
    connectionStatus.value = 'connecting';
    lastError.value = null;
    options.onConnecting && options.onConnecting();
    
    if (options.password) {
      wsUrl += (wsUrl.includes('?') ? '&' : '?') + 'password=' + encodeURIComponent(typeof options.password === 'function' ? options.password() : options.password);
    }
    
    let connectionTimeout: number | undefined;
    
    try {
      ws.value = new WebSocket(wsUrl);
      
      // 设置连接超时
      connectionTimeout = setTimeout(() => {
        if (connectionStatus.value === 'connecting') {
          lastError.value = '连接超时';
          connectionStatus.value = 'failed';
          isConnecting.value = false;
          options.onConnectFailed && options.onConnectFailed('连接超时');
          if (ws.value) {
            ws.value.close();
          }
        }
      }, 10000); // 10秒超时

      ws.value.onopen = () => {
        if (connectionTimeout) {
          clearTimeout(connectionTimeout);
        }
        connectionStatus.value = 'connected';
        isConnecting.value = false;
        lastError.value = null;
        console.log('WebSocket connected successfully'); // 调试日志
        options.onOpen && options.onOpen();
        
        // 发送连接成功消息
        if (window && (window as any).__TERMINAL_BUS__) {
          (window as any).__TERMINAL_BUS__.emit('terminal-message', 'WebSocket连接成功');
        }
        
        // 移除自动状态查询，让Home.vue统一管理状态查询
        // 这样可以避免重复查询和潜在的冲突
      };
      
      ws.value.onmessage = (event) => {
        console.log('WebSocket raw message:', event.data); // 调试日志
        let data = event.data;
        try {
          data = JSON.parse(event.data);
          console.log('WebSocket parsed message:', data); // 调试日志
        } catch (error) {
          console.log('WebSocket message is not JSON, using as string:', event.data); // 调试日志
          data = event.data;
        }
        
        // 调用消息处理器
        if (onMessage) {
          console.log('Calling onMessage handler with:', data); // 调试日志
          onMessage(data);
        }
        if (options.onMessage) {
          console.log('Calling options.onMessage handler with:', data); // 调试日志
          options.onMessage(data);
        }
      };
      
      ws.value.onclose = (event) => {
        console.log('WebSocket onclose event, code:', event.code, 'reason:', event.reason);
        clearTimeout(connectionTimeout);
        isConnecting.value = false;
        
        if (connectionStatus.value === 'connecting') {
          // 如果还在连接状态就关闭了，说明连接失败
          connectionStatus.value = 'failed';
          lastError.value = event.code === 1006 ? '连接被拒绝或网络错误' : `连接关闭 (代码: ${event.code})`;
          console.log('WebSocket connection failed:', lastError.value);
          options.onConnectFailed && options.onConnectFailed(lastError.value);
        } else {
          connectionStatus.value = 'disconnected';
          console.log('WebSocket connection closed normally');
        }
        
        options.onClose && options.onClose();
      };
      
      ws.value.onerror = (err) => {
        clearTimeout(connectionTimeout);
        isConnecting.value = false;
        connectionStatus.value = 'failed';
        lastError.value = 'WebSocket 连接发生错误';
        options.onError && options.onError(err);
        options.onConnectFailed && options.onConnectFailed('WebSocket 连接发生错误');
      };
      
    } catch (error) {
      if (connectionTimeout) {
        clearTimeout(connectionTimeout);
      }
      isConnecting.value = false;
      connectionStatus.value = 'failed';
      lastError.value = '创建 WebSocket 连接失败';
      options.onConnectFailed && options.onConnectFailed('创建 WebSocket 连接失败');
    }
  }

  function disconnect() {
    console.log('WebSocket disconnect called, current status:', connectionStatus.value);
    
    if (ws.value) {
      console.log('Closing WebSocket connection');
      ws.value.close();
      ws.value = null;
    }
    
    connectionStatus.value = 'disconnected';
    isConnecting.value = false;
    lastError.value = null;
    
    console.log('WebSocket disconnected, status reset to:', connectionStatus.value);
  }

  function send(data: any) {
    console.log('WebSocket send called with:', data); // 调试日志
    console.log('Current connection status:', connectionStatus.value); // 调试日志
    
    if (ws.value && connectionStatus.value === 'connected') {
      const msg = {
        command: data.command || '',
        content: data.content || '',
        output: data.output || '',
        error: data.error || '',
        status: data.status || ''
      };
      console.log('Sending message:', msg); // 调试日志
      ws.value.send(JSON.stringify(msg));
    } else {
      console.error('Cannot send message: WebSocket not connected or not ready');
      console.error('ws.value:', !!ws.value, 'connectionStatus:', connectionStatus.value);
    }
  }

  // 兼容性：保持原有的 isConnected 属性
  const isConnected = computed(() => connectionStatus.value === 'connected');

  Object.defineProperty(connect, 'onMessage', {
    get() { return onMessage; },
    set(fn) {
      onMessage = fn;
      if (ws.value) {
        ws.value.onmessage = (event) => {
          console.log('WebSocket raw message (from setter):', event.data); // 调试日志
          let data = event.data;
          try {
            data = JSON.parse(event.data);
            console.log('WebSocket parsed message (from setter):', data); // 调试日志
          } catch (error) {
            console.log('WebSocket message is not JSON (from setter), using as string:', event.data); // 调试日志
            data = event.data;
          }
          
          if (onMessage) {
            console.log('Calling onMessage handler (from setter) with:', data); // 调试日志
            onMessage(data);
          }
          if (options.onMessage) {
            console.log('Calling options.onMessage handler (from setter) with:', data); // 调试日志
            options.onMessage(data);
          }
        };
      }
    }
  });

  onUnmounted(() => {
    disconnect();
  });

  return {
    ws,
    isConnected,
    connectionStatus,
    isConnecting,
    lastError,
    connect,
    disconnect,
    send,
    get onMessage() { return onMessage; },
    set onMessage(fn) {
      onMessage = fn;
      if (ws.value) {
        ws.value.onmessage = (event) => {
          let data = event.data;
          try {
            data = JSON.parse(event.data);
          } catch {}
          if (onMessage) onMessage(data);
          options.onMessage && options.onMessage(data);
        };
      }
    },
  };
} 