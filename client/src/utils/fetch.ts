export const baseURL = "http://192.168.0.107:8000";
// export const baseURL = "";

export function API_URL(...endpoints: string[]) {
	return `${baseURL}${endpoints.join("/")}`;
}

type ApiResponse = {
	response?: Response,
	err: boolean,
	msg: string
};

export async function fetchGet<T>(endpoint: string, abortController?: AbortController): Promise<ApiResponse & T> {
	try {
		const response = await fetch(API_URL(`/api/web${endpoint}`), {
			method: "GET",
			signal: abortController?.signal,
			headers: {
				"Content-Type": "application/json"
			},
		})

		const data: { err: boolean, msg: string } & T = await response.json();

		return { response, ...data };
	} catch(e: any) {
		let msg = String(e);

		if(e.message === "Failed to fetch") {
			msg = "Please check your internet connection. Or the server is down."
		}

		return { err: true, msg } as ApiResponse & T;
	}
}

export async function fetchPostRaw<T>(endpoint: string, body?: any, abortController?: AbortController): Promise<ApiResponse & T> {
	try {
		const response = await fetch(API_URL(`/api/web${endpoint}`), {
			method: "POST",
			signal: abortController?.signal,
			headers: {
				"Content-Type": "application/json"
			},
			body
		})

		const data: { err: boolean, msg: string } & T = await response.json();

		return { response, ...data };
	} catch(e) {
		return { err: true, msg: String(e) } as ApiResponse & T;
	}
}

export async function fetchPostWithHeader<T>(
	endpoint: string,
	headers: any,
	body?: any,
	abortController?: AbortController
): Promise<ApiResponse & T> {
	try {
		const response = await fetch(API_URL(`/api/web${endpoint}`), {
			method: "POST",
			signal: abortController?.signal,
			headers,
			body
		})

		const data: { err: boolean, msg: string } & T = await response.json();

		return { response, ...data };
	} catch(e) {
		return { err: true, msg: String(e) } as ApiResponse & T;
	}
}

export async function fetchPost<T>(endpoint: string, body?: any, abortController?: AbortController): Promise<ApiResponse & T> {
	if(typeof body !== "string") {
		body = JSON.stringify(body);
	}

	return fetchPostRaw(endpoint, body, abortController);
}

export async function fetchDelete<T>(endpoint: string, abortController?: AbortController): Promise<ApiResponse & T> {
	try {
		const response = await fetch(API_URL(`/api/web${endpoint}`), {
			method: "DELETE",
			signal: abortController?.signal,
			headers: {
				"Content-Type": "application/json"
			},
		})

		const data: { err: boolean, msg: string } & T = await response.json();

		return { response, ...data };
	} catch(e) {
		return { err: true, msg: String(e) } as ApiResponse & T;
	}
}
