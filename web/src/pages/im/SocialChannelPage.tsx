import { useState, useEffect } from 'react';
import { Table, Tag, Button, Modal, Form, Input, Select, Space, Switch, message, Popconfirm, Badge } from 'antd';
import {
  WechatOutlined, WeiboOutlined, PlusOutlined, EditOutlined, DeleteOutlined,
  ReloadOutlined,
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import api from '../../api/client';

interface SocialChannel {
  id: number;
  name: string;
  platform: 'wechat' | 'weibo';
  app_id: string;
  status: 'connected' | 'disconnected' | 'error';
  enabled: boolean;
  follower_count: number;
  last_sync: string;
  created_at: string;
}

const platformConfig = {
  wechat: { label: '微信公众号', icon: <WechatOutlined />, color: '#07c160' },
  weibo:  { label: '微博',      icon: <WeiboOutlined />,  color: '#ff8200' },
};

export default function SocialChannelPage() {
  const [channels, setChannels] = useState<SocialChannel[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalOpen, setModalOpen] = useState(false);
  const [editing, setEditing] = useState<SocialChannel | null>(null);
  const [form] = Form.useForm();

  const load = async () => {
    setLoading(true);
    try {
      const res = await api.get('/social-channels');
      setChannels(Array.isArray(res.data) ? res.data : res.data?.items || []);
    } catch { /* ignore */ }
    setLoading(false);
  };

  useEffect(() => { load(); }, []);

  const handleSave = async () => {
    const values = await form.validateFields();
    try {
      if (editing) {
        await api.put(`/social-channels/${editing.id}`, values);
        message.success('更新成功');
      } else {
        await api.post('/social-channels', values);
        message.success('创建成功');
      }
      setModalOpen(false);
      form.resetFields();
      setEditing(null);
      load();
    } catch {
      message.error('操作失败');
    }
  };

  const handleDelete = async (id: number) => {
    try {
      await api.delete(`/social-channels/${id}`);
      message.success('已删除');
      load();
    } catch {
      message.error('删除失败');
    }
  };

  const handleToggle = async (id: number, enabled: boolean) => {
    try {
      await api.put(`/social-channels/${id}`, { enabled });
      load();
    } catch {
      message.error('操作失败');
    }
  };

  const columns: ColumnsType<SocialChannel> = [
    {
      title: '平台', dataIndex: 'platform', width: 120,
      render: (v: keyof typeof platformConfig) => {
        const cfg = platformConfig[v];
        return <Tag icon={cfg.icon} color={cfg.color}>{cfg.label}</Tag>;
      },
    },
    { title: '名称', dataIndex: 'name' },
    { title: 'App ID', dataIndex: 'app_id', width: 180, ellipsis: true },
    {
      title: '状态', dataIndex: 'status', width: 100,
      render: (v) => v === 'connected'
        ? <Badge status="success" text="已连接" />
        : v === 'error'
          ? <Badge status="error" text="异常" />
          : <Badge status="default" text="未连接" />,
    },
    {
      title: '启用', dataIndex: 'enabled', width: 80,
      render: (v, record) => <Switch checked={v} onChange={(checked) => handleToggle(record.id, checked)} />,
    },
    { title: '粉丝数', dataIndex: 'follower_count', width: 80 },
    { title: '最后同步', dataIndex: 'last_sync', width: 160 },
    {
      title: '操作', key: 'actions', width: 160,
      render: (_, record) => (
        <Space>
          <Button size="small" icon={<EditOutlined />} onClick={() => { setEditing(record); form.setFieldsValue(record); setModalOpen(true); }}>编辑</Button>
          <Popconfirm title="确定删除？" onConfirm={() => handleDelete(record.id)}>
            <Button size="small" danger icon={<DeleteOutlined />}>删除</Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <>
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
        <h2 style={{ margin: 0 }}>社交渠道管理</h2>
        <Space>
          <Button icon={<ReloadOutlined />} onClick={load}>刷新</Button>
          <Button type="primary" icon={<PlusOutlined />} onClick={() => { setEditing(null); form.resetFields(); setModalOpen(true); }}>新建</Button>
        </Space>
      </div>
      <Table<SocialChannel> columns={columns} dataSource={channels} rowKey="id" loading={loading} size="middle" pagination={{ pageSize: 20 }} />

      <Modal title={editing ? '编辑渠道' : '新建渠道'} open={modalOpen} onOk={handleSave} onCancel={() => { setModalOpen(false); setEditing(null); form.resetFields(); }} destroyOnClose>
        <Form form={form} layout="vertical" preserve={false}>
          <Form.Item name="platform" label="平台" rules={[{ required: true }]}>
            <Select options={[
              { value: 'wechat', label: '微信公众号' },
              { value: 'weibo', label: '微博' },
            ]} />
          </Form.Item>
          <Form.Item name="name" label="名称" rules={[{ required: true }]}>
            <Input placeholder="渠道名称" />
          </Form.Item>
          <Form.Item name="app_id" label="App ID" rules={[{ required: true }]}>
            <Input placeholder="应用ID" />
          </Form.Item>
          <Form.Item name="app_secret" label="App Secret" rules={[{ required: true }]}>
            <Input.Password placeholder="应用密钥" />
          </Form.Item>
          <Form.Item name="token" label="Token">
            <Input placeholder="验证 Token" />
          </Form.Item>
          <Form.Item name="encoding_aes_key" label="EncodingAESKey">
            <Input placeholder="消息加解密密钥 (43位)" />
          </Form.Item>
          <Form.Item name="enabled" label="启用" valuePropName="checked">
            <Switch />
          </Form.Item>
        </Form>
      </Modal>
    </>
  );
}
