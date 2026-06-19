package mentor

import "strings"

// The template backend: deterministic, instant, offline. Serves ALL tier-1
// nudges (by design — nudges never reach an AI) and stands in for every
// other kind when no AI backend is configured or one misbehaves.

// nudges: six escalating pointers per bench category. Attempt picks the
// row (clamped), so repeated nudges on one problem get progressively warmer
// without ever giving the approach away.
var nudges = map[string][6]string{
	"Arrays & Hashing": {
		"What information have you already walked past that you're throwing away? A container that remembers could help.",
		"Think about what a hash map or set gives you: one-pass lookups of things you've already seen.",
		"Try storing each element (or its count/index) in a map as you scan, and ask the map a question at every step.",
		"If you had to solve this with a second pass, what would you wish you had written down during the first?",
		"What is the relationship between each element and something else in the array — and can a map represent that relationship instantly?",
		"When your loop is at index i, what is the one fact about previous elements that makes or breaks the answer at i?",
	},
	"Two Pointers & Sliding Window": {
		"Do you really need to re-examine the whole range every time, or does the answer move in one direction?",
		"Picture two markers on the data — when can you safely move the left one forward without missing an answer?",
		"Maintain a window (or pair of pointers) and a running invariant; move one side per step and update the invariant instead of recomputing.",
		"What condition must the window satisfy right now, and which end do you tighten when it's violated?",
		"Are you accidentally shrinking the window past a valid answer — what is the earliest the left pointer can safely advance?",
		"What one piece of state captures everything about the window's current health, and how does it change when you shift a boundary?",
	},
	"Stack": {
		"The most recent unfinished thing matters most here. What structure hands you the latest item first?",
		"When you meet a closer/smaller/blocker, what was the last thing still waiting? A stack remembers exactly that.",
		"Push items as you scan; when the current item resolves the one on top, pop and combine. The stack holds everything still unresolved.",
		"What does it mean for something to be 'unresolved' in this problem, and what event closes it?",
		"If you removed the stack and used a plain variable instead, what case would you get wrong first?",
		"After processing every element, what should remain on the stack, and is an empty stack a success or a failure here?",
	},
	"Binary Search": {
		"The data (or the answer space) is ordered. Can one comparison rule out half of everything?",
		"You don't have to search positions — you can binary-search the ANSWER itself if you can test 'is X feasible?' quickly.",
		"Set low/high bounds, probe the middle, and write the condition that decides which half survives. Be careful which side keeps mid.",
		"What is the predicate that flips from false to true exactly once across the sorted range — can you write it in plain English first?",
		"Are your bounds inclusive or exclusive, and does your termination condition match? A one-off boundary kills a correct binary search.",
		"After the loop ends, does your pointer land on the answer, one past it, or sometimes neither — and does your code account for that?",
	},
	"Linked List": {
		"You can't index backwards — but what can two walkers moving at different speeds or offsets tell you?",
		"A slow and a fast pointer (or one pointer given a head start) often replaces the need to know the length.",
		"Track prev/curr/next explicitly and rewire one link per step; a dummy head node removes the special cases.",
		"Before rewriting any pointers, can you state exactly which three nodes are involved in each rewire step?",
		"What happens to your algorithm if the list has exactly one node or exactly two — does your pointer logic still hold?",
		"Are you saving the next pointer before you overwrite curr.next, or could you lose part of the list mid-operation?",
	},
	"Trees & Graphs": {
		"Solve it for one node assuming the children already gave you their answers. What would you ask them for?",
		"Think recursively: what value must each subtree report upward so the parent can decide?",
		"Write the base case (empty node), recurse on children, combine their results — and decide what global state, if any, you carry along.",
		"Is the answer entirely local to each node, or does it span across the root — and does that change what you return versus what you accumulate?",
		"For a graph, have you accounted for revisiting nodes — what prevents your traversal from looping forever?",
		"What information do you need to carry downward from parent to child, versus what bubbles back upward from child to parent?",
	},
	"Dynamic Programming": {
		"Smaller versions of this exact question hide inside it. What's the smallest one you can answer instantly?",
		"Define 'the best answer ending at / using position i' — can you build i's answer from earlier ones?",
		"Write the recurrence dp[i] in words first, set the base case, and fill the table in an order where dependencies are ready.",
		"Are you solving the same subproblem more than once in your current approach — what would it cost to just remember it?",
		"What are the dimensions of your state — is one variable enough, or does the problem have a second axis that changes the answer?",
		"Could you recover the actual solution (not just the optimal value) from your table, and does the way you fill it make that possible?",
	},
	"Backtracking": {
		"You're choosing, exploring, and un-choosing. What does one 'choice' look like here?",
		"Build the answer one decision at a time; when a partial answer can't possibly work, abandon it early.",
		"Recurse with the current partial solution, try each legal option, and undo it after the recursive call returns.",
		"At what point can you tell a branch is hopeless before reaching a leaf — what is the earliest you can prune it?",
		"Are you generating duplicates? Think about whether the order you try choices in can eliminate symmetric paths.",
		"What constraint separates a 'legal option' from an illegal one at each step, and how cheaply can you check it?",
	},
	"Heap / Priority Queue": {
		"You repeatedly need the smallest/largest of a changing collection. Which structure serves that in log time?",
		"A heap keeps the extreme element on top while you push and pop — you never need the rest sorted.",
		"Push candidates into a heap as you scan; pop when it grows past k (or when the top is no longer useful).",
		"When you pop an element, is it because you're done with it, or because it's stale — and does your code distinguish those two reasons?",
		"Should this be a min-heap or a max-heap? Double-check by tracing what sits on top after two or three pushes.",
		"Is every candidate going into the heap, or only ones that could plausibly be in the answer — limiting what you push keeps the heap small.",
	},
	"Intervals": {
		"Chaos becomes order if you line the intervals up first. By which endpoint?",
		"Sort by start (or end) and sweep: only the previous interval's edge matters at each step.",
		"After sorting, compare each interval with the last kept one — overlap means merge/count, otherwise move on.",
		"What is the exact condition for two intervals to overlap, and does it change if you sort by start versus by end?",
		"When you merge two intervals, what are the correct new start and end — is the new end always the larger of the two?",
		"Are there edge cases where intervals share exactly one endpoint — does your overlap test handle touching intervals correctly?",
	},
	"Greedy": {
		"Is there a locally obvious best move that can never hurt you later? Try to argue why.",
		"Sort or scan so the safest choice comes first, then commit to it without looking back.",
		"Take the best immediate option at each step and maintain just enough state to know what 'best' means next.",
		"Can you construct a counterexample to your greedy rule? If you can, the rule needs tightening.",
		"What ordering or sorting makes the greedy choice obvious — and does that ordering actually correspond to what you want to minimize or maximize?",
		"After committing to a greedy choice, what is the one piece of state you must update so the next step stays consistent?",
	},
	"Math & Bit": {
		"There's a property of the numbers themselves (parity, digits, bits) doing the heavy lifting here.",
		"Play with small cases and watch the pattern — especially what XOR, shifts, or modulo do to it.",
		"Express the answer as an identity or bit trick first; the code is usually three lines once the math is right.",
		"What happens to the bits when you add one or subtract one — does that suggest a connection between n and a neighboring value?",
		"If you write out a few values in binary, do you see a pattern in when the property holds and when it doesn't?",
		"Is there a mathematical invariant that is preserved by every valid operation here — and what does violating it imply?",
	},
	"Strings": {
		"Characters repeat and cluster — counting or remembering positions usually beats re-scanning.",
		"A frequency map or last-seen-index map turns nested scans into one pass.",
		"Scan once, maintain counts/indices in a map, and update the answer at each character instead of afterwards.",
		"What is the 'window property' you need to maintain, and does a character entering or leaving the window change a count by exactly one?",
		"Are you handling duplicate characters inside a window correctly — what does the count reaching zero actually mean?",
		"Could two different substrings give the same frequency map? If so, what does that tell you about the problem's structure?",
	},
	"Tries": {
		"Many words share prefixes. What shape stores shared beginnings exactly once?",
		"A tree where each edge is one character makes prefix questions O(length), not O(words).",
		"Build nodes with a children map and an end-of-word flag; walking the trie IS the algorithm.",
		"What should happen when you reach a node that has no child for the next character — is that a definite miss, or does the problem allow wildcards?",
		"What information beyond end-of-word might each node need to store to answer this particular query efficiently?",
		"Trace an insert and then a search for the same word — does your node structure round-trip correctly, including the terminal flag?",
	},
	"Advanced Graphs": {
		"This is a graph problem in costume — what are the nodes, and what connects them?",
		"Think about which classic applies: shortest path (priority queue), ordering (topological), connectivity (union-find).",
		"Build the adjacency structure explicitly first, then run the classic algorithm; most of the difficulty is the modeling.",
		"Is the graph directed or undirected, and does that change which algorithm is valid — for instance, can you still use union-find if edges have direction?",
		"What is the weight (or cost) of an edge here, and is it always one, or does it vary — because that choice determines which shortest-path algorithm to reach for?",
		"Have you considered whether the graph can have cycles, and if so, does the algorithm you chose handle them or silently give the wrong answer?",
	},
}

