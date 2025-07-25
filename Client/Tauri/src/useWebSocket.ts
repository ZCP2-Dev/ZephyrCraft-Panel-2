import { ref, onUnmounted } from 'vue';

export interface UseWebSocketOptions {
  url: string;
  password?: string;
  onMessage?: (data: any) => void;
  onOpen?: () => void;
  onClose?: () => void;
  onError?: (err: Event) => void;
}

export function useWebSocket(options: UseWebSocketOptions) {
  const ws = ref<WebSocket | null>(null);
  const isConnected = ref(false);
  const lastError = ref<string | null>(null);
  let onMessage: ((data: any) => void) | undefined = options.onMessage;

  function connect() {
    if (ws.value) {
      ws.value.close();
    }
    let wsUrl = typeof options.url === 'function' ? options.url() : options.url;
    if (options.password) {
      wsUrl += (wsUrl.includes('?') ? '&' : '?') + 'password=' + encodeURIComponent(typeof options.password === 'function' ? options.password() : options.password);
    }
    ws.value = new WebSocket(wsUrl);

    ws.value.onopen = () => {
      isConnected.value = true;
      options.onOpen && options.onOpen();
    };
    ws.value.onmessage = (event) => {
      let data = event.data;
      try {
        data = JSON.parse(event.data);
      } catch {}
      if (onMessage) onMessage(data);
      options.onMessage && options.onMessage(data);
    };
    ws.value.onclose = () => {
      isConnected.value = false;
      options.onClose && options.onClose();
    };
    ws.value.onerror = (err) => {
      lastError.value = 'WebSocket 连接发生错误';
      options.onError && options.onError(err);
    };
  }

  function disconnect() {
    if (ws.value) {
      ws.value.close();
      ws.value = null;
    }
  }

  function send(data: any) {
    if (ws.value && isConnected.value) {
      const msg = {
        command: data.command || '',
        content: data.content || '',
        output: data.output || '',
        error: data.error || '',
        status: data.status || ''
      };
      ws.value.send(JSON.stringify(msg));
    }
  }

  Object.defineProperty(connect, 'onMessage', {
    get() { return onMessage; },
    set(fn) { onMessage = fn; }
  });

  onUnmounted(() => {
    disconnect();
  });

  return {
    ws,
    isConnected,
    lastError,
    connect,
    disconnect,
    send,
    get onMessage() { return onMessage; },
    set onMessage(fn) { onMessage = fn; },
  };
} 