package mentor

// The template backend: deterministic, instant, offline. Serves ALL tier-1
// nudges (by design — nudges never reach an AI) and stands in for every
// other kind when no AI backend is configured or one misbehaves.

// nudges: three escalating pointers per bench category. Attempt picks the
// row (clamped), so repeated nudges on one problem get progressively warmer
// without ever giving the approach away.
var nudges = map[string][3]string{
	"Arrays & Hashing": {
		"What information have you already walked past that you're throwing away? A container that remembers could help.",
		"Think about what a hash map or set gives you: one-pass lookups of things you've already seen.",
		"Try storing each element (or its count/index) in a map as you scan, and ask the map a question at every step.",
	},
	"Two Pointers & Sliding Window": {
		"Do you really need to re-examine the whole range every time, or does the answer move in one direction?",
		"Picture two markers on the data — when can you safely move the left one forward without missing an answer?",
		"Maintain a window (or pair of pointers) and a running invariant; move one side per step and update the invariant instead of recomputing.",
	},
	"Stack": {
		"The most recent unfinished thing matters most here. What structure hands you the latest item first?",
		"When you meet a closer/smaller/blocker, what was the last thing still waiting? A stack remembers exactly that.",
		"Push items as you scan; when the current item resolves the one on top, pop and combine. The stack holds everything still unresolved.",
	},
	"Binary Search": {
		"The data (or the answer space) is ordered. Can one comparison rule out half of everything?",
		"You don't have to search positions — you can binary-search the ANSWER itself if you can test 'is X feasible?' quickly.",
		"Set low/high bounds, probe the middle, and write the condition that decides which half survives. Be careful which side keeps mid.",
	},
	"Linked List": {
		"You can't index backwards — but what can two walkers moving at different speeds or offsets tell you?",
		"A slow and a fast pointer (or one pointer given a head start) often replaces the need to know the length.",
		"Track prev/curr/next explicitly and rewire one link per step; a dummy head node removes the special cases.",
	},
	"Trees & Graphs": {
		"Solve it for one node assuming the children already gave you their answers. What would you ask them for?",
		"Think recursively: what value must each subtree report upward so the parent can decide?",
		"Write the base case (empty node), recurse on children, combine their results — and decide what global state, if any, you carry along.",
	},
	"Dynamic Programming": {
		"Smaller versions of this exact question hide inside it. What's the smallest one you can answer instantly?",
		"Define 'the best answer ending at / using position i' — can you build i's answer from earlier ones?",
		"Write the recurrence dp[i] in words first, set the base case, and fill the table in an order where dependencies are ready.",
	},
	"Backtracking": {
		"You're choosing, exploring, and un-choosing. What does one 'choice' look like here?",
		"Build the answer one decision at a time; when a partial answer can't possibly work, abandon it early.",
		"Recurse with the current partial solution, try each legal option, and undo it after the recursive call returns.",
	},
	"Heap / Priority Queue": {
		"You repeatedly need the smallest/largest of a changing collection. Which structure serves that in log time?",
		"A heap keeps the extreme element on top while you push and pop — you never need the rest sorted.",
		"Push candidates into a heap as you scan; pop when it grows past k (or when the top is no longer useful).",
	},
	"Intervals": {
		"Chaos becomes order if you line the intervals up first. By which endpoint?",
		"Sort by start (or end) and sweep: only the previous interval's edge matters at each step.",
		"After sorting, compare each interval with the last kept one — overlap means merge/count, otherwise move on.",
	},
	"Greedy": {
		"Is there a locally obvious best move that can never hurt you later? Try to argue why.",
		"Sort or scan so the safest choice comes first, then commit to it without looking back.",
		"Take the best immediate option at each step and maintain just enough state to know what 'best' means next.",
	},
	"Math & Bit": {
		"There's a property of the numbers themselves (parity, digits, bits) doing the heavy lifting here.",
		"Play with small cases and watch the pattern — especially what XOR, shifts, or modulo do to it.",
		"Express the answer as an identity or bit trick first; the code is usually three lines once the math is right.",
	},
	"Strings": {
		"Characters repeat and cluster — counting or remembering positions usually beats re-scanning.",
		"A frequency map or last-seen-index map turns nested scans into one pass.",
		"Scan once, maintain counts/indices in a map, and update the answer at each character instead of afterwards.",
	},
	"Tries": {
		"Many words share prefixes. What shape stores shared beginnings exactly once?",
		"A tree where each edge is one character makes prefix questions O(length), not O(words).",
		"Build nodes with a children map and an end-of-word flag; walking the trie IS the algorithm.",
	},
	"Advanced Graphs": {
		"This is a graph problem in costume — what are the nodes, and what connects them?",
		"Think about which classic applies: shortest path (priority queue), ordering (topological), connectivity (union-find).",
		"Build the adjacency structure explicitly first, then run the classic algorithm; most of the difficulty is the modeling.",
	},
}

var genericNudges = [3]string{
	"Re-read the problem and restate it in your own words — what is ACTUALLY being asked?",
	"Work one small example by hand and watch what you do; your hand-steps are the algorithm.",
	"Name the pattern: scanning? searching? building up answers from smaller ones? That name picks the tool.",
}

