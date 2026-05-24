import { useState } from 'react';
import { Card, Tabs, Button, Descriptions, Tag, Spin, Space, message } from 'antd';
import { RobotOutlined, SmileOutlined, FileTextOutlined, FormOutlined } from '@ant-design/icons';
import api from '../../api/client';

interface AiAnalysis {
  sentiment: { label: string; score: number; emoji: string };
  summary: string;
  tags: string[];
  ticketDraft: Record<string, string>;
}

const sentimentMap: Record<string, { color: string; emoji: string }> = {
  positive: { color: 'green', emoji: '😊' },
  neutral:  { color: 'default', emoji: '😐' },
  negative: { color: 'red', emoji: '😠' },
};

export default function AiAssistPanel({ callId }: { callId?: number }) {
  const [analysis, setAnalysis] = useState<AiAnalysis | null>(null);
  const [loading, setLoading] = useState(false);

  const runAnalysis = async () => {
    if (!callId) return;
    setLoading(true);
    try {
      const res = await api.post('/ai/analysis/realtime', { call_id: callId });
      setAnalysis(res.data);
    } catch {
      message.error('分析失败');
    } finally {
      setLoading(false);
    }
  };

  const createTicket = async () => {
    if (!analysis?.ticketDraft) return;
    try {
      await api.post('/tickets', analysis.ticketDraft);
      message.success('工单已创建');
    } catch {
      message.error('创建失败');
    }
  };

  return (
    <Card
      title={<Space><RobotOutlined /> AI 辅助</Space>}
      size="small"
      extra={<Button size="small" type="primary" onClick={runAnalysis} loading={loading}>分析</Button>}
    >
      {loading ? (
        <div style={{ textAlign: 'center', padding: 40 }}><Spin tip="AI 分析中..." /></div>
      ) : !analysis ? (
        <div style={{ textAlign: 'center', padding: 40, color: '#999' }}>点击"分析"获取 AI 辅助信息</div>
      ) : (
        <Tabs size="small" items={[
          {
            key: 'sentiment',
            label: <Space><SmileOutlined />情绪</Space>,
            children: (
              <div style={{ textAlign: 'center', padding: 16 }}>
                <div style={{ fontSize: 48 }}>{sentimentMap[analysis.sentiment.label]?.emoji || '😐'}</div>
                <Tag color={sentimentMap[analysis.sentiment.label]?.color}>{analysis.sentiment.label}</Tag>
                <div style={{ marginTop: 8, color: '#666' }}>置信度: {Math.round(analysis.sentiment.score * 100)}%</div>
                {analysis.tags.length > 0 && (
                  <div style={{ marginTop: 12 }}>
                    {analysis.tags.map((t) => <Tag key={t} color="blue">{t}</Tag>)}
                  </div>
                )}
              </div>
            ),
          },
          {
            key: 'summary',
            label: <Space><FileTextOutlined />摘要</Space>,
            children: (
              <div style={{ whiteSpace: 'pre-wrap', lineHeight: 1.8, padding: 8 }}>
                {analysis.summary || '暂无摘要'}
              </div>
            ),
          },
          {
            key: 'ticket',
            label: <Space><FormOutlined />填单</Space>,
            children: (
              <>
                <Descriptions size="small" column={1} bordered>
                  {Object.entries(analysis.ticketDraft || {}).map(([k, v]) => (
                    <Descriptions.Item key={k} label={k}>{v}</Descriptions.Item>
                  ))}
                </Descriptions>
                <Button type="primary" block style={{ marginTop: 12 }} onClick={createTicket}>
                  创建工单
                </Button>
              </>
            ),
          },
        ]} />
      )}
    </Card>
  );
}
