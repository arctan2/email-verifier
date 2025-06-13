<script setup lang="ts">
import { type EmailDetails } from '../types/dbTypes';
import { computed, onMounted, ref } from 'vue';

import Pagination from '../common-components/pagination.vue';
import { fetchGet } from '../utils/fetch';
import { setPopupError } from '../utils/popup-data';

const props = defineProps<{
	totalEmailCount: number,
	fileId: number
}>();

const emailList = ref<EmailDetails[]>([]);
const showLoadingMessage = ref<string>("");
const limit = ref<number>(500);
const pagesCount = computed(() => Math.ceil(props.totalEmailCount / limit.value))

async function fetchRecords(from: number) {
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
}

async function refresh() {
	showLoadingMessage.value = "Refreshing...";
	await Promise.all([
		fetchRecords(0),
	])
	showLoadingMessage.value = "";
}

onMounted(() => {
	refresh();
})

</script>

<template>
<div id="email-details">
	<div class="main-section">
		<div class="email-details-container">
			<div class="email-list scroll-bar">
				<div class="details-title-container">
					<div class="details-title">
						<div>syntax</div>
						<div>reachable</div>
						<div>deliverable</div>
						<div>host-exists</div>
						<div>mx-records</div>
						<div>disposable</div>
						<div>catch-all</div>
						<div>inbox-full</div>
					</div>
				</div>
				<div v-for="email, idx in emailList" :key="idx" class="email">
					<div>
						{{ email.emailId }}
					</div>
					<div class="state">
						<div :class="[email.isValidSyntax ? 'green' : 'red']"></div>
						<div :class="[email.isValidSyntax ? 'green' : 'red']"></div>
						<div :class="[email.isValidSyntax ? 'green' : 'red']"></div>
						<div :class="[email.isValidSyntax ? 'green' : 'red']"></div>
						<div :class="[email.isValidSyntax ? 'green' : 'red']"></div>
						<div :class="[email.isValidSyntax ? 'green' : 'red']"></div>
						<div :class="[email.isValidSyntax ? 'green' : 'red']"></div>
						<div :class="[email.isValidSyntax ? 'green' : 'red']"></div>
					</div>
				</div>
			</div>
			<Pagination :length="pagesCount" :to-show="5" :on-change="fetchRecords"/>
		</div>
		<div class="refreshing" v-if="showLoadingMessage !== ''">{{ showLoadingMessage }}</div>
	</div>
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

.main-section {
	height: calc(100dvh - 2rem);
	overflow: hidden;
}

.email-details-container{
	width: 100%;
	height: calc(100dvh - 16dvh - 1rem);
	background-color: rgb(50, 50, 50);
	border-radius: 6px;
	position: relative;
}

.email-list{
	display: flex;
	flex-direction: column;
	position: relative;
	overflow-y: auto;
	height: 92%;
	margin-bottom: 1rem;
	padding: 0.5rem;
	background-color: rgb(40, 40, 40);
	border-radius: 6px;
	width: 100%;
}

.email{
	display: flex;
	flex-direction: row;
	justify-content: space-between;
	width: 100%;
	position: relative;
	border-radius: 6px;
	margin-bottom: 0.5rem;
	background-color: rgb(60, 60, 60);
	padding: 0.5rem;
	font-family: monospace;
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

.details-title-container{
	display: flex;
	flex-direction: column;
	align-items: flex-end;
	width: 100%;
	position: sticky;
	top: -1rem;
	z-index: 100;
	margin-bottom: 1rem;
	background-color: rgb(40, 40, 40);
	padding: 0.5rem;
	padding-top: 1rem;
}

.details-title, .state{
	display: flex;
	flex-direction: row;
	justify-content: space-around;
	width: 40%;
	max-width: 20rem;
}

.details-title div {
	writing-mode: vertical-rl;
	transform: rotate(180deg);
	text-orientation: sideways-right;
}

.state div {
	width: 0.8rem;
	height: 0.8rem;
}

.state .green {
	background-color: rgb(0, 255, 0);
}

.state .red {
	background-color: rgb(255, 0, 0);
}

@media (max-width: 55rem){
.email-details-container{
	height: calc(100dvh - 20dvh - 3.5rem);
}
}
</style>
