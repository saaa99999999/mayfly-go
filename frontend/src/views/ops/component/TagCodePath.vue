<template>
    <el-row v-for="(path, idx) in codePaths?.slice(0, 1)" :key="idx">
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

        <!-- 展示剩余的标签信息 -->
        <el-popover :show-after="300" v-if="paths.length > 1 && idx == 0" placement="bottom" :width="popoverWidth" trigger="hover">
            <template #reference>
                <SvgIcon :size="iconSize" :class="moreIconMarginClass" color="var(--el-color-primary)" name="MoreFilled" />
            </template>

            <el-row v-for="(opath, oi) in codePaths.slice(1)" :key="oi" class="mb-2">
                <span v-for="item in opath" :key="item.code">
                    <SvgIcon
                        :name="EnumValue.getEnumByValue(TagResourceTypeEnum, item.type)?.extra.icon"
                        :color="EnumValue.getEnumByValue(TagResourceTypeEnum, item.type)?.extra.iconColor"
                        class="mr-0.5"
                        :size="iconSize"
                    />
                    <span :class="textClass"> {{ item.name ? item.name : item.code }}</span>
                    <SvgIcon v-if="!item.isEnd" color="var(--el-text-color-placeholder)" :size="iconSize" :class="arrowMarginClass" name="arrow-right" />
                </span>
            </el-row>
        </el-popover>
    </el-row>
</template>

<script lang="ts" setup>
import { TagResourceTypeEnum } from '@/common/commonEnum';
import EnumValue from '@/common/Enum';
import { tagApi } from '@/views/ops/tag/api';
import { computed, onMounted, ref, watch } from 'vue';

const props = defineProps({
    path: {
        type: [String, Array<string>, Array<Object>],
    },
    // code，可直接设置该值展示路径信息
    code: {
        type: String,
    },
    // 尺寸: normal(默认) | small
    size: {
        type: String,
        default: 'small',
    },
});

const codePath = ref(props.path);
const codePaths: any = ref([]);
let allTagInfos: any = {};

const iconSize = computed(() => {
    return props.size === 'small' ? 14 : 15;
});

const textClass = computed(() => {
    return props.size === 'small' ? 'text-sm' : '';
});

const arrowMarginClass = computed(() => {
    return props.size === 'small' ? 'mx-0.5' : 'mx-1';
});

const moreIconMarginClass = computed(() => {
    return props.size === 'small' ? 'mt-2 ml-1' : 'mt-1 ml-1';
});

const popoverWidth = computed(() => {
    return props.size === 'small' ? 400 : 500;
});

const paths = computed(() => {
    if (Array.isArray(codePath.value)) {
        const ps = [];
        // 兼容["default/test1/test2/"] 与 [{id: 1, codePath: "default/test1/test2/"}]
        for (let p of codePath.value as any) {
            if (typeof p === 'string') {
                ps.push(p);
            } else {
                ps.push(p.codePath);
            }
        }
        return ps;
    }

    return [codePath.value];
});

onMounted(() => {
    codePath.value = props.path;
    setCodePaths();
});

watch(
    () => props.path,
    () => {
        codePath.value = props.path;
        setCodePaths();
    }
);

watch(
    () => props.code,
    () => {
        if (!props.code) {
            clear();
            return;
        }
        setCodePaths();
    }
);

const setCodePaths = async () => {
    if (props.code) {
        const tagInfos = await tagApi.listByQuery.request({ codes: props.code });
        if (tagInfos.length == 0) {
            clear();
            return;
        }
        codePath.value = tagInfos[0].codePath;
    }

    if (!paths.value) {
        clear();
        return;
    }

    allTagInfos = await getAllCodePaths(paths.value as any);
    codePaths.value = paths.value.map((p) => parseTagPath(p));
};

const clear = () => {
    codePath.value = '';
    codePaths.value = [];
};

const parseTagPath = (tagPath: string = '') => {
    if (!tagPath) {
        return [];
    }
    const res = [] as any;
    let codePath = '';
    const codes = tagPath.split('/');
    for (let code of codes) {
        codePath += code + '/';
        const typeAndCode = code.split('|');

        let tagInfo;
        if (typeAndCode.length == 1) {
            const tagCode = typeAndCode[0];
            if (!tagCode) {
                continue;
            }

            tagInfo = {
                type: TagResourceTypeEnum.Tag.value,
                code: typeAndCode[0],
                codePath: codePath,
                name: '',
            };
        } else {
            tagInfo = {
                type: typeAndCode[0],
                code: typeAndCode[1],
                codePath: codePath,
                name: '',
            };
        }

        const ti = allTagInfos[codePath];
        if (ti) {
            tagInfo.name = ti.name;
        }

        res.push(tagInfo);
    }

    res[res.length - 1].isEnd = true;
    return res;
};

/**
 * 获取所有标签路径信息，如 default/test1/1|machinecode -> ['default/', 'default/test1/', 'default/test1/1|machinecode']
 * @param codePath 标签路径
 * @returns 所有层级路径数组
 */
function getAllCodePath(codePath: string) {
    if (!codePath) return [];

    const parts = codePath.split('/');
    const result: string[] = [];
    let currentPath = '';

    for (const part of parts) {
        if (!part) {
            continue;
        }
        currentPath += part + '/';
        result.push(currentPath);
    }

    return result;
}

/**
 * 完善标签路径信息
 * @param codePaths 标签路径
 * @returns
 */
async function getAllCodePaths(codePaths: string[]) {
    if (!codePaths) return;
    const allCodePaths: string[] = [];

    // 收集所有层级路径并去重
    for (const codePath of codePaths) {
        allCodePaths.push(...getAllCodePath(codePath));
    }

    const codepath2CodeInfo: any = {};
    // 去重
    const uniqueCodePaths = [...new Set(allCodePaths)];
    if (uniqueCodePaths.length == 0) {
        return codepath2CodeInfo;
    }

    const tagInfos = await tagApi.listByQuery.request({ tagPaths: uniqueCodePaths });

    for (const tagInfo of tagInfos) {
        codepath2CodeInfo[tagInfo.codePath] = tagInfo;
    }

    return codepath2CodeInfo;
}
</script>
<style lang="scss" scoped></style>
