<script setup lang="ts">
import { ref } from 'vue';
import { fetchPost, fetchGet, fetchDelete, fetchPut } from '../../utils/fetch';
import { getUserId } from '../../utils/local-storage';
import { ProxyDetails } from "../../types/dbTypes";
import { setPopupError, setPopupInfo } from '../../utils/popup-data';
import EditProxy from "./edit-proxy.vue";
import DeleteIcon from "../../assets/icons/delete.svg";
import EditIcon from "../../assets/icons/edit.svg";
import Modal from '../../common-components/modal.vue';

enum Msg {
	None = ""
}

const deleteIdx = ref<number>(-1);
const editIdx = ref<number>(-1);
const showLoadingMessage = ref<string>("");

const msg = ref<string>(Msg.None);
const list = ref<ProxyDetails[]>([]);

async function fetchProxyList() {
	const data = await fetchGet<{proxyList: ProxyDetails[]}>(`/get-proxy-list?userId=${getUserId()}`);

	if(data.err) {
		msg.value = data.msg;
	} else {
		list.value = data.proxyList;
		msg.value = Msg.None;
	}
}

async function refresh() {
	showLoadingMessage.value = "Refreshing...";
	await Promise.all([
		fetchProxyList()
	])
	showLoadingMessage.value = "";
}

async function deleteProxy() {
	showLoadingMessage.value = "Deleting...";
	const toDeleteId = list.value[deleteIdx.value].id;
	deleteIdx.value = -1;
	const res = await fetchDelete<{keywords: string[]}>(`/${toDeleteId}/delete-proxy?userId=${getUserId()}`);
	if(res.err) {
		setPopupError(res.msg);
	} else {
		setPopupInfo("Proxy deleted successfully.");
	}
	showLoadingMessage.value = "";
	refresh();
}

async function insertProxy(proxy: ProxyDetails) {
	proxy.userId = getUserId();
	showLoadingMessage.value = "Adding...";
	const res = await fetchPost(`/insert-proxy`, proxy);
	if(res.err) {
		setPopupError(res.msg);
	} else {
		setPopupInfo("Proxy inserted successfully.");
	}
	showLoadingMessage.value = "";
	refresh();
}

async function updateProxy(proxy: ProxyDetails) {
	showLoadingMessage.value = "Updating...";
	const res = await fetchPut(`/${proxy.id}/update-proxy`, proxy);
	if(res.err) {
		setPopupError(res.msg);
	} else {
		setPopupInfo("Proxy updated successfully.");
	}
	showLoadingMessage.value = "";
	refresh();
}

function saveProxy(proxy: ProxyDetails) {
	if(proxy.id === -1) {
		insertProxy(proxy);
	} else {
		updateProxy(proxy);
	}
	editIdx.value = -1;
}

async function toggleEnable(idx: number) {
	showLoadingMessage.value = "Updating...";
	const proxy = list.value[idx];
	const res = await fetchPut(`/${proxy.id}/update-proxy-is-enabled?isEnabled=${!proxy.isEnabled}&userId=${getUserId()}`);
	if(res.err) {
		setPopupError(res.msg);
	}
	showLoadingMessage.value = "";
	refresh();
}

fetchProxyList();

</script>

