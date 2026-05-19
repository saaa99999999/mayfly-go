<template>
    <div class="h-full">
        <page-table
            ref="pageTableRef"
            :page-api="mqApi.kafkaList"
            :before-query-fn="checkRouteTagPath"
            :search-items="searchItems"
            v-model:query-form="query"
            :show-selection="true"
            v-model:selection-data="selectionData"
            :columns="columns"
            lazy
        >
            <template #tableHeader>
                <el-button v-auth="'mq:kafka:save'" type="primary" icon="plus" @click="editKafka(false)" plain>{{ $t('common.create') }}</el-button>
                <el-button v-auth="'mq:kafka:del'" type="danger" icon="delete" :disabled="selectionData.length < 1" @click="deleteKafka" plain>
                    {{ $t('common.delete') }}
                </el-button>
            </template>

            <template #name="{ data }">
                <TagCodePath :code="data.code" show-popover />
                {{ data.name }}
            </template>

            <template #action="{ data }">
                <el-button v-auth="'mq:kafka:save'" @click="editKafka(data)" link type="primary">{{ $t('common.edit') }}</el-button>
            </template>
        </page-table>

        <kafka-edit @val-change="search()" :title="kafkaEditDialog.title" v-model:visible="kafkaEditDialog.visible" v-model:kafka="kafkaEditDialog.data" />
    </div>
</template>

<script lang="ts" setup>
import { TagResourceTypeEnum } from '@/common/commonEnum';
import { TableColumn } from '@/components/pagetable';
import PageTable from '@/components/pagetable/PageTable.vue';
import { SearchItem } from '@/components/pagetable/SearchForm';
import { Msg, useI18nCreateTitle, useI18nDeleteConfirm, useI18nEditTitle } from '@/hooks/useI18n';
import { getTagPathSearchItem } from '@/views/ops/component/tag';
import { mqApi } from '@/views/ops/mq/api';
import { defineAsyncComponent, onMounted, reactive, ref, Ref, toRefs } from 'vue';
import { useRoute } from 'vue-router';
import TagCodePath from '../../component/TagCodePath.vue';

const KafkaEdit = defineAsyncComponent(() => import('./KafkaEdit.vue'));

const props = defineProps({
    lazy: {
        type: [Boolean],
        default: false,
    },
});

const route = useRoute();
const pageTableRef: Ref<any> = ref(null);

const searchItems = [
    SearchItem.input('keyword', 'common.keyword').withPlaceholder('mq.kafka.keywordPlaceholder'),
    getTagPathSearchItem(TagResourceTypeEnum.MqKafka.value),
];

const columns = [
    TableColumn.new('name', 'common.name').isSlot('name').setAddWidth(15),
    TableColumn.new('hosts', 'Hosts'),
    TableColumn.new('username', 'mq.kafka.username'),
    TableColumn.new('password', 'common.password'),
    TableColumn.new('createTime', 'common.createTime').isTime(),
    TableColumn.new('creator', 'common.creator'),
    TableColumn.new('code', 'common.code'),
    TableColumn.new('action', 'common.operation').isSlot().setMinWidth(170).fixedRight().alignCenter(),
];

const state = reactive({
    dbOps: {
        dbId: 0,
        db: '',
    },
    selectionData: [],
    query: {
        pageNum: 1,
        pageSize: 0,
        tagPath: '',
    },
    kafkaEditDialog: {
        visible: false,
        data: null as any,
        title: '',
    },
});

const { selectionData, query, kafkaEditDialog } = toRefs(state);

const checkRouteTagPath = (query: any) => {
    if (route.query.tagPath) {
        query.tagPath = route.query.tagPath as string;
    }
    return query;
};

const deleteKafka = async () => {
    try {
        await useI18nDeleteConfirm(state.selectionData.map((x: any) => x.name).join('、'));
        await mqApi.kafkaDel.request({ id: state.selectionData.map((x: any) => x.id).join(',') });
        Msg.deleteSuccess();
        search();
    } catch (err) {
        //
    }
};

const search = async (tagPath: string = '') => {
    if (tagPath) {
        state.query.tagPath = tagPath;
    }
    pageTableRef.value.search();
};

const editKafka = async (data: any) => {
    if (!data) {
        state.kafkaEditDialog.data = null;
        state.kafkaEditDialog.title = useI18nCreateTitle('Kafka');
    } else {
        state.kafkaEditDialog.data = data;
        state.kafkaEditDialog.title = useI18nEditTitle('Kafka');
    }
    state.kafkaEditDialog.visible = true;
};

onMounted(() => {
    if (!props.lazy) {
        search();
    }
});

defineExpose({ search });
</script>

<style></style>
