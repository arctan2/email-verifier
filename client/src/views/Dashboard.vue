<script setup lang="ts">
import { computed, onMounted } from 'vue';
import router from '../router';
import { RouterLink, RouterView } from 'vue-router';

function logout() {
	localStorage.clear();
	router.replace("/")
}

const curRoute = computed(() => {
	const routesParts = router.currentRoute.value.path.split("/")
	return routesParts[routesParts.length - 1];
})

onMounted(() => {
	router.replace("/dashboard/files");
})

</script>

<template>
	<div id="dashboard">
		<div class="header">
			<h1>Email Verifier</h1>
			<button class="btn btn-nobg logout-btn" @click="logout">Logout</button>
		</div>
		<div class="dashboard-body">
			<nav>
				<RouterLink to="/dashboard/files" :replace="true" :class="[curRoute === 'files' ? 'cur-route' : '']">
					My Files
				</RouterLink>
				<RouterLink to="/dashboard/quick-verify" :replace="true" :class="[curRoute === 'quick-verify' ? 'cur-route' : '']">
					Quick Verify
				</RouterLink>
			</nav>
			<router-view />
		</div>
	</div>
</template>

<style scoped>
@import "@css/common.css";

#dashboard{
	width: 100%;
	height: 100%;
	background-color: rgb(30, 30, 30);
	color: white;
}

h1{
	font-family: monospace;
}

.header{
	display: flex;
	flex-direction: row;
	justify-content: space-between;
	align-items: center;
	width: 100%;
	padding: 1rem;
	border-bottom: 2px solid #8f07f7;
}

.dashboard-body{
	background-color: rgb(35, 35, 35);
	width: 100%;
	height: 100%;
}

nav{
	display: flex;
	flex-direction: row;
	align-items: center;
	justify-content: center;
	width: 100%;
	padding: 1rem;
}

nav > *{
	text-decoration: none;
	color: white;
	margin-right: 2rem;
	font-size: 1.35rem;
	position: relative;
	width: 8rem;
	text-align: center;
}

nav > *:hover{
	color: rgb(200, 200, 255);
}

.cur-route::after{
	display: block;
	content: "";
	background-color: white;
	position: absolute;
	width: 100%;
	height: 3px;
	bottom: -0.3rem;
}

</style>
