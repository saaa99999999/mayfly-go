import syssocket from '@/common/syssocket';
import { createOrUpdateNotification, completeNotification, activeNotifications } from '../global-notification-manager';
import MachineFileUploadProgress from './MachineFileUploadProgress.vue';
import { reactive, nextTick } from 'vue';

// 存储上传任务的取消方法
const uploadAborters = new Map<string, { abort: () => void; progress?: any }>();

// 存储待注册的 abort 方法（等待 WebSocket 消息到达）
const pendingAborters = new Map<string, () => void>();

/**
 * 注册机器文件上传进度消息处理
 */
export async function registerMachineFileUploadProgress() {
    await syssocket.registerMsgHandler('machineFileUploadProgress', function (message: any) {
        const content = message.params;
        const uploadId = content.uploadId;

        // 上传完成或失败
        if (content.status === 'complete' || content.status === 'error') {
            completeNotification(uploadId, 1000);
            uploadAborters.delete(uploadId);
            return;
        }

        // 构建组件props
        const props = {
            progress: reactive({
                authCertName: content.authCertName || '',
                path: content.path || '',
                filename: content.filename || '',
                uploadedSize: content.uploadedSize || 0,
                totalSize: content.totalSize || 0,
                timestamp: content.timestamp || 0,
                status: content.status || 'uploading',
            }),
            onCancel: () => {
                const aborter = uploadAborters.get(uploadId);
                if (aborter) {
                    aborter.abort();

                    // 更新通知状态为取消
                    if (aborter.progress) {
                        nextTick(() => {
                            aborter.progress.status = 'error';
                            aborter.progress.filename = '已取消: ' + (aborter.progress.filename || '');
                        });

                        // 延迟后关闭通知
                        setTimeout(() => {
                            completeNotification(uploadId, 1000);
                            uploadAborters.delete(uploadId);
                        }, 1500);
                    } else {
                        uploadAborters.delete(uploadId);
                    }
                }
            },
        };

        // 创建或更新上传通知
        createOrUpdateNotification(uploadId, 'machineFileUpload', content, MachineFileUploadProgress, props, {
            title: message.title || 'machine.fileUpload',
        });

        // 如果有待注册的 abort 方法，现在注册
        const pendingAbort = pendingAborters.get(uploadId);
        if (pendingAbort) {
            console.log('[MachineFileUpload] Registering pending aborter for uploadId:', uploadId);
            uploadAborters.set(uploadId, { abort: pendingAbort, progress: props.progress });
            pendingAborters.delete(uploadId);
        }
    });
}

/**
 * 注册上传任务的取消方法
 * @param uploadId 上传ID
 * @param abort 取消方法
 */
export function registerUploadAborter(uploadId: string, abort: () => void) {
    // 先检查通知是否已经存在
    const task = activeNotifications.get(uploadId);
    const progress = task?.componentProps?.progress || null;

    if (progress) {
        // 通知已存在，直接注册
        uploadAborters.set(uploadId, { abort, progress });
    } else {
        // 通知还未创建，保存为 pending
        pendingAborters.set(uploadId, abort);
    }
}
