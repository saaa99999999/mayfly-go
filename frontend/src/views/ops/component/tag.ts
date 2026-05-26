import { ContextmenuItem } from '@/components/contextmenu';
import { markRaw } from 'vue';

// 资源配置
export interface ResourceConfig {
    order?: number;
    resourceType: number; // 资源类型
    rootNodeType: NodeType; // 资源根节点类型

    // 资源管理组件配置
    manager?: {
        componentConf: ResourceComponentConfig; // 组件
        countKey?: string; // 统计数key，tab展示的数字对象key
        permCode?: string; // 权限码
    };
}

export interface ResourceComponentConfig {
    name: string; // 名称
    component?: any; // 组件
    icon?: {
        name: string;
        color?: string;
    };
}

export interface ResourceOpCtx {
    /**
     * 添加资源相关组件
     * @param component 资源相关组件配置
     * @returns 组件引用
     */
    addResourceComponent(component: ResourceComponentConfig): Promise<any>;

    /**
     * 获取树节点
     * @param nodeKey 节点key
     */
    getTreeNode(nodeKey: string): any;

    setCurrentTreeKey(nodeKey: string): void;

    reloadTreeNode(nodeKey: string): void;
}

export class TagTreeNode {
    /**
     * 节点id
     */
    key: any;

    /**
     * 节点名称
     */
    label: string;

    /**
     * 节点名称备注（用于元素title属性）
     */
    labelRemark: string;

    /**
     * 树节点类型
     */
    type: NodeType;

    /**
     * 是否为叶子节点
     */
    isLeaf: boolean = false;

    /**
     * 是否禁用状态
     */
    disabled: boolean = false;

    /**
     * 额外需要传递的参数
     */
    params: any;

    icon: any;

    // 节点组件
    nodeComponent?: any;

    /**
     * 节点上下文
     */
    ctx?: ResourceOpCtx;

    static TagPath = -1;

    constructor(key: any, label: string, type?: NodeType) {
        this.key = key;
        this.label = label;
        this.type = type || new NodeType(TagTreeNode.TagPath);
    }

    static new(parent: TagTreeNode, key: any, label: string, type?: NodeType) {
        return new TagTreeNode(key, label, type).withContext(parent.ctx);
    }

    withLabelRemark(labelRemark: any) {
        this.labelRemark = labelRemark;
        return this;
    }

    withIsLeaf(isLeaf: boolean) {
        this.isLeaf = isLeaf;
        return this;
    }

    withDisabled(disabled: boolean) {
        this.disabled = disabled;
        return this;
    }

    withParams(params: any) {
        this.params = params;
        return this;
    }

    withIcon(icon: any) {
        this.icon = icon;
        return this;
    }

    withNodeComponent(component: any) {
        this.nodeComponent = markRaw(component);
        return this;
    }

    withContext(ctx: ResourceOpCtx | undefined | null) {
        if (!ctx) {
            return this;
        }
        this.ctx = ctx;
        return this;
    }

    /**
     * 加载子节点，使用节点类型的loadNodesFunc去加载子节点
     * @returns 子节点信息
     */
    async loadChildren() {
        if (this.isLeaf) {
            return [];
        }
        if (this.type && this.type.loadNodesFunc) {
            return await this.type.loadNodesFunc(this);
        }
        return [];
    }
}

/**
 * 节点类型，用于加载子节点及点击事件等
 */
export class NodeType {
    /**
     * 节点类型值
     */
    value: number;

    contextMenuItems: ContextmenuItem[];

    loadNodesFunc: (parentNode: TagTreeNode) => Promise<TagTreeNode[]>;

    /**
     * 节点点击事件
     */
    nodeClickFunc: (node: TagTreeNode) => void;

    // 节点双击事件
    nodeDblclickFunc?: (node: TagTreeNode) => void;

    constructor(value: number) {
        this.value = value;
    }

    /**
     * 赋值加载子节点回调函数
     * @param func 加载子节点回调函数
     * @returns this
     */
    withLoadNodesFunc(func: (parentNode: TagTreeNode) => Promise<TagTreeNode[]>) {
        this.loadNodesFunc = func;
        return this;
    }

    /**
     * 赋值节点点击事件回调函数
     * @param func 节点点击事件回调函数
     * @returns this
     */
    withNodeClickFunc(func: (node: TagTreeNode) => void) {
        this.nodeClickFunc = func;
        return this;
    }

    /**
     * 赋值节点双击事件回调函数
     * @param func 节点双击事件回调函数
     * @returns this
     */
    withNodeDblclickFunc(func: (node: TagTreeNode) => void) {
        this.nodeDblclickFunc = func;
        return this;
    }

    /**
     * 赋值右击菜单按钮选项
     * @param contextMenuItems 右击菜单按钮选项
     * @returns this
     */
    withContextMenuItems(contextMenuItems: ContextmenuItem[]) {
        this.contextMenuItems = contextMenuItems;
        return this;
    }
}

export function expandCodePath(codePath: string) {
    const parts = codePath.split('/');
    const result = [];
    let currentPath = '';

    for (let i = 0; i < parts.length - 1; i++) {
        currentPath += parts[i] + '/';
        result.push(currentPath);
    }

    return result;
}
