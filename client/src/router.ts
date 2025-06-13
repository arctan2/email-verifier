import { createRouter, createWebHistory, type RouteRecordRaw } from "vue-router";

import RootRoute from "./root-route.vue";

const routes: Array<RouteRecordRaw> = [
	{ path: "/", component: RootRoute },
];

const router = createRouter({
	history: createWebHistory("/"),
	routes
});

export default router;
