import { useState, useEffect, useCallback, useRef } from 'react';
import { Card, Select, Form, Input, InputNumber, Button, Space, Row, Col, message, Tag, Divider, Tooltip } from 'antd';
import { PlusOutlined, DeleteOutlined, SaveOutlined, ZoomInOutlined, ZoomOutOutlined, UndoOutlined } from '@ant-design/icons';
import { ivrFlowApi } from '../../api/endpoints';

const NODE_TYPES = [
  { value: 'start', label: '开始', color: 'green' },
  { value: 'play', label: '播放语音', color: 'blue' },
  { value: 'collect_dtmf', label: '按键收集', color: 'blue' },
  { value: 'condition', label: '条件分支', color: 'orange' },
  { value: 'time_condition', label: '时间条件', color: 'orange' },
  { value: 'variable_assign', label: '变量赋值', color: 'purple' },
  { value: 'transfer_to_skill_group', label: '转技能组', color: 'cyan' },
  { value: 'transfer_to_agent', label: '转坐席', color: 'cyan' },
  { value: 'transfer_to_external', label: '转外线', color: 'cyan' },
  { value: 'blind_transfer', label: '直接转接', color: 'cyan' },
  { value: 'voicemail', label: '语音信箱', color: 'gold' },
  { value: 'callback', label: '排队回呼', color: 'gold' },
  { value: 'hangup', label: '挂机', color: 'red' },
  { value: 'http_request', label: 'HTTP请求', color: 'geekblue' },
  { value: 'json_parser', label: 'JSON解析', color: 'geekblue' },
  { value: 'sms', label: '发短信', color: 'lime' },
  { value: 'satisfaction_rating', label: '满意度', color: 'magenta' },
  { value: 'asr', label: '语音识别', color: 'volcano' },
  { value: 'tts', label: '语音合成', color: 'volcano' },
  { value: 'end', label: '结束', color: 'red' },
];

interface IvrNode {
  id: string;
  type: string;
  label: string;
  config: Record<string, unknown>;
  next: string[];
  x: number;
  y: number;
}

interface CanvasState {
  offsetX: number;
  offsetY: number;
  scale: number;
}

const NODE_W = 160;
const NODE_H = 48;

