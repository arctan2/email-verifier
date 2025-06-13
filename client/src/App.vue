<script setup lang="ts">
import InfoErrorPopup from './common-components/info-error-popup.vue';
import { popupData } from './utils/popup-data';

</script>

<template>
	<div id="app-body">
		<router-view></router-view>
		<InfoErrorPopup v-if="popupData.type !== ''" />
	</div>
</template>

<style scoped>
#app-body{
	width: 100%;
	height: 100%;
	background-color: rgb(235, 235, 235);
	position: relative;
}
</style>
