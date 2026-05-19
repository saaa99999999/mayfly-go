<template>
    <div class="kafka-consume-message h-full card !p-1">
        <el-form ref="consumeFormRef" :model="form" label-width="auto" size="small">
            <el-row :gutter="10">
                <el-col :span="10">
                    <el-form-item :label="$t('mq.kafka.selectTopic')" required>
                        <el-select v-model="form.topic" filterable :placeholder="$t('mq.kafka.selectTopicPlaceholder')" clearable>
                            <el-option v-for="topic in topics" :key="topic" :label="topic" :value="topic" />
                        </el-select>
                    </el-form-item>
                </el-col>
                <el-col :span="5">
                    <el-form-item :label="$t('mq.kafka.messageNumber')" required>
                        <el-input-number v-model="form.number" :min="1" :max="1000" />
                    </el-form-item>
                </el-col>
                <el-col :span="9">
                    <el-form-item :label="$t('mq.kafka.consumerGroup')">
                        <el-select v-model="form.group" filterable :placeholder="$t('mq.kafka.consumerGroupPlaceholder')" clearable allow-create>
                            <el-option label="(auto generate)" value="" />
                            <el-option v-for="g in groups" :key="g.Group" :label="g.Group" :value="g.Group" />
                        </el-select>
                    </el-form-item>
                </el-col>
            </el-row>

            <el-row :gutter="10">
                <el-col :span="5">
                    <el-form-item :label="$t('mq.kafka.pullTimeout')">
                        <el-input-number v-model="form.pullTimeout" :min="1" :max="100" />
                    </el-form-item>
                </el-col>
                <el-col :span="5">
                    <el-form-item :label="$t('mq.kafka.decompression')">
                        <el-select v-model="form.decompression" :placeholder="$t('mq.kafka.decompressionPlaceholder')" clearable>
                            <el-option label="none" value="" />
                            <el-option label="gzip" value="gzip" />
                            <el-option label="lz4" value="lz4" />
                            <el-option label="zstd" value="zstd" />
                            <el-option label="snappy" value="snappy" />
                        </el-select>
                    </el-form-item>
                </el-col>
                <el-col :span="5">
                    <el-form-item :label="$t('mq.kafka.decode')">
                        <el-select v-model="form.decode" :placeholder="$t('mq.kafka.decodePlaceholder')" clearable>
                            <el-option label="None" value="" />
                            <el-option label="Base64" value="base64" />
                        </el-select>
                    </el-form-item>
                </el-col>
                <el-col :span="5">
                    <el-form-item :label="$t('mq.kafka.isolationLevel')">
                        <el-select v-model="form.isolationLevel" :placeholder="$t('mq.kafka.isolationLevelPlaceholder')">
                            <el-option :label="$t('mq.kafka.readUncommitted')" value="read_uncommitted" />
                            <el-option :label="$t('mq.kafka.readCommitted')" value="read_committed" />
                        </el-select>
                    </el-form-item>
                </el-col>
            </el-row>

            <el-row :gutter="10">
                <el-col :span="5">
                    <el-form-item :label="$t('mq.kafka.commitOffset')">
                        <el-switch v-model="form.commitOffset" />
                    </el-form-item>
                </el-col>
                <el-col :span="7">
                    <el-tooltip :content="$t('mq.kafka.consumerOnlyTip')">
                        <el-form-item :label="$t('mq.kafka.defaultConsumePosition')">
                            <el-switch v-model="form.earliest" :active-text="$t('mq.kafka.earliest')" :inactive-text="$t('mq.kafka.latest')" />
                        </el-form-item>
                    </el-tooltip>
                </el-col>
                <el-col :span="8">
                    <el-tooltip :content="$t('mq.kafka.consumerOnlyTip')">
                        <el-form-item :label="$t('mq.kafka.defaultConsumeStartTime')">
                            <el-date-picker
                                v-model="form.startTime"
                                type="datetime"
                                :placeholder="$t('mq.kafka.selectDateTime')"
                                value-format="YYYY-MM-DD HH:mm:ss"
                                size="small"
                            />
                        </el-form-item>
                    </el-tooltip>
                </el-col>
            </el-row>

            <el-form-item>
                <el-button @click="resetForm" icon="refresh">{{ $t('common.reset') }}</el-button>
                <el-button @click="consumeMessage" type="primary" icon="download" :loading="consuming">
                    {{ $t('mq.kafka.consumeMessage') }}
                </el-button>
            </el-form-item>
        </el-form>

        <el-table :data="messages" stripe style="width: 100%" v-loading="consuming" max-height="700">
            <el-table-column prop="offset" :label="$t('mq.kafka.offset')" min-width="100" />
            <el-table-column prop="partition" :label="$t('mq.kafka.partition')" min-width="80" />
            <el-table-column prop="key" :label="$t('mq.kafka.key')" min-width="150" />
            <el-table-column prop="timestamp" :label="$t('mq.kafka.timestamp')" min-width="180" />
            <el-table-column prop="value" :label="$t('mq.kafka.messageBody')" min-width="300">
                <template #default="{ row }">
                    <div class="flex items-center">
                        <el-input v-model="row.displayValue" type="textarea" :rows="1" size="small" class="flex-1" />
                        <SvgIcon
                            v-if="row.value && row.value.length > 50"
                            @click="viewMessageDetail(row)"
                            class="string-input-container-icon ml-1 cursor-pointer"
                            name="FullScreen"
                            :size="10"
                        />
                    </div>
                </template>
            </el-table-column>
            <el-table-column prop="headers" :label="$t('mq.kafka.headers')" min-width="150">
                <template #default="{ row }">
                    {{ JSON.stringify(row.headers) }}
                </template>
            </el-table-column>
        </el-table>
    </div>
