import { reactive } from 'vue';

// 活跃通知任务映射表
export const activeNotifications = reactive<Map<string, any>>(new Map());

// 悬浮通知状态
export const globalNotificationState = reactive({
    hasActiveNotifications: false,
    activeCount: 0,
    // 按类别统计
    categoryCount: reactive<Map<string, number>>(new Map()),
});

/**
 * 更新悬浮通知状态
 */
const updateNotificationState = () => {
    globalNotificationState.activeCount = activeNotifications.size;
    globalNotificationState.hasActiveNotifications = activeNotifications.size > 0;

    // 按类别统计
    const categoryMap = new Map<string, number>();
    for (const [_, task] of activeNotifications) {
        const category = task.category || 'default';
        categoryMap.set(category, (categoryMap.get(category) || 0) + 1);
    }
    globalNotificationState.categoryCount.clear();
    for (const [key, value] of categoryMap) {
        globalNotificationState.categoryCount.set(key, value);
    }
};

/**
 * 创建或更新通知
 * @param id 通知唯一ID
 * @param category 通知类别(如:machineFileUpload, machineFolderUpload, sqlScript等)
 * @param content 通知内容
 * @param component 通知组件
 * @param componentProps 组件props
 * @param options 通知选项
 */
export const createOrUpdateNotification = (
    id: string,
    category: string,
    content: any,
    component: any,
    componentProps: any,
    options: {
        title: string;
        onCancel?: () => void; // 取消回调
    }
) => {
    // 添加到活跃任务
    activeNotifications.set(id, {
        id,
        category,
        content,
        component,
        componentProps,
        options,
        timestamp: Date.now(),
    });

    updateNotificationState();
};

/**
 * 完成通知
 * @param id 通知唯一ID
 * @param closeDelay 延迟关闭时间（毫秒）
 */
export const completeNotification = (id: string, closeDelay: number = 1000) => {
    // 延迟从活跃列表中移除
    setTimeout(() => {
        activeNotifications.delete(id);
        updateNotificationState();
    }, closeDelay);
};

/**
 * 关闭指定通知
 * @param id 通知唯一ID
 */
export const closeNotification = (id: string) => {
    activeNotifications.delete(id);
    updateNotificationState();
};

/**
 * 关闭指定类别的所有通知
 * @param category 通知类别
 */
export const closeCategoryNotifications = (category: string) => {
    for (const [id, task] of activeNotifications) {
        if (task.category === category) {
            closeNotification(id);
        }
    }
};

/**
 * 关闭所有通知
 */
export const closeAllNotifications = () => {
    activeNotifications.clear();
    updateNotificationState();
};
