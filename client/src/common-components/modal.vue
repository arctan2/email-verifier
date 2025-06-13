<script setup lang="ts">
const props = defineProps<{
	dismiss: () => void,
}>();

</script>

<template>
<div class="modal-container" @click.self="props.dismiss">
	<slot></slot>
</div>
</template>

<style>
@import "@css/common.css";

.modal-container{
	display: flex;
	align-items: center;
	justify-content: center;
	position: fixed;
	width: 100%;
	height: 100%;
	top: 0;
	left: 0;
	z-index: 100;
	background-color: rgba(0, 0, 0, 0.5);
}

</style>