// Nudge serves a tier-1 nudge — always template, always free of AI.
func Nudge(category string, attempt int) string {
	row := attempt
	if row < 0 {
		row = 0
	}
	if row > 2 {
		row = 2
	}
	if n, ok := nudges[category]; ok {
		return n[row]
	}
	return genericNudges[row]
}

// strategyTemplates: the offline tier-2 answer, per category.
var strategyTemplates = map[string]string{
	"Arrays & Hashing":              "Scan once while a hash map/set remembers what you've seen (values, counts, or indices). At each element, ask the map whether the answer just became completable. One pass, O(n) — the map replaces the inner loop.",
	"Two Pointers & Sliding Window": "Keep a window [left, right] and an invariant (sum, counts, uniqueness). Grow right; when the invariant breaks, shrink left until it holds. Every element enters and leaves the window once.",
	"Stack":                         "Scan and push unresolved items. When the current item resolves the top (matches it, is warmer/taller/closer), pop and combine. The stack always holds exactly the still-open items in order.",
	"Binary Search":                 "Identify the sorted axis — the array itself or the space of possible answers. Write a predicate that's false…false,true…true along that axis, then binary-search the boundary.",
	"Linked List":                   "Use pointer choreography instead of indexing: a dummy head for edge cases, slow/fast pointers for middles and cycles, and an offset pointer for nth-from-end. Rewire one link per step.",
	"Trees & Graphs":                "Recurse: define what each subtree must report (depth, validity range, best path through), combine the children's reports at the node, and let the base case be the empty tree. For graphs, that's DFS/BFS with a visited set.",
	"Dynamic Programming":           "Define the state in words — 'best answer for prefix i (and choice j)'. Write the recurrence relating it to smaller states, set the base case, and fill in dependency order. Optimize space last.",
	"Backtracking":                  "Frame it as a decision tree: at each step enumerate legal choices, apply one, recurse, undo. Prune branches that can't reach a valid answer. Collect results at the leaves.",
	"Heap / Priority Queue":         "Keep a heap of the candidates that still matter (often capped at size k). Push as you scan, pop when it outgrows k or the top is stale — the top is your running answer.",
	"Intervals":                     "Sort by start. Sweep left to right comparing only with the last kept interval: overlap → merge or count a conflict; no overlap → commit and move on.",
	"Greedy":                        "Sort so the safest choice appears first, then commit step by step, keeping just enough state to evaluate the next choice. Justify: exchanging any choice for the greedy one never makes things worse.",
	"Math & Bit":                    "Look for the invariant in the numbers: XOR cancels pairs, n&(n-1) drops the lowest set bit, digit loops run in O(log n). Verify on two or three small cases before coding.",
	"Strings":                       "One pass with a frequency or last-position map. Windows handle 'longest substring with property X'; counting handles anagram-style equivalence.",
	"Tries":                         "Build a character tree: nodes with a children map and a word-end flag. Insert and query are both 'walk the word one char at a time'; prefix questions fall out for free.",
	"Advanced Graphs":               "Model nodes and edges explicitly first. Then pick the classic: Dijkstra/priority queue for weighted shortest paths, topological order for dependencies, union-find for connectivity/cycles.",
}

const genericStrategy = "Restate the problem, pick the dominant pattern (scan with memory, two pointers, recursion over structure, search over answers), and write the loop invariant in a comment before coding it."

// strategyTemplate is the offline tier-2 fallback.
func strategyTemplate(category string) string {
	if s, ok := strategyTemplates[category]; ok {
		return s
	}
	return genericStrategy
}

// walkthroughTemplate is the offline tier-3 fallback: the strategy plus a
// generic decomposition (templates can't see the player's code, so the AI
// version is strictly better — this keeps the tier functional offline).
func walkthroughTemplate(category string) string {
	return strategyTemplate(category) + "\n\nWork it as steps:\n" +
		"1. Write the brute force in comments — name the redundant work it repeats.\n" +
		"2. Pick the structure/technique above that eliminates exactly that redundancy.\n" +
		"3. Code the main loop around its invariant; handle the empty/size-1 cases first.\n" +
		"4. Trace your code on the example in the prompt BEFORE running the grader.\n" +
		"5. If a test still fails, diff your trace against the expected output at the first divergence."
}

const followupTemplate = "In one sentence: why does your approach not miss any valid answer, and what input shape would stress it most?"

const reviewTemplate = "Banked. Now re-read your solution once as a stranger: is every name honest, and is there one branch you could delete by handling the general case better?"

// templateAnswer renders the offline answer for any kind.
func templateAnswer(req Request) string {
	switch req.Kind {
	case KindStrategy:
		return strategyTemplate(req.Category)
	case KindWalkthrough:
		return walkthroughTemplate(req.Category)
	case KindFollowup:
		return followupTemplate
	case KindReview:
		return reviewTemplate
	default:
		return Nudge(req.Category, req.Attempt)
	}
}
