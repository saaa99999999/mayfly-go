import DbSqlExecProgress from './DbSqlExecProgress.vue';
import { createOrUpdateNotification, completeNotification, activeNotifications } from '../global-notification-manager';
import syssocket from '@/common/syssocket';
import { reactive, nextTick } from 'vue';

// 存储SQL执行任务的取消方法
const sqlExecAborters = new Map<string, { abort: () => void; progress?: any }>();

// 存储待注册的 abort 方法（等待 WebSocket 消息到达）
const pendingSqlExecAborters = new Map<string, () => void>();

export async function registerDbSqlExecProgress() {
    await syssocket.registerMsgHandler('sqlScriptRunProgress', function (message: any) {
        const content = message.params;
        const id = content.id;

        // SQL执行完成
        if (content.terminated) {
            completeNotification(id, 1000);
            sqlExecAborters.delete(id);
            return;
        }

        // 构建组件props
        const props = {
            progress: reactive({
                title: content.title || '',
                executedStatements: content.executedStatements || 0,
                terminated: content.terminated || false,
                status: content.status || '',
                dbCode: content.dbCode || '',
                dbName: content.dbName || '',
            }),
            onCancel: () => {
                const aborter = sqlExecAborters.get(id);
                if (aborter) {
                    aborter.abort();

                    // 更新通知状态为取消
                    if (aborter.progress) {
                        nextTick(() => {
                            aborter.progress.status = 'cancelled';
                            aborter.progress.terminated = true;
                        });

                        // 延迟后关闭通知
                        setTimeout(() => {
                            completeNotification(id, 1000);
                            sqlExecAborters.delete(id);
                        }, 1500);
                    } else {
                        sqlExecAborters.delete(id);
                    }
                }
            },
        };

        // 创建或更新通知
        createOrUpdateNotification(id, 'sqlScriptRun', content, DbSqlExecProgress, props, {
            title: message.title || 'db.sqlExecute',
            onCancel: props.onCancel,
        });

        // 如果有待注册的 abort 方法，现在注册
        const pendingAbort = pendingSqlExecAborters.get(id);
        if (pendingAbort) {
            sqlExecAborters.set(id, { abort: pendingAbort, progress: props.progress });
            pendingSqlExecAborters.delete(id);
        }
    });
}

/**
 * 注册SQL执行任务的取消方法
 * @param execId 执行ID
 * @param abort 取消方法
 */
export function registerSqlExecAborter(execId: string, abort: () => void) {
    // 先检查通知是否已经存在
    const task = activeNotifications.get(execId);
    const progress = task?.componentProps?.progress || null;

    if (progress) {
        // 通知已存在，直接注册
        sqlExecAborters.set(execId, { abort, progress });
    } else {
        // 通知还未创建，保存为 pending
        pendingSqlExecAborters.set(execId, abort);
    }
}
