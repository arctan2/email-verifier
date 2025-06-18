<script setup lang="ts">
import { type EmailDetails } from '../../types/dbTypes';
import { computed, onMounted, ref } from 'vue';

import Pagination from '../../common-components/pagination.vue';
import { fetchGet, fetchPost } from '../../utils/fetch';
import { setPopupError } from '../../utils/popup-data';
import { emailListProps } from "../../state/email-list-route-state";
import EmailList from "./email-list.vue";
import { filterers, getToFilter, resetFilters } from '../../utils/email-filter-state';

const props = emailListProps.value;

const totalEmailCount = ref(props.totalEmailCount);
const emailList = ref<EmailDetails[]>([]);
const showLoadingMessage = ref<string>("");
const limit = ref<number>(500);
const pagesCount = computed(() => Math.ceil(totalEmailCount.value / limit.value))
let prevFrom = 0;
let prevFilterFrom = 0;

async function fetchRecords(from: number) {
	if(prevFilterFrom !== 0) {
		from = 0;
		prevFilterFrom = 0;
	}
	prevFrom = from;
	emailList.value = [];
	showLoadingMessage.value = "Loading...";

	const res = await fetchGet<{emailDetailsList: EmailDetails[]}>(
		`/${props.fileId}/get-email-details-list?from=${from * limit.value}&limit=${limit.value}`
	);

	showLoadingMessage.value = "";
	if(res.err) {
		setPopupError(res.msg)
		return;
	}

	emailList.value = res.emailDetailsList;
	totalEmailCount.value = props.totalEmailCount;
}

async function refresh() {
	showLoadingMessage.value = "Refreshing...";
	await Promise.all([
		fetchRecords(0),
	])
	showLoadingMessage.value = "";
}

async function onFilterChange(toFilterList: string[], from: number) {
	if(toFilterList.length === 0) {
		fetchRecords(prevFrom);
		return;
	}

	if(prevFrom !== 0) {
		from = 0;
		prevFrom = 0;
	}
	
	let toFilter: any = {};

	for(const f of toFilterList) {
		toFilter[f] = filterers[f].cmpVal;
	}

	emailList.value = [];
	showLoadingMessage.value = "Loading...";

	const body = {
		from: from * limit.value,
		limit: limit.value,
		filterFields: toFilter,
		fileId: props.fileId
	}

	const res = await fetchPost<{emailDetailsList: EmailDetails[], totalEmailCount: number}>("/filter-emails", body);

	showLoadingMessage.value = "";
	if(res.err) {
		setPopupError(res.msg)
		return;
	}

	emailList.value = res.emailDetailsList;
	totalEmailCount.value = res.totalEmailCount;
}

function handlePageChange(newPageIdx: number) {
	const toFilterList = getToFilter();

	if(toFilterList.length) {
		onFilterChange(toFilterList, newPageIdx);
	} else {
		fetchRecords(newPageIdx);
	}
}

onMounted(() => {
	resetFilters();
	refresh();
})

</script>

<template>
<div id="email-details">
	<div class="email-details-container">
		<EmailList :emailList="emailList" :filterers="filterers" @filter-change="() => onFilterChange(getToFilter(), prevFilterFrom)" />
		<Pagination :length="pagesCount" :to-show="5" :on-change="handlePageChange" />
	</div>
	<div class="refreshing" v-if="showLoadingMessage !== ''">{{ showLoadingMessage }}</div>
</div>
</template>

<style scoped>
@import "@css/common.css";

#email-details{
	color: white;
	max-width: 100rem;
	height: 100%;
}

h1{
	height: 2rem;
	margin-bottom: 0.5rem;
}

.email-details-container{
	width: 100%;
	height: calc(100% - 3.5rem);
	background-color: rgb(50, 50, 50);
	border-radius: 6px;
	position: relative;
}

.icon-btns{
	display: flex;
	flex-direction: row;
	align-items: center;
	position: absolute;
	top: 0.5rem;
	right: 0.5rem;
}

.icon-btns > div {
	--size: 2rem;
	width: var(--size);
	height: var(--size);
	margin: 0 0.3rem;
	padding: 0.2rem;
	border-radius: 6px;
	cursor: pointer;
	user-select: none;
}

.icon-btns > div:hover{
	background-color: rgb(100, 100, 100);
}

.icon-btns > div img {
	width: 100%;
	height: 100%;
	object-fit: contain;
	pointer-events: none;
}

.refreshing{
	position: absolute;
	display: flex;
	align-items: center;
	justify-content: center;
	top: 0;
	left: 0;
	width: 100%;
	height: 100%;
	background-color: rgba(0, 0, 0, 0.7);
	border-radius: 6px;
}

@media (max-width: 55rem){
.email-details-container{
	height: calc(100dvh - 20dvh - 3.5rem);
}
}
</style>
