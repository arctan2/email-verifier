import type { ProgressData } from "../../types/verifierTypes";

export interface Batch {
	num: number,
	from: number,
	to: number,
	classNames: string[],
	retries: ProgressData[],
}

export function computeBatchesList(
	batchSize: number,
	emailCount: number,
	curBatchNum: number,
	completedBatches: { [_:number]: ProgressData[] }
) {
	const batches: Batch[] = [];
	const totalBatches = Math.ceil(emailCount / batchSize);

	for(let i = 0; i < totalBatches; i++) {
		const batch: Batch = {
			num: i,
			from: (i * batchSize) + 1,
			to: Math.min(emailCount, i * batchSize + batchSize),
			classNames: [],
			retries: [] as ProgressData[]
		}

		const completed = completedBatches[i];

		if(completed) {
			batch.retries = completed;
		} else if(i !== curBatchNum) {
			batch.classNames.push("not-started-batch");
		} else {
			batch.classNames.push("current-batch");
		}

		batches.push(batch);
	}

	return batches;
}

