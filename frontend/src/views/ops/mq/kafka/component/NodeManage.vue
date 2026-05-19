<template>
    <div class="kafka-node-manage h-full card !p-1">
        <div class="toolbar flex items-center mb-2">
            <el-button @click="refreshBrokers" icon="refresh" :loading="loading" size="small" plain>
                {{ $t('common.refresh') }}
            </el-button>
        </div>

        <el-table :data="brokers" stripe style="width: 100%" v-loading="loading">
            <el-table-column prop="id" :label="$t('mq.kafka.nodeId')" min-width="100" />
            <el-table-column prop="addr" :label="$t('mq.kafka.addr')" min-width="100" />
            <el-table-column prop="rack" :label="$t('mq.kafka.rack')" min-width="150" />
            <el-table-column :label="$t('common.operation')" width="120" fixed="right">
                <template #default="{ row }">
                    <el-button @click="viewBrokerConfig(row)" type="primary" size="small" icon="setting" link>
                        {{ $t('mq.kafka.viewConfig') }}
                    </el-button>
                </template>
            </el-table-column>
        </el-table>

        <el-drawer
            v-model="openDrawer"
            :before-close="cancel"
            :destroy-on-close="true"
            :close-on-click-modal="false"
            size="80%"
            :title="$t('mq.kafka.brokerConfig') + selectedBroker?.addr"
        >
            <div class="toolbar">
                <div class="">
                    <el-input v-model="searchConfig" :placeholder="$t('mq.kafka.configName')" clearable size="small" class="w-60 mb-2" />
                </div>
                <span class="text-sm text-gray-500">{{ `count: ${filteredBrokerConfigs.length}` }}</span>
            </div>

            <el-table :data="filteredBrokerConfigs" stripe style="width: 100%" v-loading="loading">
                <el-table-column type="index" label="#" width="50" />
                <el-table-column prop="Key" :label="$t('mq.kafka.configName')" min-width="200" />
                <el-table-column prop="Value" :label="$t('mq.kafka.configValue')" min-width="300" />
                <el-table-column prop="Source" :label="$t('mq.kafka.configSource')" min-width="150" />
                <el-table-column prop="Sensitive" :label="$t('mq.kafka.configSensitive')" min-width="150" />
            </el-table>
        </el-drawer>
    </div>
</template>

<script lang="ts" setup>
import { Msg } from '@/hooks/useI18n';
import { computed, onMounted, reactive, ref, toRefs } from 'vue';
import { mqApi } from '../../api';

interface Broker {
    id: number;
    addr: string;
    rack: string;
}

interface BrokerConfig {
    Key: string;
    Value: string;
    Source: number;
    Sensitive: boolean;
}

const props = defineProps({
    kafkaId: {
        type: Number,
        required: true,
    },
});

const loading = ref(false);
const selectedBroker = ref<Broker | null>(null);
const openDrawer = ref(false);
const searchConfig = ref('');

const state = reactive({
    brokers: [] as Broker[],
    brokerConfigs: [] as BrokerConfig[],
});

const cancel = () => {
    state.brokerConfigs = [];
    openDrawer.value = false;
    searchConfig.value = '';
};

const { brokers, brokerConfigs } = toRefs(state);

const filteredBrokerConfigs = computed(() => {
    if (!searchConfig.value) {
        return state.brokerConfigs;
    }
    return state.brokerConfigs.filter((config: BrokerConfig) => config.Key.toLowerCase().includes(searchConfig.value.toLowerCase()));
});

onMounted(() => {
    refreshBrokers();
});

const refreshBrokers = async () => {
    loading.value = true;
    try {
        const res = await mqApi.kafkaTopicBrokers.request({ id: props.kafkaId });
        state.brokers = res || [];
    } catch (error: any) {
        Msg.error(error.message || 'common.requestFail');
    } finally {
        loading.value = false;
    }
};

const viewBrokerConfig = async (broker: Broker) => {
    selectedBroker.value = broker;
    openDrawer.value = true;
    loading.value = true;
    try {
        const res = await mqApi.kafkaTopicBrokerConfig.request({
            id: props.kafkaId,
            brokerId: broker.id,
        });

        if (res && res[broker.id].Configs) {
            res[broker.id].Configs.sort((a: any, b: any) => (a['Key'] > b['Key'] ? 1 : -1));
            state.brokerConfigs = res && res[broker.id].Configs;
        } else {
            state.brokerConfigs = [];
        }
    } catch (error: any) {
        Msg.error(error.message || 'common.requestFail');
    } finally {
        loading.value = false;
    }
};
</script>

<style lang="scss" scoped>
.toolbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
}
</style>
