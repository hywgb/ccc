import { useState, useEffect, useRef, useCallback } from 'react';
import { Input, Button, Avatar, Badge, Spin } from 'antd';
import {
  MessageOutlined, CloseOutlined, SendOutlined, PaperClipOutlined,
  CustomerServiceOutlined, MinusOutlined,
} from '@ant-design/icons';

interface ChatMessage {
  id: string;
  role: 'customer' | 'agent' | 'system';
  content: string;
  timestamp: string;
  type: 'text' | 'image' | 'file';
  fileUrl?: string;
}

interface WebchatWidgetProps {
  tenantId: string;
  channelId: string;
  apiBase?: string;
  title?: string;
  greeting?: string;
  primaryColor?: string;
}

export default function WebchatWidget({
  tenantId, channelId, apiBase = '/api/v1', title = '在线客服',
  greeting = '您好！有什么可以帮您？', primaryColor = '#1677ff',
}: WebchatWidgetProps) {
  const [open, setOpen] = useState(false);
  const [minimized, setMinimized] = useState(false);
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [input, setInput] = useState('');
  const [sessionId, setSessionId] = useState<string>('');
  const [connected, setConnected] = useState(false);
  const [sending, setSending] = useState(false);
  const [unread, setUnread] = useState(0);
  const bottomRef = useRef<HTMLDivElement>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const scrollToBottom = () => bottomRef.current?.scrollIntoView({ behavior: 'smooth' });
  useEffect(scrollToBottom, [messages]);

  const initSession = useCallback(async () => {
    try {
      const res = await fetch(`${apiBase}/webchat/sessions`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ tenant_id: tenantId, channel_id: channelId }),
      });
      const data = await res.json();
      setSessionId(data.id || data.session_id);
      setConnected(true);
      setMessages([{
        id: 'greeting', role: 'system', content: greeting,
        timestamp: new Date().toLocaleTimeString(), type: 'text',
      }]);
    } catch {
      setMessages([{
        id: 'error', role: 'system', content: '连接失败，请稍后再试',
        timestamp: new Date().toLocaleTimeString(), type: 'text',
      }]);
    }
  }, [tenantId, channelId, apiBase, greeting]);

  useEffect(() => { if (open && !sessionId) initSession(); }, [open, sessionId, initSession]);

  // Poll for new messages
  useEffect(() => {
    if (!sessionId || !connected) return;
    const poll = setInterval(async () => {
      try {
        const res = await fetch(`${apiBase}/webchat/sessions/${sessionId}/messages`);
        const data = await res.json();
        const msgs: ChatMessage[] = Array.isArray(data) ? data : data?.items || [];
        if (msgs.length > messages.length) {
          setMessages(msgs);
          if (minimized || !open) setUnread((u) => u + msgs.length - messages.length);
        }
      } catch { /* ignore */ }
    }, 3000);
    return () => clearInterval(poll);
  }, [sessionId, connected, messages.length, minimized, open, apiBase]);

  const sendMessage = async () => {
    const text = input.trim();
    if (!text || !sessionId) return;
    setSending(true);
    setInput('');
    const tempMsg: ChatMessage = {
      id: `temp_${Date.now()}`, role: 'customer', content: text,
      timestamp: new Date().toLocaleTimeString(), type: 'text',
    };
    setMessages((m) => [...m, tempMsg]);
    try {
      await fetch(`${apiBase}/webchat/sessions/${sessionId}/messages`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ content: text, type: 'text' }),
      });
    } catch { /* ignore */ }
    setSending(false);
  };

  const handleFile = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file || !sessionId) return;
    const formData = new FormData();
    formData.append('file', file);
    try {
      await fetch(`${apiBase}/webchat/sessions/${sessionId}/messages`, { method: 'POST', body: formData });
    } catch { /* ignore */ }
  };

  // Floating button
  if (!open) {
    return (
      <div
        style={{
          position: 'fixed', bottom: 24, right: 24, zIndex: 9999,
          width: 56, height: 56, borderRadius: '50%', background: primaryColor,
          display: 'flex', alignItems: 'center', justifyContent: 'center',
          cursor: 'pointer', boxShadow: '0 4px 12px rgba(0,0,0,0.2)',
        }}
        onClick={() => { setOpen(true); setUnread(0); }}
      >
        <Badge count={unread} offset={[-4, 4]}>
          <MessageOutlined style={{ fontSize: 24, color: '#fff' }} />
        </Badge>
      </div>
    );
  }

  // Minimized state
  if (minimized) {
    return (
      <div
        style={{
          position: 'fixed', bottom: 24, right: 24, zIndex: 9999,
          background: primaryColor, color: '#fff', padding: '8px 16px', borderRadius: 24,
          cursor: 'pointer', boxShadow: '0 4px 12px rgba(0,0,0,0.2)', display: 'flex', alignItems: 'center', gap: 8,
        }}
        onClick={() => { setMinimized(false); setUnread(0); }}
      >
        <Badge count={unread} size="small"><CustomerServiceOutlined style={{ fontSize: 18, color: '#fff' }} /></Badge>
        <span>{title}</span>
      </div>
    );
  }

  return (
    <div style={{
      position: 'fixed', bottom: 24, right: 24, zIndex: 9999,
      width: 380, height: 520, borderRadius: 12, overflow: 'hidden',
      boxShadow: '0 8px 32px rgba(0,0,0,0.2)', display: 'flex', flexDirection: 'column', background: '#fff',
    }}>
      {/* Header */}
      <div style={{ background: primaryColor, color: '#fff', padding: '12px 16px', display: 'flex', alignItems: 'center', gap: 8 }}>
        <CustomerServiceOutlined style={{ fontSize: 20 }} />
        <span style={{ flex: 1, fontWeight: 500 }}>{title}</span>
        <MinusOutlined style={{ cursor: 'pointer', fontSize: 16 }} onClick={() => setMinimized(true)} />
        <CloseOutlined style={{ cursor: 'pointer', fontSize: 16, marginLeft: 8 }} onClick={() => setOpen(false)} />
      </div>

      {/* Messages */}
      <div style={{ flex: 1, overflowY: 'auto', padding: '12px 16px' }}>
        {messages.map((msg) => (
          <div key={msg.id} style={{ marginBottom: 12, display: 'flex', flexDirection: msg.role === 'customer' ? 'row-reverse' : 'row', gap: 8 }}>
            <Avatar size={32} style={{ background: msg.role === 'agent' ? primaryColor : msg.role === 'customer' ? '#87d068' : '#999', flexShrink: 0 }}>
              {msg.role === 'agent' ? '客' : msg.role === 'customer' ? '我' : '系'}
            </Avatar>
            <div style={{
              maxWidth: '70%', padding: '8px 12px', borderRadius: 8, fontSize: 14, lineHeight: 1.5,
              background: msg.role === 'customer' ? primaryColor : '#f0f0f0',
              color: msg.role === 'customer' ? '#fff' : '#333',
            }}>
              {msg.type === 'image' ? <img src={msg.fileUrl} alt="" style={{ maxWidth: '100%', borderRadius: 4 }} /> : msg.content}
              <div style={{ fontSize: 11, opacity: 0.7, marginTop: 4, textAlign: 'right' }}>{msg.timestamp}</div>
            </div>
          </div>
        ))}
        <div ref={bottomRef} />
      </div>

      {/* Input */}
      <div style={{ padding: '8px 12px', borderTop: '1px solid #f0f0f0', display: 'flex', gap: 8, alignItems: 'center' }}>
        <input ref={fileInputRef} type="file" hidden onChange={handleFile} />
        <Button size="small" type="text" icon={<PaperClipOutlined />} onClick={() => fileInputRef.current?.click()} />
        <Input
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onPressEnter={sendMessage}
          placeholder="输入消息..."
          style={{ flex: 1 }}
          suffix={sending ? <Spin size="small" /> : undefined}
        />
        <Button type="primary" icon={<SendOutlined />} onClick={sendMessage} disabled={!input.trim()} style={{ background: primaryColor }} />
      </div>
    </div>
  );
}
