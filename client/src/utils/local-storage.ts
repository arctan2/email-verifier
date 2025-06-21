export function getUserId() {
	const userId = localStorage.getItem("userId") || "";
	return userId;
}
