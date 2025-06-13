export interface ProgressData {
	total: number,
	progress: number,
	success: number
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
	currentBatchNumber: 0,
	currentBatchSize: 0,
	currentProgressList: ProgressData[]
	delayMs: number
	emailCount: number,
	proxies: string[],
	retryCount: number,
	state: Status
}