</template>

<script lang="ts" setup>
import { ref, reactive, toRefs, onMounted, defineAsyncComponent, watch } from 'vue';
import { mqApi } from '../../api';
import { ElMessage } from 'element-plus';
import { useI18n } from 'vue-i18n';
import SvgIcon from '@/components/svgIcon/index.vue';
import MonacoEditorBox from '@/components/monaco/MonacoEditorBox';
import { ConsumerGroup } from '@/views/ops/mq/kafka/component/ConsumerGroup.vue';
import { randomUuid } from '@/common/utils/string';
import { Msg } from '@/hooks/useI18n';

const { t } = useI18n();

const props = defineProps({
    kafkaId: {
        type: Number,
        required: true,
    },
    defaultTopic: {
        type: String,
        default: '',
    },
    topics: {
        type: Array as () => string[],
        default: () => [],
    },
    groups: {
        type: Array as () => ConsumerGroup[],
        default: () => [],
    },
});

const consumeFormRef = ref();
const consuming = ref(false);

const state = reactive({
    form: {
        topic: '',
        number: 10,
        group: '',
        pullTimeout: 10,
        decompression: '',
        decode: '',
        isolationLevel: 'read_uncommitted',
        commitOffset: false,
        earliest: true,
        startTime: '',
    },
    messages: [] as any[],
});

const { form, messages } = toRefs(state);

onMounted(() => {
    if (props.defaultTopic) {
        state.form.topic = props.defaultTopic;
    }
});

watch(
    () => props.defaultTopic,
    (newTopic) => {
        state.form.topic = newTopic || '';
    }
);

const resetForm = () => {
    state.form = {
        topic: props.defaultTopic || '',
        number: 10,
        group: '',
        pullTimeout: 10,
        decompression: '',
        decode: '',
        isolationLevel: 'read_uncommitted',
        commitOffset: false,
        earliest: true,
        startTime: '',
    };
    state.messages = [];
};

const consumeMessage = async () => {
    if (!consumeFormRef.value) return;
    consuming.value = true;
    if (!state.form.group) {
        state.form.group = '__mayfly-server__' + randomUuid();
    }

    try {
        const param = {
            id: props.kafkaId,
            ...state.form,
        };

        const res = await mqApi.kafkaTopicConsume.request(param);
        state.messages = (res || []).map((msg: any, index: number) => ({
            ...msg,
            displayValue: typeof msg.value === 'object' ? JSON.stringify(msg.value, null, 2) : String(msg.value),
        }));
    } catch (error: any) {
        Msg.error(error.message || 'common.requestFail');
    } finally {
        consuming.value = false;
    }
};

const viewMessageDetail = (row: any) => {
    const value = typeof row.value === 'object' ? JSON.stringify(row.value, null, 2) : String(row.value);
    const editorLang = getEditorLangByValue(value);

    MonacoEditorBox({
        content: value,
        title: `${t('mq.kafka.messageBody')} - Offset ${row.offset}`,
        language: editorLang,
        showConfirmButton: false,
        closeFn: () => {},
    });
};

const getEditorLangByValue = (value: any) => {
    try {
        if (typeof JSON.parse(value) === 'object') {
            return 'json';
        }
    } catch (e) {
        /* empty */
    }

    try {
        const doc = new DOMParser().parseFromString(value, 'text/html');
        if (Array.from(doc.body.childNodes).some((node) => node.nodeType === 1)) {
            return 'html';
        }
    } catch (e) {
        /* empty */
    }

    return 'text';
};
</script>

<style lang="scss" scoped>
.kafka-consume-message {
    overflow: auto;

    .string-input-container-icon {
        color: var(--el-color-primary);
        &:hover {
            color: var(--el-color-success);
        }
    }
}
</style>
