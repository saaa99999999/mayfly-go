import { reactive } from 'vue';

// 活跃通知任务映射表
export const activeNotifications = reactive<Map<string, any>>(new Map());

// 悬浮通知状态
export const globalNotificationState = reactive({
    activeCount: 0,
});

/**
 * 更新悬浮通知状态
 */
const updateNotificationState = () => {
    globalNotificationState.activeCount = activeNotifications.size;
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
