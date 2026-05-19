<template>
    <div class="kafka-consumer-group h-full card !p-1">
        <div class="toolbar flex items-center justify-between mb-2">
            <div class="flex items-center">
                <el-input v-model="searchGroup" :placeholder="$t('mq.kafka.searchGroup')" clearable size="small" class="w-60" @clear="loadGroups" />
                <el-button @click="loadGroups" icon="refresh" :loading="loading" size="small" plain class="ml-2">
                    {{ $t('common.refresh') }}
                </el-button>
            </div>
            <div class="flex items-center">
                <span class="text-sm text-gray-500 mr-2">{{ $t('count') + ` ${groups.length}` }}</span>
            </div>
        </div>

        <el-table :data="filteredGroups" stripe style="width: 100%" v-loading="loading">
            <el-table-column prop="Group" :label="$t('mq.kafka.groupId')" min-width="250" />
            <el-table-column prop="Coordinator" :label="$t('mq.kafka.coordinator')" min-width="150" />
            <el-table-column prop="State" :label="$t('mq.kafka.state')" min-width="120">
                <template #default="{ row }">
                    <el-tag :type="getStateTagType(row.State)" size="small">{{ row.State }}</el-tag>
                </template>
            </el-table-column>
            <el-table-column prop="ProtocolType" :label="$t('mq.kafka.protocolType')" min-width="150" />
            <el-table-column :label="$t('common.operation')" width="150" fixed="right">
                <template #default="{ row }">
                    <el-button @click="handleGetGroupMembers(row)" size="small" icon="setting" link>
                        {{ $t('mq.kafka.Members') }}
                    </el-button>
                    <el-button @click="handleDeleteGroup(row)" type="danger" size="small" icon="delete" link v-auth="'kafka:group:delete'">
                        {{ $t('common.delete') }}
                    </el-button>
                </template>
            </el-table-column>
        </el-table>
    </div>
</template>

<script lang="ts" setup>
import { Msg, useI18nDeleteConfirm } from '@/hooks/useI18n';
import { computed, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { mqApi } from '../../api';

export interface ConsumerGroup {
    Coordinator: number;
    State: string;
    ProtocolType: string;
    Group: string;
}

const { t } = useI18n();

const props = defineProps({
    kafkaId: {
        type: Number,
        required: true,
    },
    groups: {
        type: Array as () => ConsumerGroup[],
        default: () => [],
    },
    loading: {
        type: Boolean,
        default: false,
    },
});

const emits = defineEmits(['refresh']);

const searchGroup = ref('');

const filteredGroups = computed(() => {
    if (!searchGroup.value) {
        return props.groups;
    }
    return props.groups.filter((group: ConsumerGroup) => group.Group.toLowerCase().includes(searchGroup.value.toLowerCase()));
});

const loadGroups = () => {
    emits('refresh');
};

const handleDeleteGroup = async (group: ConsumerGroup) => {
    await useI18nDeleteConfirm(`Group: ${group.Group}`);
    try {
        await mqApi.kafkaDeleteGroup.request({
            id: props.kafkaId,
            group: group.Group,
        });
        Msg.saveSuccess();
        emits('refresh');
    } catch (error: any) {
        Msg.error(error.message || 'common.requestFail');
    }
};
const handleGetGroupMembers = async (group: ConsumerGroup) => {
    try {
        let res = await mqApi.kafkaGetGroupMembers.request({
            id: props.kafkaId,
            group: group.Group,
        });
        console.log(res);
    } catch (error: any) {
        Msg.error(error.message || 'common.requestFail');
    }
};

const getStateTagType = (state: string) => {
    switch (state?.toLowerCase()) {
        case 'stable':
            return 'success';
        default:
            return '';
    }
};
</script>

<style lang="scss" scoped>
.kafka-consumer-group {
    .toolbar {
        display: flex;
        justify-content: space-between;
        align-items: center;
    }
}
</style>