var genericNudges = [6]string{
	"Re-read the problem and restate it in your own words — what is ACTUALLY being asked?",
	"Work one small example by hand and watch what you do; your hand-steps are the algorithm.",
	"Name the pattern: scanning? searching? building up answers from smaller ones? That name picks the tool.",
	"What is the brute-force solution you already understand, and where exactly is it doing redundant work?",
	"Which input would break your current approach first — try the empty case, a single element, and all-identical elements.",
	"What invariant must hold after every step of your loop, and can you check it holds on your small example before running the grader?",
}

// Nudge serves a tier-1 nudge — always template, always free of AI.
func Nudge(category string, attempt int) string {
	row := attempt
	if row < 0 {
		row = 0
	}
	if row > 5 {
		row = 5
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
	case KindStandup:
		return standupTemplate(req)
	case KindDiscuss:
		return discussTemplate
	default:
		return Nudge(req.Category, req.Attempt)
	}
}

// standupTemplate is the offline standup: the pre-computed status lines, as-is.
func standupTemplate(req Request) string {
	if len(req.Status) == 0 {
		return "Nothing in flight — let's pick up the next tickets and keep the board moving."
	}
	return strings.Join(req.Status, "\n")
}

// discussTemplate is a delegated teammate's offline plan + estimate.
const discussTemplate = "Here's my plan: I'll nail the acceptance criteria first, make the smallest change that satisfies them, and add a test before I call it done. I'll keep you posted at standup."
