import { useState, useEffect } from 'react';
import { Card, List, Button, Tag, Empty, Space, message } from 'antd';
import { BulbOutlined, CopyOutlined, LikeOutlined } from '@ant-design/icons';
import api from '../../api/client';

interface ScriptItem {
  id: string;
  title: string;
  content: string;
  category: string;
  confidence: number;
}

export default function ScriptRecommendPanel({ callId }: { callId?: number }) {
  const [scripts, setScripts] = useState<ScriptItem[]>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (!callId) return;
    setLoading(true);
    api.get(`/ai/script-recommend/${callId}`)
      .then((res) => setScripts(Array.isArray(res.data) ? res.data : res.data?.items || []))
      .catch(() => {})
      .finally(() => setLoading(false));
  }, [callId]);

  const copyText = (text: string) => {
    navigator.clipboard.writeText(text);
    message.success('已复制');
  };

  return (
    <Card
      title={<Space><BulbOutlined /> 话术推荐</Space>}
      size="small"
      style={{ height: '100%' }}
      bodyStyle={{ maxHeight: 400, overflowY: 'auto' }}
    >
      {scripts.length === 0 ? (
        <Empty description={loading ? '加载中...' : '暂无推荐'} image={Empty.PRESENTED_IMAGE_SIMPLE} />
      ) : (
        <List
          size="small"
          dataSource={scripts}
          renderItem={(item) => (
            <List.Item
              actions={[
                <Button key="copy" size="small" icon={<CopyOutlined />} onClick={() => copyText(item.content)}>复制</Button>,
                <Button key="like" size="small" icon={<LikeOutlined />} type="text" />,
              ]}
            >
              <List.Item.Meta
                title={<Space>{item.title}<Tag color="blue">{item.category}</Tag><Tag>{Math.round(item.confidence * 100)}%</Tag></Space>}
                description={<div style={{ whiteSpace: 'pre-wrap', fontSize: 13 }}>{item.content}</div>}
              />
            </List.Item>
          )}
        />
      )}
    </Card>
  );
}
