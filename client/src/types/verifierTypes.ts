export interface ProgressData {
	total: number,
	progress: number,
	success: number,
	failed: number,
	retry: number
}

export enum Status {
	NotCreated = "not created",
	Created = "created",
	Running = "running",
	Done = "done"
}

export interface VerifierDetails {
	batchSize: number,
	completedBatches: { [_:number]: ProgressData[] },
	currentBatchNumber: number,
	currentBatchSize: number,
	currentProgressList: ProgressData[]
	delayMs: number
	emailCount: number,
	proxies: string[],
	retryCount: number,
	state: Status
}

