import { useState, useEffect, useMemo } from 'react';
import { Card, Input, List, Tag, Button, Empty, Space, message } from 'antd';
import { SearchOutlined, SendOutlined, CopyOutlined, ThunderboltOutlined } from '@ant-design/icons';
import api from '../../api/client';

interface QuickReply {
  id: number;
  title: string;
  content: string;
  category: string;
  shortcut?: string;
}

interface QuickReplyPickerProps {
  onSelect?: (content: string) => void;
}

export default function QuickReplyPicker({ onSelect }: QuickReplyPickerProps) {
  const [replies, setReplies] = useState<QuickReply[]>([]);
  const [search, setSearch] = useState('');
  const [category, setCategory] = useState<string>('');

  useEffect(() => {
    api.get('/quick-replies')
      .then((res) => setReplies(Array.isArray(res.data) ? res.data : res.data?.items || []))
      .catch(() => {});
  }, []);

  const categories = useMemo(() => [...new Set(replies.map((r) => r.category))].filter(Boolean), [replies]);

  const filtered = useMemo(() => {
    let list = replies;
    if (category) list = list.filter((r) => r.category === category);
    if (search) {
      const q = search.toLowerCase();
      list = list.filter((r) => r.title.toLowerCase().includes(q) || r.content.toLowerCase().includes(q));
    }
    return list;
  }, [replies, search, category]);

  const handleCopy = (text: string) => {
    navigator.clipboard.writeText(text);
    message.success('已复制');
  };

  return (
    <Card title={<Space><ThunderboltOutlined /> 快捷回复</Space>} size="small" bodyStyle={{ padding: '8px 12px' }}>
      <Input
        prefix={<SearchOutlined />}
        placeholder="搜索快捷回复..."
        value={search}
        onChange={(e) => setSearch(e.target.value)}
        allowClear
        style={{ marginBottom: 8 }}
      />
      <div style={{ marginBottom: 8 }}>
        <Tag color={!category ? 'blue' : undefined} style={{ cursor: 'pointer' }} onClick={() => setCategory('')}>全部</Tag>
        {categories.map((c) => (
          <Tag key={c} color={category === c ? 'blue' : undefined} style={{ cursor: 'pointer' }} onClick={() => setCategory(c)}>{c}</Tag>
        ))}
      </div>
      <div style={{ maxHeight: 300, overflowY: 'auto' }}>
        {filtered.length === 0 ? (
          <Empty description="无匹配" image={Empty.PRESENTED_IMAGE_SIMPLE} />
        ) : (
          <List
            size="small"
            dataSource={filtered}
            renderItem={(item) => (
              <List.Item
                style={{ padding: '6px 0', cursor: 'pointer' }}
                actions={[
                  <Button key="copy" size="small" type="text" icon={<CopyOutlined />} onClick={() => handleCopy(item.content)} />,
                  onSelect ? <Button key="send" size="small" type="text" icon={<SendOutlined />} onClick={() => onSelect(item.content)} /> : null,
                ].filter(Boolean)}
                onClick={() => onSelect?.(item.content)}
              >
                <List.Item.Meta
                  title={<Space size={4}>{item.title}{item.shortcut && <Tag style={{ fontSize: 11 }}>{item.shortcut}</Tag>}</Space>}
                  description={<div style={{ fontSize: 12, color: '#666', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis', maxWidth: 220 }}>{item.content}</div>}
                />
              </List.Item>
            )}
          />
        )}
      </div>
    </Card>
  );
}
