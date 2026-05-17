<template>
    <div class="w-full py-2 max-w-[500px]">
        <el-row>
            <TagCodePath :code="progress.authCertName" />
        </el-row>

        <!-- 文件路径 -->
        <div v-if="progress.path" class="mb-3 px-1">
            <span class="text-xs text-gray-500 dark:text-gray-400 truncate block" :title="progress.path">
                {{ progress.path }}
            </span>
        </div>

        <!-- 文件夹信息 -->
        <div class="flex justify-between items-center mb-2">
            <span class="font-semibold text-sm text-gray-700 dark:text-gray-200">{{ progress.folderName }}</span>
            <span class="text-xs text-gray-500 dark:text-gray-400">{{ progress.uploadedFiles }}/{{ progress.totalFiles }}</span>
        </div>

        <!-- 整体进度条 -->
        <el-progress :percentage="percent" :status="progress.status" :stroke-width="10" />

        <!-- 整体进度信息 -->
        <div class="mt-1.5 flex justify-between items-center">
            <span class="text-xs text-gray-500 dark:text-gray-400">{{ formatSize(progress.uploadedSize) }} / {{ formatSize(progress.totalSize) }}</span>
            <span class="text-xs font-semibold text-gray-700 dark:text-gray-200">{{ percent }}%</span>
        </div>

        <!-- 正在上传的文件列表 -->
        <div v-if="progress.uploadingFiles && progress.uploadingFiles.length > 0" class="mt-3 pt-3 border-t border-gray-200 dark:border-gray-700">
            <div class="text-xs font-semibold text-primary mb-2">{{ t('machine.uploading') }} ({{ t('machine.concurrentFiles', { count: progress.uploadingFiles.length }) }}):</div>
            <div v-for="(file, index) in progress.uploadingFiles" :key="index" class="flex items-center gap-1.5 py-1 text-xs text-gray-600 dark:text-gray-300">
                <el-icon class="animate-[rotating_2s_linear_infinite] text-primary"><Loading /></el-icon>
                <span class="flex-1 truncate">{{ file }}</span>
            </div>
        </div>

        <!-- 最后完成的文件 -->
        <div
            v-if="progress.lastFile && progress.status === 'uploading'"
            class="mt-2 flex items-center gap-1.5 px-2 py-1.5 bg-green-50 dark:bg-green-900/20 rounded text-xs"
        >
            <el-icon class="text-green-500"><Check /></el-icon>
            <span class="text-green-600 dark:text-green-400 truncate">{{ progress.lastFile }}</span>
        </div>
    </div>
</template>

<script lang="ts" setup>
import TagCodePath from '@/views/ops/component/TagCodePath.vue';
import { Loading, Check } from '@element-plus/icons-vue';
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

const props = defineProps({
    progress: {
        type: Object,
        required: true,
    },
});

// 计算百分比
const percent = computed(() => {
    if (!props.progress.totalSize || !props.progress.uploadedSize) {
        return 0;
    }
    return Math.min(100, Math.floor((props.progress.uploadedSize / props.progress.totalSize) * 100));
});

// 格式化文件大小
const formatSize = (bytes: number): string => {
    if (!bytes || bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return (bytes / Math.pow(k, i)).toFixed(1) + ' ' + sizes[i];
};
</script>
