export function convertToDate(dateStr: string) {
	let d = dateStr.split("/");
	return new Date(d[2] + '/' + d[1] + '/' + d[0]);
}
