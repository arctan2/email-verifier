import { ref } from "vue";

export class FilterCycler {
	curIdx = ref(-1);
	classNames: string[];
	cmpValues: any[];

	constructor(classNames: string[], cmpValues: any[]) {
		this.classNames = classNames;
		this.cmpValues = cmpValues;
	}

	get cmpVal() {
		return this.curIdx.value === -1 ? null : this.cmpValues[this.curIdx.value];
	}

	get className() {
		return this.curIdx.value === -1 ? '' : this.classNames[this.curIdx.value];
	}

	reset() {
		this.curIdx.value = -1;
	}

	cycle(emitter: () => void) {
		this.curIdx.value++;
		if(this.curIdx.value >= this.cmpValues.length) {
			this.curIdx.value = -1;
		}
		emitter();
	}
}

const twoStateCmp = [true, false];
const twoStateClassNames = ['green', 'red'];

export const filterers: {[_:string]: FilterCycler} = {
	isValidSyntax: new FilterCycler(twoStateClassNames, twoStateCmp),
	reachable: new FilterCycler([...twoStateClassNames, 'yellow'], ['yes', 'no', 'unknown']),
	isDeliverable: new FilterCycler(twoStateClassNames, twoStateCmp),
	isHostExists: new FilterCycler(twoStateClassNames, twoStateCmp),
	hasMxRecords: new FilterCycler(twoStateClassNames, twoStateCmp),
	isDisposable: new FilterCycler(twoStateClassNames, twoStateCmp),
	isCatchAll: new FilterCycler(twoStateClassNames, twoStateCmp),
	isInboxFull: new FilterCycler(twoStateClassNames, twoStateCmp),
}

export function getToFilter() {
	const toFilter = [];
	for(const k in filterers) {
		if(filterers[k].curIdx.value !== -1) {
			toFilter.push(k);
		}
	}
	return toFilter;
}

export function resetFilters() {
	for(const k in filterers) {
		filterers[k].reset();
	}
}
