<script setup lang="ts">

import { ref } from 'vue';
import { ProxyDetails } from '../../types/dbTypes';
import { getUserId } from '../../utils/local-storage';

const props = defineProps<{
	proxy: ProxyDetails,
	onSave: (p: ProxyDetails) => void,
	onCancel: () => void
}>();

const proto = ref(props.proxy.proto);
const host = ref(props.proxy.host);
const port = ref(props.proxy.port);
const name = ref(props.proxy.name);
const password = ref(props.proxy.password);
const errMsg = ref("");

function onSave() {
	const p = new ProxyDetails;

	const values: any = {
		id: props.proxy.id,
		userId: props.proxy.userId,
		proto: proto.value.trim(),
		host: host.value.trim(),
		port: port.value.trim(),
		name: name.value.trim(),
		password: password.value.trim(),
		isInUse: props.proxy.isInUse
	}

	const mandatory = ["proto", "host", "port"];

	for(const m of mandatory) {
		if(values[m] === "") {
			errMsg.value = `${m} is mandatory.`;
			return;
		}
	}

	if(Number.isNaN(Number(values.port))) {
		errMsg.value = `port should be a number.`;
		return;
	}

	for(const k in values) {
		(p as any)[k] = values[k];
	}

	p.userId = getUserId();

	props.onSave(p);
}

</script>

<template>
	<div class="edit-proxy-container">
		<div class="edit-proxy-input-container">
			<div>
				<input type="text" v-model="proto" placeholder="proto" />
			</div>
			<div>
				<input type="text" v-model="host" placeholder="host" />
			</div>
			<div>
				<input type="text" v-model="port" placeholder="port" />
			</div>
			<div>
				<input type="text" v-model="name" placeholder="name" />
			</div>
			<div>
				<input type="text" v-model="password" placeholder="password" />
			</div>
		</div>

		<div v-if="errMsg !== ''" class="err-msg">
			{{ errMsg }}
		</div>

		<div class="btn-container">
			<button class="btn btn-nobg" @click="props.onCancel">cancel</button>
			<button class="btn btn-nobg clr-green" @click="onSave">save</button>
		</div>
	</div>
</template>

<style scoped>
@import "@css/common.css";

.edit-proxy-container{
	padding: 0.5rem;
	background-color: rgb(80, 80, 80);
	border-radius: 6px;
}

.edit-proxy-input-container{
	display: flex;
	flex-direction: row;
	justify-content: space-around;
	margin-bottom: 1rem;
}

.edit-proxy-input-container div{
	width: 100%;
	margin: 0 1rem;
}

.edit-proxy-input-container input{
	width: 100%;
}

.edit-proxy-input-container input::placeholder{
	color: rgb(140, 140, 140);
}

.err-msg{
	margin-bottom: 1rem;
	margin-left: 1rem;
	color: #ff7575;
}

.btn-container{
	margin-left: 1rem;
}


</style>
