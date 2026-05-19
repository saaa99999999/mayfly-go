import syssocket from '@/common/syssocket';
import { createOrUpdateNotification, completeNotification, activeNotifications } from '../global-notification-manager';
import MachineFolderUploadProgress from './MachineFolderUploadProgress.vue';
import { reactive, nextTick } from 'vue';

// 存储上传任务的取消方法
const folderUploadAborters = new Map<string, { abort: () => void; progress?: any }>();

// 存储待注册的 abort 方法（等待 WebSocket 消息到达）
const pendingFolderAborters = new Map<string, () => void>();

/**
 * 注册文件夹上传进度消息处理
 */
export async function registerFolderUploadProgressHandler() {
    await syssocket.registerMsgHandler('machineFolderUploadProgress', function (message: any) {
        const content = message.params;
        const uploadId = content.uploadId;

        if (!uploadId) {
            return;
        }

        // 上传完成或失败
        if (content.status === 'complete' || content.status === 'error') {
            completeNotification(uploadId, 1000);
            folderUploadAborters.delete(uploadId);
            return;
        }

        // 构建组件props
        const props = {
            progress: reactive({
                authCertName: content.authCertName || '',
                path: content.path || '',
                folderName: content.folderName || '',
                totalFiles: content.totalFiles || 0,
                uploadedFiles: content.uploadedFiles || 0,
                totalSize: content.totalSize || 0,
                uploadedSize: content.uploadedSize || 0,
                uploadingFiles: content.uploadingFiles || [],
                timestamp: content.timestamp || 0,
                status: content.status || 'uploading',
            }),
            onCancel: () => {
                const aborter = folderUploadAborters.get(uploadId);
                if (aborter) {
                    aborter.abort();

                    // 更新通知状态为取消
                    if (aborter.progress) {
                        nextTick(() => {
                            aborter.progress.status = 'exception';
                            aborter.progress.folderName = '已取消: ' + (aborter.progress.folderName || '');
                        });

                        // 延迟后关闭通知
                        setTimeout(() => {
                            completeNotification(uploadId, 1000);
                            folderUploadAborters.delete(uploadId);
                        }, 1500);
                    } else {
                        folderUploadAborters.delete(uploadId);
                    }
                }
            },
        };

        // 创建或更新上传通知
        if (content.status === 'uploading') {
            createOrUpdateNotification(uploadId, 'machineFolderUpload', content, MachineFolderUploadProgress, props, {
                title: message.title || 'machine.folderUpload',
            });
        }

        // 如果有待注册的 abort 方法，现在注册
        const pendingAbort = pendingFolderAborters.get(uploadId);
        if (pendingAbort) {
            folderUploadAborters.set(uploadId, { abort: pendingAbort, progress: props.progress });
            pendingFolderAborters.delete(uploadId);
        }
    });
}

/**
 * 注册文件夹上传任务的取消方法
 * @param uploadId 上传ID
 * @param abort 取消方法
 */
export function registerFolderUploadAborter(uploadId: string, abort: () => void) {
    // 先检查通知是否已经存在
    const task = activeNotifications.get(uploadId);
    const progress = task?.componentProps?.progress || null;

    if (progress) {
        // 通知已存在，直接注册
        folderUploadAborters.set(uploadId, { abort, progress });
    } else {
        // 通知还未创建，保存为 pending
        pendingFolderAborters.set(uploadId, abort);
    }
}