<template>
	<div id="proxy-list-container">
		<div class="scroll-bar">
			<Modal v-if="deleteIdx >= 0" :dismiss="() => deleteIdx = -1">
				<div class="modal delete-modal">
					<h2 style="text-align: center;">Delete</h2>

					<div>Are you sure you want to delete the proxy?</div>
								
					<div class="btn-container" style="align-self: center;">
						<button class="btn btn-nobg clr-red" @click="deleteProxy">delete</button>
						<button class="btn btn-nobg clr-white" @click="deleteIdx = -1">cancel</button>
					</div>
				</div>
			</Modal>

			<div v-if="msg !== Msg.None" class="msg">{{ msg }}</div>
			<div class="proxy-list" v-else>
				<div class="headers-container">
					<div class="titles">
						<div>proto</div>
						<div>host</div>
						<div>port</div>
						<div>name</div>
						<div>password</div>
					</div>
					<div class="pad"></div>
				</div>

				<div :class="['proxy-item', proxy.isEnabled ? '' : 'disabled']" v-for="proxy, idx in list">
					<EditProxy
						v-if="editIdx === idx"
						:proxy="{...proxy}"
						:onSave="saveProxy"
						:onCancel="() => editIdx = -1"
					/>

					<template v-else>
						<div class="details">
							<div>{{ proxy.proto }}</div>
							<div>{{ proxy.host }}</div>
							<div>{{ proxy.port }}</div>
							<div>{{ proxy.name ? proxy.name : "---" }}</div>
							<div>{{ proxy.password ? proxy.password : "---" }}</div>
						</div>

						<div class="icon-btns" v-if="editIdx < 0">
							<div @click="editIdx = idx">
								<img :src="EditIcon" alt="">
							</div>
							<div @click="deleteIdx = idx">
								<img :src="DeleteIcon" alt="">
							</div>
							<div>
								<div class="block-icon" @click="toggleEnable(idx)"></div>
							</div>
						</div>
					</template>
				</div>

				<button v-if="editIdx === -1" class="btn btn-nobg add-btn" @click="editIdx = -2">+</button>

				<EditProxy
					v-if="editIdx === -2"
					:proxy="new ProxyDetails"
					:onSave="saveProxy"
					:onCancel="() => editIdx = -1"
				/>

			</div>
			<div class="loading" v-if="showLoadingMessage !== ''">{{ showLoadingMessage }}</div>
		</div>
	</div>
</template>

<style scoped>
@import "@css/common.css";

#proxy-list-container{
	width: 100%;
	height: 84dvh;
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	color: white;
	position: relative;
}

#proxy-list-container > div{
	display: flex;
	flex-direction: column;
	padding: 1rem;
	padding-top: 0;
	border-radius: 6px;
	height: 90%;
	width: 96%;
	background-color: rgb(60, 60, 60);
	overflow-x: auto;
	max-width: max-content;
}

.msg{
	text-align: center;
}

.proxy-list{
	max-height: 80%;
	position: relative;
	width: max-content;
}

.proxy-item{
	display: flex;
	flex-direction: row;
	position: relative;
	background-color: rgb(80, 80, 80);
	border-radius: 6px;
	margin-bottom: 0.5rem;
	padding: 0.5rem 0;
	min-width: max-content;
}

.proxy-item.disabled{
	background-color: rgba(80, 80, 80, 0.2);
	color: rgba(255, 255, 255, 0.3);
}

.proxy-item.disabled .icon-btns div {
	opacity: 0.2;
}

.proxy-item.disabled .icon-btns div:last-child {
	opacity: 1;
}

.icon-btns{
	display: flex;
	flex-direction: row;
	align-items: center;
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

.add-btn{
	min-width: 100%;
	font-size: 1.5rem;
	margin-bottom: 0.5rem;
}

.delete-modal{
	justify-content: space-between;
	background-color: rgb(45, 35, 35);
	border: 2px solid #ff6161;
}

.headers-container{
	display: flex;
	flex-direction: row;
	margin-bottom: 0.5rem;
	position: sticky;
	top: 0;
	z-index: 20;
	background-color: rgb(60, 60, 60);
	padding: 1rem 0;
}

.titles, .details{
	display: flex;
	flex-direction: row;
	align-items: center;
	justify-content: space-around;
	padding: 0.5rem;
	border-radius: 6px;
	min-width: max-content;
}

.titles{
	background-color: #6200ff;
}

.titles div, .details div{
	min-width: 10rem;
	width: 100%;
	text-align: center;
	margin: 0 1rem;
	position: relative;
	max-width: 10ch;
	text-overflow: ellipsis;
	overflow: hidden;
}

.titles div::after{
	content: "";
	display: block;
	min-width: 2px;
	height: 100%;
	background-color: #2d0075;
	position: absolute;
	right: -1rem;
	top: 0px
}

.titles div:last-child::after{
	display: none;
}

.pad{
	width: 100%;
}

.block-icon{
	--clr: rgb(255, 100, 100);
	position: relative;
	width: 100%;
	height: 100%;
	border: 4px solid var(--clr);
	border-radius: 5rem;
	display: flex;
	align-items: center;
	justify-content: center;
}

.block-icon::after{
	content: "";
	position: absolute;
	height: 4px;
	width: 100%;
	background-color: var(--clr);
	transform: rotate(-45deg);
}

</style>

