import { useState, useEffect, useRef } from 'react';
import { Card, Tag, Empty, Switch, Space } from 'antd';
import { AudioOutlined, UserOutlined, RobotOutlined } from '@ant-design/icons';

interface TranscriptLine {
  id: string;
  role: 'agent' | 'customer' | 'system';
  text: string;
  timestamp: string;
  sentiment?: 'positive' | 'neutral' | 'negative';
}

const roleConfig = {
  agent:    { label: '坐席', color: 'blue',    icon: <UserOutlined /> },
  customer: { label: '客户', color: 'green',   icon: <AudioOutlined /> },
  system:   { label: '系统', color: 'default', icon: <RobotOutlined /> },
};

const sentimentColor = { positive: '#52c41a', neutral: '#999', negative: '#ff4d4f' };

export default function RealtimeTranscriptPanel() {
  const [lines, setLines] = useState<TranscriptLine[]>([]);
  const [autoScroll, setAutoScroll] = useState(true);
  const bottomRef = useRef<HTMLDivElement>(null);

  // Simulate receiving transcript lines via polling or WebSocket.
  useEffect(() => {
    const demo: TranscriptLine[] = [
      { id: '1', role: 'system', text: '通话已接通', timestamp: new Date().toLocaleTimeString() },
    ];
    setLines(demo);
  }, []);

  useEffect(() => {
    if (autoScroll) bottomRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [lines, autoScroll]);

  return (
    <Card
      title={<Space><AudioOutlined /> 实时转写</Space>}
      size="small"
      extra={<Switch checkedChildren="自动滚动" unCheckedChildren="暂停" checked={autoScroll} onChange={setAutoScroll} />}
      style={{ height: '100%' }}
      bodyStyle={{ maxHeight: 400, overflowY: 'auto', padding: '8px 12px' }}
    >
      {lines.length === 0 ? (
        <Empty description="等待通话开始..." image={Empty.PRESENTED_IMAGE_SIMPLE} />
      ) : (
        lines.map((line) => {
          const cfg = roleConfig[line.role];
          return (
            <div key={line.id} style={{ marginBottom: 8, display: 'flex', gap: 8, alignItems: 'flex-start' }}>
              <Tag icon={cfg.icon} color={cfg.color} style={{ flexShrink: 0 }}>{cfg.label}</Tag>
              <div style={{ flex: 1 }}>
                <span style={{ color: line.sentiment ? sentimentColor[line.sentiment] : undefined }}>
                  {line.text}
                </span>
                <span style={{ fontSize: 11, color: '#999', marginLeft: 8 }}>{line.timestamp}</span>
              </div>
            </div>
          );
        })
      )}
      <div ref={bottomRef} />
    </Card>
  );
}