export default function IvrFlowEditor({ flowId }: { flowId: number }) {
  const [nodes, setNodes] = useState<IvrNode[]>([]);
  const [selectedNode, setSelectedNode] = useState<string | null>(null);
  const [form] = Form.useForm();
  const [history, setHistory] = useState<IvrNode[][]>([]);
  const [canvas, setCanvas] = useState<CanvasState>({ offsetX: 0, offsetY: 0, scale: 1 });
  const [connecting, setConnecting] = useState<string | null>(null);
  const canvasRef = useRef<HTMLDivElement>(null);
  const dragRef = useRef<{ nodeId: string; startX: number; startY: number; origX: number; origY: number } | null>(null);

  const pushHistory = (prev: IvrNode[]) => setHistory((h) => [...h.slice(-20), prev]);

  const loadFlow = useCallback(async () => {
    try {
      const res = await ivrFlowApi.get(flowId);
      const graph = res.data.graph || res.data.flow_graph;
      if (graph?.nodes) {
        const loaded = graph.nodes.map((n: IvrNode, i: number) => ({
          ...n,
          x: n.x ?? 80 + (i % 4) * 200,
          y: n.y ?? 60 + Math.floor(i / 4) * 100,
        }));
        setNodes(loaded);
      }
    } catch { /* */ }
  }, [flowId]);

  useEffect(() => { loadFlow(); }, [loadFlow]);

  const addNode = (type: string) => {
    pushHistory(nodes);
    const id = `node_${Date.now()}`;
    const typeDef = NODE_TYPES.find((t) => t.value === type);
    const x = 80 + (nodes.length % 4) * 200;
    const y = 60 + Math.floor(nodes.length / 4) * 100;
    setNodes([...nodes, { id, type, label: typeDef?.label || type, config: {}, next: [], x, y }]);
  };

  const removeNode = (id: string) => {
    pushHistory(nodes);
    setNodes(nodes.filter((n) => n.id !== id).map((n) => ({ ...n, next: n.next.filter((nid) => nid !== id) })));
    if (selectedNode === id) setSelectedNode(null);
  };

  const updateNodeConfig = (id: string, config: Record<string, unknown>) => {
    setNodes(nodes.map((n) => n.id === id ? { ...n, config } : n));
  };

  const handleSave = async () => {
    try {
      await ivrFlowApi.update(flowId, { graph: { nodes } });
      message.success('保存成功');
    } catch {
      message.error('保存失败');
    }
  };

  const undo = () => {
    if (history.length === 0) return;
    setNodes(history[history.length - 1]);
    setHistory((h) => h.slice(0, -1));
  };

  // Drag node on canvas
  const handleMouseDown = (e: React.MouseEvent, nodeId: string) => {
    if (connecting) {
      // Finish connection
      if (connecting !== nodeId) {
        pushHistory(nodes);
        setNodes(nodes.map((n) => n.id === connecting ? { ...n, next: [...new Set([...n.next, nodeId])] } : n));
      }
      setConnecting(null);
      return;
    }
    e.stopPropagation();
    const node = nodes.find((n) => n.id === nodeId);
    if (!node) return;
    dragRef.current = { nodeId, startX: e.clientX, startY: e.clientY, origX: node.x, origY: node.y };

    const handleMove = (me: MouseEvent) => {
      if (!dragRef.current) return;
      const dx = (me.clientX - dragRef.current.startX) / canvas.scale;
      const dy = (me.clientY - dragRef.current.startY) / canvas.scale;
      setNodes((prev) => prev.map((n) => n.id === dragRef.current!.nodeId
        ? { ...n, x: dragRef.current!.origX + dx, y: dragRef.current!.origY + dy }
        : n
      ));
    };
    const handleUp = () => {
      dragRef.current = null;
      window.removeEventListener('mousemove', handleMove);
      window.removeEventListener('mouseup', handleUp);
    };
    window.addEventListener('mousemove', handleMove);
    window.addEventListener('mouseup', handleUp);
  };

  const selected = nodes.find((n) => n.id === selectedNode);

  // Render edges as SVG lines
  const edges: { from: IvrNode; to: IvrNode }[] = [];
  nodes.forEach((n) => {
    n.next.forEach((nid) => {
      const target = nodes.find((t) => t.id === nid);
      if (target) edges.push({ from: n, to: target });
    });
  });

  return (
    <Row gutter={16} style={{ minHeight: 560 }}>
      {/* Node palette */}
      <Col span={5}>
        <Card title="节点面板" size="small">
          <div style={{ display: 'flex', flexWrap: 'wrap', gap: 4 }}>
            {NODE_TYPES.map((t) => (
              <Tag key={t.value} color={t.color} style={{ cursor: 'pointer', marginBottom: 4 }} onClick={() => addNode(t.value)}>
                <PlusOutlined /> {t.label}
              </Tag>
            ))}
          </div>
        </Card>
      </Col>

      {/* Canvas */}
      <Col span={11}>
        <Card
          title="流程画布"
          size="small"
          extra={
            <Space size={4}>
              <Tooltip title="撤销"><Button size="small" icon={<UndoOutlined />} onClick={undo} disabled={history.length === 0} /></Tooltip>
              <Tooltip title="放大"><Button size="small" icon={<ZoomInOutlined />} onClick={() => setCanvas((c) => ({ ...c, scale: Math.min(c.scale + 0.1, 2) }))} /></Tooltip>
              <Tooltip title="缩小"><Button size="small" icon={<ZoomOutOutlined />} onClick={() => setCanvas((c) => ({ ...c, scale: Math.max(c.scale - 0.1, 0.4) }))} /></Tooltip>
              <Button icon={<SaveOutlined />} type="primary" size="small" onClick={handleSave}>保存</Button>
            </Space>
          }
        >
          <div
            ref={canvasRef}
            style={{
              position: 'relative', height: 480, overflow: 'hidden', background: '#fafafa',
              backgroundImage: 'radial-gradient(circle, #ddd 1px, transparent 1px)',
              backgroundSize: '20px 20px', borderRadius: 4, cursor: connecting ? 'crosshair' : 'default',
            }}
            onClick={() => { if (connecting) setConnecting(null); }}
          >
            <svg style={{ position: 'absolute', inset: 0, width: '100%', height: '100%', pointerEvents: 'none' }}>
              <g transform={`translate(${canvas.offsetX},${canvas.offsetY}) scale(${canvas.scale})`}>
                <defs>
                  <marker id="arrow" viewBox="0 0 10 10" refX="10" refY="5" markerWidth="6" markerHeight="6" orient="auto-start-reverse">
                    <path d="M 0 0 L 10 5 L 0 10 z" fill="#999" />
                  </marker>
                </defs>
                {edges.map((e, i) => {
                  const x1 = e.from.x + NODE_W / 2;
                  const y1 = e.from.y + NODE_H;
                  const x2 = e.to.x + NODE_W / 2;
                  const y2 = e.to.y;
                  const midY = (y1 + y2) / 2;
                  return (
                    <path
                      key={i}
                      d={`M${x1},${y1} C${x1},${midY} ${x2},${midY} ${x2},${y2}`}
                      fill="none" stroke="#999" strokeWidth={1.5} markerEnd="url(#arrow)"
                    />
                  );
                })}
              </g>
            </svg>
            <div style={{ position: 'absolute', inset: 0, transform: `translate(${canvas.offsetX}px,${canvas.offsetY}px) scale(${canvas.scale})`, transformOrigin: '0 0' }}>
              {nodes.map((node) => {
                const typeDef = NODE_TYPES.find((t) => t.value === node.type);
                const isSelected = selectedNode === node.id;
                return (
                  <div
                    key={node.id}
                    style={{
                      position: 'absolute', left: node.x, top: node.y, width: NODE_W,
                      padding: '6px 10px', background: '#fff', border: `2px solid ${isSelected ? '#1677ff' : '#d9d9d9'}`,
                      borderRadius: 6, cursor: 'grab', userSelect: 'none', fontSize: 13,
                      boxShadow: isSelected ? '0 0 8px rgba(22,119,255,0.3)' : '0 1px 3px rgba(0,0,0,0.1)',
                    }}
                    onMouseDown={(e) => handleMouseDown(e, node.id)}
                    onClick={(e) => { e.stopPropagation(); setSelectedNode(node.id); form.setFieldsValue(node.config); }}
                  >
                    <Space size={4}>
                      <Tag color={typeDef?.color} style={{ marginRight: 0, fontSize: 11 }}>{typeDef?.label}</Tag>
                      <span style={{ fontSize: 12 }}>{node.label !== typeDef?.label ? node.label : ''}</span>
                    </Space>
                    {/* Connect handle */}
                    <div
                      style={{
                        position: 'absolute', bottom: -6, left: '50%', transform: 'translateX(-50%)',
                        width: 12, height: 12, borderRadius: '50%', background: connecting === node.id ? '#1677ff' : '#d9d9d9',
                        border: '2px solid #fff', cursor: 'pointer',
                      }}
                      onMouseDown={(e) => { e.stopPropagation(); setConnecting(node.id); }}
                    />
                    <div
                      style={{ position: 'absolute', top: 2, right: 4, cursor: 'pointer', fontSize: 11, color: '#ff4d4f' }}
                      onClick={(e) => { e.stopPropagation(); removeNode(node.id); }}
                    >
                      <DeleteOutlined />
                    </div>
                  </div>
                );
              })}
            </div>
            {connecting && (
              <div style={{ position: 'absolute', top: 8, left: '50%', transform: 'translateX(-50%)', background: '#1677ff', color: '#fff', padding: '2px 12px', borderRadius: 4, fontSize: 12 }}>
                点击目标节点完成连线 (ESC取消)
              </div>
            )}
          </div>
        </Card>
      </Col>

      {/* Config panel */}
      <Col span={8}>
        <Card title={selected ? `配置: ${selected.label}` : '选择节点'} size="small">
          {selected ? (
            <Form form={form} layout="vertical" onValuesChange={(_, all) => updateNodeConfig(selected.id, all)}>
              <Form.Item name="label" label="标签"><Input placeholder="节点标签" /></Form.Item>
              {(selected.type === 'play' || selected.type === 'tts') && (
                <Form.Item name="audio_file_id" label="音频/文本"><Input placeholder="音频文件ID或TTS文本" /></Form.Item>
              )}
              {selected.type === 'collect_dtmf' && (
                <>
                  <Form.Item name="max_digits" label="最大位数"><InputNumber min={1} max={20} /></Form.Item>
                  <Form.Item name="timeout_sec" label="超时(秒)"><InputNumber min={1} max={30} /></Form.Item>
                  <Form.Item name="finish_key" label="结束键"><Input placeholder="# 或 *" /></Form.Item>
                </>
              )}
              {selected.type === 'transfer_to_skill_group' && (
                <Form.Item name="skill_group_id" label="技能组ID"><InputNumber /></Form.Item>
              )}
              {selected.type === 'transfer_to_agent' && (
                <Form.Item name="agent_id" label="坐席ID"><InputNumber /></Form.Item>
              )}
              {selected.type === 'transfer_to_external' && (
                <Form.Item name="external_number" label="外线号码"><Input placeholder="外线号码" /></Form.Item>
              )}
              {selected.type === 'http_request' && (
                <>
                  <Form.Item name="url" label="URL"><Input /></Form.Item>
                  <Form.Item name="method" label="方法"><Select options={[{ value: 'GET' }, { value: 'POST' }, { value: 'PUT' }]} /></Form.Item>
                  <Form.Item name="headers" label="请求头(JSON)"><Input.TextArea rows={2} /></Form.Item>
                  <Form.Item name="body" label="请求体"><Input.TextArea rows={2} /></Form.Item>
                </>
              )}
              {selected.type === 'json_parser' && (
                <Form.Item name="json_path" label="JSON Path"><Input placeholder="$.data.result" /></Form.Item>
              )}
              {selected.type === 'condition' && (
                <Form.Item name="expression" label="条件表达式"><Input placeholder="e.g. ${dtmf} == 1" /></Form.Item>
              )}
              {selected.type === 'time_condition' && (
                <>
                  <Form.Item name="business_hours_id" label="营业时间ID"><InputNumber /></Form.Item>
                  <Form.Item name="timezone" label="时区"><Input placeholder="Asia/Shanghai" /></Form.Item>
                </>
              )}
              {selected.type === 'variable_assign' && (
                <>
                  <Form.Item name="variable_name" label="变量名"><Input /></Form.Item>
                  <Form.Item name="variable_value" label="值"><Input /></Form.Item>
                </>
              )}
              {selected.type === 'sms' && (
                <>
                  <Form.Item name="template_id" label="短信模板ID"><Input /></Form.Item>
                  <Form.Item name="sign_name" label="签名"><Input /></Form.Item>
                </>
              )}
              {selected.type === 'voicemail' && (
                <>
                  <Form.Item name="max_duration_sec" label="最大录音(秒)"><InputNumber min={5} max={300} /></Form.Item>
                  <Form.Item name="greeting_audio_id" label="问候语音ID"><Input /></Form.Item>
                </>
              )}
              {selected.type === 'callback' && (
                <Form.Item name="callback_queue_id" label="回呼队列ID"><InputNumber /></Form.Item>
              )}
              {selected.type === 'satisfaction_rating' && (
                <Form.Item name="survey_type" label="调研类型">
                  <Select options={[{ value: 'dtmf', label: '按键评分' }, { value: 'voice', label: '语音评分' }]} />
                </Form.Item>
              )}
              {selected.type === 'asr' && (
                <>
                  <Form.Item name="language" label="语言"><Input placeholder="zh-CN" /></Form.Item>
                  <Form.Item name="hotword_group_id" label="热词组ID"><InputNumber /></Form.Item>
                </>
              )}
              <Divider>连线</Divider>
              <Form.Item label="连接到">
                <Select
                  mode="multiple"
                  value={selected.next}
                  onChange={(vals: string[]) => { pushHistory(nodes); setNodes(nodes.map((n) => n.id === selected.id ? { ...n, next: vals } : n)); }}
                  options={nodes.filter((n) => n.id !== selected.id).map((n) => ({ value: n.id, label: `${n.label} (${n.type})` }))}
                />
              </Form.Item>
            </Form>
          ) : (
            <div style={{ color: '#999', textAlign: 'center', padding: 40 }}>点击画布中的节点进行配置</div>
          )}
        </Card>
      </Col>
    </Row>
  );
}
