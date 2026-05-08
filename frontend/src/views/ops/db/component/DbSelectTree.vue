<template>
    <ResourceSelect
        v-bind="$attrs"
        v-model="selectNode"
        @change="changeNode"
        :resource-type="ResourceTypeEnum.Db.value"
        :leaf-node-types="[NodeTypePostgresSchema]"
        :transform-node="transformNode"
    >
        <template #iconPrefix>
            <SvgIcon v-if="dbType && getDbDialect(dbType)" :name="getDbDialect(dbType).getInfo().icon" :size="16" />
            <TagCodePath :code="dbCode" />
        </template>
    </ResourceSelect>
</template>

<script setup lang="ts">
import { ResourceTypeEnum } from '@/common/commonEnum';
import { TagTreeNode } from '@/views/ops/component/tag';
import { getDbDialect, schemaDbTypes } from '@/views/ops/db/dialect';
import { NodeTypeDb, NodeTypePostgresSchema } from '@/views/ops/db/resource';
import ResourceSelect from '@/views/ops/resource/ResourceSelect.vue';
import { computed, ref, watch } from 'vue';
import TagCodePath from '../../component/TagCodePath.vue';
import { dbApi } from '@/views/ops/db/api';

const dbId = defineModel<number>('dbId');
const instName = defineModel<string>('instName');
const dbName = defineModel<string>('dbName');
const tagPath = defineModel<string>('tagPath');
const dbType = defineModel<string>('dbType');

const dbCode = ref('');

const emits = defineEmits(['selectDb']);

const selectNode = computed({
    get: () => {
        return dbName.value;
    },
    set: () => {
        //
    },
});

watch(
    () => dbId.value,
    async (id) => {
        if (!id || id <= 0) {
            return;
        }

        const dbRes = await dbApi.dbs.request({ id: dbId.value });
        const db = dbRes.list?.[0];
        if (!db) {
            console.log('not found db: {}', id);
            return '';
        }
        dbCode.value = db.code;
    },
    { immediate: true }
);

// 节点转换函数：动态判断数据库节点是否为叶子节点
const transformNode = (node: TagTreeNode): TagTreeNode => {
    // 如果是数据库节点，根据数据库类型动态设置 isLeaf
    if (node.type.value === NodeTypeDb.value) {
        const hasSchema = schemaDbTypes.includes(node.params?.type);
        // 没有 schema 的数据库（如 MySQL），标记为叶子节点
        if (!hasSchema) {
            node.isLeaf = true;
        }
    }
    return node;
};

const changeNode = (nodeData: TagTreeNode) => {
    const params = nodeData.params;
    dbName.value = params.db;
    instName.value = params.name;
    dbId.value = params.id;
    tagPath.value = params.tagPath;
    dbType.value = params.type;

    emits('selectDb', params);
};
</script>

<style lang="scss"></style>
