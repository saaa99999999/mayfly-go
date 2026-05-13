<template>
    <span v-for="item in path" :key="item.code">
        <SvgIcon
            :name="EnumValue.getEnumByValue(TagResourceTypeEnum, item.type)?.extra.icon"
            :color="EnumValue.getEnumByValue(TagResourceTypeEnum, item.type)?.extra.iconColor"
            class="mr-0.5"
            :size="iconSize"
        />
        <span :class="textClass"> {{ item.name ? item.name : item.code }}</span>

        <SvgIcon v-if="!item.isEnd" color="var(--el-text-color-placeholder)" :size="iconSize" :class="arrowMarginClass" name="arrow-right" />
    </span>
</template>

<script lang="ts" setup>
import { TagResourceTypeEnum } from '@/common/commonEnum';
import EnumValue from '@/common/Enum';
import { computed } from 'vue';

const props = defineProps({
    // 路径数据
    path: {
        type: Array<any>,
        required: true,
        default: () => [],
    },
    // 尺寸: normal(默认) | small
    size: {
        type: String,
        default: 'small',
    },
});

const iconSize = computed(() => {
    return props.size === 'small' ? 14 : 15;
});

const textClass = computed(() => {
    return props.size === 'small' ? 'text-sm' : '';
});

const arrowMarginClass = computed(() => {
    return props.size === 'small' ? 'mx-0.5' : 'mx-1';
});
</script>
