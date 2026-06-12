export namespace guiapi {
	
	export class AdvExerciseView {
	    index: number;
	    kind: string;
	    prompt: string;
	    brokenCode: string;
	    fixedCode: string;
	    bug: string;
	    check: string;
	    gradeable: boolean;
	
	    static createFrom(source: any = {}) {
	        return new AdvExerciseView(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.index = source["index"];
	        this.kind = source["kind"];
	        this.prompt = source["prompt"];
	        this.brokenCode = source["brokenCode"];
	        this.fixedCode = source["fixedCode"];
	        this.bug = source["bug"];
	        this.check = source["check"];
	        this.gradeable = source["gradeable"];
	    }
	}
	export class PrimerOpView {
	    label: string;
	    code: string;
	
	    static createFrom(source: any = {}) {
	        return new PrimerOpView(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.label = source["label"];
	        this.code = source["code"];
	    }
	}
	export class PrimerSectionView {
	    title: string;
	    ops: PrimerOpView[];
	
	    static createFrom(source: any = {}) {
	        return new PrimerSectionView(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.ops = this.convertValues(source["ops"], PrimerOpView);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class AdvTopicDetail {
	    found: boolean;
	    index: number;
	    lang: string;
	    group: string;
	    title: string;
	    tag: string;
	    summary: string;
	    sections: PrimerSectionView[];
	    exercises: AdvExerciseView[];
	
	    static createFrom(source: any = {}) {
	        return new AdvTopicDetail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.found = source["found"];
	        this.index = source["index"];
	        this.lang = source["lang"];
	        this.group = source["group"];
	        this.title = source["title"];
	        this.tag = source["tag"];
	        this.summary = source["summary"];
	        this.sections = this.convertValues(source["sections"], PrimerSectionView);
	        this.exercises = this.convertValues(source["exercises"], AdvExerciseView);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class AdvTopicSummary {
	    index: number;
	    group: string;
	    title: string;
	    tag: string;
	    exercises: number;
	    gradeable: number;
	
	    static createFrom(source: any = {}) {
	        return new AdvTopicSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.index = source["index"];
	        this.group = source["group"];
	        this.title = source["title"];
	        this.tag = source["tag"];
	        this.exercises = source["exercises"];
	        this.gradeable = source["gradeable"];
	    }
	}
	export class CaseResult {
	    name: string;
	    passed: boolean;
	    got: string;
	    expected: string;
	    err: string;
	
	    static createFrom(source: any = {}) {
	        return new CaseResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.passed = source["passed"];
	        this.got = source["got"];
	        this.expected = source["expected"];
	        this.err = source["err"];
	    }
	}
	export class DevLitStep {
	    done: boolean;
	    index: number;
	    total: number;
	    category: string;
	    title: string;
	    prompt: string;
	    passed: number;
	
	    static createFrom(source: any = {}) {
	        return new DevLitStep(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.done = source["done"];
	        this.index = source["index"];
	        this.total = source["total"];
	        this.category = source["category"];
	        this.title = source["title"];
	        this.prompt = source["prompt"];
	        this.passed = source["passed"];
	    }
	}
	export class DevLitOutcome {
	    passed: boolean;
	    hint: string;
	    success: string;
	    next: DevLitStep;
	
	    static createFrom(source: any = {}) {
	        return new DevLitOutcome(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.passed = source["passed"];
	        this.hint = source["hint"];
	        this.success = source["success"];
	        this.next = this.convertValues(source["next"], DevLitStep);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class OrientationStep {
	    done: boolean;
	    index: number;
	    total: number;
	    kind: string;
	    slot: string;
	    measures: string;
	    prompt: string;
	    lang: string;
	    funcName: string;
	    starter: string;
	    choices: string[];
	    placement: string;
	    score: number;
	    level: string;
	    codingOK: number;
	    codingTotal: number;
	    machineOK: number;
	    machineTotal: number;
	    specOK: number;
	    specTotal: number;
	
	    static createFrom(source: any = {}) {
	        return new OrientationStep(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.done = source["done"];
	        this.index = source["index"];
	        this.total = source["total"];
	        this.kind = source["kind"];
	        this.slot = source["slot"];
	        this.measures = source["measures"];
	        this.prompt = source["prompt"];
	        this.lang = source["lang"];
	        this.funcName = source["funcName"];
	        this.starter = source["starter"];
	        this.choices = source["choices"];
	        this.placement = source["placement"];
	        this.score = source["score"];
	        this.level = source["level"];
	        this.codingOK = source["codingOK"];
	        this.codingTotal = source["codingTotal"];
	        this.machineOK = source["machineOK"];
	        this.machineTotal = source["machineTotal"];
	        this.specOK = source["specOK"];
	        this.specTotal = source["specTotal"];
	    }
	}
	export class GradeResult {
	    passed: boolean;
	    casesTotal: number;
	    casesFailed: number;
	    err: string;
	    results: CaseResult[];
	    banked: boolean;
	    newlyBanked: boolean;
	    saveErr: string;
	    tokensAwarded: number;
	    writeupPending: boolean;
	
	    static createFrom(source: any = {}) {
	        return new GradeResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.passed = source["passed"];
	        this.casesTotal = source["casesTotal"];
	        this.casesFailed = source["casesFailed"];
	        this.err = source["err"];
	        this.results = this.convertValues(source["results"], CaseResult);
	        this.banked = source["banked"];
	        this.newlyBanked = source["newlyBanked"];
	        this.saveErr = source["saveErr"];
	        this.tokensAwarded = source["tokensAwarded"];
	        this.writeupPending = source["writeupPending"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class DiagOutcome {
	    passed: boolean;
	    feedback: string;
	    verdict?: GradeResult;
	    next: OrientationStep;
	
	    static createFrom(source: any = {}) {
	        return new DiagOutcome(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.passed = source["passed"];
	        this.feedback = source["feedback"];
	        this.verdict = this.convertValues(source["verdict"], GradeResult);
	        this.next = this.convertValues(source["next"], OrientationStep);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class GateCategoryView {
	    category: string;
	    done: number;
	    required: number;
	    available: number;
	
	    static createFrom(source: any = {}) {
	        return new GateCategoryView(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.category = source["category"];
	        this.done = source["done"];
	        this.required = source["required"];
	        this.available = source["available"];
	    }
	}
	export class GateItemView {
	    slug: string;
	    title: string;
	    done: boolean;
	
	    static createFrom(source: any = {}) {
	        return new GateItemView(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.slug = source["slug"];
	        this.title = source["title"];
	        this.done = source["done"];
	    }
	}
	export class GateView {
	    full: number;
	    provisional: number;
	    target: number;
	    categories: GateCategoryView[];
	    mandatory: GateItemView[];
	    countMet: boolean;
	    catsMet: boolean;
	    mandatoryOk: boolean;
	    met: boolean;
	
	    static createFrom(source: any = {}) {
	        return new GateView(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.full = source["full"];
	        this.provisional = source["provisional"];
	        this.target = source["target"];
	        this.categories = this.convertValues(source["categories"], GateCategoryView);
	        this.mandatory = this.convertValues(source["mandatory"], GateItemView);
	        this.countMet = source["countMet"];
	        this.catsMet = source["catsMet"];
	        this.mandatoryOk = source["mandatoryOk"];
	        this.met = source["met"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class WalletView {
	    tokens: number;
	    nudgeCharges: number;
	    nudgeMax: number;
	    nextRechargeSec: number;
	
	    static createFrom(source: any = {}) {
	        return new WalletView(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.tokens = source["tokens"];
	        this.nudgeCharges = source["nudgeCharges"];
	        this.nudgeMax = source["nudgeMax"];
	        this.nextRechargeSec = source["nextRechargeSec"];
	    }
	}
	export class HintResult {
	    text: string;
	    source: string;
	    tier: number;
	    pity: boolean;
	    refunded: boolean;
	    wallet: WalletView;
	    err: string;
	
	    static createFrom(source: any = {}) {
	        return new HintResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.text = source["text"];
	        this.source = source["source"];
	        this.tier = source["tier"];
	        this.pity = source["pity"];
	        this.refunded = source["refunded"];
	        this.wallet = this.convertValues(source["wallet"], WalletView);
	        this.err = source["err"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class InstallGuideView {
	    found: boolean;
	    lang: string;
	    label: string;
	    notes: string;
	    os: string;
	    link: string;
	    steps: string[];
	    verify: string;
	
	    static createFrom(source: any = {}) {
	        return new InstallGuideView(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.found = source["found"];
	        this.lang = source["lang"];
	        this.label = source["label"];
	        this.notes = source["notes"];
	        this.os = source["os"];
	        this.link = source["link"];
	        this.steps = source["steps"];
	        this.verify = source["verify"];
	    }
	}
	export class LangStatus {
	    lang: string;
	    status: string;
	    verified: boolean;
	    version: string;
	    reason: string;
	
	    static createFrom(source: any = {}) {
	        return new LangStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.lang = source["lang"];
	        this.status = source["status"];
	        this.verified = source["verified"];
	        this.version = source["version"];
	        this.reason = source["reason"];
	    }
	}
	export class LessonStageView {
	    kind: string;
	    title: string;
	    body: string;
	    hasTask: boolean;
	    prompt: string;
	    funcName: string;
	    starter: string;
	
	    static createFrom(source: any = {}) {
	        return new LessonStageView(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.kind = source["kind"];
	        this.title = source["title"];
	        this.body = source["body"];
	        this.hasTask = source["hasTask"];
	        this.prompt = source["prompt"];
	        this.funcName = source["funcName"];
	        this.starter = source["starter"];
	    }
	}
	export class LessonView {
	    id: string;
	    title: string;
	    index: number;
	    total: number;
	    found: boolean;
	    stages: LessonStageView[];
	
	    static createFrom(source: any = {}) {
	        return new LessonView(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.index = source["index"];
	        this.total = source["total"];
	        this.found = source["found"];
	        this.stages = this.convertValues(source["stages"], LessonStageView);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	export class PrimerView {
	    found: boolean;
	    category: string;
	    title: string;
	    summary: string;
	    sections: PrimerSectionView[];
	    example: string;
	
	    static createFrom(source: any = {}) {
	        return new PrimerView(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.found = source["found"];
	        this.category = source["category"];
	        this.title = source["title"];
	        this.summary = source["summary"];
	        this.sections = this.convertValues(source["sections"], PrimerSectionView);
	        this.example = source["example"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ProblemDetail {
	    id: string;
	    title: string;
	    difficulty: string;
	    category: string;
	    prompt: string;
	    funcName: string;
	    lang: string;
	    starter: string;
	    found: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ProblemDetail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.difficulty = source["difficulty"];
	        this.category = source["category"];
	        this.prompt = source["prompt"];
	        this.funcName = source["funcName"];
	        this.lang = source["lang"];
	        this.starter = source["starter"];
	        this.found = source["found"];
	    }
	}
	export class ProblemSummary {
	    id: string;
	    title: string;
	    difficulty: string;
	    category: string;
	    lists: string[];
	    solved: boolean;
	    writeup: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ProblemSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.difficulty = source["difficulty"];
	        this.category = source["category"];
	        this.lists = source["lists"];
	        this.solved = source["solved"];
	        this.writeup = source["writeup"];
	    }
	}
	export class ProfileView {
	    lang: string;
	    stage: string;
	    placement: string;
	    level: string;
	    banked: number;
	    updatedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new ProfileView(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.lang = source["lang"];
	        this.stage = source["stage"];
	        this.placement = source["placement"];
	        this.level = source["level"];
	        this.banked = source["banked"];
	        this.updatedAt = source["updatedAt"];
	    }
	}
	export class Progress {
	    lang: string;
	    banked: number;
	    cats: number;
	    hard: number;
	    bankTarget: number;
	    catTarget: number;
	    hardTarget: number;
	    step0Met: boolean;
	    problemSolving: number;
	    langProf: number;
	    track: string;
	    placement: string;
	    level: string;
	    totalProblems: number;
	
	    static createFrom(source: any = {}) {
	        return new Progress(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.lang = source["lang"];
	        this.banked = source["banked"];
	        this.cats = source["cats"];
	        this.hard = source["hard"];
	        this.bankTarget = source["bankTarget"];
	        this.catTarget = source["catTarget"];
	        this.hardTarget = source["hardTarget"];
	        this.step0Met = source["step0Met"];
	        this.problemSolving = source["problemSolving"];
	        this.langProf = source["langProf"];
	        this.track = source["track"];
	        this.placement = source["placement"];
	        this.level = source["level"];
	        this.totalProblems = source["totalProblems"];
	    }
	}
	export class TutorialPos {
	    lesson: number;
	    stage: number;
	    done: boolean;
	
	    static createFrom(source: any = {}) {
	        return new TutorialPos(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.lesson = source["lesson"];
	        this.stage = source["stage"];
	        this.done = source["done"];
	    }
	}
	
	export class WriteupResult {
	    accepted: boolean;
	    mcqCorrect: boolean;
	    tokensAwarded: number;
	    wallet: WalletView;
	    followup: string;
	    err: string;
	
	    static createFrom(source: any = {}) {
	        return new WriteupResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.accepted = source["accepted"];
	        this.mcqCorrect = source["mcqCorrect"];
	        this.tokensAwarded = source["tokensAwarded"];
	        this.wallet = this.convertValues(source["wallet"], WalletView);
	        this.followup = source["followup"];
	        this.err = source["err"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class WriteupView {
	    problemId: string;
	    title: string;
	    solved: boolean;
	    done: boolean;
	    question: string;
	    options: string[];
	    hasMcq: boolean;
	    minLen: number;
	
	    static createFrom(source: any = {}) {
	        return new WriteupView(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.problemId = source["problemId"];
	        this.title = source["title"];
	        this.solved = source["solved"];
	        this.done = source["done"];
	        this.question = source["question"];
	        this.options = source["options"];
	        this.hasMcq = source["hasMcq"];
	        this.minLen = source["minLen"];
	    }
	}

}

export namespace main {
	
	export class MentorStatusView {
	    id: string;
	    name: string;
	    present: boolean;
	    info: string;
	    selected: boolean;
	    probed: boolean;
	    probeOk: boolean;
	    probeErr: string;
	
	    static createFrom(source: any = {}) {
	        return new MentorStatusView(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.present = source["present"];
	        this.info = source["info"];
	        this.selected = source["selected"];
	        this.probed = source["probed"];
	        this.probeOk = source["probeOk"];
	        this.probeErr = source["probeErr"];
	    }
	}

}

